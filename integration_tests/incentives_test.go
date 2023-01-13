package integration_tests

import (
	"context"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	disttypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	paramsproposal "github.com/cosmos/cosmos-sdk/x/params/types/proposal"
	"github.com/peggyjv/sommelier/v4/app/params"
	incentivestypes "github.com/peggyjv/sommelier/v4/x/incentives/types"
)

func (s *IntegrationTestSuite) TestIncentives() {
	s.Run("Bring up chain, observe incentives distribution after param proposal executes", func() {
		ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
		defer cancel()

		val := s.chain.validators[0]
		kb, err := val.keyring()
		s.Require().NoError(err)

		clientCtx, err := s.chain.clientContext("tcp://localhost:26657", &kb, "val", val.keyInfo.GetAddress())
		s.Require().NoError(err)

		incentivesQueryClient := incentivestypes.NewQueryClient(clientCtx)
		distQueryClient := disttypes.NewQueryClient(clientCtx)

		s.T().Log("verifying that the base distribution rate is zero")
		beforeAmount := s.queryDelegationRewardAmount(ctx, val.keyInfo.GetAddress().String(), distQueryClient)
		time.Sleep(time.Second * 12)
		afterAmount := s.queryDelegationRewardAmount(ctx, val.keyInfo.GetAddress().String(), distQueryClient)
		s.Require().Equal(beforeAmount, afterAmount)

		s.T().Log("submitting proposal to enable incentives")
		proposer := s.chain.proposer
		proposerCtx, err := s.chain.clientContext("tcp://localhost:26657", proposer.keyring, "proposer", proposer.keyInfo.GetAddress())
		s.Require().NoError(err)
		orch := s.chain.orchestrators[0]
		orchClientCtx, err := s.chain.clientContext("tcp://localhost:26657", orch.keyring, "orch", orch.keyInfo.GetAddress())
		s.Require().NoError(err)

		proposal := paramsproposal.ParameterChangeProposal{
			Title:       "enable incentives",
			Description: "enables incentives",
			Changes: []paramsproposal.ParamChange{
				{
					Subspace: "incentives",
					Key:      "IncentivesCutoffHeight",
					Value:    "\"1000\"",
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
		incentivesParamsRes, err := incentivesQueryClient.QueryParams(ctx, &incentivestypes.QueryParamsRequest{})
		s.Require().NoError(err)
		s.Require().Equal(incentivesParamsRes.Params.IncentivesCutoffHeight, uint64(1000))

		s.T().Log("verifying that the base distribution rate is greater than zero")
		beforeAmount = s.queryDelegationRewardAmount(ctx, val.keyInfo.GetAddress().String(), distQueryClient)
		time.Sleep(time.Second * 12)
		afterAmount = s.queryDelegationRewardAmount(ctx, val.keyInfo.GetAddress().String(), distQueryClient)
		s.Require().Greater(afterAmount, beforeAmount)
	})
}

func (s *IntegrationTestSuite) queryDelegationRewardAmount(ctx context.Context, delegatorAddress string, distQueryClient disttypes.QueryClient) sdk.Dec {
	rewardsRes, err := distQueryClient.DelegationRewards(ctx, &disttypes.QueryDelegationRewardsRequest{
		DelegatorAddress: delegatorAddress,
		ValidatorAddress: "sommvaloper199sjfhaw3hempwzljw0lgwsm9kk6r8e5ef3hmp",
	})
	s.Require().NoError(err)
	return rewardsRes.Rewards.AmountOf(params.BaseCoinUnit)
}
