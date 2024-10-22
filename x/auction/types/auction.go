package types

import (
	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/peggyjv/sommelier/v8/app/params"
)

func (a *Auction) ValidateBasic() error {
	if a.Id == 0 {
		return errorsmod.Wrapf(ErrAuctionIDMustBeNonZero, "id: %d", a.Id)
	}

	if !a.StartingTokensForSale.IsValid() || !a.StartingTokensForSale.IsPositive() {
		return errorsmod.Wrapf(ErrAuctionStartingAmountMustBePositve, "Starting tokens for sale: %s", a.StartingTokensForSale.String())
	}

	if a.StartingTokensForSale.Denom == params.BaseCoinUnit {
		return errorsmod.Wrapf(ErrCannotAuctionUsomm, "Starting denom tokens for sale: %s", params.BaseCoinUnit)
	}

	if a.StartBlock == 0 {
		return errorsmod.Wrapf(ErrInvalidStartBlock, "start block: %d", a.StartBlock)
	}

	if a.InitialPriceDecreaseRate.LTE(sdk.NewDec(0)) || a.InitialPriceDecreaseRate.GTE(sdk.NewDec(1)) {
		return errorsmod.Wrapf(ErrInvalidInitialDecreaseRate, "Initial price decrease rate %s", a.InitialPriceDecreaseRate.String())
	}

	if a.CurrentPriceDecreaseRate.LTE(sdk.NewDec(0)) || a.CurrentPriceDecreaseRate.GTE(sdk.NewDec(1)) {
		return errorsmod.Wrapf(ErrInvalidCurrentDecreaseRate, "Current price decrease rate %s", a.CurrentPriceDecreaseRate.String())
	}

	if a.PriceDecreaseBlockInterval == 0 {
		return errorsmod.Wrapf(ErrInvalidBlockDecreaseInterval, "price decrease block interval: %d", a.PriceDecreaseBlockInterval)
	}

	if !a.InitialUnitPriceInUsomm.IsPositive() {
		return errorsmod.Wrapf(ErrPriceMustBePositive, "initial unit price in usomm: %s", a.InitialUnitPriceInUsomm.String())
	}

	if !a.CurrentUnitPriceInUsomm.IsPositive() {
		return errorsmod.Wrapf(ErrPriceMustBePositive, "current unit price in usomm: %s", a.CurrentUnitPriceInUsomm.String())
	}

	if a.FundingModuleAccount == "" {
		return errorsmod.Wrapf(ErrUnauthorizedFundingModule, "funding module account: %s", a.FundingModuleAccount)
	}

	if a.ProceedsModuleAccount == "" {
		return errorsmod.Wrapf(ErrUnauthorizedFundingModule, "proceeds module account: %s", a.ProceedsModuleAccount)
	}

	return nil
}

func (b *Bid) ValidateBasic() error {
	if b.Id == 0 {
		return errorsmod.Wrapf(ErrBidIDMustBeNonZero, "id: %d", b.Id)
	}

	if b.AuctionId == 0 {
		return errorsmod.Wrapf(ErrAuctionIDMustBeNonZero, "id: %d", b.AuctionId)
	}

	if b.Bidder == "" {
		return errorsmod.Wrapf(ErrAddressExpected, "bidder: %s", b.Bidder)
	}

	if _, err := sdk.AccAddressFromBech32(b.Bidder); err != nil {
		return errorsmod.Wrap(sdkerrors.ErrInvalidAddress, err.Error())
	}

	if !b.MaxBidInUsomm.IsValid() || !b.MaxBidInUsomm.IsPositive() {
		return errorsmod.Wrapf(ErrBidAmountMustBePositive, "bid amount in usomm: %s", b.MaxBidInUsomm.String())
	}

	if b.MaxBidInUsomm.Denom != params.BaseCoinUnit {
		return errorsmod.Wrapf(ErrBidMustBeInUsomm, "bid: %s", b.MaxBidInUsomm.String())
	}

	if !b.SaleTokenMinimumAmount.IsValid() || !b.SaleTokenMinimumAmount.IsPositive() {
		return errorsmod.Wrapf(ErrMinimumAmountMustBePositive, "sale token amount: %s", b.SaleTokenMinimumAmount.String())
	}

	if !b.TotalFulfilledSaleTokens.IsValid() {
		return errorsmod.Wrapf(ErrBidFulfilledSaleTokenAmountMustBeNonNegative, "fulfilled sale token amount: %s", b.TotalFulfilledSaleTokens.String())
	}

	if !b.SaleTokenUnitPriceInUsomm.IsPositive() {
		return errorsmod.Wrapf(ErrBidUnitPriceInUsommMustBePositive, "sale token unit price: %s", b.SaleTokenUnitPriceInUsomm.String())
	}

	if b.TotalUsommPaid.IsNegative() {
		return errorsmod.Wrapf(ErrBidPaymentCannotBeNegative, "payment in usomm: %s", b.TotalUsommPaid.String())
	}

	if b.TotalUsommPaid.Denom != params.BaseCoinUnit {
		return errorsmod.Wrapf(ErrBidMustBeInUsomm, "payment denom: %s", b.TotalUsommPaid.Denom)
	}

	return nil
}

func (t *TokenPrice) ValidateBasic() error {
	if t.Denom == "" {
		return errorsmod.Wrapf(ErrDenomCannotBeEmpty, "price denom: %s", t.Denom)
	}

	if !t.UsdPrice.IsPositive() {
		return errorsmod.Wrapf(ErrPriceMustBePositive, "usd price: %s", t.UsdPrice.String())
	}

	if t.LastUpdatedBlock == 0 {
		return errorsmod.Wrapf(ErrInvalidLastUpdatedBlock, "block: %d", t.LastUpdatedBlock)
	}

	if t.Exponent > 18 {
		return errorsmod.Wrapf(ErrTokenPriceExponentTooHigh, "exponent: %d", t.Exponent)
	}

	return nil
}

func (t *ProposedTokenPrice) ValidateBasic() error {
	if t.Denom == "" {
		return errorsmod.Wrapf(ErrDenomCannotBeEmpty, "price denom: %s", t.Denom)
	}

	if !t.UsdPrice.IsPositive() {
		return errorsmod.Wrapf(ErrPriceMustBePositive, "usd price: %s", t.UsdPrice.String())
	}

	if t.Exponent > 18 {
		return errorsmod.Wrapf(ErrTokenPriceExponentTooHigh, "exponent: %d", t.Exponent)
	}

	return nil
}
