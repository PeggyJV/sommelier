package client

import (
	govclient "github.com/cosmos/cosmos-sdk/x/gov/client"
	"github.com/peggyjv/sommelier/v3/x/cork/client/cli"
)

var (
	AddProposalHandler    = govclient.NewProposalHandler(cli.GetCmdSubmitAddProposal, nil)
	RemoveProposalHandler = govclient.NewProposalHandler(cli.GetCmdSubmitRemoveProposal, nil)
)
