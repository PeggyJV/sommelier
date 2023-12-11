package types

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	govtypesv1beta1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1beta1"
)

const (
	ProposalTypeSetTokenPrices = "SetTokenPrices"
)

var _ govtypesv1beta1.Content = &SetTokenPricesProposal{}

func init() {
	govtypesv1beta1.RegisterProposalType(ProposalTypeSetTokenPrices)
	// The RegisterProposalTypeCodec function was mysteriously removed by in 0.46.0 even though
	// the claim was that the old API would be preserved in .../x/gov/types/v1beta1 so we have
	// to interact with the codec directly.
	//
	// The PR that removed it: https://github.com/cosmos/cosmos-sdk/pull/11240
	// This PR was later reverted, but RegisterProposalTypeCodec was still left out. Not sure if
	// this was intentional or not.
	govtypesv1beta1.ModuleCdc.RegisterConcrete(&SetTokenPricesProposal{}, "sommelier/SetTokenPricesProposal", nil)
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
