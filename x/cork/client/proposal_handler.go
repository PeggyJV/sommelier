package client

import (
	govclient "github.com/cosmos/cosmos-sdk/x/gov/client"
	"github.com/peggyjv/sommelier/v6/x/cork/client/cli"
	"github.com/peggyjv/sommelier/v6/x/cork/client/rest"
)

var (
	AddProposalHandler    = govclient.NewProposalHandler(cli.GetCmdSubmitAddProposal, rest.AddProposalRESTHandler)
	RemoveProposalHandler = govclient.NewProposalHandler(cli.GetCmdSubmitRemoveProposal, rest.RemoveProposalRESTHandler)
)
