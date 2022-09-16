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
		return sdkerrors.Wrapf(ErrAuctionIDMustBeNonZero, "id: %d", a.Id)
	}

	if !a.StartingTokensForSale.IsPositive() {
		return sdkerrors.Wrapf(ErrAuctionStartingAmountMustBePositve, "Starting tokens for sale: %s", a.StartingTokensForSale.String())
	}

	if a.StartingTokensForSale.Denom == UsommDenom {
		return sdkerrors.Wrapf(ErrCannotAuctionUsomm, "Starting denom tokens for sale: %s", UsommDenom)
	}

	if a.StartBlock == 0 {
		return sdkerrors.Wrapf(ErrInvalidStartBlock, "start block: %d", a.StartBlock)
	}

	if a.InitialPriceDecreaseRate.LTE(sdk.NewDec(0)) || a.InitialPriceDecreaseRate.GTE(sdk.NewDec(1)) {
		return sdkerrors.Wrapf(ErrInvalidInitialDecreaseRate, "Inital price decrease rate %s", a.InitialPriceDecreaseRate.String())
	}

	if a.CurrentPriceDecreaseRate.LTE(sdk.NewDec(0)) || a.CurrentPriceDecreaseRate.GTE(sdk.NewDec(1)) {
		return sdkerrors.Wrapf(ErrInvalidCurrentDecreaseRate, "Current price decrease rate %s", a.CurrentPriceDecreaseRate.String())
	}

	if a.PriceDecreaseBlockInterval == 0 {
		return sdkerrors.Wrapf(ErrInvalidBlockDecreaseInterval, "price decrease block interval: %d", a.PriceDecreaseBlockInterval)
	}

	if !a.InitialUnitPriceInUsomm.IsPositive() {
		return sdkerrors.Wrapf(ErrPriceMustBePositive, "initial unit price in usomm: %s", a.InitialUnitPriceInUsomm.String())
	}

	if !a.CurrentUnitPriceInUsomm.IsPositive() {
		return sdkerrors.Wrapf(ErrPriceMustBePositive, "current unit price in usomm: %s", a.CurrentUnitPriceInUsomm.String())
	}

	if a.FundingModuleAccount == "" {
		return sdkerrors.Wrapf(ErrUnauthorizedFundingModule, "funding module account: %s", a.FundingModuleAccount)
	}

	if a.ProceedsModuleAccount == "" {
		return sdkerrors.Wrapf(ErrUnauthorizedFundingModule, "proceeds module account: %s", a.ProceedsModuleAccount)
	}

	return nil
}

func (b *Bid) ValidateBasic() error {
	if b.Id == 0 {
		return sdkerrors.Wrapf(ErrBidIDMustBeNonZero, "id: %d", b.Id)
	}

	if b.AuctionId == 0 {
		return sdkerrors.Wrapf(ErrAuctionIDMustBeNonZero, "id: %d", b.AuctionId)
	}

	if b.Bidder == "" {
		return sdkerrors.Wrapf(ErrAddressExpected, "bidder: %s", b.Bidder)
	}

	if _, err := sdk.AccAddressFromBech32(b.Bidder); err != nil {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, err.Error())
	}

	if !b.MaxBidInUsomm.IsPositive() {
		return sdkerrors.Wrapf(ErrBidAmountMustBePositive, "bid amount in usomm: %s", b.MaxBidInUsomm.String())
	}

	if b.MaxBidInUsomm.Denom != UsommDenom {
		return sdkerrors.Wrapf(ErrBidMustBeInUsomm, "bid: %s", b.MaxBidInUsomm.String())
	}

	if !strings.HasPrefix(b.SaleTokenMinimumAmount.Denom, gravitytypes.GravityDenomPrefix) {
		return sdkerrors.Wrapf(ErrInvalidTokenBeingBidOn, "sale token: %s", b.SaleTokenMinimumAmount.String())
	}

	if !b.SaleTokenMinimumAmount.IsPositive() {
		return sdkerrors.Wrapf(ErrMinimumAmountMustBePositive, "sale token amount: %s", b.SaleTokenMinimumAmount.String())
	}

	if !b.SaleTokenUnitPriceInUsomm.IsPositive() {
		return sdkerrors.Wrapf(ErrBidUnitPriceInUsommMustBePositive, "sale token unit price: %s", b.SaleTokenUnitPriceInUsomm.String())
	}

	if b.TotalUsommPaid.Denom != UsommDenom {
		return sdkerrors.Wrapf(ErrBidMustBeInUsomm, "payment denom: %s", b.TotalUsommPaid.Denom)
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

	if !strings.HasPrefix(t.Denom, gravitytypes.GravityDenomPrefix) && t.Denom != UsommDenom {
		return sdkerrors.Wrapf(ErrInvalidTokenPriceDenom, "denom: %s", t.Denom)
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

	if !strings.HasPrefix(t.Denom, gravitytypes.GravityDenomPrefix) && t.Denom != UsommDenom {
		return sdkerrors.Wrapf(ErrInvalidTokenPriceDenom, "denom: %s", t.Denom)
	}

	return nil
}
