package types

import (
	fmt "fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var (
	_ sdk.Msg = &MsgAddPublisherIntent{}
	_ sdk.Msg = &MsgAddSubscriberIntent{}
	_ sdk.Msg = &MsgAddSubscriber{}
	_ sdk.Msg = &MsgRemovePublisherIntent{}
	_ sdk.Msg = &MsgRemoveSubscriberIntent{}
	_ sdk.Msg = &MsgRemoveSubscriber{}
	_ sdk.Msg = &MsgRemovePublisher{}
)

const (
	TypeMsgAddPublisherIntent     = "add_publisher_intent"
	TypeMsgAddSubscriberIntent    = "add_subscriber_intent"
	TypeMsgAddSubscriber          = "add_subscriber"
	TypeMsgRemovePublisherIntent  = "remove_publisher_intent"
	TypeMsgRemoveSubscriberIntent = "remove_subscriber_intent"
	TypeMsgRemoveSubscriber       = "remove_subscriber"
	TypeMsgRemovePublisher        = "remove_publisher"
)

///////////////////////////
// MsgAddPublisherIntent //
///////////////////////////

// NewMsgAddPublisherIntent returns a new MsgAddPublisherIntent
func NewMsgAddPublisherIntent(publisherIntent PublisherIntent, signer sdk.AccAddress) (*MsgAddPublisherIntent, error) {
	if signer == nil {
		return nil, fmt.Errorf("no signer provided")
	}

	return &MsgAddPublisherIntent{
		PublisherIntent: &publisherIntent,
		Signer:          signer.String(),
	}, nil
}

// Route implements sdk.Msg
func (m *MsgAddPublisherIntent) Route() string { return ModuleName }

// Type implements sdk.Msg
func (m *MsgAddPublisherIntent) Type() string { return TypeMsgAddPublisherIntent }

// ValidateBasic implements sdk.Msg
func (m *MsgAddPublisherIntent) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(m.Signer); err != nil {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, err.Error())
	}

	if m.PublisherIntent == nil {
		return fmt.Errorf("empty PublisherIntent")
	}

	return m.PublisherIntent.ValidateBasic()
}

// GetSignBytes implements sdk.Msg
func (m *MsgAddPublisherIntent) GetSignBytes() []byte {
	panic("amino support disabled")
}

// GetSigners implements sdk.Msg
func (m *MsgAddPublisherIntent) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{m.MustGetSigner()}
}

// MustGetSigner returns the signer address
func (m *MsgAddPublisherIntent) MustGetSigner() sdk.AccAddress {
	addr, err := sdk.AccAddressFromBech32(m.Signer)
	if err != nil {
		panic(err)
	}
	return addr
}

////////////////////////////
// MsgAddSubscriberIntent //
////////////////////////////

// NewMsgAddSubscriberIntent returns a new MsgAddSubscriberIntent
func NewMsgAddSubscriberIntent(subscriberIntent SubscriberIntent, signer sdk.AccAddress) (*MsgAddSubscriberIntent, error) {
	if signer == nil {
		return nil, fmt.Errorf("no signer provided")
	}

	return &MsgAddSubscriberIntent{
		SubscriberIntent: &subscriberIntent,
		Signer:           signer.String(),
	}, nil
}

// Route implements sdk.Msg
func (m *MsgAddSubscriberIntent) Route() string { return ModuleName }

// Type implements sdk.Msg
func (m *MsgAddSubscriberIntent) Type() string { return TypeMsgAddSubscriberIntent }

// ValidateBasic implements sdk.Msg
func (m *MsgAddSubscriberIntent) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(m.Signer); err != nil {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, err.Error())
	}

	if m.SubscriberIntent == nil {
		return fmt.Errorf("empty SubscriberIntent")
	}

	return m.SubscriberIntent.ValidateBasic()
}

// GetSignBytes implements sdk.Msg
func (m *MsgAddSubscriberIntent) GetSignBytes() []byte {
	panic("amino support disabled")
}

// GetSigners implements sdk.Msg
func (m *MsgAddSubscriberIntent) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{m.MustGetSigner()}
}

// MustGetSigner returns the signer address
func (m *MsgAddSubscriberIntent) MustGetSigner() sdk.AccAddress {
	addr, err := sdk.AccAddressFromBech32(m.Signer)
	if err != nil {
		panic(err)
	}
	return addr
}

//////////////////////
// MsgAddSubscriber //
//////////////////////

// NewMsgAddSubscriber returns a new MsgAddSubscriber
func NewMsgAddSubscriber(subscriber Subscriber, signer sdk.AccAddress) (*MsgAddSubscriber, error) {
	if signer == nil {
		return nil, fmt.Errorf("no signer provided")
	}

	return &MsgAddSubscriber{
		Subscriber: &subscriber,
		Signer:     signer.String(),
	}, nil
}

// Route implements sdk.Msg
func (m *MsgAddSubscriber) Route() string { return ModuleName }

// Type implements sdk.Msg
func (m *MsgAddSubscriber) Type() string { return TypeMsgAddSubscriber }

// ValidateBasic implements sdk.Msg
func (m *MsgAddSubscriber) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(m.Signer); err != nil {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, err.Error())
	}

	if m.Subscriber == nil {
		return fmt.Errorf("empty Subscriber")
	}

	return m.Subscriber.ValidateBasic()
}

// GetSignBytes implements sdk.Msg
func (m *MsgAddSubscriber) GetSignBytes() []byte {
	panic("amino support disabled")
}

// GetSigners implements sdk.Msg
func (m *MsgAddSubscriber) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{m.MustGetSigner()}
}

// MustGetSigner returns the signer address
func (m *MsgAddSubscriber) MustGetSigner() sdk.AccAddress {
	addr, err := sdk.AccAddressFromBech32(m.Signer)
	if err != nil {
		panic(err)
	}
	return addr
}

//////////////////////////////
// MsgRemovePublisherIntent //
//////////////////////////////

// NewMsgRemovePublisherIntent returns a new MsgRemovePublisherIntent
func NewMsgRemovePublisherIntent(publisherIntent PublisherIntent, signer sdk.AccAddress) (*MsgRemovePublisherIntent, error) {
	if signer == nil {
		return nil, fmt.Errorf("no signer provided")
	}

	return &MsgRemovePublisherIntent{
		PublisherIntent: &publisherIntent,
		Signer:          signer.String(),
	}, nil
}

// Route implements sdk.Msg
func (m *MsgRemovePublisherIntent) Route() string { return ModuleName }

// Type implements sdk.Msg
func (m *MsgRemovePublisherIntent) Type() string { return TypeMsgRemovePublisherIntent }

// ValidateBasic implements sdk.Msg
func (m *MsgRemovePublisherIntent) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(m.Signer); err != nil {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, err.Error())
	}

	if m.PublisherIntent == nil {
		return fmt.Errorf("empty PublisherIntent")
	}

	return m.PublisherIntent.ValidateBasic()
}

// GetSignBytes implements sdk.Msg
func (m *MsgRemovePublisherIntent) GetSignBytes() []byte {
	panic("amino support disabled")
}

// GetSigners implements sdk.Msg
func (m *MsgRemovePublisherIntent) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{m.MustGetSigner()}
}

// MustGetSigner returns the signer address
func (m *MsgRemovePublisherIntent) MustGetSigner() sdk.AccAddress {
	addr, err := sdk.AccAddressFromBech32(m.Signer)
	if err != nil {
		panic(err)
	}
	return addr
}

///////////////////////////////
// MsgRemoveSubscriberIntent //
///////////////////////////////

// NewMsgRemoveSubscriberIntent returns a new MsgRemoveSubscriberIntent
func NewMsgRemoveSubscriberIntent(subscriberIntent SubscriberIntent, signer sdk.AccAddress) (*MsgRemoveSubscriberIntent, error) {
	if signer == nil {
		return nil, fmt.Errorf("no signer provided")
	}

	return &MsgRemoveSubscriberIntent{
		SubscriberIntent: &subscriberIntent,
		Signer:           signer.String(),
	}, nil
}

// Route implements sdk.Msg
func (m *MsgRemoveSubscriberIntent) Route() string { return ModuleName }

// Type implements sdk.Msg
func (m *MsgRemoveSubscriberIntent) Type() string { return TypeMsgRemoveSubscriberIntent }

// ValidateBasic implements sdk.Msg
func (m *MsgRemoveSubscriberIntent) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(m.Signer); err != nil {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, err.Error())
	}

	if m.SubscriberIntent == nil {
		return fmt.Errorf("empty SubscriberIntent")
	}

	return m.SubscriberIntent.ValidateBasic()
}

// GetSignBytes implements sdk.Msg
func (m *MsgRemoveSubscriberIntent) GetSignBytes() []byte {
	panic("amino support disabled")
}

// GetSigners implements sdk.Msg
func (m *MsgRemoveSubscriberIntent) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{m.MustGetSigner()}
}

// MustGetSigner returns the signer address
func (m *MsgRemoveSubscriberIntent) MustGetSigner() sdk.AccAddress {
	addr, err := sdk.AccAddressFromBech32(m.Signer)
	if err != nil {
		panic(err)
	}
	return addr
}

/////////////////////////
// MsgRemoveSubscriber //
/////////////////////////

// NewMsgRemoveSubscriber returns a new MsgRemoveSubscriber
func NewMsgRemoveSubscriber(subscriber Subscriber, signer sdk.AccAddress) (*MsgRemoveSubscriber, error) {
	if signer == nil {
		return nil, fmt.Errorf("no signer provided")
	}

	return &MsgRemoveSubscriber{
		Subscriber: &subscriber,
		Signer:     signer.String(),
	}, nil
}

// Route implements sdk.Msg
func (m *MsgRemoveSubscriber) Route() string { return ModuleName }

// Type implements sdk.Msg
func (m *MsgRemoveSubscriber) Type() string { return TypeMsgRemoveSubscriber }

// ValidateBasic implements sdk.Msg
func (m *MsgRemoveSubscriber) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(m.Signer); err != nil {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, err.Error())
	}

	if m.Subscriber == nil {
		return fmt.Errorf("empty Subscriber")
	}

	return m.Subscriber.ValidateBasic()
}

// GetSignBytes implements sdk.Msg
func (m *MsgRemoveSubscriber) GetSignBytes() []byte {
	panic("amino support disabled")
}

// GetSigners implements sdk.Msg
func (m *MsgRemoveSubscriber) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{m.MustGetSigner()}
}

// MustGetSigner returns the signer address
func (m *MsgRemoveSubscriber) MustGetSigner() sdk.AccAddress {
	addr, err := sdk.AccAddressFromBech32(m.Signer)
	if err != nil {
		panic(err)
	}
	return addr
}

////////////////////////
// MsgRemovePublisher //
////////////////////////

// NewMsgRemovePublisher returns a new MsgRemovePublisher
func NewMsgRemovePublisher(publisher Publisher, signer sdk.AccAddress) (*MsgRemovePublisher, error) {
	if signer == nil {
		return nil, fmt.Errorf("no signer provided")
	}

	return &MsgRemovePublisher{
		Publisher: &publisher,
		Signer:    signer.String(),
	}, nil
}

// Route implements sdk.Msg
func (m *MsgRemovePublisher) Route() string { return ModuleName }

// Type implements sdk.Msg
func (m *MsgRemovePublisher) Type() string { return TypeMsgRemovePublisher }

// ValidateBasic implements sdk.Msg
func (m *MsgRemovePublisher) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(m.Signer); err != nil {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, err.Error())
	}

	if m.Publisher == nil {
		return fmt.Errorf("empty Publisher")
	}

	return m.Publisher.ValidateBasic()
}

// GetSignBytes implements sdk.Msg
func (m *MsgRemovePublisher) GetSignBytes() []byte {
	panic("amino support disabled")
}

// GetSigners implements sdk.Msg
func (m *MsgRemovePublisher) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{m.MustGetSigner()}
}

// MustGetSigner returns the signer address
func (m *MsgRemovePublisher) MustGetSigner() sdk.AccAddress {
	addr, err := sdk.AccAddressFromBech32(m.Signer)
	if err != nil {
		panic(err)
	}
	return addr
}
