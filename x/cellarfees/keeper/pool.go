package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/peggyjv/sommelier/v4/app/params"
	"github.com/peggyjv/sommelier/v4/x/cellarfees/types"
)

// Appends the coin to the pool coins if denom isn't already present, otherwise add the amount to the
// existing balance
func (k Keeper) addCoinToPool(ctx sdk.Context, coin sdk.Coin) {
	// (Collin) should an empty coin error?
	if (coin == sdk.Coin{} || coin.IsZero()) {
		return
	}

	if coin.Denom == params.BaseCoinUnit {
		panic("Cannot add SOMM to cellar fee pool")
	}

	pool := k.GetCellarFeePool(ctx)
	pool.Pool = pool.Pool.Add(coin)
	k.SetCellarFeePool(ctx, pool)
}

// Getter for module account that holds the fee pool funds
func (k Keeper) GetFeesAccount(ctx sdk.Context) authtypes.ModuleAccountI {
	return k.accountKeeper.GetModuleAccount(ctx, types.ModuleName)
}
