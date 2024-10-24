package v1

import (
	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/ethereum/go-ethereum/common"
	corktypes "github.com/peggyjv/sommelier/v8/x/cork/types"
)

var (
	_ sdk.Msg = &MsgSubmitCorkRequest{}
	_ sdk.Msg = &MsgScheduleCorkRequest{}
)

const (
	TypeMsgSubmitCorkRequest   = "cork_submit"
	TypeMsgScheduleCorkRequest = "cork_schedule"
)

//////////////////////////
// MsgSubmitCorkRequest //
//////////////////////////

// NewMsgSubmitCorkRequest return a new MsgSubmitCorkRequest
func NewMsgSubmitCorkRequest(body []byte, address common.Address, signer sdk.AccAddress) (*MsgSubmitCorkRequest, error) {
	return &MsgSubmitCorkRequest{
		Cork: &Cork{
			EncodedContractCall:   body,
			TargetContractAddress: address.String(),
		},
		Signer: signer.String(),
	}, nil
}

// Route implements sdk.Msg
func (m *MsgSubmitCorkRequest) Route() string { return corktypes.ModuleName }

// Type implements sdk.Msg
func (m *MsgSubmitCorkRequest) Type() string { return TypeMsgSubmitCorkRequest }

// ValidateBasic implements sdk.Msg
func (m *MsgSubmitCorkRequest) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(m.Signer); err != nil {
		return errorsmod.Wrap(sdkerrors.ErrInvalidAddress, err.Error())
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

//////////////////////////
// MsgScheduleCorkRequest //
//////////////////////////

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
func (m *MsgScheduleCorkRequest) Route() string { return corktypes.ModuleName }

// Type implements sdk.Msg
func (m *MsgScheduleCorkRequest) Type() string { return TypeMsgScheduleCorkRequest }

// ValidateBasic implements sdk.Msg
func (m *MsgScheduleCorkRequest) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(m.Signer); err != nil {
		return errorsmod.Wrap(sdkerrors.ErrInvalidAddress, err.Error())
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
