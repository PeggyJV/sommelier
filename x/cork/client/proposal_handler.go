package client

import (
	govclient "github.com/cosmos/cosmos-sdk/x/gov/client"
	"github.com/peggyjv/sommelier/v7/x/cork/client/cli"
)

var (
	AddProposalHandler           = govclient.NewProposalHandler(cli.GetCmdSubmitAddProposal)
	RemoveProposalHandler        = govclient.NewProposalHandler(cli.GetCmdSubmitRemoveProposal)
	ScheduledCorkProposalHandler = govclient.NewProposalHandler(cli.GetCmdSubmitScheduledCorkProposal)
)
