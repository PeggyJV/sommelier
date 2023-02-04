package rest

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/rest"
)

type (
	// AddPublisherProposalReq defines an add publisher proposal request body.
	AddPublisherProposalReq struct {
		BaseReq rest.BaseReq `json:"base_req" yaml:"base_req"`

		Title       string         `json:"title" yaml:"title"`
		Description string         `json:"description" yaml:"description"`
		Domain      string         `json:"domain" yaml:"domain"`
		Address     string         `json:"address" yaml:"address"`
		ProofURL    string         `json:"proof_url" yaml:"proof_url"`
		CaCert      string         `json:"ca_cert" yaml:"ca_cert"`
		Proposer    sdk.AccAddress `json:"proposer" yaml:"proposer"`
		Deposit     sdk.Coins      `json:"deposit" yaml:"deposit"`
	}

	// RemovePublisherProposalReq defines a remove publisher proposal request body.
	RemovePublisherProposalReq struct {
		BaseReq rest.BaseReq `json:"base_req" yaml:"base_req"`

		Title       string         `json:"title" yaml:"title"`
		Description string         `json:"description" yaml:"description"`
		Domain      string         `json:"domain" yaml:"domain"`
		Proposer    sdk.AccAddress `json:"proposer" yaml:"proposer"`
		Deposit     sdk.Coins      `json:"deposit" yaml:"deposit"`
	}

	// AddDefaultSubscriptionProposalReq defines an add default subscription proposal request body.
	AddDefaultSubscriptionProposalReq struct {
		BaseReq rest.BaseReq `json:"base_req" yaml:"base_req"`

		Title           string         `json:"title" yaml:"title"`
		Description     string         `json:"description" yaml:"description"`
		SubscriptionID  string         `json:"subscription_id" yaml:"subscription_id"`
		PublisherDomain string         `json:"publisher_domain" yaml:"publisher_domain"`
		Proposer        sdk.AccAddress `json:"proposer" yaml:"proposer"`
		Deposit         sdk.Coins      `json:"deposit" yaml:"deposit"`
	}

	// RemoveDefaultSubscriptionProposalReq defines a remove default subscription proposal request body.
	RemoveDefaultSubscriptionProposalReq struct {
		BaseReq rest.BaseReq `json:"base_req" yaml:"base_req"`

		Title          string         `json:"title" yaml:"title"`
		Description    string         `json:"description" yaml:"description"`
		SubscriptionID string         `json:"subscription_id" yaml:"subscription_id"`
		Proposer       sdk.AccAddress `json:"proposer" yaml:"proposer"`
		Deposit        sdk.Coins      `json:"deposit" yaml:"deposit"`
	}
)
