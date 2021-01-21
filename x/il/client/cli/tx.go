package cli

import (
	"strconv"
	"strings"

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
	}...)

	return oracleTxCmd
}

// GetCmdCreateStoploss
func GetCmdCreateStoploss() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create-stoploss [uniswap-pair] [liquidity-pool-shares] [max-slippage] [ref-pair-ratio]",
		Args:  cobra.ExactArgs(4),
		Short: "Delegate the permission to vote for the oracle to an address",
		Long: strings.TrimSpace(`
Delegate the permission to submit exchange rate votes for the oracle to an address.

Delegation can keep your validator operator key offline and use a separate replaceable key online.

$ sommelier tx oracle set-feeder terra1...

where "terra1..." is the address you want to delegate your voting rights to.
`),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			uniswapPairID := args[0]

			shares, err := strconv.ParseInt(args[1], 10, 64)
			if err != nil {
				return err
			}

			maxSlippage, err := sdk.NewDecFromStr(args[2])
			if err != nil {
				return err
			}

			ratio, err := sdk.NewDecFromStr(args[3])
			if err != nil {
				return err
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
