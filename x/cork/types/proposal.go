package types

import (
	"fmt"

	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
)

const (
	ProposalTypeAddManagedCellarIDs    = "AddManagedCellarIDs"
	ProposalTypeRemoveManagedCellarIDs = "RemoveManagedCellarIDs"
	ProposalTypeScheduledCork          = "ScheduledCork"
)

var _ govtypes.Content = &AddManagedCellarIDsProposal{}
var _ govtypes.Content = &RemoveManagedCellarIDsProposal{}
var _ govtypes.Content = &ScheduledCorkProposal{}

func init() {
	govtypes.RegisterProposalType(ProposalTypeAddManagedCellarIDs)
	govtypes.RegisterProposalTypeCodec(&AddManagedCellarIDsProposal{}, "sommelier/AddManagedCellarIDsProposal")

	govtypes.RegisterProposalType(ProposalTypeRemoveManagedCellarIDs)
	govtypes.RegisterProposalTypeCodec(&RemoveManagedCellarIDsProposal{}, "sommelier/RemoveManagedCellarIDsProposal")

	govtypes.RegisterProposalType(ProposalTypeScheduledCork)
	govtypes.RegisterProposalTypeCodec(&ScheduledCorkProposal{}, "sommelier/cheduledCorkProposal")
}

func NewAddManagedCellarIDsProposal(title string, description string, cellarIds *CellarIDSet) *AddManagedCellarIDsProposal {
	return &AddManagedCellarIDsProposal{
		Title:       title,
		Description: description,
		CellarIds:   cellarIds,
	}
}

func (m *AddManagedCellarIDsProposal) ProposalRoute() string {
	return RouterKey
}

func (m *AddManagedCellarIDsProposal) ProposalType() string {
	return ProposalTypeAddManagedCellarIDs
}

func (m *AddManagedCellarIDsProposal) ValidateBasic() error {
	if err := govtypes.ValidateAbstract(m); err != nil {
		return err
	}

	if len(m.CellarIds.Ids) == 0 {
		return fmt.Errorf("can't have an add prosoposal with no cellars")
	}

	return nil
}

func NewRemoveManagedCellarIDsProposal(title string, description string, cellarIds *CellarIDSet) *RemoveManagedCellarIDsProposal {
	return &RemoveManagedCellarIDsProposal{
		Title:       title,
		Description: description,
		CellarIds:   cellarIds,
	}
}

func (m *RemoveManagedCellarIDsProposal) ProposalRoute() string {
	return RouterKey
}

func (m *RemoveManagedCellarIDsProposal) ProposalType() string {
	return ProposalTypeRemoveManagedCellarIDs
}

func (m *RemoveManagedCellarIDsProposal) ValidateBasic() error {
	if err := govtypes.ValidateAbstract(m); err != nil {
		return err
	}

	if len(m.CellarIds.Ids) == 0 {
		return fmt.Errorf("can't have a remove prosoposal with no cellars")
	}

	return nil
}

func NewScheduledCorkProposal(title string, description string, block_height uint64, target_contract_address string, contract_call_proto_json string) *ScheduledCorkProposal {
	return &ScheduledCorkProposal{
		Title:                 title,
		Description:           description,
		BlockHeight:           block_height,
		TargetContractAddress: target_contract_address,
		ContractCallProtoJson: contract_call_proto_json,
	}
}

func (m *ScheduledCorkProposal) ProposalRoute() string {
	return RouterKey
}

func (m *ScheduledCorkProposal) ProposalType() string {
	return ProposalTypeScheduledCork
}

func (m *ScheduledCorkProposal) ValidateBasic() error {
	if err := govtypes.ValidateAbstract(m); err != nil {
		return err
	}

	if len(m.ContractCallProtoJson) == 0 {
		return fmt.Errorf("can't have an empty command")
	}

	if len(m.TargetContractAddress) == 0 {
		return fmt.Errorf("can't have an empty contract address")
	}

	return nil
}
