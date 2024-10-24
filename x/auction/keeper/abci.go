package keeper

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/peggyjv/sommelier/v8/x/auction/types"
)

// BeginBlocker is called at the beginning of every block
func (k Keeper) BeginBlocker(ctx sdk.Context) {}

// EndBlocker is called at the end of every block
func (k Keeper) EndBlocker(ctx sdk.Context) {
	// Auction price updates
	for _, auction := range k.GetActiveAuctions(ctx) {
		if ctx.BlockHeight() != int64(auction.StartBlock) && ((ctx.BlockHeight()-int64(auction.StartBlock))%int64(auction.PriceDecreaseBlockInterval)) == 0 {
			decreaseAccelerationFactor := k.GetParamSet(ctx).AuctionPriceDecreaseAccelerationRate
			// Constant decrease rate if acceleration factor is 0
			// Otherwise decrease rate with acceleration factor that accelerates the decrease amount until a bid is seen

			// Cycle through bids first to see if one was found in the last decrease interval
			bids := k.GetBidsByAuctionID(ctx, auction.Id)

			// Only look at the most recent bid, with the highest ID we can see if we had any bids in the last block decrease period or not
			totalBids := len(bids)

			// Reset decrease rate if we've seen at least 1 bid in the last interval
			if totalBids > 0 && bids[totalBids-1].BlockHeight >= uint64(ctx.BlockHeight())-auction.PriceDecreaseBlockInterval {
				auction.CurrentPriceDecreaseRate = auction.InitialPriceDecreaseRate
			} else { // Otherwise add in the acceleration factor
				auction.CurrentPriceDecreaseRate = auction.CurrentPriceDecreaseRate.Mul(sdk.MustNewDecFromStr("1.0").Add(decreaseAccelerationFactor))
			}

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

	//  Prune bids and auctions -- keep last inactive auction per denom (+ bids) at minimum -- PLUS auctions still in the param freshness window
	auctionMaxBlockAge := k.GetParamSet(ctx).AuctionMaxBlockAge

	// Don't prune if 0
	if auctionMaxBlockAge != 0 {
		denomsFound := make(map[string]bool)
		endedAuctions := k.GetEndedAuctions(ctx)

		for i := 0; i < len(endedAuctions); i++ {
			// Iterate in reverse to guarantee we keep at least the most recent denom auction
			auction := endedAuctions[len(endedAuctions)-1-i]
			denom := auction.StartingTokensForSale.Denom

			// Check if denom already exists, if so and is prune-able delete both it and it's bids
			if _, ok := denomsFound[denom]; ok && auction.EndBlock < uint64(ctx.BlockHeight())-auctionMaxBlockAge {
				bids := k.GetBidsByAuctionID(ctx, auction.GetId())

				for _, bid := range bids {
					k.deleteBid(ctx, auction.GetId(), bid.GetId())
				}

				k.deleteEndedAuction(ctx, auction.GetId())
			} else { // If it doesnt exist/is fresh enough note and skip deletion (since we're iterating in reverse this includes at least the most recent auction for a denom)
				denomsFound[denom] = true
			}
		}
	}
}
