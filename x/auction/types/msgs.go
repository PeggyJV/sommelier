package types

import (
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var (
	_ sdk.Msg = &MsgSubmitBidRequest{}
)

const (
	TypeMsgSubmitBidRequest = "submit_bid"
)

/////////////////////////
// MsgSubmitBidRequest //
/////////////////////////

// NewMsgSubmitBidRequest return a new MsgSubmitBidRequest
func NewMsgSubmitBidRequest(auctionID uint32, maxBidInUsomm sdk.Coin, saleTokenMinimumAmount sdk.Coin, signer sdk.AccAddress) (*MsgSubmitBidRequest, error) {
	return &MsgSubmitBidRequest{
		AuctionId:              auctionID,
		MaxBidInUsomm:          maxBidInUsomm,
		SaleTokenMinimumAmount: saleTokenMinimumAmount,
		Bidder:                 signer.String(),
		Signer:                 signer.String(),
	}, nil
}

// Route implements sdk.Msg
func (m *MsgSubmitBidRequest) Route() string { return ModuleName }

// Type implements sdk.Msg
func (m *MsgSubmitBidRequest) Type() string { return TypeMsgSubmitBidRequest }

// ValidateBasic implements sdk.Msg
func (m *MsgSubmitBidRequest) ValidateBasic() error {
	if m.AuctionId == 0 {
		return sdkerrors.Wrapf(ErrAuctionIDMustBeNonZero, "id: %d", m.AuctionId)
	}

	if m.MaxBidInUsomm.Denom != UsommDenom {
		return sdkerrors.Wrapf(ErrBidMustBeInUsomm, "bid: %s", m.MaxBidInUsomm.String())
	}

	if !m.MaxBidInUsomm.IsPositive() {
		return sdkerrors.Wrapf(ErrBidAmountMustBePositive, "bid amount in usomm: %s", m.MaxBidInUsomm.String())
	}

	if !strings.HasPrefix(m.SaleTokenMinimumAmount.Denom, "gravity0x") {
		return sdkerrors.Wrapf(ErrInvalidTokenBeingBidOn, "sale token: %s", m.SaleTokenMinimumAmount.String())
	}

	if !m.SaleTokenMinimumAmount.IsPositive() {
		return sdkerrors.Wrapf(ErrMinimumAmountMustBePositive, "sale token amount: %s", m.SaleTokenMinimumAmount.String())
	}

	if _, err := sdk.AccAddressFromBech32(m.Bidder); err != nil {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, err.Error())
	}

	if m.Signer != m.Bidder {
		return sdkerrors.Wrapf(ErrSignerDifferentFromBidder, "signer: %s, bidder: %s", m.Signer, m.Bidder)
	}

	return nil
}

// GetSignBytes implements sdk.Msg
func (m *MsgSubmitBidRequest) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(m))
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
