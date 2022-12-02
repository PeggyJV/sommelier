package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	auctiontypes "github.com/peggyjv/sommelier/v4/x/auction/types"
)

// Attempts to start an auction for the provided denom
func (k Keeper) beginAuction(ctx sdk.Context, denom string) (started bool) {
	activeAuctions := k.auctionKeeper.GetActiveAuctions(ctx)

	// Don't start an auction if the denom has an active one
	for _, auction := range activeAuctions {
		if denom == auction.StartingTokensForSale.Denom {
			return false
		}
	}

	// We auction the entire balance in the cellarfees module account
	cellarfeesParams := k.GetParams(ctx)
	cellarfeesAccountAddr := k.GetFeesAccount(ctx).GetAddress()
	balance := k.bankKeeper.GetBalance(ctx, cellarfeesAccountAddr, denom)
	if balance.IsZero() {
		k.Logger(ctx).Error("Attempted to begin auction for denom %s with a zero balance.", denom)
		return false
	}

	err := k.auctionKeeper.BeginAuction(
		ctx,
		balance,
		cellarfeesParams.InitialPriceDecreaseRate,
		cellarfeesParams.PriceDecreaseBlockInterval,
		string(cellarfeesAccountAddr),
		string(cellarfeesAccountAddr),
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
