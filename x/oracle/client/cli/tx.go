package cli

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/peggyjv/sommelier/x/oracle/types"

	"github.com/spf13/cobra"
)

// GetTxCmd returns the transaction commands for this module
func GetTxCmd() *cobra.Command {
	oracleTxCmd := &cobra.Command{
		Use:                        "oracle",
		Short:                      "Oracle transaction subcommands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	oracleTxCmd.AddCommand(
		txDelegateFeedPermission(),
	)

	return oracleTxCmd
}

func txDelegateFeedPermission() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "delegate-feeder [address]",
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			del, err := sdk.AccAddressFromBech32(args[0])
			if err != nil {
				return err
			}

			msg := types.NewMsgDelegateFeedConsent(ctx.GetFromAddress(), del)
			if err = msg.ValidateBasic(); err != nil {
				return fmt.Errorf("message validation failed: %w", err)
			}

			return tx.GenerateOrBroadcastTxCLI(ctx, cmd.Flags(), msg)

		},
	}
	flags.AddTxFlagsToCmd(cmd)
	return cmd
}
