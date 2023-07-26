package cli

import (
	"fmt"
	"os"
	"strings"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/version"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	"github.com/ethereum/go-ethereum/common"
	types "github.com/peggyjv/sommelier/v6/x/axelarcork/types"
	"github.com/spf13/cobra"
)

// GetTxCmd returns the transaction commands for this module
func GetTxCmd() *cobra.Command {
	corkTxCmd := &cobra.Command{
		Use:                        "axelar-cork",
		Short:                      "AxelarCork transaction subcommands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	return corkTxCmd
}

// GetCmdSubmitAddProposal implements the command to submit a cellar id addition proposal
func GetCmdSubmitAddProposal() *cobra.Command {

	cmd := &cobra.Command{
		Use:   "add-cellar-id [proposal-file]",
		Args:  cobra.ExactArgs(1),
		Short: "Submit a cellar id addition proposal",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Submit a cellar addition proposal along with an initial deposit.
The proposal details must be supplied via a JSON file.

Example:
$ %s tx gov submit-proposal add-cellar-id <path/to/proposal.json> --from=<key_or_address>

Where proposal.json contains:

{
  "title": "Dollary-doos LP Cellar Proposal",
  "description": "I have a hunch",
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

			for _, id := range proposal.CellarIds.Ids {
				if !common.IsHexAddress(id) {
					return fmt.Errorf("%s is not a valid ethereum address", id)
				}
			}

			chainID, err := GetChainInfoFromFlags(cmd)
			if err != nil {
				return err
			}

			content := types.NewAddManagedCellarIDsProposal(
				proposal.Title,
				proposal.Description,
				chainID,
				&types.CellarIDSet{Ids: proposal.CellarIds.Ids})

			from := clientCtx.GetFromAddress()
			msg, err := govtypes.NewMsgSubmitProposal(content, deposit, from)
			if err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	AddChainFlagsToCmd(cmd)

	return cmd
}

// GetCmdSubmitRemoveProposal implements the command to submit a cellar id removal proposal
func GetCmdSubmitRemoveProposal() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "remove-cellar-id [proposal-file]",
		Args:  cobra.ExactArgs(1),
		Short: "Submit a cellar id removal proposal",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Submit a cellar removal proposal along with an initial deposit.
The proposal details must be supplied via a JSON file.

Example:
$ %s tx gov submit-proposal remove-cellar-id <path/to/proposal.json> --from=<key_or_address>

Where proposal.json contains:

{
  "title": "Dollary-doos LP Cellar Removal Proposal",
  "description": "I don't trust them",
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

			for _, id := range proposal.CellarIds.Ids {
				if !common.IsHexAddress(id) {
					return fmt.Errorf("%s is not a valid ethereum address", id)
				}
			}

			chainID, err := GetChainInfoFromFlags(cmd)
			if err != nil {
				return err
			}

			content := types.NewRemoveManagedCellarIDsProposal(
				proposal.Title,
				proposal.Description,
				chainID,
				&types.CellarIDSet{Ids: proposal.CellarIds.Ids})

			from := clientCtx.GetFromAddress()
			msg, err := govtypes.NewMsgSubmitProposal(content, deposit, from)
			if err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	AddChainFlagsToCmd(cmd)

	return cmd
}

// GetCmdSubmitScheduledCorkProposal implements the command to submit scheduled cork proposal
func GetCmdSubmitScheduledCorkProposal() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "schedule-cork [proposal-file]",
		Args:  cobra.ExactArgs(1),
		Short: "Submit a scheduled cork proposal",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Submit a scheduled cork proposal along with an initial deposit.
The proposal details must be supplied via a JSON file.

Example:
$ %s tx gov submit-proposal schedule-cork <path/to/proposal.json> --from=<key_or_address>

Where proposal.json contains:

{
  "title": "Dollary-doos LP Scheduled AxelarCork Proposal",
  "description": "I trust them, approve cork",
  "block_height": 100000,
  "target_contract_address": "0x123801a7D398351b8bE11C439e05C5B3259aeC9B",
  "contract_call_proto_json": "{\"cellar_id\":\"0x123801a7D398351b8bE11C439e05C5B3259aeC9B\",\"<cellar_type_name>\":{\"some_fuction\":{\"function_args\":{}},\"block_height\":12345}}",
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

			if !common.IsHexAddress(proposal.TargetContractAddress) {
				return fmt.Errorf("%s is not a valid contract address", proposal.TargetContractAddress)
			}

			chainID, err := GetChainInfoFromFlags(cmd)
			if err != nil {
				return err
			}

			content := types.NewAxelarScheduledCorkProposal(proposal.Title, proposal.Description, proposal.BlockHeight,
				chainID, proposal.TargetContractAddress, proposal.ContractCallProtoJson)
			if err := content.ValidateBasic(); err != nil {
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

	AddChainFlagsToCmd(cmd)

	return cmd
}

func CmdSubmitCommunityPoolEthereumSpendProposal() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "community-pool-spend [proposal-file]",
		Args:  cobra.ExactArgs(1),
		Short: "Submit a community pool spend proposal",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Submit a community pool spend proposal along with an initial deposit.
The proposal details must be supplied via a JSON file. The funds from the community pool
will be bridged to the target EVM via Axelar to the supplied recipient address. Only one denomination
of Cosmos token can be sent. Fees will be removed from the balance by Axelar automatically.

Example:
$ %s tx gov submit-proposal community-pool-spend <path/to/proposal.json> --from=<key_or_address>

Where proposal.json contains:

{
	"title": "Community Pool Ethereum Spend",
	"description": "Bridge me some tokens to Ethereum!",
	"recipient": "0x0000000000000000000000000000000000000000",
	"chain_name": "Avalanche",
	"amount": "20000stake",
	"deposit": "1000stake"
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

			proposal, err := ParseCommunityPoolSpendProposal(clientCtx.Codec, args[0])
			if err != nil {
				return err
			}

			if len(proposal.Title) == 0 {
				return fmt.Errorf("title is empty")
			}

			if len(proposal.Description) == 0 {
				return fmt.Errorf("description is empty")
			}

			if !common.IsHexAddress(proposal.Recipient) {
				return fmt.Errorf("recipient is not a valid Ethereum address")
			}

			amount, err := sdk.ParseCoinNormalized(proposal.Amount)
			if err != nil {
				return err
			}

			deposit, err := sdk.ParseCoinsNormalized(proposal.Deposit)
			if err != nil {
				return err
			}

			from := clientCtx.GetFromAddress()

			content := types.NewAxelarCommunitySpendProposal(proposal.Title, proposal.Description, proposal.Recipient, proposal.ChainId, amount)

			msg, err := govtypes.NewMsgSubmitProposal(content, deposit, from)
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
		Use:   "add-chain-config [proposal-file]",
		Args:  cobra.ExactArgs(1),
		Short: "Submit a chain configuration addition proposal",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Submit a chain configuration addition proposal along with an initial deposit.
The proposal details must be supplied via a JSON file.

Example:
$ %s tx gov submit-proposal add-chain-config <path/to/proposal.json> --from=<key_or_address>

Where proposal.json contains:

{
  "title": "Enable Fake EVM proposal",
  "description": "Fake it 'til you make it'",
  "chain_congifuration": {
	"name": "FakeEVM",
	"id": 1000,
	"vote_threshold": "0.333",
	"proxy_address": "0x1234..."
  },
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
			msg, err := govtypes.NewMsgSubmitProposal(content, deposit, from)
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
		Use:   "remove-chain-configuration [proposal-file]",
		Args:  cobra.ExactArgs(1),
		Short: "Submit a chain configuration removal proposal",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Submit a chain configuration removal proposal along with an initial deposit.
The proposal details must be supplied via a JSON file.

Example:
$ %s tx gov submit-proposal remove-chain-configuration <path/to/proposal.json> --from=<key_or_address>

Where proposal.json contains:

{
  "title": "Dollary-doos LP Cellar Removal Proposal",
  "description": "I don't trust them",
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
				proposal.ChainId)

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
