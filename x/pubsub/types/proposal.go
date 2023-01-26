package types

import (
	fmt "fmt"
	"net/url"

	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
)

const (
	ProposalTypeAddPublisher              = "AddPublisher"
	ProposalTypeRemovePublisher           = "RemovePublisher"
	ProposalTypeAddDefaultSubscription    = "AddDefaultSubscription"
	ProposalTypeRemoveDefaultSubscription = "RemoveDefaultSubscription"
)

var _ govtypes.Content = &AddPublisherProposal{}
var _ govtypes.Content = &RemovePublisherProposal{}
var _ govtypes.Content = &AddDefaultSubscriptionProposal{}
var _ govtypes.Content = &RemoveDefaultSubscriptionProposal{}

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
		Domain:  p.Domain,
		Address: p.Address,
		CaCert:  p.CaCert,
	}

	if err := publisher.ValidateBasic(); err != nil {
		return err
	}

	return ValidateProofUrl(p.ProofUrl, p.Domain, p.Address)
}

func ValidateProofUrl(proofUrl string, domain string, address string) error {
	if proofUrl == "" {
		return fmt.Errorf("empty proof URL")
	}

	if err := ValidateGenericUrl(proofUrl); err != nil {
		return err
	}

	validProofUrl := fmt.Sprintf("https://%s/%s/cacert.pem", domain, address)
	if proofUrl != validProofUrl {
		return fmt.Errorf("invalid proof URL format, should be: %s", validProofUrl)
	}

	return nil
}

func ValidateGenericUrl(urlString string) error {
	if urlString == "" {
		return fmt.Errorf("empty URL")
	}

	if len(urlString) > MaxUrlLength {
		return fmt.Errorf("URL over max length of %d: %d", MaxUrlLength, len(urlString))
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
	if err := govtypes.ValidateAbstract(p); err != nil {
		return err
	}

	return ValidateDomain(p.Domain)
}

////////////////////////////////////
// AddDefaultSubscriptionProposal //
///////////////////////////////////

func NewAddDefaultSubscriptionProposal(title string, description string, subscriptionId string, publisherDomain string) *AddDefaultSubscriptionProposal {
	return &AddDefaultSubscriptionProposal{
		Title:           title,
		Description:     description,
		SubscriptionId:  subscriptionId,
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
	if err := govtypes.ValidateAbstract(p); err != nil {
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

func NewRemoveDefaultSubscriptionProposal(title string, description string, subscriptionId string) *RemoveDefaultSubscriptionProposal {
	return &RemoveDefaultSubscriptionProposal{
		Title:          title,
		Description:    description,
		SubscriptionId: subscriptionId,
	}
}

func (p *RemoveDefaultSubscriptionProposal) ProposalRoute() string {
	return RouterKey
}

func (p *RemoveDefaultSubscriptionProposal) ProposalType() string {
	return ProposalTypeRemoveDefaultSubscription
}

func (p *RemoveDefaultSubscriptionProposal) ValidateBasic() error {
	if err := govtypes.ValidateAbstract(p); err != nil {
		return err
	}

	return ValidateSubscriptionId(p.SubscriptionId)
}
