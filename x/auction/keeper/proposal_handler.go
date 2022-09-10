package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/peggyjv/sommelier/v4/x/auction/types"
)

// HandleSetTokenPricesProposal is a handler for executing a passed community token price update proposal
func HandleSetTokenPricesProposal(ctx sdk.Context, k Keeper, p types.SetTokenPricesProposal) error {
	for _, tokenPrice := range p.TokenPrices {
		k.SetTokenPrice(ctx, types.TokenPrice{
			Denom:            tokenPrice.Denom,
			UsdPrice:         tokenPrice.UsdPrice,
			LastUpdatedBlock: uint64(ctx.BlockHeight()),
		})
	}

	return nil
}
