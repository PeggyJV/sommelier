package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	distributionTypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	"github.com/golang/mock/gomock"
	"github.com/peggyjv/sommelier/v4/app/params"
	inventivesTypes "github.com/peggyjv/sommelier/v4/x/incentives/types"
)

func (suite *KeeperTestSuite) TestBeginBlockerZeroRewardsBalance() {
	ctx, incentivesKeeper := suite.ctx, suite.incentivesKeeper
	require := suite.Require()
	distributionPerBlock := sdk.NewCoin(params.BaseCoinUnit, sdk.OneInt())

	incentivesParams := inventivesTypes.DefaultParams()
	incentivesParams.DistributionPerBlock = distributionPerBlock
	incentivesKeeper.SetParams(ctx, incentivesParams)

	// mocks
	pool := distributionTypes.FeePool{
		CommunityPool: sdk.NewDecCoins(sdk.NewDecCoin(params.BaseCoinUnit, sdk.NewInt(10_000_000))),
	}
	suite.distributionKeepr.EXPECT().GetFeePool(ctx).Return(pool)
	suite.bankKeeper.EXPECT().SendCoinsFromModuleToModule(ctx, gomock.Any(), gomock.Any(), gomock.Any()).Times(1)

	expectedAmount := pool.CommunityPool.AmountOf(params.BaseCoinUnit).Sub(distributionPerBlock.Amount.ToDec())
	expectedPool := distributionTypes.FeePool{
		CommunityPool: sdk.NewDecCoins(sdk.NewDecCoinFromDec(params.BaseCoinUnit, expectedAmount)),
	}
	suite.distributionKeepr.EXPECT().SetFeePool(ctx, expectedPool).Times(1)

	// // Note BeginBlocker is only run once for completeness, since it has no code in it
	require.NotPanics(func() { incentivesKeeper.BeginBlocker(ctx) })
	require.NotPanics(func() { incentivesKeeper.EndBlocker(ctx) })
}
