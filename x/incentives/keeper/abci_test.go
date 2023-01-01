package keeper

func (suite *KeeperTestSuite) TestBeginBlockerZeroRewardsBalance() {
	// ctx, cellarfeesKeeper := suite.ctx, suite.cellarfeesKeeper
	// require := suite.Require()

	// params := cellarfeesTypes.DefaultParams()
	// cellarfeesKeeper.SetParams(ctx, params)

	// // mocks
	// suite.accountKeeper.EXPECT().GetModuleAccount(ctx, cellarfeesTypes.ModuleName).Return(feesAccount)
	// suite.bankKeeper.EXPECT().GetBalance(ctx, feesAccount.GetAddress(), appParams.BaseCoinUnit).Return(sdk.NewCoin(appParams.BaseCoinUnit, sdk.ZeroInt()))
	// suite.bankKeeper.EXPECT().SendCoinsFromModuleToModule(ctx, gomock.Any(), gomock.Any(), gomock.Any()).Times(0)

	// // Note EndBlocker is only run once for completeness, since it has no code in it
	// require.NotPanics(func() { cellarfeesKeeper.BeginBlocker(ctx) })
	// require.NotPanics(func() { cellarfeesKeeper.EndBlocker(ctx) })
}
