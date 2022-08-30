package keeper

import (
	"context"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"

	"github.com/peggyjv/sommelier/v4/x/auction/types"
)

var _ types.MsgServer = Keeper{}

// SubmitBid implements types.MsgServer
func (k Keeper) SubmitBid(c context.Context, msg *types.MsgSubmitBidRequest) (*types.MsgSubmitBidResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)

	// Verify signer is the same as the bidder
	signer := msg.MustGetSigner()

	if !signer.Equals(sdk.AccAddress(msg.Bidder)) {
		return &types.MsgSubmitBidResponse{}, types.ErrSignerDifferentFromBidder
	}

	// Verify auction
	auction, found := k.GetActiveAuctionById(ctx, msg.GetAuctionId())
	if !found {
		return &types.MsgSubmitBidResponse{}, types.ErrAuctionNotFound
	}

	// Verify bidder address
	if _, err := sdk.AccAddressFromBech32(msg.GetBidder()); err != nil {
		return &types.MsgSubmitBidResponse{}, types.ErrInvalidAddress
	}

	// Verify auction coin type and bidder coin type are equal
	if auction.StartingAmount.Denom != msg.MinimumSaleTokenPurchaseAmount.GetDenom() {
		return &types.MsgSubmitBidResponse{}, types.ErrBidAuctionDenomMismatch
	}

	// Query our module address for funds
	totalSaleTokenBalanceForSale := k.bankKeeper.GetBalance(ctx, authtypes.NewModuleAddress(types.ModuleName), auction.StartingAmount.Denom)

	// Calculate minimum purchase price
	currentSaleUnitPriceInUsomm, err := auction.CurrentUnitPriceInUsomm.Float64()
	if err != nil {
		return &types.MsgSubmitBidResponse{}, types.ErrConvertingTokenPriceToFloat
	}

	minimumSaleTokenPurchaseAmount, err := msg.MinimumSaleTokenPurchaseAmount.Amount.ToDec().Float64()
	if err != nil {
		return &types.MsgSubmitBidResponse{}, types.ErrConvertingTokenPriceToFloat
	}

	maxBidInUsomm, err := msg.MaxBidInUsomm.Amount.ToDec().Float64()
	if err != nil {
		return &types.MsgSubmitBidResponse{}, types.ErrConvertingTokenPriceToFloat
	}

	minimumPurchasePriceInUsomm := currentSaleUnitPriceInUsomm * minimumSaleTokenPurchaseAmount

	// Verify minimum price is <= bid
	if minimumPurchasePriceInUsomm > maxBidInUsomm {
		return &types.MsgSubmitBidResponse{}, types.ErrMinimumPurchaseLargerThanBid
	}

	// Verify bid is >= auction base denom price (aka minimum purchase price); note in theory this can be < 1
	if msg.MaxBidInUsomm.Amount.ToDec().LT(auction.CurrentUnitPriceInUsomm) {
		return &types.MsgSubmitBidResponse{}, types.ErrBidSmallerThanMinimumPurchasePrice
	}

	// Start off fulfilled sale token amount at 0
	totalFulfilledSaleTokenAmount := sdk.NewCoin(auction.StartingAmount.GetDenom(), sdk.NewInt(0))

	// Note this is the quotient of the divison operation so fractions are truncated appropriately
	largestSaleTokenAmountPossibleToPurchaseForBid := msg.MaxBidInUsomm.Amount.Quo(sdk.Int(auction.CurrentUnitPriceInUsomm))

	// Verify you can actually purchase at least 1 sale token denom with this amount; note this is important as sale token unit price can be < 1
	if !largestSaleTokenAmountPossibleToPurchaseForBid.IsPositive() {
		return &types.MsgSubmitBidResponse{}, types.ErrInsufficientBid
	}

	// Figure out how much of bid we can fulfill
	if totalSaleTokenBalanceForSale.Amount.GTE(largestSaleTokenAmountPossibleToPurchaseForBid) {
		totalFulfilledSaleTokenAmount.Amount = largestSaleTokenAmountPossibleToPurchaseForBid

	} else if totalSaleTokenBalanceForSale.Amount.GTE(msg.MinimumSaleTokenPurchaseAmount.Amount) {
		totalFulfilledSaleTokenAmount.Amount = totalSaleTokenBalanceForSale.Amount
	} else {
		return &types.MsgSubmitBidResponse{}, types.ErrMinimumPurchaseAmountLargerThanTokensRemaining
	}

	totalFulfilledSaleTokenAmountFloat, err := totalFulfilledSaleTokenAmount.Amount.ToDec().Float64()
	if err != nil {
		return &types.MsgSubmitBidResponse{}, types.ErrConvertingTokenPriceToFloat
	}

	uSommToBePaid := int64(totalFulfilledSaleTokenAmountFloat * currentSaleUnitPriceInUsomm)

	// Transfer payment first
	if err := k.bankKeeper.SendCoinsFromAccountToModule(ctx, sdk.AccAddress(msg.GetBidder()), types.ModuleName, sdk.NewCoins(sdk.NewCoin("usomm", sdk.NewInt(uSommToBePaid)))); err != nil {
		return &types.MsgSubmitBidResponse{}, err
	}

	// Transfer purchase to bidder
	if err := k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, sdk.AccAddress(msg.GetBidder()), sdk.NewCoins(totalFulfilledSaleTokenAmount)); err != nil {
		return &types.MsgSubmitBidResponse{}, err
	}

	// Update amount remaining in auction
	auction.AmountRemaining.Amount = totalSaleTokenBalanceForSale.Amount.Sub(totalFulfilledSaleTokenAmount.Amount)
	k.setActiveAuction(ctx, auction)

	newBidId := k.GetLastBidId(ctx) + 1
	bid := types.Bid{
		Id:                            newBidId,
		AuctionId:                     msg.GetAuctionId(),
		Bidder:                        msg.GetBidder(),
		MaxBid:                        msg.GetMaxBidInUsomm(),
		MinimumAmount:                 msg.GetMinimumSaleTokenPurchaseAmount(),
		TotalFulfilledSaleTokenAmount: totalFulfilledSaleTokenAmount,
		UnitPriceOfSaleTokenInUsomm:   auction.CurrentUnitPriceInUsomm,
	}

	// Create bid in store
	k.setBid(ctx, bid)

	// Update latest bid id
	k.setLastBidId(ctx, newBidId)

	// Verify auction still has supply to see if we need to finish it
	if totalSaleTokenBalanceForSale.Amount.IsZero() {
		// Move to ended auctions if so
		k.FinishAuction(ctx, &auction)
	}

	// Emit Event to signal bid was made
	ctx.EventManager().EmitEvents(
		sdk.Events{
			sdk.NewEvent(
				sdk.EventTypeMessage,
				sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			),
			sdk.NewEvent(
				types.EventTypeBid,
				sdk.NewAttribute(types.AttributeKeyAuctionId, fmt.Sprint(msg.GetAuctionId())),
				sdk.NewAttribute(types.AttributeKeyBidId, fmt.Sprint(newBidId)),
				sdk.NewAttribute(types.AttributeKeyBidder, fmt.Sprint(msg.GetBidder())),
				sdk.NewAttribute(types.AttributeKeyMinimumAmount, msg.GetMinimumSaleTokenPurchaseAmount().String()),
				sdk.NewAttribute(types.AttributeKeySigner, msg.GetSigner()),
				sdk.NewAttribute(types.AttributeKeyFulfilledPrice, fmt.Sprintf("%f", currentSaleUnitPriceInUsomm)),
				sdk.NewAttribute(types.AttributeKeyFulfilledAmount, totalFulfilledSaleTokenAmount.String()),
			),
		},
	)

	return &types.MsgSubmitBidResponse{Bid: &bid}, nil
}
