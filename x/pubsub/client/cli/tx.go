package cli

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/peggyjv/sommelier/v4/x/pubsub/types"
)

// TODO(bolten): fill out tx commands

// GetTxCmd returns the transaction commands for this module
func GetTxCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      fmt.Sprintf("%s transactions subcommands", types.ModuleName),
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	return cmd
}

func GetCmdSubmitAddPublisherProposal() *cobra.Command {
	return nil
}

func GetCmdSubmitRemovePublisherProposal() *cobra.Command {
	return nil
}

func GetCmdSubmitAddDefaultSubscriptionProposal() *cobra.Command {
	return nil
}

func GetCmdSubmitRemoveDefaultSubscriptionProposal() *cobra.Command {
	return nil
}
