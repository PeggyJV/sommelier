package rest

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/rest"
	"github.com/peggyjv/sommelier/v4/x/auction/types"
)

type (
	// UpdateTokenPricesProposalReq defines a token price update proposal request body.
	UpdateTokenPricesProposalReq struct {
		BaseReq rest.BaseReq 				`json:"base_req" yaml:"base_req"`

		Title         string         		`json:"title" yaml:"title"`
		Description   string         		`json:"description" yaml:"description"`
		TokenPrices   []*types.TokenPrice   `json:"token_prices" yaml:"token_prices"`
		Proposer      sdk.AccAddress 		`json:"proposer" yaml:"proposer"`
		Deposit       sdk.Coins     		`json:"deposit" yaml:"deposit"`
	}
)
