package v10

import (
	"testing"

	"github.com/cosmos/cosmos-sdk/store"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"

	tmdb "github.com/cometbft/cometbft-db"
	"github.com/cometbft/cometbft/libs/log"
	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"

	axelarcorktypes "github.com/peggyjv/sommelier/v10/x/axelarcork/types"
	corktypes "github.com/peggyjv/sommelier/v10/x/cork/types"
)

// drainTestStores mounts the two module stores this migration touches.
func drainTestStores(t *testing.T) (sdk.Context, storetypes.StoreKey, storetypes.StoreKey) {
	t.Helper()

	db := tmdb.NewMemDB()
	cms := store.NewCommitMultiStore(db)

	corkKey := sdk.NewKVStoreKey(corktypes.StoreKey)
	axelarKey := sdk.NewKVStoreKey(axelarcorktypes.StoreKey)
	cms.MountStoreWithDB(corkKey, storetypes.StoreTypeIAVL, db)
	cms.MountStoreWithDB(axelarKey, storetypes.StoreTypeIAVL, db)
	require.NoError(t, cms.LoadLatestVersion())

	ctx := sdk.NewContext(cms, tmproto.Header{Height: 100}, false, log.NewNopLogger())
	return ctx, corkKey, axelarKey
}

func countPrefix(ctx sdk.Context, key storetypes.StoreKey, prefix []byte) int {
	n := 0
	iter := sdk.KVStorePrefixIterator(ctx.KVStore(key), prefix)
	defer iter.Close()
	for ; iter.Valid(); iter.Next() {
		n++
	}
	return n
}

// seedLegacy writes entries under the retired prefixes, imitating the state a
// live chain carries into the upgrade.
func seedLegacy(ctx sdk.Context, corkKey, axelarKey storetypes.StoreKey) (corkCount, axelarCount int) {
	corkStore := ctx.KVStore(corkKey)
	axelarStore := ctx.KVStore(axelarKey)

	val := sdk.ValAddress([]byte("12345678901234567890"))
	contract := common.HexToAddress("0x1111111111111111111111111111111111111111")
	id := make([]byte, 32)
	for i := range id {
		id[i] = byte(i)
	}

	// three cork corks at differing heights
	for _, h := range []uint64{100, 101, 102} {
		corkStore.Set(corktypes.GetScheduledCorkKey(h, id, val, contract), []byte{0x01})
		corkCount++
	}
	corkStore.Set(corktypes.GetValidatorCorkCountKey(val), []byte{0x03})
	corkCount++

	// two axelar corks on each of two chains
	for _, chainID := range []uint64{42161, 10} {
		for _, h := range []uint64{100, 101} {
			axelarStore.Set(axelarcorktypes.GetScheduledAxelarCorkKey(chainID, h, id, val, contract), []byte{0x01})
			axelarCount++
		}
	}
	axelarStore.Set(axelarcorktypes.GetValidatorAxelarCorkCountKey(val), []byte{0x04})
	axelarCount++

	return corkCount, axelarCount
}

func TestDrainLegacyCorkQueuesRemovesEverything(t *testing.T) {
	ctx, corkKey, axelarKey := drainTestStores(t)
	corkSeeded, axelarSeeded := seedLegacy(ctx, corkKey, axelarKey)

	// An authority cork must NOT be touched: the drain drops legacy state, it
	// does not migrate it, and it must not disturb the new queue.
	authorityKey := corktypes.GetAuthorityCorkKey(100, make([]byte, 32), common.HexToAddress("0x2222222222222222222222222222222222222222"))
	ctx.KVStore(corkKey).Set(authorityKey, []byte{0x09})

	require.Equal(t, corkSeeded+axelarSeeded, 4+5)

	drained := DrainLegacyCorkQueues(ctx, corkKey, axelarKey)
	require.Equal(t, corkSeeded+axelarSeeded, drained,
		"every seeded legacy key must be counted as drained")

	require.Zero(t, countPrefix(ctx, corkKey, []byte{corktypes.ScheduledCorkKeyPrefix}),
		"cork scheduled queue must be empty")
	require.Zero(t, countPrefix(ctx, corkKey, []byte{corktypes.ValidatorCorkCountKey}),
		"cork validator counters must be cleared")
	require.Zero(t, countPrefix(ctx, axelarKey, []byte{axelarcorktypes.ScheduledCorkKeyPrefix}),
		"axelarcork scheduled queue must be empty")
	require.Zero(t, countPrefix(ctx, axelarKey, []byte{axelarcorktypes.ValidatorAxelarCorkCountKey}),
		"axelarcork validator counters must be cleared")

	require.Equal(t, []byte{0x09}, ctx.KVStore(corkKey).Get(authorityKey),
		"the authority queue must be left untouched")
}

func TestDrainLegacyCorkQueuesIsIdempotent(t *testing.T) {
	ctx, corkKey, axelarKey := drainTestStores(t)
	seedLegacy(ctx, corkKey, axelarKey)

	first := DrainLegacyCorkQueues(ctx, corkKey, axelarKey)
	require.NotZero(t, first)

	second := DrainLegacyCorkQueues(ctx, corkKey, axelarKey)
	require.Zero(t, second, "a second run on drained state must remove nothing and not panic")
}

func TestDrainLegacyCorkQueuesOnEmptyStoreIsNoop(t *testing.T) {
	ctx, corkKey, axelarKey := drainTestStores(t)
	require.Zero(t, DrainLegacyCorkQueues(ctx, corkKey, axelarKey))
}
