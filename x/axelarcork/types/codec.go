package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	cryptocodec "github.com/cosmos/cosmos-sdk/crypto/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/msgservice"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
)

// RegisterLegacyAminoCodec registers the vesting interfaces and concrete types on the
// provided LegacyAmino codec. These types are used for Amino JSON serialization
func RegisterLegacyAminoCodec(cdc *codec.LegacyAmino) {
	cdc.RegisterConcrete(&MsgScheduleAxelarCorkRequest{}, "axelarcork/MsgScheduleAxelarCorkRequest", nil)
	cdc.RegisterConcrete(&MsgRelayAxelarCorkRequest{}, "axelarcork/MsgRelayAxelarCorkRequest", nil)
	cdc.RegisterConcrete(&MsgBumpAxelarCorkGasRequest{}, "axelarcork/MsgBumpAxelarCorkGasRequest", nil)
}

var (
	amino = codec.NewLegacyAmino()
	// ModuleCdc Note, the codec should ONLY be used in certain instances of tests and for
	// JSON encoding as Amino is still used for that purpose.
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

// RegisterInterfaces registers the cork proto files
func RegisterInterfaces(registry codectypes.InterfaceRegistry) {
	registry.RegisterImplementations(
		(*sdk.Msg)(nil),
		&MsgScheduleAxelarCorkRequest{},
		&MsgRelayAxelarCorkRequest{},
		&MsgBumpAxelarCorkGasRequest{},
		&MsgCancelAxelarCorkRequest{},
	)

	registry.RegisterImplementations((*govtypes.Content)(nil),
		&AxelarScheduledCorkProposal{},
		&AddChainConfigurationProposal{},
		&RemoveChainConfigurationProposal{},
		&AxelarCommunityPoolSpendProposal{},
		&AddAxelarManagedCellarIDsProposal{},
		&RemoveAxelarManagedCellarIDsProposal{},
	)

	msgservice.RegisterMsgServiceDesc(registry, &_Msg_serviceDesc)
}
