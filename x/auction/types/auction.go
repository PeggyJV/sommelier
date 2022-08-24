package types

import (
	fmt "fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

func (a *Auction) Equals(other Auction) bool {
	if a.Id != other.Id {
		return false
	}

	if !a.StartingAmount.IsEqual(other.StartingAmount) {
		return false
	}
	
	if a.StartBlock != other.StartBlock {
		return false
	}

	if a.EndBlock != other.EndBlock {
		return false
	}

	if a.InitialDecreaseRate != other.InitialDecreaseRate {
		return false
	}

	if a.CurrentDecreaseRate != other.CurrentDecreaseRate {
		return false
	}

	if a.BlockDecreaseInterval != other.BlockDecreaseInterval {
		return false
	}

	if !a.CurrentPrice.IsEqual(other.CurrentPrice) {
		return false
	}

	if !a.AmountRemaining.IsEqual(other.AmountRemaining) {
		return false
	}

	return true
}

func (a *Auction) ValidateBasic() error {
	if a.Id == 0 {
		return fmt.Errorf("auction IDs must be non-zero")
	}

	if !a.StartingAmount.IsPositive() {
		return fmt.Errorf("minimum amount must be a positive amount of coins")
	}

	if a.StartBlock <= 0 {
		return fmt.Errorf("start block must be a positive")
	}

	if a.EndBlock <= 0 {
		return fmt.Errorf("end block must be a positive")
	}

	if a.StartBlock > a.EndBlock {
		return fmt.Errorf("start block must be smaller or equal to end block")
	}

	if a.InitialDecreaseRate <= 0 || a.InitialDecreaseRate >= 1 {
		return fmt.Errorf("initial decrease rate must be a float less than one and greater than zero")
	}

	if a.CurrentDecreaseRate < 0 || a.CurrentDecreaseRate > 1 {
		return fmt.Errorf("current decrease rate must be a float less than or equal to one and greater than or equal to zero")
	}

	return nil
}

func (b *Bid) Equals(other Bid) bool {
	if b.Id != other.Id {
		return false
	}

	if b.AuctionId != other.AuctionId {
		return false
	}

	if !b.MaxBid.IsEqual(other.MaxBid) {
		return false
	}

	if !b.MinimumAmount.IsEqual(other.MinimumAmount) {
		return false
	}

	if b.Bidder != other.Bidder {
		return false
	}

	if !b.FulfilledAmount.IsEqual(other.FulfilledAmount) {
		return false
	}

	if !b.FulfillmentPrice.IsEqual(other.FulfillmentPrice) {
		return false
	}

	return true
}

func (b *Bid) ValidateBasic() error {
	if b.Id == 0 {
		return fmt.Errorf("bid IDs must be non-zero")
	}

	if b.AuctionId == 0 {
		return fmt.Errorf("auction IDs must be non-zero")
	}

	if !b.MaxBid.IsPositive() {
		return fmt.Errorf("bids must be a positive amount of SOMM")
	}

	if !b.MinimumAmount.IsPositive() {
		return fmt.Errorf("minimum amount must be a positive amount of auctioned coins")
	}

	// TODO: fulfilled bid updates

	// TODO(bolten): is it possible to check the denom correctly here?

	if _, err := sdk.AccAddressFromBech32(b.Bidder); err != nil {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, err.Error())
	}

	return nil
}

func (t *TokenPrice) Equals(other TokenPrice) bool {
	if t.Denom != other.Denom {
		return false
	}

	if t.UsdPrice != other.UsdPrice {
		return false
	}

	return true
}

func (t *TokenPrice) ValidateBasic() error {
	if t.Denom == "" {
		return fmt.Errorf("denom must be a non empty string")
	}

	if t.UsdPrice.IsNegative() {
		return fmt.Errorf("price must be greater than or equal to 0")
	}

	return nil
}
