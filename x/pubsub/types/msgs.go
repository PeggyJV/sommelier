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

////////////////////////////
// MsgAddSubscriberIntent //
////////////////////////////

//////////////////////
// MsgAddSubscriber //
//////////////////////

//////////////////////////////
// MsgRemovePublisherIntent //
//////////////////////////////

///////////////////////////////
// MsgRemoveSubscriberIntent //
///////////////////////////////

/////////////////////////
// MsgRemoveSubscriber //
/////////////////////////

////////////////////////
// MsgRemovePublisher //
////////////////////////
