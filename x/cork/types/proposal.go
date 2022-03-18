package types

import (
	"fmt"

	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
)

const (
	ProposalTypeAddManagedCellarIDs    = "AddManagedCellarIDs"
	ProposalTypeRemoveManagedCellarIDs = "RemoveManagedCellarIDs"
)

var _ govtypes.Content = &AddManagedCellarIDsProposal{}
var _ govtypes.Content = &RemoveManagedCellarIDsProposal{}

func init() {
	govtypes.RegisterProposalType(ProposalTypeAddManagedCellarIDs)
	govtypes.RegisterProposalTypeCodec(&AddManagedCellarIDsProposal{}, "sommelier/AddManagedCellarIDsProposal")

	govtypes.RegisterProposalType(ProposalTypeRemoveManagedCellarIDs)
	govtypes.RegisterProposalTypeCodec(&RemoveManagedCellarIDsProposal{}, "sommelier/RemoveManagedCellarIDsProposal")

}

func NewAddManagedCellarIDsProposal(title string, description string, cellarIds []string) *AddManagedCellarIDsProposal {
	return &AddManagedCellarIDsProposal{
		Title:       title,
		Description: description,
		CellarIds:   &CellarIDSet{Ids: cellarIds},
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

func NewRemoveManagedCellarIDsProposal(title string, description string, cellarIds []string) *RemoveManagedCellarIDsProposal {
	return &RemoveManagedCellarIDsProposal{
		Title:       title,
		Description: description,
		CellarIds:   &CellarIDSet{Ids: cellarIds},
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
