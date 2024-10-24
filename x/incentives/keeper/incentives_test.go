package keeper

import (
	abci "github.com/cometbft/cometbft/abci/types"
	"github.com/cosmos/cosmos-sdk/crypto/keys/ed25519"
	ccrypto "github.com/cosmos/cosmos-sdk/crypto/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	distributiontypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/peggyjv/sommelier/v8/x/incentives/types"
)

var (
	// ConsPrivKeys generate ed25519 ConsPrivKeys to be used for validator operator keys
	ConsPrivKeys = []ccrypto.PrivKey{
		ed25519.GenPrivKey(),
		ed25519.GenPrivKey(),
		ed25519.GenPrivKey(),
		ed25519.GenPrivKey(),
		ed25519.GenPrivKey(),
	}

	// ConsPubKeys holds the consensus public keys to be used for validator operator keys
	ConsPubKeys = []ccrypto.PubKey{
		ConsPrivKeys[0].PubKey(),
		ConsPrivKeys[1].PubKey(),
		ConsPrivKeys[2].PubKey(),
		ConsPrivKeys[3].PubKey(),
		ConsPrivKeys[4].PubKey(),
	}
)

func (suite *KeeperTestSuite) getMockValidators() []*stakingtypes.Validator {
	validator1, err := stakingtypes.NewValidator(sdk.ValAddress([]byte("val1val1val1val1val1")), ConsPubKeys[0], stakingtypes.Description{})
	suite.Require().NoError(err)
	validator2, err := stakingtypes.NewValidator(sdk.ValAddress([]byte("val2val2val2val2val2")), ConsPubKeys[1], stakingtypes.Description{})
	suite.Require().NoError(err)
	validator3, err := stakingtypes.NewValidator(sdk.ValAddress([]byte("val3val3val3val3val3")), ConsPubKeys[2], stakingtypes.Description{})
	suite.Require().NoError(err)
	return []*stakingtypes.Validator{&validator1, &validator2, &validator3}
}

func (suite *KeeperTestSuite) TestGetValidatorInfos() {
	ctx, incentivesKeeper := suite.ctx, suite.incentivesKeeper

	// Create mock validators
	validators := suite.getMockValidators()
	validator1, validator2, validator3 := validators[0], validators[1], validators[2]

	consAddr1, err := validator1.GetConsAddr()
	suite.Require().NoError(err)
	consAddr2, err := validator2.GetConsAddr()
	suite.Require().NoError(err)
	consAddr3, err := validator3.GetConsAddr()
	suite.Require().NoError(err)

	// Create mock RequestBeginBlock
	req := abci.RequestBeginBlock{
		LastCommitInfo: abci.CommitInfo{
			Votes: []abci.VoteInfo{
				{
					Validator:       abci.Validator{Address: consAddr1, Power: 10},
					SignedLastBlock: true,
				},
				{
					Validator:       abci.Validator{Address: consAddr2, Power: 20},
					SignedLastBlock: true,
				},
				{
					Validator:       abci.Validator{Address: consAddr3, Power: 30},
					SignedLastBlock: false,
				},
			},
		},
	}

	// Set up expectations for the mock StakingKeeper
	suite.stakingKeeper.EXPECT().ValidatorByConsAddr(ctx, consAddr1).Return(validator1)
	suite.stakingKeeper.EXPECT().ValidatorByConsAddr(ctx, consAddr2).Return(validator2)

	// Call the function being tested
	validatorInfos := incentivesKeeper.getValidatorInfos(ctx, req)

	// Assert the results
	suite.Require().Len(validatorInfos, 2)
	suite.Require().Equal(validator1, validatorInfos[0].Validator)
	suite.Require().Equal(int64(10), validatorInfos[0].Power)
	suite.Require().Equal(validator2, validatorInfos[1].Validator)
	suite.Require().Equal(int64(20), validatorInfos[1].Power)
}

func (suite *KeeperTestSuite) TestGetValidatorInfosNoSigners() {
	ctx, incentivesKeeper := suite.ctx, suite.incentivesKeeper

	// Create mock RequestBeginBlock with no signers
	req := abci.RequestBeginBlock{
		LastCommitInfo: abci.CommitInfo{
			Votes: []abci.VoteInfo{
				{
					Validator:       abci.Validator{Address: []byte("val1"), Power: 10},
					SignedLastBlock: false,
				},
				{
					Validator:       abci.Validator{Address: []byte("val2"), Power: 20},
					SignedLastBlock: false,
				},
			},
		},
	}

	// Call the function being tested
	validatorInfos := incentivesKeeper.getValidatorInfos(ctx, req)

	// Assert the results
	suite.Require().Len(validatorInfos, 0)
}

func (suite *KeeperTestSuite) TestSortValidatorInfosByPower() {
	// Create a slice of ValidatorInfo with unsorted power
	valInfos := []ValidatorInfo{
		{Power: 30},
		{Power: 10},
		{Power: 50},
		{Power: 20},
		{Power: 40},
	}

	// Sort the validator infos
	sortedValInfos := sortValidatorInfosByPower(valInfos)

	// Assert the results
	suite.Require().Len(sortedValInfos, 5)
	suite.Require().Equal(int64(50), sortedValInfos[0].Power)
	suite.Require().Equal(int64(40), sortedValInfos[1].Power)
	suite.Require().Equal(int64(30), sortedValInfos[2].Power)
	suite.Require().Equal(int64(20), sortedValInfos[3].Power)
	suite.Require().Equal(int64(10), sortedValInfos[4].Power)
}

func (suite *KeeperTestSuite) TestTruncateVoters() {
	// Create a slice of ValidatorInfo
	valInfos := []ValidatorInfo{
		{Power: 30},
		{Power: 10},
		{Power: 50},
		{Power: 20},
		{Power: 40},
	}

	// Get the truncated voters
	truncatedVoters := truncateVoters(valInfos, 3)

	// Assert the results
	suite.Require().Len(truncatedVoters, 3)
	suite.Require().Equal(int64(30), truncatedVoters[0].Power)
	suite.Require().Equal(int64(10), truncatedVoters[1].Power)
	suite.Require().Equal(int64(50), truncatedVoters[2].Power)
}

func (suite *KeeperTestSuite) TestSortValidatorInfosByPowerEmptySlice() {
	// Create an empty slice of ValidatorInfo
	var valInfos []ValidatorInfo

	// Sort the validator infos
	sortedValInfos := sortValidatorInfosByPower(valInfos)

	// Assert the results
	suite.Require().Len(sortedValInfos, 0)
}

func (suite *KeeperTestSuite) TestGetTotalPower() {
	// Create a slice of ValidatorInfo
	valInfos := []ValidatorInfo{
		{Power: 30},
		{Power: 10},
		{Power: 50},
		{Power: 20},
		{Power: 40},
	}

	// Get the total power
	totalPower := getTotalPower(&valInfos)

	// Assert the result
	suite.Require().Equal(int64(150), totalPower)
}

func (suite *KeeperTestSuite) TestGetTotalPowerEmptySlice() {
	// Create an empty slice of ValidatorInfo
	var valInfos []ValidatorInfo

	// Get the total power
	totalPower := getTotalPower(&valInfos)

	// Assert the result
	suite.Require().Equal(int64(0), totalPower)
}

func (suite *KeeperTestSuite) TestAllocateTokensToValidator() {
	ctx, incentivesKeeper := suite.ctx, suite.incentivesKeeper

	// Create a mock validator
	valAddr := sdk.ValAddress([]byte("validatorvalidatorva"))
	validator, err := stakingtypes.NewValidator(valAddr, ConsPubKeys[0], stakingtypes.Description{})
	suite.Require().NoError(err)

	// Create mock tokens to allocate
	tokens := sdk.NewDecCoins(sdk.NewDecCoin("usom", sdk.NewInt(100)))

	// Set up expectations for the mock DistributionKeeper
	currentRewards := distributiontypes.ValidatorCurrentRewards{Rewards: sdk.DecCoins{}}
	outstandingRewards := distributiontypes.ValidatorOutstandingRewards{Rewards: sdk.DecCoins{}}

	suite.distributionKeeper.EXPECT().
		GetValidatorCurrentRewards(ctx, valAddr).
		Return(currentRewards)
	suite.distributionKeeper.EXPECT().
		SetValidatorCurrentRewards(ctx, valAddr, distributiontypes.ValidatorCurrentRewards{Rewards: tokens})
	suite.distributionKeeper.EXPECT().
		GetValidatorOutstandingRewards(ctx, valAddr).
		Return(outstandingRewards)
	suite.distributionKeeper.EXPECT().
		SetValidatorOutstandingRewards(ctx, valAddr, distributiontypes.ValidatorOutstandingRewards{Rewards: tokens})

	// Call the function being tested
	incentivesKeeper.AllocateTokensToValidator(ctx, validator, tokens)

	// Verify that the event was emitted
	events := ctx.EventManager().Events()
	suite.Require().Len(events, 1)
	event := events[0]
	suite.Require().Equal(types.EventTypeValidatorIncentivesReward, event.Type)
	suite.Require().Len(event.Attributes, 2)
	suite.Require().Equal(sdk.AttributeKeyAmount, event.Attributes[0].Key)
	suite.Require().Equal(tokens.String(), event.Attributes[0].Value)
	suite.Require().Equal(types.AttributeKeyValidator, event.Attributes[1].Key)
	suite.Require().Equal(valAddr.String(), event.Attributes[1].Value)
}

func (suite *KeeperTestSuite) TestAllocateTokensToValidatorWithExistingRewards() {
	ctx, incentivesKeeper := suite.ctx, suite.incentivesKeeper

	// Create a mock validator
	valAddr := sdk.ValAddress([]byte("validatorvalidatorva"))
	validator, err := stakingtypes.NewValidator(valAddr, ConsPubKeys[0], stakingtypes.Description{})
	suite.Require().NoError(err)

	// Create mock tokens to allocate
	existingRewards := sdk.NewDecCoins(sdk.NewDecCoin("usom", sdk.NewInt(50)))
	newTokens := sdk.NewDecCoins(sdk.NewDecCoin("usom", sdk.NewInt(100)))
	expectedTotalRewards := existingRewards.Add(newTokens...)

	// Set up expectations for the mock DistributionKeeper
	currentRewards := distributiontypes.ValidatorCurrentRewards{Rewards: existingRewards}
	outstandingRewards := distributiontypes.ValidatorOutstandingRewards{Rewards: existingRewards}

	suite.distributionKeeper.EXPECT().
		GetValidatorCurrentRewards(ctx, valAddr).
		Return(currentRewards)
	suite.distributionKeeper.EXPECT().
		SetValidatorCurrentRewards(ctx, valAddr, distributiontypes.ValidatorCurrentRewards{Rewards: expectedTotalRewards})
	suite.distributionKeeper.EXPECT().
		GetValidatorOutstandingRewards(ctx, valAddr).
		Return(outstandingRewards)
	suite.distributionKeeper.EXPECT().
		SetValidatorOutstandingRewards(ctx, valAddr, distributiontypes.ValidatorOutstandingRewards{Rewards: expectedTotalRewards})

	// Call the function being tested
	incentivesKeeper.AllocateTokensToValidator(ctx, validator, newTokens)

	// Verify that the event was emitted
	events := ctx.EventManager().Events()
	suite.Require().Len(events, 1)
	event := events[0]
	suite.Require().Equal(types.EventTypeValidatorIncentivesReward, event.Type)
	suite.Require().Len(event.Attributes, 2)
	suite.Require().Equal(sdk.AttributeKeyAmount, event.Attributes[0].Key)
	suite.Require().Equal(newTokens.String(), event.Attributes[0].Value)
	suite.Require().Equal(types.AttributeKeyValidator, event.Attributes[1].Key)
	suite.Require().Equal(valAddr.String(), event.Attributes[1].Value)
}

func (suite *KeeperTestSuite) TestAllocateTokens() {
	ctx, incentivesKeeper := suite.ctx, suite.incentivesKeeper

	// Create mock validators
	validators := suite.getMockValidators()
	validator1, validator2, validator3 := validators[0], validators[1], validators[2]

	// Set up qualifying voters
	qualifyingVoters := []ValidatorInfo{
		{Validator: validator1, Power: 30},
		{Validator: validator2, Power: 20},
		{Validator: validator3, Power: 10},
	}

	totalPreviousPower := int64(60)
	totalDistribution := sdk.NewDecCoins(sdk.NewDecCoin("usom", sdk.NewInt(100)))
	maxFraction := sdk.NewDecWithPrec(5, 1) // 0.5

	// Set up expectations for the mock DistributionKeeper
	totalExpectedRewards := sdk.NewDecCoins()
	for _, voter := range qualifyingVoters {
		powerFraction := sdk.NewDecFromInt(sdk.NewInt(voter.Power)).QuoInt64(totalPreviousPower)
		expectedReward := totalDistribution.MulDecTruncate(powerFraction)
		if powerFraction.GT(maxFraction) {
			expectedReward = totalDistribution.MulDecTruncate(maxFraction)
		}

		totalExpectedRewards = totalExpectedRewards.Add(expectedReward...)

		suite.distributionKeeper.EXPECT().
			GetValidatorCurrentRewards(ctx, voter.Validator.GetOperator()).
			Return(distributiontypes.ValidatorCurrentRewards{Rewards: sdk.DecCoins{}})
		suite.distributionKeeper.EXPECT().
			SetValidatorCurrentRewards(ctx, voter.Validator.GetOperator(), distributiontypes.ValidatorCurrentRewards{Rewards: expectedReward})
		suite.distributionKeeper.EXPECT().
			GetValidatorOutstandingRewards(ctx, voter.Validator.GetOperator()).
			Return(distributiontypes.ValidatorOutstandingRewards{Rewards: sdk.DecCoins{}})
		suite.distributionKeeper.EXPECT().
			SetValidatorOutstandingRewards(ctx, voter.Validator.GetOperator(), distributiontypes.ValidatorOutstandingRewards{Rewards: expectedReward})
	}

	// Call the function being tested
	remaining := incentivesKeeper.AllocateTokens(ctx, totalPreviousPower, totalDistribution, qualifyingVoters, maxFraction)

	// Verify that the sum of remaining and distributed rewards equals totalDistribution
	totalAllocated := remaining.Add(totalExpectedRewards...)
	suite.Require().Equal(totalDistribution, totalAllocated, "Sum of remaining and distributed rewards should equal total distribution")

	// Verify that events were emitted
	totalEvents := len(qualifyingVoters) + 1 // One event for each validator plus one for the total distribution
	events := ctx.EventManager().Events()
	suite.Require().Len(events, totalEvents)
	for i, event := range events {
		if i == totalEvents-1 {
			suite.Require().Equal(types.EventTypeTotalValidatorIncentivesRewards, event.Type)
			suite.Require().Equal(totalDistribution.Sub(remaining).String(), event.Attributes[0].Value)
			continue
		} else {
			suite.Require().Equal(types.EventTypeValidatorIncentivesReward, event.Type)
			suite.Require().Len(event.Attributes, 2)
			suite.Require().Equal(sdk.AttributeKeyAmount, event.Attributes[0].Key)
			suite.Require().Equal(types.AttributeKeyValidator, event.Attributes[1].Key)
			suite.Require().Equal(qualifyingVoters[i].Validator.GetOperator().String(), event.Attributes[1].Value)
		}
	}
}

func (suite *KeeperTestSuite) TestAllocateTokensNoQualifyingVoters() {
	ctx, incentivesKeeper := suite.ctx, suite.incentivesKeeper

	totalPreviousPower := int64(100)
	totalDistribution := sdk.NewDecCoins(sdk.NewDecCoin("usom", sdk.NewInt(100)))
	maxFraction := sdk.NewDecWithPrec(5, 1) // 0.5

	// Call the function being tested with empty qualifyingVoters
	remaining := incentivesKeeper.AllocateTokens(ctx, totalPreviousPower, totalDistribution, []ValidatorInfo{}, maxFraction)

	// Verify that all tokens remain unallocated
	suite.Require().Equal(totalDistribution, remaining, "All tokens should remain unallocated when there are no qualifying voters")

	// Verify that no events were emitted
	events := ctx.EventManager().Events()
	suite.Require().Len(events, 0, "No events should be emitted when there are no qualifying voters")
}
