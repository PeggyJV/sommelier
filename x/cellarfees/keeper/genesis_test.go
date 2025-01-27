package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	cellarfeestypesv2 "github.com/peggyjv/sommelier/v9/x/cellarfees/types/v2"
)

func (suite *KeeperTestSuite) TestImportingEmptyGenesis() {
	ctx, cellarfeesKeeper := suite.ctx, suite.cellarfeesKeeper
	require := suite.Require()

	testGenesis := cellarfeestypesv2.GenesisState{}

	// Canary to make sure validate basic is being run
	require.Panics(func() { cellarfeesKeeper.InitGenesis(ctx, testGenesis) })

	testGenesis = cellarfeestypesv2.DefaultGenesisState()
	require.NotPanics(func() {
		suite.accountKeeper.EXPECT().GetModuleAccount(ctx, feesAccount.GetName()).Return(feesAccount)
		cellarfeesKeeper.InitGenesis(ctx, testGenesis)
	})

	require.Zero(cellarfeesKeeper.GetLastRewardSupplyPeak(ctx).Int64())
	require.Equal(cellarfeesKeeper.GetParams(ctx), cellarfeestypesv2.DefaultParams())
}

func (suite *KeeperTestSuite) TestImportingPopulatedGenesis() {
	ctx, cellarfeesKeeper := suite.ctx, suite.cellarfeesKeeper
	require := suite.Require()

	testGenesis := cellarfeestypesv2.GenesisState{}

	// Canary to make sure validate basic is being run
	require.Panics(func() { cellarfeesKeeper.InitGenesis(ctx, testGenesis) })
	testGenesis.LastRewardSupplyPeak = sdk.NewInt(1337)
	testGenesis.Params.InitialPriceDecreaseRate = sdk.MustNewDecFromStr("0.01")
	testGenesis.Params.PriceDecreaseBlockInterval = 10
	testGenesis.Params.RewardEmissionPeriod = 600
	testGenesis.Params.AuctionInterval = 1000
	testGenesis.Params.AuctionThresholdUsdValue = sdk.NewDec(1000000)
	testGenesis.Params.ProceedsPortion = sdk.MustNewDecFromStr("0.5")

	require.NotPanics(func() {
		suite.accountKeeper.EXPECT().GetModuleAccount(ctx, feesAccount.GetName()).Return(feesAccount)
		cellarfeesKeeper.InitGenesis(ctx, testGenesis)
	})

	require.Equal(testGenesis.LastRewardSupplyPeak, cellarfeesKeeper.GetLastRewardSupplyPeak(ctx))
	require.Equal(testGenesis.Params, cellarfeesKeeper.GetParams(ctx))
}

func (suite *KeeperTestSuite) TestExportingEmptyGenesis() {
	ctx, cellarfeesKeeper := suite.ctx, suite.cellarfeesKeeper
	require := suite.Require()

	cellarfeesKeeper.SetParams(ctx, cellarfeestypesv2.DefaultParams())
	cellarfeesKeeper.SetLastRewardSupplyPeak(ctx, sdk.ZeroInt())

	require.Equal(cellarfeestypesv2.DefaultGenesisState(), cellarfeesKeeper.ExportGenesis(ctx))
}

func (suite *KeeperTestSuite) TestExportingPopulatedGenesis() {
	ctx, cellarfeesKeeper := suite.ctx, suite.cellarfeesKeeper
	require := suite.Require()

	var params cellarfeestypesv2.Params
	params.InitialPriceDecreaseRate = sdk.MustNewDecFromStr("0.01")
	params.PriceDecreaseBlockInterval = 10
	params.RewardEmissionPeriod = 600
	params.AuctionInterval = 1000
	params.AuctionThresholdUsdValue = sdk.NewDec(1000000)
	params.ProceedsPortion = sdk.MustNewDecFromStr("0.5")
	cellarfeesKeeper.SetParams(ctx, params)
	peak := sdk.NewInt(1337)
	cellarfeesKeeper.SetLastRewardSupplyPeak(ctx, peak)

	export := cellarfeesKeeper.ExportGenesis(ctx)
	require.Equal(params, export.Params)
	require.Equal(peak, export.LastRewardSupplyPeak)
}
