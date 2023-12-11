package rest

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/rest"
	"github.com/peggyjv/sommelier/v7/x/axelarcork/types"
)

type (
	// AddManagedCellarIDsProposalReq defines a managed cellar ID addition proposal request body.
	AddManagedCellarIDsProposalReq struct {
		BaseReq rest.BaseReq `json:"base_req" yaml:"base_req"`

		Title       string         `json:"title" yaml:"title"`
		Description string         `json:"description" yaml:"description"`
		ChainID     uint64         `json:"chain_id" yaml:"chain_id"`
		CellarIDs   []string       `json:"cellar_ids" yaml:"cellar_ids"`
		Proposer    sdk.AccAddress `json:"proposer" yaml:"proposer"`
		Deposit     sdk.Coins      `json:"deposit" yaml:"deposit"`
	}
	// RemoveManagedCellarIDsProposalReq defines a managed cellar ID removal proposal request body.
	RemoveManagedCellarIDsProposalReq struct {
		BaseReq rest.BaseReq `json:"base_req" yaml:"base_req"`

		Title       string         `json:"title" yaml:"title"`
		Description string         `json:"description" yaml:"description"`
		ChainID     uint64         `json:"chain_id" yaml:"chain_id"`
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
		ChainName             string         `json:"chain_name" yaml:"chain_name"`
		ChainID               uint64         `json:"chain_id" yaml:"chain_id"`
		TargetContractAddress string         `json:"target_contract_address" yaml:"target_contract_address"`
		ContractCallProtoJSON string         `json:"contract_call_proto_json" yaml:"contract_call_proto_json"`
		Proposer              sdk.AccAddress `json:"proposer" yaml:"proposer"`
		Deposit               sdk.Coins      `json:"deposit" yaml:"deposit"`
	}
	// CommunityPoolSpendProposalReq defines a community pool spend proposal request body.
	CommunityPoolSpendProposalReq struct {
		BaseReq rest.BaseReq `json:"base_req" yaml:"base_req"`

		Title       string         `json:"title" yaml:"title"`
		Description string         `json:"description" yaml:"description"`
		Recipient   string         `json:"recipient" yaml:"recipient"`
		Amount      sdk.Coin       `json:"amount" yaml:"amount"`
		ChainID     uint64         `json:"chain_id" yaml:"chain_id"`
		ChainName   string         `json:"chain_name" yaml:"chain_name"`
		Proposer    sdk.AccAddress `json:"proposer" yaml:"proposer"`
		Deposit     sdk.Coins      `json:"deposit" yaml:"deposit"`
	}
	// AddChainConfigurationProposalReq defines a chain configuration addition proposal request body.
	AddChainConfigurationProposalReq struct {
		BaseReq rest.BaseReq `json:"base_req" yaml:"base_req"`

		Title              string                   `json:"title" yaml:"title"`
		Description        string                   `json:"description" yaml:"description"`
		ChainConfiguration types.ChainConfiguration `json:"chain_configuration" yaml:"chain_configuration"`
		Proposer           sdk.AccAddress           `json:"proposer" yaml:"proposer"`
		Deposit            sdk.Coins                `json:"deposit" yaml:"deposit"`
	}
	// RemoveChainConfigurationProposalReq defines a chain configuration removal proposal request body.
	RemoveChainConfigurationProposalReq struct {
		BaseReq rest.BaseReq `json:"base_req" yaml:"base_req"`

		Title       string         `json:"title" yaml:"title"`
		Description string         `json:"description" yaml:"description"`
		ChainID     uint64         `json:"chain_id" yaml:"chain_id"`
		Proposer    sdk.AccAddress `json:"proposer" yaml:"proposer"`
		Deposit     sdk.Coins      `json:"deposit" yaml:"deposit"`
	}
	// UpgradeAxelarProxyContractProposalReq defines a upgrade axelar proxy contract proposal request body.
	UpgradeAxelarProxyContractProposalReq struct {
		BaseReq rest.BaseReq `json:"base_req" yaml:"base_req"`

		Title           string         `json:"title" yaml:"title"`
		Description     string         `json:"description" yaml:"description"`
		ChainID         uint64         `json:"chain_id" yaml:"chain_id"`
		NewProxyAddress string         `json:"new_proxy_address" yaml:"new_proxy_address"`
		Proposer        sdk.AccAddress `json:"proposer" yaml:"proposer"`
		Deposit         sdk.Coins      `json:"deposit" yaml:"deposit"`
	}

	// CancelAxelarProxyContractUpgradeProposalReq defines a cancel axelar proxy contract upgrade proposal request body.
	CancelAxelarProxyContractUpgradeProposalReq struct {
		BaseReq rest.BaseReq `json:"base_req" yaml:"base_req"`

		Title       string         `json:"title" yaml:"title"`
		Description string         `json:"description" yaml:"description"`
		ChainID     uint64         `json:"chain_id" yaml:"chain_id"`
		Proposer    sdk.AccAddress `json:"proposer" yaml:"proposer"`
		Deposit     sdk.Coins      `json:"deposit" yaml:"deposit"`
	}
)
