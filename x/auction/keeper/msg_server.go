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

	// Verify signer is the same as the bidder (this validates both the bidder and signer addresses)
	signer := msg.MustGetSigner()
	if !signer.Equals(sdk.AccAddress(msg.Bidder)) {
		return &types.MsgSubmitBidResponse{}, sdkerrors.Wrapf(types.ErrSignerDifferentFromBidder, "Signer: %s, Bidder: %s", msg.GetSigner(), msg.GetBidder())
	}

	// Verify auction
	auction, found := k.GetActiveAuctionByID(ctx, msg.GetAuctionId())
	if !found {
		return &types.MsgSubmitBidResponse{}, sdkerrors.Wrapf(types.ErrAuctionNotFound, "Auction id: %d", msg.GetAuctionId())
	}

	// Verify auction coin type and bidder coin type are equal
	if auction.GetStartingTokensForSale().Denom != msg.GetSaleTokenMinimumAmount().Denom {
		return &types.MsgSubmitBidResponse{}, sdkerrors.Wrapf(types.ErrBidAuctionDenomMismatch, "Bid denom: %s, Auction denom: %s", msg.GetSaleTokenMinimumAmount().Denom, auction.GetStartingTokensForSale().Denom)
	}

	// Query our module address for funds
	totalSaleTokenBalance := k.bankKeeper.GetBalance(ctx, k.GetAuctionAccount(ctx).GetAddress(), auction.StartingTokensForSale.Denom)

	// Convert & standardize types for use below
	minimumSaleTokenPurchaseAmount := msg.SaleTokenMinimumAmount.Amount
	maxBidInUsomm := msg.MaxBidInUsomm.Amount

	// Calculate minimum purchase price
	// Note we round up, thus making the price more expensive to prevent this rounding from being exploited
	// TODO(pbal): consider adding minimum amount of usomm being bid as a global param
	minimumPurchasePriceInUsomm := sdk.NewInt(auction.CurrentUnitPriceInUsomm.Mul(minimumSaleTokenPurchaseAmount.ToDec()).Ceil().TruncateInt64())

	// Verify minimum price is <= bid, note this also checks the max bid is enough to purchase at least one sale token
	if minimumPurchasePriceInUsomm.GT(maxBidInUsomm) {
		return &types.MsgSubmitBidResponse{}, sdkerrors.Wrapf(types.ErrInsufficientBid, "minimum purchase price: %s, max bid: %s", minimumPurchasePriceInUsomm.String(), maxBidInUsomm.String())
	}

	// Start off fulfilled sale token amount at 0
	totalFulfilledSaleTokens := sdk.NewCoin(auction.StartingTokensForSale.GetDenom(), sdk.NewInt(0))

	// See how many whole tokens of base denom can be purchased
	// Round down, can't purchase fractional sale tokens
	saleTokensToPurchase := msg.MaxBidInUsomm.Amount.ToDec().Quo(auction.CurrentUnitPriceInUsomm).TruncateInt()

	// Figure out how much of bid we can fulfill
	if totalSaleTokenBalance.Amount.GTE(saleTokensToPurchase) {
		totalFulfilledSaleTokens.Amount = saleTokensToPurchase

	} else if totalSaleTokenBalance.Amount.GTE(msg.SaleTokenMinimumAmount.Amount) {
		totalFulfilledSaleTokens.Amount = totalSaleTokenBalance.Amount

	} else {
		return &types.MsgSubmitBidResponse{}, sdkerrors.Wrapf(types.ErrMinimumPurchaseAmountLargerThanTokensRemaining, "Minimum purchase: %s, amount remaining: %s", msg.SaleTokenMinimumAmount.String(), auction.RemainingTokensForSale.String())
	}

	// Round up to prevent exploitability; ensure you can't get more than you pay for
	usommAmount := sdk.NewInt(totalFulfilledSaleTokens.Amount.ToDec().Mul(auction.CurrentUnitPriceInUsomm).Ceil().TruncateInt64())
	totalUsommPaid := sdk.NewCoin(types.UsommDenom, usommAmount)

	newBidID := k.GetLastBidID(ctx) + 1
	bid := types.Bid{
		Id:                        newBidID,
		AuctionId:                 msg.GetAuctionId(),
		Bidder:                    msg.GetBidder(),
		MaxBidInUsomm:             msg.GetMaxBidInUsomm(),
		SaleTokenMinimumAmount:    msg.GetSaleTokenMinimumAmount(),
		TotalFulfilledSaleTokens:  totalFulfilledSaleTokens,
		SaleTokenUnitPriceInUsomm: auction.CurrentUnitPriceInUsomm,
		TotalUsommPaid:            totalUsommPaid,
	}

	err := bid.ValidateBasic()
	if err != nil {
		return &types.MsgSubmitBidResponse{}, err
	}

	// Transfer payment first
	if err := k.bankKeeper.SendCoinsFromAccountToModule(ctx, sdk.AccAddress(msg.GetBidder()), types.ModuleName, sdk.NewCoins(totalUsommPaid)); err != nil {
		return &types.MsgSubmitBidResponse{}, err
	}

	// Transfer purchase to bidder
	if err := k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, sdk.AccAddress(msg.GetBidder()), sdk.NewCoins(totalFulfilledSaleTokens)); err != nil {
		// TODO(pbal): Audit if we should panic here
		return &types.MsgSubmitBidResponse{}, err
	}

	// Update amount remaining in auction
	// Note we don't need to check negative here as we're checking totalFulfilledSaleTokens against the balance earlier
	auction.RemainingTokensForSale = totalSaleTokenBalance.Sub(totalFulfilledSaleTokens)
	k.setActiveAuction(ctx, auction)

	// Create bid in store
	k.setBid(ctx, bid)

	// Update latest bid id
	k.setLastBidID(ctx, newBidID)

	// Verify auction still has supply to see if we need to finish it
	if auction.RemainingTokensForSale.IsZero() {
		// Finish auction if so
		err := k.FinishAuction(ctx, &auction)

		// Since we use FinishAuction in EndBlocker, and the user can't & shouldn't do anything with this error, we panic
		// This will only occur if a module to module token transfer fails while finishing the auction
		if err != nil {
			panic(err)
		}
	}

	// Emit Event to signal bid was made
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeBid,
			sdk.NewAttribute(types.AttributeKeyAuctionID, fmt.Sprint(msg.GetAuctionId())),
			sdk.NewAttribute(types.AttributeKeyBidID, fmt.Sprint(newBidID)),
			sdk.NewAttribute(types.AttributeKeyBidder, msg.GetBidder()),
			sdk.NewAttribute(types.AttributeKeyMinimumAmount, msg.GetSaleTokenMinimumAmount().String()),
			sdk.NewAttribute(types.AttributeKeySigner, msg.GetSigner()),
			sdk.NewAttribute(types.AttributeKeyFulfilledPrice, auction.CurrentUnitPriceInUsomm.String()),
			sdk.NewAttribute(types.AttributeKeyTotalPayment, totalUsommPaid.String()),
			sdk.NewAttribute(types.AttributeKeyFulfilledAmount, totalFulfilledSaleTokens.String()),
		),
	)

	return &types.MsgSubmitBidResponse{Bid: &bid}, nil
}
