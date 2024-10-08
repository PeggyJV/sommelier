package types

import (
	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/ethereum/go-ethereum/common"
)

var (
	_ sdk.Msg = &MsgAddAddressMapping{}
	_ sdk.Msg = &MsgRemoveAddressMapping{}
)

const (
	TypeMsgAddAddressMapping    = "add_address_mapping"
	TypeMsgRemoveAddressMapping = "remove_address_mapping"
)

//////////////////////////
// MsgAddAddressMapping //
//////////////////////////

// NewMsgAddAddressMapping return a new MsgAddAddressMapping
func NewMsgAddAddressMapping(evmAddres common.Address, signer sdk.AccAddress) (*MsgAddAddressMapping, error) {
	msg := &MsgAddAddressMapping{
		EvmAddress: evmAddres.Hex(),
		Signer:     signer.String(),
	}

	if err := msg.ValidateBasic(); err != nil {
		return nil, err
	}

	return msg, nil
}

// Route implements sdk.Msg
func (m *MsgAddAddressMapping) Route() string { return ModuleName }

// Type implements sdk.Msg
func (m *MsgAddAddressMapping) Type() string { return TypeMsgAddAddressMapping }

// ValidateBasic implements sdk.Msg
func (m *MsgAddAddressMapping) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(m.Signer); err != nil {
		return errorsmod.Wrap(sdkerrors.ErrInvalidAddress, err.Error())
	}

	if !common.IsHexAddress(m.EvmAddress) {
		return errorsmod.Wrapf(ErrInvalidEvmAddress, "%s is not a valid hex address", m.EvmAddress)
	}

	return nil
}

// GetSignBytes implements sdk.Msg
func (m *MsgAddAddressMapping) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(m))
}

// GetSigners implements sdk.Msg
func (m *MsgAddAddressMapping) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{m.MustGetSigner()}
}

// MustGetSigner returns the signer address (which is also the bidder)
func (m *MsgAddAddressMapping) MustGetSigner() sdk.AccAddress {
	addr, err := sdk.AccAddressFromBech32(m.Signer)
	if err != nil {
		panic(err)
	}
	return addr
}

/////////////////////////////
// MsgRemoveAddressMapping //
/////////////////////////////

// NewMsgRemoveAddressMapping return a new MsgRemoveAddressMapping
func NewMsgRemoveAddressMapping(signer sdk.AccAddress) (*MsgRemoveAddressMapping, error) {
	msg := &MsgRemoveAddressMapping{
		Signer: signer.String(),
	}

	if err := msg.ValidateBasic(); err != nil {
		return nil, err
	}

	return msg, nil
}

// Route implements sdk.Msg
func (m *MsgRemoveAddressMapping) Route() string { return ModuleName }

// Type implements sdk.Msg
func (m *MsgRemoveAddressMapping) Type() string { return TypeMsgRemoveAddressMapping }

// ValidateBasic implements sdk.Msg
func (m *MsgRemoveAddressMapping) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(m.Signer); err != nil {
		return errorsmod.Wrap(sdkerrors.ErrInvalidAddress, err.Error())
	}

	return nil
}

// GetSignBytes implements sdk.Msg
func (m *MsgRemoveAddressMapping) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(m))
}

// GetSigners implements sdk.Msg
func (m *MsgRemoveAddressMapping) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{m.MustGetSigner()}
}

// MustGetSigner returns the signer address (which is also the bidder)
func (m *MsgRemoveAddressMapping) MustGetSigner() sdk.AccAddress {
	addr, err := sdk.AccAddressFromBech32(m.Signer)
	if err != nil {
		panic(err)
	}
	return addr
}
