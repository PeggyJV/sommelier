package keeper

import (
	"github.com/peggyjv/sommelier/v7/app/params"
	"github.com/peggyjv/sommelier/v7/x/incentives/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	mintTypes "github.com/cosmos/cosmos-sdk/x/mint/types"
)

func (suite *KeeperTestSuite) TestQueriesHappyPath() {
	ctx, incentivesKeeper := suite.ctx, suite.incentivesKeeper
	require := suite.Require()
	blocksPerYear := 365 * 6
	bondedRatio := sdk.MustNewDecFromStr("0.2")
	stakingTotalSupply := sdk.NewInt(100_000_000_000)

	incentivesParams := types.DefaultParams()
	incentivesParams.DistributionPerBlock = sdk.NewCoin(params.BaseCoinUnit, sdk.NewInt(2_000_000))
	incentivesParams.IncentivesCutoffHeight = 1500
	incentivesKeeper.SetParams(ctx, incentivesParams)

	// mocks
	suite.mintKeeper.EXPECT().GetParams(ctx).Return(mintTypes.Params{
		BlocksPerYear: uint64(blocksPerYear),
		MintDenom:     params.BaseCoinUnit,
	})
	suite.mintKeeper.EXPECT().BondedRatio(ctx).Return(bondedRatio)
	suite.mintKeeper.EXPECT().StakingTokenSupply(ctx).Return(stakingTotalSupply)

	paramsResult, err := incentivesKeeper.QueryParams(sdk.WrapSDKContext(ctx), &types.QueryParamsRequest{})
	require.Nil(err)
	require.Equal(incentivesParams, paramsResult.Params)

	APYResult, err := incentivesKeeper.QueryAPY(sdk.WrapSDKContext(ctx), &types.QueryAPYRequest{})
	require.Nil(err)
	expectedAPY := incentivesParams.DistributionPerBlock.Amount.Mul(sdk.NewInt(int64(blocksPerYear))).ToDec().Quo(stakingTotalSupply.ToDec()).Quo(bondedRatio)
	require.Equal(expectedAPY.String(), APYResult.Apy)
}
