package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/peggyjv/sommelier/v4/x/auction/types"
)

// InitGenesis initializes the module's state from a provided genesis
// state.
func InitGenesis(ctx sdk.Context, k Keeper, gs types.GenesisState) {
	k.SetParams(ctx, gs.Params)
	k.SetAuctions(ctx, gs.Auctions)
	k.SetBids(ctx, gs.Bids)
	k.SetTokenPrices(ctx, gs.TokenPrices)
}

// ExportGenesis returns the module's exported genesis.
func ExportGenesis(ctx sdk.Context, k Keeper) types.GenesisState {
	return types.GenesisState{
		Params:      k.GetParamSet(ctx),
		Auctions:    k.GetAllAuctions(ctx),
		Bids:        k.GetBids(ctx),
		TokenPrices: k.GetTokenPrices(ctx),
	}
}
