package keeper

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/peggyjv/sommelier/v9/x/auction/types"
)

// InitGenesis initializes the module's state from a provided genesis
// state.
func InitGenesis(ctx sdk.Context, k Keeper, gs types.GenesisState) {
	k.setParams(ctx, gs.Params)

	auctionAccount := k.GetAuctionAccount(ctx)
	if auctionAccount == nil {
		panic(fmt.Sprintf("%s module account has not been set", types.ModuleName))
	}

	for _, auction := range gs.Auctions {
		if auction.EndBlock > 0 {
			k.setEndedAuction(ctx, *auction)
		} else {
			k.setActiveAuction(ctx, *auction)
		}
	}
	for _, bid := range gs.Bids {
		k.setBid(ctx, *bid)
	}

	for _, tokenPrice := range gs.TokenPrices {
		k.setTokenPrice(ctx, *tokenPrice)
	}

	k.setLastAuctionID(ctx, gs.LastAuctionId)
	k.setLastBidID(ctx, gs.LastBidId)
}

// ExportGenesis returns the module's exported genesis.
func ExportGenesis(ctx sdk.Context, k Keeper) types.GenesisState {
	var auctions []*types.Auction
	auctions = append(auctions, k.GetActiveAuctions(ctx)...)
	auctions = append(auctions, k.GetEndedAuctions(ctx)...)

	return types.GenesisState{
		Params:        k.GetParamSet(ctx),
		Auctions:      auctions,
		Bids:          k.GetBids(ctx),
		TokenPrices:   k.GetTokenPrices(ctx),
		LastAuctionId: k.GetLastAuctionID(ctx),
		LastBidId:     k.GetLastBidID(ctx),
	}
}
