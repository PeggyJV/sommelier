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
	addr := val.String()
	for _, e := range s.Entries {
		if e == nil {
			continue
		}
		if e.OperatorAddress == addr {
			d, err := sdk.NewDecFromStr(e.Multiplier)
			if err != nil {
				return sdk.OneDec(), true
			}
			return d, true
		}
	}
	return sdk.OneDec(), true
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
