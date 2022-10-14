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
		if ctx.BlockHeight() != int64(auction.StartBlock) && ((ctx.BlockHeight()-int64(auction.StartBlock))%int64(auction.PriceDecreaseBlockInterval)) == 0 {
			// TODO post MVP (pbal) Make a more intricate & responsive step function for auction price updates -- come back to this after other todos
			// Need accelerationFactor param, reset to initialDecreaseRate after some time/criteria, handler for 0 factor case, just incase
			decreaseAccelerationFactor := k.GetParamSet(ctx).AuctionPriceDecreaseAccelerationRate
			newUnitPriceInUsomm := sdk.MustNewDecFromStr("0.0")

			// Constant decrease rate if acceleration factor is 0
			if decreaseAccelerationFactor.Equal(sdk.MustNewDecFromStr("0.0")) {
				priceDecreaseAmountInUsomm := auction.InitialUnitPriceInUsomm.Mul(auction.CurrentPriceDecreaseRate)
				newUnitPriceInUsomm= auction.CurrentUnitPriceInUsomm.Sub(priceDecreaseAmountInUsomm)
			} else { // Decrease rate with acceleration factor that accelerates the decrease amount until a bid is seen
				// Cycle through bids first to see if one was found in the last decrease interval
				bids := k.GetBidsByAuctionID(ctx, auction.Id)

				// Only look at the most recent bid, with the highest ID we can see if we had any bids in the last block decrease period or not
				bid := bids[len(bids)-1]
				
				if bid





			}






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
	denomsFound := make(map[string]bool)
	blocksToNotPrune := k.GetParamSet(ctx).BlocksToNotPrune
	auctionsToDelete := make(map[uint32]bool)

	// Iterate in reverse to guarantee we keep at least the most recent denom auction
	endedAuctions := k.GetEndedAuctions(ctx)
	for i := 0; i < len(endedAuctions); i++ {
		auction := endedAuctions[len(endedAuctions)-1-i]
		denom := auction.StartingTokensForSale.Denom

		// Check if denom already exists, if so and is prune-able, slate for deletion
		if _, ok := denomsFound[denom]; ok && auction.EndBlock < uint64(ctx.BlockHeight())-blocksToNotPrune {
			auctionsToDelete[k.GetLastAuctionID(ctx)] = true
		} else { // If it doesnt exist/is fresh enough note and skip deletion (since we're iterating in reverse this includes at least the most recent auction for a denom)
			denomsFound[denom] = true
		}
	}

	// Delete selected auctions and their associated bids
	for auctionID := range auctionsToDelete {
		bids := k.GetBidsByAuctionID(ctx, auctionID)

		for _, bid := range bids {
			k.deleteBid(ctx, auctionID, bid.Id)
		}

		k.deleteEndedAuction(ctx, auctionID)
	}
}
