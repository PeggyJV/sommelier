package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	paramstypes "github.com/peggyjv/sommelier/v7/app/params"
	"github.com/peggyjv/sommelier/v7/x/cellarfees/types"
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

	for _, tokenPrice := range k.auctionKeeper.GetTokenPrices(ctx) {
		// skip usomm
		if tokenPrice.Denom == paramstypes.BaseCoinUnit {
			continue
		}

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
