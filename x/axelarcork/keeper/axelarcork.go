package keeper

import (
	"math/big"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	distributionTypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	"github.com/peggyjv/sommelier/v9/x/axelarcork/types"
)

var maxSafeInt math.Int

// Taken from https://github.com/cosmos/cosmos-sdk/blob/25b14c3caa2ecdc99840dbb88fdb3a2d8ac02158/math/dec.go#L32
const maxDecBitLen = 315

func init() {
	// Define the maximum integer value that can be safely converted to a LegacyDec
	var ok bool
	maxSafeInt, ok = sdk.NewIntFromString("115792089237316195423570985008687907853269984665640564039457")
	if !ok {
		panic("failed to parse max safe integer for DecCoin conversion")
	}

}

// SweepModuleAccountBalances sweeps all safe axelarcork sender module account balances to the community pool. Because this account is the
// sender for transfers created by RelayCork calls, funds will not be returned to the caller if the IBC
// transfer fails or gas is refunded.
func (k Keeper) SweepModuleAccountBalances(ctx sdk.Context) {
	moduleAcct := k.GetSenderAccount(ctx)
	balances := k.bankKeeper.GetAllBalances(ctx, moduleAcct.GetAddress())

	if balances.IsZero() {
		return
	}

	safeBalancesToSend := sdk.Coins{}
	balancesToBurn := sdk.Coins{}
	feePool := k.distributionKeeper.GetFeePool(ctx)
	communityPool := &feePool.CommunityPool

	// Only sweep coins that will not cause overflows
	for _, coin := range balances {
		safe := AddToCommunityPoolIfSafe(communityPool, coin)
		if safe {
			safeBalancesToSend = safeBalancesToSend.Add(coin)
		} else {
			balancesToBurn = balancesToBurn.Add(coin)
		}
	}

	if len(balancesToBurn) > 0 {
		ctx.Logger().With("module", "x/"+types.ModuleName).Info("burning unsafe coins from axelarcork module account", "coins", balancesToBurn.String())
		k.bankKeeper.BurnCoins(ctx, types.ModuleName, balancesToBurn)
	}

	if safeBalancesToSend.IsZero() {
		return
	}

	ctx.Logger().With("module", "x/"+types.ModuleName).Info("sweeping funds from axelarcork module account to community pool", "coins", safeBalancesToSend.String())

	if err := k.bankKeeper.SendCoinsFromModuleToModule(ctx, types.ModuleName, distributionTypes.ModuleName, safeBalancesToSend); err != nil {
		panic(err)
	}

	k.distributionKeeper.SetFeePool(ctx, feePool)
}

// AddToCommunityPoolIfSafe checks if adding a coin to the community pool will cause an overflow.
// If so, it returns false. Otherwise, it adds the coin to the community pool and returns true.
func AddToCommunityPoolIfSafe(communityPool *sdk.DecCoins, coin sdk.Coin) bool {
	// Check if the Int -> LegacyDec conversion will overflow
	if coin.Amount.GT(maxSafeInt) {
		return false
	}

	decCoin := sdk.NewDecCoinFromCoin(coin)

	// Check if the LegacyDec addition will overflow
	existingAmount := communityPool.AmountOf(coin.Denom)

	if !IsSafeToAdd(existingAmount, decCoin.Amount) {
		return false
	}

	*communityPool = communityPool.Add(decCoin)
	return true
}

// IsSafeToAdd checks if adding two LegacyDecs will overflow.
func IsSafeToAdd(a, b math.LegacyDec) bool {
	// Check if sum would exceed max bit length. This is probably overkill because the max bit length of sdkmath.Int which sdk.Coin uses is 256 bits while maxDecBitLen is 315, but oh well.
	sum := new(big.Int).Add(a.BigInt(), b.BigInt())
	return sum.BitLen() <= maxDecBitLen
}
