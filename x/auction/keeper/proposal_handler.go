package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/peggyjv/sommelier/v4/x/auction/types"
)

// HandleSetTokenPricesProposal is a handler for executing a passed community token price update proposal
func HandleSetTokenPricesProposal(ctx sdk.Context, k Keeper, p types.SetTokenPricesProposal) error {
	err := p.ValidateBasic()
	if err != nil {
		return err
	}

	seenDenomPrices := []string{}

	for _, tokenPrice := range p.TokenPrices {
		// Check if this price proposal attempts to update the same denom price twice
		for _, seenDenom := range seenDenomPrices {
			if seenDenom == tokenPrice.Denom {
				return sdkerrors.Wrapf(types.ErrTokenPriceProposalAttemptsToUpdateTokenPriceMoreThanOnce, "denom: %s", tokenPrice.Denom)
			}
		}

		seenDenomPrices = append(seenDenomPrices, tokenPrice.Denom)

		k.setTokenPrice(ctx, types.TokenPrice{
			Denom:            tokenPrice.Denom,
			UsdPrice:         tokenPrice.UsdPrice,
			LastUpdatedBlock: uint64(ctx.BlockHeight()),
		})
	}

	return nil
}
