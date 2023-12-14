package types

import (
	"fmt"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/ethereum/go-ethereum/common"
)

var (
	_ sdk.Msg = &MsgScheduleAxelarCorkRequest{}
)

const (
	TypeMsgScheduleCorkRequest            = "axelar_cork_schedule"
	TypeMsgRelayCorkRequest               = "axelar_cork_relay"
	TypeMsgBumpCorkGasRequest             = "axelar_cork_bump_gas"
	TypeMsgCancelAxelarCorkRequest        = "axelar_cancel_cork"
	TypeMsgRelayAxelarProxyUpgradeRequest = "axelar_proxy_upgrade_relay"
)

////////////////////////////
// MsgScheduleAxelarCorkRequest //
////////////////////////////

// NewMsgScheduleCorkRequest return a new MsgScheduleAxelarCorkRequest
func NewMsgScheduleCorkRequest(chainID uint64, body []byte, address common.Address, deadline uint64, blockHeight uint64, signer sdk.AccAddress) (*MsgScheduleAxelarCorkRequest, error) {
	return &MsgScheduleAxelarCorkRequest{
		Cork: &AxelarCork{
			ChainId:               chainID,
			EncodedContractCall:   body,
			TargetContractAddress: address.String(),
			Deadline:              deadline,
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
		return errorsmod.Wrap(sdkerrors.ErrInvalidAddress, err.Error())
	}

	return m.Cork.ValidateBasic()
}

// GetSignBytes implements sdk.Msg
func (m *MsgScheduleAxelarCorkRequest) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(m))
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
		return errorsmod.Wrap(sdkerrors.ErrInvalidAddress, err.Error())
	}

	if m.TargetContractAddress == "" {
		return fmt.Errorf("cannot relay a cork to an empty address")
	}

	return nil
}

// GetSignBytes implements sdk.Msg
func (m *MsgRelayAxelarCorkRequest) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(m))
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

///////////////////////////////
// MsgRelayAxelarProxyUpgradeRequest //
///////////////////////////////

// Route implements sdk.Msg
func (m *MsgRelayAxelarProxyUpgradeRequest) Route() string { return ModuleName }

// Type implements sdk.Msg
func (m *MsgRelayAxelarProxyUpgradeRequest) Type() string {
	return TypeMsgRelayAxelarProxyUpgradeRequest
}

// ValidateBasic implements sdk.Msg
func (m *MsgRelayAxelarProxyUpgradeRequest) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(m.Signer); err != nil {
		return errorsmod.Wrap(sdkerrors.ErrInvalidAddress, err.Error())
	}

	return nil
}

// GetSignBytes implements sdk.Msg
func (m *MsgRelayAxelarProxyUpgradeRequest) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(m))
}

// GetSigners implements sdk.Msg
func (m *MsgRelayAxelarProxyUpgradeRequest) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{m.MustGetSigner()}
}

// MustGetSigner returns the signer address
func (m *MsgRelayAxelarProxyUpgradeRequest) MustGetSigner() sdk.AccAddress {
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
		return errorsmod.Wrap(sdkerrors.ErrInvalidAddress, err.Error())
	}

	return nil
}

// GetSignBytes implements sdk.Msg
func (m *MsgBumpAxelarCorkGasRequest) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(m))
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

///////////////////////////
// MsgCancelAxelarCorkRequest //
///////////////////////////

// Route implements sdk.Msg
func (m *MsgCancelAxelarCorkRequest) Route() string { return ModuleName }

// Type implements sdk.Msg
func (m *MsgCancelAxelarCorkRequest) Type() string { return TypeMsgBumpCorkGasRequest }

// ValidateBasic implements sdk.Msg
func (m *MsgCancelAxelarCorkRequest) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(m.Signer); err != nil {
		return errorsmod.Wrap(sdkerrors.ErrInvalidAddress, err.Error())
	}

	return nil
}

// GetSignBytes implements sdk.Msg
func (m *MsgCancelAxelarCorkRequest) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(m))
}

// GetSigners implements sdk.Msg
func (m *MsgCancelAxelarCorkRequest) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{m.MustGetSigner()}
}

// MustGetSigner returns the signer address
func (m *MsgCancelAxelarCorkRequest) MustGetSigner() sdk.AccAddress {
	addr, err := sdk.AccAddressFromBech32(m.Signer)
	if err != nil {
		panic(err)
	}
	return addr
}
