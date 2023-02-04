package cli

import (
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/version"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	"github.com/peggyjv/sommelier/v4/x/pubsub/types"
)

// GetTxCmd returns the transaction commands for this module
func GetTxCmd() *cobra.Command {
	pubsubTxCommand := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      fmt.Sprintf("%s transaction subcommands", types.ModuleName),
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	return pubsubTxCommand
}

func GetCmdSubmitAddPublisherProposal() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "add-publisher [proposal-file]",
		Args:  cobra.ExactArgs(1),
		Short: "Submit an add publisher proposal",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Submit an add publisher proposal along with an initial deposit.
The proposal details must be supplied via a JSON file.

Example:
$ %s tx gov submit-proposal add-publisher <path/to/proposal.json> --from=<key_or_address>

Where proposal.json contains:

{
  "title": "Add us as a publisher",
  "description": "We do strategies",
  "domain": "sommelier.example.com",
  "address": "<bech 32 somm address owned by the publisher>",
  "proof_url": "https://sommelier.example.com/<the same bech 32 somm address>/cacert.pem",
  "ca_cert": "<text contents of the self-signed cacert.pem file served above>",
  "deposit": "10000000usomm"
}
`,
				version.AppName,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			proposal := types.AddPublisherProposalWithDeposit{}
			contents, err := ioutil.ReadFile(args[0])
			if err != nil {
				return err
			}

			if err = clientCtx.Codec.UnmarshalJSON(contents, &proposal); err != nil {
				return err
			}

			deposit, err := sdk.ParseCoinsNormalized(proposal.Deposit)
			if err != nil {
				return err
			}

			content := types.NewAddPublisherProposal(
				proposal.Title,
				proposal.Description,
				proposal.Domain,
				proposal.Address,
				proposal.ProofUrl,
				proposal.CaCert,
			)

			if err = content.ValidateBasic(); err != nil {
				return err
			}

			from := clientCtx.GetFromAddress()
			msg, err := govtypes.NewMsgSubmitProposal(content, deposit, from)
			if err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	return cmd
}

func GetCmdSubmitRemovePublisherProposal() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "remove-publisher [proposal-file]",
		Args:  cobra.ExactArgs(1),
		Short: "Submit a remove publisher proposal",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Submit a remove publisher proposal along with an initial deposit.
The proposal details must be supplied via a JSON file.

Example:
$ %s tx gov submit-proposal remove-publisher <path/to/proposal.json> --from=<key_or_address>

Where proposal.json contains:

{
  "title": "Remove this publisher",
  "description": "Get rid of them",
  "domain": "sommelier.example.com",
  "deposit": "10000000usomm"
}
`,
				version.AppName,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			proposal := types.RemovePublisherProposalWithDeposit{}
			contents, err := ioutil.ReadFile(args[0])
			if err != nil {
				return err
			}

			if err = clientCtx.Codec.UnmarshalJSON(contents, &proposal); err != nil {
				return err
			}

			deposit, err := sdk.ParseCoinsNormalized(proposal.Deposit)
			if err != nil {
				return err
			}

			content := types.NewRemovePublisherProposal(
				proposal.Title,
				proposal.Description,
				proposal.Domain,
			)

			if err = content.ValidateBasic(); err != nil {
				return err
			}

			from := clientCtx.GetFromAddress()
			msg, err := govtypes.NewMsgSubmitProposal(content, deposit, from)
			if err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	return cmd
}

func GetCmdSubmitAddDefaultSubscriptionProposal() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "add-default-subscription [proposal-file]",
		Args:  cobra.ExactArgs(1),
		Short: "Submit an add default subscription proposal",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Submit an add default subscription proposal along with an initial deposit.
The proposal details must be supplied via a JSON file.

Example:
$ %s tx gov submit-proposal add-default-subscription <path/to/proposal.json> --from=<key_or_address>

Where proposal.json contains:

{
  "title": "Add this default subscription",
  "description": "Use this publisher unless you override it",
  "subscription_id": "0x123801a7D398351b8bE11C439e05C5B3259aeC9B",
  "publisher_domain": "sommelier.example.com",
  "deposit": "10000000usomm"
}
`,
				version.AppName,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			proposal := types.AddDefaultSubscriptionProposalWithDeposit{}
			contents, err := ioutil.ReadFile(args[0])
			if err != nil {
				return err
			}

			if err = clientCtx.Codec.UnmarshalJSON(contents, &proposal); err != nil {
				return err
			}

			deposit, err := sdk.ParseCoinsNormalized(proposal.Deposit)
			if err != nil {
				return err
			}

			content := types.NewAddDefaultSubscriptionProposal(
				proposal.Title,
				proposal.Description,
				proposal.SubscriptionId,
				proposal.PublisherDomain,
			)

			if err = content.ValidateBasic(); err != nil {
				return err
			}

			from := clientCtx.GetFromAddress()
			msg, err := govtypes.NewMsgSubmitProposal(content, deposit, from)
			if err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	return cmd
}

func GetCmdSubmitRemoveDefaultSubscriptionProposal() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "remove-default-subscription [proposal-file]",
		Args:  cobra.ExactArgs(1),
		Short: "Submit a remove default subscription proposal",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Submit a remove default subscription proposal along with an initial deposit.
The proposal details must be supplied via a JSON file.

Example:
$ %s tx gov submit-proposal add-default-subscription <path/to/proposal.json> --from=<key_or_address>

Where proposal.json contains:

{
  "title": "Remove this default subscription",
  "description": "Remove the defeault for this subscription ID",
  "subscription_id": "0x123801a7D398351b8bE11C439e05C5B3259aeC9B",
  "deposit": "10000000usomm"
}
`,
				version.AppName,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			proposal := types.RemoveDefaultSubscriptionProposalWithDeposit{}
			contents, err := ioutil.ReadFile(args[0])
			if err != nil {
				return err
			}

			if err = clientCtx.Codec.UnmarshalJSON(contents, &proposal); err != nil {
				return err
			}

			deposit, err := sdk.ParseCoinsNormalized(proposal.Deposit)
			if err != nil {
				return err
			}

			content := types.NewAddDefaultSubscriptionProposal(
				proposal.Title,
				proposal.Description,
				proposal.SubscriptionId,
			)

			if err = content.ValidateBasic(); err != nil {
				return err
			}

			from := clientCtx.GetFromAddress()
			msg, err := govtypes.NewMsgSubmitProposal(content, deposit, from)
			if err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	return cmd
}
