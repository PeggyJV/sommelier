package v2

import (
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	govtypesv1beta1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1beta1"
)

func RegisterInterfaces(registry codectypes.InterfaceRegistry) {
	registry.RegisterImplementations(
		(*sdk.Msg)(nil),
	)
	registry.RegisterImplementations(
		(*govtypesv1beta1.Content)(nil),
	)

	//msgservice.RegisterMsgServiceDesc(registry, &_Msg_serviceDesc)
}
