package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/ethereum/go-ethereum/common"
)

var (
	_ sdk.Msg = &MsgSubmitReinvestRequest{}
)

const (
	TypeMsgSubmitReinvestRequest = "reinvest_submit"
)

//////////////////////////
// MsgSubmitReinvestRequest //
//////////////////////////

// NewMsgSubmitReinvestRequest return a new MsgSubmitReinvestRequest
func NewMsgSubmitReinvestRequest(body []byte, address common.Address, signer sdk.AccAddress) (*MsgSubmitReinvestRequest, error) {
	return &MsgSubmitReinvestRequest{
		Reinvestment: &Reinvestment{
			Body:    body,
			Address: address.String(),
		},
		Signer: signer.String(),
	}, nil
}

// Route implements sdk.Msg
func (m *MsgSubmitReinvestRequest) Route() string { return ModuleName }

// Type implements sdk.Msg
func (m *MsgSubmitReinvestRequest) Type() string { return TypeMsgSubmitReinvestRequest }

// ValidateBasic implements sdk.Msg
func (m *MsgSubmitReinvestRequest) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(m.Signer); err != nil {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, err.Error())
	}

	return nil
}

// GetSignBytes implements sdk.Msg
func (m *MsgSubmitReinvestRequest) GetSignBytes() []byte {
	panic("amino support disabled")
}

// GetSigners implements sdk.Msg
func (m *MsgSubmitReinvestRequest) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{m.MustGetSigner()}
}

// MustGetSigner returns the signer address
func (m *MsgSubmitReinvestRequest) MustGetSigner() sdk.AccAddress {
	addr, err := sdk.AccAddressFromBech32(m.Signer)
	if err != nil {
		panic(err)
	}
	return addr
}
