package keeper

import (
	"fmt"

	"github.com/peggyjv/sommelier/x/oracle/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
)

// EndBlocker is called at the end of every block
func (k Keeper) EndBlocker(ctx sdk.Context) {
	params := k.GetParamSet(ctx)
	votePeriodStart, found := k.GetVotePeriodStart(ctx)
	if !found {
		panic("vote period start not set")
	}

	// if the vote period has ended, tally the votes
	periodEnded := (ctx.BlockHeight() - votePeriodStart) >= params.VotePeriod
	if !periodEnded {
		return
	}

	k.Logger(ctx).Info("tallying oracle votes", "height", fmt.Sprintf("%d", ctx.BlockHeight()))

	votedPower := int64(0)

	var (
		// validators who submitted their vote
		validatorVotesMap = make(map[string]bool)
		// oracle data fed by a validator by type: <validator_address><data_type><data_id> --> oracle data
		submittedFeedMap = make(map[string]types.OracleData)
		// all oracle data submitted on the current voting period, by type
		oracleDataByTypeMap = make(map[string][]types.OracleData)

		// all aggregated oracle data for the current voting period
		aggregates = make([]types.OracleData, 0)
	)

	// iterate over the data votes
	k.IterateOracleDataVotes(ctx, func(validatorAddr sdk.AccAddress, vote types.OracleVote) bool {
		// save a voted array
		validatorVotesMap[validatorAddr.String()] = true

		// remove the miss counter in the current voting window for validators who have already voted
		// TODO: Set this to 0 instead?
		k.DeleteMissCounter(ctx, validatorAddr)

		// find total voting votedPower
		validator := k.stakingKeeper.Validator(ctx, sdk.ValAddress(validatorAddr))
		if validator == nil {
			// validator nor found, continue with next vote
			return false
		}

		votedPower += validator.GetConsensusPower()

		// save the oracle data for later processing
		for _, oracleDataAny := range vote.Feed.OracleData {
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
		// TODO: consider keeping the votes for a few voting windows in order to
		// be able to submit evidence of inaccurate data feed
		k.DeleteOracleDataVote(ctx, validatorAddr)
		return false
	})

	// iterate over the full list of validators to increment miss counters if they didn't vote
	totalPower := int64(0)
	k.stakingKeeper.IterateBondedValidatorsByPower(ctx, func(_ int64, validator stakingtypes.ValidatorI) bool {
		totalPower += validator.GetConsensusPower()
		validatorAddr := sdk.AccAddress(validator.GetOperator())

		if !validatorVotesMap[validatorAddr.String()] {
			// TODO: this is wrong because the feeder could submit a single oracle data uniswap pair instead of all the required
			// ones and still be counted as if they voted correctly. Maybe consider adding each uniswap pair id to the params?
			// TODO: we need to define what is the exact data that we want the validators to submit.
			k.IncrementMissCounter(ctx, validatorAddr)
		}

		return false
	})

	// if the voted_power/total_power > params.VoteThreshold then we store the averages in the store
	// TODO: check for sample size too as there could be only a few feeders with a large voting power that submitted the
	// data to the oracle.
	storeAverages := sdk.NewDec(votedPower).Quo(sdk.NewDec(totalPower)).GT(params.VoteThreshold)

	// Now, compute the aggregated data (eg: avg, median, etc.) for each type of data tracked by the oracle

	// NOTE: we iterate over the params data types to avoid using map iteration
	// which is non-deterministic.
	for _, dt := range params.DataTypes {
		oracleData, ok := oracleDataByTypeMap[dt]
		if !ok {
			// no oracle data type was submitted
			continue
		}

		// TODO: delete the oracle old data
		// k.DeleteOracleData(ctx, dataType, dataType, oracleData)

		// aggregate the date using the handler set up during app initialization
		aggregatedData, err := k.oracleHandler(ctx, oracleData)
		if err != nil {
			// TODO: ensure correctness or consider logging instead?
			panic(err)
		}

		// once we have an "average" we set it in the store
		if storeAverages {
			k.SetAggregatedOracleData(ctx, aggregatedData)
		}

		// store the "average" for scoring validators later
		aggregates = append(aggregates, aggregatedData)

	}

	// Compare each validators vote for each data type against the
	// averages to define which are eligable for rewards
	for dataType, vals := range detailedMap {
		for val, data := range vals {
			rewardEligable[val] = false
			// TODO: check if the data provided is within x% of the average? I would assume x would have to be low:
			// eg: 0.1?
			// TODO: define X on the parameters
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
			panic(err)
		}

		// TODO: make sure the validator is not unbonded before slashing

		k.stakingKeeper.Slash(ctx, consensusAddr, ctx.BlockHeight(), validator.GetConsensusPower(), params.SlashFraction)
		return false
	})

	// TODO: reward validators.
	// NOTE: the reward amount should be less than the slashed amount

	// TODO: Setup module account for oracle module.

	// Reset state prior to next round
	k.DeleteAllPrevotes(ctx)

	// After the tallying is done, reset the vote period start
	k.SetVotePeriodStart(ctx, ctx.BlockHeight())

	k.Logger(ctx).Info("vote period set", "height", fmt.Sprintf("%d", votePeriodStart))
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeVotePeriod,
			sdk.NewAttribute(types.AttributeKeyVotePeriodStart, fmt.Sprintf("%d", votePeriodStart)),
			sdk.NewAttribute(types.AttributeKeyVotePeriodEnd, fmt.Sprintf("%d", votePeriodStart+params.VotePeriod)),
		),
	)
}
