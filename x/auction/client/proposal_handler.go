package client

import (
	govclient "github.com/cosmos/cosmos-sdk/x/gov/client"
	"github.com/peggyjv/sommelier/v4/x/auction/client/cli"
	"github.com/peggyjv/sommelier/v4/x/auction/client/rest"
)

var (
	AddProposalHandler    = govclient.NewProposalHandler(cli.GetCmdSubmitAddProposal, rest.AddProposalRESTHandler)
)
