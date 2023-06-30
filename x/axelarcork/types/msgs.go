package types

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/ethereum/go-ethereum/common"
)

var (
	_ sdk.Msg = &MsgScheduleAxelarCorkRequest{}
)

const (
	TypeMsgScheduleCorkRequest = "cork_schedule"
	TypeMsgRelayCorkRequest    = "cork_relay"
	TypeMsgBumpCorkGasRequest  = "cork_bump_gas"
)

////////////////////////////
// MsgScheduleAxelarCorkRequest //
////////////////////////////

// NewMsgScheduleCorkRequest return a new MsgScheduleAxelarCorkRequest
func NewMsgScheduleCorkRequest(body []byte, address common.Address, blockHeight uint64, signer sdk.AccAddress) (*MsgScheduleAxelarCorkRequest, error) {
	return &MsgScheduleAxelarCorkRequest{
		Cork: &Cork{
			EncodedContractCall:   body,
			TargetContractAddress: address.String(),
		},
		BlockHeight: blockHeight,
		Signer:      signer.String(),
	}, nil
}

// Route implements sdk.Msg
func (m *MsgScheduleAxelarCorkRequest) Route() string { return ModuleName }

// Type implements sdk.Msg
func (m *MsgScheduleAxelarCorkRequest) Type() string { return TypeMsgScheduleCorkRequest }

// ValidateBasic implements sdk.Msg
func (m *MsgScheduleAxelarCorkRequest) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(m.Signer); err != nil {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, err.Error())
	}

	return m.Cork.ValidateBasic()
}

// GetSignBytes implements sdk.Msg
func (m *MsgScheduleAxelarCorkRequest) GetSignBytes() []byte {
	panic("amino support disabled")
}

// GetSigners implements sdk.Msg
func (m *MsgScheduleAxelarCorkRequest) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{m.MustGetSigner()}
}

// MustGetSigner returns the signer address
func (m *MsgScheduleAxelarCorkRequest) MustGetSigner() sdk.AccAddress {
	addr, err := sdk.AccAddressFromBech32(m.Signer)
	if err != nil {
		panic(err)
	}
	return addr
}

/////////////////////////
// MsgRelayAxelarCorkRequest //
/////////////////////////

// Route implements sdk.Msg
func (m *MsgRelayAxelarCorkRequest) Route() string { return ModuleName }

// Type implements sdk.Msg
func (m *MsgRelayAxelarCorkRequest) Type() string { return TypeMsgRelayCorkRequest }

// ValidateBasic implements sdk.Msg
func (m *MsgRelayAxelarCorkRequest) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(m.Signer); err != nil {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, err.Error())
	}

	if m.TargetContractAddress == "" {
		return fmt.Errorf("cannot relay a cork to an empty address")
	}

	return nil
}

// GetSignBytes implements sdk.Msg
func (m *MsgRelayAxelarCorkRequest) GetSignBytes() []byte {
	panic("amino support disabled")
}

// GetSigners implements sdk.Msg
func (m *MsgRelayAxelarCorkRequest) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{m.MustGetSigner()}
}

// MustGetSigner returns the signer address
func (m *MsgRelayAxelarCorkRequest) MustGetSigner() sdk.AccAddress {
	addr, err := sdk.AccAddressFromBech32(m.Signer)
	if err != nil {
		panic(err)
	}
	return addr
}

///////////////////////////
// MsgBumpAxelarCorkGasRequest //
///////////////////////////

// Route implements sdk.Msg
func (m *MsgBumpAxelarCorkGasRequest) Route() string { return ModuleName }

// Type implements sdk.Msg
func (m *MsgBumpAxelarCorkGasRequest) Type() string { return TypeMsgBumpCorkGasRequest }

// ValidateBasic implements sdk.Msg
func (m *MsgBumpAxelarCorkGasRequest) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(m.Signer); err != nil {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, err.Error())
	}

	return nil
}

// GetSignBytes implements sdk.Msg
func (m *MsgBumpAxelarCorkGasRequest) GetSignBytes() []byte {
	panic("amino support disabled")
}

// GetSigners implements sdk.Msg
func (m *MsgBumpAxelarCorkGasRequest) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{m.MustGetSigner()}
}

// MustGetSigner returns the signer address
func (m *MsgBumpAxelarCorkGasRequest) MustGetSigner() sdk.AccAddress {
	addr, err := sdk.AccAddressFromBech32(m.Signer)
	if err != nil {
		panic(err)
	}
	return addr
}
