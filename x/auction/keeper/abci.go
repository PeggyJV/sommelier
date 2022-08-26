package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// BeginBlocker is called at the beginning of every block
func (k Keeper) BeginBlocker(ctx sdk.Context) {
	for _, auction := range k.GetActiveAuctions(ctx) {
		// TODO Step function for auction price updates

	}


	// TODO: do we need to do anything with proposal voting here?
}

// EndBlocker is called at the end of every block
func (k Keeper) EndBlocker(ctx sdk.Context) {
	// TODO End Auction that have no funds left
	for _, auction := range k.GetActiveAuctions(ctx) {
		// TODO: Move auction to ended auctions list with updated fields

		// TODO: Send proceeds to their appropriate destination module

	}

	// TODO: anything else? Trim down old auctions and bids maybe?
}
