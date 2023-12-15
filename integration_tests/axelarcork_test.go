package integration_tests

import (
	"context"
	"encoding/hex"
	"sort"
	"time"

	"cosmossdk.io/math"
	"github.com/cosmos/cosmos-sdk/client"
	sdk "github.com/cosmos/cosmos-sdk/types"
	govtypesv1beta1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1beta1"
	"github.com/ethereum/go-ethereum/common"
	"github.com/golang/protobuf/proto"
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

		sortedValidators := make([]string, 0, 4)
		for _, validator := range s.chain.validators {
			sortedValidators = append(sortedValidators, validator.validatorAddress().String())
		}
		sort.Sort(sort.StringSlice(sortedValidators))

		axelarcorkQueryClient := types.NewQueryClient(val0ClientCtx)
		govQueryClient := govtypesv1beta1.NewQueryClient(orch0ClientCtx)

		////////////////
		// Happy path //
		////////////////

		arbitrumChainName := "arbitrum"
		arbitrumChainID := uint64(42161)
		proxyAddress := "0xEe75bA2C81C04DcA4b0ED6d1B7077c188FEde4d2"

		// add chain configuration
		s.T().Log("Creating AddChainConfigurationProposal")
		addChainConfigurationProp := types.AddChainConfigurationProposal{
			Title:       "add a chain configuration",
			Description: "adding an arbitrum chain config",
			ChainConfiguration: &types.ChainConfiguration{
				Name:         arbitrumChainName,
				Id:           arbitrumChainID,
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
		s.Require().Equal(chainConfig.Id, arbitrumChainID)
		s.Require().Equal(chainConfig.ProxyAddress, proxyAddress)

		// add managed cellar
		s.T().Log("Creating AddAxelarManagedCellarIDsProposal for counter contract")
		addAxelarManagedCellarIDsProp := types.AddAxelarManagedCellarIDsProposal{
			Title:       "add the counter contract as axelar cellar",
			Description: "arbitrum counter contract",
			ChainId:     arbitrumChainID,
			CellarIds: &types.CellarIDSet{
				Ids: []string{
					counterContract.Hex(),
				},
			},
		}

		addAxelarManagedCellarIDsPropMsg, err := govtypesv1beta1.NewMsgSubmitProposal(
			&addAxelarManagedCellarIDsProp,
			sdk.Coins{
				{
					Denom:  testDenom,
					Amount: math.NewInt(2),
				},
			},
			proposer.address(),
		)
		s.Require().NoError(err, "Unable to create AddAxelarManagedCellarIDsProposal")

		s.submitAndVoteForAxelarProposal(proposerCtx, orch0ClientCtx, propID, addAxelarManagedCellarIDsPropMsg)
		propID++

		s.T().Log("Verifying CellarID correctly added")
		cellarIDsResponse, err := axelarcorkQueryClient.QueryCellarIDsByChainID(context.Background(), &types.QueryCellarIDsByChainIDRequest{ChainId: arbitrumChainID})
		s.Require().NoError(err)
		s.Require().Len(cellarIDsResponse.CellarIds, 1)
		s.Require().Equal(cellarIDsResponse.CellarIds[0], counterContract.Hex())

		s.T().Log("Schedule an axelar cork for the future")
		node, err := proposerCtx.GetNode()
		s.Require().NoError(err)
		status, err := node.Status(context.Background())
		s.Require().NoError(err)
		currentBlockHeight := status.SyncInfo.LatestBlockHeight
		targetBlockHeight := uint64(currentBlockHeight + 15)
		deadline := uint64(time.Now().Unix() + int64(time.Hour.Seconds()))

		s.T().Logf("Scheduling axelar cork calls for height %d", targetBlockHeight)
		axelarCork := types.AxelarCork{
			EncodedContractCall:   ABIEncodedInc(),
			ChainId:               arbitrumChainID,
			TargetContractAddress: counterContract.Hex(),
			Deadline:              deadline,
		}
		axelarCorkID := axelarCork.IDHash(targetBlockHeight)
		axelarCorkIDHex := hex.EncodeToString(axelarCorkID)
		s.T().Logf("Axelar cork ID is %s", axelarCorkIDHex)
		for i, orch := range s.chain.orchestrators {
			i := i
			orch := orch
			clientCtx, err := s.chain.clientContext("tcp://localhost:26657", orch.keyring, "orch", orch.address())
			s.Require().NoError(err)
			axelarCorkMsg, err := types.NewMsgScheduleAxelarCorkRequest(
				arbitrumChainID,
				ABIEncodedInc(),
				counterContract,
				deadline,
				targetBlockHeight,
				orch.address())
			s.Require().NoError(err, "Failed to construct axelar cork")
			response, err := s.chain.sendMsgs(*clientCtx, axelarCorkMsg)
			s.Require().NoError(err, "Failed to send axelar cork to node %d", i)
			if response.Code != 0 {
				if response.Code != 32 {
					s.T().Log(response)
				}
			}

			s.T().Logf("Axelar cork msg for orch %d sent successfully", i)
		}

		s.T().Log("Verifying scheduled axelar corks were created")
		scheduledCorksResponse, err := axelarcorkQueryClient.QueryScheduledCorks(context.Background(), &types.QueryScheduledCorksRequest{ChainId: arbitrumChainID})
		s.Require().NoError(err)
		s.Require().Len(scheduledCorksResponse.Corks, 4)
		cork0 := scheduledCorksResponse.Corks[0]
		cork1 := scheduledCorksResponse.Corks[1]
		cork2 := scheduledCorksResponse.Corks[2]
		cork3 := scheduledCorksResponse.Corks[3]
		s.Require().Equal(cork0.Cork.EncodedContractCall, ABIEncodedInc())
		s.Require().Equal(cork0.Cork.ChainId, arbitrumChainID)
		s.Require().Equal(cork0.Cork.TargetContractAddress, counterContract.Hex())
		s.Require().Equal(cork0.Cork.Deadline, deadline)
		s.Require().Equal(cork0.BlockHeight, targetBlockHeight)
		s.Require().Equal(cork0.Id, axelarCorkIDHex)
		s.Require().Equal(cork1.Cork.EncodedContractCall, ABIEncodedInc())
		s.Require().Equal(cork1.Cork.ChainId, arbitrumChainID)
		s.Require().Equal(cork1.Cork.TargetContractAddress, counterContract.Hex())
		s.Require().Equal(cork1.Cork.Deadline, deadline)
		s.Require().Equal(cork1.BlockHeight, targetBlockHeight)
		s.Require().Equal(cork1.Id, axelarCorkIDHex)
		s.Require().Equal(cork2.Cork.EncodedContractCall, ABIEncodedInc())
		s.Require().Equal(cork2.Cork.ChainId, arbitrumChainID)
		s.Require().Equal(cork2.Cork.TargetContractAddress, counterContract.Hex())
		s.Require().Equal(cork2.Cork.Deadline, deadline)
		s.Require().Equal(cork2.BlockHeight, targetBlockHeight)
		s.Require().Equal(cork2.Id, axelarCorkIDHex)
		s.Require().Equal(cork3.Cork.EncodedContractCall, ABIEncodedInc())
		s.Require().Equal(cork3.Cork.ChainId, arbitrumChainID)
		s.Require().Equal(cork3.Cork.TargetContractAddress, counterContract.Hex())
		s.Require().Equal(cork3.Cork.Deadline, deadline)
		s.Require().Equal(cork3.BlockHeight, targetBlockHeight)
		s.Require().Equal(cork3.Id, axelarCorkIDHex)

		corkValidators := []string{cork0.Validator, cork1.Validator, cork2.Validator, cork3.Validator}
		sort.Sort(sort.StringSlice(corkValidators))
		s.Require().Equal(corkValidators, sortedValidators)

		s.T().Log("Waiting for scheduled height")
		s.Require().Eventuallyf(func() bool {
			node, err := val0ClientCtx.GetNode()
			s.Require().NoError(err)
			status, err := node.Status(context.Background())
			s.Require().NoError(err)

			currentHeight := status.SyncInfo.LatestBlockHeight
			if currentHeight > int64(targetBlockHeight+1) {
				return true
			} else if currentHeight < int64(targetBlockHeight) {
				scheduledCorksResponse, err := axelarcorkQueryClient.QueryScheduledCorks(context.Background(), &types.QueryScheduledCorksRequest{ChainId: arbitrumChainID})
				if err != nil {
					s.T().Logf("error: %s", err)
					return false
				}

				// verify that the scheduled corks have not yet been consumed
				s.Require().Len(scheduledCorksResponse.Corks, len(s.chain.validators))
			}

			return false
		}, 3*time.Minute, 1*time.Second, "never reached scheduled height")

		s.T().Log("Verifying axelar cork was approved")
		corkResultResponse, err := axelarcorkQueryClient.QueryCorkResult(context.Background(), &types.QueryCorkResultRequest{Id: axelarCorkIDHex, ChainId: arbitrumChainID})
		s.Require().NoError(err)
		s.Require().True(corkResultResponse.CorkResult.Approved)
		s.Require().True(sdk.MustNewDecFromStr(corkResultResponse.CorkResult.ApprovalPercentage).GT(corkVoteThreshold))
		s.Require().Equal(counterContract, common.HexToAddress(corkResultResponse.CorkResult.Cork.TargetContractAddress))

		s.T().Log("Verifying scheduled axelar corks were deleted")
		scheduledCorksByHeightResponse, err := axelarcorkQueryClient.QueryScheduledCorksByBlockHeight(context.Background(), &types.QueryScheduledCorksByBlockHeightRequest{BlockHeight: targetBlockHeight, ChainId: arbitrumChainID})
		s.Require().NoError(err)
		s.Require().Len(scheduledCorksByHeightResponse.Corks, 0)

		protoJSON := "{\"cellar_id\":\"0x123801a7D398351b8bE11C439e05C5B3259aeC9B\",\"cellar_v1\":{\"some_fuction\":{\"function_args\":{}},\"block_height\":12345}}"
		s.T().Log("Creating AxelarScheduledCorkProposal for counter contract")
		addAxelarScheduledCorkProp := types.AxelarScheduledCorkProposal{
			Title:                 "schedule a cork to the counter contract",
			Description:           "arbitrum counter contract scheduled cork",
			BlockHeight:           targetBlockHeight,
			ChainId:               arbitrumChainID,
			TargetContractAddress: counterContract.Hex(),
			ContractCallProtoJson: protoJSON,
			Deadline:              deadline,
		}

		addAxelarScheduledCorkPropMsg, err := govtypesv1beta1.NewMsgSubmitProposal(
			&addAxelarScheduledCorkProp,
			sdk.Coins{
				{
					Denom:  testDenom,
					Amount: math.NewInt(2),
				},
			},
			proposer.address(),
		)
		s.Require().NoError(err, "Unable to create AxelarScheduledCorkProposal")

		s.submitAndVoteForAxelarProposal(proposerCtx, orch0ClientCtx, propID, addAxelarScheduledCorkPropMsg)

		s.T().Log("Verifying the details of the AxelarScheduledCorkProposal proposal")
		proposalResponse, err := govQueryClient.Proposal(context.Background(), &govtypesv1beta1.QueryProposalRequest{ProposalId: propID})
		s.Require().NoError(err)
		propContent := &types.AxelarScheduledCorkProposal{}
		err = proto.Unmarshal(proposalResponse.Proposal.Content.Value, propContent)
		s.Require().NoError(err, "couldn't unmarshal into proposal")
		s.Require().Equal(propContent.Title, addAxelarScheduledCorkProp.Title)
		s.Require().Equal(propContent.Description, addAxelarScheduledCorkProp.Description)
		s.Require().Equal(propContent.BlockHeight, addAxelarScheduledCorkProp.BlockHeight)
		s.Require().Equal(propContent.ChainId, addAxelarScheduledCorkProp.ChainId)
		s.Require().Equal(propContent.TargetContractAddress, addAxelarScheduledCorkProp.TargetContractAddress)
		s.Require().Equal(propContent.ContractCallProtoJson, addAxelarScheduledCorkProp.ContractCallProtoJson)
		s.Require().Equal(propContent.Deadline, addAxelarScheduledCorkProp.Deadline)
		propID++

		s.T().Log("Creating UpgradeAxelarProxyContractProposal")
		newProxyAddress := "0x438087f7c226A89762a791F187d7c3D4a0e95ae6"
		upgradeAxelarProxyContractProp := types.UpgradeAxelarProxyContractProposal{
			Title:           "upgrade an axelar proxy contract",
			Description:     "arbitrum is getting a new proxy",
			ChainId:         arbitrumChainID,
			NewProxyAddress: newProxyAddress,
		}

		upgradeAxelarProxyContractPropMsg, err := govtypesv1beta1.NewMsgSubmitProposal(
			&upgradeAxelarProxyContractProp,
			sdk.Coins{
				{
					Denom:  testDenom,
					Amount: math.NewInt(2),
				},
			},
			proposer.address(),
		)
		s.Require().NoError(err, "Unable to create UpgradeAxelarProxyContractProposal")

		s.submitAndVoteForAxelarProposal(proposerCtx, orch0ClientCtx, propID, upgradeAxelarProxyContractPropMsg)
		//propID++

		s.T().Log("Verifying upgrade data added correctly")
		node, err = val0ClientCtx.GetNode()
		s.Require().NoError(err)
		status, err = node.Status(context.Background())
		s.Require().NoError(err)

		threshold := int64(types.DefaultExecutableHeightThreshold)
		currentHeight := status.SyncInfo.LatestBlockHeight
		upgradeResponse, err := axelarcorkQueryClient.QueryAxelarProxyUpgradeData(context.Background(), &types.QueryAxelarProxyUpgradeDataRequest{})
		s.Require().NoError(err)
		s.Require().Len(upgradeResponse.ProxyUpgradeData, 1)
		upgradeData := upgradeResponse.ProxyUpgradeData[0]
		s.Require().Equal(upgradeData.ChainId, arbitrumChainID)
		// an approximation but timing is difficult
		s.Require().Less(upgradeData.ExecutableHeightThreshold, currentHeight+threshold+5)
		s.Require().Greater(upgradeData.ExecutableHeightThreshold, currentHeight+threshold-5)
		encodedProxy, _, err := types.DecodeUpgradeArgs(upgradeData.Payload)
		s.Require().NoError(err)
		s.Require().Equal(encodedProxy, newProxyAddress)

		// cancel proxy upgrade
		// remove managed cellar
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
