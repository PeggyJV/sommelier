package client

import (
	govclient "github.com/cosmos/cosmos-sdk/x/gov/client"
	"github.com/peggyjv/sommelier/v7/x/axelarcork/client/cli"
)

var (
	AddProposalHandler                        = govclient.NewProposalHandler(cli.GetCmdSubmitAddCellarIDProposal)
	RemoveProposalHandler                     = govclient.NewProposalHandler(cli.GetCmdSubmitRemoveCellarIDProposal)
	ScheduledCorkProposalHandler              = govclient.NewProposalHandler(cli.GetCmdSubmitScheduledAxelarCorkProposal)
	CommunityPoolEthereumSpendProposalHandler = govclient.NewProposalHandler(cli.CmdSubmitAxelarCommunityPoolEthereumSpendProposal)
	AddChainConfigurationHandler              = govclient.NewProposalHandler(cli.GetCmdSubmitAddChainConfigurationProposal)
	RemoveChainConfigurationHandler           = govclient.NewProposalHandler(cli.GetCmdSubmitRemoveChainConfigurationProposal)
	UpgradeAxelarProxyContractHandler         = govclient.NewProposalHandler(cli.GetCmdSubmitUpgradeAxelarProxyContractProposal)
	CancelAxelarProxyContractUpgradeHandler   = govclient.NewProposalHandler(cli.GetCmdSubmitCancelAxelarProxyContractUpgradeProposal)
)
