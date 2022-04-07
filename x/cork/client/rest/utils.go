package rest

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/rest"
)

type (
	// AddManagedCellarIDsProposalReq defines a managed cellar ID addition proposal request body.
	AddManagedCellarIDsProposalReq struct {
		BaseReq rest.BaseReq `json:"base_req" yaml:"base_req"`

		Title       string         `json:"title" yaml:"title"`
		Description string         `json:"description" yaml:"description"`
		CellarIDs   []string       `json:"cellar_ids" yaml:"cellar_ids"`
		Proposer    sdk.AccAddress `json:"proposer" yaml:"proposer"`
		Deposit     sdk.Coins      `json:"deposit" yaml:"deposit"`
	}
	// RemoveManagedCellarIDsProposalReq defines a managed cellar ID removal proposal request body.
	RemoveManagedCellarIDsProposalReq struct {
		BaseReq rest.BaseReq `json:"base_req" yaml:"base_req"`

		Title       string         `json:"title" yaml:"title"`
		Description string         `json:"description" yaml:"description"`
		CellarIDs   []string       `json:"cellar_ids" yaml:"cellar_ids"`
		Proposer    sdk.AccAddress `json:"proposer" yaml:"proposer"`
		Deposit     sdk.Coins      `json:"deposit" yaml:"deposit"`
	}
)
