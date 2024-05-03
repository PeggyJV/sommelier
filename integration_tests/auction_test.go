package integration_tests

import (
	"context"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	govtypesv1beta1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1beta1"
	"github.com/peggyjv/sommelier/v7/x/auction/types"
	cellarfees "github.com/peggyjv/sommelier/v7/x/cellarfees/types"
)

func (s *IntegrationTestSuite) TestAuction() {
	s.Run("Bring up chain, test governance proposal to set token prices, submit some bids, and finish an auction", func() {
		val := s.chain.validators[0]
		kb, err := val.keyring()
		s.Require().NoError(err)
		val0ClientCtx, err := s.chain.clientContext("tcp://localhost:26657", &kb, "val", val.address())
		s.Require().NoError(err)
		auctionQueryClient := types.NewQueryClient(val0ClientCtx)

		// Verify auction created for testing exists
		auctionQuery := types.QueryActiveAuctionRequest{
			AuctionId: uint32(1),
		}
		s.T().Log("Verifying expected auction exists ...")
		auctionResponse, err := auctionQueryClient.QueryActiveAuction(context.Background(), &auctionQuery)
		s.Require().NoError(err)
		s.Require().Equal(auctionResponse.Auction.Id, uint32(1))
		s.T().Log("Expected auction found!")

		s.T().Logf("Create governance proposal to update some token prices")
		proposer := s.chain.proposer
		proposerCtx, err := s.chain.clientContext("tcp://localhost:26657", proposer.keyring, "proposer", proposer.address())
		s.Require().NoError(err)
		orch := s.chain.orchestrators[0]
		orchClientCtx, err := s.chain.clientContext("tcp://localhost:26657", orch.keyring, "orch", orch.address())
		s.Require().NoError(err)

		proposal := types.SetTokenPricesProposal{
			Title:       "initial token price submission",
			Description: "our first token prices",
			TokenPrices: []*types.ProposedTokenPrice{
				{
					Denom:    "gravity0x3506424f91fd33084466f402d5d97f05f8e3b4af",
					Exponent: 6,
					UsdPrice: sdk.MustNewDecFromStr("1.0"),
				},
				{
					Denom:    "gravity0x5a98fcbea516cf06857215779fd812ca3bef1b32",
					Exponent: 6,
					UsdPrice: sdk.MustNewDecFromStr("0.25"),
				},
				{
					Denom:    testDenom,
					Exponent: 6,
					UsdPrice: sdk.MustNewDecFromStr("0.5"),
				},
			},
		}

		proposalMsg, err := govtypesv1beta1.NewMsgSubmitProposal(
			&proposal,
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
			s.Require().Equal(uint64(1), proposalsQueryResponse.Proposals[0].ProposalId, "not proposal id 1")
			s.Require().Equal(govtypesv1beta1.StatusVotingPeriod, proposalsQueryResponse.Proposals[0].Status, "proposal not in voting period")

			return true
		}, time.Second*30, time.Second*5, "proposal submission was never found")

		s.T().Log("Vote for proposal")
		for _, val := range s.chain.validators {
			kr, err := val.keyring()
			s.Require().NoError(err)
			localClientCtx, err := s.chain.clientContext("tcp://localhost:26657", &kr, "val", val.address())
			s.Require().NoError(err)

			voteMsg := govtypesv1beta1.NewMsgVote(val.address(), 1, govtypesv1beta1.OptionYes)
			voteResponse, err := s.chain.sendMsgs(*localClientCtx, voteMsg)
			s.Require().NoError(err)
			s.Require().Zero(voteResponse.Code, "Vote error: %s", voteResponse.RawLog)
		}

		s.T().Log("Waiting for proposal to be approved..")
		s.Require().Eventually(func() bool {
			proposalQueryResponse, _ := govQueryClient.Proposal(context.Background(), &govtypesv1beta1.QueryProposalRequest{ProposalId: 1})
			return govtypesv1beta1.StatusPassed == proposalQueryResponse.Proposal.Status
		}, time.Second*30, time.Second*5, "proposal was never accepted")
		s.T().Log("Proposal approved!")

		bankQueryClient := banktypes.NewQueryClient(val0ClientCtx)
		balanceRes, err := bankQueryClient.AllBalances(context.Background(), &banktypes.QueryAllBalancesRequest{Address: authtypes.NewModuleAddress(types.ModuleName).String()})
		s.Require().NoError(err)
		s.T().Logf("Auction module token balances before bids %v", balanceRes.Balances)

		bidderAddress := proposer.address().String()
		initialBidderBalanceRes, err := bankQueryClient.AllBalances(context.Background(), &banktypes.QueryAllBalancesRequest{Address: bidderAddress})
		s.Require().NoError(err)
		found, initialBidderSomm := balanceOfDenom(initialBidderBalanceRes.Balances, testDenom)
		s.Require().True(found)
		s.T().Logf("Bidder balances before bids %v", initialBidderBalanceRes.Balances)
		s.T().Log("Submitting first bid for auction")
		bidRequest1 := types.MsgSubmitBidRequest{
			AuctionId:              uint32(1),
			Signer:                 bidderAddress,
			MaxBidInUsomm:          sdk.NewCoin("usomm", sdk.NewIntFromUint64(5000000000)),
			SaleTokenMinimumAmount: sdk.NewCoin("gravity0x3506424f91fd33084466f402d5d97f05f8e3b4af", sdk.NewIntFromUint64(1)),
		}

		_, err = s.chain.sendMsgs(*proposerCtx, &bidRequest1)
		s.Require().NoError(err)
		s.T().Log("Bid submmitted successfully!")

		s.T().Log("Verifying auction updated as expected.")
		// Verify auction updated as expected
		expectedAuction := types.Auction{
			Id:                         uint32(1),
			StartingTokensForSale:      sdk.NewCoin("gravity0x3506424f91fd33084466f402d5d97f05f8e3b4af", sdk.NewInt(5000000000)),
			StartBlock:                 uint64(1),
			EndBlock:                   uint64(0),
			InitialPriceDecreaseRate:   sdk.MustNewDecFromStr("0.05"),
			CurrentPriceDecreaseRate:   sdk.MustNewDecFromStr("0.05"),
			PriceDecreaseBlockInterval: uint64(1000),
			InitialUnitPriceInUsomm:    sdk.MustNewDecFromStr("2"),
			CurrentUnitPriceInUsomm:    sdk.MustNewDecFromStr("2"),
			RemainingTokensForSale:     sdk.NewCoin("gravity0x3506424f91fd33084466f402d5d97f05f8e3b4af", sdk.NewInt(2500000000)),
			FundingModuleAccount:       cellarfees.ModuleName,
			ProceedsModuleAccount:      cellarfees.ModuleName,
		}
		s.Require().Eventually(func() bool {
			auctionResponse, err = auctionQueryClient.QueryActiveAuction(context.Background(), &auctionQuery)
			s.T().Logf("auctionResponse: %v", auctionResponse)
			if err != nil {
				return false
			}

			return expectedAuction.RemainingTokensForSale.Amount.Equal(auctionResponse.Auction.RemainingTokensForSale.Amount)
		}, time.Second*30, time.Second*5, "auction was never updated")

		// Verify user has funds debited and purchase credited
		s.T().Log("Verifying user funds debited and credited appropriately.")
		balanceRes, err = bankQueryClient.AllBalances(context.Background(), &banktypes.QueryAllBalancesRequest{Address: bidderAddress})
		s.Require().NoError(err)
		s.T().Logf("Bidder token balances after first bid %v", balanceRes.Balances)

		found, latestBidderGravity := balanceOfDenom(balanceRes.Balances, "gravity0x3506424f91fd33084466f402d5d97f05f8e3b4af")
		s.Require().True(found, "gravity denom balance not present in bidder wallet")
		s.Require().Equal(int64(2500000000), latestBidderGravity.Amount.Int64())
		found, latestBidderSomm := balanceOfDenom(balanceRes.Balances, testDenom)
		s.Require().True(found, "SOMM balance not present in bidder wallet")
		s.Require().Equal(initialBidderSomm.Amount.Sub(sdk.NewInt(5000000000)).Sub(sdk.NewInt(feeAmount)).Int64(), latestBidderSomm.Amount.Int64())
		s.T().Log("Bidder funds updated correctly!")

		node, err := orchClientCtx.GetNode()
		s.Require().NoError(err)
		status, err := node.Status(context.Background())
		s.Require().NoError(err)
		currentBlockHeight := status.SyncInfo.LatestBlockHeight

		// Verify bid is stored as expected
		expectedBid1 := types.Bid{
			Id:                        uint64(1),
			AuctionId:                 uint32(1),
			Bidder:                    bidderAddress,
			MaxBidInUsomm:             sdk.NewCoin(testDenom, sdk.NewIntFromUint64(5000000000)),
			SaleTokenMinimumAmount:    sdk.NewCoin("gravity0x3506424f91fd33084466f402d5d97f05f8e3b4af", sdk.NewIntFromUint64(1)),
			TotalFulfilledSaleTokens:  sdk.NewCoin("gravity0x3506424f91fd33084466f402d5d97f05f8e3b4af", sdk.NewIntFromUint64(2500000000)),
			SaleTokenUnitPriceInUsomm: sdk.MustNewDecFromStr("2"),
			TotalUsommPaid:            sdk.NewCoin(testDenom, sdk.NewIntFromUint64(5000000000)),
			BlockHeight:               uint64(currentBlockHeight),
		}

		s.T().Log("Verifying bid stored as expected.")
		s.Require().Eventually(func() bool {
			actualBid, err := auctionQueryClient.QueryBid(context.Background(), &types.QueryBidRequest{BidId: uint64(1), AuctionId: uint32(1)})
			if err != nil {
				return false
			}

			return expectedBid1.Bidder == actualBid.Bid.Bidder && expectedBid1.MaxBidInUsomm.Amount.Equal(actualBid.Bid.MaxBidInUsomm.Amount)
		}, time.Second*30, time.Second*5, "bid was never stored")
		s.T().Log("Bid stored correctly!")

		s.T().Log("Submitting another bid...")
		// Submit another bid to be partially fulfilled
		bidRequest2 := types.MsgSubmitBidRequest{
			AuctionId:              uint32(1),
			Signer:                 bidderAddress,
			MaxBidInUsomm:          sdk.NewCoin(testDenom, sdk.NewIntFromUint64(10000000000)),
			SaleTokenMinimumAmount: sdk.NewCoin("gravity0x3506424f91fd33084466f402d5d97f05f8e3b4af", sdk.NewIntFromUint64(50000000)),
		}

		_, err = s.chain.sendMsgs(*proposerCtx, &bidRequest2)
		s.Require().NoError(err)
		s.T().Log("Bid submitted successfully!")

		s.T().Log("Verifying expected bid stored correctly...")
		node, err = orchClientCtx.GetNode()
		s.Require().NoError(err)
		status, err = node.Status(context.Background())
		s.Require().NoError(err)
		currentBlockHeight = status.SyncInfo.LatestBlockHeight

		// Verify bid is stored as expected
		expectedBid2 := types.Bid{
			Id:                        uint64(2),
			AuctionId:                 uint32(1),
			Bidder:                    bidderAddress,
			MaxBidInUsomm:             sdk.NewCoin(testDenom, sdk.NewIntFromUint64(10000000000)),
			SaleTokenMinimumAmount:    sdk.NewCoin("gravity0x3506424f91fd33084466f402d5d97f05f8e3b4af", sdk.NewIntFromUint64(50000000)),
			TotalFulfilledSaleTokens:  sdk.NewCoin("gravity0x3506424f91fd33084466f402d5d97f05f8e3b4af", sdk.NewIntFromUint64(2500000000)),
			SaleTokenUnitPriceInUsomm: sdk.MustNewDecFromStr("2"),
			TotalUsommPaid:            sdk.NewCoin(testDenom, sdk.NewIntFromUint64(5000000000)),
			BlockHeight:               uint64(currentBlockHeight),
		}

		s.Require().Eventually(func() bool {
			actualBid2, err := auctionQueryClient.QueryBid(context.Background(), &types.QueryBidRequest{BidId: uint64(2), AuctionId: uint32(1)})
			if err != nil {
				return false
			}
			return expectedBid2.Bidder == actualBid2.Bid.Bidder && expectedBid2.MaxBidInUsomm.Amount.Equal(actualBid2.Bid.MaxBidInUsomm.Amount)
		}, time.Second*30, time.Second*5, "bid was never stored")
		s.T().Log("Bid stored correctly!")

		// Verify user has funds debited and purchase credited
		s.T().Log("Verifying bidder funds debited and credited appropriately.")
		balanceRes, err = bankQueryClient.AllBalances(context.Background(), &banktypes.QueryAllBalancesRequest{Address: bidderAddress})
		s.Require().NoError(err)
		s.T().Logf("Bidder token balances after first bid %v", balanceRes.Balances)

		found, bidderGravityBalance := balanceOfDenom(balanceRes.Balances, "gravity0x3506424f91fd33084466f402d5d97f05f8e3b4af")
		s.Require().True(found, "gravity denom balance not present in bidder wallet")
		s.Require().Equal(latestBidderGravity.Amount.Add(expectedBid2.TotalFulfilledSaleTokens.Amount).Int64(), bidderGravityBalance.Amount.Int64())
		found, bidderSommBalance := balanceOfDenom(balanceRes.Balances, testDenom)
		s.Require().True(found, "SOMM balance not present in bidder wallet")
		s.Require().Equal(latestBidderSomm.Amount.Sub(expectedBid1.TotalUsommPaid.Amount.Add(sdk.NewInt(feeAmount))).Int64(), bidderSommBalance.Amount.Int64())
		s.T().Log("Bidder funds updated correctly!")

		s.T().Log("Verifying auction has completed...")
		// Verify no active auctions
		auctions, err := auctionQueryClient.QueryActiveAuctions(context.Background(), &types.QueryActiveAuctionsRequest{})
		s.Require().NoError(err)
		s.Require().Zero(len(auctions.Auctions))
		s.T().Log("Auction completed successfully!")

		s.T().Log("Verifying ended auction stored correctly..")
		// Verify ended auction
		endedAuctionResponse, err := auctionQueryClient.QueryEndedAuction(context.Background(), &types.QueryEndedAuctionRequest{AuctionId: uint32(1)})
		s.Require().NoError(err)

		node, err = orchClientCtx.GetNode()
		s.Require().NoError(err)
		status, err = node.Status(context.Background())
		s.Require().NoError(err)
		currentBlockHeight = status.SyncInfo.LatestBlockHeight

		expectedEndedAuction := types.Auction{
			Id:                         uint32(1),
			StartingTokensForSale:      sdk.NewCoin("gravity0x3506424f91fd33084466f402d5d97f05f8e3b4af", sdk.NewInt(5000000000)),
			StartBlock:                 uint64(1),
			EndBlock:                   endedAuctionResponse.Auction.EndBlock,
			InitialPriceDecreaseRate:   sdk.MustNewDecFromStr("0.05"),
			CurrentPriceDecreaseRate:   sdk.MustNewDecFromStr("0.05"),
			PriceDecreaseBlockInterval: uint64(1000),
			InitialUnitPriceInUsomm:    sdk.MustNewDecFromStr("2"),
			CurrentUnitPriceInUsomm:    sdk.MustNewDecFromStr("2"),
			RemainingTokensForSale:     sdk.NewCoin("gravity0x3506424f91fd33084466f402d5d97f05f8e3b4af", sdk.NewInt(0)),
			FundingModuleAccount:       cellarfees.ModuleName,
			ProceedsModuleAccount:      cellarfees.ModuleName,
		}
		s.Require().Equal(expectedEndedAuction, *endedAuctionResponse.Auction)
		s.T().Log("Ended auction stored correctly!")

		s.T().Log("--Test completed successfully--")
	})
}

func balanceOfDenom(balances sdk.Coins, denom string) (found bool, balance sdk.Coin) {
	for _, balance := range balances {
		if balance.Denom == denom {
			return true, balance
		}
	}

	return false, sdk.NewCoin(denom, sdk.ZeroInt())
}
