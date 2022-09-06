package client

import (
	govclient "github.com/cosmos/cosmos-sdk/x/gov/client"
	"github.com/peggyjv/sommelier/v4/x/auction/client/cli"
	"github.com/peggyjv/sommelier/v4/x/auction/client/rest"
)

var (
	SetProposalHandler = govclient.NewProposalHandler(cli.GetCmdSubmitSetTokenPricesProposal, rest.SetProposalRESTHandler)
)
