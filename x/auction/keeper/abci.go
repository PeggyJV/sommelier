package keeper

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/peggyjv/sommelier/v4/x/auction/types"
)

// BeginBlocker is called at the beginning of every block
func (k Keeper) BeginBlocker(ctx sdk.Context) {}

// EndBlocker is called at the end of every block
func (k Keeper) EndBlocker(ctx sdk.Context) {
	// Auction price updates
	for _, auction := range k.GetActiveAuctions(ctx) {
		if ((ctx.BlockHeight() - int64(auction.StartBlock)) % int64(auction.PriceDecreaseBlockInterval)) == 0 {
			// TODO post MVP (pbal) Make a more intricate & responsive step function for auction price updates

			// Simple constant decrease rate meant for MVP
			priceDecreaseAmountInUsomm := auction.InitialUnitPriceInUsomm.Mul(auction.CurrentPriceDecreaseRate)
			newUnitPriceInUsomm := auction.CurrentUnitPriceInUsomm.Sub(priceDecreaseAmountInUsomm)

			// If the new price would be non positive, finish the auction
			if newUnitPriceInUsomm.LTE(sdk.NewDec(0)) {
				err := k.FinishAuction(ctx, auction)

				if err != nil {
					panic(err)
				}
			} else { // Otherwise update the auction with the newest price
				// Note we are not truncating the unit price (if usomm is stronger than the fee token, users will have to bid for a minimum number of fee tokens to make a purchase)
				lastPrice := auction.CurrentUnitPriceInUsomm
				auction.CurrentUnitPriceInUsomm = newUnitPriceInUsomm

				// Update stored auction
				k.setActiveAuction(ctx, *auction)

				// Emit event that auction has updated
				ctx.EventManager().EmitEvent(
					sdk.NewEvent(
						types.EventTypeAuctionUpdated,
						sdk.NewAttribute(types.AttributeKeyAuctionID, fmt.Sprint(auction.Id)),
						sdk.NewAttribute(types.AttributeKeyLastPrice, lastPrice.String()),
						sdk.NewAttribute(types.AttributeKeyNewPrice, auction.CurrentUnitPriceInUsomm.String()),
						sdk.NewAttribute(types.AttributeKeyCurrentDecreaseRate, auction.CurrentPriceDecreaseRate.String()),
						sdk.NewAttribute(types.AttributeKeyBlockDecreaseInterval, fmt.Sprint(auction.PriceDecreaseBlockInterval)),
					),
				)
			}
		}
	}
	// TODO post MVP (pbal): prune bids and auctions
}
