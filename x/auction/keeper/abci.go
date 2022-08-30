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

			initialUSDPrice, found := k.GetTokenPrice(ctx, auction.StartingAmount.Denom)
			if !found {
				panic("no price data found for starting amount denom")
			}

			usommPrice, found := k.GetTokenPrice(ctx, "usomm")

			if !found {
				panic("no price data found for usomm")
			}

			// Note we do not check price freshness above, freshness is only checked at auction start time

			initialSaleTokenUnitPriceInUsomm := usommPrice.UsdPrice.MustFloat64() / initialUSDPrice.UsdPrice.MustFloat64()

			// Simple constant decrease rate meant for MVP
			priceDecreaseAmountInUsomm := initialSaleTokenUnitPriceInUsomm * float64(auction.CurrentDecreaseRate)
			newUnitPriceInUsomm := auction.CurrentUnitPriceInUsomm.MustFloat64() - priceDecreaseAmountInUsomm

			// If the new price would be non positive, finish the auction
			if newUnitPriceInUsomm <= 0 {
				err := k.FinishAuction(ctx, auction)

				if err != nil {
					panic(err)
				}
			} else { // Otherwise update the auction with the newest price
				// Note we are not truncating the unit price (if usomm is stronger than the fee token, users will have to bid for a minimum number of fee tokens to make a purchase)
				auction.CurrentUnitPriceInUsomm = sdk.MustNewDecFromStr(fmt.Sprintf("%f", newUnitPriceInUsomm))

				// Update stored auction
				k.setActiveAuction(ctx, *auction)

				// Emit event that auction has finished
				ctx.EventManager().EmitEvents(
					sdk.Events{
						sdk.NewEvent(
							sdk.EventTypeMessage,
							sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
						),
						sdk.NewEvent(
							types.EventTypeAuctionUpdated,
							sdk.NewAttribute(types.AttributeKeyAuctionID, fmt.Sprint(auction.Id)),
							sdk.NewAttribute(types.AttributeKeyStartBlock, fmt.Sprint(auction.StartBlock)),
							sdk.NewAttribute(types.AttributeKeyEndBlock, fmt.Sprint(auction.EndBlock)),
							sdk.NewAttribute(types.AttributeKeyInitialDecreaseRate, fmt.Sprintf("%f", auction.InitialDecreaseRate)),
							sdk.NewAttribute(types.AttributeKeyCurrentDecreaseRate, fmt.Sprintf("%f", auction.CurrentDecreaseRate)),
							sdk.NewAttribute(types.AttributeKeyBlockDecreaseInterval, fmt.Sprint(auction.BlockDecreaseInterval)),
							sdk.NewAttribute(types.AttributeKeyStartingDenom, auction.StartingAmount.Denom),
							sdk.NewAttribute(types.AttributeKeyStartingAmount, auction.StartingAmount.Amount.String()),
							sdk.NewAttribute(types.AttributeKeyAmountRemaining, auction.AmountRemaining.String()),
						),
					},
				)
			}
		}
	}

	// TODO post MVP (pbal): prune bids and auctions
}
