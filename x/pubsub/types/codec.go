package types

import (
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/msgservice"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
)

func RegisterInterfaces(registry codectypes.InterfaceRegistry) {
	registry.RegisterImplementations(
		(*sdk.Msg)(nil),
		&MsgAddPublisherIntent{},
		&MsgAddSubscriberIntent{},
		&MsgAddSubscriber{},
		&MsgRemovePublisherIntent{},
		&MsgRemoveSubscriberIntent{},
		&MsgRemoveSubscriber{},
		&MsgRemovePublisher{},
	)
	registry.RegisterImplementations(
		(*govtypes.Content)(nil),
		&AddPublisherProposal{},
		&RemovePublisherProposal{},
	)

	msgservice.RegisterMsgServiceDesc(registry, &_Msg_serviceDesc)
}
