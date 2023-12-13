package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	cryptocodec "github.com/cosmos/cosmos-sdk/crypto/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/msgservice"
	govtypesv1beta1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1beta1"
)

// RegisterLegacyAminoCodec registers the vesting interfaces and concrete types on the
// provided LegacyAmino codec. These types are used for Amino JSON serialization
func RegisterLegacyAminoCodec(cdc *codec.LegacyAmino) {
	cdc.RegisterConcrete(&MsgAddPublisherIntentRequest{}, "pubsub/MsgAddPublisherIntentRequest", nil)
	cdc.RegisterConcrete(&MsgAddSubscriberIntentRequest{}, "pubsub/MsgAddSubscriberIntentRequest", nil)
	cdc.RegisterConcrete(&MsgAddSubscriberRequest{}, "pubsub/MsgAddSubscriberRequest", nil)
	cdc.RegisterConcrete(&MsgRemovePublisherIntentRequest{}, "pubsub/MsgRemovePublisherIntentRequest", nil)
	cdc.RegisterConcrete(&MsgRemoveSubscriberIntentRequest{}, "pubsub/MsgRemoveSubscriberIntentRequest", nil)
	cdc.RegisterConcrete(&MsgRemoveSubscriberRequest{}, "pubsub/MsgRemoveSubscriberRequest", nil)
	cdc.RegisterConcrete(&MsgRemovePublisherRequest{}, "pubsub/MsgRemovePublisherRequest", nil)
}

var (
	amino = codec.NewLegacyAmino()

	// ModuleCdc references the global x/bank module codec. Note, the codec should
	// ONLY be used in certain instances of tests and for JSON encoding as Amino is
	// still used for that purpose.
	//
	// The actual codec used for serialization should be provided to x/staking and
	// defined at the application level.
	ModuleCdc = codec.NewAminoCodec(amino)
)

func init() {
	RegisterLegacyAminoCodec(amino)
	cryptocodec.RegisterCrypto(amino)
	amino.Seal()
}

func RegisterInterfaces(registry codectypes.InterfaceRegistry) {
	registry.RegisterImplementations(
		(*sdk.Msg)(nil),
		&MsgAddPublisherIntentRequest{},
		&MsgAddSubscriberIntentRequest{},
		&MsgAddSubscriberRequest{},
		&MsgRemovePublisherIntentRequest{},
		&MsgRemoveSubscriberIntentRequest{},
		&MsgRemoveSubscriberRequest{},
		&MsgRemovePublisherRequest{},
	)
	registry.RegisterImplementations(
		(*govtypesv1beta1.Content)(nil),
		&AddPublisherProposal{},
		&RemovePublisherProposal{},
		&AddDefaultSubscriptionProposal{},
		&RemoveDefaultSubscriptionProposal{},
	)

	msgservice.RegisterMsgServiceDesc(registry, &_Msg_serviceDesc)
}
