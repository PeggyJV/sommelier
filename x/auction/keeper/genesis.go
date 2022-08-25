package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/peggyjv/sommelier/v4/x/auction/types"
)

// InitGenesis initializes the module's state from a provided genesis
// state.
func InitGenesis(ctx sdk.Context, k Keeper, gs types.GenesisState) {
	k.setParams(ctx, gs.Params)

	for _, auction := range gs.Auctions {
		if auction.EndBlock > 0 {
			k.setActiveAuction(ctx, *auction)
		} else {
			k.setEndedAuction(ctx, *auction)
		}
	}
	for _, bid := range gs.Bids {
		k.setBid(ctx, *bid)
	}

	for _, tokenPrice := range gs.TokenPrices {
		k.setTokenPrice(ctx, *tokenPrice)
	}

	k.setLastAuctionId(gs.LastAuctionId)
	k.setLastBidId(gs.LastBidId)
	k.setLastTokenPriceUpdateBlock(gs.LastTokenPriceUpdateBlock)
}

// ExportGenesis returns the module's exported genesis.
func ExportGenesis(ctx sdk.Context, k Keeper) types.GenesisState {
	var auctions []*types.Auction
	auctions = append(auctions, k.GetActiveAuctions(ctx)...)
	auctions = append(auctions, k.GetEndedAuctions(ctx)...)

	// TODO: Recompile proto 
	return types.GenesisState{
		Params:      				k.GetParamSet(ctx),
		Auctions:    				auctions,
		Bids:        				k.GetBids(ctx),
		TokenPrices: 				k.GetTokenPrices(ctx),
		LastAuctionId:  			k.GetLastAuctionId(),
		LastBidId:      			k.GetLastBidId(),
		LastTokenPriceUpdateBlock:  k.GetLastTokenPriceUpdateBlock(),
	}
}
