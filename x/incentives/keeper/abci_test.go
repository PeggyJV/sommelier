package keeper

import (
	inventivesTypes "github.com/peggyjv/sommelier/v4/x/incentives/types"
)

func (suite *KeeperTestSuite) TestBeginBlockerZeroRewardsBalance() {
	ctx, incentivesKeeper := suite.ctx, suite.incentivesKeeper
	require := suite.Require()

	params := inventivesTypes.DefaultParams()
	incentivesKeeper.SetParams(ctx, params)

	suite.bankKeeper.SetBalance(ctx, feesAccount.GetAddress(), sdk.NewCoin(appParams.BaseCoinUnit, sdk.ZeroInt()))

	// // mocks
	// suite.bankKeeper.EXPECT().GetBalance(ctx, feesAccount.GetAddress(), appParams.BaseCoinUnit).Return(sdk.NewCoin(appParams.BaseCoinUnit, sdk.ZeroInt()))
	// suite.bankKeeper.EXPECT().SendCoinsFromModuleToModule(ctx, gomock.Any(), gomock.Any(), gomock.Any()).Times(0)

	// // Note EndBlocker is only run once for completeness, since it has no code in it
	require.NotPanics(func() { incentivesKeeper.BeginBlocker(ctx) })
	require.NotPanics(func() { incentivesKeeper.EndBlocker(ctx) })
}
