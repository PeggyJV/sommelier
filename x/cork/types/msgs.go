package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/ethereum/go-ethereum/common"
)

var (
	_ sdk.Msg = &MsgScheduleCorkRequest{}
)

const (
	TypeMsgSubmitCorkRequest   = "cork_submit"
	TypeMsgScheduleCorkRequest = "cork_schedule"
)

////////////////////////////
// MsgScheduleCorkRequest //
////////////////////////////

// NewMsgScheduleCorkRequest return a new MsgScheduleCorkRequest
func NewMsgScheduleCorkRequest(body []byte, address common.Address, blockHeight uint64, signer sdk.AccAddress) (*MsgScheduleCorkRequest, error) {
	return &MsgScheduleCorkRequest{
		Cork: &Cork{
			EncodedContractCall:   body,
			TargetContractAddress: address.String(),
		},
		BlockHeight: blockHeight,
		Signer:      signer.String(),
	}, nil
}

// Route implements sdk.Msg
func (m *MsgScheduleCorkRequest) Route() string { return ModuleName }

// Type implements sdk.Msg
func (m *MsgScheduleCorkRequest) Type() string { return TypeMsgScheduleCorkRequest }

// ValidateBasic implements sdk.Msg
func (m *MsgScheduleCorkRequest) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(m.Signer); err != nil {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, err.Error())
	}

	return m.Cork.ValidateBasic()
}

// GetSignBytes implements sdk.Msg
func (m *MsgScheduleCorkRequest) GetSignBytes() []byte {
	panic("amino support disabled")
}

// GetSigners implements sdk.Msg
func (m *MsgScheduleCorkRequest) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{m.MustGetSigner()}
}

// MustGetSigner returns the signer address
func (m *MsgScheduleCorkRequest) MustGetSigner() sdk.AccAddress {
	addr, err := sdk.AccAddressFromBech32(m.Signer)
	if err != nil {
		panic(err)
	}
	return addr
}
