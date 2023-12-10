package cli

import (
	"os"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/peggyjv/sommelier/v7/x/axelarcork/types"
)

// ParseCommunityPoolSpendProposal reads and parses a CommunityPoolEthereumSpendProposalForCLI from a file.
func ParseCommunityPoolSpendProposal(cdc codec.JSONCodec, proposalFile string) (types.AxelarCommunityPoolSpendProposalForCLI, error) {
	proposal := types.AxelarCommunityPoolSpendProposalForCLI{}

	contents, err := os.ReadFile(proposalFile)
	if err != nil {
		return proposal, err
	}

	if err = cdc.UnmarshalJSON(contents, &proposal); err != nil {
		return proposal, err
	}

	return proposal, nil
}
