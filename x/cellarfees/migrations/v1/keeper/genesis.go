package keeper

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	v1types "github.com/peggyjv/sommelier/v9/x/cellarfees/migrations/v1/types"
	"github.com/peggyjv/sommelier/v9/x/cellarfees/types"
)

// InitGenesis initializes the module's state from a provided genesis
// state.
func (k Keeper) InitGenesis(ctx sdk.Context, gs v1types.GenesisState) {
	k.SetParams(ctx, gs.Params)
	k.SetFeeAccrualCounters(ctx, gs.FeeAccrualCounters)
	k.SetLastRewardSupplyPeak(ctx, gs.LastRewardSupplyPeak)

	feesAccount := k.GetFeesAccount(ctx)
	if feesAccount == nil {
		panic(fmt.Sprintf("%s module account has not been set", types.ModuleName))
	}
}

// ExportGenesis returns the module's exported genesis.
func (k Keeper) ExportGenesis(ctx sdk.Context) v1types.GenesisState {
	return v1types.GenesisState{
		Params:               k.GetParams(ctx),
		FeeAccrualCounters:   k.GetFeeAccrualCounters(ctx),
		LastRewardSupplyPeak: k.GetLastRewardSupplyPeak(ctx),
	}
}
