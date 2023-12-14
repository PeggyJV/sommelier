package integration_tests

import (
	"context"
	"fmt"
	"time"

	"cosmossdk.io/math"
	"github.com/cosmos/cosmos-sdk/client"
	sdk "github.com/cosmos/cosmos-sdk/types"
	govtypesv1beta1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1beta1"
	"github.com/peggyjv/sommelier/v7/x/pubsub/types"
)

func (s *IntegrationTestSuite) TestPubsub() {
	s.Run("Test the pubsub module", func() {
		// Set up validator, orchestrator, proposer, query client
		val0 := s.chain.validators[0]
		val0kb, err := val0.keyring()
		s.Require().NoError(err)
		val0ClientCtx, err := s.chain.clientContext("tcp://localhost:26657", &val0kb, "val", val0.address())
		s.Require().NoError(err)

		orch0 := s.chain.orchestrators[0]
		orch0ClientCtx, err := s.chain.clientContext("tcp://localhost:26657", orch0.keyring, "orch", orch0.address())
		s.Require().NoError(err)
		orch1 := s.chain.orchestrators[1]
		orch1ClientCtx, err := s.chain.clientContext("tcp://localhost:26657", orch1.keyring, "orch", orch1.address())
		s.Require().NoError(err)

		proposer := s.chain.proposer
		proposerCtx, err := s.chain.clientContext("tcp://localhost:26657", proposer.keyring, "proposer", proposer.address())
		s.Require().NoError(err)
		propID := uint64(1)

		pubsubQueryClient := types.NewQueryClient(val0ClientCtx)

		// add publisher (controlled by proposer)
		s.T().Log("Creating AddPublisherProposal")
		addPublisherProp := types.AddPublisherProposal{
			Title:       "add a publisher",
			Description: "example publisher",
			Domain:      "example.com",
			Address:     proposer.address().String(),
			ProofUrl:    fmt.Sprintf("https://example.com/%s/cacert.pem", proposer.address().String()),
			CaCert:      PublisherCACert,
		}

		addPublisherPropMsg, err := govtypesv1beta1.NewMsgSubmitProposal(
			&addPublisherProp,
			sdk.Coins{
				{
					Denom:  testDenom,
					Amount: math.NewInt(2),
				},
			},
			proposer.address(),
		)
		s.Require().NoError(err, "Unable to create governance proposal")

		s.submitAndVoteForProposal(proposerCtx, orch0ClientCtx, propID, addPublisherPropMsg)
		propID += 1

		s.T().Log("Verifying Publisher correctly added")
		publishersResponse, err := pubsubQueryClient.QueryPublishers(context.Background(), &types.QueryPublishersRequest{})
		s.Require().NoError(err)
		s.Require().Len(publishersResponse.Publishers, 1)
		publisher := publishersResponse.Publishers[0]
		s.Require().Equal(publisher.Address, proposer.address().String())
		s.Require().Equal(publisher.CaCert, PublisherCACert)
		s.Require().Equal(publisher.Domain, "example.com")

		// set publisher intent for cellar
		s.T().Log("Submitting PublisherIntent")
		subscriptionID := fmt.Sprintf("1:%s", unusedGenesisContract.String())
		publisherIntentMsg := types.MsgAddPublisherIntentRequest{
			PublisherIntent: &types.PublisherIntent{
				SubscriptionId:     subscriptionID,
				PublisherDomain:    publisher.Domain,
				Method:             types.PublishMethod_PUSH,
				AllowedSubscribers: types.AllowedSubscribers_VALIDATORS,
			},
			Signer: proposer.address().String(),
		}

		_, err = s.chain.sendMsgs(*proposerCtx, &publisherIntentMsg)
		s.Require().NoError(err)
		s.T().Log("PublisherIntent submitted succesfully")

		s.T().Log("Verifying PublisherIntent correctly added")
		publisherIntentsResponse, err := pubsubQueryClient.QueryPublisherIntents(context.Background(), &types.QueryPublisherIntentsRequest{})
		s.Require().NoError(err)
		s.Require().Len(publisherIntentsResponse.PublisherIntents, 1)
		publisherIntent := publisherIntentsResponse.PublisherIntents[0]
		s.Require().Equal(publisherIntent.SubscriptionId, subscriptionID)
		s.Require().Equal(publisherIntent.PublisherDomain, publisher.Domain)
		s.Require().Equal(publisherIntent.Method, types.PublishMethod_PUSH)
		s.Require().Equal(publisherIntent.AllowedSubscribers, types.AllowedSubscribers_VALIDATORS)

		// add default subscription prop
		s.T().Log("Creating AddDefaultSubscriptionProposal")
		addDefaultSubscriptionProp := types.AddDefaultSubscriptionProposal{
			Title:           "add a default subscription",
			Description:     "a default subscription!",
			SubscriptionId:  subscriptionID,
			PublisherDomain: publisher.Domain,
		}

		addDefaultSubscriptionPropMsg, err := govtypesv1beta1.NewMsgSubmitProposal(
			&addDefaultSubscriptionProp,
			sdk.Coins{
				{
					Denom:  testDenom,
					Amount: math.NewInt(2),
				},
			},
			proposer.address(),
		)
		s.Require().NoError(err, "Unable to create governance proposal")

		s.submitAndVoteForProposal(proposerCtx, orch0ClientCtx, propID, addDefaultSubscriptionPropMsg)
		propID += 1

		s.T().Log("Verifying DefaultSubscription correctly added")
		defaultSubscriptionsResponse, err := pubsubQueryClient.QueryDefaultSubscriptions(context.Background(), &types.QueryDefaultSubscriptionsRequest{})
		s.Require().NoError(err)
		s.Require().Len(defaultSubscriptionsResponse.DefaultSubscriptions, 1)
		defaultSubscription := defaultSubscriptionsResponse.DefaultSubscriptions[0]
		s.Require().Equal(defaultSubscription.SubscriptionId, subscriptionID)
		s.Require().Equal(defaultSubscription.PublisherDomain, publisher.Domain)

		// create subscribers
		s.T().Log("Creating Subscriber for two orchestrators")
		subscriber0PushURL := "https://steward.orch0.example.com:5734"
		addSubscriber0Msg := types.MsgAddSubscriberRequest{
			Subscriber: &types.Subscriber{
				Address: orch0.address().String(),
				CaCert:  SubscriberCACert,
				PushUrl: subscriber0PushURL,
			},
			Signer: orch0.address().String(),
		}

		subscriber1PushURL := "https://steward.orch1.example.com:5734"
		addSubscriber1Msg := types.MsgAddSubscriberRequest{
			Subscriber: &types.Subscriber{
				Address: orch1.address().String(),
				CaCert:  SubscriberCACert,
				PushUrl: subscriber1PushURL,
			},
			Signer: orch1.address().String(),
		}

		_, err = s.chain.sendMsgs(*orch0ClientCtx, &addSubscriber0Msg)
		s.Require().NoError(err)
		s.T().Log("AddSubscriber for orch 0 submitted correctly")

		_, err = s.chain.sendMsgs(*orch1ClientCtx, &addSubscriber1Msg)
		s.Require().NoError(err)
		s.T().Log("AddSubscriber for orch 1 submitted correctly")

		s.T().Log("Verifying Subscribers added correctly")
		subscribersResponse, err := pubsubQueryClient.QuerySubscribers(context.Background(), &types.QuerySubscribersRequest{})
		s.Require().NoError(err)
		s.Require().Len(subscribersResponse.Subscribers, 2)

		subscriber0 := subscribersResponse.Subscribers[0]
		subscriber1 := subscribersResponse.Subscribers[1]
		s.Require().Equal(subscriber0.Address, orch0.address().String())
		s.Require().Equal(subscriber0.CaCert, SubscriberCACert)
		s.Require().Equal(subscriber0.PushUrl, subscriber0PushURL)
		s.Require().Equal(subscriber1.Address, orch1.address().String())
		s.Require().Equal(subscriber1.CaCert, SubscriberCACert)
		s.Require().Equal(subscriber1.PushUrl, subscriber1PushURL)

		// subscribe to the cellar
		s.T().Log("Creating SubscriberIntent for both orchestrators")
		addSubscriberIntent0Msg := types.MsgAddSubscriberIntentRequest{
			SubscriberIntent: &types.SubscriberIntent{
				SubscriptionId:    subscriptionID,
				SubscriberAddress: orch0.address().String(),
				PublisherDomain:   publisher.Domain,
			},
			Signer: orch0.address().String(),
		}

		addSubscriberIntent1Msg := types.MsgAddSubscriberIntentRequest{
			SubscriberIntent: &types.SubscriberIntent{
				SubscriptionId:    subscriptionID,
				SubscriberAddress: orch1.address().String(),
				PublisherDomain:   publisher.Domain,
			},
			Signer: orch1.address().String(),
		}

		_, err = s.chain.sendMsgs(*orch0ClientCtx, &addSubscriberIntent0Msg)
		s.Require().NoError(err)
		s.T().Log("AddSubscriberIntent for orch 0 submitted correctly")

		_, err = s.chain.sendMsgs(*orch1ClientCtx, &addSubscriberIntent1Msg)
		s.Require().NoError(err)
		s.T().Log("AddSubscriberIntent for orch 1 submitted correctly")

		s.T().Log("Verifying SubscriberIntents added correctly")
		subscriberIntentsResponse, err := pubsubQueryClient.QuerySubscriberIntents(context.Background(), &types.QuerySubscriberIntentsRequest{})
		s.Require().NoError(err)
		s.Require().Len(subscriberIntentsResponse.SubscriberIntents, 2)

		subscriberIntent0 := subscriberIntentsResponse.SubscriberIntents[0]
		subscriberIntent1 := subscriberIntentsResponse.SubscriberIntents[1]
		s.Require().Equal(subscriberIntent0.SubscriptionId, subscriptionID)
		s.Require().Equal(subscriberIntent0.SubscriberAddress, orch0.address().String())
		s.Require().Equal(subscriberIntent0.PublisherDomain, publisher.Domain)
		s.Require().Equal(subscriberIntent1.SubscriptionId, subscriptionID)
		s.Require().Equal(subscriberIntent1.SubscriberAddress, orch1.address().String())
		s.Require().Equal(subscriberIntent1.PublisherDomain, publisher.Domain)

		// remove subscriptions to the cellar
		s.T().Log("Removing SubscriberIntent for orch 0")
		removeSubscriberIntent0Msg := types.MsgRemoveSubscriberIntentRequest{
			SubscriptionId:    subscriptionID,
			SubscriberAddress: orch0.address().String(),
			Signer:            orch0.address().String(),
		}

		_, err = s.chain.sendMsgs(*orch0ClientCtx, &removeSubscriberIntent0Msg)
		s.Require().NoError(err)
		s.T().Log("RemoveSubscriberIntent for orch 0 submitted correctly")

		s.T().Log("Verifying SubscriberIntent for orch 0 removed")
		subscriberIntentsResponse, err = pubsubQueryClient.QuerySubscriberIntents(context.Background(), &types.QuerySubscriberIntentsRequest{})
		s.Require().NoError(err)
		s.Require().Len(subscriberIntentsResponse.SubscriberIntents, 1)
		s.Require().Equal(subscriberIntentsResponse.SubscriberIntents[0].SubscriberAddress, orch1.address().String())

		s.T().Log("Removing SubscriberIntent for orch 1")
		removeSubscriberIntent1Msg := types.MsgRemoveSubscriberIntentRequest{
			SubscriptionId:    subscriptionID,
			SubscriberAddress: orch1.address().String(),
			Signer:            orch1.address().String(),
		}

		_, err = s.chain.sendMsgs(*orch1ClientCtx, &removeSubscriberIntent1Msg)
		s.Require().NoError(err)
		s.T().Log("RemoveSubscriberIntent for orch 1 submitted correctly")

		s.T().Log("Verifying SubscriberIntent for orch 1 removed")
		subscriberIntentsResponse, err = pubsubQueryClient.QuerySubscriberIntents(context.Background(), &types.QuerySubscriberIntentsRequest{})
		s.Require().NoError(err)
		s.Require().Len(subscriberIntentsResponse.SubscriberIntents, 0)

		// delete subscribers
		s.T().Log("Removing Subscriber for orch 0")
		removeSubscriber0Msg := types.MsgRemoveSubscriberRequest{
			SubscriberAddress: orch0.address().String(),
			Signer:            orch0.address().String(),
		}

		_, err = s.chain.sendMsgs(*orch0ClientCtx, &removeSubscriber0Msg)
		s.Require().NoError(err)
		s.T().Log("RemoveSubscriber for orch 0 submitted correctly")

		s.T().Log("Verifying Subscriber for orch 0 removed")
		subscribersResponse, err = pubsubQueryClient.QuerySubscribers(context.Background(), &types.QuerySubscribersRequest{})
		s.Require().NoError(err)
		s.Require().Len(subscribersResponse.Subscribers, 1)
		s.Require().Equal(subscribersResponse.Subscribers[0].Address, orch1.address().String())

		s.T().Log("Removing Subscriber for orch 1")
		removeSubscriber1Msg := types.MsgRemoveSubscriberRequest{
			SubscriberAddress: orch1.address().String(),
			Signer:            orch1.address().String(),
		}

		_, err = s.chain.sendMsgs(*orch1ClientCtx, &removeSubscriber1Msg)
		s.Require().NoError(err)
		s.T().Log("RemoveSubscriber for orch 1 submitted correctly")

		s.T().Log("Verifying Subscriber for orch 1 removed")
		subscribersResponse, err = pubsubQueryClient.QuerySubscribers(context.Background(), &types.QuerySubscribersRequest{})
		s.Require().NoError(err)
		s.Require().Len(subscribersResponse.Subscribers, 0)

		// remove default subscription prop
		s.T().Log("Creating RemoveDefaultSubscriptionProposal")
		removeDefaultSubscriptionProp := types.RemoveDefaultSubscriptionProposal{
			Title:          "remove a default subscription",
			Description:    "a default subscription is being removed!",
			SubscriptionId: subscriptionID,
		}

		removeDefaultSubscriptionPropMsg, err := govtypesv1beta1.NewMsgSubmitProposal(
			&removeDefaultSubscriptionProp,
			sdk.Coins{
				{
					Denom:  testDenom,
					Amount: math.NewInt(2),
				},
			},
			proposer.address(),
		)
		s.Require().NoError(err, "Unable to create governance proposal")

		s.submitAndVoteForProposal(proposerCtx, orch0ClientCtx, propID, removeDefaultSubscriptionPropMsg)
		propID += 1

		s.T().Log("Verifying DefaultSubscription correctly removed")
		defaultSubscriptionsResponse, err = pubsubQueryClient.QueryDefaultSubscriptions(context.Background(), &types.QueryDefaultSubscriptionsRequest{})
		s.Require().NoError(err)
		s.Require().Len(defaultSubscriptionsResponse.DefaultSubscriptions, 0)

		// remove publisher prop
		s.T().Log("Creating RemovePublisherProposal")
		removePublisherProp := types.RemovePublisherProposal{
			Title:       "remove a publisher",
			Description: "example publisher is being removed",
			Domain:      "example.com",
		}

		removePublisherPropMsg, err := govtypesv1beta1.NewMsgSubmitProposal(
			&removePublisherProp,
			sdk.Coins{
				{
					Denom:  testDenom,
					Amount: math.NewInt(2),
				},
			},
			proposer.address(),
		)
		s.Require().NoError(err, "Unable to create governance proposal")

		s.submitAndVoteForProposal(proposerCtx, orch0ClientCtx, propID, removePublisherPropMsg)
		propID += 1

		s.T().Log("Verifying Publisher correctly removed")
		publishersResponse, err = pubsubQueryClient.QueryPublishers(context.Background(), &types.QueryPublishersRequest{})
		s.Require().NoError(err)
		s.Require().Len(publishersResponse.Publishers, 0)
	})
}

func (s *IntegrationTestSuite) submitAndVoteForProposal(proposerCtx *client.Context, orchClientCtx *client.Context, propID uint64, proposalMsg *govtypesv1beta1.MsgSubmitProposal) {
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

const PublisherCACert = `-----BEGIN CERTIFICATE-----
MIICGzCCAaKgAwIBAgIUVYhZ4+pC7vQAf5FC6pssLk/eq5YwCgYIKoZIzj0EAwMw
RTELMAkGA1UEBhMCQVUxEzARBgNVBAgMClNvbWUtU3RhdGUxITAfBgNVBAoMGElu
dGVybmV0IFdpZGdpdHMgUHR5IEx0ZDAeFw0yMjAxMDUwNzIwMzFaFw0yNDAxMDUw
NzIwMzFaMEUxCzAJBgNVBAYTAkFVMRMwEQYDVQQIDApTb21lLVN0YXRlMSEwHwYD
VQQKDBhJbnRlcm5ldCBXaWRnaXRzIFB0eSBMdGQwdjAQBgcqhkjOPQIBBgUrgQQA
IgNiAAQ3jwZd0Xe9w55UyAxRuc4F2u/LDdo7ykCZBO34neXpLR4GRRpx5VjFdHcX
WjvM9j3DnWjptb1fe7TIKSSJRmW1skWkpktOthIPhfga9jBhU4WRUDloKk1tRuiI
e8rRSlSjUzBRMB0GA1UdDgQWBBSTyTULHT9hNAA2Wg4dCtuTuIhiXTAfBgNVHSME
GDAWgBSTyTULHT9hNAA2Wg4dCtuTuIhiXTAPBgNVHRMBAf8EBTADAQH/MAoGCCqG
SM49BAMDA2cAMGQCMEd+Eg6lhStLkWEwmJJGN3Xdh9JmNsgsdff3mI3Y7UmHOB8K
HOqHGS8ApZcunRauDAIwRtgceZpkS92KuP3QOUotAH/nnCzp7X1lVzGOSTBRTVYJ
pohf4PJrfacqpi7PoXBk
-----END CERTIFICATE-----
`

const SubscriberCACert = `-----BEGIN CERTIFICATE-----
MIICHTCCAaKgAwIBAgIUTYD5x0zSg1rOztoJK8OEgWDl+yYwCgYIKoZIzj0EAwMw
RTELMAkGA1UEBhMCQVUxEzARBgNVBAgMClNvbWUtU3RhdGUxITAfBgNVBAoMGElu
dGVybmV0IFdpZGdpdHMgUHR5IEx0ZDAeFw0yMjAxMDUwNzIwMjlaFw0yNDAxMDUw
NzIwMjlaMEUxCzAJBgNVBAYTAkFVMRMwEQYDVQQIDApTb21lLVN0YXRlMSEwHwYD
VQQKDBhJbnRlcm5ldCBXaWRnaXRzIFB0eSBMdGQwdjAQBgcqhkjOPQIBBgUrgQQA
IgNiAATi4OkAJaqyWwS1F6mBCBftwF/K02Zl07pg2C/WJxZEaGI/cRVTELt4Qsy2
7SiGcJLIIsTQXfdNkyRue20/J/SpUDPMVbWNCozC2bS4DWd1n9uHlSMT4h7gZqxf
SkkkecCjUzBRMB0GA1UdDgQWBBSngShmDy8kt2azMqFGD1ObYaXT0DAfBgNVHSME
GDAWgBSngShmDy8kt2azMqFGD1ObYaXT0DAPBgNVHRMBAf8EBTADAQH/MAoGCCqG
SM49BAMDA2kAMGYCMQCel/W4B/LB75j0WHEHrKSoED17D4w+OrXlK6wnpVRSyOmZ
A0B4pBO4uh3ldwCZnBACMQC0whN1TI8a9Ku90nfvZ+D2kKMg/p39SmCDadQJNzwc
kp4YI2VJp0zYzt/xLiBRbZc=
-----END CERTIFICATE-----
`
