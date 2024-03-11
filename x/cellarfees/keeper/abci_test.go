package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/golang/mock/gomock"
	appParams "github.com/peggyjv/sommelier/v7/app/params"
	auctiontypes "github.com/peggyjv/sommelier/v7/x/auction/types"
	cellarfeestypes "github.com/peggyjv/sommelier/v7/x/cellarfees/types"
	cellarfeestypesv2 "github.com/peggyjv/sommelier/v7/x/cellarfees/types/v2"
)

func (suite *KeeperTestSuite) TestBeginBlockerZeroRewardsBalance() {
	ctx, cellarfeesKeeper := suite.ctx, suite.cellarfeesKeeper

	require := suite.Require()

	params := cellarfeestypesv2.DefaultParams()
	cellarfeesKeeper.SetParams(ctx, params)

	// mocks
	suite.accountKeeper.EXPECT().GetModuleAccount(ctx, cellarfeestypes.ModuleName).Return(feesAccount)
	suite.bankKeeper.EXPECT().GetBalance(ctx, feesAccount.GetAddress(), appParams.BaseCoinUnit).Return(sdk.NewCoin(appParams.BaseCoinUnit, sdk.ZeroInt()))
	suite.bankKeeper.EXPECT().SendCoinsFromModuleToModule(ctx, gomock.Any(), gomock.Any(), gomock.Any()).Times(0)

	// Note EndBlocker is only run once for completeness, since it has no code in it
	require.NotPanics(func() { cellarfeesKeeper.BeginBlocker(ctx) })
	require.NotPanics(func() { cellarfeesKeeper.EndBlocker(ctx) })
}

func (suite *KeeperTestSuite) TestBeginBlockerWithRewardBalanceAndPreviousPeakZero() {
	ctx, cellarfeesKeeper := suite.ctx, suite.cellarfeesKeeper

	require := suite.Require()

	params := cellarfeestypesv2.DefaultParams()
	cellarfeesKeeper.SetParams(ctx, params)

	// mocks
	suite.accountKeeper.EXPECT().GetModuleAccount(ctx, cellarfeestypes.ModuleName).Return(feesAccount)

	rewardSupply := sdk.NewCoin(appParams.BaseCoinUnit, sdk.NewInt(1000000))
	emissionPeriod := sdk.NewInt(int64(params.RewardEmissionPeriod))
	suite.bankKeeper.EXPECT().GetBalance(ctx, feesAccount.GetAddress(), appParams.BaseCoinUnit).Return(sdk.NewCoin(appParams.BaseCoinUnit, rewardSupply.Amount))
	suite.cellarfeesKeeper.SetLastRewardSupplyPeak(ctx, sdk.NewInt(0))

	expectedEmissionAmount := rewardSupply.Amount.Quo(emissionPeriod)
	expectedEmission := sdk.NewCoin(appParams.BaseCoinUnit, expectedEmissionAmount)
	suite.bankKeeper.EXPECT().SendCoinsFromModuleToModule(ctx, gomock.Any(), gomock.Any(), sdk.NewCoins(expectedEmission)).Times(1)

	require.NotPanics(func() { cellarfeesKeeper.BeginBlocker(ctx) })
}

func (suite *KeeperTestSuite) TestBeginBlockerWithRewardBalanceAndHigherPreviousPeak() {
	ctx, cellarfeesKeeper := suite.ctx, suite.cellarfeesKeeper
	require := suite.Require()

	params := cellarfeestypesv2.DefaultParams()
	cellarfeesKeeper.SetParams(ctx, params)

	// mocks
	suite.accountKeeper.EXPECT().GetModuleAccount(ctx, cellarfeestypes.ModuleName).Return(feesAccount)

	rewardSupply := sdk.NewCoin(appParams.BaseCoinUnit, sdk.NewInt(1000000))
	emissionPeriod := sdk.NewInt(int64(params.RewardEmissionPeriod))
	suite.bankKeeper.EXPECT().GetBalance(ctx, feesAccount.GetAddress(), appParams.BaseCoinUnit).Return(sdk.NewCoin(appParams.BaseCoinUnit, rewardSupply.Amount))

	previousPeakAmount := sdk.NewInt(2000000)
	cellarfeesKeeper.SetLastRewardSupplyPeak(ctx, previousPeakAmount)

	expectedEmissionAmount := previousPeakAmount.Quo(emissionPeriod)
	expectedEmission := sdk.NewCoin(appParams.BaseCoinUnit, expectedEmissionAmount)
	suite.bankKeeper.EXPECT().SendCoinsFromModuleToModule(ctx, gomock.Any(), gomock.Any(), sdk.NewCoins(expectedEmission)).Times(1)

	require.NotPanics(func() { cellarfeesKeeper.BeginBlocker(ctx) })
}

func (suite *KeeperTestSuite) TestBeginBlockerWithRewardBalanceAndLowerPreviousPeak() {
	ctx, cellarfeesKeeper := suite.ctx, suite.cellarfeesKeeper

	require := suite.Require()

	params := cellarfeestypesv2.DefaultParams()
	cellarfeesKeeper.SetParams(ctx, params)

	// mocks
	suite.accountKeeper.EXPECT().GetModuleAccount(ctx, cellarfeestypes.ModuleName).Return(feesAccount)

	rewardSupply := sdk.NewCoin(appParams.BaseCoinUnit, sdk.NewInt(1000000))
	emissionPeriod := sdk.NewInt(int64(params.RewardEmissionPeriod))
	suite.bankKeeper.EXPECT().GetBalance(ctx, feesAccount.GetAddress(), appParams.BaseCoinUnit).Return(sdk.NewCoin(appParams.BaseCoinUnit, rewardSupply.Amount))

	previousPeakAmount := sdk.NewInt(500000)
	cellarfeesKeeper.SetLastRewardSupplyPeak(ctx, previousPeakAmount)

	expectedEmissionAmount := rewardSupply.Amount.Quo(emissionPeriod)
	expectedEmission := sdk.NewCoin(appParams.BaseCoinUnit, expectedEmissionAmount)
	suite.bankKeeper.EXPECT().SendCoinsFromModuleToModule(ctx, gomock.Any(), gomock.Any(), sdk.NewCoins(expectedEmission)).Times(1)

	require.NotPanics(func() { cellarfeesKeeper.BeginBlocker(ctx) })
}

// If the emission calculation underflows to zero, it should be set to 1
func (suite *KeeperTestSuite) TestBeginBlockerEmissionCalculationUnderflowsToZero() {
	ctx, cellarfeesKeeper := suite.ctx, suite.cellarfeesKeeper
	require := suite.Require()

	params := cellarfeestypesv2.DefaultParams()
	cellarfeesKeeper.SetParams(ctx, params)
	cellarfeesKeeper.SetLastRewardSupplyPeak(ctx, sdk.ZeroInt())

	// mocks
	suite.accountKeeper.EXPECT().GetModuleAccount(ctx, cellarfeestypes.ModuleName).Return(feesAccount)

	rewardSupply := sdk.NewCoin(appParams.BaseCoinUnit, sdk.NewInt(1))
	suite.bankKeeper.EXPECT().GetBalance(ctx, feesAccount.GetAddress(), appParams.BaseCoinUnit).Return(sdk.NewCoin(appParams.BaseCoinUnit, rewardSupply.Amount))

	expectedEmission := sdk.NewCoin(appParams.BaseCoinUnit, sdk.OneInt())
	suite.bankKeeper.EXPECT().SendCoinsFromModuleToModule(ctx, gomock.Any(), gomock.Any(), sdk.NewCoins(expectedEmission)).Times(1)

	require.NotPanics(func() { cellarfeesKeeper.BeginBlocker(ctx) })
}

// If the calculated emission is greater than the remaining supply, it should be set to the remaining supply
func (suite *KeeperTestSuite) TestBeginBlockerEmissionGreaterThanRewardSupply() {
	ctx, cellarfeesKeeper := suite.ctx, suite.cellarfeesKeeper

	require := suite.Require()
	params := cellarfeestypesv2.DefaultParams()
	cellarfeesKeeper.SetParams(ctx, params)
	cellarfeesKeeper.SetLastRewardSupplyPeak(ctx, sdk.NewInt(1000000))

	// mocks
	suite.accountKeeper.EXPECT().GetModuleAccount(ctx, cellarfeestypes.ModuleName).Return(feesAccount)

	rewardSupply := sdk.NewCoin(appParams.BaseCoinUnit, sdk.NewInt(1))
	suite.bankKeeper.EXPECT().GetBalance(ctx, feesAccount.GetAddress(), appParams.BaseCoinUnit).Return(sdk.NewCoin(appParams.BaseCoinUnit, rewardSupply.Amount))

	expectedEmission := sdk.NewCoin(appParams.BaseCoinUnit, rewardSupply.Amount)
	suite.bankKeeper.EXPECT().SendCoinsFromModuleToModule(ctx, gomock.Any(), gomock.Any(), sdk.NewCoins(expectedEmission)).Times(1)

	require.NotPanics(func() { cellarfeesKeeper.BeginBlocker(ctx) })
}

func (suite *KeeperTestSuite) TestHandleFeeAuctionsHappyPath() {
	ctx, cellarfeesKeeper := suite.ctx, suite.cellarfeesKeeper
	params := cellarfeestypesv2.DefaultParams()

	params.AuctionInterval = 1
	params.AuctionThresholdUsdValue = sdk.MustNewDecFromStr("100.00")

	cellarfeesKeeper.SetParams(ctx, params)

	denom1 := "denom1"
	denom2 := "denom2"
	denom3 := "denom3"
	amount1 := sdk.NewInt(1000000)
	amount2 := sdk.NewInt(2000000000000)
	amount3 := sdk.NewInt(3000000000000000000)
	balance1 := sdk.NewCoin(denom1, amount1)
	balance2 := sdk.NewCoin(denom2, amount2)
	balance3 := sdk.NewCoin(denom3, amount3)
	price1 := sdk.NewDec(100)
	price2 := sdk.NewDec(50)
	price3 := sdk.NewDec(33)
	tokenPrices := []*auctiontypes.TokenPrice{
		{
			Exponent: 6,
			Denom:    denom1,
			UsdPrice: price1,
		},
		{
			Exponent: 12,
			Denom:    denom2,
			UsdPrice: price2,
		},
		{
			Exponent: 18,
			Denom:    denom3,
			UsdPrice: price3,
		},
	}

	suite.auctionKeeper.EXPECT().GetTokenPrices(ctx).Return(tokenPrices)
	suite.accountKeeper.EXPECT().GetModuleAccount(ctx, cellarfeestypes.ModuleName).Return(feesAccount).Times(len(tokenPrices) * 2)
	suite.bankKeeper.EXPECT().GetBalance(ctx, feesAccount.GetAddress(), denom1).Return(balance1).Times(2)
	suite.bankKeeper.EXPECT().GetBalance(ctx, feesAccount.GetAddress(), denom2).Return(balance2).Times(2)
	suite.bankKeeper.EXPECT().GetBalance(ctx, feesAccount.GetAddress(), denom3).Return(balance3).Times(2)

	// retreiving module account
	suite.accountKeeper.EXPECT().GetModuleAccount(ctx, cellarfeestypes.ModuleName).Return(feesAccount)

	// no active auctions
	suite.auctionKeeper.EXPECT().GetActiveAuctions(ctx).Return([]*auctiontypes.Auction{})
	suite.auctionKeeper.EXPECT().BeginAuction(ctx, balance1, params.InitialPriceDecreaseRate, params.PriceDecreaseBlockInterval, cellarfeestypes.ModuleName, cellarfeestypes.ModuleName).Return(nil)
	expectedAuction1 := auctiontypes.Auction{
		Id:                         1,
		StartingTokensForSale:      balance1,
		InitialPriceDecreaseRate:   params.InitialPriceDecreaseRate,
		CurrentPriceDecreaseRate:   params.InitialPriceDecreaseRate,
		StartBlock:                 uint64(ctx.BlockHeight()),
		EndBlock:                   0,
		PriceDecreaseBlockInterval: params.PriceDecreaseBlockInterval,
		InitialUnitPriceInUsomm:    price1,
		CurrentUnitPriceInUsomm:    price1,
		FundingModuleAccount:       cellarfeestypes.ModuleName,
		ProceedsModuleAccount:      cellarfeestypes.ModuleName,
	}
	suite.auctionKeeper.EXPECT().GetActiveAuctions(ctx).Return([]*auctiontypes.Auction{&expectedAuction1})
	suite.auctionKeeper.EXPECT().BeginAuction(ctx, balance2, params.InitialPriceDecreaseRate, params.PriceDecreaseBlockInterval, cellarfeestypes.ModuleName, cellarfeestypes.ModuleName).Return(nil)

	// we only set expected calls for two auctions because the third token price is $99, so no auction should be started.

	cellarfeesKeeper.handleFeeAuctions(ctx)
}
