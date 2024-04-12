package types

import (
	errorsmod "cosmossdk.io/errors"
	govtypesv1beta1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1beta1"
)

const (
	ProposalTypeSetTokenPrices = "SetTokenPrices"
)

var _ govtypesv1beta1.Content = &SetTokenPricesProposal{}

func init() {
	govtypesv1beta1.RegisterProposalType(ProposalTypeSetTokenPrices)
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
	if err := govtypesv1beta1.ValidateAbstract(m); err != nil {
		return err
	}

	if len(m.TokenPrices) == 0 {
		return errorsmod.Wrapf(ErrTokenPriceProposalMustHaveAtLeastOnePrice, "prices: %v", m.TokenPrices)
	}

	seenDenomPrices := []string{}

	for _, tokenPrice := range m.TokenPrices {
		// Check if this price proposal attempts to update the same denom price twice
		for _, seenDenom := range seenDenomPrices {
			if seenDenom == tokenPrice.Denom {
				return errorsmod.Wrapf(ErrTokenPriceProposalAttemptsToUpdateTokenPriceMoreThanOnce, "denom: %s", tokenPrice.Denom)
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
