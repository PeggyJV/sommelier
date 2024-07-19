package keeper

import (
	"sort"

	sdk "github.com/cosmos/cosmos-sdk/types"
	types "github.com/peggyjv/sommelier/v7/x/cork/types/v2"
)

// InitGenesis initialize default parameters
// and the keeper's address to pubkey map
func InitGenesis(ctx sdk.Context, k Keeper, gs types.GenesisState) {
	k.SetParams(ctx, gs.Params)
	k.SetCellarIDs(ctx, gs.CellarIds)
	k.SetLatestInvalidationNonce(ctx, gs.InvalidationNonce)

	for _, corkResult := range gs.CorkResults {
		k.SetCorkResult(ctx, corkResult.Cork.IDHash(corkResult.BlockHeight), *corkResult)
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
	var cellarIDSet types.CellarIDSet
	existingCellarIDs := k.GetCellarIDs(ctx)
	cellarIDs := make([]string, 0)
	for _, id := range existingCellarIDs {
		cellarIDs = append(cellarIDs, id.String())
	}
	sort.Strings(cellarIDs)
	cellarIDSet.Ids = cellarIDs

	return types.GenesisState{
		Params:            k.GetParamSet(ctx),
		CellarIds:         cellarIDSet,
		InvalidationNonce: k.GetLatestInvalidationNonce(ctx),
		ScheduledCorks:    k.GetScheduledCorks(ctx),
		CorkResults:       k.GetCorkResults(ctx),
	}
}
