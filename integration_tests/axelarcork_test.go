package integration_tests

import (
	"context"
	"time"

	"cosmossdk.io/math"
	"github.com/cosmos/cosmos-sdk/client"
	sdk "github.com/cosmos/cosmos-sdk/types"
	govtypesv1beta1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1beta1"
	"github.com/peggyjv/sommelier/v7/x/axelarcork/types"
)

func (s *IntegrationTestSuite) TestAxelarCork() {
	s.Run("Test the axelarcork module", func() {
		// Set up validator, orchestrator, proposer, query client
		val0 := s.chain.validators[0]
		val0kb, err := val0.keyring()
		s.Require().NoError(err)
		val0ClientCtx, err := s.chain.clientContext("tcp://localhost:26657", &val0kb, "val", val0.address())
		s.Require().NoError(err)

		orch0 := s.chain.orchestrators[0]
		orch0ClientCtx, err := s.chain.clientContext("tcp://localhost:26657", orch0.keyring, "orch", orch0.address())
		s.Require().NoError(err)
		//orch1 := s.chain.orchestrators[1]
		//orch1ClientCtx, err := s.chain.clientContext("tcp://localhost:26657", orch1.keyring, "orch", orch1.address())
		//s.Require().NoError(err)

		proposer := s.chain.proposer
		proposerCtx, err := s.chain.clientContext("tcp://localhost:26657", proposer.keyring, "proposer", proposer.address())
		s.Require().NoError(err)
		propID := uint64(1)

		axelarcorkQueryClient := types.NewQueryClient(val0ClientCtx)

		////////////////
		// Happy path //
		////////////////

		arbitrumChainName := "arbitrum"
		arbitrumChainId := uint64(42161)
		proxyAddress := "0xEe75bA2C81C04DcA4b0ED6d1B7077c188FEde4d2"

		// add chain configuration
		s.T().Log("Creating AddChainConfigurationProposal")
		addChainConfigurationProp := types.AddChainConfigurationProposal{
			Title:       "add a chain configuration",
			Description: "adding an arbitrum chain config",
			ChainConfiguration: &types.ChainConfiguration{
				Name:         arbitrumChainName,
				Id:           arbitrumChainId,
				ProxyAddress: proxyAddress,
			},
		}

		addChainConfigurationPropMsg, err := govtypesv1beta1.NewMsgSubmitProposal(
			&addChainConfigurationProp,
			sdk.Coins{
				{
					Denom:  testDenom,
					Amount: math.NewInt(2),
				},
			},
			proposer.address(),
		)
		s.Require().NoError(err, "Unable to create AddChainConfigurationProposal")

		s.submitAndVoteForAxelarProposal(proposerCtx, orch0ClientCtx, propID, addChainConfigurationPropMsg)
		propID++

		s.T().Log("Verifying ChainConfiguration correctly added")
		chainConfigurationsResponse, err := axelarcorkQueryClient.QueryChainConfigurations(context.Background(), &types.QueryChainConfigurationsRequest{})
		s.Require().NoError(err)
		s.Require().Len(chainConfigurationsResponse.Configurations, 1)
		chainConfig := chainConfigurationsResponse.Configurations[0]
		s.Require().Equal(chainConfig.Name, arbitrumChainName)
		s.Require().Equal(chainConfig.Id, arbitrumChainId)
		s.Require().Equal(chainConfig.ProxyAddress, proxyAddress)

		// add managed cellar
		s.T().Log("Creating AddAxelarManagedCellarIDsProposal")
		addAxelarManagedCellarIDsProp := types.AddAxelarManagedCellarIDsProposal{
			Title:       "add an axelar cellar",
			Description: "arbitrum test cellar",
			ChainId:     arbitrumChainId,
		}

		// schedule a normal cork
		// scheduled cork proposal
		// remove managed cellar
		// upgrade proxy proposal
		// upgrade but then cancel proxy proposal
		// remove chain configuration

	})
}

func (s *IntegrationTestSuite) submitAndVoteForAxelarProposal(proposerCtx *client.Context, orchClientCtx *client.Context, propID uint64, proposalMsg *govtypesv1beta1.MsgSubmitProposal) {
	s.T().Log("Submit proposal")
	submitProposalResponse, err := s.chain.sendMsgs(*proposerCtx, proposalMsg)
	s.Require().NoError(err)
	s.Require().Zero(submitProposalResponse.Code, "raw log: %s", submitProposalResponse.RawLog)

	s.T().Log("Check proposal was submitted correctly")
	govQueryClient := govtypesv1beta1.NewQueryClient(orchClientCtx)

	s.Require().Eventually(func() bool {
		proposalsQueryResponse, err := govQueryClient.Proposals(context.Background(), &govtypesv1beta1.QueryProposalsRequest{})
		if err != nil {
			s.T().Logf("error querying proposals: %e", err)
			return false
		}

		s.Require().NotEmpty(proposalsQueryResponse.Proposals)
		s.Require().Equal(propID, proposalsQueryResponse.Proposals[propID-1].ProposalId, "not proposal id %d", propID)
		s.Require().Equal(govtypesv1beta1.StatusVotingPeriod, proposalsQueryResponse.Proposals[propID-1].Status, "proposal not in voting period")

		return true
	}, time.Second*30, time.Second*5, "proposal submission was never found")

	s.T().Log("Vote for proposal")
	for _, val := range s.chain.validators {
		kr, err := val.keyring()
		s.Require().NoError(err)
		localClientCtx, err := s.chain.clientContext("tcp://localhost:26657", &kr, "val", val.address())
		s.Require().NoError(err)

		voteMsg := govtypesv1beta1.NewMsgVote(val.address(), propID, govtypesv1beta1.OptionYes)
		voteResponse, err := s.chain.sendMsgs(*localClientCtx, voteMsg)
		s.Require().NoError(err)
		s.Require().Zero(voteResponse.Code, "Vote error: %s", voteResponse.RawLog)
	}

	s.T().Log("Waiting for proposal to be approved")
	s.Require().Eventually(func() bool {
		proposalQueryResponse, _ := govQueryClient.Proposal(context.Background(), &govtypesv1beta1.QueryProposalRequest{ProposalId: propID})
		return govtypesv1beta1.StatusPassed == proposalQueryResponse.Proposal.Status
	}, time.Second*30, time.Second*5, "proposal was never accepted")
	s.T().Log("Proposal approved!")
}
