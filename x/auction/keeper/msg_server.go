package keeper

import (
	"context"
	"fmt"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/peggyjv/sommelier/v9/app/params"
	"github.com/peggyjv/sommelier/v9/x/auction/types"
)

var _ types.MsgServer = Keeper{}

// SubmitBid implements types.MsgServer
func (k Keeper) SubmitBid(c context.Context, msg *types.MsgSubmitBidRequest) (*types.MsgSubmitBidResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	signer := msg.MustGetSigner()

	// Verify auction
	auction, found := k.GetActiveAuctionByID(ctx, msg.GetAuctionId())
	if !found {
		return &types.MsgSubmitBidResponse{}, errorsmod.Wrapf(types.ErrAuctionNotFound, "Auction id: %d", msg.GetAuctionId())
	}

	// Verify auction coin type and bidder coin type are equal
	if auction.GetStartingTokensForSale().Denom != msg.GetSaleTokenMinimumAmount().Denom {
		return &types.MsgSubmitBidResponse{}, errorsmod.Wrapf(types.ErrBidAuctionDenomMismatch, "Bid denom: %s, Auction denom: %s", msg.GetSaleTokenMinimumAmount().Denom, auction.GetStartingTokensForSale().Denom)
	}

	// Query our module address for funds
	totalSaleTokenBalance := k.bankKeeper.GetBalance(ctx, k.GetAuctionAccount(ctx).GetAddress(), auction.StartingTokensForSale.Denom)

	// Convert & standardize types for use below
	minimumSaleTokenPurchaseAmount := msg.SaleTokenMinimumAmount.Amount
	maxBidInUsomm := msg.MaxBidInUsomm.Amount

	// To prevent spamming of many small bids, check that minimum bid amount is satisfied (unless amount left in auction is < minimum bid req)**
	minUsommBid := sdk.NewIntFromUint64(k.GetParamSet(ctx).MinimumBidInUsomm)
	saleTokenBalanceValueInUsommRemaining := sdk.NewDecFromInt(totalSaleTokenBalance.Amount).Mul(auction.CurrentUnitPriceInUsomm)

	// **If remaining amount in auction is LT minUsommBid param, update minUsommBid to smallest possible value left in auction to prevent spamming in this edge case
	if saleTokenBalanceValueInUsommRemaining.LT(sdk.NewDecFromInt(minUsommBid)) {
		minUsommBid = sdk.NewInt(saleTokenBalanceValueInUsommRemaining.TruncateInt64())
	}

	if maxBidInUsomm.LT(minUsommBid) {
		return &types.MsgSubmitBidResponse{}, errorsmod.Wrapf(types.ErrBidAmountIsTooSmall, "bid amount: %s, minimum amount in usomm: %s", maxBidInUsomm.String(), minUsommBid.String())
	}

	// Calculate minimum purchase price
	// Note we round up, thus making the price more expensive to prevent this rounding from being exploited
	minimumPurchasePriceInUsomm := sdk.NewInt(auction.CurrentUnitPriceInUsomm.Mul(sdk.NewDecFromInt(minimumSaleTokenPurchaseAmount)).Ceil().TruncateInt64())

	// Verify minimum price is <= bid, note this also checks the max bid is enough to purchase at least one sale token
	if minimumPurchasePriceInUsomm.GT(maxBidInUsomm) {
		return &types.MsgSubmitBidResponse{}, errorsmod.Wrapf(types.ErrInsufficientBid, "minimum purchase price: %s, max bid: %s", minimumPurchasePriceInUsomm.String(), maxBidInUsomm.String())
	}

	// Start off fulfilled sale token amount at 0
	totalFulfilledSaleTokens := sdk.NewCoin(auction.StartingTokensForSale.GetDenom(), sdk.NewInt(0))

	// See how many whole tokens of base denom can be purchased
	// Round down, can't purchase fractional sale tokens

	saleTokensToPurchase := sdk.NewDecFromInt(msg.MaxBidInUsomm.Amount).Quo(auction.CurrentUnitPriceInUsomm).TruncateInt()

	// Figure out how much of bid we can fulfill
	if totalSaleTokenBalance.Amount.GTE(saleTokensToPurchase) {
		totalFulfilledSaleTokens.Amount = saleTokensToPurchase

	} else if totalSaleTokenBalance.Amount.GTE(msg.SaleTokenMinimumAmount.Amount) {
		totalFulfilledSaleTokens.Amount = totalSaleTokenBalance.Amount

	} else {
		return &types.MsgSubmitBidResponse{}, errorsmod.Wrapf(types.ErrMinimumPurchaseAmountLargerThanTokensRemaining, "Minimum purchase: %s, amount remaining: %s", minimumSaleTokenPurchaseAmount.String(), auction.RemainingTokensForSale.String())
	}

	// Round up to prevent exploitability; ensure you can't get more than you pay for
	usommAmount := sdk.NewInt(sdk.NewDecFromInt(totalFulfilledSaleTokens.Amount).Mul(auction.CurrentUnitPriceInUsomm).Ceil().TruncateInt64())
	totalUsommPaid := sdk.NewCoin(params.BaseCoinUnit, usommAmount)

	newBidID := k.GetLastBidID(ctx) + 1
	bid := types.Bid{
		Id:                        newBidID,
		AuctionId:                 msg.GetAuctionId(),
		Bidder:                    signer.String(),
		MaxBidInUsomm:             msg.GetMaxBidInUsomm(),
		SaleTokenMinimumAmount:    msg.GetSaleTokenMinimumAmount(),
		TotalFulfilledSaleTokens:  totalFulfilledSaleTokens,
		SaleTokenUnitPriceInUsomm: auction.CurrentUnitPriceInUsomm,
		TotalUsommPaid:            totalUsommPaid,
		BlockHeight:               uint64(ctx.BlockHeight()),
	}

	err := bid.ValidateBasic()
	if err != nil {
		return &types.MsgSubmitBidResponse{}, err
	}

	// Transfer payment first
	if err := k.bankKeeper.SendCoinsFromAccountToModule(ctx, signer, types.ModuleName, sdk.NewCoins(totalUsommPaid)); err != nil {
		return &types.MsgSubmitBidResponse{}, err
	}

	// Transfer purchase to bidder
	if err := k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, signer, sdk.NewCoins(totalFulfilledSaleTokens)); err != nil {
		panic(fmt.Sprintf("funds taken from bidder but purchased tokens not transferred for bid: %s,\n err: %s", bid.String(), err.Error()))
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
			sdk.NewAttribute(types.AttributeKeyBidder, signer.String()),
			sdk.NewAttribute(types.AttributeKeyMinimumAmount, msg.GetSaleTokenMinimumAmount().String()),
			sdk.NewAttribute(types.AttributeKeyFulfilledPrice, auction.CurrentUnitPriceInUsomm.String()),
			sdk.NewAttribute(types.AttributeKeyTotalPayment, totalUsommPaid.String()),
			sdk.NewAttribute(types.AttributeKeyFulfilledAmount, totalFulfilledSaleTokens.String()),
		),
	)

	return &types.MsgSubmitBidResponse{Bid: &bid}, nil
}
