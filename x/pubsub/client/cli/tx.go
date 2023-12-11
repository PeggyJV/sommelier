package cli

import (
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/version"
	govtypesv1beta1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1beta1"
	"github.com/peggyjv/sommelier/v7/x/pubsub/types"
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

	pubsubTxCommand.AddCommand([]*cobra.Command{
		GetCmdAddPublisherPullIntent(),
		GetCmdAddPublisherPushIntent(),
		GetCmdAddSubscriberIntent(),
		GetCmdAddSubscriber(),
		GetCmdRemovePublisherIntent(),
		GetCmdRemoveSubscriberIntent(),
		GetCmdRemoveSubscriber(),
		GetCmdRemovePublisher(),
	}...)

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
  "address": "somm1y6d5kasehecexf09ka6y0ggl0pxzt6dg6n8lw0",
  "proof_url": "https://sommelier.example.com/somm1y6d5kasehecexf09ka6y0ggl0pxzt6dg6n8lw0/cacert.pem",
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
			msg, err := govtypesv1beta1.NewMsgSubmitProposal(content, deposit, from)
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
			msg, err := govtypesv1beta1.NewMsgSubmitProposal(content, deposit, from)
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
			msg, err := govtypesv1beta1.NewMsgSubmitProposal(content, deposit, from)
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

			content := types.NewRemoveDefaultSubscriptionProposal(
				proposal.Title,
				proposal.Description,
				proposal.SubscriptionId,
			)

			if err = content.ValidateBasic(); err != nil {
				return err
			}

			from := clientCtx.GetFromAddress()
			msg, err := govtypesv1beta1.NewMsgSubmitProposal(content, deposit, from)
			if err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	return cmd
}

func GetCmdAddPublisherPullIntent() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "add-publisher-pull-intent [subscription-id] [publisher-domain] [pull-url] [ANY|VALIDATORS|LIST] <optional_allowed_addresses>",
		Args:  cobra.MinimumNArgs(4),
		Short: "Add a publisher intent with a pull URL",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Examples:
$ %s tx pubsub add-publisher-pull-intent 0x123801a7D398351b8bE11C439e05C5B3259aeC9B sommelier.example.com "https://sommelier.example.com/pull" VALIDATORS --from=<key_or_address>
$ %s tx pubsub add-publisher-pull-intent 0x123801a7D398351b8bE11C439e05C5B3259aeC9B sommelier.example.com "https://sommelier.example.com/pull" LIST somm1y6d5kasehecexf09ka6y0ggl0pxzt6dg6n8lw0,somm18ld4633yswcyjdklej3att6aw93nhlf7596qkk --from=<key_or_address>
`, version.AppName, version.AppName),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			if len(args) > 5 {
				return fmt.Errorf("too many arguments")
			}

			signer := clientCtx.GetFromAddress()
			if signer == nil {
				return fmt.Errorf("must include `--from` flag")
			}

			allowedSubscribers, err := parseAllowedSubscribers(args[3])
			if err != nil {
				return err
			}

			allowedAddresses := []string{}
			if len(args) == 5 {
				allowedAddresses = strings.Split(args[4], ",")
			}

			publisherIntent := types.PublisherIntent{
				SubscriptionId:     args[0],
				PublisherDomain:    args[1],
				Method:             types.PublishMethod_PULL,
				PullUrl:            args[2],
				AllowedSubscribers: allowedSubscribers,
				AllowedAddresses:   allowedAddresses,
			}

			msg, err := types.NewMsgAddPublisherIntentRequest(publisherIntent, signer)
			if err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)
	return cmd
}

func GetCmdAddPublisherPushIntent() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "add-publisher-push-intent [subscription-id] [publisher-domain] [ANY|VALIDATORS|LIST] <optional_allowed_addresses>",
		Args:  cobra.MinimumNArgs(3),
		Short: "Add a publisher intent that will push",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Examples:
$ %s tx pubsub add-publisher-push-intent 0x123801a7D398351b8bE11C439e05C5B3259aeC9B sommelier.example.com VALIDATORS --from=<key_or_address>
$ %s tx pubsub add-publisher-push-intent 0x123801a7D398351b8bE11C439e05C5B3259aeC9B sommelier.example.com LIST somm1y6d5kasehecexf09ka6y0ggl0pxzt6dg6n8lw0,somm18ld4633yswcyjdklej3att6aw93nhlf7596qkk --from=<key_or_address>
`, version.AppName, version.AppName),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			if len(args) > 4 {
				return fmt.Errorf("too many arguments")
			}

			signer := clientCtx.GetFromAddress()
			if signer == nil {
				return fmt.Errorf("must include `--from` flag")
			}

			allowedSubscribers, err := parseAllowedSubscribers(args[2])
			if err != nil {
				return err
			}

			allowedAddresses := []string{}
			if len(args) == 4 {
				allowedAddresses = strings.Split(args[3], ",")
			}

			publisherIntent := types.PublisherIntent{
				SubscriptionId:     args[0],
				PublisherDomain:    args[1],
				Method:             types.PublishMethod_PUSH,
				PullUrl:            "",
				AllowedSubscribers: allowedSubscribers,
				AllowedAddresses:   allowedAddresses,
			}

			msg, err := types.NewMsgAddPublisherIntentRequest(publisherIntent, signer)
			if err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)
	return cmd
}

func GetCmdAddSubscriberIntent() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "add-subscriber-intent [subscription-id] [subscriber-address] [publisher-domain]",
		Args:  cobra.ExactArgs(3),
		Short: "Add a subscriber intent",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Examples:
$ %s tx pubsub add-subscriber-intent 0x123801a7D398351b8bE11C439e05C5B3259aeC9B somm1y6d5kasehecexf09ka6y0ggl0pxzt6dg6n8lw0 pullpublisher.example.com --from=<key_or_address>
`, version.AppName),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			signer := clientCtx.GetFromAddress()
			if signer == nil {
				return fmt.Errorf("must include `--from` flag")
			}

			subscriberIntent := types.SubscriberIntent{
				SubscriptionId:    args[0],
				SubscriberAddress: args[1],
				PublisherDomain:   args[2],
			}

			msg, err := types.NewMsgAddSubscriberIntentRequest(subscriberIntent, signer)
			if err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)
	return cmd
}

func GetCmdAddSubscriber() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "add-subscriber [address] <optional_ca_cert> <optional_push_url>",
		Args:  cobra.MinimumNArgs(1),
		Short: "Add a subscriber",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Examples:
$ %s tx pubsub add-subscriber somm1y6d5kasehecexf09ka6y0ggl0pxzt6dg6n8lw0 --from=<key_or_address>
$ %s tx pubsub add-subscriber somm1y6d5kasehecexf09ka6y0ggl0pxzt6dg6n8lw0 <path/to/cacert.pem> "sommvalidator.example.com:5734" --from=<key_or_address>
`, version.AppName, version.AppName),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			if len(args) > 3 {
				return fmt.Errorf("too many arguments")
			}

			signer := clientCtx.GetFromAddress()
			if signer == nil {
				return fmt.Errorf("must include `--from` flag")
			}

			if len(args) > 1 && len(args) != 3 {
				return fmt.Errorf("must include both CA cert path and push URL for push subscriptions")
			}

			caCert := ""
			pushURL := ""

			if len(args) == 3 {
				caCertContent, err := ioutil.ReadFile(args[1])
				if err != nil {
					return fmt.Errorf("cannot read CA cert: %s", err)
				}
				caCert = string(caCertContent)
				pushURL = args[2]
			}

			subscriber := types.Subscriber{
				Address: args[0],
				CaCert:  caCert,
				PushUrl: pushURL,
			}

			msg, err := types.NewMsgAddSubscriberRequest(subscriber, signer)
			if err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)
	return cmd
}

func GetCmdRemovePublisherIntent() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "remove-publisher-intent [subscription-id] [publisher-domain]",
		Args:  cobra.ExactArgs(2),
		Short: "Remove a publisher intent",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Examples:
$ %s tx pubsub remove-publisher-intent 0x123801a7D398351b8bE11C439e05C5B3259aeC9B pushpublisher.example.com --from=<key_or_address>
`, version.AppName),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			signer := clientCtx.GetFromAddress()
			if signer == nil {
				return fmt.Errorf("must include `--from` flag")
			}

			msg, err := types.NewMsgRemovePublisherIntentRequest(args[0], args[1], signer)
			if err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)
	return cmd
}

func GetCmdRemoveSubscriberIntent() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "remove-subscriber-intent [subscription-id] [subscriber-address]",
		Args:  cobra.ExactArgs(2),
		Short: "Remove a subcriber intent",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Examples:
$ %s tx pubsub remove-subscriber-intent 0x123801a7D398351b8bE11C439e05C5B3259aeC9B somm1y6d5kasehecexf09ka6y0ggl0pxzt6dg6n8lw0 --from=<key_or_address>
`, version.AppName),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			signer := clientCtx.GetFromAddress()
			if signer == nil {
				return fmt.Errorf("must include `--from` flag")
			}

			msg, err := types.NewMsgRemoveSubscriberIntentRequest(args[0], args[1], signer)
			if err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)
	return cmd
}

func GetCmdRemoveSubscriber() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "remove-subscriber [subscriber-address]",
		Args:  cobra.ExactArgs(1),
		Short: "Remove a subcriber",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Examples:
$ %s tx pubsub remove-subscriber somm1y6d5kasehecexf09ka6y0ggl0pxzt6dg6n8lw0 --from=<key_or_address>
`, version.AppName),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			signer := clientCtx.GetFromAddress()
			if signer == nil {
				return fmt.Errorf("must include `--from` flag")
			}

			msg, err := types.NewMsgRemoveSubscriberRequest(args[0], signer)
			if err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)
	return cmd
}

func GetCmdRemovePublisher() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "remove-publisher [publisher-domain]",
		Args:  cobra.ExactArgs(1),
		Short: "Remove a publisher",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Examples:
$ %s tx pubsub remove-publisher publisher.example.com --from=<key_or_address>
`, version.AppName),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			signer := clientCtx.GetFromAddress()
			if signer == nil {
				return fmt.Errorf("must include `--from` flag")
			}

			msg, err := types.NewMsgRemovePublisherRequest(args[0], signer)
			if err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)
	return cmd
}

func parseAllowedSubscribers(allowedSubscribers string) (types.AllowedSubscribers, error) {
	switch allowedSubscribers {
	case types.AllowedSubscribers_ANY.String():
		return types.AllowedSubscribers_ANY, nil
	case types.AllowedSubscribers_VALIDATORS.String():
		return types.AllowedSubscribers_VALIDATORS, nil
	case types.AllowedSubscribers_LIST.String():
		return types.AllowedSubscribers_LIST, nil
	default:
		return types.AllowedSubscribers_ANY, fmt.Errorf("invalid allowed subcribers selection, choose ANY, VALIDATORS, or LIST: %s", allowedSubscribers)
	}
}
