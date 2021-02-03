package types

import (
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// Ensure Msg interface compliance at compile time
var (
	_ sdk.Msg = &MsgStoploss{}
	_ sdk.Msg = &MsgDeleteStoploss{}
)

//-------------------------------------------------
//-------------------------------------------------

// NewMsgStoploss creates a MsgStoploss instance
func NewMsgStoploss(address sdk.Address, stoploss *Stoploss) *MsgStoploss {
	if address == nil {
		return nil
	}

	return &MsgStoploss{
		Address:  address.String(),
		Stoploss: stoploss,
	}
}

// Route implements sdk.Msg
func (msg MsgStoploss) Route() string { return RouterKey }

// Type implements sdk.Msg
func (msg MsgStoploss) Type() string { return "stoploss" }

// GetSignBytes implements sdk.Msg
func (msg MsgStoploss) GetSignBytes() []byte {
	panic("impermanent loss module messages do not support amino")
}

// GetSigners implements sdk.Msg
func (msg MsgStoploss) GetSigners() []sdk.AccAddress {
	addr, err := sdk.AccAddressFromBech32(msg.Address)
	if err != nil {
		return nil
	}
	return []sdk.AccAddress{addr}
}

// ValidateBasic implements sdk.Msg
func (msg MsgStoploss) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Address)
	if err != nil {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, err.Error())
	}

	if msg.Stoploss == nil {
		return sdkerrors.Wrap(ErrStoplossInvalid, "cannot be nil")
	}

	return msg.Stoploss.Validate()
}

// NewMsgDeleteStoploss creates a MsgDeleteStoploss instance
func NewMsgDeleteStoploss(address sdk.Address, uniswapPairID string) *MsgDeleteStoploss {
	if address == nil {
		return nil
	}

	return &MsgDeleteStoploss{
		Address:       address.String(),
		UniswapPairId: uniswapPairID,
	}
}

// Route implements sdk.Msg
func (msg MsgDeleteStoploss) Route() string { return RouterKey }

// Type implements sdk.Msg
func (msg MsgDeleteStoploss) Type() string { return "delete-stoploss" }

// GetSignBytes implements sdk.Msg
func (msg MsgDeleteStoploss) GetSignBytes() []byte {
	panic("impermanent loss module messages do not support amino")
}

// GetSigners implements sdk.Msg
func (msg MsgDeleteStoploss) GetSigners() []sdk.AccAddress {
	addr, err := sdk.AccAddressFromBech32(msg.Address)
	if err != nil {
		return nil
	}
	return []sdk.AccAddress{addr}
}

// ValidateBasic implements sdk.Msg
func (msg MsgDeleteStoploss) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Address)
	if err != nil {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, err.Error())
	}

	if strings.TrimSpace(msg.UniswapPairId) == "" {
		return sdkerrors.Wrap(ErrUniswapPairInvalid, "cannot be blank")
	}

	return nil
}
