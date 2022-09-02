package keeper

import (
	"context"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
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
		return &types.MsgSubmitBidResponse{}, sdkerrors.Wrapf(types.ErrSignerDifferentFromBidder, "Signer: %s, Bidder: %s", signer, msg.GetBidder())
	}

	// Verify auction
	auction, found := k.GetActiveAuctionByID(ctx, msg.GetAuctionId())
	if !found {
		return &types.MsgSubmitBidResponse{}, sdkerrors.Wrapf(types.ErrAuctionNotFound, "Auction id: %d", msg.GetAuctionId())
	}

	// Verify bidder address
	if _, err := sdk.AccAddressFromBech32(msg.GetBidder()); err != nil {
		return &types.MsgSubmitBidResponse{}, sdkerrors.Wrapf(types.ErrInvalidAddress, "Address: %s", msg.GetBidder())
	}

	// Verify auction coin type and bidder coin type are equal
	if auction.GetStartingAmount().Denom != msg.GetMinimumSaleTokenPurchaseAmount().Denom {
		return &types.MsgSubmitBidResponse{}, sdkerrors.Wrapf(types.ErrBidAuctionDenomMismatch, "Bid denom: %s, Auction denom: %s", msg.GetMinimumSaleTokenPurchaseAmount(), auction.GetStartingAmount().Denom)
	}

	// Query our module address for funds
	totalSaleTokenBalanceForSale := k.bankKeeper.GetBalance(ctx, authtypes.NewModuleAddress(types.ModuleName), auction.StartingAmount.Denom)

	// Convert & standardize types for use below
	minimumSaleTokenPurchaseAmount := msg.MinimumSaleTokenPurchaseAmount.Amount.ToDec()
	maxBidInUsomm := msg.MaxBidInUsomm.Amount.ToDec()

	// Calculate minimum purchase price
	minimumPurchasePriceInUsomm := auction.CurrentUnitPriceInUsomm.Mul(minimumSaleTokenPurchaseAmount)

	// Verify minimum price is <= bid
	if minimumPurchasePriceInUsomm.GT(maxBidInUsomm) {
		return &types.MsgSubmitBidResponse{}, sdkerrors.Wrapf(types.ErrMinimumPurchaseLargerThanBid, "minimum: %s, max bid: %s", minimumPurchasePriceInUsomm.String(), maxBidInUsomm.String())
	}

	// Start off fulfilled sale token amount at 0
	totalFulfilledSaleTokenAmount := sdk.NewCoin(auction.StartingAmount.GetDenom(), sdk.NewInt(0))

	// See how many whole tokens of base denom can be purchased
	largestSaleTokenAmountPossibleToPurchaseForBid := msg.MaxBidInUsomm.Amount.ToDec().Quo(auction.CurrentUnitPriceInUsomm).TruncateInt()

	// Verify you can actually purchase at least 1 sale token denom with this amount
	if !largestSaleTokenAmountPossibleToPurchaseForBid.IsPositive() {
		return &types.MsgSubmitBidResponse{}, sdkerrors.Wrapf(types.ErrInsufficientBid, "Price: %s, Bid: %s", auction.CurrentUnitPriceInUsomm.String(), msg.MaxBidInUsomm.String())
	}

	// Figure out how much of bid we can fulfill
	if totalSaleTokenBalanceForSale.Amount.GTE(largestSaleTokenAmountPossibleToPurchaseForBid) {
		totalFulfilledSaleTokenAmount.Amount = largestSaleTokenAmountPossibleToPurchaseForBid

	} else if totalSaleTokenBalanceForSale.Amount.GTE(msg.MinimumSaleTokenPurchaseAmount.Amount) {
		totalFulfilledSaleTokenAmount.Amount = totalSaleTokenBalanceForSale.Amount

	} else {
		return &types.MsgSubmitBidResponse{}, sdkerrors.Wrapf(types.ErrMinimumPurchaseAmountLargerThanTokensRemaining, "Minimum purchase: %s, amount remaining: %s", msg.MinimumSaleTokenPurchaseAmount.String(), auction.AmountRemaining.String())
	}

	uSommToBePaid := totalFulfilledSaleTokenAmount.Amount.ToDec().Mul(auction.CurrentUnitPriceInUsomm).TruncateDec()

	newBidID := k.GetLastBidID(ctx) + 1
	bid := types.Bid{
		Id:                            newBidID,
		AuctionId:                     msg.GetAuctionId(),
		Bidder:                        msg.GetBidder(),
		MaxBid:                        msg.GetMaxBidInUsomm(),
		MinimumAmount:                 msg.GetMinimumSaleTokenPurchaseAmount(),
		TotalFulfilledSaleTokenAmount: totalFulfilledSaleTokenAmount,
		UnitPriceOfSaleTokenInUsomm:   auction.CurrentUnitPriceInUsomm,
		TotalAmountPaidInUsomm:        uSommToBePaid,
	}

	bidErr := bid.ValidateBasic()
	if bidErr != nil {
		return &types.MsgSubmitBidResponse{}, bidErr
	}

	// Transfer payment first
	payment := sdk.NewCoin(types.UsommDenom, uSommToBePaid.TruncateInt())
	if err := k.bankKeeper.SendCoinsFromAccountToModule(ctx, sdk.AccAddress(msg.GetBidder()), types.ModuleName, sdk.NewCoins(payment)); err != nil {
		return &types.MsgSubmitBidResponse{}, err
	}

	// Transfer purchase to bidder
	if err := k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, sdk.AccAddress(msg.GetBidder()), sdk.NewCoins(totalFulfilledSaleTokenAmount)); err != nil {
		return &types.MsgSubmitBidResponse{}, err
	}

	// Update amount remaining in auction
	auction.AmountRemaining.Amount = totalSaleTokenBalanceForSale.Amount.Sub(totalFulfilledSaleTokenAmount.Amount)
	k.setActiveAuction(ctx, auction)

	// Create bid in store
	k.setBid(ctx, bid)

	// Update latest bid id
	k.setLastBidID(ctx, newBidID)

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
				sdk.NewAttribute(types.AttributeKeyAuctionID, fmt.Sprint(msg.GetAuctionId())),
				sdk.NewAttribute(types.AttributeKeyBidID, fmt.Sprint(newBidID)),
				sdk.NewAttribute(types.AttributeKeyBidder, fmt.Sprint(msg.GetBidder())),
				sdk.NewAttribute(types.AttributeKeyMinimumAmount, msg.GetMinimumSaleTokenPurchaseAmount().String()),
				sdk.NewAttribute(types.AttributeKeySigner, msg.GetSigner()),
				sdk.NewAttribute(types.AttributeKeyFulfilledPrice, auction.CurrentUnitPriceInUsomm.String()),
				sdk.NewAttribute(types.AttributeKeyTotalPayment, payment.String()),
				sdk.NewAttribute(types.AttributeKeyFulfilledAmount, totalFulfilledSaleTokenAmount.String()),
			),
		},
	)

	return &types.MsgSubmitBidResponse{Bid: &bid}, nil
}
