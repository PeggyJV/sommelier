package keeper

import (
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"testing"

	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/testutil"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/bech32"
	paramskeeper "github.com/cosmos/cosmos-sdk/x/params/keeper"
	"github.com/ethereum/go-ethereum/common"
	"github.com/golang/mock/gomock"
	moduletestutil "github.com/peggyjv/sommelier/v6/testutil"
	corktestutil "github.com/peggyjv/sommelier/v6/x/axelarcork/tests"
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
	stakingKeeper      *corktestutil.MockStakingKeeper
	transferKeeper     *corktestutil.MockTransferKeeper
	distributionKeeper *corktestutil.MockDistributionKeeper
	ics4wrapper        *corktestutil.MockICS4Wrapper

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

	suite.stakingKeeper = corktestutil.NewMockStakingKeeper(ctrl)
	suite.transferKeeper = corktestutil.NewMockTransferKeeper(ctrl)
	suite.distributionKeeper = corktestutil.NewMockDistributionKeeper(ctrl)
	suite.ics4wrapper = corktestutil.NewMockICS4Wrapper(ctrl)
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

	corkKeeper.SetChainConfigurationByID(ctx, TestEVMChainID, types.ChainConfiguration{
		Name:          "testevm",
		Id:            TestEVMChainID,
		VoteThreshold: sdk.NewDec(0),
		ProxyAddress:  "0x123",
	})

	testHeight := uint64(123)
	val := []byte("testaddress")
	expectedCork := types.Cork{
		EncodedContractCall:   []byte("testcall"),
		TargetContractAddress: sampleCellarHex,
	}
	expectedID := expectedCork.IDHash(testHeight)
	actualID := corkKeeper.SetScheduledCork(ctx, TestEVMChainID, testHeight, val, expectedCork)
	require.Equal(expectedID, actualID)
	actualCork, found := corkKeeper.GetScheduledCork(ctx, TestEVMChainID, testHeight, actualID, val, sampleCellarAddr)
	require.True(found)
	require.Equal(expectedCork, actualCork)

	actualCorks := corkKeeper.GetScheduledCorks(ctx, TestEVMChainID)
	require.Equal(&expectedCork, actualCorks[0].Cork)

	actualCorks = corkKeeper.GetScheduledCorksByID(ctx, TestEVMChainID, actualID)
	require.Len(actualCorks, 1)
	require.Equal(&expectedCork, actualCorks[0].Cork)
	require.Equal(expectedID, actualCorks[0].Id)

	actualHeights := corkKeeper.GetScheduledBlockHeights(ctx, TestEVMChainID)
	require.Equal(actualCorks[0].BlockHeight, actualHeights[0])

	actualCorks = corkKeeper.GetScheduledCorksByBlockHeight(ctx, TestEVMChainID, testHeight)
	require.Equal(&expectedCork, actualCorks[0].Cork)
	require.Equal(testHeight, actualCorks[0].BlockHeight)
	require.Equal(expectedID, actualCorks[0].Id)

	corkKeeper.DeleteScheduledCork(ctx, TestEVMChainID, testHeight, expectedID, sdk.ValAddress(val), sampleCellarAddr)
	require.Empty(corkKeeper.GetScheduledCorks(ctx, TestEVMChainID))
}

func (suite *KeeperTestSuite) TestGetWinningVotes() {
	ctx, corkKeeper := suite.ctx, suite.corkKeeper
	require := suite.Require()

	corkKeeper.SetChainConfigurationByID(ctx, TestEVMChainID, types.ChainConfiguration{
		Name:          "testevm",
		Id:            TestEVMChainID,
		VoteThreshold: sdk.NewDec(0),
		ProxyAddress:  "0x123",
	})

	testHeight := uint64(ctx.BlockHeight())
	cork := types.Cork{
		EncodedContractCall:   []byte("testcall"),
		TargetContractAddress: sampleCellarHex,
	}
	_, bytes, err := bech32.DecodeAndConvert("somm1fcl08ymkl70dhyg3vmx4hjsqvxym7dawnp0zfp")
	require.NoError(err)
	require.Equal(20, len(bytes))
	corkKeeper.SetScheduledCork(ctx, TestEVMChainID, testHeight, bytes, cork)

	suite.stakingKeeper.EXPECT().GetLastTotalPower(ctx).Return(sdk.NewInt(100))
	//suite.stakingKeeper.EXPECT().Validator(ctx, gomock.Any()).Return(suite.validator)
	//suite.validator.EXPECT().GetConsensusPower(gomock.Any()).Return(int64(100))
	suite.stakingKeeper.EXPECT().PowerReduction(ctx).Return(sdk.OneInt())

	winningScheduledVotes := corkKeeper.GetApprovedScheduledCorks(ctx, TestEVMChainID, testHeight, sdk.ZeroDec())
	results := corkKeeper.GetCorkResults(ctx, TestEVMChainID)
	require.Len(winningScheduledVotes, 1)
	require.Equal(cork, winningScheduledVotes[0])
	require.Equal(&cork, results[0].Cork)
	require.True(results[0].Approved)
	require.Equal("100.000000000000000000", results[0].ApprovalPercentage)

	// scheduled cork should be deleted at the scheduled height
	require.Empty(corkKeeper.GetScheduledCorksByBlockHeight(ctx, TestEVMChainID, testHeight))
}

func (suite *KeeperTestSuite) TestInvalidationNonce() {
	ctx, corkKeeper := suite.ctx, suite.corkKeeper
	require := suite.Require()

	require.Zero(corkKeeper.GetLatestInvalidationNonce(ctx, TestEVMChainID))

	corkKeeper.SetLatestInvalidationNonce(ctx, TestEVMChainID, uint64(5))
	require.Equal(uint64(5), corkKeeper.GetLatestInvalidationNonce(ctx, TestEVMChainID))

	corkKeeper.IncrementInvalidationNonce(ctx, TestEVMChainID)
	require.Equal(uint64(6), corkKeeper.GetLatestInvalidationNonce(ctx, TestEVMChainID))
}

func (suite *KeeperTestSuite) TestCorkResults() {
	ctx, corkKeeper := suite.ctx, suite.corkKeeper
	require := suite.Require()

	require.Empty(corkKeeper.GetCorkResults(ctx, TestEVMChainID))

	testHeight := uint64(ctx.BlockHeight())
	cork := types.Cork{
		EncodedContractCall:   []byte("testcall"),
		TargetContractAddress: sampleCellarHex,
	}
	id := cork.IDHash(testHeight)
	result := types.CorkResult{
		Cork:               &cork,
		BlockHeight:        testHeight,
		Approved:           true,
		ApprovalPercentage: "100.00",
	}
	corkKeeper.SetCorkResult(ctx, TestEVMChainID, id, result)
	actualResult, found := corkKeeper.GetCorkResult(ctx, TestEVMChainID, id)
	require.True(found)
	require.Equal(result, actualResult)

	results := corkKeeper.GetCorkResults(ctx, TestEVMChainID)
	require.Equal(&actualResult, results[0])

	corkKeeper.DeleteCorkResult(ctx, TestEVMChainID, id)
	require.Empty(corkKeeper.GetCorkResults(ctx, TestEVMChainID))
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
	corkKeeper.setParams(ctx, params)
	require.Equal(params, corkKeeper.GetParamSet(ctx))
}
