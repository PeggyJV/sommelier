package keeper

import (
	"fmt"
	"testing"

	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/store"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	paramskeeper "github.com/cosmos/cosmos-sdk/x/params/keeper"

	"github.com/golang/mock/gomock"

	"github.com/peggyjv/sommelier/v4/x/cork/types"
	"github.com/stretchr/testify/require"

	"github.com/tendermint/tendermint/libs/log"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
	tmdb "github.com/tendermint/tm-db"
)

type mockKeepers struct {
	mockStakingKeeper *MockStakingKeeper
	mockGravityKeeper *MockGravityKeeper
}

func setupCorkKeeper(t *testing.T) (Keeper, mockKeepers, sdk.Context) {
	db := tmdb.NewMemDB()
	commitMultiStore := store.NewCommitMultiStore(db)

	// Mount the KV store with the x/cork store key
	storeKey := sdk.NewKVStoreKey(types.StoreKey)
	commitMultiStore.MountStoreWithDB(storeKey, sdk.StoreTypeIAVL, db)

	// Mount Transient store
	transientStoreKey := sdk.NewTransientStoreKey("transient" + types.StoreKey)
	commitMultiStore.MountStoreWithDB(transientStoreKey, sdk.StoreTypeTransient, nil)

	// Mount Memory store
	memStoreKey := storetypes.NewMemoryStoreKey("mem" + types.StoreKey)
	commitMultiStore.MountStoreWithDB(memStoreKey, sdk.StoreTypeMemory, nil)

	require.NoError(t, commitMultiStore.LoadLatestVersion())
	protoCodec := codec.NewProtoCodec(codectypes.NewInterfaceRegistry())
	params := initParamsKeeper(
		protoCodec, codec.NewLegacyAmino(), storeKey, memStoreKey)

	subSpace, found := params.GetSubspace(types.ModuleName)
	require.True(t, found)

	ctrl := gomock.NewController(t)
	mockStakingKeeper := NewMockStakingKeeper(ctrl)
	mockGravityKeeper := NewMockGravityKeeper(ctrl)

	k := NewKeeper(
		protoCodec,
		storeKey,
		subSpace,
		mockStakingKeeper,
		mockGravityKeeper,
	)

	ctx := sdk.NewContext(commitMultiStore, tmproto.Header{}, false, log.NewNopLogger())

	return k, mockKeepers{
		mockStakingKeeper: mockStakingKeeper,
		mockGravityKeeper: mockGravityKeeper,
	}, ctx
}

func initParamsKeeper(
	appCodec codec.BinaryCodec, legacyAmino *codec.LegacyAmino,
	key sdk.StoreKey, tkey sdk.StoreKey,
) paramskeeper.Keeper {
	paramsKeeper := paramskeeper.NewKeeper(appCodec, legacyAmino, key, tkey)
	paramsKeeper.Subspace(types.ModuleName)

	return paramsKeeper
}

func TestSetupCorkKeepers(t *testing.T) {
	testCases := []struct {
		name string
		test func()
	}{{
		name: "happy path",
		test: func() {
			k, mocks, ctx := setupCorkKeeper(t)
			fmt.Println(k, mocks, ctx)
		},
	},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			tc.test()
		})
	}

}
