package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
)

// ModuleCdc module codec
var ModuleCdc = codec.NewLegacyAmino()

// RegisterCodec registers concrete types on codec codec
func RegisterCodec(cdc *codec.LegacyAmino) {
	cdc.RegisterConcrete(MsgDelegateFeedConsent{}, "oracle/MsgDelegateFeedConsent", nil)
	cdc.RegisterConcrete(MsgAggregateExchangeRatePrevote{}, "oracle/MsgAggregateExchangeRatePrevote", nil)
	cdc.RegisterConcrete(MsgAggregateExchangeRateVote{}, "oracle/MsgAggregateExchangeRateVote", nil)
}

func init() {
	RegisterCodec(ModuleCdc)
}
