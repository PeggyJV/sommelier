package cli

import (
	"fmt"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/cosmos/cosmos-sdk/version"

	// "github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/peggyjv/sommelier/v7/x/addresses/types"
)

var (
	DefaultRelativePacketTimeoutTimestamp = uint64((time.Duration(10) * time.Minute).Nanoseconds())
)

const (
	flagPacketTimeoutTimestamp = "packet-timeout-timestamp"
	listSeparator              = ","
)

// GetTxCmd returns the transaction commands for this module
func GetTxCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      fmt.Sprintf("%s transactions subcommands", types.ModuleName),
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	cmd.AddCommand([]*cobra.Command{
		GetCmdAddAddressMapping(),
		GetCmdRemoveAddressMapping(),
	}...)

	return cmd
}

// GetCmdAddAddressMapping implements the command to submit a token price set proposal
func GetCmdAddAddressMapping() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "add-address-mapping [evm-address]",
		Args:  cobra.ExactArgs(1),
		Short: "Add a mapping from your signer address to an EVM address",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Add a mapping from your Cosmos (signer) address to an EVM address.

Example:
$ %s tx addresses add-address-mapping 0x1111111111111111111111111111111111111111 --from=<signer_key_or_address>
`,
				version.AppName,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			if !common.IsHexAddress(args[0]) {
				return fmt.Errorf("%s is not a valid EVM address", args[0])
			}

			evmAddress := common.HexToAddress(args[0])

			from := clientCtx.GetFromAddress()
			msg, err := types.NewMsgAddAddressMapping(evmAddress, from)
			if err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)
	return cmd
}

// GetCmdRemoveAddressMapping implements the command to submit a token price set proposal
func GetCmdRemoveAddressMapping() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "remove-address-mapping",
		Args:  cobra.NoArgs,
		Short: "Remove the mapping from your signer address to an EVM address",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Remove the mapping from your Cosmos (signer) address to an EVM address.

Example:
$ %s tx addresses remove-address-mapping --from=<signer_key_or_address>
`,
				version.AppName,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			from := clientCtx.GetFromAddress()
			msg, err := types.NewMsgRemoveAddressMapping(from)
			if err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)
	return cmd
}
