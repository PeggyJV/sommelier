package keeper

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	cellarfeestypes "github.com/peggyjv/sommelier/v9/x/cellarfees/types"
	types "github.com/peggyjv/sommelier/v9/x/cellarfees/types/v2"
)

// InitGenesis initializes the module's state from a provided genesis
// state.
func (k Keeper) InitGenesis(ctx sdk.Context, gs types.GenesisState) {
	k.SetParams(ctx, gs.Params)
	k.SetLastRewardSupplyPeak(ctx, gs.LastRewardSupplyPeak)

	feesAccount := k.GetFeesAccount(ctx)
	if feesAccount == nil {
		panic(fmt.Sprintf("%s module account has not been set", cellarfeestypes.ModuleName))
	}
}

// ExportGenesis returns the module's exported genesis.
func (k Keeper) ExportGenesis(ctx sdk.Context) types.GenesisState {
	return types.GenesisState{
		Params:               k.GetParams(ctx),
		LastRewardSupplyPeak: k.GetLastRewardSupplyPeak(ctx),
	}
}
