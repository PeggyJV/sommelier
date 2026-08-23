package cli

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/version"
	govtypesv1beta1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1beta1"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	types "github.com/peggyjv/sommelier/v10/x/cork/types/v2"
	pubsubtypes "github.com/peggyjv/sommelier/v10/x/pubsub/types"
	"github.com/spf13/cobra"
)

// GetTxCmd returns the transaction commands for this module
func GetTxCmd() *cobra.Command {
	corkTxCmd := &cobra.Command{
		Use:                        "cork",
		Short:                      "Cork transaction subcommands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	corkTxCmd.AddCommand(CmdScheduleCork())

	return corkTxCmd
}

// buildScheduleCorkMsg assembles a MsgScheduleCorkRequest from CLI input.
//
// Split out from the cobra command so the encoding rules are testable without a
// client context. The encoded call is the part worth guarding: the chain expects
// raw ABI-encoded bytes, so the hex must be DECODED here. Handing the message
// []byte(hexString) instead would pass the ASCII of the hex through to the
// cellar, which reverts on Ethereum with nothing locally to explain why.
func buildScheduleCorkMsg(signer, contractAddr string, blockHeight uint64, encodedCall string) (*types.MsgScheduleCorkRequest, error) {
	if !common.IsHexAddress(contractAddr) {
		return nil, fmt.Errorf("target contract address %s is invalid", contractAddr)
	}

	callBz, err := hexutil.Decode(withHexPrefix(encodedCall))
	if err != nil {
		return nil, fmt.Errorf("contract call must be hex-encoded ABI bytes: %w", err)
	}
	if len(callBz) == 0 {
		return nil, fmt.Errorf("contract call is empty; an empty body is never a valid cellar instruction")
	}

	msg := &types.MsgScheduleCorkRequest{
		Cork: &types.Cork{
			EncodedContractCall:   callBz,
			TargetContractAddress: contractAddr,
		},
		BlockHeight: blockHeight,
		Signer:      signer,
	}
	if err := msg.ValidateBasic(); err != nil {
		return nil, err
	}
	return msg, nil
}

func withHexPrefix(s string) string {
	if strings.HasPrefix(s, "0x") || strings.HasPrefix(s, "0X") {
		return s
	}
	return "0x" + s
}

// CmdScheduleCork schedules a cork against an Ethereum cellar.
//
// Only the cork_authority param may do this: v10 retires the
// validator-delegate path that steward used, so this command is the supported
// way to drive the cellar wind-down.
func CmdScheduleCork() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "schedule-cork [target-contract-address] [block-height] [encoded-contract-call]",
		Args:  cobra.ExactArgs(3),
		Short: "Schedule a cork against an Ethereum cellar (cork authority only)",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Schedule a cork to be submitted to Ethereum at a future block height.

The signing account MUST be the cork_authority param; any other signer is
rejected. The target contract must be on the managed-cellar allowlist, and the
allowlist is re-checked at execution, so removing a cellar cancels calls already
queued against it.

The contract call is ABI-encoded bytes, given as hex (0x prefix optional).

Example:
$ %s tx cork schedule-cork 0x123801a7D398351b8bE11C439e05C5B3259aeC9B 27500000 0xa9059cbb... --from cork-authority
`, version.AppName)),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			blockHeight, err := strconv.ParseUint(args[1], 10, 64)
			if err != nil {
				return fmt.Errorf("block height %s is invalid: %w", args[1], err)
			}

			msg, err := buildScheduleCorkMsg(
				clientCtx.GetFromAddress().String(), args[0], blockHeight, args[2])
			if err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)
	return cmd
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
  "publisher_domain": "example.com",
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

			proposal := types.AddManagedCellarIDsProposalWithDeposit{}
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

			for _, id := range proposal.CellarIds {
				if !common.IsHexAddress(id) {
					return fmt.Errorf("%s is not a valid ethereum address", id)
				}
			}

			if err := pubsubtypes.ValidateDomain(proposal.PublisherDomain); err != nil {
				return err
			}

			content := types.NewAddManagedCellarIDsProposal(
				proposal.Title,
				proposal.Description,
				&types.CellarIDSet{Ids: proposal.CellarIds},
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

			proposal := types.RemoveManagedCellarIDsProposalWithDeposit{}
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

			for _, id := range proposal.CellarIds {
				if !common.IsHexAddress(id) {
					return fmt.Errorf("%s is not a valid ethereum address", id)
				}
			}

			content := types.NewRemoveManagedCellarIDsProposal(
				proposal.Title,
				proposal.Description,
				&types.CellarIDSet{Ids: proposal.CellarIds})

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
  "title": "Dollary-doos LP Scheduled Cork Proposal",
  "description": "I trust them, approve cork",
  "block_height": 100000,
  "target_contract_address": "0x123801a7D398351b8bE11C439e05C5B3259aeC9B",
  "contract_call_proto_json": "{\"cellar_id\":\"0x123801a7D398351b8bE11C439e05C5B3259aeC9B\",\"<cellar_type_name>\":{\"some_fuction\":{\"function_args\":{}},\"block_height\":12345}}",
  "deposit": "10000000usomm"
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

			proposal := types.ScheduledCorkProposalWithDeposit{}
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

			content := types.NewScheduledCorkProposal(proposal.Title, proposal.Description, proposal.BlockHeight, proposal.TargetContractAddress, proposal.ContractCallProtoJson)
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
