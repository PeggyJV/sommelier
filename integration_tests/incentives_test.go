package integration_tests

import (
	"context"
	"fmt"
	"time"

	"github.com/cosmos/cosmos-sdk/client"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/bech32"
	disttypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	govtypesv1beta1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1beta1"
	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"
	paramsproposal "github.com/cosmos/cosmos-sdk/x/params/types/proposal"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/peggyjv/sommelier/v7/app/params"
	incentivestypes "github.com/peggyjv/sommelier/v7/x/incentives/types"
)

func (s *IntegrationTestSuite) TestIncentives() {
	s.Run("Bring up chain, observe incentives distribution after param proposal executes", func() {
		ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
		defer cancel()

		val := s.chain.validators[0]
		s.T().Logf("validator %s", val.address().String())
		kb, err := val.keyring()
		s.Require().NoError(err)

		_, bytes, err := bech32.DecodeAndConvert(val.address().String())
		s.Require().NoError(err)
		valOperatorAddress, err := bech32.ConvertAndEncode("sommvaloper", bytes)
		s.Require().NoError(err)

		clientCtx, err := s.chain.clientContext("tcp://localhost:26657", &kb, "val", val.address())
		s.Require().NoError(err)

		incentivesQueryClient := incentivestypes.NewQueryClient(clientCtx)
		distQueryClient := disttypes.NewQueryClient(clientCtx)
		mintQueryClient := minttypes.NewQueryClient(clientCtx)
		stakingQueryClient := stakingtypes.NewQueryClient(clientCtx)

		s.T().Log("verifying that the base distribution rate is zero")
		incentivesParamsRes, err := incentivesQueryClient.QueryParams(ctx, &incentivestypes.QueryParamsRequest{})
		s.Require().NoError(err)
		s.Require().Equal(uint64(0), incentivesParamsRes.Params.IncentivesCutoffHeight)
		s.Require().Equal(sdk.ZeroInt(), incentivesParamsRes.Params.DistributionPerBlock.Amount)

		s.T().Log("verifying APY query returns zero")
		incentivesAPYRes, err := incentivesQueryClient.QueryAPY(ctx, &incentivestypes.QueryAPYRequest{})
		s.Require().NoError(err)
		initialAPY, err := sdk.NewDecFromStr(incentivesAPYRes.Apy)
		s.Require().NoError(err)
		s.Require().True(initialAPY.IsZero())

		// this is looped because for some reason running the test too quickly after the validator
		// containers are launched results in the validator rewards pool result increasing before
		// reaching a steady state
		s.T().Log("verifying validator rewards are not increasing")
		s.Require().Eventually(func() bool {
			beforeAmount := s.queryValidatorRewards(ctx, valOperatorAddress, distQueryClient)
			time.Sleep(time.Second * 12)
			afterAmount := s.queryValidatorRewards(ctx, valOperatorAddress, distQueryClient)

			return beforeAmount.Equal(afterAmount)
		}, 120*time.Second, 12*time.Second)

		s.T().Log("submitting proposal to enable incentives")
		proposer := s.chain.proposer
		proposerCtx, err := s.chain.clientContext("tcp://localhost:26657", proposer.keyring, "proposer", proposer.address())
		s.Require().NoError(err)
		orch := s.chain.orchestrators[0]
		orchClientCtx, err := s.chain.clientContext("tcp://localhost:26657", orch.keyring, "orch", orch.address())
		s.Require().NoError(err)

		// Collin: test takes about ~110 blocks to reach the cutoff check at the end on my machine.
		// Hopefully this is high enough to avoid a false negative in CI.
		cutoffHeight := 200
		expectedUsommAmount := int64(2_000_000)
		expectedDistributionPerBlock := sdk.NewCoin(params.BaseCoinUnit, sdk.NewInt(expectedUsommAmount))
		proposal := paramsproposal.ParameterChangeProposal{
			Title:       "enable incentives",
			Description: "enables incentives",
			Changes: []paramsproposal.ParamChange{
				{
					Subspace: "incentives",
					Key:      "IncentivesCutoffHeight",
					Value:    fmt.Sprintf("\"%d\"", cutoffHeight),
				},
				{
					Subspace: "incentives",
					Key:      "DistributionPerBlock",
					Value:    fmt.Sprintf("{\"denom\":\"%s\",\"amount\":\"%d\"}", params.BaseCoinUnit, expectedUsommAmount),
				},
			},
		}

		proposalMsg, err := govtypesv1beta1.NewMsgSubmitProposal(
			&proposal,
			sdk.Coins{
				{
					Denom:  testDenom,
					Amount: stakeAmount.Quo(sdk.NewInt(2)),
				},
			},
			proposer.address(),
		)
		s.Require().NoError(err, "Unable to create governance proposal")
		submitProposalResponse, err := s.chain.sendMsgs(*proposerCtx, proposalMsg)
		s.Require().NoError(err)
		s.Require().Zero(submitProposalResponse.Code, "raw log: %s", submitProposalResponse.RawLog)

		s.T().Log("check proposal was submitted correctly")
		govQueryClient := govtypesv1beta1.NewQueryClient(orchClientCtx)

		s.Require().Eventually(func() bool {
			proposalsQueryResponse, err := govQueryClient.Proposals(context.Background(), &govtypesv1beta1.QueryProposalsRequest{})
			if err != nil {
				s.T().Logf("error querying proposals: %e", err)
				return false
			}

			s.Require().NotEmpty(proposalsQueryResponse.Proposals)
			s.Require().Equal(uint64(1), proposalsQueryResponse.Proposals[0].ProposalId, "not proposal id 1")
			s.Require().Equal(govtypesv1beta1.StatusVotingPeriod, proposalsQueryResponse.Proposals[0].Status, "proposal not in voting period")

			return true
		}, time.Second*30, time.Second*5, "proposal submission was never found")

		s.T().Log("vote for proposal")
		for _, val := range s.chain.validators {
			kr, err := val.keyring()
			s.Require().NoError(err)
			localClientCtx, err := s.chain.clientContext("tcp://localhost:26657", &kr, "val", val.address())
			s.Require().NoError(err)

			voteMsg := govtypesv1beta1.NewMsgVote(val.address(), 1, govtypesv1beta1.OptionYes)
			voteResponse, err := s.chain.sendMsgs(*localClientCtx, voteMsg)
			s.Require().NoError(err)
			s.Require().Zero(voteResponse.Code, "Vote error: %s", voteResponse.RawLog)
		}

		s.T().Log("waiting for proposal to be approved..")
		s.Require().Eventually(func() bool {
			proposalQueryResponse, _ := govQueryClient.Proposal(context.Background(), &govtypesv1beta1.QueryProposalRequest{ProposalId: 1})
			return govtypesv1beta1.StatusPassed == proposalQueryResponse.Proposal.Status
		}, time.Second*30, time.Second*5, "proposal was never accepted")
		s.T().Log("proposal approved!")

		s.T().Log("verifying parameter was changed")
		incentivesParamsRes, err = incentivesQueryClient.QueryParams(ctx, &incentivestypes.QueryParamsRequest{})
		s.Require().NoError(err)
		s.Require().Equal(incentivesParamsRes.Params.IncentivesCutoffHeight, uint64(cutoffHeight))

		s.T().Log("verifying that the base distribution rate is the expected value")
		beforeAmount, beforeHeight := s.getRewardAmountAndHeight(ctx, distQueryClient, valOperatorAddress, clientCtx)
		time.Sleep(time.Second * 12)
		afterAmount, afterHeight := s.getRewardAmountAndHeight(ctx, distQueryClient, valOperatorAddress, clientCtx)
		// Assumes each validator has equally weight bonding power
		s.T().Logf("blocks %d-%d", beforeHeight, afterHeight)
		actualDistributionPerBlock := (afterAmount.Sub(beforeAmount)).Quo(sdk.NewDec(afterHeight - beforeHeight)).Mul(sdk.NewDec(int64(len(s.chain.validators))))
		s.T().Logf("before: %s, after: %s, blocks %d-%d", beforeAmount.String(), afterAmount.String(), beforeHeight, afterHeight)
		s.Require().True(afterAmount.GT(beforeAmount))
		s.Require().Equal(sdk.NewDecFromInt(expectedDistributionPerBlock.Amount), actualDistributionPerBlock)

		s.T().Log("verifying APY query returns expected value")
		mintParams, err := mintQueryClient.Params(ctx, &minttypes.QueryParamsRequest{})
		s.Require().NoError(err)
		stakingPool, err := stakingQueryClient.Pool(ctx, &stakingtypes.QueryPoolRequest{})
		s.Require().NoError(err)
		tokensTotal := stakingPool.Pool.BondedTokens.Add(stakingPool.Pool.NotBondedTokens)
		// assumes bonded ratio is 100%
		expectedAPY := actualDistributionPerBlock.Mul(sdk.NewDec(int64(mintParams.Params.BlocksPerYear))).Quo(sdk.NewDecFromInt(tokensTotal))
		incentivesAPYRes, err = incentivesQueryClient.QueryAPY(ctx, &incentivestypes.QueryAPYRequest{})
		s.Require().NoError(err)
		APY, err := sdk.NewDecFromStr(incentivesAPYRes.Apy[:10])
		s.Require().NoError(err)
		s.Require().Equal(expectedAPY, APY)

		s.T().Log("verifying incentives end after cutoff")
		s.Require().Eventually(func() bool {
			startingAmount, startingHeight := s.getRewardAmountAndHeight(ctx, distQueryClient, valOperatorAddress, clientCtx)
			// since rewards sent to the distribution module are issued on the next block, we wait for the height after
			// the cutoff to check rewards
			if startingHeight <= int64(cutoffHeight) {
				return false
			}

			time.Sleep(time.Second * 5)
			endingAmount, _ := s.getRewardAmountAndHeight(ctx, distQueryClient, valOperatorAddress, clientCtx)

			if s.Equal(startingAmount, endingAmount) {
				return true
			}

			s.FailNow("rewards per block are nonzero after cutoff height")

			// required to satisfy the eventually block even though it's unreachable
			return false
		}, time.Minute*5, time.Second*10, "incentives did not end after cutoff height")
		s.T().Log("incentives ended!")
	})
}

func (s *IntegrationTestSuite) queryValidatorRewards(ctx context.Context, valOperatorAddress string, distQueryClient disttypes.QueryClient) sdk.Dec {
	rewardsRes, err := distQueryClient.ValidatorOutstandingRewards(ctx, &disttypes.QueryValidatorOutstandingRewardsRequest{
		ValidatorAddress: valOperatorAddress,
	})
	s.Require().NoError(err)
	return rewardsRes.Rewards.Rewards.AmountOf(params.BaseCoinUnit)
}

func (s *IntegrationTestSuite) getCurrentHeight(clientCtx *client.Context) (int64, error) {
	node, err := clientCtx.GetNode()
	if err != nil {
		return 0, err
	}
	status, err := node.Status(context.Background())
	if err != nil {
		return 0, err
	}
	return status.SyncInfo.LatestBlockHeight, nil
}

func (s *IntegrationTestSuite) getRewardAmountAndHeight(ctx context.Context, distQueryClient disttypes.QueryClient, operatorAddress string, clientCtx *client.Context) (sdk.Dec, int64) {
	var amount sdk.Dec
	var height int64

	s.Require().Eventually(func() bool {
		initialHeight, err := s.getCurrentHeight(clientCtx)
		if err != nil {
			return false
		}
		amount = s.queryValidatorRewards(ctx, operatorAddress, distQueryClient)
		height, err = s.getCurrentHeight(clientCtx)
		if err != nil {
			return false
		}
		return initialHeight == height
	}, time.Second*30, time.Second*1, "failed to reliably determine height of reward sample")

	return amount, height
}

func (s *IntegrationTestSuite) TestValidatorIncentives() {
	validator := s.chain.validators[0]
	proposer := s.chain.proposer
	proposerCtx, err := s.chain.clientContext("tcp://localhost:26657", proposer.keyring, "proposer", proposer.address())
	s.Require().NoError(err)
	orch := s.chain.orchestrators[0]
	orchClientCtx, err := s.chain.clientContext("tcp://localhost:26657", orch.keyring, "orch", orch.address())
	s.Require().NoError(err)
	ctx := context.Background()

	s.T().Log("Getting the initial community pool balance")
	queryClient := disttypes.NewQueryClient(proposerCtx)
	queryRes, err := queryClient.CommunityPool(ctx, &disttypes.QueryCommunityPoolRequest{})
	s.Require().NoError(err)
	s.T().Logf("Initial community pool balance: %s", queryRes.Pool.String())
	initialCommunityPool := queryRes.Pool

	// Wait for outstanding rewards to equal 4000000usomm. Current theory is these initial rewards come from
	// the genesis delegation tx fees.
	s.T().Log("Waiting for outstanding rewards to equal 4000000usomm")
	initialRewards := disttypes.ValidatorOutstandingRewards{
		Rewards: sdk.DecCoins{
			sdk.DecCoin{
				Denom:  params.BaseCoinUnit,
				Amount: sdk.NewDec(4000000),
			},
		},
	}
	s.Require().Eventually(func() bool {
		rewards, err := s.getValidatorOutstandingRewards(validator)
		s.Require().NoError(err)
		return rewards.Rewards.AmountOf(params.BaseCoinUnit).Equal(initialRewards.Rewards.AmountOf(params.BaseCoinUnit))
	}, time.Second*30, time.Second*1, "outstanding rewards did not reach 4000000usomm")

	// Submit proposal to enable validator incentives
	s.T().Log("Submitting proposal to enable validator incentives")
	cutoffHeight := 100
	proposal := paramsproposal.ParameterChangeProposal{
		Title:       "Enable validator incentives",
		Description: "Enable validator incentives",
		Changes: []paramsproposal.ParamChange{
			{
				Subspace: "incentives",
				Key:      "ValidatorIncentivesCutoffHeight",
				Value:    fmt.Sprintf("\"%d\"", cutoffHeight),
			},
			{
				Subspace: "incentives",
				Key:      "ValidatorMaxDistributionPerBlock",
				Value:    fmt.Sprintf("{\"denom\":\"%s\",\"amount\":\"%d\"}", params.BaseCoinUnit, 1000000),
			},
		},
	}

	proposalMsg, err := govtypesv1beta1.NewMsgSubmitProposal(
		&proposal,
		sdk.Coins{
			{
				Denom:  testDenom,
				Amount: stakeAmount.Quo(sdk.NewInt(2)),
			},
		},
		proposer.address(),
	)
	s.Require().NoError(err)
	submitProposalResponse, err := s.chain.sendMsgs(*proposerCtx, proposalMsg)
	s.Require().NoError(err)
	s.Require().Zero(submitProposalResponse.Code, "raw log: %s", submitProposalResponse.RawLog)

	s.T().Log("Checking proposal was submitted correctly")
	govQueryClient := govtypesv1beta1.NewQueryClient(orchClientCtx)
	s.Require().Eventually(func() bool {
		proposalsQueryResponse, err := govQueryClient.Proposals(context.Background(), &govtypesv1beta1.QueryProposalsRequest{})
		if err != nil {
			s.T().Logf("error querying proposals: %e", err)
			return false
		}

		s.Require().NotEmpty(proposalsQueryResponse.Proposals)
		s.Require().Equal(uint64(1), proposalsQueryResponse.Proposals[0].ProposalId, "not proposal id 1")
		s.Require().Equal(govtypesv1beta1.StatusVotingPeriod, proposalsQueryResponse.Proposals[0].Status, "proposal not in voting period")

		return true
	}, time.Second*30, time.Second*5, "proposal submission was never found")

	s.T().Log("Vote for proposal")
	for _, val := range s.chain.validators {
		kr, err := val.keyring()
		s.Require().NoError(err)
		localClientCtx, err := s.chain.clientContext("tcp://localhost:26657", &kr, "val", val.address())
		s.Require().NoError(err)

		voteMsg := govtypesv1beta1.NewMsgVote(val.address(), 1, govtypesv1beta1.OptionYes)
		voteResponse, err := s.chain.sendMsgs(*localClientCtx, voteMsg)
		s.Require().NoError(err)
		s.Require().Zero(voteResponse.Code, "Vote error: %s", voteResponse.RawLog)
	}

	// Wait for proposal to be approved
	s.T().Log("Waiting for proposal to be approved")
	s.Require().Eventually(func() bool {
		proposalQueryResponse, _ := govQueryClient.Proposal(context.Background(), &govtypesv1beta1.QueryProposalRequest{ProposalId: 1})
		return govtypesv1beta1.StatusPassed == proposalQueryResponse.Proposal.Status
	}, time.Second*30, time.Second*5, "proposal was never accepted")
	s.T().Log("proposal approved!")

	// Wait for a few blocks to pass to allow the validator rewards to increase
	s.T().Log("Waiting for a few blocks to pass")
	s.waitForBlocks(10)

	// Get the updated outstanding rewards for the validator
	s.T().Log("Getting updated validator rewards")
	currentRewards2, err := s.getValidatorOutstandingRewards(validator)
	s.Require().NoError(err)

	// Check if the validator's outstanding rewards have increased
	s.T().Logf("Initial rewards: %s, updated rewards: %s", initialRewards.Rewards, currentRewards2.Rewards)
	s.Require().True(currentRewards2.Rewards.AmountOf(params.BaseCoinUnit).GT(initialRewards.Rewards.AmountOf(params.BaseCoinUnit)),
		"Expected validator rewards to increase, got initial: %s, updated: %s", initialRewards.Rewards, currentRewards2.Rewards)

	s.T().Logf("Validator rewards increased from %s to %s", initialRewards.Rewards, currentRewards2.Rewards)

	s.T().Logf("Waiting to see validator rewards cut off at height %d", cutoffHeight)
	s.waitUntilHeight(int64(cutoffHeight))

	s.T().Log("Getting current validator rewards")
	currentRewards, err := s.getValidatorOutstandingRewards(validator)
	s.Require().NoError(err)

	s.T().Logf("Current rewards: %s", currentRewards.Rewards)

	s.T().Log("Waiting for a few blocks to pass")
	s.waitForBlocks(10)

	s.T().Log("Getting updated validator rewards")
	currentRewards2, err = s.getValidatorOutstandingRewards(validator)
	s.Require().NoError(err)

	s.T().Logf("Current rewards: %s", currentRewards2.Rewards)
	s.Require().Equal(currentRewards.Rewards, currentRewards2.Rewards, "Expected validator rewards to remain constant after cutoff height")
	s.T().Log("Validator rewards ended!")

	s.T().Log("Getting sum of all validator rewards")
	totalRewards := sdk.DecCoins{}
	for _, val := range s.chain.validators {
		rewards, err := s.getValidatorOutstandingRewards(val)
		s.Require().NoError(err)
		totalRewards = totalRewards.Add(rewards.Rewards...)
	}
	s.T().Logf("Total rewards: %s", totalRewards)

	s.T().Log("Getting community pool balance")
	queryRes, err = queryClient.CommunityPool(ctx, &disttypes.QueryCommunityPoolRequest{})
	s.Require().NoError(err)
	s.T().Logf("Community pool balance: %s", queryRes.Pool.String())

	// Subtract the initial rewards and the tx fees from the total rewards to get the incentive rewards
	s.T().Log("Total incentive rewards is current rewards minus initial rewards minus tx fees from the proposal submission and votes")
	totalIncentiveRewards := totalRewards.Sub(initialRewards.Rewards.MulDec(sdk.NewDec(4))).Sub(sdk.DecCoins{
		{
			Denom:  testDenom,
			Amount: sdk.NewDec(246913560),
		},
	}.MulDec(sdk.NewDec(5)))
	s.T().Logf("Total incentive rewards: %s", totalIncentiveRewards)

	s.T().Log("Checking that the total incentive rewards are equal to the community pool balance")
	s.T().Logf("Initial community pool: %s, updated community pool: %s", initialCommunityPool, queryRes.Pool)
	s.Require().Equal(totalIncentiveRewards, initialCommunityPool.Sub(queryRes.Pool), "Expected sum of all validator rewards to be equal to the change in community pool balance")
}

func (s *IntegrationTestSuite) getValidatorOutstandingRewards(val *validator) (disttypes.ValidatorOutstandingRewards, error) {
	ctx := context.Background()
	kb, err := val.keyring()
	s.Require().NoError(err)
	clientCtx, err := s.chain.clientContext("tcp://localhost:26657", &kb, "val", val.address())
	s.Require().NoError(err)
	queryClient := disttypes.NewQueryClient(clientCtx)
	resp, err := queryClient.ValidatorOutstandingRewards(
		ctx,
		&disttypes.QueryValidatorOutstandingRewardsRequest{
			ValidatorAddress: val.validatorAddress().String(),
		},
	)
	if err != nil {
		return disttypes.ValidatorOutstandingRewards{}, err
	}
	return resp.Rewards, nil
}

func (s *IntegrationTestSuite) waitForBlocks(numBlocks int64) error {
	validator := s.chain.validators[0]
	kb, err := validator.keyring()
	s.Require().NoError(err)
	clientCtx, err := s.chain.clientContext("tcp://localhost:26657", &kb, "val", validator.address())
	s.Require().NoError(err)

	initialHeight, err := s.getCurrentHeight(clientCtx)
	s.Require().NoError(err)
	targetHeight := initialHeight + numBlocks

	for {
		height, err := s.getCurrentHeight(clientCtx)
		if err != nil {
			return err
		}

		if height >= targetHeight {
			break
		}

		time.Sleep(time.Second)
	}

	return nil
}

func (s *IntegrationTestSuite) waitUntilHeight(height int64) error {
	validator := s.chain.validators[0]
	kb, err := validator.keyring()
	s.Require().NoError(err)
	clientCtx, err := s.chain.clientContext("tcp://localhost:26657", &kb, "val", validator.address())
	s.Require().NoError(err)

	errorsTotal := 0
	for {
		if errorsTotal > 5 {
			return fmt.Errorf("failed to get to height %d: too many errors", height)
		}

		currentHeight, err := s.getCurrentHeight(clientCtx)
		if err != nil {
			errorsTotal++
			continue
		}

		if currentHeight >= height {
			break
		}

		time.Sleep(time.Second * 3)
	}

	return nil
}
