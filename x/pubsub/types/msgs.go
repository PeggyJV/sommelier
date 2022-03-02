package types

import sdk "github.com/cosmos/cosmos-sdk/types"

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

// TODO(bolten): fill out msg boilerplate

///////////////////////////
// MsgAddPublisherIntent //
///////////////////////////

// NewMsgAddPublisherIntent returns a new MsgAddPublisherIntent
func NewMsgAddPublisherIntent() (*MsgAddPublisherIntent, error) {
	return nil, nil
}

// Route implements sdk.Msg
func (m *MsgAddPublisherIntent) Route() string { return ModuleName }

// Type implements sdk.Msg
func (m *MsgAddPublisherIntent) Type() string { return TypeMsgAddPublisherIntent }

// ValidateBasic implements sdk.Msg
func (m *MsgAddPublisherIntent) ValidateBasic() error {
	return nil
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
func NewMsgAddSubscriberIntent() (*MsgAddSubscriberIntent, error) {
	return nil, nil
}

// Route implements sdk.Msg
func (m *MsgAddSubscriberIntent) Route() string { return ModuleName }

// Type implements sdk.Msg
func (m *MsgAddSubscriberIntent) Type() string { return TypeMsgAddSubscriberIntent }

// ValidateBasic implements sdk.Msg
func (m *MsgAddSubscriberIntent) ValidateBasic() error {
	return nil
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
func NewMsgAddSubscriber() (*MsgAddSubscriber, error) {
	return nil, nil
}

// Route implements sdk.Msg
func (m *MsgAddSubscriber) Route() string { return ModuleName }

// Type implements sdk.Msg
func (m *MsgAddSubscriber) Type() string { return TypeMsgAddSubscriber }

// ValidateBasic implements sdk.Msg
func (m *MsgAddSubscriber) ValidateBasic() error {
	return nil
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
func NewMsgRemovePublisherIntent() (*MsgRemovePublisherIntent, error) {
	return nil, nil
}

// Route implements sdk.Msg
func (m *MsgRemovePublisherIntent) Route() string { return ModuleName }

// Type implements sdk.Msg
func (m *MsgRemovePublisherIntent) Type() string { return TypeMsgRemovePublisherIntent }

// ValidateBasic implements sdk.Msg
func (m *MsgRemovePublisherIntent) ValidateBasic() error {
	return nil
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
func NewMsgRemoveSubscriberIntent() (*MsgRemoveSubscriberIntent, error) {
	return nil, nil
}

// Route implements sdk.Msg
func (m *MsgRemoveSubscriberIntent) Route() string { return ModuleName }

// Type implements sdk.Msg
func (m *MsgRemoveSubscriberIntent) Type() string { return TypeMsgRemoveSubscriberIntent }

// ValidateBasic implements sdk.Msg
func (m *MsgRemoveSubscriberIntent) ValidateBasic() error {
	return nil
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
func NewMsgRemoveSubscriber() (*MsgRemoveSubscriber, error) {
	return nil, nil
}

// Route implements sdk.Msg
func (m *MsgRemoveSubscriber) Route() string { return ModuleName }

// Type implements sdk.Msg
func (m *MsgRemoveSubscriber) Type() string { return TypeMsgRemoveSubscriber }

// ValidateBasic implements sdk.Msg
func (m *MsgRemoveSubscriber) ValidateBasic() error {
	return nil
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
func NewMsgRemovePublisher() (*MsgRemovePublisher, error) {
	return nil, nil
}

// Route implements sdk.Msg
func (m *MsgRemovePublisher) Route() string { return ModuleName }

// Type implements sdk.Msg
func (m *MsgRemovePublisher) Type() string { return TypeMsgRemovePublisher }

// ValidateBasic implements sdk.Msg
func (m *MsgRemovePublisher) ValidateBasic() error {
	return nil
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
