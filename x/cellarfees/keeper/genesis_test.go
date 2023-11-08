package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	cellarfeesTypes "github.com/peggyjv/sommelier/v7/x/cellarfees/types"
)

func (suite *KeeperTestSuite) TestImportingEmptyGenesis() {
	ctx, cellarfeesKeeper := suite.ctx, suite.cellarfeesKeeper
	require := suite.Require()

	testGenesis := cellarfeesTypes.GenesisState{}

	// Canary to make sure validate basic is being run
	require.Panics(func() { cellarfeesKeeper.InitGenesis(ctx, testGenesis) })

	testGenesis = cellarfeesTypes.DefaultGenesisState()
	require.NotPanics(func() {
		suite.accountKeeper.EXPECT().GetModuleAccount(ctx, feesAccount.GetName()).Return(feesAccount)
		cellarfeesKeeper.InitGenesis(ctx, testGenesis)
	})

	require.Len(cellarfeesKeeper.GetFeeAccrualCounters(ctx).Counters, 0)
	require.Zero(cellarfeesKeeper.GetLastRewardSupplyPeak(ctx).Int64())
	require.Equal(cellarfeesKeeper.GetParams(ctx), cellarfeesTypes.DefaultParams())
}

func (suite *KeeperTestSuite) TestImportingPopulatedGenesis() {
	ctx, cellarfeesKeeper := suite.ctx, suite.cellarfeesKeeper
	require := suite.Require()

	testGenesis := cellarfeesTypes.GenesisState{}

	// Canary to make sure validate basic is being run
	require.Panics(func() { cellarfeesKeeper.InitGenesis(ctx, testGenesis) })

	testGenesis.FeeAccrualCounters = cellarfeesTypes.FeeAccrualCounters{
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
	testGenesis.LastRewardSupplyPeak = sdk.NewInt(1337)
	testGenesis.Params.FeeAccrualAuctionThreshold = 2
	testGenesis.Params.InitialPriceDecreaseRate = sdk.MustNewDecFromStr("0.01")
	testGenesis.Params.PriceDecreaseBlockInterval = 10
	testGenesis.Params.RewardEmissionPeriod = 600

	require.NotPanics(func() {
		suite.accountKeeper.EXPECT().GetModuleAccount(ctx, feesAccount.GetName()).Return(feesAccount)
		cellarfeesKeeper.InitGenesis(ctx, testGenesis)
	})

	require.Equal(testGenesis.FeeAccrualCounters, cellarfeesKeeper.GetFeeAccrualCounters(ctx))
	require.Equal(testGenesis.LastRewardSupplyPeak, cellarfeesKeeper.GetLastRewardSupplyPeak(ctx))
	require.Equal(testGenesis.Params, cellarfeesKeeper.GetParams(ctx))
}

func (suite *KeeperTestSuite) TestExportingEmptyGenesis() {
	ctx, cellarfeesKeeper := suite.ctx, suite.cellarfeesKeeper
	require := suite.Require()

	cellarfeesKeeper.SetParams(ctx, cellarfeesTypes.DefaultParams())
	cellarfeesKeeper.SetFeeAccrualCounters(ctx, cellarfeesTypes.DefaultFeeAccrualCounters())
	cellarfeesKeeper.SetLastRewardSupplyPeak(ctx, sdk.ZeroInt())

	require.Equal(cellarfeesTypes.DefaultGenesisState(), cellarfeesKeeper.ExportGenesis(ctx))
}

func (suite *KeeperTestSuite) TestExportingPopulatedGenesis() {
	ctx, cellarfeesKeeper := suite.ctx, suite.cellarfeesKeeper
	require := suite.Require()

	var params cellarfeesTypes.Params
	params.FeeAccrualAuctionThreshold = 2
	params.InitialPriceDecreaseRate = sdk.MustNewDecFromStr("0.01")
	params.PriceDecreaseBlockInterval = 10
	params.RewardEmissionPeriod = 600
	cellarfeesKeeper.SetParams(ctx, params)
	counters := cellarfeesTypes.FeeAccrualCounters{
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
	cellarfeesKeeper.SetFeeAccrualCounters(ctx, counters)
	peak := sdk.NewInt(1337)
	cellarfeesKeeper.SetLastRewardSupplyPeak(ctx, peak)

	export := cellarfeesKeeper.ExportGenesis(ctx)
	require.Equal(params, export.Params)
	require.Equal(counters, export.FeeAccrualCounters)
	require.Equal(peak, export.LastRewardSupplyPeak)
}
