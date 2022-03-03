package types

import (
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
)

const (
	ProposalTypeAddPublisher    = "AddPublisher"
	ProposalTypeRemovePublisher = "RemovePublisher"
)

var _ govtypes.Content = &AddPublisherProposal{}
var _ govtypes.Content = &RemovePublisherProposal{}

// TODO(bolten): fill out proposal boilerplate

//////////////////////////
// AddPublisherProposal //
//////////////////////////

func NewAddPublisherProposal(title string, description string, domain string, address string, proofUrl string, caCert string) *AddPublisherProposal {
	return &AddPublisherProposal{
		Title:       title,
		Description: description,
		Domain:      domain,
		Address:     address,
		ProofUrl:    proofUrl,
		CaCert:      caCert,
	}
}

func (p *AddPublisherProposal) ProposalRoute() string {
	return RouterKey
}

func (p *AddPublisherProposal) ProposalType() string {
	return ProposalTypeAddPublisher
}

func (p *AddPublisherProposal) ValidateBasic() error {
	if err := govtypes.ValidateAbstract(p); err != nil {
		return err
	}

	publisher := Publisher{
		Domain:   p.Domain,
		Address:  p.Address,
		ProofUrl: p.ProofUrl,
		CaCert:   p.CaCert,
	}

	return publisher.ValidateBasic()
}

/////////////////////////////
// RemovePublisherProposal //
/////////////////////////////

func NewRemovePublisherProposal(title string, description string, domain string) *RemovePublisherProposal {
	return &RemovePublisherProposal{
		Title:       title,
		Description: description,
		Domain:      domain,
	}
}

func (p *RemovePublisherProposal) ProposalRoute() string {
	return RouterKey
}

func (p *RemovePublisherProposal) ProposalType() string {
	return ProposalTypeRemovePublisher
}

func (p *RemovePublisherProposal) ValidateBasic() error {
	if err := govtypes.ValidateAbstract(p); err != nil {
		return err
	}

	return ValidateDomain(p.Domain)
}
