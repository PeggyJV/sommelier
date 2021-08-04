package types

import (
	"fmt"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
)

const (
	ProposalTypeManagedCellarsUpdate = "ManagedCellarsUpdate"
)

var _ govtypes.Content = &ManagedCellarsUpdateProposal{}

func NewManagedCellarsUpdateProposal(cellars []*Cellar) *ManagedCellarsUpdateProposal {
	return &ManagedCellarsUpdateProposal{Cellars: cellars}
}

func (m *ManagedCellarsUpdateProposal) ProposalRoute() string {
	return RouterKey
}

func (m *ManagedCellarsUpdateProposal) ProposalType() string {
	return ProposalTypeManagedCellarsUpdate
}

func (m *ManagedCellarsUpdateProposal) ValidateBasic() error {
	if err := govtypes.ValidateAbstract(m); err != nil {
		return err
	}

	if len(m.Cellars) == 0 {
		return fmt.Errorf("can't have a prosoposal with no cellars")
	}

	for _, c := range m.Cellars {
		if err := c.ValidateBasic(); err != nil {
			return err
		}
	}

	return nil
}