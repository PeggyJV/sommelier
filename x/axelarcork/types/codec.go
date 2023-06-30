package types

import (
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/msgservice"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
)

// RegisterInterfaces registers the cork proto files
func RegisterInterfaces(registry codectypes.InterfaceRegistry) {
	registry.RegisterImplementations(
		(*sdk.Msg)(nil),
		&MsgScheduleAxelarCorkRequest{},
		&MsgRelayAxelarCorkRequest{},
		&MsgBumpAxelarCorkGasRequest{},
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
