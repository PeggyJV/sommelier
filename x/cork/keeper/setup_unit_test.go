package keeper

import (
	"reflect"
	"testing"

	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/store"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	paramskeeper "github.com/cosmos/cosmos-sdk/x/params/keeper"

	"github.com/golang/mock/gomock"

	"github.com/peggyjv/sommelier/v9/x/cork/mock"
	"github.com/peggyjv/sommelier/v9/x/cork/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	tmdb "github.com/cometbft/cometbft-db"
	"github.com/cometbft/cometbft/libs/log"
	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
)

type mocksForCork struct {
	mockStakingKeeper *mock.MockStakingKeeper
	mockGravityKeeper *mock.MockGravityKeeper
	mockPubsubKeeper  *mock.MockPubsubKeeper
	mockValidator     *mock.MockValidatorI
}

func setupCorkKeeper(t *testing.T) (
	Keeper, sdk.Context, mocksForCork, *gomock.Controller,
) {
	db := tmdb.NewMemDB()
	commitMultiStore := store.NewCommitMultiStore(db)

	// Mount the KV store with the x/cork store key
	storeKey := sdk.NewKVStoreKey(types.StoreKey)
	commitMultiStore.MountStoreWithDB(storeKey, storetypes.StoreTypeIAVL, db)

	// Mount Transient store
	transientStoreKey := sdk.NewTransientStoreKey("transient" + types.StoreKey)
	commitMultiStore.MountStoreWithDB(transientStoreKey, storetypes.StoreTypeTransient, nil)

	// Mount Memory store
	memStoreKey := storetypes.NewMemoryStoreKey("mem" + types.StoreKey)
	commitMultiStore.MountStoreWithDB(memStoreKey, storetypes.StoreTypeMemory, nil)

	require.NoError(t, commitMultiStore.LoadLatestVersion())
	protoCodec := codec.NewProtoCodec(codectypes.NewInterfaceRegistry())
	params := initParamsKeeper(
		protoCodec, codec.NewLegacyAmino(), storeKey, memStoreKey)

	subSpace, found := params.GetSubspace(types.ModuleName)
	require.True(t, found)

	ctrl := gomock.NewController(t)
	mockStakingKeeper := mock.NewMockStakingKeeper(ctrl)
	mockGravityKeeper := mock.NewMockGravityKeeper(ctrl)
	mockPubsubKeeper := mock.NewMockPubsubKeeper(ctrl)

	k := NewKeeper(
		protoCodec,
		storeKey,
		subSpace,
		mockStakingKeeper,
		mockGravityKeeper,
		mockPubsubKeeper,
	)

	ctx := sdk.NewContext(commitMultiStore, tmproto.Header{}, false, log.NewNopLogger())

	return k, ctx, mocksForCork{
		mockStakingKeeper: mockStakingKeeper,
		mockGravityKeeper: mockGravityKeeper,
		mockPubsubKeeper:  mockPubsubKeeper,
		mockValidator:     mock.NewMockValidatorI(ctrl),
	}, ctrl
}

func initParamsKeeper(
	appCodec codec.BinaryCodec, legacyAmino *codec.LegacyAmino,
	key storetypes.StoreKey, tkey storetypes.StoreKey,
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
		name: "make sure that mocks implement expected keepers interfaces",
		test: func() {
			k, ctx, mocks, ctrl := setupCorkKeeper(t)
			require.PanicsWithError(t, "UnmarshalJSON cannot decode empty bytes",
				func() {
					params := k.GetParamSet(ctx)
					require.NoError(t, params.ValidateBasic())
				},
			)

			for _, keeperPair := range []struct {
				expected interface{}
				mock     interface{}
			}{
				{
					expected: (*types.StakingKeeper)(nil),
					mock:     mocks.mockStakingKeeper,
				},
				{
					expected: (*types.GravityKeeper)(nil),
					mock:     mocks.mockGravityKeeper,
				},
				{
					expected: (*stakingtypes.ValidatorI)(nil),
					mock:     mocks.mockValidator,
				},
			} {
				_interface := reflect.TypeOf(keeperPair.expected).Elem()
				isImplementingExpectedMethods := reflect.
					TypeOf(keeperPair.mock).Implements(_interface)
				assert.True(t, isImplementingExpectedMethods)
			}

			defer ctrl.Finish()
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
