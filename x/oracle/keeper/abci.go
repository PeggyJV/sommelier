package keeper

import (
	"fmt"

	"github.com/peggyjv/sommelier/x/oracle/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
)

// BeginBlocker is called at the beginning of every block
func (k Keeper) BeginBlocker(ctx sdk.Context) {
	// if there is not a vote period set, initialize it with current block height
	if !k.HasVotePeriodStart(ctx) {
		k.SetVotePeriodStart(ctx, ctx.BlockHeight())
	}

	// On begin block, if we are tallying, emit the new vote period data
	params := k.GetParamSet(ctx)
	if (ctx.BlockHeight() - k.GetVotePeriodStart(ctx)) >= params.VotePeriod {
		ctx.EventManager().EmitEvent(
			sdk.NewEvent(
				types.EventTypeVotePeriod,
				sdk.NewAttribute(types.AttributeKeyVotePeriodStart, fmt.Sprintf("%d", ctx.BlockHeight())),
				sdk.NewAttribute(types.AttributeKeyVotePeriodEnd, fmt.Sprintf("%d", ctx.BlockHeight()+params.VotePeriod)),
			),
		)
	}

}

// EndBlocker is called at the end of every block
func (k Keeper) EndBlocker(ctx sdk.Context) {
	params := k.GetParamSet(ctx)
	// if the vote period has ended, tally the votes
	periodEnded := (ctx.BlockHeight() - k.GetVotePeriodStart(ctx)) >= params.VotePeriod
	if !periodEnded {
		return
	}

	k.Logger(ctx).Info("tallying oracle votes")

	votedPower := int64(0)

	// TODO: add comments and explain what each map does
	var (
		// validators who submitted their vote
		validatorVotesMap   = make(map[string]bool)
		detailedMap         = make(map[string]map[string]types.OracleData)
		oracleDataByTypeMap = make(map[string][]types.OracleData)
		// average
		averageMap     = make(map[string]types.OracleData)
		rewardEligable = make(map[string]bool)
	)

	// initialize the inner maps for detailedMap
	for _, dt := range params.DataTypes {
		detailedMap[dt] = make(map[string]types.OracleData)
	}

	// iterate over the data votes
	k.IterateOracleDataVotes(ctx, func(validatorAddr sdk.AccAddress, msg *types.MsgOracleDataVote) bool {
		// save a voted array
		validatorVotesMap[validatorAddr.String()] = true

		// remove the miss counter for validators who have voted
		k.DeleteMissCounter(ctx, validatorAddr)

		// find total voting votedPower
		validator := k.stakingKeeper.Validator(ctx, sdk.ValAddress(validatorAddr))
		if validator == nil {
			// validator nor found, continue with next vote
			return false
		}

		votedPower += validator.GetConsensusPower()

		// save the oracle data for later processing
		for _, oracleDataAny := range msg.OracleData {
			oracleData, err := types.UnpackOracleData(oracleDataAny)
			if err != nil {
				// NOTE: this should never panic as the oracle data had already been processed
				panic(err)
			}

			// add oracle data to maps
			detailedMap[oracleData.Type()][validatorAddr.String()] = oracleData
			oracleDataByTypeMap[oracleData.Type()] = append(oracleDataByTypeMap[oracleData.Type()], oracleData)
		}

		// delete the vote as we no longer require it
		k.DeleteOracleDataVote(ctx, validatorAddr)
		return false
	})

	// iterate over the full list of validators to increment miss counters if they didn't vote
	totalPower := int64(0)
	k.stakingKeeper.IterateBondedValidatorsByPower(ctx, func(_ int64, validator stakingtypes.ValidatorI) bool {
		totalPower += validator.GetConsensusPower()
		validatorAddr := sdk.AccAddress(validator.GetOperator())

		if !validatorVotesMap[validatorAddr.String()] {
			k.IncrementMissCounter(ctx, validatorAddr)
		}

		return false
	})

	// if the voted_power/total_power > params.VoteThreshold then we store the averages in the store
	storeAverages := sdk.NewDec(votedPower).Quo(sdk.NewDec(totalPower)).GT(params.VoteThreshold)

	// compute the averages for each type of data tracked by the oracle
	for dataType, oracleDatas := range oracleDataByTypeMap {
		// first, lets delete the old data
		k.DeleteOracleData(ctx, dataType)

		// then we compute the average/target for the given data type
		averageFn, err := types.GetAverageFunction(dataType)
		if err != nil {
			continue
		}

		avg := averageFn(oracleDatas)

		// once we have an "average" we set it in the store
		if storeAverages {
			k.SetOracleData(ctx, avg)
		}

		// store the "average" for scoring validators later
		averageMap[dataType] = avg
	}

	// Compare each validators vote for each data type against the
	// averages to define which are eligable for rewards
	for dataType, vals := range detailedMap {
		for val, data := range vals {
			rewardEligable[val] = false
			if averageMap[dataType].Valid(data) {
				rewardEligable[val] = true
			}
		}
	}

	// slash validators who have missed to many votes
	k.IterateMissCounters(ctx, func(validatorAddr sdk.AccAddress, counter int64) bool {
		missedVotesPerWindow := sdk.NewDec(counter).Quo(sdk.NewDec(params.SlashWindow))
		if params.MinValidPerWindow.GTE(missedVotesPerWindow) {
			// continue with next counter
			return false
		}

		// slash validator below the minimum vote threshold
		validator := k.stakingKeeper.Validator(ctx, sdk.ValAddress(validatorAddr))
		if validator == nil {
			// validator not found, continue with the next counter
			// TODO: ensure correctness
			return false
		}

		consensusAddr, err := validator.GetConsAddr()
		if err != nil {
			// TODO: check why this could error and ensure correctness
			panic(err)
		}

		k.stakingKeeper.Slash(ctx, consensusAddr, ctx.BlockHeight(), validator.GetConsensusPower(), params.SlashFraction)
		return false
	})

	// TODO: reward validators
	// TODO: Setup module account for oracle module

	// Reset state prior to next round
	k.DeleteAllPrevotes(ctx)
	// After the tallying is done, reset the vote period start height and delete all the prevotes
	k.SetVotePeriodStart(ctx, ctx.BlockHeight())
}
