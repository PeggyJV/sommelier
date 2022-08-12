package types

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var (
	_ sdk.Msg = &MsgSubmitBidRequest{}
)

const (
	TypeMsgSubmitBidRequest = "bid_submit"
)

/////////////////////////
// MsgSubmitBidRequest //
/////////////////////////

// NewMsgSubmitBidRequest return a new MsgSubmitBidRequest
func NewMsgSubmitBidRequest(body []byte, auctionId uint32, maxBid sdk.Coin, minimumAmount sdk.Coin, signer sdk.AccAddress) (*MsgSubmitBidRequest, error) {
	return &MsgSubmitBidRequest{
		AuctionId:     auctionId,
		MaxBid:        maxBid,
		MinimumAmount: minimumAmount,
		Bidder:        signer.String(),
		Signer:        signer.String(),
	}, nil
}

// Route implements sdk.Msg
func (m *MsgSubmitBidRequest) Route() string { return ModuleName }

// Type implements sdk.Msg
func (m *MsgSubmitBidRequest) Type() string { return TypeMsgSubmitBidRequest }

// ValidateBasic implements sdk.Msg
func (m *MsgSubmitBidRequest) ValidateBasic() error {
	if m.AuctionId == 0 {
		return fmt.Errorf("auction IDs must be non-zero")
	}

	if !m.MaxBid.IsPositive() {
		return fmt.Errorf("bids must be a positive amount of SOMM")
	}

	if !m.MinimumAmount.IsPositive() {
		return fmt.Errorf("minimum amount must be a positive amount of auctioned coins")
	}

	// TODO(bolten): is it possible to check the denom correctly here?

	if _, err := sdk.AccAddressFromBech32(m.Bidder); err != nil {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, err.Error())
	}

	if _, err := sdk.AccAddressFromBech32(m.Signer); err != nil {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, err.Error())
	}

	return nil
}

// GetSignBytes implements sdk.Msg
func (m *MsgSubmitBidRequest) GetSignBytes() []byte {
	panic("amino support disabled")
}

// GetSigners implements sdk.Msg
func (m *MsgSubmitBidRequest) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{m.MustGetSigner()}
}

// MustGetSigner returns the signer address
func (m *MsgSubmitBidRequest) MustGetSigner() sdk.AccAddress {
	addr, err := sdk.AccAddressFromBech32(m.Signer)
	if err != nil {
		panic(err)
	}
	return addr
}
