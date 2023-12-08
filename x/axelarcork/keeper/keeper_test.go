package keeper

import (
	"encoding/hex"
	"testing"

	"github.com/peggyjv/sommelier/v7/x/axelarcork/tests/mocks"

	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/testutil"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/bech32"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	paramskeeper "github.com/cosmos/cosmos-sdk/x/params/keeper"
	"github.com/ethereum/go-ethereum/common"
	"github.com/golang/mock/gomock"
	moduletestutil "github.com/peggyjv/sommelier/v7/testutil"
	"github.com/peggyjv/sommelier/v7/x/axelarcork/types"
	"github.com/stretchr/testify/suite"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
	tmtime "github.com/tendermint/tendermint/types/time"
)

const (
	zeroAddressHex = "0x0000000000000000000000000000000000000000"
	oneAddressHex  = "0x1111111111111111111111111111111111111111"
)

var (
	sampleCellarHex  = "0xc0ffee254729296a45a3885639AC7E10F9d54979"
	sampleCellarAddr = common.HexToAddress(sampleCellarHex)
)

type KeeperTestSuite struct {
	suite.Suite

	ctx                sdk.Context
	axelarcorkKeeper   Keeper
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

	suite.axelarcorkKeeper = NewKeeper(
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
	types.RegisterQueryServer(queryHelper, suite.axelarcorkKeeper)
	queryClient := types.NewQueryClient(queryHelper)

	suite.queryClient = queryClient
	suite.encCfg = encCfg
}

func TestKeeperTestSuite(t *testing.T) {
	suite.Run(t, new(KeeperTestSuite))
}

func (suite *KeeperTestSuite) TestSetGetCellarIDsHappyPath() {
	ctx, axelarcorkKeeper := suite.ctx, suite.axelarcorkKeeper
	require := suite.Require()

	cellarIDSet := types.CellarIDSet{
		Ids: []string{sampleCellarHex},
	}
	expected := []common.Address{}
	for _, id := range cellarIDSet.Ids {
		expected = append(expected, common.HexToAddress(id))
	}
	axelarcorkKeeper.SetCellarIDs(ctx, TestEVMChainID, cellarIDSet)
	actual := axelarcorkKeeper.GetCellarIDs(ctx, TestEVMChainID)

	require.Equal(expected, actual)
	require.True(axelarcorkKeeper.HasCellarID(ctx, TestEVMChainID, sampleCellarAddr))
}

func (suite *KeeperTestSuite) TestSetGetDeleteScheduledCork() {
	ctx, axelarcorkKeeper := suite.ctx, suite.axelarcorkKeeper
	require := suite.Require()

	axelarcorkKeeper.SetChainConfiguration(ctx, TestEVMChainID, types.ChainConfiguration{
		Name:         "testevm",
		Id:           TestEVMChainID,
		ProxyAddress: "0x123",
	})

	testHeight := uint64(123)
	val := []byte("testaddress")
	deadline := uint64(10000000000)
	expectedCork := types.AxelarCork{
		EncodedContractCall:   []byte("testcall"),
		TargetContractAddress: sampleCellarHex,
		Deadline:              deadline,
	}
	expectedID := expectedCork.IDHash(testHeight)
	actualID := axelarcorkKeeper.SetScheduledAxelarCork(ctx, TestEVMChainID, testHeight, val, expectedCork)
	require.Equal(expectedID, actualID)
	actualCork, found := axelarcorkKeeper.GetScheduledAxelarCork(ctx, TestEVMChainID, testHeight, actualID, val, sampleCellarAddr, deadline)
	require.True(found)
	require.Equal(expectedCork, actualCork)

	actualCorks := axelarcorkKeeper.GetScheduledAxelarCorks(ctx, TestEVMChainID)
	require.Equal(&expectedCork, actualCorks[0].Cork)

	actualCorks = axelarcorkKeeper.GetScheduledAxelarCorksByID(ctx, TestEVMChainID, actualID)
	require.Len(actualCorks, 1)
	require.Equal(&expectedCork, actualCorks[0].Cork)
	require.Equal(hex.EncodeToString(expectedID), actualCorks[0].Id)

	actualHeights := axelarcorkKeeper.GetScheduledBlockHeights(ctx, TestEVMChainID)
	require.Equal(actualCorks[0].BlockHeight, actualHeights[0])

	actualCorks = axelarcorkKeeper.GetScheduledAxelarCorksByBlockHeight(ctx, TestEVMChainID, testHeight)
	require.Equal(&expectedCork, actualCorks[0].Cork)
	require.Equal(testHeight, actualCorks[0].BlockHeight)
	require.Equal(hex.EncodeToString(expectedID), actualCorks[0].Id)

	axelarcorkKeeper.DeleteScheduledAxelarCork(ctx, TestEVMChainID, testHeight, expectedID, sdk.ValAddress(val), sampleCellarAddr, deadline)
	require.Empty(axelarcorkKeeper.GetScheduledAxelarCorks(ctx, TestEVMChainID))
}

func (suite *KeeperTestSuite) TestGetWinningVotes() {
	ctx, axelarcorkKeeper := suite.ctx, suite.axelarcorkKeeper
	require := suite.Require()

	axelarcorkKeeper.SetChainConfiguration(ctx, TestEVMChainID, types.ChainConfiguration{
		Name:         "testevm",
		Id:           TestEVMChainID,
		ProxyAddress: "0x123",
	})

	testHeight := uint64(ctx.BlockHeight())
	deadline := uint64(100000000000000)
	cork := types.AxelarCork{
		EncodedContractCall:   []byte("testcall"),
		TargetContractAddress: sampleCellarHex,
		Deadline:              deadline,
	}
	_, bytes, err := bech32.DecodeAndConvert("somm1fcl08ymkl70dhyg3vmx4hjsqvxym7dawnp0zfp")
	require.NoError(err)
	require.Equal(20, len(bytes))
	axelarcorkKeeper.SetScheduledAxelarCork(ctx, TestEVMChainID, testHeight, bytes, cork)

	suite.stakingKeeper.EXPECT().GetLastTotalPower(ctx).Return(sdk.NewInt(100))
	suite.stakingKeeper.EXPECT().Validator(ctx, gomock.Any()).Return(suite.validator)
	suite.validator.EXPECT().GetConsensusPower(gomock.Any()).Return(int64(100))
	suite.stakingKeeper.EXPECT().PowerReduction(ctx).Return(sdk.OneInt())

	winningScheduledVotes := axelarcorkKeeper.GetApprovedScheduledAxelarCorks(ctx, TestEVMChainID)
	results := axelarcorkKeeper.GetAxelarCorkResults(ctx, TestEVMChainID)
	require.Len(winningScheduledVotes, 1)
	require.Equal(cork, winningScheduledVotes[0])
	require.Equal(&cork, results[0].Cork)
	require.True(results[0].Approved)
	require.Equal("100.000000000000000000", results[0].ApprovalPercentage)

	// scheduled cork should be deleted at the scheduled height
	require.Empty(axelarcorkKeeper.GetScheduledAxelarCorksByBlockHeight(ctx, TestEVMChainID, testHeight))
}

func (suite *KeeperTestSuite) TestCorkResults() {
	ctx, axelarcorkKeeper := suite.ctx, suite.axelarcorkKeeper
	require := suite.Require()

	require.Empty(axelarcorkKeeper.GetAxelarCorkResults(ctx, TestEVMChainID))

	testHeight := uint64(ctx.BlockHeight())
	deadline := uint64(100000000000000)
	cork := types.AxelarCork{
		EncodedContractCall:   []byte("testcall"),
		TargetContractAddress: sampleCellarHex,
		Deadline:              deadline,
	}
	id := cork.IDHash(testHeight)
	result := types.AxelarCorkResult{
		Cork:               &cork,
		BlockHeight:        testHeight,
		Approved:           true,
		ApprovalPercentage: "100.00",
	}
	axelarcorkKeeper.SetAxelarCorkResult(ctx, TestEVMChainID, id, result)
	actualResult, found := axelarcorkKeeper.GetAxelarCorkResult(ctx, TestEVMChainID, id)
	require.True(found)
	require.Equal(result, actualResult)

	results := axelarcorkKeeper.GetAxelarCorkResults(ctx, TestEVMChainID)
	require.Equal(&actualResult, results[0])

	axelarcorkKeeper.DeleteAxelarCorkResult(ctx, TestEVMChainID, id)
	require.Empty(axelarcorkKeeper.GetAxelarCorkResults(ctx, TestEVMChainID))
}

func (suite *KeeperTestSuite) TestParamSet() {
	ctx, axelarcorkKeeper := suite.ctx, suite.axelarcorkKeeper
	require := suite.Require()

	require.Panics(func() { axelarcorkKeeper.GetParamSet(ctx) })

	testGMPAccount := authtypes.NewModuleAddress("test-gmp-account")

	params := types.Params{
		IbcChannel:      "test-ibc-channel",
		IbcPort:         "test-ibc-port",
		GmpAccount:      testGMPAccount.String(),
		ExecutorAccount: "test-executor-account",
		TimeoutDuration: 10,
	}
	axelarcorkKeeper.SetParams(ctx, params)
	require.Equal(params, axelarcorkKeeper.GetParamSet(ctx))
}

func (suite *KeeperTestSuite) TestAxelarCorkContractCallNonces() {
	ctx, axelarcorkKeeper := suite.ctx, suite.axelarcorkKeeper
	require := suite.Require()

	chain1 := uint64(TestEVMChainID)
	chain2 := uint64(100)
	address1 := zeroAddressHex
	address2 := sampleCellarAddr.Hex()

	require.Equal(uint64(0), axelarcorkKeeper.GetAxelarContractCallNonce(ctx, chain1, address1))

	axelarcorkKeeper.IncrementAxelarContractCallNonce(ctx, chain1, address1)
	require.Equal(uint64(1), axelarcorkKeeper.GetAxelarContractCallNonce(ctx, chain1, address1))

	axelarcorkKeeper.SetAxelarContractCallNonce(ctx, chain1, address1, 5)
	require.Equal(uint64(5), axelarcorkKeeper.GetAxelarContractCallNonce(ctx, chain1, address1))

	axelarcorkKeeper.IncrementAxelarContractCallNonce(ctx, chain1, address2)
	axelarcorkKeeper.IncrementAxelarContractCallNonce(ctx, chain2, address2)
	axelarcorkKeeper.IncrementAxelarContractCallNonce(ctx, chain2, address2)
	require.Equal(uint64(1), axelarcorkKeeper.GetAxelarContractCallNonce(ctx, chain1, address2))
	require.Equal(uint64(2), axelarcorkKeeper.GetAxelarContractCallNonce(ctx, chain2, address2))

	var data []types.AxelarContractCallNonce
	axelarcorkKeeper.IterateAxelarContractCallNonces(ctx, func(chainID uint64, address common.Address, nonce uint64) bool {
		data = append(data, types.AxelarContractCallNonce{
			ChainId:         chainID,
			ContractAddress: address.Hex(),
			Nonce:           nonce,
		})

		return false
	})
	require.EqualValues(3, len(data))
	require.EqualValues([]types.AxelarContractCallNonce{
		{
			ChainId:         chain1,
			ContractAddress: address1,
			Nonce:           5,
		},
		{
			ChainId:         chain1,
			ContractAddress: address2,
			Nonce:           1,
		},
		{
			ChainId:         chain2,
			ContractAddress: address2,
			Nonce:           2,
		},
	}, data)
}

func (suite *KeeperTestSuite) TestAxelarProxyArgsEncoding() {
	require := suite.Require()

	cellars := []string{sampleCellarHex, oneAddressHex}

	_, err := types.EncodeLogicCallArgs(zeroAddressHex, 1, 100000000, []byte("testcall"))
	require.NoError(err)
	_, err = types.EncodeUpgradeArgs(zeroAddressHex, cellars)
	require.NoError(err)
}

// TestAxelarProxyUpgradeData tests the upgrade data set, get, and delete functions
func (suite *KeeperTestSuite) TestAxelarProxyUpgradeData() {
	ctx, axelarcorkKeeper := suite.ctx, suite.axelarcorkKeeper
	require := suite.Require()

	chainID1 := uint64(TestEVMChainID)
	chainID2 := uint64(3)
	cellars := []string{sampleCellarHex, oneAddressHex}
	payload, err := types.EncodeUpgradeArgs(zeroAddressHex, cellars)
	require.NoError(err)
	upgradeData1 := types.AxelarUpgradeData{
		ChainId:                   chainID1,
		Payload:                   payload,
		ExecutableHeightThreshold: 10000000,
	}
	upgradeData2 := types.AxelarUpgradeData{
		ChainId:                   chainID2,
		Payload:                   payload,
		ExecutableHeightThreshold: 10000001,
	}

	axelarcorkKeeper.SetAxelarProxyUpgradeData(ctx, chainID1, upgradeData1)
	actualUpgradeData1, found := axelarcorkKeeper.GetAxelarProxyUpgradeData(ctx, chainID1)
	require.True(found)
	require.Equal(upgradeData1, actualUpgradeData1)

	axelarcorkKeeper.SetAxelarProxyUpgradeData(ctx, chainID2, upgradeData2)
	var actualUpgradeData []types.AxelarUpgradeData
	axelarcorkKeeper.IterateAxelarProxyUpgradeData(ctx, func(chainID uint64, ud types.AxelarUpgradeData) (stop bool) {
		actualUpgradeData = append(actualUpgradeData, ud)

		return false
	})
	require.Equal(actualUpgradeData[0], upgradeData1)
	require.Equal(actualUpgradeData[1], upgradeData2)

	axelarcorkKeeper.DeleteAxelarProxyUpgradeData(ctx, chainID1)
	_, found = axelarcorkKeeper.GetAxelarProxyUpgradeData(ctx, chainID1)
	require.False(found)

}
