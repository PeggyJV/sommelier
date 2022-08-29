package rest

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/rest"
	"github.com/peggyjv/sommelier/v4/x/auction/types"
)

type (
	// SetTokenPricesProposalReq defines a token price set proposal request body.
	SetTokenPricesProposalReq struct {
		BaseReq rest.BaseReq `json:"base_req" yaml:"base_req"`

		Title       string                      `json:"title" yaml:"title"`
		Description string                      `json:"description" yaml:"description"`
		TokenPrices []*types.ProposedTokenPrice `json:"token_prices" yaml:"token_prices"`
		Proposer    sdk.AccAddress              `json:"proposer" yaml:"proposer"`
		Deposit     sdk.Coins                   `json:"deposit" yaml:"deposit"`
	}
)
