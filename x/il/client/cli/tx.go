package cli

import (
	"fmt"
	"strconv"

	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/peggyjv/sommelier/x/il/types"
)

// GetTxCmd returns the transaction commands for this module
func GetTxCmd() *cobra.Command {
	oracleTxCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "Impermanent Loss transaction subcommands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	oracleTxCmd.AddCommand([]*cobra.Command{
		GetCmdCreateStoploss(),
		GetCmdDeleteStoploss(),
	}...)

	return oracleTxCmd
}

// GetCmdCreateStoploss
func GetCmdCreateStoploss() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create-stoploss [uniswap-pair] [liquidity-pool-shares] [max-slippage] [ref-pair-ratio]",
		Args:  cobra.ExactArgs(4),
		Short: "Create a stoploss position for a given LP provider",
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			uniswapPairID := args[0]

			shares, err := strconv.ParseInt(args[1], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid liquidity pool shares: %w", err)
			}

			maxSlippage, err := sdk.NewDecFromStr(args[2])
			if err != nil {
				return fmt.Errorf("invalid max slippage: %w", err)
			}

			ratio, err := sdk.NewDecFromStr(args[3])
			if err != nil {
				return fmt.Errorf("invalid reference pair ratio: %w", err)
			}

			stoploss := &types.Stoploss{
				UniswapPairId:       uniswapPairID,
				LiquidityPoolShares: shares,
				MaxSlippage:         maxSlippage,
				ReferencePairRatio:  ratio,
			}

			msg := types.NewMsgStoploss(clientCtx.FromAddress, stoploss)
			err = msg.ValidateBasic()
			if err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	return cmd
}

// GetCmdDeleteStoploss
func GetCmdDeleteStoploss() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete-stoploss [uniswap-pair]",
		Args:  cobra.ExactArgs(1),
		Short: "Delete an existing stoploss position",
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			uniswapPairID := args[0]

			msg := types.NewMsgDeleteStoploss(clientCtx.FromAddress, uniswapPairID)
			err = msg.ValidateBasic()
			if err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	return cmd
}
