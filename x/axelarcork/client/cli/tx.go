package cli

import (
	"fmt"
	"os"
	"strings"

	"cosmossdk.io/math"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/version"
	govtypesv1beta1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1beta1"
	"github.com/ethereum/go-ethereum/common"
	types "github.com/peggyjv/sommelier/v7/x/axelarcork/types"
	pubsubtypes "github.com/peggyjv/sommelier/v7/x/pubsub/types"
	"github.com/spf13/cobra"
)

// GetTxCmd returns the transaction commands for this module
func GetTxCmd() *cobra.Command {
	corkTxCmd := &cobra.Command{
		Use:                        "axelarcork",
		Short:                      "AxelarCork transaction subcommands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	corkTxCmd.AddCommand(
		CmdScheduleAxelarCork(),
		CmdRelayAxelarCork(),
		CmdBumpAxelarCorkGas(),
		CmdRelayAxelarProxyUpgrade())

	return corkTxCmd
}

//////////////
// Commands //
//////////////

func CmdScheduleAxelarCork() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "schedule-axelar-cork [chain-id] [contract-address] [block-height] [contract-call]",
		Args:  cobra.ExactArgs(4),
		Short: "Schedule an Axelar cork",

		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			from := clientCtx.GetFromAddress()

			chainID, err := math.ParseUint(args[0])
			if err != nil {
				return err
			}

			contractAddr := args[1]
			if !common.IsHexAddress(contractAddr) {
				return fmt.Errorf("contract address %s is invalid", contractAddr)
			}

			blockHeight, err := math.ParseUint(args[2])
			if err != nil {
				return err
			}

			contractCallBz := []byte(args[3]) // todo: how are contract calls submitted?

			scheduleCorkMsg := types.MsgScheduleAxelarCorkRequest{
				Cork: &types.AxelarCork{
					EncodedContractCall:   contractCallBz,
					ChainId:               chainID.Uint64(),
					TargetContractAddress: contractAddr,
				},
				ChainId:     chainID.Uint64(),
				BlockHeight: blockHeight.Uint64(),
				Signer:      from.String(),
			}
			if err := scheduleCorkMsg.ValidateBasic(); err != nil {
				return err
			}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), &scheduleCorkMsg)

		},
	}

	return cmd
}

func CmdRelayAxelarCork() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "relay-axelar-cork [chain-id] [contract-address] [token] [fee]",
		Args:  cobra.ExactArgs(4),
		Short: "Relay a consensus validated Axelar cork",

		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			from := clientCtx.GetFromAddress()

			chainID, err := math.ParseUint(args[0])
			if err != nil {
				return err
			}

			contractAddr := args[1]
			if !common.IsHexAddress(contractAddr) {
				return fmt.Errorf("contract address %s is invalid", contractAddr)
			}

			token, err := sdk.ParseCoinNormalized(args[2])
			if err != nil {
				return err
			}

			fee, err := math.ParseUint(args[3])
			if err != nil {
				return err
			}

			relayCorkMsg := types.MsgRelayAxelarCorkRequest{
				Signer:                from.String(),
				Token:                 token,
				Fee:                   fee.Uint64(),
				ChainId:               chainID.Uint64(),
				TargetContractAddress: contractAddr,
			}
			if err := relayCorkMsg.ValidateBasic(); err != nil {
				return err
			}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), &relayCorkMsg)

		},
	}

	return cmd
}

func CmdRelayAxelarProxyUpgrade() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "relay-axelar-proxy-upgrade [chain-id] [token] [fee]",
		Args:  cobra.ExactArgs(4),
		Short: "Relay a proxy contract upgrade call",

		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			from := clientCtx.GetFromAddress()

			chainID, err := math.ParseUint(args[0])
			if err != nil {
				return err
			}

			token, err := sdk.ParseCoinNormalized(args[2])
			if err != nil {
				return err
			}

			fee, err := math.ParseUint(args[3])
			if err != nil {
				return err
			}

			relayProxyUpgradeMsg := types.MsgRelayAxelarProxyUpgradeRequest{
				Signer:  from.String(),
				Token:   token,
				Fee:     fee.Uint64(),
				ChainId: chainID.Uint64(),
			}
			if err := relayProxyUpgradeMsg.ValidateBasic(); err != nil {
				return err
			}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), &relayProxyUpgradeMsg)

		},
	}

	return cmd
}

func CmdBumpAxelarCorkGas() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "bump-axelar-cork-gas [token] [message-id]",
		Args:  cobra.ExactArgs(2),
		Short: "Add gas for an Axelar cork that is stuck relaying",

		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			from := clientCtx.GetFromAddress()

			token, err := sdk.ParseCoinNormalized(args[0])
			if err != nil {
				return err
			}

			messageID := args[1]

			bumpAxelarCorkGasMsg := types.MsgBumpAxelarCorkGasRequest{
				Signer:    from.String(),
				Token:     token,
				MessageId: messageID,
			}
			if err := bumpAxelarCorkGasMsg.ValidateBasic(); err != nil {
				return err
			}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), &bumpAxelarCorkGasMsg)

		},
	}

	return cmd
}

///////////////
// Proposals //
///////////////

// GetCmdSubmitAddAxelarCellarIDProposal implements the command to submit a cellar id addition proposal
func GetCmdSubmitAddAxelarCellarIDProposal() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "add-axelar-cellar-id [proposal-file]",
		Args:  cobra.ExactArgs(1),
		Short: "Submit an Axelar cellar ID addition proposal",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Submit an Axelar cellar addition proposal along with an initial deposit.
The proposal details must be supplied via a JSON file.

Example:
$ %s tx gov submit-proposal add-axelar-cellar-id <path/to/proposal.json> --from=<key_or_address>

Where proposal.json contains:

{
  "title": "Dollary-doos LP Arbitrum Cellar Proposal",
  "description": "I have a hunch",
  "chain_id": 42161,
  "cellar_ids": ["0x123801a7D398351b8bE11C439e05C5B3259aeC9B", "0x456801a7D398351b8bE11C439e05C5B3259aeC9B"],
  "publisher_domain": "example.com",
  "deposit": "10000usomm"
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

			proposal := types.AddAxelarManagedCellarIDsProposalWithDeposit{}
			contents, err := os.ReadFile(args[0])
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

			if proposal.ChainId == 0 {
				return fmt.Errorf("chain ID must be non-zero")
			}

			for _, id := range proposal.CellarIds {
				if !common.IsHexAddress(id) {
					return fmt.Errorf("%s is not a valid EVM address", id)
				}
			}

			if err := pubsubtypes.ValidateDomain(proposal.PublisherDomain); err != nil {
				return err
			}

			content := types.NewAddAxelarManagedCellarIDsProposal(
				proposal.Title,
				proposal.Description,
				proposal.ChainId,
				&types.CellarIDSet{ChainId: proposal.ChainId, Ids: proposal.CellarIds},
				proposal.PublisherDomain,
			)

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

// GetCmdSubmitRemoveAxelarCellarIDProposal implements the command to submit a cellar id removal proposal
func GetCmdSubmitRemoveAxelarCellarIDProposal() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "remove-axelar-cellar-id [proposal-file]",
		Args:  cobra.ExactArgs(1),
		Short: "Submit an Axelar cellar ID removal proposal",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Submit an Axelar cellar removal proposal along with an initial deposit.
The proposal details must be supplied via a JSON file.

Example:
$ %s tx gov submit-proposal remove-axelar-cellar-id <path/to/proposal.json> --from=<key_or_address>

Where proposal.json contains:

{
  "title": "Dollary-doos LP Arbitrum Cellar Removal Proposal",
  "description": "I don't trust them",
  "chain_id": 42161,
  "cellar_ids": ["0x123801a7D398351b8bE11C439e05C5B3259aeC9B", "0x456801a7D398351b8bE11C439e05C5B3259aeC9B"],
  "deposit": "10000usomm"
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

			proposal := types.RemoveAxelarManagedCellarIDsProposalWithDeposit{}
			contents, err := os.ReadFile(args[0])
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

			if proposal.ChainId == 0 {
				return fmt.Errorf("chain ID must be non-zero")
			}

			for _, id := range proposal.CellarIds {
				if !common.IsHexAddress(id) {
					return fmt.Errorf("%s is not a valid EVM address", id)
				}
			}

			content := types.NewRemoveAxelarManagedCellarIDsProposal(
				proposal.Title,
				proposal.Description,
				proposal.ChainId,
				&types.CellarIDSet{ChainId: proposal.ChainId, Ids: proposal.CellarIds})

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

// GetCmdSubmitScheduledAxelarCorkProposal implements the command to submit scheduled cork proposal
func GetCmdSubmitScheduledAxelarCorkProposal() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "schedule-axelar-cork [proposal-file]",
		Args:  cobra.ExactArgs(1),
		Short: "Submit a scheduled Axelar cork proposal",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Submit an Axelar scheduled cork proposal along with an initial deposit.
The proposal details must be supplied via a JSON file.

Example:
$ %s tx gov submit-proposal schedule-axelar-cork <path/to/proposal.json> --from=<key_or_address>

Where proposal.json contains:

{
  "title": "Dollary-doos LP Arbitrum Scheduled AxelarCork Proposal",
  "description": "I trust them, approve cork",
  "block_height": 100000,
  "chain_id": 42161,
  "target_contract_address": "0x123801a7D398351b8bE11C439e05C5B3259aeC9B",
  "contract_call_proto_json": "{\"cellar_id\":\"0x123801a7D398351b8bE11C439e05C5B3259aeC9B\",\"<cellar_type_name>\":{\"some_fuction\":{\"function_args\":{}},\"block_height\":12345}}",
  "deadline": 1706225320,
  "deposit": "10000usomm"
}

The contract_call_proto_json field must be the JSON representation of a ScheduleRequest, which is defined in Steward's protos. For more information, see the Steward API docs at https://github.com/peggyjv/steward.
`,
				version.AppName,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			proposal := types.AxelarScheduledCorkProposalWithDeposit{}
			contents, err := os.ReadFile(args[0])
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

			content := types.NewAxelarScheduledCorkProposal(
				proposal.Title,
				proposal.Description,
				proposal.BlockHeight,
				proposal.ChainId,
				proposal.TargetContractAddress,
				proposal.ContractCallProtoJson,
				proposal.Deadline,
			)
			if err := content.ValidateBasic(); err != nil {
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

func CmdSubmitAxelarCommunityPoolSpendProposal() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "axelar-community-pool-spend [proposal-file]",
		Args:  cobra.ExactArgs(2),
		Short: "Submit an Axelar community pool spend proposal",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Submit an Axelar community pool spend proposal along with an initial deposit.
The proposal details must be supplied via a JSON file. The funds from the community pool
will be bridged to the target EVM via Axelar to the supplied recipient address. Only one denomination
of Cosmos token can be sent. Fees will be removed from the balance by Axelar automatically.

Example:
$ %s tx gov submit-proposal axelar-community-pool-spend <path/to/proposal.json> --from=<key_or_address>

Where proposal.json contains:

{
	"title": "Community Pool Axelar Spend",
	"description": "Bridge me some tokens to Arbitrum!",
	"recipient": "0x0000000000000000000000000000000000000000",
	"chain_id": 42161,
	"amount": "20000usomm",
	"deposit": "10000usomm"
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

			proposal := types.AxelarCommunityPoolSpendProposalForCLI{}
			contents, err := os.ReadFile(args[0])
			if err != nil {
				return err
			}

			if err = clientCtx.Codec.UnmarshalJSON(contents, &proposal); err != nil {
				return err
			}

			if proposal.ChainId == 0 {
				return fmt.Errorf("chain ID must be non-zero")
			}

			if !common.IsHexAddress(proposal.Recipient) {
				return fmt.Errorf("recipient is not a valid EVM address")
			}

			amount, err := sdk.ParseCoinNormalized(proposal.Amount)
			if err != nil {
				return err
			}

			if !amount.IsPositive() {
				return fmt.Errorf("amount must be positive")
			}

			deposit, err := sdk.ParseCoinsNormalized(proposal.Deposit)
			if err != nil {
				return err
			}

			from := clientCtx.GetFromAddress()

			content := types.NewAxelarCommunitySpendProposal(
				proposal.Title,
				proposal.Description,
				proposal.Recipient,
				proposal.ChainId,
				amount,
			)

			msg, err := govtypesv1beta1.NewMsgSubmitProposal(content, deposit, from)
			if err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	return cmd
}

// GetCmdSubmitAddChainConfigurationProposal implements the command to submit a cellar id addition proposal
func GetCmdSubmitAddChainConfigurationProposal() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "add-axelar-chain-config [proposal-file]",
		Args:  cobra.ExactArgs(1),
		Short: "Submit an Axelar chain configuration addition proposal",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Submit an Axelar chain configuration addition proposal along with an initial deposit.
The proposal details must be supplied via a JSON file.

Example:
$ %s tx gov submit-proposal add-axelar-chain-config <path/to/proposal.json> --from=<key_or_address>

Where proposal.json contains:

{
  "title": "Enable Arbitrum proposal",
  "description": "allow cellars to be used on Arbitrum",
  "chain_configuration": {
	"name": "arbitrum",
	"id": 42161,
	"proxy_address": "0x0000000000000000000000000000000000000000",
	"bridge_fees": [
		{
			"denom": "usomm",
			"amount": "100000"
		}
	]
  },
  "deposit": "10000usomm"
}

Note that the "name" parameter should map to a "Chain Identifier" as defined by Axelar: https://docs.axelar.dev/dev/reference/mainnet-chain-names
Bridge fees for a given source denom and destination chain can be calculated here: https://docs.axelar.dev/resources/mainnet#cross-chain-relayer-gas-fee
`,
				version.AppName,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			proposal := types.AddChainConfigurationProposalWithDeposit{}
			contents, err := os.ReadFile(args[0])
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

			if err := proposal.ChainConfiguration.ValidateBasic(); err != nil {
				return err
			}

			content := types.NewAddChainConfigurationProposal(
				proposal.Title,
				proposal.Description,
				*proposal.ChainConfiguration,
			)

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

// GetCmdSubmitRemoveChainConfigurationProposal implements the command to submit a chain configuration removal proposal
func GetCmdSubmitRemoveChainConfigurationProposal() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "remove-axelar-chain-config [proposal-file]",
		Args:  cobra.ExactArgs(1),
		Short: "Submit an Axelear chain configuration removal proposal",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Submit an Axelar chain configuration removal proposal along with an initial deposit.
The proposal details must be supplied via a JSON file.

Example:
$ %s tx gov submit-proposal remove-axelar-chain-config <path/to/proposal.json> --from=<key_or_address>

Where proposal.json contains:

{
  "title": "Remove Arbitrum chain config",
  "description": "not using Arbitrum any more",
  "chain_id": 42161,
  "deposit": "10000usomm"
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

			proposal := types.RemoveChainConfigurationProposalWithDeposit{}
			contents, err := os.ReadFile(args[0])
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

			if proposal.ChainId == 0 {
				return fmt.Errorf("chain ID cannot be zero")
			}

			content := types.NewRemoveChainConfigurationProposal(
				proposal.Title,
				proposal.Description,
				proposal.ChainId,
			)

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

func GetCmdSubmitUpgradeAxelarProxyContractProposal() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "upgrade-axelar-proxy-contract [proposal-file]",
		Args:  cobra.ExactArgs(1),
		Short: "Submit an upgrade axelar proxy contract proposal",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Submit an Axelar proxy contract upgrade proposal along with an initial deposit.
The proposal details must be supplied via a JSON file.

Example:
$ %s tx gov submit-proposal upgrade-axelar-proxy-contract <path/to/proposal.json> --from=<key_or_address>

Where proposal.json contains:

{
  "title": "Upgrade Axelar Proxy Contract Proposal",
  "description": "New features",
  "chain_id": 1000,
  "new_proxy_address": "0x1234567890123456789012345678901234567890",
  "deposit": "10000usomm"
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

			proposal := types.UpgradeAxelarProxyContractProposalWithDeposit{}
			contents, err := os.ReadFile(args[0])
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

			if proposal.ChainId == 0 {
				return fmt.Errorf("chain ID cannot be zero")
			}

			if !common.IsHexAddress(proposal.NewProxyAddress) {
				return fmt.Errorf("new proxy address is not a valid EVM address")
			}

			content := types.NewUpgradeAxelarProxyContractProposal(
				proposal.Title,
				proposal.Description,
				proposal.ChainId,
				proposal.NewProxyAddress,
			)

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

func GetCmdSubmitCancelAxelarProxyContractUpgradeProposal() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "cancel-axelar-proxy-contract-upgrade [proposal-file]",
		Args:  cobra.ExactArgs(1),
		Short: "Submit a cancel axelar proxy contract upgrade proposal",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Submit a cancellation proposal for an existing Axelar proxy contract upgrade along with an initial deposit.
The proposal details must be supplied via a JSON file.

Example:
$ %s tx gov submit-proposal cancel-axelar-proxy-contract-upgrade <path/to/proposal.json> --from=<key_or_address>

Where proposal.json contains:

{
  "title": "Cancel Upgrade Axelar Proxy Contract Proposal",
  "description": "Cancel the new features",
  "chain_id": 1000,
  "deposit": "10000usomm"
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

			proposal := types.CancelAxelarProxyContractUpgradeProposalWithDeposit{}
			contents, err := os.ReadFile(args[0])
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

			if proposal.ChainId == 0 {
				return fmt.Errorf("chain ID cannot be zero")
			}

			content := types.NewCancelAxelarProxyContractUpgradeProposal(
				proposal.Title,
				proposal.Description,
				proposal.ChainId,
			)

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
