package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/peggyjv/sommelier/v4/app/params"
	"github.com/peggyjv/sommelier/v4/x/cellarfees/types"
)

// Set the specified denoms' balances to the respective amounts
func (k Keeper) setPoolCoins(ctx sdk.Context, coins sdk.Coins) {
	pool := k.GetCellarFeePool(ctx)

	for _, coin := range coins {
		if (coin == sdk.Coin{}) {
			panic("Attempted to set empty coin in cellar fee pool")
		}

		if coin.Amount.IsNegative() {
			panic("Attempted to set negative coin balance")
		}

		if coin.Denom == params.BaseCoinUnit {
			panic("Attempted to add SOMM to cellar fee pool")
		}

		balance := pool.Pool.AmountOf(coin.Denom)
		if !balance.IsZero() {
			pool.Pool.Sub(sdk.Coins{{Amount: balance, Denom: coin.Denom}})
		}

		pool.Pool.Add(coin)
	}

	k.SetCellarFeePool(ctx, pool)
}

// Getter for module account that holds the fee pool funds
func (k Keeper) GetFeesAccount(ctx sdk.Context) authtypes.ModuleAccountI {
	return k.accountKeeper.GetModuleAccount(ctx, types.ModuleName)
}
