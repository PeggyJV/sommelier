package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// TODO: fill out BeginBlocker and EndBlocker

// BeginBlocker is called at the beginning of every block
func (k Keeper) BeginBlocker(ctx sdk.Context) {
	// Step function for auction price updates
}

// EndBlocker is called at the end of every block
func (k Keeper) EndBlocker(ctx sdk.Context) {
	// End Auction that have no funds left
}
