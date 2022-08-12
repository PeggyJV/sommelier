package types

import (
	fmt "fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

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

	// TODO(bolten): is it possible to check the denom correctly here?

	if _, err := sdk.AccAddressFromBech32(b.Bidder); err != nil {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, err.Error())
	}

	return nil
}
