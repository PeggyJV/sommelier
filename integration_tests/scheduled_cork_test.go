package integration_tests

import (
	"bytes"
	"context"
	"encoding/hex"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	gbtypes "github.com/peggyjv/gravity-bridge/module/v2/x/gravity/types"
	"github.com/peggyjv/sommelier/v4/x/cork/types"
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

		s.T().Logf("scheduling cork calls for height %d", targetBlockHeight)
		blockHeightBytes := sdk.Uint64ToBigEndian(uint64(targetBlockHeight))
		corkID := hex.EncodeToString(crypto.Keccak256Hash(
			bytes.Join(
				[][]byte{blockHeightBytes, counterContract.Bytes(), ABIEncodedInc()},
				[]byte{},
			)).Bytes())
		s.T().Logf("cork ID is %s", corkID)
		for i, orch := range s.chain.orchestrators {
			i := i
			orch := orch
			clientCtx, err := s.chain.clientContext("tcp://localhost:26657", orch.keyring, "orch", orch.keyInfo.GetAddress())
			s.Require().NoError(err)
			corkMsg, err := types.NewMsgScheduleCorkRequest(
				ABIEncodedInc(),
				counterContract,
				uint64(targetBlockHeight),
				orch.keyInfo.GetAddress())
			s.Require().NoError(err, "failed to construct cork")
			response, err := s.chain.sendMsgs(*clientCtx, corkMsg)
			s.Require().NoError(err, "failed to send cork to node %d", i)
			if response.Code != 0 {
				if response.Code != 32 {
					s.T().Log(response)
				}
			}

			s.T().Logf("cork msg for orch %d sent successfully", i)
		}

		s.T().Log("verify scheduled corks were created")
		corkQueryClient := types.NewQueryClient(clientCtx)
		res, err := corkQueryClient.QueryScheduledCorksByBlockHeight(context.Background(), &types.QueryScheduledCorksByBlockHeightRequest{BlockHeight: uint64(targetBlockHeight)})
		s.Require().NoError(err, "failed to query scheduled corks by height")
		s.Require().Len(res.Corks, 4)

		s.T().Log("wait for scheduled height")
		gbClient := gbtypes.NewQueryClient(clientCtx)
		s.Require().Eventuallyf(func() bool {
			kb, err := val.keyring()
			s.Require().NoError(err)
			clientCtx, err := s.chain.clientContext("tcp://localhost:26657", &kb, "val", val.keyInfo.GetAddress())
			s.Require().NoError(err)

			node, err := clientCtx.GetNode()
			s.Require().NoError(err)
			status, err := node.Status(context.Background())
			s.Require().NoError(err)

			currentHeight := status.SyncInfo.LatestBlockHeight
			if currentHeight > (targetBlockHeight + 1) {
				// blockHeight = uint64(status.SyncInfo.LatestBlockHeight)
				return true
			} else if currentHeight < targetBlockHeight {
				res, err := corkQueryClient.QueryScheduledCorks(context.Background(), &types.QueryScheduledCorksRequest{})
				if err != nil {
					s.T().Logf("error: %s", err)
					return false
				}

				// verify that the scheduled corks have not yet been consumed
				s.Require().Len(res.Corks, len(s.chain.validators))
			}

			return false
		}, 3*time.Minute, 1*time.Second, "never reached scheduled height")

		s.T().Log("verify the cork was approved")
		resultRes, err := corkQueryClient.QueryCorkResult(context.Background(), &types.QueryCorkResultRequest{Id: corkID})
		s.Require().NoError(err, "failed to query cork result")
		s.Require().True(resultRes.CorkResult.Approved, "cork was not approved")
		s.Require().True(sdk.MustNewDecFromStr(resultRes.CorkResult.ApprovalPercentage).GT(corkVoteThreshold))
		s.Require().Equal(counterContract, common.HexToAddress(resultRes.CorkResult.Cork.TargetContractAddress))

		s.T().Log("verify scheduled corks were deleted")
		res, err = corkQueryClient.QueryScheduledCorksByBlockHeight(context.Background(), &types.QueryScheduledCorksByBlockHeightRequest{BlockHeight: uint64(targetBlockHeight)})
		s.Require().NoError(err, "failed to query scheduled corks by height")
		s.Require().Len(res.Corks, 0)

		s.T().Log("verify gravity contract call tx was created")
		gbRes, err := gbClient.ContractCallTxs(context.Background(), &gbtypes.ContractCallTxsRequest{
			Pagination: nil,
		})
		s.Require().NoError(err, "failed to query gravity contract call txs")
		s.Require().Len(gbRes.Calls, 1)

		s.T().Log("verify count was updated")
		s.Require().Eventuallyf(func() bool {
			count, err = s.getCurrentCount()
			if err != nil {
				return false
			}
			return int64(1) == count.Int64()
		}, time.Minute, 10*time.Second, "count was never updated")

		s.T().Log("Test governance-scheduled cork creation")
		corkProposal := types.ScheduledCorkProposal{
			Title:                 "initial token price submission",
			Description:           "our first token prices",
			TargetContractAddress: "0x0000000000000000000000000000000000000000",
			ContractCallProtoJson: "{\"cellar_id\":\"0x123801a7D398351b8bE11C439e05C5B3259aeC9B\",\"cellar_v1\":{\"some_fuction\":{\"function_args\":{}},\"block_height\":12345}}",
			BlockHeight:           uint64(targetBlockHeight),
		}

		corkProposalMsg, err := govtypes.NewMsgSubmitProposal(
			&corkProposal,
			sdk.Coins{
				{
					Denom:  testDenom,
					Amount: stakeAmount.Quo(sdk.NewInt(2)),
				},
			},
			orch.keyInfo.GetAddress(),
		)
		s.Require().NoError(err, "Unable to create governance proposal")

		s.T().Log("Submit proposal")
		submitProposalResponse, err = s.chain.sendMsgs(*clientCtx, corkProposalMsg)
		s.Require().NoError(err)
		s.Require().Zero(submitProposalResponse.Code, "raw log: %s", submitProposalResponse.RawLog)

		s.T().Log("Check proposal was submitted correctly")

		s.Require().Eventually(func() bool {
			proposalsQueryResponse, err := govQueryClient.Proposals(context.Background(), &govtypes.QueryProposalsRequest{})
			if err != nil {
				s.T().Logf("error querying proposals: %e", err)
				return false
			}

			s.Require().NotEmpty(proposalsQueryResponse.Proposals)
			s.Require().Equal(uint64(2), proposalsQueryResponse.Proposals[1].ProposalId, "not proposal id 2")
			s.Require().Equal(govtypes.StatusVotingPeriod, proposalsQueryResponse.Proposals[1].Status, "proposal not in voting period")

			return true
		}, time.Second*30, time.Second*5, "proposal submission was never found")

		s.T().Log("Vote for proposal")
		for _, val := range s.chain.validators {
			kr, err := val.keyring()
			s.Require().NoError(err)
			localClientCtx, err := s.chain.clientContext("tcp://localhost:26657", &kr, "val", val.keyInfo.GetAddress())
			s.Require().NoError(err)

			voteMsg := govtypes.NewMsgVote(val.keyInfo.GetAddress(), 2, govtypes.OptionYes)
			voteResponse, err := s.chain.sendMsgs(*localClientCtx, voteMsg)
			s.Require().NoError(err)
			s.Require().Zero(voteResponse.Code, "Vote error: %s", voteResponse.RawLog)
		}

		s.T().Log("Waiting for proposal to be approved..")
		s.Require().Eventually(func() bool {
			proposalQueryResponse, _ := govQueryClient.Proposal(context.Background(), &govtypes.QueryProposalRequest{ProposalId: 2})
			return govtypes.StatusPassed == proposalQueryResponse.Proposal.Status
		}, time.Second*30, time.Second*5, "proposal was never accepted")
		s.T().Log("Proposal approved!")
	})
}
