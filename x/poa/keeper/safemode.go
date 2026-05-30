package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/peggyjv/sommelier/v9/x/poa/types"
)

// SetSafeMode persists the authority-empty safe-mode flag.
func (k Keeper) SetSafeMode(ctx sdk.Context, active bool) {
	store := ctx.KVStore(k.storeKey)
	if active {
		store.Set(types.SafeModeKey, []byte{1})
		return
	}
	store.Delete(types.SafeModeKey)
}

// SafeModeActive reports whether the chain is in authority-empty safe mode.
// While true, value-bearing modules (gravity, cork, axelarcork) freeze their
// operations so they are never committed under an untrusted, community-only
// validator set. Reads are cheap (single key lookup) and safe to call per-tx
// from the ante handler and per-block from consumer EndBlockers.
func (k Keeper) SafeModeActive(ctx sdk.Context) bool {
	return ctx.KVStore(k.storeKey).Has(types.SafeModeKey)
}

// enterSafeMode sets the flag and emits the transition event only on the
// edge into safe mode, so the event is not re-emitted every block.
func (k Keeper) enterSafeMode(ctx sdk.Context) {
	if k.SafeModeActive(ctx) {
		return
	}
	k.SetSafeMode(ctx, true)
	ctx.Logger().Error("poa: authority set empty — entering safe mode; value-bearing modules frozen until a trusted authority set is restored via governance")
	ctx.EventManager().EmitEvent(sdk.NewEvent(types.EventTypeSafeModeEntered))
}

// exitSafeModeIfActive clears the flag and emits the transition event only on
// the edge out of safe mode.
func (k Keeper) exitSafeModeIfActive(ctx sdk.Context) {
	if !k.SafeModeActive(ctx) {
		return
	}
	k.SetSafeMode(ctx, false)
	ctx.Logger().Info("poa: authority set restored — exiting safe mode; value-bearing modules resume")
	ctx.EventManager().EmitEvent(sdk.NewEvent(types.EventTypeSafeModeExited))
}
