package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/peggyjv/sommelier/v9/x/auction/types"
)

// HandleSetTokenPricesProposal is a handler for executing a passed community token price update proposal
func HandleSetTokenPricesProposal(ctx sdk.Context, k Keeper, p types.SetTokenPricesProposal) error {
	err := p.ValidateBasic()
	if err != nil {
		return err
	}

	for _, tokenPrice := range p.TokenPrices {
		k.setTokenPrice(ctx, types.TokenPrice{
			Denom:            tokenPrice.Denom,
			UsdPrice:         tokenPrice.UsdPrice,
			Exponent:         tokenPrice.Exponent,
			LastUpdatedBlock: uint64(ctx.BlockHeight()),
		})
	}

	return nil
}
