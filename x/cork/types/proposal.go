package types

import (
	"fmt"

	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
)

const (
	ProposalTypeAddManagedCellars    = "AddManagedCellars"
	ProposalTypeRemoveManagedCellars = "RemoveManagedCellars"
)

var _ govtypes.Content = &AddManagedCellarsProposal{}
var _ govtypes.Content = &RemoveManagedCellarsProposal{}

func init() {
	govtypes.RegisterProposalType(ProposalTypeAddManagedCellars)
	govtypes.RegisterProposalTypeCodec(&AddManagedCellarsProposal{}, "sommelier/AddManagedCellarsProposal")

	govtypes.RegisterProposalType(ProposalTypeRemoveManagedCellars)
	govtypes.RegisterProposalTypeCodec(&RemoveManagedCellarsProposal{}, "sommelier/RemoveManagedCellarsProposal")

}

func NewAddManagedCellarsProposal(title string, description string, cellarIds []string) *AddManagedCellarsProposal {
	return &AddManagedCellarsProposal{
		Title:       title,
		Description: description,
		CellarIds:   &CellarIDSet{Ids: cellarIds},
	}
}

func (m *AddManagedCellarsProposal) ProposalRoute() string {
	return RouterKey
}

func (m *AddManagedCellarsProposal) ProposalType() string {
	return ProposalTypeAddManagedCellars
}

func (m *AddManagedCellarsProposal) ValidateBasic() error {
	if err := govtypes.ValidateAbstract(m); err != nil {
		return err
	}

	if len(m.CellarIds.Ids) == 0 {
		return fmt.Errorf("can't have an add prosoposal with no cellars")
	}

	return nil
}

func NewRemoveManagedCellarsProposal(title string, description string, cellarIds []string) *RemoveManagedCellarsProposal {
	return &RemoveManagedCellarsProposal{
		Title:       title,
		Description: description,
		CellarIds:   &CellarIDSet{Ids: cellarIds},
	}
}

func (m *RemoveManagedCellarsProposal) ProposalRoute() string {
	return RouterKey
}

func (m *RemoveManagedCellarsProposal) ProposalType() string {
	return ProposalTypeRemoveManagedCellars
}

func (m *RemoveManagedCellarsProposal) ValidateBasic() error {
	if err := govtypes.ValidateAbstract(m); err != nil {
		return err
	}

	if len(m.CellarIds.Ids) == 0 {
		return fmt.Errorf("can't have a remove prosoposal with no cellars")
	}

	return nil
}
