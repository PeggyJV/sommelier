package client

import (
	govclient "github.com/cosmos/cosmos-sdk/x/gov/client"
	"github.com/peggyjv/sommelier/v7/x/auction/client/cli"
)

var (
	SetProposalHandler = govclient.NewProposalHandler(cli.GetCmdSubmitSetTokenPricesProposal)
)
