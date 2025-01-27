package client

import (
	govclient "github.com/cosmos/cosmos-sdk/x/gov/client"
	"github.com/peggyjv/sommelier/v9/x/axelarcork/client/cli"
)

var (
	AddProposalHandler                        = govclient.NewProposalHandler(cli.GetCmdSubmitAddAxelarCellarIDProposal)
	RemoveProposalHandler                     = govclient.NewProposalHandler(cli.GetCmdSubmitRemoveAxelarCellarIDProposal)
	ScheduledCorkProposalHandler              = govclient.NewProposalHandler(cli.GetCmdSubmitScheduledAxelarCorkProposal)
	CommunityPoolEthereumSpendProposalHandler = govclient.NewProposalHandler(cli.CmdSubmitAxelarCommunityPoolSpendProposal)
	AddChainConfigurationHandler              = govclient.NewProposalHandler(cli.GetCmdSubmitAddChainConfigurationProposal)
	RemoveChainConfigurationHandler           = govclient.NewProposalHandler(cli.GetCmdSubmitRemoveChainConfigurationProposal)
	UpgradeAxelarProxyContractHandler         = govclient.NewProposalHandler(cli.GetCmdSubmitUpgradeAxelarProxyContractProposal)
	CancelAxelarProxyContractUpgradeHandler   = govclient.NewProposalHandler(cli.GetCmdSubmitCancelAxelarProxyContractUpgradeProposal)
)
