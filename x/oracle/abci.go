package oracle

import (
	"fmt"

	"github.com/peggyjv/sommelier/x/oracle/keeper"
	"github.com/peggyjv/sommelier/x/oracle/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// BeginBlocker is called at the beginning of every block
func BeginBlocker(ctx sdk.Context, k keeper.Keeper) {
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
func EndBlocker(ctx sdk.Context, k keeper.Keeper) {
	params := k.GetParamSet(ctx)

	// if the vote period has ended, tally the votes
	if (ctx.BlockHeight() - k.GetVotePeriodStart(ctx)) >= params.VotePeriod {
		k.Logger(ctx).Info("tallying oracle votes")
		voted := []string{}
		votedPower := int64(0)
		detailedMap := make(map[string]map[string]types.OracleData)
		collectionMap := make(map[string][]types.OracleData)
		averageMap := make(map[string]types.OracleData)
		rewardEligable := make(map[string]bool)

		// initialize the inner maps for detailedMap
		for _, dt := range params.DataTypes {
			detailedMap[dt] = make(map[string]types.OracleData)
		}

		// iterate over the data votes
		k.IterateOracleDataVotes(ctx, func(val sdk.AccAddress, msg *types.MsgOracleDataVote) bool {
			// save a voted array
			voted = append(voted, val.String())

			// remove the miss counter for validators who have voted
			k.DeleteMissCounter(ctx, val)

			// find total voting votedPower
			votedPower += k.StakingKeeper.Validator(ctx, sdk.ValAddress(val)).GetConsensusPower()

			// save the oracle data for later processing
			for _, oda := range msg.OracleData {
				od, err := types.UnpackOracleData(oda)
				if err != nil {
					panic(err)
				}
				detailedMap[od.Type()][val.String()] = od
				collectionMap[od.Type()] = append(collectionMap[od.Type()], od)
			}

			// delete the vote as we no longer require it
			k.DeleteOracleDataVote(ctx, val)
			return false
		})

		// iterate over the full list of validators to increment miss counters
		totalPower := int64(0)
		for _, val := range k.StakingKeeper.GetBondedValidatorsByPower(ctx) {
			totalPower += val.GetConsensusPower()
			valaddr := sdk.AccAddress(val.GetOperator())
			if !contains(voted, valaddr.String()) {
				k.IncrementMissCounter(ctx, valaddr)
			}
		}

		// if the voted_power/total_power > params.VoteThreshold then we store the averages in the store
		storeAverages := sdk.NewDec(votedPower).Quo(sdk.NewDec(totalPower)).GT(params.VoteThreshold)

		// compute the averages for each type of data tracked by the oracle
		for typ, dataCollection := range collectionMap {
			// first, lets delete the old data
			k.DeleteOracleData(ctx, typ)

			// then we compute the "average"
			avg := types.GetAverageFunction(typ)(dataCollection)

			// once we have an "average" we set it in the store
			if storeAverages {
				k.SetOracleData(ctx, avg)
			}

			// store the "average" for scoring validators later
			averageMap[typ] = avg
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
		k.IterateMissCounters(ctx, func(val sdk.AccAddress, counter int64) bool {

			// TODO: Review this logic please! cc @fedekunze @zmanain
			if params.MinValidPerWindow.LT(sdk.NewDec(counter).Quo(sdk.NewDec(params.SlashWindow))) {
				sval := k.StakingKeeper.Validator(ctx, sdk.ValAddress(val))
				cons, _ := sval.GetConsAddr()
				k.StakingKeeper.Slash(ctx, cons, ctx.BlockHeight(), sval.GetConsensusPower(), params.SlashFraction)
			}
			return false
		})

		// TODO: reward validators
		// TODO: Setup module account for oracle module

		// Reset state prior to next round
		// After the tallying is done, reset the vote period start height and delete all the prevotes
		k.SetVotePeriodStart(ctx, ctx.BlockHeight())
		k.IterateOracleDataPrevotes(ctx, func(val sdk.AccAddress, _ *types.MsgOracleDataPrevote) bool {
			k.DeleteOracleDataPrevote(ctx, val)
			return false
		})
	}
}

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}
