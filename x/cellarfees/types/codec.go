package types

import (
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
)

func RegisterInterfaces(registry codectypes.InterfaceRegistry) {
	registry.RegisterImplementations(
		(*sdk.Msg)(nil),
	)
	registry.RegisterImplementations(
		(*govtypes.Content)(nil),
	)

	//msgservice.RegisterMsgServiceDesc(registry, &_Msg_serviceDesc)
}
