package integration_tests

import (
	"context"
	"encoding/hex"
	"fmt"
	"strings"
	"time"

	errorsmod "cosmossdk.io/errors"
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	govtypesv1beta1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1beta1"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	gbtypes "github.com/peggyjv/gravity-bridge/module/v6/x/gravity/types"
	types "github.com/peggyjv/sommelier/v9/x/cork/types/v2"
	pubsubtypes "github.com/peggyjv/sommelier/v9/x/pubsub/types"
)

func (s *IntegrationTestSuite) TestScheduledCork() {
	s.Run("Bring up chain, and schedule a cork call to ethereum", func() {
		s.T().Log("submitting a scheduled cork porposal with unsupported cellar ID to verify rejection")
		proposer := s.chain.proposer
		orch0 := s.chain.orchestrators[0]
		orch0Ctx, err := s.chain.clientContext("tcp://localhost:26657", orch0.keyring, "orch", s.chain.orchestrators[0].address())
		s.Require().NoError(err)
		proposerCtx, err := s.chain.clientContext("tcp://localhost:26657", proposer.keyring, "proposer", proposer.address())
		s.Require().NoError(err)
		invalidProposal := types.ScheduledCorkProposal{
			Title:                 "invalid proposal",
			Description:           "proposal for cellar ID that doesn't exist",
			BlockHeight:           100,
			TargetContractAddress: "0x0000000000000000000000000000000000000000",
			ContractCallProtoJson: "{\"thing\": 1}",
		}
		proposalMsg, err := govtypesv1beta1.NewMsgSubmitProposal(
			&invalidProposal,
			sdk.Coins{
				{
					Denom:  testDenom,
					Amount: stakeAmount.Quo(sdk.NewInt(2)),
				},
			},
			proposer.address(),
		)
		s.Require().NoError(err, "unable to create governance proposal")

		submitProposalResponse, err := s.chain.sendMsgs(*orch0Ctx, proposalMsg)
		s.Require().NoError(err)
		s.Require().Zero(submitProposalResponse.Code, "raw log: %s", submitProposalResponse.RawLog)
		govQueryClient := govtypesv1beta1.NewQueryClient(proposerCtx)
		proposalsQueryResponse, err := govQueryClient.Proposals(context.Background(), &govtypesv1beta1.QueryProposalsRequest{})
		s.Require().NoError(err)
		s.Require().Empty(proposalsQueryResponse.Proposals)
		s.T().Log("proposal rejected as expected!")

		// makes sure ethereum can be contacted and counter contract is working
		count, err := s.getCurrentCount()
		s.Require().NoError(err)
		s.Require().Equal(int64(0), count.Int64())

		s.T().Logf("verify no corks are scheduled")
		val := s.chain.validators[0]
		s.Require().Eventuallyf(func() bool {
			kb, err := val.keyring()
			s.Require().NoError(err)
			clientCtx, err := s.chain.clientContext("tcp://localhost:26657", &kb, "val", val.address())
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
		proposal := types.AddManagedCellarIDsProposal{
			Title:       "add counter contract in test",
			Description: "test description",
			CellarIds: &types.CellarIDSet{
				Ids: []string{counterContract.Hex()},
			},
			PublisherDomain: "example.com",
		}
		proposalMsg, err = govtypesv1beta1.NewMsgSubmitProposal(
			&proposal,
			sdk.Coins{
				{
					Denom:  testDenom,
					Amount: stakeAmount.Quo(sdk.NewInt(2)),
				},
			},
			proposer.address(),
		)
		s.Require().NoError(err, "unable to create governance proposal")

		s.T().Log("submit proposal adding test cellar ID")
		submitProposalResponse, err = s.chain.sendMsgs(*proposerCtx, proposalMsg)
		s.Require().NoError(err)
		s.Require().Zero(submitProposalResponse.Code, "raw log: %s", submitProposalResponse.RawLog)

		s.T().Log("check proposal was submitted correctly")
		s.Require().Eventually(func() bool {
			proposalsQueryResponse, err = govQueryClient.Proposals(context.Background(), &govtypesv1beta1.QueryProposalsRequest{})
			if err != nil {
				return false
			}
			if len(proposalsQueryResponse.Proposals) == 0 {
				return false
			}
			return govtypesv1beta1.StatusVotingPeriod == proposalsQueryResponse.Proposals[0].Status
		}, time.Second*30, time.Second*5, "proposal submission was never found")

		s.T().Log("vote for proposal allowing contract")
		for _, val := range s.chain.validators {
			kr, err := val.keyring()
			s.Require().NoError(err)
			clientCtx, err := s.chain.clientContext("tcp://localhost:26657", &kr, "val", val.address())
			s.Require().NoError(err)

			voteMsg := govtypesv1beta1.NewMsgVote(val.address(), 1, govtypesv1beta1.OptionYes)
			voteResponse, err := s.chain.sendMsgs(*clientCtx, voteMsg)
			s.Require().NoError(err)
			s.Require().Zero(voteResponse.Code, "vote error: %s", voteResponse.RawLog)
		}

		s.T().Log("wait for proposal to be approved")
		s.Require().Eventuallyf(func() bool {
			proposalQueryResponse, err := govQueryClient.Proposal(context.Background(), &govtypesv1beta1.QueryProposalRequest{ProposalId: 1})
			s.Require().NoError(err)
			return govtypesv1beta1.StatusPassed == proposalQueryResponse.Proposal.Status
		}, time.Second*30, time.Second*5, "proposal was never accepted")

		s.T().Log("verify that contract exists in allowed addresses")
		val = s.chain.validators[0]
		s.Require().Eventuallyf(func() bool {
			kb, err := val.keyring()
			s.Require().NoError(err)
			clientCtx, err := s.chain.clientContext("tcp://localhost:26657", &kb, "val", val.address())
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

		s.T().Log("verify a default subscription was created")
		pubsubQueryClient := pubsubtypes.NewQueryClient(proposerCtx)
		subscriptionID := fmt.Sprintf("1:%s", counterContract.String())
		pubsubResponse, err := pubsubQueryClient.QueryDefaultSubscription(context.Background(), &pubsubtypes.QueryDefaultSubscriptionRequest{SubscriptionId: subscriptionID})
		s.Require().NoError(err)
		s.Require().Equal(pubsubResponse.DefaultSubscription.SubscriptionId, subscriptionID)
		s.Require().Equal(pubsubResponse.DefaultSubscription.PublisherDomain, "example.com")

		s.T().Log("schedule a cork for the future")
		node, err := proposerCtx.GetNode()
		s.Require().NoError(err)
		status, err := node.Status(context.Background())
		s.Require().NoError(err)
		currentBlockHeight := status.SyncInfo.LatestBlockHeight
		targetBlockHeight := currentBlockHeight + 15

		s.T().Logf("scheduling cork calls for height %d", targetBlockHeight)
		cork := types.Cork{
			EncodedContractCall:   ABIEncodedInc(),
			TargetContractAddress: counterContract.Hex(),
		}
		corkID := cork.IDHash(uint64(targetBlockHeight))
		s.T().Logf("cork ID is %s", hex.EncodeToString(corkID))
		for i, orch := range s.chain.orchestrators {
			i := i
			orch := orch
			clientCtx, err := s.chain.clientContext("tcp://localhost:26657", orch.keyring, "orch", orch.address())
			s.Require().NoError(err)
			corkMsg, err := types.NewMsgScheduleCorkRequest(
				ABIEncodedInc(),
				counterContract,
				uint64(targetBlockHeight),
				orch.address())
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
		corkQueryClient := types.NewQueryClient(proposerCtx)

		s.Require().Eventually(func() bool {
			res, err := corkQueryClient.QueryScheduledCorksByBlockHeight(context.Background(), &types.QueryScheduledCorksByBlockHeightRequest{BlockHeight: uint64(targetBlockHeight)})
			if err != nil {
				return false
			}
			return len(res.Corks) == 4
		}, 10*time.Second, 1*time.Second, "scheduled corks were not created")

		s.T().Log("wait for scheduled height")
		gbClient := gbtypes.NewQueryClient(proposerCtx)
		s.Require().Eventuallyf(func() bool {
			kb, err := val.keyring()
			s.Require().NoError(err)
			clientCtx, err := s.chain.clientContext("tcp://localhost:26657", &kb, "val", val.address())
			s.Require().NoError(err)

			node, err := clientCtx.GetNode()
			s.Require().NoError(err)
			status, err := node.Status(context.Background())
			s.Require().NoError(err)

			currentHeight := status.SyncInfo.LatestBlockHeight
			if currentHeight > (targetBlockHeight + 1) {
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
		resultRes, err := corkQueryClient.QueryCorkResult(context.Background(), &types.QueryCorkResultRequest{Id: hex.EncodeToString(corkID)})
		s.Require().NoError(err, "failed to query cork result")
		s.Require().True(resultRes.CorkResult.Approved, "cork was not approved")
		s.Require().True(sdk.MustNewDecFromStr(resultRes.CorkResult.ApprovalPercentage).GT(corkVoteThreshold))
		s.Require().Equal(counterContract, common.HexToAddress(resultRes.CorkResult.Cork.TargetContractAddress))

		s.T().Log("verify scheduled corks were deleted")
		res, err := corkQueryClient.QueryScheduledCorksByBlockHeight(context.Background(), &types.QueryScheduledCorksByBlockHeightRequest{BlockHeight: uint64(targetBlockHeight)})
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
			TargetContractAddress: unusedGenesisContract.Hex(),
			ContractCallProtoJson: "{\"cellar_id\":\"0x123801a7D398351b8bE11C439e05C5B3259aeC9B\",\"cellar_v1\":{\"some_fuction\":{\"function_args\":{}},\"block_height\":12345}}",
			BlockHeight:           uint64(targetBlockHeight),
		}
		corkProposalMsg, err := govtypesv1beta1.NewMsgSubmitProposal(
			&corkProposal,
			sdk.Coins{
				{
					Denom:  testDenom,
					Amount: stakeAmount.Quo(sdk.NewInt(2)),
				},
			},
			proposer.address(),
		)
		s.Require().NoError(err, "Unable to create governance proposal")

		s.T().Log("Submit proposal")
		submitProposalResponse, err = s.chain.sendMsgs(*proposerCtx, corkProposalMsg)
		s.Require().NoError(err)
		s.Require().Zero(submitProposalResponse.Code, "raw log: %s", submitProposalResponse.RawLog)

		s.T().Log("Check proposal was submitted correctly")

		s.Require().Eventually(func() bool {
			proposalsQueryResponse, err := govQueryClient.Proposals(context.Background(), &govtypesv1beta1.QueryProposalsRequest{})
			if err != nil {
				s.T().Logf("error querying proposals: %e", err)
				return false
			}

			s.Require().NotEmpty(proposalsQueryResponse.Proposals)
			s.Require().Equal(uint64(2), proposalsQueryResponse.Proposals[1].ProposalId, "not proposal id 2")
			s.Require().Equal(govtypesv1beta1.StatusVotingPeriod, proposalsQueryResponse.Proposals[1].Status, "proposal not in voting period")

			return true
		}, time.Second*30, time.Second*5, "proposal submission was never found")

		s.T().Log("Vote for proposal")
		for _, val := range s.chain.validators {
			kr, err := val.keyring()
			s.Require().NoError(err)
			localClientCtx, err := s.chain.clientContext("tcp://localhost:26657", &kr, "val", val.address())
			s.Require().NoError(err)

			voteMsg := govtypesv1beta1.NewMsgVote(val.address(), 2, govtypesv1beta1.OptionYes)
			voteResponse, err := s.chain.sendMsgs(*localClientCtx, voteMsg)
			s.Require().NoError(err)
			s.Require().Zero(voteResponse.Code, "Vote error: %s", voteResponse.RawLog)
		}

		s.T().Log("Waiting for proposal to be approved..")
		s.Require().Eventually(func() bool {
			proposalQueryResponse, _ := govQueryClient.Proposal(context.Background(), &govtypesv1beta1.QueryProposalRequest{ProposalId: 2})
			return govtypesv1beta1.StatusPassed == proposalQueryResponse.Proposal.Status
		}, time.Second*30, time.Second*5, "proposal was never accepted")
		s.T().Log("Proposal approved!")

		s.T().Logf("create governance proposal to remove counter contract")
		removalProposal := types.RemoveManagedCellarIDsProposal{
			Title:       "add counter contract in test",
			Description: "test description",
			CellarIds: &types.CellarIDSet{
				Ids: []string{counterContract.Hex()},
			},
		}
		proposalMsg, err = govtypesv1beta1.NewMsgSubmitProposal(
			&removalProposal,
			sdk.Coins{
				{
					Denom:  testDenom,
					Amount: stakeAmount.Quo(sdk.NewInt(2)),
				},
			},
			proposer.address(),
		)
		s.Require().NoError(err, "unable to create governance proposal")

		s.T().Log("submit proposal adding test cellar ID")
		submitProposalResponse, err = s.chain.sendMsgs(*proposerCtx, proposalMsg)
		s.Require().NoError(err)
		s.Require().Zero(submitProposalResponse.Code, "raw log: %s", submitProposalResponse.RawLog)

		s.T().Log("check proposal was submitted correctly")
		s.Require().Eventually(func() bool {
			proposalsQueryResponse, err = govQueryClient.Proposals(context.Background(), &govtypesv1beta1.QueryProposalsRequest{})
			if err != nil {
				return false
			}
			if len(proposalsQueryResponse.Proposals) == 0 {
				return false
			}
			s.Require().Equal(uint64(3), proposalsQueryResponse.Proposals[2].ProposalId, "not proposal id 3")
			return govtypesv1beta1.StatusVotingPeriod == proposalsQueryResponse.Proposals[2].Status
		}, time.Second*30, time.Second*5, "proposal submission was never found")

		s.T().Log("vote for proposal allowing contract")
		for _, val := range s.chain.validators {
			kr, err := val.keyring()
			s.Require().NoError(err)
			clientCtx, err := s.chain.clientContext("tcp://localhost:26657", &kr, "val", val.address())
			s.Require().NoError(err)

			voteMsg := govtypesv1beta1.NewMsgVote(val.address(), 3, govtypesv1beta1.OptionYes)
			voteResponse, err := s.chain.sendMsgs(*clientCtx, voteMsg)
			s.Require().NoError(err)
			s.Require().Zero(voteResponse.Code, "vote error: %s", voteResponse.RawLog)
		}

		s.T().Log("wait for proposal to be approved")
		s.Require().Eventuallyf(func() bool {
			proposalQueryResponse, err := govQueryClient.Proposal(context.Background(), &govtypesv1beta1.QueryProposalRequest{ProposalId: 3})
			s.Require().NoError(err)
			return govtypesv1beta1.StatusPassed == proposalQueryResponse.Proposal.Status
		}, time.Second*30, time.Second*5, "proposal was never accepted")

		s.T().Log("verify cellar ID was removed")
		queryClient := types.NewQueryClient(proposerCtx)
		cellarIDsResponse, err := queryClient.QueryCellarIDs(context.Background(), &types.QueryCellarIDsRequest{})
		s.Require().NoError(err)
		s.Require().NotContains(cellarIDsResponse.CellarIds, counterContract.String())

		s.T().Log("verify default subscription was removed")
		subscriptionID = fmt.Sprintf("1:%s", counterContract.String())
		_, err = pubsubQueryClient.QueryDefaultSubscription(context.Background(), &pubsubtypes.QueryDefaultSubscriptionRequest{SubscriptionId: subscriptionID})
		s.Require().Error(err)
	})
}

const CounterABI = `
  [
    {
      "inputs": [],
      "name": "count",
      "outputs": [
        {
          "internalType": "uint256",
          "name": "",
          "type": "uint256"
        }
      ],
      "stateMutability": "view",
      "type": "function"
    },
    {
      "inputs": [],
      "name": "dec",
      "outputs": [],
      "stateMutability": "nonpayable",
      "type": "function"
    },
    {
      "inputs": [],
      "name": "get",
      "outputs": [
        {
          "internalType": "uint256",
          "name": "",
          "type": "uint256"
        }
      ],
      "stateMutability": "view",
      "type": "function"
    },
    {
      "inputs": [],
      "name": "inc",
      "outputs": [],
      "stateMutability": "nonpayable",
      "type": "function"
    }
  ]
`

func ABIEncodedGet() []byte {
	encodedCall, err := abi.JSON(strings.NewReader(CounterABI))
	if err != nil {
		panic(errorsmod.Wrap(err, "bad ABI definition in code"))
	}

	abiEncodedCall, err := encodedCall.Pack("get")
	if err != nil {
		panic(err)
	}

	return abiEncodedCall
}

func ABIEncodedInc() []byte {
	encodedCall, err := abi.JSON(strings.NewReader(CounterABI))
	if err != nil {
		panic(errorsmod.Wrap(err, "bad ABI definition in code"))
	}

	abiEncodedCall, err := encodedCall.Pack("inc")
	if err != nil {
		panic(err)
	}

	return abiEncodedCall
}

func (s *IntegrationTestSuite) getCurrentCount() (*math.Int, error) {
	ethClient, err := ethclient.Dial(fmt.Sprintf("http://%s", s.ethResource.GetHostPort("8545/tcp")))
	if err != nil {
		return nil, err
	}

	bz, err := ethClient.CallContract(context.Background(), ethereum.CallMsg{
		From: common.HexToAddress(s.chain.validators[0].ethereumKey.address),
		To:   &counterContract,
		Gas:  0,
		Data: ABIEncodedGet(),
	}, nil)
	if err != nil {
		return nil, err
	}

	count := UnpackEthUInt(bz)

	return &count, nil
}
