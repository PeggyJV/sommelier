package integration_tests

import (
	"context"
	"fmt"
	"time"

	"github.com/cosmos/cosmos-sdk/client"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/bech32"
	disttypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"
	paramsproposal "github.com/cosmos/cosmos-sdk/x/params/types/proposal"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/peggyjv/sommelier/v6/app/params"
	incentivestypes "github.com/peggyjv/sommelier/v6/x/incentives/types"
)

func (s *IntegrationTestSuite) TestIncentives() {
	s.Run("Bring up chain, observe incentives distribution after param proposal executes", func() {
		ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
		defer cancel()

		val := s.chain.validators[0]
		s.T().Logf("validator %s", val.keyInfo.GetAddress().String())
		kb, err := val.keyring()
		s.Require().NoError(err)

		_, bytes, err := bech32.DecodeAndConvert(val.keyInfo.GetAddress().String())
		s.Require().NoError(err)
		valOperatorAddress, err := bech32.ConvertAndEncode("sommvaloper", bytes)
		s.Require().NoError(err)

		clientCtx, err := s.chain.clientContext("tcp://localhost:26657", &kb, "val", val.keyInfo.GetAddress())
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

		beforeAmount := s.queryValidatorRewards(ctx, valOperatorAddress, distQueryClient)
		time.Sleep(time.Second * 12)
		afterAmount := s.queryValidatorRewards(ctx, valOperatorAddress, distQueryClient)
		s.Require().Equal(beforeAmount, afterAmount)

		s.T().Log("submitting proposal to enable incentives")
		proposer := s.chain.proposer
		proposerCtx, err := s.chain.clientContext("tcp://localhost:26657", proposer.keyring, "proposer", proposer.keyInfo.GetAddress())
		s.Require().NoError(err)
		orch := s.chain.orchestrators[0]
		orchClientCtx, err := s.chain.clientContext("tcp://localhost:26657", orch.keyring, "orch", orch.keyInfo.GetAddress())
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

		proposalMsg, err := govtypes.NewMsgSubmitProposal(
			&proposal,
			sdk.Coins{
				{
					Denom:  testDenom,
					Amount: stakeAmount.Quo(sdk.NewInt(2)),
				},
			},
			proposer.keyInfo.GetAddress(),
		)
		s.Require().NoError(err, "Unable to create governance proposal")
		submitProposalResponse, err := s.chain.sendMsgs(*proposerCtx, proposalMsg)
		s.Require().NoError(err)
		s.Require().Zero(submitProposalResponse.Code, "raw log: %s", submitProposalResponse.RawLog)

		s.T().Log("check proposal was submitted correctly")
		govQueryClient := govtypes.NewQueryClient(orchClientCtx)

		s.Require().Eventually(func() bool {
			proposalsQueryResponse, err := govQueryClient.Proposals(context.Background(), &govtypes.QueryProposalsRequest{})
			if err != nil {
				s.T().Logf("error querying proposals: %e", err)
				return false
			}

			s.Require().NotEmpty(proposalsQueryResponse.Proposals)
			s.Require().Equal(uint64(1), proposalsQueryResponse.Proposals[0].ProposalId, "not proposal id 1")
			s.Require().Equal(govtypes.StatusVotingPeriod, proposalsQueryResponse.Proposals[0].Status, "proposal not in voting period")

			return true
		}, time.Second*30, time.Second*5, "proposal submission was never found")

		s.T().Log("vote for proposal")
		for _, val := range s.chain.validators {
			kr, err := val.keyring()
			s.Require().NoError(err)
			localClientCtx, err := s.chain.clientContext("tcp://localhost:26657", &kr, "val", val.keyInfo.GetAddress())
			s.Require().NoError(err)

			voteMsg := govtypes.NewMsgVote(val.keyInfo.GetAddress(), 1, govtypes.OptionYes)
			voteResponse, err := s.chain.sendMsgs(*localClientCtx, voteMsg)
			s.Require().NoError(err)
			s.Require().Zero(voteResponse.Code, "Vote error: %s", voteResponse.RawLog)
		}

		s.T().Log("waiting for proposal to be approved..")
		s.Require().Eventually(func() bool {
			proposalQueryResponse, _ := govQueryClient.Proposal(context.Background(), &govtypes.QueryProposalRequest{ProposalId: 1})
			return govtypes.StatusPassed == proposalQueryResponse.Proposal.Status
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
		s.Require().Equal(expectedDistributionPerBlock.Amount.ToDec(), actualDistributionPerBlock)

		s.T().Log("verifying APY query returns expected value")
		mintParams, err := mintQueryClient.Params(ctx, &minttypes.QueryParamsRequest{})
		s.Require().NoError(err)
		stakingPool, err := stakingQueryClient.Pool(ctx, &stakingtypes.QueryPoolRequest{})
		s.Require().NoError(err)
		tokensTotal := stakingPool.Pool.BondedTokens.Add(stakingPool.Pool.NotBondedTokens)
		// assumes bonded ratio is 100%
		expectedAPY := actualDistributionPerBlock.Mul(sdk.NewDec(int64(mintParams.Params.BlocksPerYear))).Quo(tokensTotal.ToDec())
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

func (s *IntegrationTestSuite) getCurrentHeight(clientCtx *client.Context) int64 {
	node, err := clientCtx.GetNode()
	s.Require().NoError(err)
	status, err := node.Status(context.Background())
	s.Require().NoError(err)

	return status.SyncInfo.LatestBlockHeight
}

func (s *IntegrationTestSuite) getRewardAmountAndHeight(ctx context.Context, distQueryClient disttypes.QueryClient, operatorAddress string, clientCtx *client.Context) (sdk.Dec, int64) {
	var amount sdk.Dec
	var height int64

	s.Require().Eventually(func() bool {
		initialHeight := s.getCurrentHeight(clientCtx)
		amount = s.queryValidatorRewards(ctx, operatorAddress, distQueryClient)
		height = s.getCurrentHeight(clientCtx)
		return initialHeight == height
	}, time.Second*30, time.Second*1, "failed to reliably determine height of reward sample")

	return amount, height
}
