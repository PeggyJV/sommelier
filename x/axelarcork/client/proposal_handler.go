package client

import (
	govclient "github.com/cosmos/cosmos-sdk/x/gov/client"
	"github.com/peggyjv/sommelier/v6/x/axelarcork/client/cli"
	"github.com/peggyjv/sommelier/v6/x/axelarcork/client/rest"
)

var (
	AddProposalHandler                        = govclient.NewProposalHandler(cli.GetCmdSubmitAddProposal, rest.AddProposalRESTHandler)
	RemoveProposalHandler                     = govclient.NewProposalHandler(cli.GetCmdSubmitRemoveProposal, rest.RemoveProposalRESTHandler)
	ScheduledCorkProposalHandler              = govclient.NewProposalHandler(cli.GetCmdSubmitScheduledCorkProposal, rest.ScheduledCorkProposalRESTHandler)
	CommunityPoolEthereumSpendProposalHandler = govclient.NewProposalHandler(cli.CmdSubmitCommunityPoolEthereumSpendProposal, rest.CommunitySpendProposalRESTHandler)
	AddChainConfigurationHandler              = govclient.NewProposalHandler(cli.GetCmdSubmitAddChainConfigurationProposal, rest.AddChainConfigurationProposalRESTHandler)
	RemoveChainConfigurationHandler           = govclient.NewProposalHandler(cli.GetCmdSubmitRemoveChainConfigurationProposal, rest.RemoveChainConfigurationProposalRESTHandler)
)
