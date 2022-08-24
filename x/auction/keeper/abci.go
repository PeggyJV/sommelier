package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// TODO: fill out BeginBlocker and EndBlocker

// BeginBlocker is called at the beginning of every block
func (k Keeper) BeginBlocker(ctx sdk.Context) {
}

// EndBlocker is called at the end of every block
func (k Keeper) EndBlocker(ctx sdk.Context) {
}