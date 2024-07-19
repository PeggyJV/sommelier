package keeper

import (
	"encoding/hex"

	sdk "github.com/cosmos/cosmos-sdk/types"
	types "github.com/peggyjv/sommelier/v7/x/cork/types/v2"
)

func (suite *KeeperTestSuite) TestQueriesHappyPath() {
	ctx, corkKeeper := suite.ctx, suite.corkKeeper
	require := suite.Require()

	params := types.DefaultParams()
	corkKeeper.SetParams(ctx, params)

	testHeight := uint64(ctx.BlockHeight())
	cork := types.Cork{
		EncodedContractCall:   []byte("testcall"),
		TargetContractAddress: sampleCellarHex,
	}
	id := cork.IDHash(testHeight)
	val := sdk.ValAddress("12345678901234567890")
	expectedScheduledCork := types.ScheduledCork{
		Cork:        &cork,
		BlockHeight: testHeight,
		Validator:   "cosmosvaloper1xyerxdp4xcmnswfsxyerxdp4xcmnswfs008wpw",
		Id:          id,
	}
	corkKeeper.SetScheduledCork(ctx, testHeight, val, cork)

	corkResult := types.CorkResult{
		Cork:               &cork,
		BlockHeight:        testHeight,
		Approved:           true,
		ApprovalPercentage: "100.00",
	}
	corkKeeper.SetCorkResult(ctx, id, corkResult)

	corkKeeper.SetCellarIDs(ctx, types.CellarIDSet{Ids: []string{"0x0000000000000000000000000000000000000000", "0x1111111111111111111111111111111111111111"}})

	paramsResult, err := corkKeeper.QueryParams(sdk.WrapSDKContext(ctx), &types.QueryParamsRequest{})
	require.Nil(err)
	require.Equal(params, paramsResult.Params)

	scheduledCorksResult, err := corkKeeper.QueryScheduledCorks(sdk.WrapSDKContext(ctx), &types.QueryScheduledCorksRequest{})
	require.Nil(err)
	require.Equal(&expectedScheduledCork, scheduledCorksResult.Corks[0])

	scheduledCorksByHeightResult, err := corkKeeper.QueryScheduledCorksByBlockHeight(sdk.WrapSDKContext(ctx), &types.QueryScheduledCorksByBlockHeightRequest{BlockHeight: testHeight})
	require.Nil(err)
	require.Equal(&expectedScheduledCork, scheduledCorksByHeightResult.Corks[0])

	scheduledCorksByIDResult, err := corkKeeper.QueryScheduledCorksByID(sdk.WrapSDKContext(ctx), &types.QueryScheduledCorksByIDRequest{Id: hex.EncodeToString(id)})
	require.Nil(err)
	require.Equal(&expectedScheduledCork, scheduledCorksByIDResult.Corks[0])

	blockHeightResult, err := corkKeeper.QueryScheduledBlockHeights(sdk.WrapSDKContext(ctx), &types.QueryScheduledBlockHeightsRequest{})
	require.Nil(err)
	require.Equal(testHeight, blockHeightResult.BlockHeights[0])

	corkResultResult, err := corkKeeper.QueryCorkResult(sdk.WrapSDKContext(ctx), &types.QueryCorkResultRequest{Id: hex.EncodeToString(id)})
	require.Nil(err)
	require.Equal(&corkResult, corkResultResult.CorkResult)

	corkResultsResult, err := corkKeeper.QueryCorkResults(sdk.WrapSDKContext(ctx), &types.QueryCorkResultsRequest{})
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
