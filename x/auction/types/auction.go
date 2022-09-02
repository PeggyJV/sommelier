package types

import (
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	gravitytypes "github.com/peggyjv/gravity-bridge/module/v2/x/gravity/types"
)

const UsommDenom = "usomm"

func (a *Auction) ValidateBasic() error {
	if a.Id == 0 {
		return sdkerrors.Wrapf(ErrAuctionIdMustBeNonZero, "id: %d", a.Id)
	}

	if !a.StartingAmount.IsPositive() {
		return sdkerrors.Wrapf(ErrAuctionStartingAmountMustBePositve, "Starting amount: %s", a.StartingAmount.String())
	}

	if a.StartingAmount.Denom == "" {
		return sdkerrors.Wrapf(ErrAuctionDenomInvalid, "Starting denom: %s", a.StartingAmount.String())
	}

	if a.StartingAmount.Denom == UsommDenom {
		return sdkerrors.Wrapf(ErrCannotAuctionUsomm, "Starting denom is: %s", UsommDenom)
	}

	if a.StartBlock == 0 {
		return sdkerrors.Wrapf(ErrInvalidStartBlock, "block: %d", a.StartBlock)
	}

	if a.InitialDecreaseRate <= 0 || a.InitialDecreaseRate >= 1 {
		return sdkerrors.Wrapf(ErrInvalidInitialDecreaseRate, "Inital decrease rate %f", a.InitialDecreaseRate)
	}

	if a.CurrentDecreaseRate <= 0 || a.CurrentDecreaseRate >= 1 {
		return sdkerrors.Wrapf(ErrInvalidCurrentDecreaseRate, "Current decrease rate %f", a.CurrentDecreaseRate)
	}

	if a.BlockDecreaseInterval == 0 {
		return sdkerrors.Wrapf(ErrInvalidBlockDecreaeInterval, "Block Decrease interval: %d", a.BlockDecreaseInterval)
	}

	if !a.CurrentUnitPriceInUsomm.IsPositive() {
		return sdkerrors.Wrapf(ErrPriceMustBePositive, "current unit price: %s", a.CurrentUnitPriceInUsomm.String())
	}

	if a.AmountRemaining.Denom == "" {
		return sdkerrors.Wrapf(ErrDenomCannotBeEmpty, "amount remaining denom: %s", a.AmountRemaining.String())
	}

	if a.FundingModuleAccount == "" {
		return sdkerrors.Wrapf(ErrUnauthorizedFundingModule, "Account: %s", a.FundingModuleAccount)
	}

	if a.ProceedsModuleAccount == "" {
		return sdkerrors.Wrapf(ErrUnauthorizedFundingModule, "Account: %s", a.ProceedsModuleAccount)
	}

	return nil
}

func (b *Bid) ValidateBasic() error {
	if b.Id == 0 {
		return sdkerrors.Wrapf(ErrBidIdMustBeNonZero, "id: %d", b.Id)
	}

	if b.AuctionId == 0 {
		return sdkerrors.Wrapf(ErrAuctionIdMustBeNonZero, "id: %d", b.AuctionId)
	}

	if b.Bidder == "" {
		return sdkerrors.Wrapf(ErrAddressExpected, "bidder: %s", b.Bidder)
	}

	if _, err := sdk.AccAddressFromBech32(b.Bidder); err != nil {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, err.Error())
	}

	if !b.MaxBid.IsPositive() {
		return sdkerrors.Wrapf(ErrBidIdAmountMustBePositive, "bid amount: %s", b.MaxBid.String())
	}

	if b.MaxBid.Denom != UsommDenom {
		return sdkerrors.Wrapf(ErrBidMustBeInUsomm, "bid: %s", b.MaxBid.String())
	}

	if !strings.HasPrefix(b.MinimumAmount.Denom, gravitytypes.GravityDenomPrefix) {
		return sdkerrors.Wrapf(ErrInvalidTokenBeingBidOn, "token: %s", b.MinimumAmount)
	}

	if !b.MinimumAmount.IsPositive() {
		return sdkerrors.Wrapf(ErrMinimumAmountMustBePositive, "amount: %s", b.MinimumAmount.String())
	}

	if b.TotalFulfilledSaleTokenAmount.Amount.IsNegative() {
		return sdkerrors.Wrapf(ErrBidFulfilledSaleTokenAmountMustBeNonNegative, "amount: %s", b.TotalFulfilledSaleTokenAmount.String())
	}

	if !b.UnitPriceOfSaleTokenInUsomm.IsPositive() {
		return sdkerrors.Wrapf(ErrBidUnitPriceInUsommMustBePositive, "unit price: %s", b.UnitPriceOfSaleTokenInUsomm.String())
	}

	if b.TotalAmountPaidInUsomm.IsNegative() {
		return sdkerrors.Wrapf(ErrBidPaymentCannotBeNegative, "payment: %s", b.TotalAmountPaidInUsomm.String())
	}

	return nil
}

func (t *TokenPrice) ValidateBasic() error {
	if t.Denom == "" {
		return sdkerrors.Wrapf(ErrDenomCannotBeEmpty, "price denom: %s", t.Denom)
	}

	if !t.UsdPrice.IsPositive() {
		return sdkerrors.Wrapf(ErrPriceMustBePositive, "usd price: %s", t.UsdPrice.String())
	}

	if t.LastUpdatedBlock == 0 {
		return sdkerrors.Wrapf(ErrInvalidLastUpdatedBlock, "block: %d", t.LastUpdatedBlock)
	}

	return nil
}

func (t *ProposedTokenPrice) ValidateBasic() error {
	if t.Denom == "" {
		return sdkerrors.Wrapf(ErrDenomCannotBeEmpty, "price denom: %s", t.Denom)
	}

	if !t.UsdPrice.IsPositive() {
		return sdkerrors.Wrapf(ErrPriceMustBePositive, "usd price: %s", t.UsdPrice.String())
	}

	return nil
}
