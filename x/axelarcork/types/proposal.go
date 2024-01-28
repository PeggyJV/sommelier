package types

import (
	"encoding/json"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"

	errorsmod "cosmossdk.io/errors"
	govtypesv1beta1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1beta1"
	"github.com/ethereum/go-ethereum/common"
	pubsubtypes "github.com/peggyjv/sommelier/v7/x/pubsub/types"
)

const (
	ProposalTypeAddManagedCellarIDs = "AddAxelarManagedCellarIDs"

	ProposalTypeRemoveManagedCellarIDs           = "RemoveAxelarManagedCellarIDs"
	ProposalTypeScheduledCork                    = "AxelarScheduledCork"
	ProposalTypeCommunitySpend                   = "AxelarCommunitySpend"
	ProposalTypeAddChainConfiguration            = "AddAxelarChainConfiguration"
	ProposalTypeRemoveChainConfiguration         = "RemoveAxelarChainConfiguration"
	ProposalTypeUpgradeAxelarProxyContract       = "UpgradeAxelarProxyContract"
	ProposalTypeCancelAxelarProxyContractUpgrade = "CancelAxelarProxyContractUpgrade"
)

var _ govtypesv1beta1.Content = &AddAxelarManagedCellarIDsProposal{}
var _ govtypesv1beta1.Content = &RemoveAxelarManagedCellarIDsProposal{}
var _ govtypesv1beta1.Content = &AxelarScheduledCorkProposal{}
var _ govtypesv1beta1.Content = &AxelarCommunityPoolSpendProposal{}
var _ govtypesv1beta1.Content = &AddChainConfigurationProposal{}
var _ govtypesv1beta1.Content = &RemoveChainConfigurationProposal{}
var _ govtypesv1beta1.Content = &UpgradeAxelarProxyContractProposal{}
var _ govtypesv1beta1.Content = &CancelAxelarProxyContractUpgradeProposal{}

func init() {
	// The RegisterProposalTypeCodec function was mysteriously removed by in 0.46.0 even though
	// the claim was that the old API would be preserved in .../x/gov/types/v1beta1 so we have
	// to interact with the codec directly.
	//
	// The PR that removed it: https://github.com/cosmos/cosmos-sdk/pull/11240
	// This PR was later reverted, but RegisterProposalTypeCodec was still left out. Not sure if
	// this was intentional or not.
	govtypesv1beta1.RegisterProposalType(ProposalTypeAddManagedCellarIDs)
	govtypesv1beta1.ModuleCdc.RegisterConcrete(&AddAxelarManagedCellarIDsProposal{}, "sommelier/AddAxelarManagedCellarIDsProposal", nil)

	govtypesv1beta1.RegisterProposalType(ProposalTypeRemoveManagedCellarIDs)
	govtypesv1beta1.ModuleCdc.RegisterConcrete(&RemoveAxelarManagedCellarIDsProposal{}, "sommelier/RemoveAxelarManagedCellarIDsProposal", nil)

	govtypesv1beta1.RegisterProposalType(ProposalTypeScheduledCork)
	govtypesv1beta1.ModuleCdc.RegisterConcrete(&AxelarScheduledCorkProposal{}, "sommelier/AxelarScheduledCorkProposal", nil)

	govtypesv1beta1.RegisterProposalType(ProposalTypeAddChainConfiguration)
	govtypesv1beta1.ModuleCdc.RegisterConcrete(&AddChainConfigurationProposal{}, "sommelier/AddAxelarChainConfigurationProposal", nil)

	govtypesv1beta1.RegisterProposalType(ProposalTypeRemoveChainConfiguration)
	govtypesv1beta1.ModuleCdc.RegisterConcrete(&RemoveChainConfigurationProposal{}, "sommelier/RemoveAxelarChainConfigurationProposal", nil)

	govtypesv1beta1.RegisterProposalType(ProposalTypeCommunitySpend)
	govtypesv1beta1.ModuleCdc.RegisterConcrete(&AxelarCommunityPoolSpendProposal{}, "sommelier/AxelarCommunitySpendProposal", nil)

	govtypesv1beta1.RegisterProposalType(ProposalTypeUpgradeAxelarProxyContract)
	govtypesv1beta1.ModuleCdc.RegisterConcrete(&UpgradeAxelarProxyContractProposal{}, "sommelier/UpgradeAxelarProxyContractProposal", nil)

	govtypesv1beta1.RegisterProposalType(ProposalTypeCancelAxelarProxyContractUpgrade)
	govtypesv1beta1.ModuleCdc.RegisterConcrete(&CancelAxelarProxyContractUpgradeProposal{}, "sommelier/CancelAxelarProxyContractUpgradeProposal", nil)
}

func NewAddAxelarManagedCellarIDsProposal(title string, description string, chainID uint64, cellarIds *CellarIDSet, publisherDomain string) *AddAxelarManagedCellarIDsProposal {
	return &AddAxelarManagedCellarIDsProposal{
		Title:           title,
		Description:     description,
		CellarIds:       cellarIds,
		ChainId:         chainID,
		PublisherDomain: publisherDomain,
	}
}

func (m *AddAxelarManagedCellarIDsProposal) ProposalRoute() string {
	return RouterKey
}

func (m *AddAxelarManagedCellarIDsProposal) ProposalType() string {
	return ProposalTypeAddManagedCellarIDs
}

func (m *AddAxelarManagedCellarIDsProposal) ValidateBasic() error {
	if err := govtypesv1beta1.ValidateAbstract(m); err != nil {
		return err
	}

	if err := m.CellarIds.ValidateBasic(); err != nil {
		return err
	}

	if m.ChainId == 0 {
		return fmt.Errorf("chain ID must be non-zero")
	}

	if err := pubsubtypes.ValidateDomain(m.PublisherDomain); err != nil {
		return err
	}

	return nil
}

func NewRemoveAxelarManagedCellarIDsProposal(title string, description string, chainID uint64, cellarIds *CellarIDSet) *RemoveAxelarManagedCellarIDsProposal {
	return &RemoveAxelarManagedCellarIDsProposal{
		Title:       title,
		Description: description,
		CellarIds:   cellarIds,
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
	if err := govtypesv1beta1.ValidateAbstract(m); err != nil {
		return err
	}

	if err := m.CellarIds.ValidateBasic(); err != nil {
		return err
	}

	if m.ChainId == 0 {
		return fmt.Errorf("chain ID must be non-zero")
	}

	return nil
}

func NewAxelarScheduledCorkProposal(title string, description string, blockHeight uint64, chainID uint64, targetContractAddress string, contractCallProtoJSON string, deadline uint64) *AxelarScheduledCorkProposal {
	return &AxelarScheduledCorkProposal{
		Title:                 title,
		Description:           description,
		BlockHeight:           blockHeight,
		ChainId:               chainID,
		TargetContractAddress: targetContractAddress,
		ContractCallProtoJson: contractCallProtoJSON,
		Deadline:              deadline,
	}
}

func (m *AxelarScheduledCorkProposal) ProposalRoute() string {
	return RouterKey
}

func (m *AxelarScheduledCorkProposal) ProposalType() string {
	return ProposalTypeScheduledCork
}

func (m *AxelarScheduledCorkProposal) ValidateBasic() error {
	if err := govtypesv1beta1.ValidateAbstract(m); err != nil {
		return err
	}

	if m.BlockHeight == 0 {
		return fmt.Errorf("block height must be non-zero")
	}

	if m.ChainId == 0 {
		return fmt.Errorf("chain ID must be non-zero")
	}

	if !common.IsHexAddress(m.TargetContractAddress) {
		return errorsmod.Wrapf(ErrInvalidEVMAddress, "%s", m.TargetContractAddress)
	}

	if len(m.ContractCallProtoJson) == 0 {
		return errorsmod.Wrapf(ErrInvalidJSON, "cannot have empty contract call")
	}

	if !json.Valid([]byte(m.ContractCallProtoJson)) {
		return errorsmod.Wrapf(ErrInvalidJSON, "%s", m.ContractCallProtoJson)
	}

	if m.Deadline == 0 {
		return fmt.Errorf("deadline must be non-zero")
	}

	return nil
}

func NewAxelarCommunitySpendProposal(title string, description string, recipient string, chainID uint64, amount sdk.Coin) *AxelarCommunityPoolSpendProposal {
	return &AxelarCommunityPoolSpendProposal{
		Title:       title,
		Description: description,
		Recipient:   recipient,
		ChainId:     chainID,
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
	if err := govtypesv1beta1.ValidateAbstract(m); err != nil {
		return err
	}

	if !common.IsHexAddress(m.Recipient) {
		return errorsmod.Wrapf(ErrInvalidEVMAddress, "%s", m.Recipient)
	}

	if m.ChainId == 0 {
		return fmt.Errorf("chain ID must be non-zero")
	}

	if !m.Amount.IsValid() || !m.Amount.IsPositive() {
		return ErrValuelessSend
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
	if err := govtypesv1beta1.ValidateAbstract(m); err != nil {
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
	if err := govtypesv1beta1.ValidateAbstract(m); err != nil {
		return err
	}

	if m.ChainId == 0 {
		return fmt.Errorf("chain ID must be non-zero")
	}

	return nil
}

func NewUpgradeAxelarProxyContractProposal(title string, description string, chainID uint64, newProxyAddress string) *UpgradeAxelarProxyContractProposal {
	return &UpgradeAxelarProxyContractProposal{
		Title:           title,
		Description:     description,
		ChainId:         chainID,
		NewProxyAddress: newProxyAddress,
	}
}

func (m *UpgradeAxelarProxyContractProposal) ProposalRoute() string {
	return RouterKey
}

func (m *UpgradeAxelarProxyContractProposal) ProposalType() string {
	return ProposalTypeUpgradeAxelarProxyContract
}

func (m *UpgradeAxelarProxyContractProposal) ValidateBasic() error {
	if err := govtypesv1beta1.ValidateAbstract(m); err != nil {
		return err
	}

	if m.ChainId == 0 {
		return fmt.Errorf("chain ID must be non-zero")
	}

	if !common.IsHexAddress(m.NewProxyAddress) {
		return errorsmod.Wrapf(ErrInvalidEVMAddress, "%s", m.NewProxyAddress)
	}

	return nil
}

func NewCancelAxelarProxyContractUpgradeProposal(title string, description string, chainID uint64) *CancelAxelarProxyContractUpgradeProposal {
	return &CancelAxelarProxyContractUpgradeProposal{
		Title:       title,
		Description: description,
		ChainId:     chainID,
	}
}

func (m *CancelAxelarProxyContractUpgradeProposal) ProposalRoute() string {
	return RouterKey
}

func (m *CancelAxelarProxyContractUpgradeProposal) ProposalType() string {
	return ProposalTypeCancelAxelarProxyContractUpgrade
}

func (m *CancelAxelarProxyContractUpgradeProposal) ValidateBasic() error {
	if err := govtypesv1beta1.ValidateAbstract(m); err != nil {
		return err
	}

	if m.ChainId == 0 {
		return fmt.Errorf("chain ID must be non-zero")
	}

	return nil
}
