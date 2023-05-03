package cli

import (
	"os"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/peggyjv/sommelier/v6/x/axelarcork/types"
	"github.com/spf13/cobra"
)

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

// ParseCommunityPoolSpendProposal reads and parses a CommunityPoolEthereumSpendProposalForCLI from a file.
func ParseCommunityPoolSpendProposal(cdc codec.JSONCodec, proposalFile string) (types.CommunityPoolSpendProposalForCLI, error) {
	proposal := types.CommunityPoolSpendProposalForCLI{}

	contents, err := os.ReadFile(proposalFile)
	if err != nil {
		return proposal, err
	}

	if err = cdc.UnmarshalJSON(contents, &proposal); err != nil {
		return proposal, err
	}

	return proposal, nil
}
