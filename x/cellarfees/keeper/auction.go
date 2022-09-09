package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/peggyjv/sommelier/v4/x/cellarfees/types"
)

func (k Keeper) HandleAuctions(ctx sdk.Context) {
	pool := k.GetCellarFeePool(ctx).Pool
	if pool.Empty() {
		// Schedule next auction a short delay away until fees are found to auction
		k.SetScheduledAuctionHeight(ctx, k.GetScheduledAuctionHeight(ctx).Add(sdk.NewInt(100)))
		return
	}

	// Get coin denominations that don't have an active auction
	activeAuctions := k.auctionKeeper.GetActiveAuctions(ctx)
	eligibleCoins := sdk.Coins{}
	newPool := sdk.Coins{}
	for _, coin := range pool {
		found := false

		for _, auction := range activeAuctions {
			if coin.Denom == auction.StartingTokensForSale.Denom {
				found = true
				break
			}
		}

		if !found {
			eligibleCoins = append(eligibleCoins, coin)
		} else {
			newPool = append(newPool, coin)
		}
	}

	if eligibleCoins.Empty() {
		k.ScheduleNextAuction(ctx)
		return
	}

	// Start auctions
	params := k.GetParams(ctx)
	cellarfeesAccount := string(k.accountKeeper.GetModuleAddress(types.ModuleName))
	for _, coin := range eligibleCoins {
		k.auctionKeeper.BeginAuction(
			ctx,
			coin,
			params.InitialPriceDecreaseRate,
			params.PriceDecreaseBlockInterval,
			cellarfeesAccount,
			cellarfeesAccount,
		)
	}

	// Update the pool to remove coins sent to auction and schedule the next auction.
	k.SetCellarFeePool(ctx, types.CellarFeePool{Pool: newPool})
	k.ScheduleNextAuction(ctx)
}

func (k Keeper) ScheduleNextAuction(ctx sdk.Context) {
	lastAuctionHeight := k.GetScheduledAuctionHeight(ctx)
	// next = last + delay param
	nextAuctionHeight := lastAuctionHeight.Add(sdk.NewInt(int64(k.GetParams(ctx).AuctionBlockDelay)))
	k.SetScheduledAuctionHeight(ctx, nextAuctionHeight)
}
