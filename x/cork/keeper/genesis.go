package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/peggyjv/sommelier/v4/x/cork/types"
)

// InitGenesis initialize default parameters
// and the keeper's address to pubkey map
func InitGenesis(ctx sdk.Context, k Keeper, gs types.GenesisState) {
	k.setParams(ctx, gs.Params)
	// Set the vote period at initialization
	k.SetCellarIDs(ctx, gs.CellarIds)
	k.SetLatestInvalidationNonce(ctx, gs.InvalidationNonce)
	k.SetLatestScheduledCorkID(ctx, gs.LatestCorkId)

	for _, corkResults := range gs.CorkResults {
		k.SetCorkResult(ctx, gs.LatestCorkId, *corkResults)
	}

	for _, scheduledCork := range gs.ScheduledCorks {
		valAddr, err := sdk.ValAddressFromHex(scheduledCork.Validator)
		if err != nil {
			panic(err)
		}

		k.SetScheduledCork(ctx, scheduledCork.BlockHeight, valAddr, *scheduledCork.Cork)
	}
}

// ExportGenesis writes the current store values
// to a genesis file, which can be imported again
// with InitGenesis
func ExportGenesis(ctx sdk.Context, k Keeper) types.GenesisState {
	var ids types.CellarIDSet
	for _, id := range k.GetCellarIDs(ctx) {
		ids.Ids = append(ids.Ids, id.Hex())
	}

	return types.GenesisState{
		Params:            k.GetParamSet(ctx),
		CellarIds:         ids,
		InvalidationNonce: k.GetLatestInvalidationNonce(ctx),
		LatestCorkId:      k.GetLatestScheduledCorkID(ctx),
		ScheduledCorks:    k.GetScheduledCorks(ctx),
		CorkResults:       k.GetCorkResults(ctx),
	}
}
