package integration_tests

import (
	"context"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	"github.com/ethereum/go-ethereum/common"
	gbtypes "github.com/peggyjv/gravity-bridge/module/v3/x/gravity/types"
	"github.com/peggyjv/sommelier/v6/x/cork/types"
)

func (s *IntegrationTestSuite) TestScheduledCork() {
	s.Run("Bring up chain, and schedule a cork call to ethereum", func() {

		// makes sure ethereum can be contacted and counter contract is working
		count, err := s.getCurrentCount()
		s.Require().NoError(err)
		s.Require().Equal(int64(0), count.Int64())

		s.T().Logf("verify no corks are scheduled")
		val := s.chain.validators[0]
		s.Require().Eventuallyf(func() bool {
			kb, err := val.keyring()
			s.Require().NoError(err)
			clientCtx, err := s.chain.clientContext("tcp://localhost:26657", &kb, "val", val.keyInfo.GetAddress())
			s.Require().NoError(err)

			queryClient := types.NewQueryClient(clientCtx)
			res, err := queryClient.QueryScheduledCorks(context.Background(), &types.QueryScheduledCorksRequest{})
			if err != nil {
				s.T().Logf("error: %s", err)
				return false
			}

			return len(res.Corks) == 0
		}, 20*time.Second, 1*time.Second, "got a non-empty result for scheduled corks")

		s.T().Logf("create governance proposal to add counter contract")
		orch := s.chain.orchestrators[0]
		clientCtx, err := s.chain.clientContext("tcp://localhost:26657", orch.keyring, "orch", orch.keyInfo.GetAddress())
		s.Require().NoError(err)

		proposal := types.AddManagedCellarIDsProposal{
			Title:       "add counter contract in test",
			Description: "test description",
			CellarIds: &types.CellarIDSet{
				Ids: []string{counterContract.Hex()},
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
			orch.keyInfo.GetAddress(),
		)
		s.Require().NoError(err, "unable to create governance proposal")

		s.T().Log("submit proposal adding test cellar ID")
		submitProposalResponse, err := s.chain.sendMsgs(*clientCtx, proposalMsg)
		s.Require().NoError(err)
		s.Require().Zero(submitProposalResponse.Code, "raw log: %s", submitProposalResponse.RawLog)

		s.T().Log("check proposal was submitted correctly")
		govQueryClient := govtypes.NewQueryClient(clientCtx)
		proposalsQueryResponse, err := govQueryClient.Proposals(context.Background(), &govtypes.QueryProposalsRequest{})
		s.Require().NoError(err)
		s.Require().NotEmpty(proposalsQueryResponse.Proposals)
		s.Require().Equal(uint64(1), proposalsQueryResponse.Proposals[0].ProposalId, "not proposal id 1")
		s.Require().Equal(govtypes.StatusVotingPeriod, proposalsQueryResponse.Proposals[0].Status, "proposal not in voting period")

		s.T().Log("vote for proposal allowing contract")
		for _, val := range s.chain.validators {
			kr, err := val.keyring()
			s.Require().NoError(err)
			clientCtx, err := s.chain.clientContext("tcp://localhost:26657", &kr, "val", val.keyInfo.GetAddress())
			s.Require().NoError(err)

			voteMsg := govtypes.NewMsgVote(val.keyInfo.GetAddress(), 1, govtypes.OptionYes)
			voteResponse, err := s.chain.sendMsgs(*clientCtx, voteMsg)
			s.Require().NoError(err)
			s.Require().Zero(voteResponse.Code, "vote error: %s", voteResponse.RawLog)
		}

		s.T().Log("wait for proposal to be approved")
		s.Require().Eventuallyf(func() bool {
			proposalQueryResponse, err := govQueryClient.Proposal(context.Background(), &govtypes.QueryProposalRequest{ProposalId: 1})
			s.Require().NoError(err)
			return govtypes.StatusPassed == proposalQueryResponse.Proposal.Status
		}, time.Second*30, time.Second*5, "proposal was never accepted")

		s.T().Log("verify that contract exists in allowed addresses")
		val = s.chain.validators[0]
		s.Require().Eventuallyf(func() bool {
			kb, err := val.keyring()
			s.Require().NoError(err)
			clientCtx, err := s.chain.clientContext("tcp://localhost:26657", &kb, "val", val.keyInfo.GetAddress())
			s.Require().NoError(err)

			queryClient := types.NewQueryClient(clientCtx)
			res, err := queryClient.QueryCellarIDs(context.Background(), &types.QueryCellarIDsRequest{})
			if err != nil {
				s.T().Logf("error: %s", err)
				return false
			}

			found := false
			for _, id := range res.CellarIds {
				s.T().Logf("managed addresses: %v", res.CellarIds)
				if common.HexToAddress(id) == counterContract {
					found = true
					break
				}
			}

			return found
		}, 10*time.Second, 2*time.Second, "did not find address in managed cellars")

		s.T().Log("schedule a cork for the future")
		node, err := clientCtx.GetNode()
		s.Require().NoError(err)
		status, err := node.Status(context.Background())
		s.Require().NoError(err)
		currentBlockHeight := status.SyncInfo.LatestBlockHeight
		targetBlockHeight := currentBlockHeight + 15

		s.T().Log("scheduling cork calls")
		for i, orch := range s.chain.orchestrators {
			i := i
			orch := orch
			s.Require().Eventuallyf(func() bool {
				clientCtx, err := s.chain.clientContext("tcp://localhost:26657", orch.keyring, "orch", orch.keyInfo.GetAddress())
				s.Require().NoError(err)

				corkMsg, err := types.NewMsgScheduleCorkRequest(
					ABIEncodedInc(),
					counterContract,
					uint64(targetBlockHeight),
					orch.keyInfo.GetAddress())
				s.Require().NoError(err, "unable to create cork schedule msg")

				response, err := s.chain.sendMsgs(*clientCtx, corkMsg)
				if err != nil {
					s.T().Logf("error: %s", err)
					return false
				}
				if response.Code != 0 {
					if response.Code != 32 {
						s.T().Log(response)
					}
					return false
				}

				s.T().Logf("cork msg for %d node sent successfully", i)
				return true
			}, 10*time.Second, 500*time.Millisecond, "unable to deploy cork schedule msg for node %d", i)
		}

		s.T().Logf("wait for the block to pass, monitoring it as it goes")
		s.Require().Eventuallyf(func() bool {
			kb, err := val.keyring()
			s.Require().NoError(err)
			clientCtx, err := s.chain.clientContext("tcp://localhost:26657", &kb, "val", val.keyInfo.GetAddress())
			s.Require().NoError(err)

			queryClient := types.NewQueryClient(clientCtx)
			res, err := queryClient.QueryScheduledCorks(context.Background(), &types.QueryScheduledCorksRequest{})
			if err != nil {
				s.T().Logf("error: %s", err)
				return false
			}

			node, err := clientCtx.GetNode()
			s.Require().NoError(err)
			status, err := node.Status(context.Background())
			s.Require().NoError(err)
			blockHeight := status.SyncInfo.LatestBlockHeight

			gbClient := gbtypes.NewQueryClient(clientCtx)
			gbRes, err := gbClient.ContractCallTxs(context.Background(), &gbtypes.ContractCallTxsRequest{
				Pagination: nil,
			})
			s.Require().NoError(err)

			if blockHeight < (targetBlockHeight - 2) {
				// verify that tbe scheduled cork has not yet been consumed, and that the counter has not been incremented
				s.Require().Len(res.Corks, len(s.chain.validators))
				s.Require().Equal(counterContract, common.HexToAddress(res.Corks[0].Cork.TargetContractAddress))
				s.Require().Len(gbRes.Calls, 0)
			} else if blockHeight > (targetBlockHeight + 1) {
				// verify that block height has been passed, cork consumed, contractcalltx created
				s.Require().Len(res.Corks, 0)
				s.Require().Len(gbRes.Calls, 1)

				// this is the only situation where this loop will complete
				return true
			}

			return false
		}, 3*time.Minute, 1*time.Second, "count was never updated")

		s.T().Logf("verify count was updated")
		s.Require().Eventuallyf(func() bool {
			count, err = s.getCurrentCount()
			if err != nil {
				return false
			}
			return int64(1) == count.Int64()
		}, time.Minute, 10*time.Second, "count was never updated")
	})
}
