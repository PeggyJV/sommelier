package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	auctiontypes "github.com/peggyjv/sommelier/v7/x/auction/types"
	"github.com/peggyjv/sommelier/v7/x/cellarfees/types"
)

// Attempts to start an auction for the provided denom
func (k Keeper) beginAuction(ctx sdk.Context, balance sdk.Coin) bool {
	activeAuctions := k.auctionKeeper.GetActiveAuctions(ctx)

	// Don't start an auction if the denom has an active one
	for _, auction := range activeAuctions {
		if balance.Denom == auction.StartingTokensForSale.Denom {
			return false
		}
	}

	cellarfeesParams := k.GetParams(ctx)
	if balance.IsZero() {
		k.Logger(ctx).Error("Attempted to begin auction for denom %s with a zero balance.", balance.Denom)
		return false
	}

	err := k.auctionKeeper.BeginAuction(
		ctx,
		balance,
		cellarfeesParams.InitialPriceDecreaseRate,
		cellarfeesParams.PriceDecreaseBlockInterval,
		types.ModuleName,
		types.ModuleName,
	)
	if err != nil {
		switch err {
		case auctiontypes.ErrUnauthorizedFundingModule:
			panic("Attempted to start an auction with an unauthorized funding module")
		case auctiontypes.ErrUnauthorizedProceedsModule:
			panic("Attempted to start an auction with an unauthorized proceeds module")
		default:
			k.Logger(ctx).Error(err.Error())
		}

		return false
	}

	return true
}
