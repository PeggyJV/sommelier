package cli

import "github.com/spf13/cobra"

const (
	FlagAxelarChainID   = "axelar-chain-id"
	FlagAxelarChainName = "axelar-chain-name"
)

// AddChainFlagsToCmd adds common chain flags to a module command.
func AddChainFlagsToCmd(cmd *cobra.Command) {
	cmd.Flags().String(FlagAxelarChainName, "", "The case sensitive name of the Axelar target chain")
	cmd.Flags().Uint64(FlagAxelarChainID, 0, "The Chain ID for the Axelar target EVM chain")
}

func GetChainInfoFromFlags(cmd *cobra.Command) (string, uint64, error) {
	name, err := cmd.Flags().GetString(FlagAxelarChainName)
	if err != nil {
		return "", 0, err
	}
	chainID, err := cmd.Flags().GetUint64(FlagAxelarChainID)
	if err != nil {
		return "", 0, err
	}

	return name, chainID, nil
}
