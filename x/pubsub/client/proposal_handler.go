package client

import (
	govclient "github.com/cosmos/cosmos-sdk/x/gov/client"
	"github.com/peggyjv/sommelier/v7/x/pubsub/client/cli"
	"github.com/peggyjv/sommelier/v7/x/pubsub/client/rest"
)

var (
	AddPublisherProposalHandler              = govclient.NewProposalHandler(cli.GetCmdSubmitAddPublisherProposal, rest.AddPublisherProposalRESTHandler)
	RemovePublisherProposalHandler           = govclient.NewProposalHandler(cli.GetCmdSubmitRemovePublisherProposal, rest.RemovePublisherProposalRESTHandler)
	AddDefaultSubscriptionProposalHandler    = govclient.NewProposalHandler(cli.GetCmdSubmitAddDefaultSubscriptionProposal, rest.AddDefaultSubscriptionProposalRESTHandler)
	RemoveDefaultSubscriptionProposalHandler = govclient.NewProposalHandler(cli.GetCmdSubmitRemoveDefaultSubscriptionProposal, rest.RemoveDefaultSubscriptionProposalRESTHandler)
)
