package types

import (
	"encoding/json"
	"fmt"

	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	govtypesv1beta1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1beta1"
	"github.com/ethereum/go-ethereum/common"
)

const (
	ProposalTypeAddManagedCellarIDs    = "AddManagedCellarIDs"
	ProposalTypeRemoveManagedCellarIDs = "RemoveManagedCellarIDs"
	ProposalTypeScheduledCork          = "ScheduledCork"
)

var _ govtypesv1beta1.Content = &AddManagedCellarIDsProposal{}
var _ govtypesv1beta1.Content = &RemoveManagedCellarIDsProposal{}
var _ govtypesv1beta1.Content = &ScheduledCorkProposal{}

func init() {
	// The RegisterProposalTypeCodec function was mysteriously removed by in 0.46.0 even though
	// the claim was that the old API would be preserved in .../x/gov/types/v1beta1 so we have
	// to interact with the codec directly.
	//
	// The PR that removed it: https://github.com/cosmos/cosmos-sdk/pull/11240
	// This PR was later reverted, but RegisterProposalTypeCodec was still left out. Not sure if
	// this was intentional or not.
	govtypesv1beta1.RegisterProposalType(ProposalTypeAddManagedCellarIDs)
	govtypesv1beta1.ModuleCdc.RegisterConcrete(&AddManagedCellarIDsProposal{}, "sommelier/AddManagedCellarIDsProposal", nil)

	govtypesv1beta1.RegisterProposalType(ProposalTypeRemoveManagedCellarIDs)
	govtypesv1beta1.ModuleCdc.RegisterConcrete(&RemoveManagedCellarIDsProposal{}, "sommelier/RemoveManagedCellarIDsProposal", nil)

	govtypesv1beta1.RegisterProposalType(ProposalTypeScheduledCork)
	govtypesv1beta1.ModuleCdc.RegisterConcrete(&ScheduledCorkProposal{}, "sommelier/ScheduledCorkProposal", nil)
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
	if err := govtypesv1beta1.ValidateAbstract(m); err != nil {
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
	if err := govtypesv1beta1.ValidateAbstract(m); err != nil {
		return err
	}

	if len(m.CellarIds.Ids) == 0 {
		return fmt.Errorf("can't have a remove prosoposal with no cellars")
	}

	return nil
}

func NewScheduledCorkProposal(title string, description string, blockHeight uint64, targetContractAddress string, contractCallProtoJSON string) *ScheduledCorkProposal {
	return &ScheduledCorkProposal{
		Title:                 title,
		Description:           description,
		BlockHeight:           blockHeight,
		TargetContractAddress: targetContractAddress,
		ContractCallProtoJson: contractCallProtoJSON,
	}
}

func (m *ScheduledCorkProposal) ProposalRoute() string {
	return RouterKey
}

func (m *ScheduledCorkProposal) ProposalType() string {
	return ProposalTypeScheduledCork
}

func (m *ScheduledCorkProposal) ValidateBasic() error {
	if err := govtypesv1beta1.ValidateAbstract(m); err != nil {
		return err
	}

	if len(m.ContractCallProtoJson) == 0 {
		return sdkerrors.Wrapf(ErrInvalidJSON, "cannot have empty contract call")
	}

	if !json.Valid([]byte(m.ContractCallProtoJson)) {
		return sdkerrors.Wrapf(ErrInvalidJSON, "%s", m.ContractCallProtoJson)
	}

	if !common.IsHexAddress(m.TargetContractAddress) {
		return sdkerrors.Wrapf(ErrInvalidEthereumAddress, "%s", m.TargetContractAddress)
	}

	return nil
}
