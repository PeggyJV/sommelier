package types

import (
	"encoding/json"
	"fmt"

	errorsmod "cosmossdk.io/errors"
	govtypesv1beta1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1beta1"
	"github.com/ethereum/go-ethereum/common"
	pubsubtypes "github.com/peggyjv/sommelier/v7/x/pubsub/types"
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

func NewAddManagedCellarIDsProposal(title string, description string, cellarIds *CellarIDSet, publisherDomain string) *AddManagedCellarIDsProposal {
	return &AddManagedCellarIDsProposal{
		Title:           title,
		Description:     description,
		CellarIds:       cellarIds,
		PublisherDomain: publisherDomain,
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

	if err := pubsubtypes.ValidateDomain(m.PublisherDomain); err != nil {
		return err
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

	if m.BlockHeight == 0 {
		return fmt.Errorf("block height must be non-zero")
	}

	if len(m.ContractCallProtoJson) == 0 {
		return errorsmod.Wrapf(ErrInvalidJSON, "cannot have empty contract call")
	}

	if !json.Valid([]byte(m.ContractCallProtoJson)) {
		return errorsmod.Wrapf(ErrInvalidJSON, "%s", m.ContractCallProtoJson)
	}

	if !common.IsHexAddress(m.TargetContractAddress) {
		return errorsmod.Wrapf(ErrInvalidEthereumAddress, "%s", m.TargetContractAddress)
	}

	return nil
}
