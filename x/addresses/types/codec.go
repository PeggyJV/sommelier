package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"

	"github.com/cosmos/cosmos-sdk/types/msgservice"
)

func RegisterCodec(cdc *codec.LegacyAmino) {
	cdc.RegisterConcrete(&MsgAddAddressMapping{}, "addresses/MsgAddAddressMapping", nil)
	cdc.RegisterConcrete(&MsgRemoveAddressMapping{}, "addresses/MsgRemoveAddressMapping", nil)
}

func RegisterInterfaces(registry cdctypes.InterfaceRegistry) {
	msgservice.RegisterMsgServiceDesc(registry, &_Msg_serviceDesc)
}

var (
	Amino     = codec.NewLegacyAmino()
	ModuleCdc = codec.NewProtoCodec(cdctypes.NewInterfaceRegistry())
)
