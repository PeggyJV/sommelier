package types

import (
	fmt "fmt"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var (
	_ sdk.Msg = &MsgAddPublisherIntentRequest{}
	_ sdk.Msg = &MsgAddSubscriberIntentRequest{}
	_ sdk.Msg = &MsgAddSubscriberRequest{}
	_ sdk.Msg = &MsgRemovePublisherIntentRequest{}
	_ sdk.Msg = &MsgRemoveSubscriberIntentRequest{}
	_ sdk.Msg = &MsgRemoveSubscriberRequest{}
	_ sdk.Msg = &MsgRemovePublisherRequest{}
)

const (
	TypeMsgAddPublisherIntentRequest     = "add_publisher_intent"
	TypeMsgAddSubscriberIntentRequest    = "add_subscriber_intent"
	TypeMsgAddSubscriberRequest          = "add_subscriber"
	TypeMsgRemovePublisherIntentRequest  = "remove_publisher_intent"
	TypeMsgRemoveSubscriberIntentRequest = "remove_subscriber_intent"
	TypeMsgRemoveSubscriberRequest       = "remove_subscriber"
	TypeMsgRemovePublisherRequest        = "remove_publisher"
)

//////////////////////////////////
// MsgAddPublisherIntentRequest //
//////////////////////////////////

// NewMsgAddPublisherIntentRequest returns a new MsgAddPublisherIntentRequest
func NewMsgAddPublisherIntentRequest(publisherIntent PublisherIntent, signer sdk.AccAddress) (*MsgAddPublisherIntentRequest, error) {
	if signer == nil {
		return nil, fmt.Errorf("no signer provided")
	}

	return &MsgAddPublisherIntentRequest{
		PublisherIntent: &publisherIntent,
		Signer:          signer.String(),
	}, nil
}

// Route implements sdk.Msg
func (m *MsgAddPublisherIntentRequest) Route() string { return ModuleName }

// Type implements sdk.Msg
func (m *MsgAddPublisherIntentRequest) Type() string { return TypeMsgAddPublisherIntentRequest }

// ValidateBasic implements sdk.Msg
func (m *MsgAddPublisherIntentRequest) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(m.Signer); err != nil {
		return errorsmod.Wrap(sdkerrors.ErrInvalidAddress, err.Error())
	}

	if m.PublisherIntent == nil {
		return fmt.Errorf("empty PublisherIntent")
	}

	return m.PublisherIntent.ValidateBasic()
}

// GetSignBytes implements sdk.Msg
func (m *MsgAddPublisherIntentRequest) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(m))
}

// GetSigners implements sdk.Msg
func (m *MsgAddPublisherIntentRequest) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{m.MustGetSigner()}
}

// MustGetSigner returns the signer address
func (m *MsgAddPublisherIntentRequest) MustGetSigner() sdk.AccAddress {
	addr, err := sdk.AccAddressFromBech32(m.Signer)
	if err != nil {
		panic(err)
	}
	return addr
}

///////////////////////////////////
// MsgAddSubscriberIntentRequest //
///////////////////////////////////

// NewMsgAddSubscriberIntentRequest returns a new MsgAddSubscriberIntentRequest
func NewMsgAddSubscriberIntentRequest(subscriberIntent SubscriberIntent, signer sdk.AccAddress) (*MsgAddSubscriberIntentRequest, error) {
	if signer == nil {
		return nil, fmt.Errorf("no signer provided")
	}

	return &MsgAddSubscriberIntentRequest{
		SubscriberIntent: &subscriberIntent,
		Signer:           signer.String(),
	}, nil
}

// Route implements sdk.Msg
func (m *MsgAddSubscriberIntentRequest) Route() string { return ModuleName }

// Type implements sdk.Msg
func (m *MsgAddSubscriberIntentRequest) Type() string { return TypeMsgAddSubscriberIntentRequest }

// ValidateBasic implements sdk.Msg
func (m *MsgAddSubscriberIntentRequest) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(m.Signer); err != nil {
		return errorsmod.Wrap(sdkerrors.ErrInvalidAddress, err.Error())
	}

	if m.SubscriberIntent == nil {
		return fmt.Errorf("empty SubscriberIntent")
	}

	return m.SubscriberIntent.ValidateBasic()
}

// GetSignBytes implements sdk.Msg
func (m *MsgAddSubscriberIntentRequest) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(m))
}

// GetSigners implements sdk.Msg
func (m *MsgAddSubscriberIntentRequest) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{m.MustGetSigner()}
}

// MustGetSigner returns the signer address
func (m *MsgAddSubscriberIntentRequest) MustGetSigner() sdk.AccAddress {
	addr, err := sdk.AccAddressFromBech32(m.Signer)
	if err != nil {
		panic(err)
	}
	return addr
}

/////////////////////////////
// MsgAddSubscriberRequest //
/////////////////////////////

// NewMsgAddSubscriberRequest returns a new MsgAddSubscriberRequest
func NewMsgAddSubscriberRequest(subscriber Subscriber, signer sdk.AccAddress) (*MsgAddSubscriberRequest, error) {
	if signer == nil {
		return nil, fmt.Errorf("no signer provided")
	}

	return &MsgAddSubscriberRequest{
		Subscriber: &subscriber,
		Signer:     signer.String(),
	}, nil
}

// Route implements sdk.Msg
func (m *MsgAddSubscriberRequest) Route() string { return ModuleName }

// Type implements sdk.Msg
func (m *MsgAddSubscriberRequest) Type() string { return TypeMsgAddSubscriberRequest }

// ValidateBasic implements sdk.Msg
func (m *MsgAddSubscriberRequest) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(m.Signer); err != nil {
		return errorsmod.Wrap(sdkerrors.ErrInvalidAddress, err.Error())
	}

	if m.Subscriber == nil {
		return fmt.Errorf("empty Subscriber")
	}

	return m.Subscriber.ValidateBasic()
}

// GetSignBytes implements sdk.Msg
func (m *MsgAddSubscriberRequest) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(m))
}

// GetSigners implements sdk.Msg
func (m *MsgAddSubscriberRequest) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{m.MustGetSigner()}
}

// MustGetSigner returns the signer address
func (m *MsgAddSubscriberRequest) MustGetSigner() sdk.AccAddress {
	addr, err := sdk.AccAddressFromBech32(m.Signer)
	if err != nil {
		panic(err)
	}
	return addr
}

/////////////////////////////////////
// MsgRemovePublisherIntentRequest //
/////////////////////////////////////

// NewMsgRemovePublisherIntentRequest returns a new MsgRemovePublisherIntentRequest
func NewMsgRemovePublisherIntentRequest(subscriptionID string, publisherDomain string, signer sdk.AccAddress) (*MsgRemovePublisherIntentRequest, error) {
	if signer == nil {
		return nil, fmt.Errorf("no signer provided")
	}

	return &MsgRemovePublisherIntentRequest{
		SubscriptionId:  subscriptionID,
		PublisherDomain: publisherDomain,
		Signer:          signer.String(),
	}, nil
}

// Route implements sdk.Msg
func (m *MsgRemovePublisherIntentRequest) Route() string { return ModuleName }

// Type implements sdk.Msg
func (m *MsgRemovePublisherIntentRequest) Type() string { return TypeMsgRemovePublisherIntentRequest }

// ValidateBasic implements sdk.Msg
func (m *MsgRemovePublisherIntentRequest) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(m.Signer); err != nil {
		return errorsmod.Wrap(sdkerrors.ErrInvalidAddress, err.Error())
	}

	if err := ValidateSubscriptionID(m.SubscriptionId); err != nil {
		return fmt.Errorf("invalid subscription ID: %s", err.Error())
	}

	if err := ValidateDomain(m.PublisherDomain); err != nil {
		return fmt.Errorf("invalid publisher domain: %s", err.Error())
	}

	return nil
}

// GetSignBytes implements sdk.Msg
func (m *MsgRemovePublisherIntentRequest) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(m))
}

// GetSigners implements sdk.Msg
func (m *MsgRemovePublisherIntentRequest) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{m.MustGetSigner()}
}

// MustGetSigner returns the signer address
func (m *MsgRemovePublisherIntentRequest) MustGetSigner() sdk.AccAddress {
	addr, err := sdk.AccAddressFromBech32(m.Signer)
	if err != nil {
		panic(err)
	}
	return addr
}

//////////////////////////////////////
// MsgRemoveSubscriberIntentRequest //
//////////////////////////////////////

// NewMsgRemoveSubscriberIntentRequest returns a new MsgRemoveSubscriberIntentRequest
func NewMsgRemoveSubscriberIntentRequest(subscriptionID string, subscriberAddress string, signer sdk.AccAddress) (*MsgRemoveSubscriberIntentRequest, error) {
	if signer == nil {
		return nil, fmt.Errorf("no signer provided")
	}

	return &MsgRemoveSubscriberIntentRequest{
		SubscriptionId:    subscriptionID,
		SubscriberAddress: subscriberAddress,
		Signer:            signer.String(),
	}, nil
}

// Route implements sdk.Msg
func (m *MsgRemoveSubscriberIntentRequest) Route() string { return ModuleName }

// Type implements sdk.Msg
func (m *MsgRemoveSubscriberIntentRequest) Type() string { return TypeMsgRemoveSubscriberIntentRequest }

// ValidateBasic implements sdk.Msg
func (m *MsgRemoveSubscriberIntentRequest) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(m.Signer); err != nil {
		return errorsmod.Wrap(sdkerrors.ErrInvalidAddress, err.Error())
	}

	if err := ValidateSubscriptionID(m.SubscriptionId); err != nil {
		return fmt.Errorf("invalid subscription ID: %s", err.Error())
	}

	if _, err := sdk.AccAddressFromBech32(m.SubscriberAddress); err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid subscriber address: %s", err.Error())
	}

	return nil
}

// GetSignBytes implements sdk.Msg
func (m *MsgRemoveSubscriberIntentRequest) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(m))
}

// GetSigners implements sdk.Msg
func (m *MsgRemoveSubscriberIntentRequest) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{m.MustGetSigner()}
}

// MustGetSigner returns the signer address
func (m *MsgRemoveSubscriberIntentRequest) MustGetSigner() sdk.AccAddress {
	addr, err := sdk.AccAddressFromBech32(m.Signer)
	if err != nil {
		panic(err)
	}
	return addr
}

////////////////////////////////
// MsgRemoveSubscriberRequest //
////////////////////////////////

// NewMsgRemoveSubscriberRequest returns a new MsgRemoveSubscriberRequest
func NewMsgRemoveSubscriberRequest(subscriberAddress string, signer sdk.AccAddress) (*MsgRemoveSubscriberRequest, error) {
	if signer == nil {
		return nil, fmt.Errorf("no signer provided")
	}

	return &MsgRemoveSubscriberRequest{
		SubscriberAddress: subscriberAddress,
		Signer:            signer.String(),
	}, nil
}

// Route implements sdk.Msg
func (m *MsgRemoveSubscriberRequest) Route() string { return ModuleName }

// Type implements sdk.Msg
func (m *MsgRemoveSubscriberRequest) Type() string { return TypeMsgRemoveSubscriberRequest }

// ValidateBasic implements sdk.Msg
func (m *MsgRemoveSubscriberRequest) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(m.Signer); err != nil {
		return errorsmod.Wrap(sdkerrors.ErrInvalidAddress, err.Error())
	}

	if _, err := sdk.AccAddressFromBech32(m.SubscriberAddress); err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid subscriber address: %s", err.Error())
	}

	return nil
}

// GetSignBytes implements sdk.Msg
func (m *MsgRemoveSubscriberRequest) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(m))
}

// GetSigners implements sdk.Msg
func (m *MsgRemoveSubscriberRequest) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{m.MustGetSigner()}
}

// MustGetSigner returns the signer address
func (m *MsgRemoveSubscriberRequest) MustGetSigner() sdk.AccAddress {
	addr, err := sdk.AccAddressFromBech32(m.Signer)
	if err != nil {
		panic(err)
	}
	return addr
}

///////////////////////////////
// MsgRemovePublisherRequest //
///////////////////////////////

// NewMsgRemovePublisherRequest returns a new MsgRemovePublisherRequest
func NewMsgRemovePublisherRequest(publisherDomain string, signer sdk.AccAddress) (*MsgRemovePublisherRequest, error) {
	if signer == nil {
		return nil, fmt.Errorf("no signer provided")
	}

	return &MsgRemovePublisherRequest{
		PublisherDomain: publisherDomain,
		Signer:          signer.String(),
	}, nil
}

// Route implements sdk.Msg
func (m *MsgRemovePublisherRequest) Route() string { return ModuleName }

// Type implements sdk.Msg
func (m *MsgRemovePublisherRequest) Type() string { return TypeMsgRemovePublisherRequest }

// ValidateBasic implements sdk.Msg
func (m *MsgRemovePublisherRequest) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(m.Signer); err != nil {
		return errorsmod.Wrap(sdkerrors.ErrInvalidAddress, err.Error())
	}

	if err := ValidateDomain(m.PublisherDomain); err != nil {
		return fmt.Errorf("invalid publisher domain: %s", err.Error())
	}

	return nil
}

// GetSignBytes implements sdk.Msg
func (m *MsgRemovePublisherRequest) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(m))
}

// GetSigners implements sdk.Msg
func (m *MsgRemovePublisherRequest) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{m.MustGetSigner()}
}

// MustGetSigner returns the signer address
func (m *MsgRemovePublisherRequest) MustGetSigner() sdk.AccAddress {
	addr, err := sdk.AccAddressFromBech32(m.Signer)
	if err != nil {
		panic(err)
	}
	return addr
}
