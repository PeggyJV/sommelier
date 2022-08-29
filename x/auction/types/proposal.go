package types

import (
	"fmt"

	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
)

const (
	ProposalTypeSetTokenPrices = "SetTokenPrices"
)

var _ govtypes.Content = &SetTokenPricesProposal{}

func init() {
	govtypes.RegisterProposalType(ProposalTypeSetTokenPrices)
	govtypes.RegisterProposalTypeCodec(&SetTokenPricesProposal{}, "sommelier/SetTokenPricesProposal")
}

func NewSetTokenPricesProposal(title string, description string, proposedTokenPrices []*ProposedTokenPrice) *SetTokenPricesProposal {
	return &SetTokenPricesProposal{
		Title:       title,
		Description: description,
		TokenPrices: proposedTokenPrices,
	}
}

func (m *SetTokenPricesProposal) ProposalRoute() string {
	return RouterKey
}

func (m *SetTokenPricesProposal) ProposalType() string {
	return ProposalTypeSetTokenPrices
}

func (m *SetTokenPricesProposal) ValidateBasic() error {
	if err := govtypes.ValidateAbstract(m); err != nil {
		return err
	}

	if len(m.TokenPrices) == 0 {
		return fmt.Errorf("list of token prices must be non-zero")
	}

	for _, tokenPrice := range m.TokenPrices {
		tokenPriceError := tokenPrice.ValidateBasic()
		if tokenPriceError != nil {
			return tokenPriceError
		}
	}

	return nil
}
