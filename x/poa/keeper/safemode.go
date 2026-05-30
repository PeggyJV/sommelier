package keeper

import (
	"encoding/binary"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/peggyjv/sommelier/v9/x/poa/types"
)

// safeModeThawDelayBlocks is how long value-bearing modules stay frozen after
// the authority set is restored, before they may resume. CometBFT applies the
// validator updates emitted in block H's EndBlock at height H+2, so the
// re-boosted authority set only secures consensus from H+2 onward. Thawing
// earlier would let cork/axelarcork/gravity act in a block still validated by
// the old community-only set. 2 = that propagation delay.
const safeModeThawDelayBlocks int64 = 2

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
// validator set. The flag stays true through the post-restoration thaw delay.
// Reads are cheap (single key lookup) and safe to call per-tx from the ante
// handler and per-block from consumer EndBlockers.
func (k Keeper) SafeModeActive(ctx sdk.Context) bool {
	return ctx.KVStore(k.storeKey).Has(types.SafeModeKey)
}

func (k Keeper) setThawHeight(ctx sdk.Context, height int64) {
	bz := make([]byte, 8)
	binary.BigEndian.PutUint64(bz, uint64(height))
	ctx.KVStore(k.storeKey).Set(types.SafeModeThawHeightKey, bz)
}

func (k Keeper) getThawHeight(ctx sdk.Context) (int64, bool) {
	bz := ctx.KVStore(k.storeKey).Get(types.SafeModeThawHeightKey)
	if bz == nil {
		return 0, false
	}
	return int64(binary.BigEndian.Uint64(bz)), true
}

func (k Keeper) clearThawHeight(ctx sdk.Context) {
	ctx.KVStore(k.storeKey).Delete(types.SafeModeThawHeightKey)
}

// enterSafeMode sets the flag (and cancels any pending thaw, in case the
// authority set re-emptied during the thaw window). The transition event is
// emitted only on the edge into safe mode, not every block.
func (k Keeper) enterSafeMode(ctx sdk.Context) {
	alreadyActive := k.SafeModeActive(ctx)
	k.SetSafeMode(ctx, true)
	k.clearThawHeight(ctx)
	if !alreadyActive {
		ctx.Logger().Error("poa: authority set empty — entering safe mode; value-bearing modules frozen until a trusted authority set is restored via governance")
		ctx.EventManager().EmitEvent(sdk.NewEvent(types.EventTypeSafeModeEntered))
	}
}

// maybeThawSafeMode is called each block the authority set is healthy. If safe
// mode is active it schedules the thaw on the first such block and completes it
// once the restored validator set is securing consensus, so value-bearing
// modules resume only after the freeze is genuinely safe to lift.
func (k Keeper) maybeThawSafeMode(ctx sdk.Context) {
	if !k.SafeModeActive(ctx) {
		return
	}
	thaw, ok := k.getThawHeight(ctx)
	if !ok {
		k.setThawHeight(ctx, ctx.BlockHeight()+safeModeThawDelayBlocks)
		return
	}
	if ctx.BlockHeight() >= thaw {
		k.exitSafeMode(ctx)
	}
}

// exitSafeModeIfActive clears safe mode immediately (flag + pending thaw),
// emitting the transition event on the edge. Used when PoA is disabled, where
// no thaw delay applies.
func (k Keeper) exitSafeModeIfActive(ctx sdk.Context) {
	if !k.SafeModeActive(ctx) {
		k.clearThawHeight(ctx)
		return
	}
	k.exitSafeMode(ctx)
}

// exitSafeMode clears the flag and any pending thaw and emits the exit event.
func (k Keeper) exitSafeMode(ctx sdk.Context) {
	k.SetSafeMode(ctx, false)
	k.clearThawHeight(ctx)
	ctx.Logger().Info("poa: authority set restored — exiting safe mode; value-bearing modules resume")
	ctx.EventManager().EmitEvent(sdk.NewEvent(types.EventTypeSafeModeExited))
}
