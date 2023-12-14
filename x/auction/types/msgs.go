package types

import (
	"strings"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/peggyjv/sommelier/v7/app/params"
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
		return errorsmod.Wrapf(ErrAuctionIDMustBeNonZero, "id: %d", m.AuctionId)
	}

	if m.MaxBidInUsomm.Denom != params.BaseCoinUnit {
		return errorsmod.Wrapf(ErrBidMustBeInUsomm, "bid: %s", m.MaxBidInUsomm.String())
	}

	if !m.MaxBidInUsomm.IsPositive() {
		return errorsmod.Wrapf(ErrBidAmountMustBePositive, "bid amount in usomm: %s", m.MaxBidInUsomm.String())
	}

	if !strings.HasPrefix(m.SaleTokenMinimumAmount.Denom, "gravity0x") {
		return errorsmod.Wrapf(ErrInvalidTokenBeingBidOn, "sale token: %s", m.SaleTokenMinimumAmount.String())
	}

	if !m.SaleTokenMinimumAmount.IsPositive() {
		return errorsmod.Wrapf(ErrMinimumAmountMustBePositive, "sale token amount: %s", m.SaleTokenMinimumAmount.String())
	}

	if _, err := sdk.AccAddressFromBech32(m.Signer); err != nil {
		return errorsmod.Wrap(sdkerrors.ErrInvalidAddress, err.Error())
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

// MustGetSigner returns the signer address (which is also the bidder)
func (m *MsgSubmitBidRequest) MustGetSigner() sdk.AccAddress {
	addr, err := sdk.AccAddressFromBech32(m.Signer)
	if err != nil {
		panic(err)
	}
	return addr
}
