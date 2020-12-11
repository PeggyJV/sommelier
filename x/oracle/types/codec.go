package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/msgservice"
)

// ModuleCdc module codec
var ModuleCdc = codec.NewLegacyAmino()

// RegisterLegacyAminoCodec registers concrete types on codec codec
func RegisterLegacyAminoCodec(cdc *codec.LegacyAmino) {
	cdc.RegisterConcrete(MsgDelegateFeedConsent{}, "oracle/MsgDelegateFeedConsent", nil)
	cdc.RegisterConcrete(MsgAggregateExchangeRatePrevote{}, "oracle/MsgAggregateExchangeRatePrevote", nil)
	cdc.RegisterConcrete(MsgAggregateExchangeRateVote{}, "oracle/MsgAggregateExchangeRateVote", nil)
}

// RegisterInterfaces registers the x/staking interfaces types with the interface registry
func RegisterInterfaces(registry types.InterfaceRegistry) {
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgDelegateFeedConsent{},
		&MsgAggregateExchangeRatePrevote{},
		&MsgAggregateExchangeRateVote{},
	)

	msgservice.RegisterMsgServiceDesc(registry, &_Msg_serviceDesc)
}

func init() {
	RegisterLegacyAminoCodec(ModuleCdc)
}
