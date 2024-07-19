package types

import (
	fmt "fmt"
	"net/url"

	govtypesv1beta1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1beta1"
)

const (
	ProposalTypeAddPublisher              = "AddPublisher"
	ProposalTypeRemovePublisher           = "RemovePublisher"
	ProposalTypeAddDefaultSubscription    = "AddDefaultSubscription"
	ProposalTypeRemoveDefaultSubscription = "RemoveDefaultSubscription"
)

var _ govtypesv1beta1.Content = &AddPublisherProposal{}
var _ govtypesv1beta1.Content = &RemovePublisherProposal{}
var _ govtypesv1beta1.Content = &AddDefaultSubscriptionProposal{}
var _ govtypesv1beta1.Content = &RemoveDefaultSubscriptionProposal{}

func init() {
	// The RegisterProposalTypeCodec function was mysteriously removed by in 0.46.0 even though
	// the claim was that the old API would be preserved in .../x/gov/types/v1beta1 so we have
	// to interact with the codec directly.
	//
	// The PR that removed it: https://github.com/cosmos/cosmos-sdk/pull/11240
	// This PR was later reverted, but RegisterProposalTypeCodec was still left out. Not sure if
	// this was intentional or not.
	govtypesv1beta1.RegisterProposalType(ProposalTypeAddPublisher)

	govtypesv1beta1.RegisterProposalType(ProposalTypeRemovePublisher)

	govtypesv1beta1.RegisterProposalType(ProposalTypeAddDefaultSubscription)

	govtypesv1beta1.RegisterProposalType(ProposalTypeRemoveDefaultSubscription)
}

//////////////////////////
// AddPublisherProposal //
//////////////////////////

func NewAddPublisherProposal(title string, description string, domain string, address string, proofURL string, caCert string) *AddPublisherProposal {
	return &AddPublisherProposal{
		Title:       title,
		Description: description,
		Domain:      domain,
		Address:     address,
		ProofUrl:    proofURL,
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
	if err := govtypesv1beta1.ValidateAbstract(p); err != nil {
		return err
	}

	publisher := Publisher{
		Domain:  p.Domain,
		Address: p.Address,
		CaCert:  p.CaCert,
	}

	if err := publisher.ValidateBasic(); err != nil {
		return err
	}

	return ValidateProofURL(p.ProofUrl, p.Domain, p.Address)
}

func ValidateProofURL(proofURL string, domain string, address string) error {
	if proofURL == "" {
		return fmt.Errorf("empty proof URL")
	}

	if err := ValidateGenericURL(proofURL); err != nil {
		return err
	}

	validProofURL := fmt.Sprintf("%s/%s/cacert.pem", domain, address)
	validProofURLWithHTTPS := fmt.Sprintf("https://%s", validProofURL)
	if proofURL != validProofURL && proofURL != validProofURLWithHTTPS {
		return fmt.Errorf("invalid proof URL format, should be: %s", validProofURLWithHTTPS)
	}

	return nil
}

func ValidateGenericURL(urlString string) error {
	if urlString == "" {
		return fmt.Errorf("empty URL")
	}

	if len(urlString) > MaxURLLength {
		return fmt.Errorf("URL over max length of %d: %d", MaxURLLength, len(urlString))
	}

	if _, err := url.Parse(urlString); err != nil {
		return fmt.Errorf("invalid URL format: %s", err.Error())
	}

	return nil
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
	if err := govtypesv1beta1.ValidateAbstract(p); err != nil {
		return err
	}

	return ValidateDomain(p.Domain)
}

////////////////////////////////////
// AddDefaultSubscriptionProposal //
///////////////////////////////////

func NewAddDefaultSubscriptionProposal(title string, description string, subscriptionID string, publisherDomain string) *AddDefaultSubscriptionProposal {
	return &AddDefaultSubscriptionProposal{
		Title:           title,
		Description:     description,
		SubscriptionId:  subscriptionID,
		PublisherDomain: publisherDomain,
	}
}

func (p *AddDefaultSubscriptionProposal) ProposalRoute() string {
	return RouterKey
}

func (p *AddDefaultSubscriptionProposal) ProposalType() string {
	return ProposalTypeAddDefaultSubscription
}

func (p *AddDefaultSubscriptionProposal) ValidateBasic() error {
	if err := govtypesv1beta1.ValidateAbstract(p); err != nil {
		return err
	}

	defaultSubscription := DefaultSubscription{
		SubscriptionId:  p.SubscriptionId,
		PublisherDomain: p.PublisherDomain,
	}

	return defaultSubscription.ValidateBasic()
}

///////////////////////////////////////
// RemoveDefaultSubscriptionProposal //
///////////////////////////////////////

func NewRemoveDefaultSubscriptionProposal(title string, description string, subscriptionID string) *RemoveDefaultSubscriptionProposal {
	return &RemoveDefaultSubscriptionProposal{
		Title:          title,
		Description:    description,
		SubscriptionId: subscriptionID,
	}
}

func (p *RemoveDefaultSubscriptionProposal) ProposalRoute() string {
	return RouterKey
}

func (p *RemoveDefaultSubscriptionProposal) ProposalType() string {
	return ProposalTypeRemoveDefaultSubscription
}

func (p *RemoveDefaultSubscriptionProposal) ValidateBasic() error {
	if err := govtypesv1beta1.ValidateAbstract(p); err != nil {
		return err
	}

	return ValidateSubscriptionID(p.SubscriptionId)
}
