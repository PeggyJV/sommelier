package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/peggyjv/sommelier/v3/x/cork/types"
)

// InitGenesis initialize default parameters
// and the keeper's address to pubkey map
func InitGenesis(ctx sdk.Context, k Keeper, gs types.GenesisState) {
	k.setParams(ctx, gs.Params)
	// Set the vote period at initialization
	k.SetCommitPeriodStart(ctx, ctx.BlockHeight())
	k.SetCellarIDs(ctx, *gs.CellarIds)
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
		Params:    k.GetParamSet(ctx),
		CellarIds: &ids,
	}
}
