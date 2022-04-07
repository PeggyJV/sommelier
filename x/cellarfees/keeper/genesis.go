package keeper

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/peggyjv/sommelier/v4/x/cellarfees/types"
)

// InitGenesis initializes the module's state from a provided genesis
// state.
func InitGenesis(ctx sdk.Context, k Keeper, gs types.GenesisState) {
	k.SetParams(ctx, gs.Params)

	feesAccount := k.GetFeesAccount(ctx)
	if feesAccount == nil {
		panic(fmt.Sprintf("%s module account has not been set", types.ModuleName))
	}

	// TODO(bolten): I added this check since other cosmos-sdk modules were doing this,
	// but it seems to me that SetModuleAccount is getting called when GetModuleAccount
	// has to create a new module account anyway, which will happen by proxy when we call
	// GetFeesAccount, so not sure how necessary this is
	balances := k.bankKeeper.GetAllBalances(ctx, feesAccount.GetAddress())
	if balances.IsZero() {
		k.accountKeeper.SetModuleAccount(ctx, feesAccount)
	}
}

// ExportGenesis returns the module's exported genesis.
func ExportGenesis(ctx sdk.Context, k Keeper) types.GenesisState {
	return types.GenesisState{
		Params: k.GetParams(ctx),
	}
}
