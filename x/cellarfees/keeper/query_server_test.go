package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"
	"github.com/golang/mock/gomock"
	"github.com/peggyjv/sommelier/v7/app/params"
	cellarfeesTypes "github.com/peggyjv/sommelier/v7/x/cellarfees/types"
)

func (suite *KeeperTestSuite) TestQueriesHappyPath() {
	ctx, cellarfeesKeeper := suite.ctx, suite.cellarfeesKeeper
	require := suite.Require()

	// mock
	suite.accountKeeper.EXPECT().GetModuleAccount(ctx, cellarfeesTypes.ModuleName).Return(feesAccount)

	cellarfeesParams := cellarfeesTypes.DefaultParams()
	cellarfeesKeeper.SetParams(ctx, cellarfeesParams)

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
	require.Equal(&cellarfeesTypes.QueryParamsResponse{Params: cellarfeesParams}, paramsResponse)

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

	// QueryAPY
	blocksPerYear := 365 * 6
	bondedRatio := sdk.MustNewDecFromStr("0.2")
	stakingTotalSupply := sdk.NewInt(2_500_000_000_000)
	lastPeak := sdk.NewInt(10_000_000)
	cellarfeesKeeper.SetLastRewardSupplyPeak(ctx, lastPeak)
	cellarfeesParams.RewardEmissionPeriod = 10
	cellarfeesKeeper.SetParams(ctx, cellarfeesParams)

	//// mocks for QueryAPY
	suite.mintKeeper.EXPECT().GetParams(ctx).Return(minttypes.Params{
		BlocksPerYear: uint64(blocksPerYear),
		MintDenom:     params.BaseCoinUnit,
	})
	suite.accountKeeper.EXPECT().GetModuleAccount(ctx, gomock.Any()).Return(feesAccount)
	suite.bankKeeper.EXPECT().GetBalance(ctx, gomock.Any(), params.BaseCoinUnit).Return(sdk.NewCoin(params.BaseCoinUnit, sdk.NewInt(9_000_000)))
	suite.mintKeeper.EXPECT().BondedRatio(ctx).Return(bondedRatio)
	suite.mintKeeper.EXPECT().StakingTokenSupply(ctx).Return(stakingTotalSupply)

	APYResult, err := cellarfeesKeeper.QueryAPY(sdk.WrapSDKContext(ctx), &cellarfeesTypes.QueryAPYRequest{})
	require.Nil(err)
	require.Equal("0.004380000000000000", APYResult.Apy)
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
