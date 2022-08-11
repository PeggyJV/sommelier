package keeper

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/peggyjv/sommelier/v4/x/cellarfees/types"
)

func RegisterInvariants(ir sdk.InvariantRegistry, k Keeper) {
	ir.RegisterRoute(types.ModuleName, "pool-balances", PoolBalanceInvariants(k))
}

func AllInvariants(k Keeper) sdk.Invariant {
	return func(ctx sdk.Context) (string, bool) {
		return PoolBalanceInvariants(k)(ctx)
	}
}

// 1. All coin denominations in the cellar fee pool are present in the cellarfees module account
// 2. The individual balances in the fee pool are not greater than their corresponding module account balances
// 3. Pool balances are non-negative
func PoolBalanceInvariants(k Keeper) sdk.Invariant {
	return func(ctx sdk.Context) (string, bool) {
		var msg string

		balances := k.bankKeeper.GetAllBalances(ctx, k.accountKeeper.GetModuleAddress(types.ModuleName))
		balanceMap := make(map[string]sdk.Coin)
		for _, balance := range balances {
			balanceMap[balance.Denom] = balance
		}

		pool := k.GetCellarFeePool(ctx).Pool
		for _, coin := range pool {
			val, ok := balanceMap[coin.Denom]
			if !ok {
				msg += fmt.Sprintf("Cellar fee pool contains a coin that isn't present in the module account: %#v\n", coin)
			}
			if coin.Amount.GT(val.Amount) {
				msg += fmt.Sprintf("Cellar fee pool contains a balance that is greater than the coins balance in the module account. Pool Coin: %+v, Account Coin: %+v", coin, val)
			}
			if coin.Amount.IsNegative() {
				msg += fmt.Sprintf("Coin balance in pool is negative: %+v", coin)
			}
		}

		broken := len(msg) != 0

		return sdk.FormatInvariant(types.ModuleName, "pool balances", msg), broken
	}
}
