package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	paramstypes "github.com/peggyjv/sommelier/v9/app/params"
	"github.com/peggyjv/sommelier/v9/x/cellarfees/types"
)

// BeginBlocker emits rewards each block they are available by sending them to the distribution module's fee collector
// account. Emissions are a constant value based on the last peak supply of distributable fees so that the reward supply
// will decrease linearly until exhausted.
func (k Keeper) BeginBlocker(ctx sdk.Context) {
	k.handleRewardEmission(ctx)
	k.handleFeeAuctions(ctx)
}

// EndBlocker is called at the end of every block
func (k Keeper) EndBlocker(ctx sdk.Context) {}

func (k Keeper) handleRewardEmission(ctx sdk.Context) {
	moduleAccount := k.GetFeesAccount(ctx)
	remainingRewardsSupply := k.bankKeeper.GetBalance(ctx, moduleAccount.GetAddress(), paramstypes.BaseCoinUnit).Amount

	if remainingRewardsSupply.IsZero() {
		return
	}

	emission := k.GetEmission(ctx, remainingRewardsSupply)

	// sanity check. the upcoming bank keeper call will error causing a panic if this is zero
	if emission.IsZero() {
		return
	}

	// Send to fee collector for distribution
	ctx.Logger().Info("Sending rewards to fee collector", "module", types.ModuleName)
	err := k.bankKeeper.SendCoinsFromModuleToModule(ctx, moduleAccount.GetName(), authtypes.FeeCollectorName, emission)
	if err != nil {
		panic(err)
	}
}

func (k Keeper) handleFeeAuctions(ctx sdk.Context) {
	params := k.GetParams(ctx)

	if uint64(ctx.BlockHeight())%params.AuctionInterval != 0 {
		return
	}

	for _, balance := range k.bankKeeper.GetAllBalances(ctx, k.GetFeesAccount(ctx).GetAddress()) {
		// skip usomm
		if balance.Denom == paramstypes.BaseCoinUnit {
			continue
		}

		if balance.IsZero() {
			continue
		}

		tokenPrice, found := k.auctionKeeper.GetTokenPrice(ctx, balance.Denom)
		if !found {
			continue
		}

		usdValue := k.GetBalanceUsdValue(ctx, balance, tokenPrice)

		if usdValue.GTE(params.AuctionThresholdUsdValue) {
			// Send portion to proceeds module
			auctionBalance := k.handleProceeds(ctx, balance)

			if auctionBalance.IsZero() {
				continue
			}

			// Begin auction with remaining balance
			k.beginAuction(ctx, auctionBalance)
		}
	}
}

// handleProceeds transfers the proceeds portion of the balance to the proceeds account and returns the remaining balance for auction
func (k Keeper) handleProceeds(ctx sdk.Context, balance sdk.Coin) sdk.Coin {
	portion := k.GetParams(ctx).ProceedsPortion

	if portion.IsZero() || balance.IsZero() {
		return balance
	}

	zeroBalance := sdk.NewCoin(balance.Denom, sdk.ZeroInt())
	proceedsAccount := k.GetProceedsAccount(ctx)
	if proceedsAccount == nil {
		ctx.Logger().Error("Proceeds account not found", "address", proceedsAddress)
		// Don't auction the funds so that the account can be created and balance processed later
		return zeroBalance
	}

	// Special case: if proceeds portion is 100%, send entire balance to proceeds
	if portion.Equal(sdk.OneDec()) {
		err := k.bankKeeper.SendCoinsFromModuleToAccount(
			ctx,
			types.ModuleName,
			proceedsAccount.GetAddress(),
			sdk.NewCoins(balance),
		)
		if err != nil {
			ctx.Logger().Error("Error sending proceeds to proceeds account", "error", err)
		}

		return zeroBalance
	}

	// Normal case: calculate proceeds amount
	proceedsAmount := balance.Amount.ToLegacyDec().Mul(portion).TruncateInt()
	if proceedsAmount.IsZero() {
		return balance
	}

	auctionBalance := balance.Sub(sdk.NewCoin(balance.Denom, proceedsAmount))
	proceedsCoin := sdk.NewCoin(balance.Denom, proceedsAmount)

	err := k.bankKeeper.SendCoinsFromModuleToAccount(
		ctx,
		types.ModuleName,
		proceedsAccount.GetAddress(),
		sdk.NewCoins(proceedsCoin),
	)
	if err != nil {
		ctx.Logger().Error("Error sending proceeds to proceeds account", "error", err)
	}

	return auctionBalance
}
