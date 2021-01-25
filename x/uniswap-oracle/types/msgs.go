package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// ensure Msg interface compliance at compile time
var (
	_ sdk.Msg = &MsgDelegateFeedConsent{}
	_ sdk.Msg = &MsgAggregateExchangeRatePrevote{}
	_ sdk.Msg = &MsgAggregateExchangeRateVote{}
)

//-------------------------------------------------
//-------------------------------------------------

// NewMsgDelegateFeedConsent creates a MsgDelegateFeedConsent instance
func NewMsgDelegateFeedConsent(operatorAddress sdk.ValAddress, feederAddress sdk.AccAddress) *MsgDelegateFeedConsent {
	return &MsgDelegateFeedConsent{
		Operator: operatorAddress.String(),
		Delegate: feederAddress.String(),
	}
}

// Route implements sdk.Msg
func (msg MsgDelegateFeedConsent) Route() string { return RouterKey }

// Type implements sdk.Msg
func (msg MsgDelegateFeedConsent) Type() string { return "delegatefeeder" }

// GetSignBytes implements sdk.Msg
func (msg MsgDelegateFeedConsent) GetSignBytes() []byte {
	panic("oracle messages do not support amino")
}

// GetSigners implements sdk.Msg
func (msg MsgDelegateFeedConsent) GetSigners() []sdk.AccAddress {
	addr, err := sdk.AccAddressFromBech32(msg.Operator)
	if err != nil {
		return nil
	}
	return []sdk.AccAddress{addr}
}

// ValidateBasic implements sdk.Msg
func (msg MsgDelegateFeedConsent) ValidateBasic() error {
	operator, err := sdk.ValAddressFromBech32(msg.Operator)
	if err != nil || operator.Empty() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "must give valid validator address")
	}

	delegate, err := sdk.AccAddressFromBech32(msg.Delegate)
	if err != nil || delegate.Empty() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "must give valid delegate address")
	}

	return nil
}

// NewMsgAggregateExchangeRatePrevote returns MsgAggregateExchangeRatePrevote instance
func NewMsgAggregateExchangeRatePrevote(hash AggregateVoteHash, feeder sdk.AccAddress, validator sdk.ValAddress) *MsgAggregateExchangeRatePrevote {
	return &MsgAggregateExchangeRatePrevote{
		Hash:      hash,
		Feeder:    feeder.String(),
		Validator: validator.String(),
	}
}

// Route implements sdk.Msg
func (msg MsgAggregateExchangeRatePrevote) Route() string { return RouterKey }

// Type implements sdk.Msg
func (msg MsgAggregateExchangeRatePrevote) Type() string { return "aggregateexchangerateprevote" }

// GetSignBytes implements sdk.Msg
func (msg MsgAggregateExchangeRatePrevote) GetSignBytes() []byte {
	panic("oracle messages do not support amino")
}

// GetSigners implements sdk.Msg
func (msg MsgAggregateExchangeRatePrevote) GetSigners() []sdk.AccAddress {
	addr, err := sdk.AccAddressFromBech32(msg.Feeder)
	if err != nil {
		return nil
	}
	return []sdk.AccAddress{addr}
}

// ValidateBasic Implements sdk.Msg
func (msg MsgAggregateExchangeRatePrevote) ValidateBasic() error {

	// TODO: validate hash here
	// if len(msg.Hash) != tmhash.TruncatedSize {
	// 	return ErrInvalidHashLength
	// }

	feeder, err := sdk.AccAddressFromBech32(msg.Feeder)
	if err != nil || feeder.Empty() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "must give valid feeder address")
	}

	validator, err := sdk.ValAddressFromBech32(msg.Validator)
	if err != nil || validator.Empty() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "must give valid validator address")
	}

	return nil
}

// NewMsgAggregateExchangeRateVote returns MsgAggregateExchangeRateVote instance
func NewMsgAggregateExchangeRateVote(salt string, exchangeRates string, feeder sdk.AccAddress, validator sdk.ValAddress) *MsgAggregateExchangeRateVote {
	return &MsgAggregateExchangeRateVote{
		Salt:          salt,
		ExchangeRates: exchangeRates,
		Feeder:        feeder.String(),
		Validator:     validator.String(),
	}
}

// Route implements sdk.Msg
func (msg MsgAggregateExchangeRateVote) Route() string { return RouterKey }

// Type implements sdk.Msg
func (msg MsgAggregateExchangeRateVote) Type() string { return "aggregateexchangeratevote" }

// GetSignBytes implements sdk.Msg
func (msg MsgAggregateExchangeRateVote) GetSignBytes() []byte {
	panic("oracle messages do not support amino")
}

// GetSigners implements sdk.Msg
func (msg MsgAggregateExchangeRateVote) GetSigners() []sdk.AccAddress {
	addr, err := sdk.AccAddressFromBech32(msg.Feeder)
	if err != nil {
		return nil
	}
	return []sdk.AccAddress{addr}
}

// ValidateBasic implements sdk.Msg
func (msg MsgAggregateExchangeRateVote) ValidateBasic() error {
	feeder, err := sdk.AccAddressFromBech32(msg.Feeder)
	if err != nil || feeder.Empty() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "must give valid feeder address")
	}

	validator, err := sdk.ValAddressFromBech32(msg.Validator)
	if err != nil || validator.Empty() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "must give valid validator address")
	}

	if l := len(msg.ExchangeRates); l == 0 {
		return sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "must provide at least one oracle exchange rate")
	} else if l > 4096 {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "exchange rates string can not exceed 4096 characters")
	}

	exchangeRateTuples, err := sdk.ParseDecCoins(msg.ExchangeRates)
	if err != nil {
		return sdkerrors.Wrap(err, "failed to parse exchange rates string")
	}

	if exchangeRateTuples.Empty() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidCoins, "exchange rate coins cannot be empty")
	}

	for _, tuple := range exchangeRateTuples {
		// Check overflow bit length
		if tuple.Amount.BigInt().BitLen() > 100+sdk.DecimalPrecisionBits {
			return sdkerrors.Wrap(ErrInvalidExchangeRate, "overflow")
		}
	}

	if len(msg.Salt) > 4 || len(msg.Salt) < 1 {
		return sdkerrors.Wrap(ErrInvalidSaltLength, "salt length must be [1, 4]")
	}

	return nil
}
