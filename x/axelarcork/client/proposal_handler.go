package client

import (
	govclient "github.com/cosmos/cosmos-sdk/x/gov/client"
	"github.com/peggyjv/sommelier/v6/x/axelarcork/client/cli"
	"github.com/peggyjv/sommelier/v6/x/axelarcork/client/rest"
)

var (
	AddProposalHandler                        = govclient.NewProposalHandler(cli.GetCmdSubmitAddCellarIDProposal, rest.AddProposalRESTHandler)
	RemoveProposalHandler                     = govclient.NewProposalHandler(cli.GetCmdSubmitRemoveCellarIDProposal, rest.RemoveProposalRESTHandler)
	ScheduledCorkProposalHandler              = govclient.NewProposalHandler(cli.GetCmdSubmitScheduledAxelarCorkProposal, rest.ScheduledCorkProposalRESTHandler)
	CommunityPoolEthereumSpendProposalHandler = govclient.NewProposalHandler(cli.CmdSubmitAxelarCommunityPoolEthereumSpendProposal, rest.CommunitySpendProposalRESTHandler)
	AddChainConfigurationHandler              = govclient.NewProposalHandler(cli.GetCmdSubmitAddChainConfigurationProposal, rest.AddChainConfigurationProposalRESTHandler)
	RemoveChainConfigurationHandler           = govclient.NewProposalHandler(cli.GetCmdSubmitRemoveChainConfigurationProposal, rest.RemoveChainConfigurationProposalRESTHandler)
)
