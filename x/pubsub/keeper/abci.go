package keeper

import sdk "github.com/cosmos/cosmos-sdk/types"

func (k Keeper) BeginBlocker(ctx sdk.Context) {}
func (k Keeper) EndBlocker(ctx sdk.Context)   {}
