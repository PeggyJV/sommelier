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
	// ScheduledCorkProposalReq defines a schedule cork proposal request body.
	ScheduledCorkProposalReq struct {
		BaseReq rest.BaseReq `json:"base_req" yaml:"base_req"`

		Title                 string         `json:"title" yaml:"title"`
		Description           string         `json:"description" yaml:"description"`
		BlockHeight           uint64         `json:"block_height" yaml:"block_height"`
		TargetContractAddress string         `json:"target_contract_address" yaml:"target_contract_address"`
		ContractCallProtoJSON string         `json:"contract_call_proto_json" yaml:"contract_call_proto_json"`
		Proposer              sdk.AccAddress `json:"proposer" yaml:"proposer"`
		Deposit               sdk.Coins      `json:"deposit" yaml:"deposit"`
	}
)
