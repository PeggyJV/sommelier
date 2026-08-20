package integration_tests

import (
	"context"
	"encoding/hex"
	"fmt"
	"time"

	"cosmossdk.io/math"
	"github.com/cosmos/cosmos-sdk/client"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	distributiontypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	govtypesv1beta1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1beta1"
	"github.com/ethereum/go-ethereum/common"
	"github.com/golang/protobuf/proto" //nolint:staticcheck
	"github.com/peggyjv/sommelier/v10/x/axelarcork/types"
	pubsubtypes "github.com/peggyjv/sommelier/v10/x/pubsub/types"
)

func (s *IntegrationTestSuite) TestAxelarCork() {
	s.Run("Test the axelarcork module", func() { ///////////
		// Setup //
		///////////

		val0 := s.chain.validators[0]
		val0kb, err := val0.keyring()
		s.Require().NoError(err)
		val0ClientCtx, err := s.chain.clientContext("tcp://localhost:26657", &val0kb, "val", val0.address())
		s.Require().NoError(err)

		orch0 := s.chain.orchestrators[0]
		orch0ClientCtx, err := s.chain.clientContext("tcp://localhost:26657", orch0.keyring, "orch", orch0.address())
		s.Require().NoError(err)

		proposer := s.chain.proposer
		proposerCtx, err := s.chain.clientContext("tcp://localhost:26657", proposer.keyring, "proposer", proposer.address())
		s.Require().NoError(err)
		propID := uint64(1)

		axelarcorkQueryClient := types.NewQueryClient(val0ClientCtx)
		pubsubQueryClient := pubsubtypes.NewQueryClient(orch0ClientCtx)
		govQueryClient := govtypesv1beta1.NewQueryClient(orch0ClientCtx)

		/////////////////////////////
		// Add chain configuration //
		/////////////////////////////

		arbitrumChainName := "arbitrum"
		arbitrumChainID := uint64(42161)
		proxyAddress := "0xEe75bA2C81C04DcA4b0ED6d1B7077c188FEde4d2"
		bridgeFees := sdk.NewCoins(sdk.NewCoin("usomm", sdk.NewIntFromUint64(33670000)))

		s.T().Log("Creating AddChainConfigurationProposal")
		addChainConfigurationProp := types.AddChainConfigurationProposal{
			Title:       "add a chain configuration",
			Description: "adding an arbitrum chain config",
			ChainConfiguration: &types.ChainConfiguration{
				Name:         arbitrumChainName,
				Id:           arbitrumChainID,
				ProxyAddress: proxyAddress,
				BridgeFees:   bridgeFees,
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
		s.Require().Equal(chainConfig.BridgeFees, bridgeFees)

		//////////////////
		// Add a cellar //
		//////////////////

		s.T().Log("Creating AddAxelarManagedCellarIDsProposal for counter contract")
		addAxelarManagedCellarIDsProp := types.AddAxelarManagedCellarIDsProposal{
			Title:       "add the counter contract as axelar cellar",
			Description: "arbitrum counter contract",
			ChainId:     arbitrumChainID,
			CellarIds: &types.CellarIDSet{
				ChainId: arbitrumChainID,
				Ids: []string{
					counterContract.Hex(),
				},
			},
			PublisherDomain: "example.com",
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

		// x/axelarcork no longer creates a pubsub default subscription when a
		// cellar is added: the modules are decoupled as of v10. Asserting one
		// exists here would be asserting the retired coupling.
		s.T().Log("Verifying no default subscription is created (axelarcork/pubsub decoupled)")
		subscriptionID := fmt.Sprintf("%d:%s", arbitrumChainID, counterContract.String())
		_, err = pubsubQueryClient.QueryDefaultSubscription(context.Background(), &pubsubtypes.QueryDefaultSubscriptionRequest{SubscriptionId: subscriptionID})
		s.Require().Error(err, "adding a managed cellar must no longer create a pubsub subscription")

		/////////////////////////////
		// Schedule an Axelar cork //
		/////////////////////////////

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
		// A single cork from the cork authority. The validator-supermajority
		// path is retired: one authorized message schedules the cork outright.
		authority := s.chain.orchestrators[0]
		authorityCtx, err := s.chain.clientContext("tcp://localhost:26657", authority.keyring, "orch", authority.address())
		s.Require().NoError(err)
		axelarCorkMsg, err := types.NewMsgScheduleAxelarCorkRequest(
			arbitrumChainID,
			ABIEncodedInc(),
			counterContract,
			deadline,
			targetBlockHeight,
			authority.address())
		s.Require().NoError(err, "Failed to construct axelar cork")
		_, err = s.chain.sendMsgs(*authorityCtx, axelarCorkMsg)
		s.Require().NoError(err, "Failed to send axelar cork from the cork authority")
		s.T().Log("Axelar cork msg sent successfully by the cork authority")

		// A non-authority orchestrator must be refused. It targets a DIFFERENT
		// height so its cork occupies a distinct store key; an identical cork
		// would collapse onto the authority's entry and prove nothing. The
		// rejection is not observable from the broadcast response (BroadcastSync
		// returns the CheckTx code, while the authority check runs in the msg
		// server at DeliverTx), so state is the witness.
		unauthorizedHeight := targetBlockHeight + 1
		other := s.chain.orchestrators[1]
		otherCtx, err := s.chain.clientContext("tcp://localhost:26657", other.keyring, "orch", other.address())
		s.Require().NoError(err)
		otherMsg, err := types.NewMsgScheduleAxelarCorkRequest(
			arbitrumChainID,
			ABIEncodedInc(),
			counterContract,
			deadline,
			unauthorizedHeight,
			other.address())
		s.Require().NoError(err, "Failed to construct axelar cork")
		_, _ = s.chain.sendMsgs(*otherCtx, otherMsg)
		s.T().Log("Non-authority axelar cork submitted; verifying it was refused")

		s.T().Log("Verifying scheduled axelar corks were created")
		corks := []*types.ScheduledAxelarCork{}
		s.Require().Eventually(func() bool {
			res, err := axelarcorkQueryClient.QueryScheduledCorks(context.Background(), &types.QueryScheduledCorksRequest{ChainId: arbitrumChainID})
			if err != nil {
				return false
			}

			if len(res.Corks) == 1 {
				corks = res.Corks
				return true
			}

			return false
		}, time.Second*30, time.Second*5, "scheduled corks never created")

		// The non-authority cork must never appear at its target height.
		s.Require().Never(func() bool {
			res, err := axelarcorkQueryClient.QueryScheduledCorksByBlockHeight(context.Background(), &types.QueryScheduledCorksByBlockHeightRequest{ChainId: arbitrumChainID, BlockHeight: unauthorizedHeight})
			if err != nil {
				return false
			}
			return len(res.Corks) > 0
		}, time.Second*8, time.Second*1,
			"a non-authority orchestrator must not be able to schedule an axelar cork")
		s.T().Log("Non-authority axelar cork correctly refused")

		s.T().Log("Checking that the cork has expected values")
		cork0 := corks[0]
		s.Require().Equal(cork0.Cork.EncodedContractCall, ABIEncodedInc())
		s.Require().Equal(cork0.Cork.ChainId, arbitrumChainID)
		s.Require().Equal(cork0.Cork.TargetContractAddress, counterContract.Hex())
		s.Require().Equal(cork0.Cork.Deadline, deadline)
		s.Require().Equal(cork0.BlockHeight, targetBlockHeight)
		s.Require().Equal(cork0.Id, axelarCorkIDHex)
		// Authority corks carry no scheduling validator.
		s.Require().Empty(cork0.Validator)

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

		// the corks are deleted when it's converted into a WinningAxelarCork and is relayable
		s.T().Log("Verifying scheduled axelar corks were deleted")
		scheduledCorksByHeightResponse, err := axelarcorkQueryClient.QueryScheduledCorksByBlockHeight(context.Background(), &types.QueryScheduledCorksByBlockHeightRequest{BlockHeight: targetBlockHeight, ChainId: arbitrumChainID})
		s.Require().NoError(err)
		s.Require().Empty(scheduledCorksByHeightResponse.Corks)

		///////////////////////////////////////////////
		// Create a governance scheduled Axelar cork //
		///////////////////////////////////////////////

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

		//////////////////////////////////////
		// Upgrade an Axelar proxy contract //
		//////////////////////////////////////

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
		propID++

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

		/////////////////////////////////////////////
		// Cancel an Axelar proxy contract upgrade //
		/////////////////////////////////////////////

		s.T().Log("Creating CancelAxelarProxyContractUpgradeProposal")
		cancelAxelarProxyContractUpgradeProp := types.CancelAxelarProxyContractUpgradeProposal{
			Title:       "cancel upgrade ofan axelar proxy contract",
			Description: "arbitrum is not getting a new proxy",
			ChainId:     arbitrumChainID,
		}

		cancelAxelarProxyContractUpgradePropMsg, err := govtypesv1beta1.NewMsgSubmitProposal(
			&cancelAxelarProxyContractUpgradeProp,
			sdk.Coins{
				{
					Denom:  testDenom,
					Amount: math.NewInt(2),
				},
			},
			proposer.address(),
		)
		s.Require().NoError(err, "Unable to create CancelAxelarProxyContractUpgradeProposal")

		s.submitAndVoteForAxelarProposal(proposerCtx, orch0ClientCtx, propID, cancelAxelarProxyContractUpgradePropMsg)
		propID++

		s.T().Log("Verifying upgrade data removed correctly")
		upgradeResponse, err = axelarcorkQueryClient.QueryAxelarProxyUpgradeData(context.Background(), &types.QueryAxelarProxyUpgradeDataRequest{})
		s.Require().NoError(err)
		s.Require().Empty(upgradeResponse.ProxyUpgradeData)

		/////////////////////
		// Remove a cellar //
		/////////////////////

		s.T().Log("Creating RemoveAxelarManagedCellarIDsProposal for counter contract")
		removeAxelarManagedCellarIDsProp := types.RemoveAxelarManagedCellarIDsProposal{
			Title:       "add the counter contract as axelar cellar",
			Description: "arbitrum counter contract",
			ChainId:     arbitrumChainID,
			CellarIds: &types.CellarIDSet{
				ChainId: arbitrumChainID,
				Ids: []string{
					counterContract.Hex(),
				},
			},
		}

		removeAxelarManagedCellarIDsPropMsg, err := govtypesv1beta1.NewMsgSubmitProposal(
			&removeAxelarManagedCellarIDsProp,
			sdk.Coins{
				{
					Denom:  testDenom,
					Amount: math.NewInt(2),
				},
			},
			proposer.address(),
		)
		s.Require().NoError(err, "Unable to create RemoveAxelarManagedCellarIDsProposal")

		s.submitAndVoteForAxelarProposal(proposerCtx, orch0ClientCtx, propID, removeAxelarManagedCellarIDsPropMsg)
		propID++

		s.T().Log("Verifying CellarID correctly removed")
		cellarIDsResponse, err = axelarcorkQueryClient.QueryCellarIDsByChainID(context.Background(), &types.QueryCellarIDsByChainIDRequest{ChainId: arbitrumChainID})
		s.Require().NoError(err)
		s.Require().Empty(cellarIDsResponse.CellarIds)

		// Still absent after removal -- nothing created it in the first place.
		s.T().Log("Verifying no default subscription exists after cellar removal")
		subscriptionID = fmt.Sprintf("%d:%s", arbitrumChainID, counterContract.String())
		_, err = pubsubQueryClient.QueryDefaultSubscription(context.Background(), &pubsubtypes.QueryDefaultSubscriptionRequest{SubscriptionId: subscriptionID})
		s.Require().Error(err)

		//////////////////////////////////
		// Remove a chain configuration //
		//////////////////////////////////

		s.T().Log("Creating RemoveChainConfigurationProposal")
		removeChainConfigurationProp := types.RemoveChainConfigurationProposal{
			Title:       "add a chain configuration",
			Description: "adding an arbitrum chain config",
			ChainId:     arbitrumChainID,
		}

		removeChainConfigurationPropMsg, err := govtypesv1beta1.NewMsgSubmitProposal(
			&removeChainConfigurationProp,
			sdk.Coins{
				{
					Denom:  testDenom,
					Amount: math.NewInt(2),
				},
			},
			proposer.address(),
		)
		s.Require().NoError(err, "Unable to create RemoveChainConfigurationProposal")

		s.submitAndVoteForAxelarProposal(proposerCtx, orch0ClientCtx, propID, removeChainConfigurationPropMsg)

		s.T().Log("Verifying ChainConfiguration correctly added")
		chainConfigurationsResponse, err = axelarcorkQueryClient.QueryChainConfigurations(context.Background(), &types.QueryChainConfigurationsRequest{})
		s.Require().NoError(err)
		s.Require().Empty(chainConfigurationsResponse.Configurations)

		//////////////////////////////////////////
		// Test module account balance sweeping //
		//////////////////////////////////////////

		// Get the baseline balance of the distribution community pool, send funds to the axelarcork module account,
		// verify they are received, and that the balance is zero on the next block.
		s.T().Log("Querying distribution community pool balance")
		distributionQueryClient := distributiontypes.NewQueryClient(orch0ClientCtx)
		distributionCommunityPoolResponse, err := distributionQueryClient.CommunityPool(context.Background(), &distributiontypes.QueryCommunityPoolRequest{})
		s.Require().NoError(err)
		initialPool := distributionCommunityPoolResponse.Pool

		// Send all of orchestrator's sweep denom and some usomm
		s.T().Log("Querying orchestrator account balances")
		bankQueryClient := banktypes.NewQueryClient(orch0ClientCtx)
		orch0AccountResponse, err := bankQueryClient.AllBalances(context.Background(), &banktypes.QueryAllBalancesRequest{Address: orch0.address().String()})
		s.Require().NoError(err)
		orch0Balances := orch0AccountResponse.Balances
		usommToSend := sdk.NewCoin(testDenom, math.NewInt(1000))
		found, sweepDenomToSend := orch0Balances.Find(axelarSweepDenom)
		s.Require().True(found, "orch0 doesn't have any sweep test denom funds")

		orch0SweepFunds := sdk.Coins{
			sweepDenomToSend,
			usommToSend,
		}

		s.T().Log("Sending funds to axelarcork module account")
		axelarcorkModuleAddress := authtypes.NewModuleAddress(types.ModuleName)
		sendFundsToAxelarcorkMsg := banktypes.NewMsgSend(
			orch0.address(),
			axelarcorkModuleAddress,
			orch0SweepFunds,
		)
		sendResponse, err := s.chain.sendMsgs(*orch0ClientCtx, sendFundsToAxelarcorkMsg)
		s.Require().NoError(err)
		s.Require().Zero(sendResponse.Code, "raw log: %s", sendResponse.RawLog)

		s.T().Log("Verifying distribution community pool balances includes the swept funds (excluding overflow)")
		poolAfterSweep := initialPool.Add(sdk.NewDecCoinsFromCoins(usommToSend)...).Add(sdk.NewDecCoinsFromCoins(sweepDenomToSend)...)
		s.Require().Eventually(func() bool {
			distributionCommunityPoolResponse, err := distributionQueryClient.CommunityPool(context.Background(), &distributiontypes.QueryCommunityPoolRequest{})
			s.Require().NoError(err)
			return poolAfterSweep.IsEqual(distributionCommunityPoolResponse.Pool)
		}, time.Second*60, time.Second*5, "swept funds never reached community pool")

		s.T().Log("Verifying overflow denom was burned")
		moduleBalanceResponse, err := bankQueryClient.Balance(context.Background(), &banktypes.QueryBalanceRequest{
			Address: axelarcorkModuleAddress.String(),
			Denom:   "overflow",
		})
		s.Require().NoError(err)
		s.Require().True(moduleBalanceResponse.Balance.IsZero(), "overflow denom was not burned")
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

		foundProps := len(proposalsQueryResponse.Proposals) > 0
		expectedID := propID == proposalsQueryResponse.Proposals[propID-1].ProposalId
		inVotingPeriod := govtypesv1beta1.StatusVotingPeriod == proposalsQueryResponse.Proposals[propID-1].Status

		return foundProps && expectedID && inVotingPeriod
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
