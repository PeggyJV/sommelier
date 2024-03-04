package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/peggyjv/sommelier/v7/app/params"
)

// BeginBlocker emits rewards each block they are available by sending them to the distribution module's fee collector
// account. Emissions are a constant value based on the last peak supply of distributable fees so that the reward supply
// will decrease linearly until exhausted.
func (k Keeper) BeginBlocker(ctx sdk.Context) {
	// Handle reward emissions
	moduleAccount := k.GetFeesAccount(ctx)
	remainingRewardsSupply := k.bankKeeper.GetBalance(ctx, moduleAccount.GetAddress(), params.BaseCoinUnit).Amount

	if remainingRewardsSupply.IsZero() {
		return
	}

	emission := k.GetEmission(ctx, remainingRewardsSupply)

	// Send to fee collector for distribution
	err := k.bankKeeper.SendCoinsFromModuleToModule(ctx, moduleAccount.GetName(), authtypes.FeeCollectorName, emission)
	if err != nil {
		panic(err)
	}

	// Handle fee auctions
	params := k.GetParams(ctx)

	if uint64(ctx.BlockHeight())%params.AuctionInterval != 0 {
		return
	}

	tokenPrices := k.auctionKeeper.GetTokenPrices(ctx)

	for _, tokenPrice := range tokenPrices {
		balance := k.GetFeeBalance(ctx, tokenPrice.Denom)

		if balance.IsZero() {
			continue
		}

		usdValue := k.GetBalanceUsdValue(ctx, balance, tokenPrice)

		if usdValue.GTE(params.AuctionThresholdUsdValue) {
			k.beginAuction(ctx, tokenPrice.Denom)
		}
	}
}

// EndBlocker is called at the end of every block
func (k Keeper) EndBlocker(ctx sdk.Context) {}
