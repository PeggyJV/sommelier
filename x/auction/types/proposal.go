package types

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
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
		return sdkerrors.Wrapf(ErrTokenPriceProposalMustHaveAtLeastOnePrice, "prices: %v", m.TokenPrices)
	}

	seenDenomPrices := []string{}

	for _, tokenPrice := range m.TokenPrices {
		// Check if this price proposal attempts to update the same denom price twice
		for _, seenDenom := range seenDenomPrices {
			if seenDenom == tokenPrice.Denom {
				return sdkerrors.Wrapf(ErrTokenPriceProposalAttemptsToUpdateTokenPriceMoreThanOnce, "denom: %s", tokenPrice.Denom)
			}
		}

		seenDenomPrices = append(seenDenomPrices, tokenPrice.Denom)

		tokenPriceError := tokenPrice.ValidateBasic()
		if tokenPriceError != nil {
			return tokenPriceError
		}
	}

	return nil
}
