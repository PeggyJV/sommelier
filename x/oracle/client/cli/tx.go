package cli

import (
	"github.com/cosmos/cosmos-sdk/client"

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

	// oracleTxCmd.AddCommand([]*cobra.Command{
	// 	GetCmdDelegateFeederPermission(),
	// }...)

	return oracleTxCmd
}
