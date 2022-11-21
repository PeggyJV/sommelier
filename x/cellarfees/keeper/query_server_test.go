package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	cellarfeesTypes "github.com/peggyjv/sommelier/v4/x/cellarfees/types"
)

func (suite *KeeperTestSuite) TestQueriesHappyPath() {
	ctx, cellarfeesKeeper := suite.ctx, suite.cellarfeesKeeper
	require := suite.Require()

	// mock
	suite.accountKeeper.EXPECT().GetModuleAccount(ctx, cellarfeesTypes.ModuleName).Return(feesAccount)

	params := cellarfeesTypes.DefaultParams()
	cellarfeesKeeper.SetParams(ctx, params)

	expectedLastRewardSupplyPeakAmount := sdk.NewInt(25000)
	cellarfeesKeeper.SetLastRewardSupplyPeak(ctx, expectedLastRewardSupplyPeakAmount)

	expectedFeeAccrualCounters := cellarfeesTypes.FeeAccrualCounters{
		Counters: []cellarfeesTypes.FeeAccrualCounter{
			{
				Denom: "denom1",
				Count: 2,
			},
			{
				Denom: "denom2",
				Count: 0,
			},
		},
	}
	cellarfeesKeeper.SetFeeAccrualCounters(ctx, expectedFeeAccrualCounters)

	// QueryParams
	paramsResponse, err := cellarfeesKeeper.QueryParams(sdk.WrapSDKContext(ctx), &cellarfeesTypes.QueryParamsRequest{})
	require.Nil(err)
	require.Equal(&cellarfeesTypes.QueryParamsResponse{Params: params}, paramsResponse)

	// QueryModuleAccounts
	moduleAccountsResponse, err := cellarfeesKeeper.QueryModuleAccounts(sdk.WrapSDKContext(ctx), &cellarfeesTypes.QueryModuleAccountsRequest{})
	require.Nil(err)
	require.Equal(&cellarfeesTypes.QueryModuleAccountsResponse{FeesAddress: feesAccount.GetAddress().String()}, moduleAccountsResponse)

	// QueryLastRewardSupplyPeak
	lastRewardSupplyPeakResponse, err := cellarfeesKeeper.QueryLastRewardSupplyPeak(sdk.WrapSDKContext(ctx), &cellarfeesTypes.QueryLastRewardSupplyPeakRequest{})
	require.Nil(err)
	require.Equal(&cellarfeesTypes.QueryLastRewardSupplyPeakResponse{LastRewardSupplyPeak: expectedLastRewardSupplyPeakAmount}, lastRewardSupplyPeakResponse)

	// QueryFeeAccrualCounters
	feeAccrualCountersResponse, err := cellarfeesKeeper.QueryFeeAccrualCounters(sdk.WrapSDKContext(ctx), &cellarfeesTypes.QueryFeeAccrualCountersRequest{})
	require.Nil(err)
	require.Equal(&cellarfeesTypes.QueryFeeAccrualCountersResponse{FeeAccrualCounters: expectedFeeAccrualCounters}, feeAccrualCountersResponse)
}

func (suite *KeeperTestSuite) TestQueriesUnhappyPath() {
	ctx, cellarfeesKeeper := suite.ctx, suite.cellarfeesKeeper
	require := suite.Require()

	// QueryParams
	paramsResponse, err := cellarfeesKeeper.QueryParams(sdk.WrapSDKContext(ctx), nil)
	require.Nil(paramsResponse)
	require.NotNil(err)

	// QueryModuleAccounts
	moduleAccountsResponse, err := cellarfeesKeeper.QueryModuleAccounts(sdk.WrapSDKContext(ctx), nil)
	require.Nil(moduleAccountsResponse)
	require.NotNil(err)

	// QueryLastRewardSupplyPeak
	lastRewardSupplyPeakResponse, err := cellarfeesKeeper.QueryLastRewardSupplyPeak(sdk.WrapSDKContext(ctx), nil)
	require.Nil(lastRewardSupplyPeakResponse)
	require.NotNil(err)

	// QueryFeeAccrualCounters
	feeAccrualCountersResponse, err := cellarfeesKeeper.QueryFeeAccrualCounters(sdk.WrapSDKContext(ctx), nil)
	require.Nil(feeAccrualCountersResponse)
	require.NotNil(err)
}
