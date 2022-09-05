package upgrade_tests

import (
	"context"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"
	"github.com/peggyjv/sommelier/x/allocation/types"
)

func (s *UpgradeTestSuite) TestSommChainUpgrade() {
	s.Run("Bring up chain, and test Sommelier chain upgrade", func() {
		s.T().Logf("create governance proposal for chain upgrade")
		orch := s.chain.orchestrators[0]
		clientCtx, err := s.chain.clientContext("tcp://localhost:26657", orch.keyring, "orch", orch.keyInfo.GetAddress())
		s.Require().NoError(err)

		proposal := upgradetypes.SoftwareUpgradeProposal{
			Title:       "Chain Upgrade 1",
			Description: "First chain software upgrade",
			Plan:        upgradetypes.Plan{Name: "v4", Height: 100}}

		proposalMsg, err := govtypes.NewMsgSubmitProposal(
			&proposal,
			sdk.Coins{
				{
					Denom:  testDenom,
					Amount: sdk.NewInt(2),
				},
			},
			orch.keyInfo.GetAddress(),
		)

		s.Require().NoError(err, "unable to create governance proposal")

		// Submit proposal
		s.T().Log("submit proposal upgrading chain")
		res, err := s.chain.sendMsgs(*clientCtx, proposalMsg)
		s.T().Logf("submit proposal response:%s", res)
		s.Require().NoError(err)

		// Query proposal was submitted successfully
		s.T().Log("check proposal was submitted correctly")
		govQueryClient := govtypes.NewQueryClient(clientCtx)
		proposalsQueryResponse, err := govQueryClient.Proposals(context.Background(), &govtypes.QueryProposalsRequest{})
		s.Require().NoError(err)
		s.Require().NotEmpty(proposalsQueryResponse.Proposals)
		s.Require().Equal(uint64(1), proposalsQueryResponse.Proposals[0].ProposalId, "not proposal id 1")
		s.Require().Equal(govtypes.StatusVotingPeriod, proposalsQueryResponse.Proposals[0].Status, "proposal not in voting period")

		// Vote YES on proposal
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

		// Wait for proposal to be accepted
		s.T().Log("wait for proposal to be approved")
		s.Require().Eventuallyf(func() bool {
			proposalQueryResponse, err := govQueryClient.Proposal(context.Background(), &govtypes.QueryProposalRequest{ProposalId: 1})
			s.Require().NoError(err)
			return govtypes.StatusPassed == proposalQueryResponse.Proposal.Status
		}, time.Second*30, time.Second*5, "proposal was never accepted")

		// Query block height to ensure chain halted
		s.T().Log("Query block height to ensure chain halted")
		s.Require().Eventuallyf(func() bool {
			node, err := clientCtx.GetNode()
			s.Require().NoError(err)
			status, err := node.Status(context.Background())
			s.Require().NoError(err)
			currentBlockHeight := status.SyncInfo.LatestBlockHeight
			s.T().Logf("Current block height:%d", currentBlockHeight)
			if currentBlockHeight == 100 {
				// Stop all nodes
				s.StopAllNodes()
				return true
			}
			return false
		}, time.Minute*10, time.Second*1, "An error occurred when querying block height before upgrade")
		s.Require().Eventuallyf(func() bool {
			// Start all nodes
			// initialization

			// continue generating node genesis
			s.initGenesis()
			s.initValidatorConfigs()

			// container infrastructure
			s.runValidators("upgrade")
			s.runOrchestrators("3.1.0")
			return true
		}, time.Minute*5, time.Second*30, "An error occurred when querying block height before upgrade")

		// Query block height to ensure chain successfully updated
		s.T().Log("Query block height to ensure chain successfully updated")
		s.Require().Eventuallyf(func() bool {
			node, err := clientCtx.GetNode()
			s.Require().NoError(err)
			status, err := node.Status(context.Background())
			s.Require().NoError(err)
			currentBlockHeight := status.SyncInfo.LatestBlockHeight
			s.T().Logf("Current block height:%d", currentBlockHeight)
			return true
		}, time.Minute*5, time.Minute*1, "An error occurred when querying block height after upgrade")

		// Test upgrade by making sure allocation module was removed
		s.T().Log("Checkout chain upgrade")
		tickRange, err := s.getFirstTickRange()
		s.Require().NoError(err)
		s.Require().Equal(int32(600), tickRange.Upper)
		s.Require().Equal(int32(300), tickRange.Lower)
		s.Require().Equal(uint32(900), tickRange.Weight)

		commit := types.Allocation{
			Vote: &types.RebalanceVote{
				Cellar: &types.Cellar{
					Id: hardhatCellar.String(),
					TickRanges: []*types.TickRange{
						{Upper: 198840, Lower: 192180, Weight: 100},
					},
				},
				CurrentPrice: 100,
			},
			Salt: "testsalt",
		}

		s.T().Logf("checking that test cellar exists in the chain")
		val := s.chain.validators[0]
		s.Require().Eventuallyf(func() bool {
			kb, err := val.keyring()
			s.Require().NoError(err)
			clientCtx, err := s.chain.clientContext("tcp://localhost:26657", &kb, "val", val.keyInfo.GetAddress())
			s.Require().NoError(err)
			// 	This query should return an error, since allocation module has been removed"
			queryClient := types.NewQueryClient(clientCtx)
			res, err := queryClient.QueryCellars(context.Background(), &types.QueryCellarsRequest{})
			s.Require().EqualError(err, "rpc error: code = Unknown desc = unknown query path: unknown request")
			if err != nil {
				return true
			}
			if res == nil {
				return true
			}
			for _, c := range res.Cellars {
				if c.Id == commit.Vote.Cellar.Id {
					return true
				}
			}
			return true
		}, 30*time.Second, 2*time.Second, "hardhat cellar not found in chain")
	})
}
