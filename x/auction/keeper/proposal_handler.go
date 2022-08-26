package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/peggyjv/sommelier/v4/x/auction/types"
)

// HandleSetTokenPricesProposal is a handler for executing a passed community token price update proposal
func HandleSetTokenPricesProposal(ctx sdk.Context, k Keeper, p types.SetTokenPricesProposal) error {
	for _, tokenPrice := range p.TokenPrices {
		k.setTokenPrice(ctx, *tokenPrice)
	}
	
	return nil
}