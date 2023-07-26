package keeper

import (
	"encoding/hex"

	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/peggyjv/sommelier/v6/x/axelarcork/types"
)

const TestEVMChainID = 2

func (suite *KeeperTestSuite) TestQueriesHappyPath() {
	ctx, corkKeeper := suite.ctx, suite.corkKeeper
	require := suite.Require()

	testGMPAccount := authtypes.NewModuleAddress("test-gmp-account")

	params := types.Params{
		IbcChannel:      "test-ibc-channel",
		IbcPort:         "test-ibc-port",
		GmpAccount:      testGMPAccount.String(),
		ExecutorAccount: "test-executor-account",
		TimeoutDuration: 10,
	}
	corkKeeper.SetParams(ctx, params)

	corkKeeper.SetChainConfiguration(ctx, TestEVMChainID, types.ChainConfiguration{
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
	corkKeeper.SetScheduledAxelarCork(ctx, TestEVMChainID, testHeight, val, cork)

	corkResult := types.AxelarCorkResult{
		Cork:               &cork,
		BlockHeight:        testHeight,
		Approved:           true,
		ApprovalPercentage: "100.00",
	}
	corkKeeper.SetAxelarCorkResult(ctx, TestEVMChainID, id, corkResult)

	corkKeeper.SetCellarIDs(ctx, TestEVMChainID, types.CellarIDSet{Ids: []string{"0x0000000000000000000000000000000000000000", "0x1111111111111111111111111111111111111111"}})

	paramsResult, err := corkKeeper.QueryParams(sdk.WrapSDKContext(ctx), &types.QueryParamsRequest{})
	require.Nil(err)
	require.Equal(params, paramsResult.Params)

	scheduledCorksResult, err := corkKeeper.QueryScheduledCorks(sdk.WrapSDKContext(ctx), &types.QueryScheduledCorksRequest{ChainId: TestEVMChainID})
	require.Nil(err)
	require.Equal(&expectedScheduledCork, scheduledCorksResult.Corks[0])

	scheduledCorksByHeightResult, err := corkKeeper.QueryScheduledCorksByBlockHeight(sdk.WrapSDKContext(ctx),
		&types.QueryScheduledCorksByBlockHeightRequest{
			BlockHeight: testHeight,
			ChainId:     TestEVMChainID,
		})
	require.Nil(err)
	require.Equal(&expectedScheduledCork, scheduledCorksByHeightResult.Corks[0])

	scheduledCorksByIDResult, err := corkKeeper.QueryScheduledCorksByID(sdk.WrapSDKContext(ctx),
		&types.QueryScheduledCorksByIDRequest{
			Id:      hex.EncodeToString(id),
			ChainId: TestEVMChainID,
		})
	require.Nil(err)
	require.Equal(&expectedScheduledCork, scheduledCorksByIDResult.Corks[0])

	blockHeightResult, err := corkKeeper.QueryScheduledBlockHeights(sdk.WrapSDKContext(ctx), &types.QueryScheduledBlockHeightsRequest{ChainId: TestEVMChainID})
	require.Nil(err)
	require.Equal(testHeight, blockHeightResult.BlockHeights[0])

	corkResultResult, err := corkKeeper.QueryCorkResult(sdk.WrapSDKContext(ctx), &types.QueryCorkResultRequest{Id: hex.EncodeToString(id), ChainId: TestEVMChainID})
	require.Nil(err)
	require.Equal(&corkResult, corkResultResult.CorkResult)

	corkResultsResult, err := corkKeeper.QueryCorkResults(sdk.WrapSDKContext(ctx), &types.QueryCorkResultsRequest{ChainId: TestEVMChainID})
	require.Nil(err)
	require.Equal(&corkResult, corkResultsResult.CorkResults[0])
}

func (suite *KeeperTestSuite) TestQueriesUnhappyPath() {
	ctx, corkKeeper := suite.ctx, suite.corkKeeper
	require := suite.Require()

	paramsResult, err := corkKeeper.QueryParams(sdk.WrapSDKContext(ctx), nil)
	require.Nil(paramsResult)
	require.NotNil(err)

	scheduledCorksResult, err := corkKeeper.QueryScheduledCorks(sdk.WrapSDKContext(ctx), nil)
	require.Nil(scheduledCorksResult)
	require.NotNil(err)

	scheduledCorksByHeightResult, err := corkKeeper.QueryScheduledCorksByBlockHeight(sdk.WrapSDKContext(ctx), nil)
	require.Nil(scheduledCorksByHeightResult)
	require.NotNil(err)

	scheduledCorksByIDResult, err := corkKeeper.QueryScheduledCorksByID(sdk.WrapSDKContext(ctx), nil)
	require.Nil(scheduledCorksByIDResult)
	require.NotNil(err)

	blockHeightResult, err := corkKeeper.QueryScheduledBlockHeights(sdk.WrapSDKContext(ctx), nil)
	require.Nil(blockHeightResult)
	require.NotNil(err)

	corkResultResult, err := corkKeeper.QueryCorkResult(sdk.WrapSDKContext(ctx), nil)
	require.Nil(corkResultResult)
	require.NotNil(err)

	corkResultsResult, err := corkKeeper.QueryCorkResults(sdk.WrapSDKContext(ctx), nil)
	require.Nil(corkResultsResult)
	require.NotNil(err)
}
