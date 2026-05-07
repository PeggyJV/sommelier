package keeper

import (
	"testing"

	tmdb "github.com/cometbft/cometbft-db"
	"github.com/cometbft/cometbft/libs/log"
	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/store"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	paramskeeper "github.com/cosmos/cosmos-sdk/x/params/keeper"
	paramstypes "github.com/cosmos/cosmos-sdk/x/params/types"
	"github.com/stretchr/testify/require"

	"github.com/peggyjv/sommelier/v9/x/poa/types"
)

// NewTestKeeper builds a PoA keeper backed by an in-memory store. Staking and
// slashing keepers are set to nil; tests that exercise paths touching them
// must use NewTestKeeperWithMocks instead.
func NewTestKeeper(t *testing.T) (Keeper, sdk.Context) {
	t.Helper()
	db := tmdb.NewMemDB()
	cms := store.NewCommitMultiStore(db)

	storeKey := sdk.NewKVStoreKey(types.StoreKey)
	tStoreKey := sdk.NewTransientStoreKey("transient_" + types.StoreKey)
	memStoreKey := storetypes.NewMemoryStoreKey("mem_" + types.StoreKey)

	cms.MountStoreWithDB(storeKey, storetypes.StoreTypeIAVL, db)
	cms.MountStoreWithDB(tStoreKey, storetypes.StoreTypeTransient, nil)
	cms.MountStoreWithDB(memStoreKey, storetypes.StoreTypeMemory, nil)

	// Params subspace setup.
	paramsStoreKey := sdk.NewKVStoreKey(paramstypes.StoreKey)
	paramsTStoreKey := sdk.NewTransientStoreKey(paramstypes.TStoreKey)
	cms.MountStoreWithDB(paramsStoreKey, storetypes.StoreTypeIAVL, db)
	cms.MountStoreWithDB(paramsTStoreKey, storetypes.StoreTypeTransient, nil)

	require.NoError(t, cms.LoadLatestVersion())

	cdc := codec.NewProtoCodec(codectypes.NewInterfaceRegistry())
	pk := paramskeeper.NewKeeper(cdc, codec.NewLegacyAmino(), paramsStoreKey, paramsTStoreKey)
	subspace := pk.Subspace(types.ModuleName).WithKeyTable(types.ParamKeyTable())

	k := NewKeeper(cdc, storeKey, subspace, nil, nil, "cosmos1zkmrn5j2t9k3s7n9z2c5w0n9c0w8e8mq8w8mq8")
	ctx := sdk.NewContext(cms, tmproto.Header{}, false, log.NewNopLogger())

	// Initialise default params so reads do not panic.
	k.SetParams(ctx, types.DefaultParams())

	return k, ctx
}
