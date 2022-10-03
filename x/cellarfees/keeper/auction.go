package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k Keeper) HandleAuctions(ctx sdk.Context) {
	pool := k.GetCellarFeePool(ctx).Pool
	if pool.Empty() {
		// Schedule next auction a short delay away until fees are found to auction
		k.SetScheduledAuctionHeight(ctx, k.GetScheduledAuctionHeight(ctx)+100)
		return
	}

	// Get coin denominations that don't have an active auction
	activeAuctions := k.auctionKeeper.GetActiveAuctions(ctx)
	eligibleCoins := sdk.Coins{}
	coinsToZero := sdk.Coins{}
	for _, coin := range pool {
		found := false

		for _, auction := range activeAuctions {
			if coin.Denom == auction.StartingTokensForSale.Denom {
				found = true
				break
			}
		}

		if !found && !coin.IsZero() {
			eligibleCoins = append(eligibleCoins, coin)
			coinsToZero = append(coinsToZero, sdk.Coin{Amount: sdk.ZeroInt(), Denom: coin.Denom})
		}
	}

	if eligibleCoins.Empty() {
		k.ScheduleNextAuction(ctx)
		return
	}

	// Start auctions
	params := k.GetParams(ctx)
	cellarfeesAccount := string(k.GetFeesAccount(ctx).GetAddress())
	for _, coin := range eligibleCoins {
		// TO-DO (Collin): handle this error
		k.auctionKeeper.BeginAuction(
			ctx,
			coin,
			params.InitialPriceDecreaseRate,
			params.PriceDecreaseBlockInterval,
			cellarfeesAccount,
			cellarfeesAccount,
		)
	}

	// Update the pool to zero out coins sent to auction and schedule the next auction.
	k.setPoolCoins(ctx, coinsToZero)
	k.ScheduleNextAuction(ctx)
}

func (k Keeper) ScheduleNextAuction(ctx sdk.Context) {
	lastAuctionHeight := k.GetScheduledAuctionHeight(ctx)
	// next = last + delay param
	nextAuctionHeight := lastAuctionHeight + k.GetParams(ctx).AuctionBlockDelay
	k.SetScheduledAuctionHeight(ctx, nextAuctionHeight)
}
