package keeper

import (
	abci "github.com/cometbft/cometbft/abci/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	distributionTypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	"github.com/golang/mock/gomock"
	"github.com/peggyjv/sommelier/v9/app/params"
	incentivesTypes "github.com/peggyjv/sommelier/v9/x/incentives/types"
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

func (suite *KeeperTestSuite) TestBeginBlockerIncentivesDisabled() {
	ctx, incentivesKeeper := suite.ctx, suite.incentivesKeeper
	require := suite.Require()

	incentivesParams := incentivesTypes.DefaultParams()
	incentivesParams.ValidatorIncentivesCutoffHeight = 100
	incentivesParams.ValidatorMaxDistributionPerBlock = sdk.NewCoin(params.BaseCoinUnit, sdk.NewInt(1000))
	incentivesKeeper.SetParams(ctx, incentivesParams)

	// Set block height above cutoff
	ctx = ctx.WithBlockHeight(101)

	// By not mocking any other calls, the test will panic and fail if an unmocked keeper function is called,
	// implying that the function isn't exiting early as designed.
	require.NotPanics(func() { incentivesKeeper.BeginBlocker(ctx, abci.RequestBeginBlock{}) })

	incentivesParams.ValidatorIncentivesCutoffHeight = 200
	incentivesParams.ValidatorMaxDistributionPerBlock = sdk.NewCoin(params.BaseCoinUnit, sdk.ZeroInt())
	incentivesKeeper.SetParams(ctx, incentivesParams)

	require.NotPanics(func() { incentivesKeeper.BeginBlocker(ctx, abci.RequestBeginBlock{}) })
}

func (suite *KeeperTestSuite) TestBeginBlockerInsufficientCommunityPoolBalance() {
	ctx, incentivesKeeper := suite.ctx, suite.incentivesKeeper
	require := suite.Require()

	incentivesParams := incentivesTypes.DefaultParams()
	incentivesParams.ValidatorIncentivesCutoffHeight = 100
	incentivesParams.ValidatorMaxDistributionPerBlock = sdk.NewCoin(params.BaseCoinUnit, sdk.NewInt(1000))
	incentivesKeeper.SetParams(ctx, incentivesParams)

	// Set block height below cutoff
	ctx = ctx.WithBlockHeight(99)

	// Mock insufficient community pool balance
	pool := distributionTypes.FeePool{
		CommunityPool: sdk.NewDecCoins(sdk.NewDecCoin(params.BaseCoinUnit, sdk.NewInt(999))),
	}
	suite.distributionKeeper.EXPECT().GetFeePool(ctx).Return(pool)

	// By not mocking any further calls, the test will panic and fail if the community pool balance
	// check branch isn't taken as intended.
	require.NotPanics(func() { incentivesKeeper.BeginBlocker(ctx, abci.RequestBeginBlock{}) })
}

func (suite *KeeperTestSuite) TestBeginBlockerSuccess() {
	ctx, incentivesKeeper := suite.ctx, suite.incentivesKeeper
	require := suite.Require()

	incentivesParams := incentivesTypes.DefaultParams()
	incentivesParams.ValidatorIncentivesCutoffHeight = 100
	incentivesParams.ValidatorMaxDistributionPerBlock = sdk.NewCoin(params.BaseCoinUnit, sdk.NewInt(1000))
	incentivesKeeper.SetParams(ctx, incentivesParams)

	// Set block height below cutoff
	ctx = ctx.WithBlockHeight(99)

	// Mock sufficient community pool balance
	pool := distributionTypes.FeePool{
		CommunityPool: sdk.NewDecCoins(sdk.NewDecCoin(params.BaseCoinUnit, sdk.NewInt(2000))),
	}
	suite.distributionKeeper.EXPECT().GetFeePool(ctx).Return(pool)

	// Mock validators
	validators := suite.getMockValidators()
	validator1, validator2 := validators[0], validators[1]

	consAddr1, err := validator1.GetConsAddr()
	require.NoError(err)
	consAddr2, err := validator2.GetConsAddr()
	require.NoError(err)

	// Mock RequestBeginBlock
	req := abci.RequestBeginBlock{
		LastCommitInfo: abci.CommitInfo{
			Votes: []abci.VoteInfo{
				{Validator: abci.Validator{Address: consAddr1, Power: 10}, SignedLastBlock: true},
				{Validator: abci.Validator{Address: consAddr2, Power: 20}, SignedLastBlock: true},
			},
		},
	}

	// Mock StakingKeeper expectations
	suite.stakingKeeper.EXPECT().ValidatorByConsAddr(ctx, consAddr1).Return(validator1)
	suite.stakingKeeper.EXPECT().ValidatorByConsAddr(ctx, consAddr2).Return(validator2)

	// Mock DistributionKeeper expectations for AllocateTokens
	suite.distributionKeeper.EXPECT().GetValidatorCurrentRewards(ctx, validator1.GetOperator()).Return(distributionTypes.ValidatorCurrentRewards{Rewards: sdk.DecCoins{}})
	suite.distributionKeeper.EXPECT().SetValidatorCurrentRewards(ctx, validator1.GetOperator(), gomock.Any())
	suite.distributionKeeper.EXPECT().GetValidatorOutstandingRewards(ctx, validator1.GetOperator()).Return(distributionTypes.ValidatorOutstandingRewards{Rewards: sdk.DecCoins{}})
	suite.distributionKeeper.EXPECT().SetValidatorOutstandingRewards(ctx, validator1.GetOperator(), gomock.Any())

	suite.distributionKeeper.EXPECT().GetValidatorCurrentRewards(ctx, validator2.GetOperator()).Return(distributionTypes.ValidatorCurrentRewards{Rewards: sdk.DecCoins{}})
	suite.distributionKeeper.EXPECT().SetValidatorCurrentRewards(ctx, validator2.GetOperator(), gomock.Any())
	suite.distributionKeeper.EXPECT().GetValidatorOutstandingRewards(ctx, validator2.GetOperator()).Return(distributionTypes.ValidatorOutstandingRewards{Rewards: sdk.DecCoins{}})
	suite.distributionKeeper.EXPECT().SetValidatorOutstandingRewards(ctx, validator2.GetOperator(), gomock.Any())

	// Mock setting the updated fee pool
	suite.distributionKeeper.EXPECT().SetFeePool(ctx, gomock.Any())

	require.NotPanics(func() { incentivesKeeper.BeginBlocker(ctx, req) })

	// You can add more specific assertions here if needed, such as checking emitted events
}
