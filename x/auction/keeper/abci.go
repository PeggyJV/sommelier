package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// BeginBlocker is called at the beginning of every block
func (k Keeper) BeginBlocker(ctx sdk.Context) {}

// EndBlocker is called at the end of every block
func (k Keeper) EndBlocker(ctx sdk.Context) {
	// End Auctions that have no funds left
	for _, auction := range k.GetActiveAuctions(ctx) {
		if auction.AmountRemaining.Amount.IsZero() {
			err := k.FinishAuction(ctx, auction)

			if err != nil {
				panic(err)
			}
		}
	}

	for _, auction := range k.GetActiveAuctions(ctx) {
		if ((ctx.BlockHeight() - int64(auction.StartBlock)) % int64(auction.BlockDecreaseInterval)) == 0 {
			// TODO post MVP (pbal) Make a more intricate & responsive step function for auction price updates

			// Simple constant decrease rate meant for MVP
			decreaseMultiplier := float64(1 - auction.CurrentDecreaseRate)
			//grossPriceDecrease :=

			//newUnitPriceInUsomm :=

			// If the new price would be <= 0, finish the auction
			if !newUnitPriceInUsomm.IsPositive() {
				err := k.FinishAuction(ctx, auction)

				if err != nil {
					panic(err)
				}
			} else { // Otherwise update the auction with the newest price
				// Note we are not truncating the unit price (if usomm is stronger than the fee token, users will have to bid for a minimum number of fee tokens to make a purchase)
				auction.CurrentUnitPriceInUsomm = newUnitPriceInUsomm

				// Update stored auction
				k.setActiveAuction(ctx, *auction)
			}
		}
	}

	// TODO post MVP (pbal): prune bids and auctions
}
