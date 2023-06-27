package types

import (
	"encoding/json"
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"

	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	"github.com/ethereum/go-ethereum/common"
)

const (
	ProposalTypeAddManagedCellarIDs = "AddAxelarManagedCellarIDs"

	ProposalTypeRemoveManagedCellarIDs   = "RemoveAxelarManagedCellarIDs"
	ProposalTypeScheduledCork            = "AxelarScheduledCork"
	ProposalTypeCommunitySpend           = "AxelarCommunitySpend"
	ProposalTypeAddChainConfiguration    = "AddAxelarChainConfiguration"
	ProposalTypeRemoveChainConfiguration = "RemoveAxelarChainConfiguration"
)

var _ govtypes.Content = &AddAxelarManagedCellarIDsProposal{}
var _ govtypes.Content = &RemoveAxelarManagedCellarIDsProposal{}
var _ govtypes.Content = &AxelarScheduledCorkProposal{}
var _ govtypes.Content = &AxelarCommunityPoolSpendProposal{}
var _ govtypes.Content = &AddChainConfigurationProposal{}
var _ govtypes.Content = &RemoveChainConfigurationProposal{}

func init() {
	govtypes.RegisterProposalType(ProposalTypeAddManagedCellarIDs)
	govtypes.RegisterProposalTypeCodec(&AddAxelarManagedCellarIDsProposal{}, "sommelier/AddAxelarManagedCellarIDsProposal")

	govtypes.RegisterProposalType(ProposalTypeRemoveManagedCellarIDs)
	govtypes.RegisterProposalTypeCodec(&RemoveAxelarManagedCellarIDsProposal{}, "sommelier/RemoveAxelarManagedCellarIDsProposal")

	govtypes.RegisterProposalType(ProposalTypeScheduledCork)
	govtypes.RegisterProposalTypeCodec(&AxelarScheduledCorkProposal{}, "sommelier/AxelarScheduledCorkProposal")

	govtypes.RegisterProposalType(ProposalTypeAddChainConfiguration)
	govtypes.RegisterProposalTypeCodec(&AddChainConfigurationProposal{}, "sommelier/AddAxelarChainConfigurationProposal")

	govtypes.RegisterProposalType(ProposalTypeRemoveChainConfiguration)
	govtypes.RegisterProposalTypeCodec(&RemoveChainConfigurationProposal{}, "sommelier/RemoveAxelarChainConfigurationProposal")

}

func NewAddManagedCellarIDsProposal(title string, description string, chainName string, chainID uint64, cellarIds *CellarIDSet) *AddAxelarManagedCellarIDsProposal {
	return &AddAxelarManagedCellarIDsProposal{
		Title:       title,
		Description: description,
		CellarIds:   cellarIds,
		ChainName:   chainName,
		ChainId:     chainID,
	}
}

func (m *AddAxelarManagedCellarIDsProposal) ProposalRoute() string {
	return RouterKey
}

func (m *AddAxelarManagedCellarIDsProposal) ProposalType() string {
	return ProposalTypeAddManagedCellarIDs
}

func (m *AddAxelarManagedCellarIDsProposal) ValidateBasic() error {
	if err := govtypes.ValidateAbstract(m); err != nil {
		return err
	}

	if len(m.CellarIds.Ids) == 0 {
		return fmt.Errorf("can't have an add prosoposal with no cellars")
	}

	return nil
}

func NewRemoveManagedCellarIDsProposal(title string, description string, chainName string, chainID uint64, cellarIds *CellarIDSet) *RemoveAxelarManagedCellarIDsProposal {
	return &RemoveAxelarManagedCellarIDsProposal{
		Title:       title,
		Description: description,
		CellarIds:   cellarIds,
		ChainName:   chainName,
		ChainId:     chainID,
	}
}

func (m *RemoveAxelarManagedCellarIDsProposal) ProposalRoute() string {
	return RouterKey
}

func (m *RemoveAxelarManagedCellarIDsProposal) ProposalType() string {
	return ProposalTypeRemoveManagedCellarIDs
}

func (m *RemoveAxelarManagedCellarIDsProposal) ValidateBasic() error {
	if err := govtypes.ValidateAbstract(m); err != nil {
		return err
	}

	if len(m.CellarIds.Ids) == 0 {
		return fmt.Errorf("can't have a remove prosoposal with no cellars")
	}

	return nil
}

func NewScheduledCorkProposal(title string, description string, blockHeight uint64, chainName string, chainID uint64, targetContractAddress string, contractCallProtoJSON string) *AxelarScheduledCorkProposal {
	return &AxelarScheduledCorkProposal{
		Title:                 title,
		Description:           description,
		BlockHeight:           blockHeight,
		ChainName:             chainName,
		ChainId:               chainID,
		TargetContractAddress: targetContractAddress,
		ContractCallProtoJson: contractCallProtoJSON,
	}
}

func (m *AxelarScheduledCorkProposal) ProposalRoute() string {
	return RouterKey
}

func (m *AxelarScheduledCorkProposal) ProposalType() string {
	return ProposalTypeScheduledCork
}

func (m *AxelarScheduledCorkProposal) ValidateBasic() error {
	if err := govtypes.ValidateAbstract(m); err != nil {
		return err
	}

	if len(m.ContractCallProtoJson) == 0 {
		return sdkerrors.Wrapf(ErrInvalidJSON, "cannot have empty contract call")
	}

	if !json.Valid([]byte(m.ContractCallProtoJson)) {
		return sdkerrors.Wrapf(ErrInvalidJSON, "%s", m.ContractCallProtoJson)
	}

	if !common.IsHexAddress(m.TargetContractAddress) {
		return sdkerrors.Wrapf(ErrInvalidEVMAddress, "%s", m.TargetContractAddress)
	}

	return nil
}

func NewCommunitySpendProposal(title string, description string, recipient string, chainID uint64, chainName string, amount sdk.Coin) *AxelarCommunityPoolSpendProposal {
	return &AxelarCommunityPoolSpendProposal{
		Title:       title,
		Description: description,
		Recipient:   recipient,
		ChainId:     chainID,
		ChainName:   chainName,
		Amount:      amount,
	}
}

func (m *AxelarCommunityPoolSpendProposal) ProposalRoute() string {
	return RouterKey
}

func (m *AxelarCommunityPoolSpendProposal) ProposalType() string {
	return ProposalTypeCommunitySpend
}

func (m *AxelarCommunityPoolSpendProposal) ValidateBasic() error {
	if err := govtypes.ValidateAbstract(m); err != nil {
		return err
	}

	if m.Amount.Amount.IsZero() {
		return ErrValuelessSend
	}

	if m.Recipient == "" {
		return sdkerrors.Wrapf(ErrInvalidEVMAddress, "empty recipient")
	}

	if !common.IsHexAddress(m.Recipient) {
		return sdkerrors.Wrapf(ErrInvalidEVMAddress, "%s", m.Recipient)
	}

	return nil
}

func NewAddChainConfigurationProposal(title string, description string, configuration ChainConfiguration) *AddChainConfigurationProposal {
	return &AddChainConfigurationProposal{
		Title:              title,
		Description:        description,
		ChainConfiguration: &configuration,
	}
}

func (m *AddChainConfigurationProposal) ProposalRoute() string {
	return RouterKey
}

func (m *AddChainConfigurationProposal) ProposalType() string {
	return ProposalTypeAddChainConfiguration
}

func (m *AddChainConfigurationProposal) ValidateBasic() error {
	if err := govtypes.ValidateAbstract(m); err != nil {
		return err
	}

	if err := m.ChainConfiguration.ValidateBasic(); err != nil {
		return err
	}

	return nil
}

func NewRemoveChainConfigurationProposal(title string, description string, chainID uint64) *RemoveChainConfigurationProposal {
	return &RemoveChainConfigurationProposal{
		Title:       title,
		Description: description,
		ChainId:     chainID,
	}
}

func (m *RemoveChainConfigurationProposal) ProposalRoute() string {
	return RouterKey
}

func (m *RemoveChainConfigurationProposal) ProposalType() string {
	return ProposalTypeRemoveChainConfiguration
}

func (m *RemoveChainConfigurationProposal) ValidateBasic() error {
	if err := govtypes.ValidateAbstract(m); err != nil {
		return err
	}

	return nil
}
