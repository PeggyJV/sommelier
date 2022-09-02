package types

import (
	fmt "fmt"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	gravitytypes "github.com/peggyjv/gravity-bridge/module/v2/x/gravity/types"
)

const UsommDenom = "usomm"

func (a *Auction) ValidateBasic() error {
	if a.Id == 0 {
		return fmt.Errorf("auction IDs must be non-zero")
	}

	if !a.StartingAmount.IsPositive() {
		return sdkerrors.Wrapf(ErrAuctionStartingAmountMustBePositve, "Starting amount: %s", a.StartingAmount.String())
	}

	if a.StartingAmount.Denom == "" {
		return sdkerrors.Wrapf(ErrAuctionDenomInvalid, "Starting denom: %s", a.StartingAmount.String())
	}

	if a.StartingAmount.Denom == UsommDenom {
		return sdkerrors.Wrapf(ErrCannotAuctionUsomm,"Starting denom is: %s", UsommDenom)
	}

	if a.StartBlock == 0 {
		return fmt.Errorf("start block must be non-zero")
	}

	if a.InitialDecreaseRate <= 0 || a.InitialDecreaseRate >= 1 {
		return sdkerrors.Wrapf(ErrInvalidInitialDecreaseRate, "Inital decrease rate %f", a.InitialDecreaseRate)
	}

	if a.CurrentDecreaseRate <= 0 || a.CurrentDecreaseRate >= 1 {
		return fmt.Errorf("current decrease rate must be a float less than or equal to one and greater than or equal to zero")
	}

	if a.BlockDecreaseInterval == 0 {
		return sdkerrors.Wrapf(ErrInvalidBlockDecreaeInterval, "Block Decrease interval: %d", a.BlockDecreaseInterval)
	}

	if !a.CurrentUnitPriceInUsomm.IsPositive() {
		return fmt.Errorf("current price must be positive")
	}

	if a.AmountRemaining.Denom == "" {
		return fmt.Errorf("amount remaining denom cannot be empty")
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
		return fmt.Errorf("bid IDs must be non-zero")
	}

	if b.AuctionId == 0 {
		return fmt.Errorf("auction IDs must be non-zero")
	}

	if b.Bidder == "" {
		return fmt.Errorf("bidder cannot be empty")
	}

	if _, err := sdk.AccAddressFromBech32(b.Bidder); err != nil {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, err.Error())
	}

	if !b.MaxBid.IsPositive() {
		return fmt.Errorf("bids must be a positive amount of %s", UsommDenom)
	}

	if !strings.HasPrefix(b.MinimumAmount.Denom, gravitytypes.GravityDenomPrefix) {
		return fmt.Errorf("bids may only be placed for gravity tokens")
	}

	if !b.MinimumAmount.IsPositive() {
		return fmt.Errorf("minimum amount must be a positive amount of auctioned coins")
	}

	if b.TotalFulfilledSaleTokenAmount.Amount.IsNegative() {
		return fmt.Errorf("total sale token fulfillment amount must be non-negative")
	}

	if !b.UnitPriceOfSaleTokenInUsomm.IsPositive() {
		return fmt.Errorf("unit price of sale tokens in usomm must be positive")
	}

	if b.TotalAmountPaidInUsomm.IsNegative() {
		return fmt.Errorf("total amount paid in usomm cannot be negative")
	}

	return nil
}

func (t *TokenPrice) ValidateBasic() error {
	if t.Denom == "" {
		return fmt.Errorf("denom cannot be empty")
	}

	if !t.UsdPrice.IsPositive() {
		return fmt.Errorf("price must be greater than 0")
	}

	if t.LastUpdatedBlock == 0 {
		return fmt.Errorf("last updated block must be greater than 0")
	}

	return nil
}

func (t *ProposedTokenPrice) ValidateBasic() error {
	if t.Denom == "" {
		return fmt.Errorf("denom cannot be empty")
	}

	if !t.UsdPrice.IsPositive() {
		return fmt.Errorf("price must be greater than 0")
	}

	return nil
}
