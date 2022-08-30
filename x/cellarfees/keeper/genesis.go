package keeper

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/peggyjv/sommelier/v4/x/cellarfees/types"
)

// InitGenesis initializes the module's state from a provided genesis
// state.
func (k Keeper) InitGenesis(ctx sdk.Context, gs types.GenesisState) {
	k.SetCellarFeePool(ctx, gs.CellarFeePool)
	k.SetLastRewardSupplyPeak(ctx, sdk.ZeroInt())
	k.SetScheduledAuctionHeight(ctx, sdk.ZeroInt())
	k.SetParams(ctx, gs.Params)

	feesAccount := k.GetFeesAccount(ctx)
	if feesAccount == nil {
		panic(fmt.Sprintf("%s module account has not been set", types.ModuleName))
	}

	k.accountKeeper.SetModuleAccount(ctx, feesAccount)
}

// ExportGenesis returns the module's exported genesis.
func (k Keeper) ExportGenesis(ctx sdk.Context) types.GenesisState {
	return types.GenesisState{
		Params:        k.GetParams(ctx),
		CellarFeePool: k.GetCellarFeePool(ctx),
	}
}
