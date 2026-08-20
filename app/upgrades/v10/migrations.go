package v10

import (
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"

	axelarcorktypes "github.com/peggyjv/sommelier/v10/x/axelarcork/types"
	corktypes "github.com/peggyjv/sommelier/v10/x/cork/types"
)

// DrainLegacyCorkQueues deletes every cork left in the pre-v10
// validator-scheduled queues, along with the per-validator cork counters.
//
// This MUST run in the v10 handler. Before v10, the power tally in each
// module's EndBlocker was the only site that deleted from these prefixes. v10
// removes the tally, so anything left behind becomes permanently undeletable
// state: no code path reads it, no code path executes it, and no code path can
// ever remove it.
//
// Corks are DROPPED rather than migrated into the authority queue. They were
// scheduled under validator consent that no longer carries meaning after the
// authorization model changes, and re-scheduling under the authority key costs
// one transaction. Migrating them would silently convert a supermajority
// decision into a unilateral one.
//
// Raw store prefixes are iterated rather than the typed keeper helpers,
// because those helpers are deleted in the same release.
//
// Returns the number of keys removed.
func DrainLegacyCorkQueues(ctx sdk.Context, corkStoreKey, axelarcorkStoreKey storetypes.StoreKey) int {
	drained := 0

	corkStore := ctx.KVStore(corkStoreKey)
	drained += deletePrefix(corkStore, []byte{corktypes.ScheduledCorkKeyPrefix})
	drained += deletePrefix(corkStore, []byte{corktypes.ValidatorCorkCountKey})

	axelarcorkStore := ctx.KVStore(axelarcorkStoreKey)
	drained += deletePrefix(axelarcorkStore, []byte{axelarcorktypes.ScheduledCorkKeyPrefix})
	drained += deletePrefix(axelarcorkStore, []byte{axelarcorktypes.ValidatorAxelarCorkCountKey})

	return drained
}

// deletePrefix removes every key under prefix and returns how many it removed.
//
// Keys are collected before any deletion: mutating a store through a live
// iterator is undefined behaviour in the SDK. The key bytes are copied because
// the iterator may reuse its buffer between iterations.
func deletePrefix(store storetypes.KVStore, prefix []byte) int {
	var keys [][]byte

	iter := sdk.KVStorePrefixIterator(store, prefix)
	for ; iter.Valid(); iter.Next() {
		key := make([]byte, len(iter.Key()))
		copy(key, iter.Key())
		keys = append(keys, key)
	}
	iter.Close()

	for _, k := range keys {
		store.Delete(k)
	}

	return len(keys)
}
