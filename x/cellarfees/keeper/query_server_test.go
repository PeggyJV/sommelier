package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"
	"github.com/golang/mock/gomock"
	"github.com/peggyjv/sommelier/v9/app/params"
	auctiontypes "github.com/peggyjv/sommelier/v9/x/auction/types"
	cellarfeestypes "github.com/peggyjv/sommelier/v9/x/cellarfees/types"
	cellarfeestypesv2 "github.com/peggyjv/sommelier/v9/x/cellarfees/types/v2"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (suite *KeeperTestSuite) TestQueriesHappyPath() {
	ctx, cellarfeesKeeper := suite.ctx, suite.cellarfeesKeeper
	require := suite.Require()

	// mock
	suite.accountKeeper.EXPECT().GetModuleAccount(ctx, cellarfeestypes.ModuleName).Return(feesAccount)

	cellarfeesParams := cellarfeestypesv2.DefaultParams()
	cellarfeesKeeper.SetParams(ctx, cellarfeesParams)

	expectedLastRewardSupplyPeakAmount := sdk.NewInt(25000)
	cellarfeesKeeper.SetLastRewardSupplyPeak(ctx, expectedLastRewardSupplyPeakAmount)
	// QueryParams
	paramsResponse, err := cellarfeesKeeper.QueryParams(sdk.WrapSDKContext(ctx), &cellarfeestypesv2.QueryParamsRequest{})
	require.Nil(err)
	require.Equal(&cellarfeestypesv2.QueryParamsResponse{Params: cellarfeesParams}, paramsResponse)

	// QueryModuleAccounts
	moduleAccountsResponse, err := cellarfeesKeeper.QueryModuleAccounts(sdk.WrapSDKContext(ctx), &cellarfeestypesv2.QueryModuleAccountsRequest{})
	require.Nil(err)
	require.Equal(&cellarfeestypesv2.QueryModuleAccountsResponse{FeesAddress: feesAccount.GetAddress().String()}, moduleAccountsResponse)

	// QueryLastRewardSupplyPeak
	lastRewardSupplyPeakResponse, err := cellarfeesKeeper.QueryLastRewardSupplyPeak(sdk.WrapSDKContext(ctx), &cellarfeestypesv2.QueryLastRewardSupplyPeakRequest{})
	require.Nil(err)
	require.Equal(&cellarfeestypesv2.QueryLastRewardSupplyPeakResponse{LastRewardSupplyPeak: expectedLastRewardSupplyPeakAmount}, lastRewardSupplyPeakResponse)

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

	APYResult, err := cellarfeesKeeper.QueryAPY(sdk.WrapSDKContext(ctx), &cellarfeestypesv2.QueryAPYRequest{})
	require.Nil(err)
	require.Equal("0.004380000000000000", APYResult.Apy)

	// QueryFeeTokenBalance
	denom := feeDenom
	amount := sdk.NewInt(1000000)
	suite.bankKeeper.EXPECT().GetDenomMetaData(ctx, denom).Return(banktypes.Metadata{}, true)
	suite.accountKeeper.EXPECT().GetModuleAccount(ctx, cellarfeestypes.ModuleName).Return(feesAccount).Times(1)
	suite.auctionKeeper.EXPECT().GetTokenPrice(ctx, denom).Return(auctiontypes.TokenPrice{
		Exponent: 6,
		UsdPrice: sdk.NewDec(100),
	}, true)
	suite.bankKeeper.EXPECT().GetAllBalances(ctx, feesAccount.GetAddress()).Return(sdk.Coins{sdk.NewCoin(denom, amount)})

	expectedFeeTokenBalance := cellarfeestypesv2.FeeTokenBalance{
		Balance:  sdk.NewCoin(denom, amount),
		UsdValue: 100.00,
	}
	feeTokenBalanceResponse, err := cellarfeesKeeper.QueryFeeTokenBalance(sdk.WrapSDKContext(ctx), &cellarfeestypesv2.QueryFeeTokenBalanceRequest{
		Denom: denom,
	})
	require.Nil(err)
	require.Equal(&expectedFeeTokenBalance, feeTokenBalanceResponse.Balance)

	// QueryFeeTokenBalances
	suite.SetupTest()
	ctx, cellarfeesKeeper = suite.ctx, suite.cellarfeesKeeper
	denom1 := "testdenom1"
	denom2 := "testdenom2"
	denom3 := "testdenom3"
	amount1 := sdk.NewInt(1000000)
	amount2 := sdk.NewInt(2000000)
	amount3 := sdk.NewInt(3000000)
	balance1 := sdk.NewCoin(denom1, amount1)
	balance2 := sdk.NewCoin(denom2, amount2)
	balance3 := sdk.NewCoin(denom3, amount3)
	tokenPrice1 := auctiontypes.TokenPrice{
		Exponent: 6,
		UsdPrice: sdk.NewDec(100),
		Denom:    denom1,
	}
	tokenPrice2 := auctiontypes.TokenPrice{
		Exponent: 12,
		UsdPrice: sdk.NewDec(50),
		Denom:    denom2,
	}
	tokenPrice3 := auctiontypes.TokenPrice{
		Exponent: 18,
		UsdPrice: sdk.NewDec(25),
		Denom:    denom3,
	}
	tokenPrices := []*auctiontypes.TokenPrice{&tokenPrice1, &tokenPrice2, &tokenPrice3}
	suite.accountKeeper.EXPECT().GetModuleAccount(ctx, cellarfeestypes.ModuleName).Return(feesAccount).Times(3)
	suite.bankKeeper.EXPECT().GetAllBalances(ctx, feesAccount.GetAddress()).Return(sdk.Coins{balance1, balance2, balance3}).Times(3)
	suite.bankKeeper.EXPECT().GetBalance(ctx, feesAccount.GetAddress(), denom1).Return(balance1)
	suite.bankKeeper.EXPECT().GetBalance(ctx, feesAccount.GetAddress(), denom2).Return(balance2)
	suite.bankKeeper.EXPECT().GetBalance(ctx, feesAccount.GetAddress(), denom3).Return(balance3)
	suite.auctionKeeper.EXPECT().GetTokenPrice(ctx, denom1).Return(*tokenPrices[0], true)
	suite.auctionKeeper.EXPECT().GetTokenPrice(ctx, denom2).Return(*tokenPrices[1], true)
	suite.auctionKeeper.EXPECT().GetTokenPrice(ctx, denom3).Return(*tokenPrices[2], true)

	expectedFeeTokenBalances := []*cellarfeestypesv2.FeeTokenBalance{
		{
			Balance:  balance1,
			UsdValue: cellarfeesKeeper.GetBalanceUsdValue(ctx, balance1, tokenPrice1).MustFloat64(),
		},
		{
			Balance:  balance2,
			UsdValue: cellarfeesKeeper.GetBalanceUsdValue(ctx, balance2, tokenPrice2).MustFloat64(),
		},
		{
			Balance:  balance3,
			UsdValue: cellarfeesKeeper.GetBalanceUsdValue(ctx, balance3, tokenPrice3).MustFloat64(),
		},
	}

	feeTokenBalancesResponse, err := cellarfeesKeeper.QueryFeeTokenBalances(sdk.WrapSDKContext(ctx), &cellarfeestypesv2.QueryFeeTokenBalancesRequest{})
	require.Nil(err)
	require.Equal(expectedFeeTokenBalances, feeTokenBalancesResponse.Balances)
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

	// QueryFeeTokenBalance
	denom := feeDenom
	suite.accountKeeper.EXPECT().GetModuleAccount(ctx, cellarfeestypes.ModuleName).Return(feesAccount)
	suite.bankKeeper.EXPECT().GetAllBalances(ctx, feesAccount.GetAddress()).Return(sdk.Coins{sdk.NewCoin(denom, sdk.NewInt(1000000))})
	suite.auctionKeeper.EXPECT().GetTokenPrice(ctx, denom).Return(auctiontypes.TokenPrice{}, false).Times(2)
	feeTokenBalanceResponse, err := cellarfeesKeeper.QueryFeeTokenBalance(sdk.WrapSDKContext(ctx), nil)
	require.Nil(feeTokenBalanceResponse)
	require.NotNil(err)
	require.Equal(status.Code(err), codes.InvalidArgument)

	feeTokenBalanceResponse, err = cellarfeesKeeper.QueryFeeTokenBalance(sdk.WrapSDKContext(ctx), &cellarfeestypesv2.QueryFeeTokenBalanceRequest{})
	require.Nil(feeTokenBalanceResponse)
	require.NotNil(err)
	require.Equal(status.Code(err), codes.InvalidArgument)

	feeTokenBalanceResponse, err = cellarfeesKeeper.QueryFeeTokenBalance(sdk.WrapSDKContext(ctx), &cellarfeestypesv2.QueryFeeTokenBalanceRequest{Denom: denom})
	require.Nil(feeTokenBalanceResponse)
	require.NotNil(err)
	require.Equal(status.Code(err), codes.NotFound)

	suite.accountKeeper.EXPECT().GetModuleAccount(ctx, gomock.Any()).Return(feesAccount)
	suite.bankKeeper.EXPECT().GetAllBalances(ctx, feesAccount.GetAddress()).Return(sdk.Coins{sdk.NewCoin(denom, sdk.NewInt(1000000))})
	feeTokenBalanceResponse, err = cellarfeesKeeper.QueryFeeTokenBalance(sdk.WrapSDKContext(ctx), &cellarfeestypesv2.QueryFeeTokenBalanceRequest{Denom: denom})
	require.Nil(feeTokenBalanceResponse)
	require.NotNil(err)
	require.Equal(status.Code(err), codes.NotFound)
}
