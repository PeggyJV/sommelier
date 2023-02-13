package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/peggyjv/sommelier/v4/x/incentives/types"
)

// InitGenesis initialize default parameters
// and the keeper's address to pubkey map
func InitGenesis(ctx sdk.Context, k Keeper, gs types.GenesisState) {
	k.SetParams(ctx, gs.Params)
}

func ExportGenesis(ctx sdk.Context, k Keeper) types.GenesisState {
	return types.GenesisState{
		Params: k.GetParamSet(ctx),
	}
}
