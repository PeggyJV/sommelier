package oracle

import (
	"github.com/peggyjv/sommelier/x/oracle/keeper"
	"github.com/peggyjv/sommelier/x/oracle/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// EndBlocker is called at the end of every block
func EndBlocker(ctx sdk.Context, k keeper.Keeper) {
	params := k.GetParamSet(ctx)
	// valset := k.StakingKeeper.GetBondedValidatorsByPower(ctx)

	// if there is not a vote period start, set it to current block height
	if !k.HasVotePeriodStart(ctx) {
		k.SetVotePeriodStart(ctx, ctx.BlockHeight())
	}

	// If the vote period has ended its time to tally votes
	if (ctx.BlockHeight() - k.GetVotePeriodStart(ctx)) >= params.VotePeriod {
		voted := []sdk.AccAddress{}
		power := int64(0)
		detailedMap := make(map[string]map[string]types.OracleData)
		collectionMap := make(map[string][]types.OracleData)
		// initialize the inner maps
		for _, dt := range params.DataTypes {
			detailedMap[dt] = make(map[string]types.OracleData)
		}

		// iterate over the data votes
		k.IterateOracleDataVotes(ctx, func(val sdk.AccAddress, msg *types.MsgOracleDataVote) bool {
			// save a voted array
			voted = append(voted, val)

			// find total voting power
			power += k.StakingKeeper.Validator(ctx, sdk.ValAddress(val)).GetConsensusPower()

			// save the oracle data for later processing
			for _, oda := range msg.OracleData {
				od, err := types.UnpackOracleData(oda)
				if err != nil {
					panic(err)
				}
				detailedMap[od.Type()][val.String()] = od
				collectionMap[od.Type()] = append(collectionMap[od.Type()], od)
			}

			return false
		})

		// After the tallying is done, set the vote period start height
		k.SetVotePeriodStart(ctx, ctx.BlockHeight())

		// ... and delete all the prevotes
		k.IterateOracleDataPrevotes(ctx, func(val sdk.AccAddress, _ [][]byte) bool {
			k.DeleteOracleDataPrevote(ctx, val)
			return false
		})

		// ... and delete all the votes
		k.IterateOracleDataVotes(ctx, func(val sdk.AccAddress, _ *types.MsgOracleDataVote) bool {
			k.DeleteOracleDataVote(ctx, val)
			return false
		})
	}
}
