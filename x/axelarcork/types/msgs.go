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
	_ sdk.Msg = &MsgRelayAxelarCorkRequest{}
	_ sdk.Msg = &MsgBumpAxelarCorkGasRequest{}
	_ sdk.Msg = &MsgCancelAxelarCorkRequest{}
	_ sdk.Msg = &MsgRelayAxelarCorkRequest{}
)

const (
	TypeMsgScheduleAxelarCorkRequest      = "axelar_cork_schedule"
	TypeMsgRelayAxelarCorkRequest         = "axelar_cork_relay"
	TypeMsgBumpAxelarCorkGasRequest       = "axelar_cork_bump_gas"
	TypeMsgCancelAxelarCorkRequest        = "axelar_cancel_cork"
	TypeMsgRelayAxelarProxyUpgradeRequest = "axelar_proxy_upgrade_relay"
)

//////////////////////////////////
// MsgScheduleAxelarCorkRequest //
//////////////////////////////////

// NewMsgScheduleAxelarCorkRequest return a new MsgScheduleAxelarCorkRequest
func NewMsgScheduleAxelarCorkRequest(chainID uint64, body []byte, address common.Address, deadline uint64, blockHeight uint64, signer sdk.AccAddress) (*MsgScheduleAxelarCorkRequest, error) {
	return &MsgScheduleAxelarCorkRequest{
		Cork: &AxelarCork{
			ChainId:               chainID,
			EncodedContractCall:   body,
			TargetContractAddress: address.String(),
			Deadline:              deadline,
		},
		ChainId:     chainID,
		BlockHeight: blockHeight,
		Signer:      signer.String(),
	}, nil
}

// Route implements sdk.Msg
func (m *MsgScheduleAxelarCorkRequest) Route() string { return ModuleName }

// Type implements sdk.Msg
func (m *MsgScheduleAxelarCorkRequest) Type() string { return TypeMsgScheduleAxelarCorkRequest }

// ValidateBasic implements sdk.Msg
func (m *MsgScheduleAxelarCorkRequest) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(m.Signer); err != nil {
		return errorsmod.Wrap(sdkerrors.ErrInvalidAddress, err.Error())
	}

	if m.ChainId == 0 {
		return fmt.Errorf("chain ID must be non-zero")
	}

	if m.BlockHeight == 0 {
		return fmt.Errorf("block height must be greater than zero")
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

///////////////////////////////
// MsgRelayAxelarCorkRequest //
///////////////////////////////

// NewMsgRelayAxelarCorkRequest returns a new MsgRelayAxelarCorkRequest
func NewMsgRelayAxelarCorkRequest(signer sdk.Address, token sdk.Coin, fee uint64, chainID uint64, address common.Address) (*MsgRelayAxelarCorkRequest, error) {
	return &MsgRelayAxelarCorkRequest{
		Signer:                signer.String(),
		Token:                 token,
		Fee:                   fee,
		ChainId:               chainID,
		TargetContractAddress: address.String(),
	}, nil
}

// Route implements sdk.Msg
func (m *MsgRelayAxelarCorkRequest) Route() string { return ModuleName }

// Type implements sdk.Msg
func (m *MsgRelayAxelarCorkRequest) Type() string { return TypeMsgRelayAxelarCorkRequest }

// ValidateBasic implements sdk.Msg
func (m *MsgRelayAxelarCorkRequest) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(m.Signer); err != nil {
		return errorsmod.Wrap(sdkerrors.ErrInvalidAddress, err.Error())
	}

	if !m.Token.IsValid() || !m.Token.IsPositive() {
		return fmt.Errorf("token amount must be positive")
	}

	if m.Fee == 0 {
		return fmt.Errorf("fee must be greather than zero")
	}

	if m.ChainId == 0 {
		return fmt.Errorf("chain ID must be greater than zero")
	}

	if m.TargetContractAddress == "" {
		return fmt.Errorf("cannot relay a cork to an empty address")
	}

	if !common.IsHexAddress(m.TargetContractAddress) {
		return errorsmod.Wrapf(ErrInvalidEVMAddress, "%s", m.TargetContractAddress)
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

///////////////////////////////////////
// MsgRelayAxelarProxyUpgradeRequest //
///////////////////////////////////////

// NewMsgRelayAxelarProxyUpgradeRequest returns a new MsgRelayAxelarProxyUpgradeRequest
func NewMsgRelayAxelarProxyUpgradeRequest(signer sdk.AccAddress, token sdk.Coin, fee uint64, chainID uint64) (*MsgRelayAxelarProxyUpgradeRequest, error) {
	return &MsgRelayAxelarProxyUpgradeRequest{
		Signer:  signer.String(),
		Token:   token,
		Fee:     fee,
		ChainId: chainID,
	}, nil
}

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

	if !m.Token.IsValid() || !m.Token.IsPositive() {
		return fmt.Errorf("token amount must be positive")
	}

	if m.Fee == 0 {
		return fmt.Errorf("fee must be greather than zero")
	}

	if m.ChainId == 0 {
		return fmt.Errorf("chain ID must be greater than zero")
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

/////////////////////////////////
// MsgBumpAxelarCorkGasRequest //
/////////////////////////////////

// NewMsgBumpAxelarCorkGasRequest returns a new MsgBumpAxelarCorkGasRequest
func NewMsgBumpAxelarCorkGasRequest(signer sdk.AccAddress, token sdk.Coin, messageID string) (*MsgBumpAxelarCorkGasRequest, error) {
	return &MsgBumpAxelarCorkGasRequest{
		Signer:    signer.String(),
		Token:     token,
		MessageId: messageID,
	}, nil
}

// Route implements sdk.Msg
func (m *MsgBumpAxelarCorkGasRequest) Route() string { return ModuleName }

// Type implements sdk.Msg
func (m *MsgBumpAxelarCorkGasRequest) Type() string { return TypeMsgBumpAxelarCorkGasRequest }

// ValidateBasic implements sdk.Msg
func (m *MsgBumpAxelarCorkGasRequest) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(m.Signer); err != nil {
		return errorsmod.Wrap(sdkerrors.ErrInvalidAddress, err.Error())
	}

	if !m.Token.IsValid() || !m.Token.IsPositive() {
		return fmt.Errorf("token amount must be positive")
	}

	if m.MessageId == "" {
		return fmt.Errorf("message ID must be present")
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

////////////////////////////////
// MsgCancelAxelarCorkRequest //
////////////////////////////////

// NewMsgCancelAxelarCorkRequest returns a new MsgCancelAxelarCorkRequest
func NewMsgCancelAxelarCorkRequest(signer sdk.AccAddress, chainID uint64, address common.Address) (*MsgCancelAxelarCorkRequest, error) {
	return &MsgCancelAxelarCorkRequest{
		Signer:                signer.String(),
		ChainId:               chainID,
		TargetContractAddress: address.String(),
	}, nil
}

// Route implements sdk.Msg
func (m *MsgCancelAxelarCorkRequest) Route() string { return ModuleName }

// Type implements sdk.Msg
func (m *MsgCancelAxelarCorkRequest) Type() string { return TypeMsgCancelAxelarCorkRequest }

// ValidateBasic implements sdk.Msg
func (m *MsgCancelAxelarCorkRequest) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(m.Signer); err != nil {
		return errorsmod.Wrap(sdkerrors.ErrInvalidAddress, err.Error())
	}

	if m.TargetContractAddress == "" {
		return fmt.Errorf("cannot cancel a cork to an empty address")
	}

	if !common.IsHexAddress(m.TargetContractAddress) {
		return errorsmod.Wrapf(ErrInvalidEVMAddress, "%s", m.TargetContractAddress)
	}

	if m.ChainId == 0 {
		return fmt.Errorf("chain ID must be greater than zero")
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
