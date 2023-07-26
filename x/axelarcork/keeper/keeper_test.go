package keeper

import (
	"testing"

	"github.com/peggyjv/sommelier/v6/x/axelarcork/tests/mocks"

	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/testutil"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/bech32"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	paramskeeper "github.com/cosmos/cosmos-sdk/x/params/keeper"
	"github.com/ethereum/go-ethereum/common"
	"github.com/golang/mock/gomock"
	moduletestutil "github.com/peggyjv/sommelier/v6/testutil"
	"github.com/peggyjv/sommelier/v6/x/axelarcork/types"
	"github.com/stretchr/testify/suite"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
	tmtime "github.com/tendermint/tendermint/types/time"
)

var (
	sampleCellarHex  = "0xc0ffee254729296a45a3885639AC7E10F9d54979"
	sampleCellarAddr = common.HexToAddress(sampleCellarHex)
)

type KeeperTestSuite struct {
	suite.Suite

	ctx                sdk.Context
	corkKeeper         Keeper
	stakingKeeper      *mocks.MockStakingKeeper
	transferKeeper     *mocks.MockTransferKeeper
	distributionKeeper *mocks.MockDistributionKeeper
	ics4wrapper        *mocks.MockICS4Wrapper
	gravityKeeper      *mocks.MockGravityKeeper

	validator *mocks.MockValidatorI

	queryClient types.QueryClient

	encCfg moduletestutil.TestEncodingConfig
}

func (suite *KeeperTestSuite) SetupTest() {
	key := sdk.NewKVStoreKey(types.StoreKey)
	tkey := sdk.NewTransientStoreKey("transient_test")
	testCtx := testutil.DefaultContext(key, tkey)
	ctx := testCtx.WithBlockHeader(tmproto.Header{Height: 5, Time: tmtime.Now()})
	encCfg := moduletestutil.MakeTestEncodingConfig()

	// gomock initializations
	ctrl := gomock.NewController(suite.T())
	defer ctrl.Finish()

	suite.stakingKeeper = mocks.NewMockStakingKeeper(ctrl)
	suite.transferKeeper = mocks.NewMockTransferKeeper(ctrl)
	suite.distributionKeeper = mocks.NewMockDistributionKeeper(ctrl)
	suite.ics4wrapper = mocks.NewMockICS4Wrapper(ctrl)
	suite.validator = mocks.NewMockValidatorI(ctrl)
	suite.ctx = ctx

	params := paramskeeper.NewKeeper(
		encCfg.Codec,
		codec.NewLegacyAmino(),
		key,
		tkey,
	)

	params.Subspace(types.ModuleName)
	subSpace, found := params.GetSubspace(types.ModuleName)
	suite.Assertions.True(found)

	suite.corkKeeper = NewKeeper(
		encCfg.Codec,
		key,
		subSpace,
		suite.stakingKeeper,
		suite.transferKeeper,
		suite.distributionKeeper,
		suite.ics4wrapper,
		suite.gravityKeeper,
	)

	//types.RegisterInterfaces(encCfg.InterfaceRegistry)

	queryHelper := baseapp.NewQueryServerTestHelper(ctx, encCfg.InterfaceRegistry)
	types.RegisterQueryServer(queryHelper, suite.corkKeeper)
	queryClient := types.NewQueryClient(queryHelper)

	suite.queryClient = queryClient
	suite.encCfg = encCfg
}

func TestKeeperTestSuite(t *testing.T) {
	suite.Run(t, new(KeeperTestSuite))
}

func (suite *KeeperTestSuite) TestSetGetCellarIDsHappyPath() {
	ctx, corkKeeper := suite.ctx, suite.corkKeeper
	require := suite.Require()

	cellarIDSet := types.CellarIDSet{
		Ids: []string{sampleCellarHex},
	}
	expected := []common.Address{}
	for _, id := range cellarIDSet.Ids {
		expected = append(expected, common.HexToAddress(id))
	}
	corkKeeper.SetCellarIDs(ctx, TestEVMChainID, cellarIDSet)
	actual := corkKeeper.GetCellarIDs(ctx, TestEVMChainID)

	require.Equal(expected, actual)
	require.True(corkKeeper.HasCellarID(ctx, TestEVMChainID, sampleCellarAddr))
}

func (suite *KeeperTestSuite) TestSetGetDeleteScheduledCork() {
	ctx, corkKeeper := suite.ctx, suite.corkKeeper
	require := suite.Require()

	corkKeeper.SetChainConfiguration(ctx, TestEVMChainID, types.ChainConfiguration{
		Name:         "testevm",
		Id:           TestEVMChainID,
		ProxyAddress: "0x123",
	})

	testHeight := uint64(123)
	val := []byte("testaddress")
	expectedCork := types.AxelarCork{
		EncodedContractCall:   []byte("testcall"),
		TargetContractAddress: sampleCellarHex,
	}
	expectedID := expectedCork.IDHash(TestEVMChainID, testHeight)
	actualID := corkKeeper.SetScheduledAxelarCork(ctx, TestEVMChainID, testHeight, val, expectedCork)
	require.Equal(expectedID, actualID)
	actualCork, found := corkKeeper.GetScheduledAxelarCork(ctx, TestEVMChainID, testHeight, actualID, val, sampleCellarAddr)
	require.True(found)
	require.Equal(expectedCork, actualCork)

	actualCorks := corkKeeper.GetScheduledAxelarCorks(ctx, TestEVMChainID)
	require.Equal(&expectedCork, actualCorks[0].Cork)

	actualCorks = corkKeeper.GetScheduledAxelarCorksByID(ctx, TestEVMChainID, actualID)
	require.Len(actualCorks, 1)
	require.Equal(&expectedCork, actualCorks[0].Cork)
	require.Equal(expectedID, actualCorks[0].Id)

	actualHeights := corkKeeper.GetScheduledBlockHeights(ctx, TestEVMChainID)
	require.Equal(actualCorks[0].BlockHeight, actualHeights[0])

	actualCorks = corkKeeper.GetScheduledAxelarCorksByBlockHeight(ctx, TestEVMChainID, testHeight)
	require.Equal(&expectedCork, actualCorks[0].Cork)
	require.Equal(testHeight, actualCorks[0].BlockHeight)
	require.Equal(expectedID, actualCorks[0].Id)

	corkKeeper.DeleteScheduledAxelarCork(ctx, TestEVMChainID, testHeight, expectedID, sdk.ValAddress(val), sampleCellarAddr)
	require.Empty(corkKeeper.GetScheduledAxelarCorks(ctx, TestEVMChainID))
}

func (suite *KeeperTestSuite) TestGetWinningVotes() {
	ctx, corkKeeper := suite.ctx, suite.corkKeeper
	require := suite.Require()

	corkKeeper.SetChainConfiguration(ctx, TestEVMChainID, types.ChainConfiguration{
		Name:         "testevm",
		Id:           TestEVMChainID,
		ProxyAddress: "0x123",
	})

	testHeight := uint64(ctx.BlockHeight())
	cork := types.AxelarCork{
		EncodedContractCall:   []byte("testcall"),
		TargetContractAddress: sampleCellarHex,
	}
	_, bytes, err := bech32.DecodeAndConvert("somm1fcl08ymkl70dhyg3vmx4hjsqvxym7dawnp0zfp")
	require.NoError(err)
	require.Equal(20, len(bytes))
	corkKeeper.SetScheduledAxelarCork(ctx, TestEVMChainID, testHeight, bytes, cork)

	suite.stakingKeeper.EXPECT().GetLastTotalPower(ctx).Return(sdk.NewInt(100))
	suite.stakingKeeper.EXPECT().Validator(ctx, gomock.Any()).Return(suite.validator)
	suite.validator.EXPECT().GetConsensusPower(gomock.Any()).Return(int64(100))
	suite.stakingKeeper.EXPECT().PowerReduction(ctx).Return(sdk.OneInt())

	winningScheduledVotes := corkKeeper.GetApprovedScheduledAxelarCorks(ctx, TestEVMChainID)
	results := corkKeeper.GetAxelarCorkResults(ctx, TestEVMChainID)
	require.Len(winningScheduledVotes, 1)
	require.Equal(cork, winningScheduledVotes[0])
	require.Equal(&cork, results[0].Cork)
	require.True(results[0].Approved)
	require.Equal("100.000000000000000000", results[0].ApprovalPercentage)

	// scheduled cork should be deleted at the scheduled height
	require.Empty(corkKeeper.GetScheduledAxelarCorksByBlockHeight(ctx, TestEVMChainID, testHeight))
}

func (suite *KeeperTestSuite) TestCorkResults() {
	ctx, corkKeeper := suite.ctx, suite.corkKeeper
	require := suite.Require()

	require.Empty(corkKeeper.GetAxelarCorkResults(ctx, TestEVMChainID))

	testHeight := uint64(ctx.BlockHeight())
	cork := types.AxelarCork{
		EncodedContractCall:   []byte("testcall"),
		TargetContractAddress: sampleCellarHex,
	}
	id := cork.IDHash(TestEVMChainID, testHeight)
	result := types.AxelarCorkResult{
		Cork:               &cork,
		BlockHeight:        testHeight,
		Approved:           true,
		ApprovalPercentage: "100.00",
	}
	corkKeeper.SetAxelarCorkResult(ctx, TestEVMChainID, id, result)
	actualResult, found := corkKeeper.GetAxelarCorkResult(ctx, TestEVMChainID, id)
	require.True(found)
	require.Equal(result, actualResult)

	results := corkKeeper.GetAxelarCorkResults(ctx, TestEVMChainID)
	require.Equal(&actualResult, results[0])

	corkKeeper.DeleteAxelarCorkResult(ctx, TestEVMChainID, id)
	require.Empty(corkKeeper.GetAxelarCorkResults(ctx, TestEVMChainID))
}

func (suite *KeeperTestSuite) TestParamSet() {
	ctx, corkKeeper := suite.ctx, suite.corkKeeper
	require := suite.Require()

	require.Panics(func() { corkKeeper.GetParamSet(ctx) })

	testGMPAccount := authtypes.NewModuleAddress("test-gmp-account")

	params := types.Params{
		IbcChannel:      "test-ibc-channel",
		IbcPort:         "test-ibc-port",
		GmpAccount:      testGMPAccount.String(),
		ExecutorAccount: "test-executor-account",
		TimeoutDuration: 10,
	}
	corkKeeper.SetParams(ctx, params)
	require.Equal(params, corkKeeper.GetParamSet(ctx))
}
