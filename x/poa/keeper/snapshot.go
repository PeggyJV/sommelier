package keeper

import (
	"encoding/binary"

	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/peggyjv/sommelier/v9/x/poa/types"
)

// SetMultiplierSnapshot persists a per-block snapshot of authority boost
// multipliers, used at slash time to convert boosted consensus power back to
// raw stake.
func (k Keeper) SetMultiplierSnapshot(ctx sdk.Context, s types.MultiplierSnapshot) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.MultiplierSnapshotKey(s.Height), k.cdc.MustMarshal(&s))
}

// GetMultiplierSnapshot returns the snapshot at `height` if one exists.
func (k Keeper) GetMultiplierSnapshot(ctx sdk.Context, height int64) (types.MultiplierSnapshot, bool) {
	bz := ctx.KVStore(k.storeKey).Get(types.MultiplierSnapshotKey(height))
	if bz == nil {
		return types.MultiplierSnapshot{}, false
	}
	var s types.MultiplierSnapshot
	k.cdc.MustUnmarshal(bz, &s)
	return s, true
}

// MultiplierForValidatorWithStatus returns the boost recorded for `val` at
// `height` along with a found-bool indicating whether the snapshot exists.
//   - snapFound=false: no snapshot at that height (likely pruned)
//   - snapFound=true, multiplier=1: snapshot exists but the validator was not
//     boosted on that block
func (k Keeper) MultiplierForValidatorWithStatus(ctx sdk.Context, val sdk.ValAddress, height int64) (sdk.Dec, bool) {
	s, ok := k.GetMultiplierSnapshot(ctx, height)
	if !ok {
		return sdk.OneDec(), false
	}
	return multiplierFromSnapshot(s, val.String())
}

// multiplierFromSnapshot extracts the boost for `addr` from an already-loaded
// snapshot. The returned bool reports whether the snapshot is usable:
//   - true, multiplier=m>1: validator was boosted
//   - true, multiplier=1: snapshot present, validator not boosted
//   - false: the validator's entry is corrupt (treat as missing snapshot so
//     the slash is refused rather than normalizing with a wrong multiplier)
func multiplierFromSnapshot(s types.MultiplierSnapshot, addr string) (sdk.Dec, bool) {
	for _, e := range s.Entries {
		if e == nil {
			continue
		}
		if e.OperatorAddress == addr {
			d, err := sdk.NewDecFromStr(e.Multiplier)
			if err != nil {
				return sdk.OneDec(), false
			}
			return d, true
		}
	}
	return sdk.OneDec(), true
}

// LatestMultiplierForValidator returns the boost for `val` from the most recent
// snapshot at height <= maxHeight. This is the phase-independent lookup used by
// boosted reads: during BeginBlock/DeliverTx of block N the snapshot for N is
// not yet written, so the latest committed snapshot (N-1) is returned, matching
// the LastValidatorPower the staking store still holds from N-1's EndBlocker;
// during EndBlock of N (after PoA has run) the snapshot for N is present and
// returned, matching the freshly-overwritten LastValidatorPower. Querying a
// fixed `BlockHeight()` instead would miss the boost for the whole of block N.
func (k Keeper) LatestMultiplierForValidator(ctx sdk.Context, val sdk.ValAddress, maxHeight int64) (sdk.Dec, bool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.MultiplierSnapshotPrefix)
	// Exclusive upper bound of maxHeight+1 so the iterator includes maxHeight.
	end := make([]byte, 8)
	binary.BigEndian.PutUint64(end, uint64(maxHeight)+1)
	iter := store.ReverseIterator(nil, end)
	defer iter.Close()
	if !iter.Valid() {
		return sdk.OneDec(), false
	}
	var s types.MultiplierSnapshot
	k.cdc.MustUnmarshal(iter.Value(), &s)
	return multiplierFromSnapshot(s, val.String())
}

// MultiplierForValidator is the snapshot-found-true variant for callers that
// just want the boost value without inspecting found-status.
func (k Keeper) MultiplierForValidator(ctx sdk.Context, val sdk.ValAddress, height int64) sdk.Dec {
	m, _ := k.MultiplierForValidatorWithStatus(ctx, val, height)
	return m
}

// PruneSnapshotsBefore deletes all snapshots with height < `height`.
func (k Keeper) PruneSnapshotsBefore(ctx sdk.Context, height int64) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.MultiplierSnapshotPrefix)
	iter := store.Iterator(nil, nil)
	defer iter.Close()
	for ; iter.Valid(); iter.Next() {
		key := iter.Key()
		if len(key) != 8 {
			continue
		}
		h := int64(binary.BigEndian.Uint64(key))
		if h >= height {
			break // big-endian iteration: keys are height-sorted
		}
		store.Delete(key)
	}
}
