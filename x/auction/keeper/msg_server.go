package keeper

import (
	"context"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/peggyjv/sommelier/v4/x/auction/types"
)

var _ types.MsgServer = Keeper{}

// SubmitBid implements types.MsgServer
func (k Keeper) SubmitBid(c context.Context, msg *types.MsgSubmitBidRequest) (*types.MsgSubmitBidResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)

	// Verify auction
	auction, found := k.GetActiveAuctionById(ctx, msg.GetAuctionId())
	if !found {
		return &types.MsgSubmitBidResponse{}, fmt.Errorf("auction not found")
	}

	// Verify bidder
	if _, err := sdk.AccAddressFromBech32(msg.GetBidder()); err != nil {
		return &types.MsgSubmitBidResponse{}, sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, err.Error())
	}

	// Verify auction coin type and bidder coin type are equal
	if auction.StartingAmount.Denom != msg.MinimumAmount.GetDenom() {
		return &types.MsgSubmitBidResponse{}, fmt.Errorf("auction denom different from bid requested denom")
	}

	// Query our module address for funds
	supply := k.bankKeeper.GetSupply(ctx, msg.MinimumAmount.GetDenom())

	// Verify auction still has supply (edge case where it can be 0)
	if supply.Amount.IsZero() {
		return &types.MsgSubmitBidResponse{}, fmt.Errorf("auction has ended")
	}

	// Calculate minimum purchase price
	salePriceFloat := auction.CurrentPrice.Amount.ToDec().MustFloat64()
	minimumPurchase := msg.MinimumAmount.Amount.ToDec().MustFloat64()
	maxBid := msg.MaxBid.Amount.ToDec().MustFloat64()
	minimumPurchasePrice := int64(salePriceFloat / minimumPurchase)

	// Verify minimum price is <= bid
	if minimumPurchasePrice > msg.GetMaxBid().Amount.Int64() {
		return &types.MsgSubmitBidResponse{}, fmt.Errorf("minimum purchase is larger than allocated bid amount")
	}

	// Verify bid is >= auction price
	if msg.MaxBid.Amount.Int64() < auction.CurrentPrice.Amount.Int64() {
		return &types.MsgSubmitBidResponse{}, fmt.Errorf("bid smaller than purchase price")
	}

	fulfilledAmt := sdk.NewCoin(auction.StartingAmount.GetDenom(), sdk.NewInt(0))
	fulfilledPrice := sdk.NewCoin(msg.MaxBid.GetDenom(), sdk.NewInt(0))

	// Figure out how much of bid we can fulfill
	if supply.Amount.GTE(msg.GetMaxBid().Amount) {
		fulfilledAmt.Amount = sdk.NewInt(int64(maxBid / salePriceFloat))
		fulfilledPrice.Amount = fulfilledAmt.Amount.Mul(sdk.NewInt(int64(salePriceFloat)))
	} else if supply.Amount.GTE(msg.GetMinimumAmount().Amount) {
		partialPurchaseAmt := supply.Amount.ToDec().MustFloat64()
		fulfilledAmt.Amount = sdk.NewInt(int64(partialPurchaseAmt / salePriceFloat))
		fulfilledPrice.Amount = fulfilledAmt.Amount.Mul(sdk.NewInt(int64(salePriceFloat)))
	} else {
		return &types.MsgSubmitBidResponse{}, fmt.Errorf("fewer tokens left in auction than bid minimum amount")
	}

	// Transfer payment first
	if err := k.bankKeeper.SendCoinsFromAccountToModule(ctx, sdk.AccAddress(msg.GetBidder()), types.ModuleName, sdk.NewCoins(fulfilledPrice)); err != nil {
		return &types.MsgSubmitBidResponse{}, err
	}

	// Transfer purchase to bidder
	if err := k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, sdk.AccAddress(msg.GetBidder()), sdk.NewCoins(fulfilledAmt)); err != nil {
		// TODO: is it worth panicing here? We basically took a users funds and didnt give them anything in return, seems like a really bad regression
		return &types.MsgSubmitBidResponse{}, err
	}

	// Update amount remaining in auction
	auction.AmountRemaining.Amount = supply.Amount.Sub(fulfilledAmt.Amount)
	k.setActiveAuction(ctx, auction)
	
	// Create bid in store
	k.setBid(ctx, types.Bid{
		Id: k.GetLastBidId(ctx) + 1,
		AuctionId: msg.GetAuctionId(),
		Bidder: msg.GetBidder(),
		MaxBid: msg.GetMaxBid(),
		MinimumAmount: msg.GetMinimumAmount(),
		FulfilledAmount:fulfilledAmt,
		FulfillmentPrice: fulfilledPrice,
	})

	// Update latest bid id
	k.setLastBidId(ctx, k.GetLastBidId(ctx) + 1)

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
				sdk.NewAttribute(types.AttributeKeyBidId, string(k.GetLastBidId(ctx))),
				sdk.NewAttribute(types.AttributeKeyBidder,  fmt.Sprint(msg.GetBidder())),
				sdk.NewAttribute(types.AttributeKeyMinimumAmount, msg.GetMinimumAmount().String()),
				sdk.NewAttribute(types.AttributeKeySigner, msg.GetSigner()),
				sdk.NewAttribute(types.AttributeKeyFulfilledPrice, fulfilledPrice.String()),
				sdk.NewAttribute(types.AttributeKeyFulfilledAmount, fulfilledAmt.String()),
			),
		},
	)

	bid, _ := k.GetBid(ctx, msg.GetAuctionId(), k.GetLastBidId(ctx))
	return &types.MsgSubmitBidResponse{&bid}, nil
}