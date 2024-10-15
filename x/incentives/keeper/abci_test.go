package keeper

import (
	abci "github.com/cometbft/cometbft/abci/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	distributionTypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	"github.com/peggyjv/sommelier/v7/app/params"
	incentivesTypes "github.com/peggyjv/sommelier/v7/x/incentives/types"
)

func (suite *KeeperTestSuite) TestEndBlockerIncentivesDisabledDoesNothing() {
	ctx, incentivesKeeper := suite.ctx, suite.incentivesKeeper
	require := suite.Require()

	incentivesParams := incentivesTypes.DefaultParams()
	incentivesKeeper.SetParams(ctx, incentivesParams)

	// By not mocking any other calls, the test will panic and fail if an unmocked keeper function is called,
	// implying that the function isn't exiting early as designed.

	require.NotPanics(func() { incentivesKeeper.BeginBlocker(ctx, abci.RequestBeginBlock{}) })
	require.NotPanics(func() { incentivesKeeper.EndBlocker(ctx) })

	incentivesParams.DistributionPerBlock = sdk.NewCoin(params.BaseCoinUnit, sdk.OneInt())
	incentivesKeeper.SetParams(ctx, incentivesParams)

	require.NotPanics(func() { incentivesKeeper.BeginBlocker(ctx, abci.RequestBeginBlock{}) })
	require.NotPanics(func() { incentivesKeeper.EndBlocker(ctx) })

	incentivesParams.DistributionPerBlock = sdk.NewCoin(params.BaseCoinUnit, sdk.ZeroInt())
	incentivesParams.IncentivesCutoffHeight = 1500
	incentivesKeeper.SetParams(ctx, incentivesParams)

	require.NotPanics(func() { incentivesKeeper.BeginBlocker(ctx, abci.RequestBeginBlock{}) })
	require.NotPanics(func() { incentivesKeeper.EndBlocker(ctx) })
}

func (suite *KeeperTestSuite) TestEndBlockerInsufficientCommunityPoolBalance() {
	ctx, incentivesKeeper := suite.ctx, suite.incentivesKeeper
	require := suite.Require()
	distributionPerBlock := sdk.NewCoin(params.BaseCoinUnit, sdk.NewInt(1000))

	incentivesParams := incentivesTypes.DefaultParams()
	incentivesParams.DistributionPerBlock = distributionPerBlock
	incentivesParams.IncentivesCutoffHeight = 100
	incentivesKeeper.SetParams(ctx, incentivesParams)

	// mocks
	pool := distributionTypes.FeePool{
		CommunityPool: sdk.NewDecCoins(sdk.NewDecCoin(params.BaseCoinUnit, sdk.NewInt(999))),
	}
	suite.distributionKeeper.EXPECT().GetFeePool(ctx).Return(pool)

	// By not mocking the bank SendModuleToModule call, the test will panic and fail if the community pool balance
	// check branch isn't taken as intended.

	require.NotPanics(func() { incentivesKeeper.BeginBlocker(ctx, abci.RequestBeginBlock{}) })
	require.NotPanics(func() { incentivesKeeper.EndBlocker(ctx) })
}
