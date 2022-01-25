package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/ethereum/go-ethereum/common"
)

var (
	_ sdk.Msg = &MsgSubmitCorkRequest{}
)

const (
	TypeMsgSubmitCorkRequest = "cork_submit"
)

//////////////////////////
// MsgSubmitCorkRequest //
//////////////////////////

// NewMsgSubmitCorkRequest return a new MsgSubmitCorkRequest
func NewMsgSubmitCorkRequest(body []byte, address common.Address, signer sdk.AccAddress) (*MsgSubmitCorkRequest, error) {
	return &MsgSubmitCorkRequest{
		Cork: &Cork{
			Body:    body,
			Address: address.String(),
		},
		Signer: signer.String(),
	}, nil
}

// Route implements sdk.Msg
func (m *MsgSubmitCorkRequest) Route() string { return ModuleName }

// Type implements sdk.Msg
func (m *MsgSubmitCorkRequest) Type() string { return TypeMsgSubmitCorkRequest }

// ValidateBasic implements sdk.Msg
func (m *MsgSubmitCorkRequest) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(m.Signer); err != nil {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, err.Error())
	}

	return nil
}

// GetSignBytes implements sdk.Msg
func (m *MsgSubmitCorkRequest) GetSignBytes() []byte {
	panic("amino support disabled")
}

// GetSigners implements sdk.Msg
func (m *MsgSubmitCorkRequest) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{m.MustGetSigner()}
}

// MustGetSigner returns the signer address
func (m *MsgSubmitCorkRequest) MustGetSigner() sdk.AccAddress {
	addr, err := sdk.AccAddressFromBech32(m.Signer)
	if err != nil {
		panic(err)
	}
	return addr
}
