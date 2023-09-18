package keeper

import (
	"encoding/hex"

	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/peggyjv/sommelier/v6/x/axelarcork/types"
)

const TestEVMChainID = 2

func (suite *KeeperTestSuite) TestQueriesHappyPath() {
	ctx, axelarcorkKeeper := suite.ctx, suite.axelarcorkKeeper
	require := suite.Require()

	testGMPAccount := authtypes.NewModuleAddress("test-gmp-account")

	params := types.Params{
		IbcChannel:      "test-ibc-channel",
		IbcPort:         "test-ibc-port",
		GmpAccount:      testGMPAccount.String(),
		ExecutorAccount: "test-executor-account",
		TimeoutDuration: 10,
	}
	axelarcorkKeeper.SetParams(ctx, params)

	axelarcorkKeeper.SetChainConfiguration(ctx, TestEVMChainID, types.ChainConfiguration{
		Name:         "testevm",
		Id:           TestEVMChainID,
		ProxyAddress: "0x123",
	})

	testHeight := uint64(ctx.BlockHeight())
	cork := types.AxelarCork{
		EncodedContractCall:   []byte("testcall"),
		TargetContractAddress: sampleCellarHex,
		ChainId:               TestEVMChainID,
	}
	id := cork.IDHash(TestEVMChainID, testHeight)

	val := sdk.ValAddress("12345678901234567890")
	expectedScheduledCork := types.ScheduledAxelarCork{
		Cork:        &cork,
		BlockHeight: testHeight,
		Validator:   "cosmosvaloper1xyerxdp4xcmnswfsxyerxdp4xcmnswfs008wpw",
		Id:          hex.EncodeToString(id),
	}
	axelarcorkKeeper.SetScheduledAxelarCork(ctx, TestEVMChainID, testHeight, val, cork)

	corkResult := types.AxelarCorkResult{
		Cork:               &cork,
		BlockHeight:        testHeight,
		Approved:           true,
		ApprovalPercentage: "100.00",
	}
	axelarcorkKeeper.SetAxelarCorkResult(ctx, TestEVMChainID, id, corkResult)

	axelarcorkKeeper.SetCellarIDs(ctx, TestEVMChainID, types.CellarIDSet{Ids: []string{"0x0000000000000000000000000000000000000000", "0x1111111111111111111111111111111111111111"}})

	paramsResult, err := axelarcorkKeeper.QueryParams(sdk.WrapSDKContext(ctx), &types.QueryParamsRequest{})
	require.Nil(err)
	require.Equal(params, paramsResult.Params)

	scheduledCorksResult, err := axelarcorkKeeper.QueryScheduledCorks(sdk.WrapSDKContext(ctx), &types.QueryScheduledCorksRequest{ChainId: TestEVMChainID})
	require.Nil(err)
	require.Equal(&expectedScheduledCork, scheduledCorksResult.Corks[0])

	scheduledCorksByHeightResult, err := axelarcorkKeeper.QueryScheduledCorksByBlockHeight(sdk.WrapSDKContext(ctx),
		&types.QueryScheduledCorksByBlockHeightRequest{
			BlockHeight: testHeight,
			ChainId:     TestEVMChainID,
		})
	require.Nil(err)
	require.Equal(&expectedScheduledCork, scheduledCorksByHeightResult.Corks[0])

	scheduledCorksByIDResult, err := axelarcorkKeeper.QueryScheduledCorksByID(sdk.WrapSDKContext(ctx),
		&types.QueryScheduledCorksByIDRequest{
			Id:      hex.EncodeToString(id),
			ChainId: TestEVMChainID,
		})
	require.Nil(err)
	require.Equal(&expectedScheduledCork, scheduledCorksByIDResult.Corks[0])

	blockHeightResult, err := axelarcorkKeeper.QueryScheduledBlockHeights(sdk.WrapSDKContext(ctx), &types.QueryScheduledBlockHeightsRequest{ChainId: TestEVMChainID})
	require.Nil(err)
	require.Equal(testHeight, blockHeightResult.BlockHeights[0])

	corkResultResult, err := axelarcorkKeeper.QueryCorkResult(sdk.WrapSDKContext(ctx), &types.QueryCorkResultRequest{Id: hex.EncodeToString(id), ChainId: TestEVMChainID})
	require.Nil(err)
	require.Equal(&corkResult, corkResultResult.CorkResult)

	corkResultsResult, err := axelarcorkKeeper.QueryCorkResults(sdk.WrapSDKContext(ctx), &types.QueryCorkResultsRequest{ChainId: TestEVMChainID})
	require.Nil(err)
	require.Equal(&corkResult, corkResultsResult.CorkResults[0])
}

func (suite *KeeperTestSuite) TestQueriesUnhappyPath() {
	ctx, axelarcorkKeeper := suite.ctx, suite.axelarcorkKeeper
	require := suite.Require()

	paramsResult, err := axelarcorkKeeper.QueryParams(sdk.WrapSDKContext(ctx), nil)
	require.Nil(paramsResult)
	require.NotNil(err)

	scheduledCorksResult, err := axelarcorkKeeper.QueryScheduledCorks(sdk.WrapSDKContext(ctx), nil)
	require.Nil(scheduledCorksResult)
	require.NotNil(err)

	scheduledCorksByHeightResult, err := axelarcorkKeeper.QueryScheduledCorksByBlockHeight(sdk.WrapSDKContext(ctx), nil)
	require.Nil(scheduledCorksByHeightResult)
	require.NotNil(err)

	scheduledCorksByIDResult, err := axelarcorkKeeper.QueryScheduledCorksByID(sdk.WrapSDKContext(ctx), nil)
	require.Nil(scheduledCorksByIDResult)
	require.NotNil(err)

	blockHeightResult, err := axelarcorkKeeper.QueryScheduledBlockHeights(sdk.WrapSDKContext(ctx), nil)
	require.Nil(blockHeightResult)
	require.NotNil(err)

	corkResultResult, err := axelarcorkKeeper.QueryCorkResult(sdk.WrapSDKContext(ctx), nil)
	require.Nil(corkResultResult)
	require.NotNil(err)

	corkResultsResult, err := axelarcorkKeeper.QueryCorkResults(sdk.WrapSDKContext(ctx), nil)
	require.Nil(corkResultsResult)
	require.NotNil(err)
}