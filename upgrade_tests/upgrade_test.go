package upgrade_test

import (
	"context"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"
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
			Plan:        upgradetypes.Plan{Name: "chain-upgrade", Height: 100}}

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
			return true
		}, time.Second*30, time.Second*5, "An error occurred when querying block height before upgrade")

		// Stop all nodes
		s.StopAllNodes()

		// Start all nodes
		// initialization

		// run the eth container so that the contract addresses are available
		s.runEthContainer("prebuilt")

		// continue generating node genesis
		s.initGenesis()
		s.initValidatorConfigs()

		// container infrastructure
		s.runValidators("prebuilt")
		s.runOrchestrators("prebuilt")

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
	})
}