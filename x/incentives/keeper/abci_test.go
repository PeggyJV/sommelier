package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	distributionTypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	"github.com/golang/mock/gomock"
	inventivesTypes "github.com/peggyjv/sommelier/v4/x/incentives/types"
)

func (suite *KeeperTestSuite) TestBeginBlockerZeroRewardsBalance() {
	ctx, incentivesKeeper := suite.ctx, suite.incentivesKeeper
	require := suite.Require()
	denom := "usomm"
	distributionPerBlock := sdk.NewCoin(denom, sdk.OneInt())

	params := inventivesTypes.DefaultParams()
	params.DistributionPerBlock = distributionPerBlock
	incentivesKeeper.SetParams(ctx, params)

	// mocks
	pool := distributionTypes.FeePool{
		CommunityPool: sdk.NewDecCoins(sdk.NewDecCoin(denom, sdk.NewInt(10_000_000))),
	}
	suite.distributionKeepr.EXPECT().GetFeePool(ctx).Return(pool)
	suite.bankKeeper.EXPECT().SendCoinsFromModuleToModule(ctx, gomock.Any(), gomock.Any(), gomock.Any()).Times(1)

	expectedAmount := pool.CommunityPool.AmountOf(denom).Sub(distributionPerBlock.Amount.ToDec())
	expectedPool := distributionTypes.FeePool{
		CommunityPool: sdk.NewDecCoins(sdk.NewDecCoinFromDec(denom, expectedAmount)),
	}
	suite.distributionKeepr.EXPECT().SetFeePool(ctx, expectedPool).Times(1)

	// // Note BeginBlocker is only run once for completeness, since it has no code in it
	require.NotPanics(func() { incentivesKeeper.BeginBlocker(ctx) })
	require.NotPanics(func() { incentivesKeeper.EndBlocker(ctx) })
}
