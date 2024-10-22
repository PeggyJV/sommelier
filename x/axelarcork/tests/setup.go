package tests

import (
	"testing"

	tmdb "github.com/cometbft/cometbft-db"
	"github.com/cometbft/cometbft/libs/log"
	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"
	"github.com/cosmos/cosmos-sdk/codec"
	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/store"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	paramskeeper "github.com/cosmos/cosmos-sdk/x/params/keeper"
	paramstypes "github.com/cosmos/cosmos-sdk/x/params/types"
	porttypes "github.com/cosmos/ibc-go/v7/modules/core/05-port/types"
	"github.com/golang/mock/gomock"
	"github.com/peggyjv/sommelier/v8/x/axelarcork"
	"github.com/peggyjv/sommelier/v8/x/axelarcork/keeper"
	"github.com/peggyjv/sommelier/v8/x/axelarcork/tests/mocks"
	"github.com/peggyjv/sommelier/v8/x/axelarcork/types"
	"github.com/stretchr/testify/require"
)

var TestGMPAccount = authtypes.NewModuleAddress("test-gmp-account")

// DefaultParams returns default oracle parameters
func TestParams() types.Params {

	return types.Params{
		IbcChannel:      "channel-2",
		IbcPort:         "axelar",
		GmpAccount:      TestGMPAccount.String(),
		ExecutorAccount: "test-executor-account",
		TimeoutDuration: 100,
	}
}

func NewTestSetup(t *testing.T, ctl *gomock.Controller) *Setup {
	t.Helper()
	initializer := newInitializer()

	accountKeeperMock := mocks.NewMockAccountKeeper(ctl)
	bankKeeperMock := mocks.NewMockBankKeeper(ctl)
	transferKeeperMock := mocks.NewMockTransferKeeper(ctl)
	stakingKeeper := mocks.NewMockStakingKeeper(ctl)
	distributionKeeperMock := mocks.NewMockDistributionKeeper(ctl)
	ibcModuleMock := mocks.NewMockIBCModule(ctl)
	ics4WrapperMock := mocks.NewMockICS4Wrapper(ctl)
	gravityKeeper := mocks.NewMockGravityKeeper(ctl)
	pubsubKeeper := mocks.NewMockPubsubKeeper(ctl)

	paramsKeeper := initializer.paramsKeeper()
	acKeeper := initializer.axelarcorkKeeper(
		paramsKeeper, accountKeeperMock, bankKeeperMock, stakingKeeper, transferKeeperMock, distributionKeeperMock,
		ics4WrapperMock, gravityKeeper, pubsubKeeper)

	require.NoError(t, initializer.StateStore.LoadLatestVersion())

	acKeeper.SetParams(initializer.Ctx, TestParams())
	//acKeeper.SetChainConfiguration()

	return &Setup{
		Initializer: initializer,

		Keepers: &testKeepers{
			ParamsKeeper:     &paramsKeeper,
			AxelarCorkKeeper: acKeeper,
		},

		Mocks: &testMocks{
			AccountKeeperMock:      accountKeeperMock,
			BankKeeperMock:         bankKeeperMock,
			TransferKeeperMock:     transferKeeperMock,
			DistributionKeeperMock: distributionKeeperMock,
			IBCModuleMock:          ibcModuleMock,
			ICS4WrapperMock:        ics4WrapperMock,
		},

		AxelarCorkMiddleware: initializer.axelarMiddleware(ibcModuleMock, acKeeper),
	}
}

type Setup struct {
	Initializer initializer

	Keepers *testKeepers
	Mocks   *testMocks

	AxelarCorkMiddleware axelarcork.IBCMiddleware
}

type testKeepers struct {
	ParamsKeeper     *paramskeeper.Keeper
	AxelarCorkKeeper *keeper.Keeper
}

type testMocks struct {
	AccountKeeperMock      *mocks.MockAccountKeeper
	BankKeeperMock         *mocks.MockBankKeeper
	TransferKeeperMock     *mocks.MockTransferKeeper
	DistributionKeeperMock *mocks.MockDistributionKeeper
	IBCModuleMock          *mocks.MockIBCModule
	ICS4WrapperMock        *mocks.MockICS4Wrapper
}

type initializer struct {
	DB         *tmdb.MemDB
	StateStore store.CommitMultiStore
	Ctx        sdk.Context
	Marshaler  codec.Codec
	Amino      *codec.LegacyAmino
}

// Create an initializer with in memory database and default codecs
func newInitializer() initializer {
	logger := log.TestingLogger()
	logger.Debug("initializing test setup")

	db := tmdb.NewMemDB()
	stateStore := store.NewCommitMultiStore(db)

	ctx := sdk.NewContext(stateStore, tmproto.Header{}, false, logger)
	interfaceRegistry := cdctypes.NewInterfaceRegistry()
	amino := codec.NewLegacyAmino()
	marshaler := codec.NewProtoCodec(interfaceRegistry)

	return initializer{
		DB:         db,
		StateStore: stateStore,
		Ctx:        ctx,
		Marshaler:  marshaler,
		Amino:      amino,
	}
}

func (i initializer) paramsKeeper() paramskeeper.Keeper {
	storeKey := sdk.NewKVStoreKey(paramstypes.StoreKey)
	transientStoreKey := sdk.NewTransientStoreKey(paramstypes.TStoreKey)
	i.StateStore.MountStoreWithDB(storeKey, storetypes.StoreTypeIAVL, i.DB)
	i.StateStore.MountStoreWithDB(transientStoreKey, storetypes.StoreTypeTransient, i.DB)

	paramsKeeper := paramskeeper.NewKeeper(i.Marshaler, i.Amino, storeKey, transientStoreKey)

	return paramsKeeper
}

func (i initializer) axelarcorkKeeper(
	paramsKeeper paramskeeper.Keeper,
	accountKeeper types.AccountKeeper,
	bankKeeper types.BankKeeper,
	stakingKeeper types.StakingKeeper,
	transferKeeper types.TransferKeeper,
	distributionKeeper types.DistributionKeeper,
	ics4Wrapper porttypes.ICS4Wrapper,
	gravityKeeper types.GravityKeeper,
	pubsubKeeper types.PubsubKeeper,
) *keeper.Keeper {
	storeKey := sdk.NewKVStoreKey(types.StoreKey)
	i.StateStore.MountStoreWithDB(storeKey, storetypes.StoreTypeIAVL, i.DB)

	subspace := paramsKeeper.Subspace(types.ModuleName)
	routerKeeper := keeper.NewKeeper(
		i.Marshaler,
		storeKey,
		subspace,
		accountKeeper,
		bankKeeper,
		stakingKeeper,
		transferKeeper,
		distributionKeeper,
		ics4Wrapper,
		gravityKeeper,
		pubsubKeeper,
	)

	return &routerKeeper
}

func (i initializer) axelarMiddleware(app porttypes.IBCModule, k *keeper.Keeper) axelarcork.IBCMiddleware {
	return axelarcork.NewIBCMiddleware(k, app)
}
