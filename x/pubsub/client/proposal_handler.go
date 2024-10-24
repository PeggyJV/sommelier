package client

import (
	govclient "github.com/cosmos/cosmos-sdk/x/gov/client"
	"github.com/peggyjv/sommelier/v8/x/pubsub/client/cli"
)

var (
	AddPublisherProposalHandler              = govclient.NewProposalHandler(cli.GetCmdSubmitAddPublisherProposal)
	RemovePublisherProposalHandler           = govclient.NewProposalHandler(cli.GetCmdSubmitRemovePublisherProposal)
	AddDefaultSubscriptionProposalHandler    = govclient.NewProposalHandler(cli.GetCmdSubmitAddDefaultSubscriptionProposal)
	RemoveDefaultSubscriptionProposalHandler = govclient.NewProposalHandler(cli.GetCmdSubmitRemoveDefaultSubscriptionProposal)
)
