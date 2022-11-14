package integration_tests

import (
	"context"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	"github.com/peggyjv/sommelier/v4/app/params"
	"github.com/peggyjv/sommelier/v4/x/auction/types"
	cellarfees "github.com/peggyjv/sommelier/v4/x/cellarfees/types"
)

func (s *IntegrationTestSuite) TestAuctionModule() {
	s.Run("Bring up chain, test governance proposal to set token prices, submit some bids, and finish an auction", func() {
		val := s.chain.validators[0]
		kb, err := val.keyring()
		s.Require().NoError(err)
		val0ClientCtx, err := s.chain.clientContext("tcp://localhost:26657", &kb, "val", val.keyInfo.GetAddress())
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
		orch := s.chain.orchestrators[0]
		orchClientCtx, err := s.chain.clientContext("tcp://localhost:26657", orch.keyring, "orch", orch.keyInfo.GetAddress())
		s.Require().NoError(err)
		orch0Address := s.chain.orchestrators[0].keyInfo.GetAddress().String()

		proposal := types.SetTokenPricesProposal{
			Title:       "initial token price submission",
			Description: "our first token prices",
			TokenPrices: []*types.ProposedTokenPrice{
				{
					Denom:    "gravity0x3506424f91fd33084466f402d5d97f05f8e3b4af",
					UsdPrice: sdk.MustNewDecFromStr("1.0"),
				},
				{
					Denom:    "gravity0x5a98fcbea516cf06857215779fd812ca3bef1b32",
					UsdPrice: sdk.MustNewDecFromStr("0.25"),
				},
				{
					Denom:    params.BaseCoinUnit,
					UsdPrice: sdk.MustNewDecFromStr("0.5"),
				},
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
		s.Require().NoError(err, "Unable to create governance proposal")

		s.T().Log("Submit proposal")
		submitProposalResponse, err := s.chain.sendMsgs(*orchClientCtx, proposalMsg)
		s.Require().NoError(err)
		s.Require().Zero(submitProposalResponse.Code, "raw log: %s", submitProposalResponse.RawLog)

		s.T().Log("Check proposal was submitted correctly")
		govQueryClient := govtypes.NewQueryClient(orchClientCtx)

		s.Require().Eventually(func() bool {
			proposalsQueryResponse, err := govQueryClient.Proposals(context.Background(), &govtypes.QueryProposalsRequest{})
			if err != nil {
				s.T().Logf("error querying proposals: %e", err)
				return false
			}

			s.Require().NotEmpty(proposalsQueryResponse.Proposals)
			s.Require().Equal(uint64(1), proposalsQueryResponse.Proposals[0].ProposalId, "not proposal id 1")
			s.Require().Equal(govtypes.StatusVotingPeriod, proposalsQueryResponse.Proposals[0].Status, "proposal not in voting period")

			return true
		}, time.Second*30, time.Second*5, "proposal submission was never found")

		s.T().Log("Vote for proposal")
		for _, val := range s.chain.validators {
			kr, err := val.keyring()
			s.Require().NoError(err)
			localClientCtx, err := s.chain.clientContext("tcp://localhost:26657", &kr, "val", val.keyInfo.GetAddress())
			s.Require().NoError(err)

			voteMsg := govtypes.NewMsgVote(val.keyInfo.GetAddress(), 1, govtypes.OptionYes)
			voteResponse, err := s.chain.sendMsgs(*localClientCtx, voteMsg)
			s.Require().NoError(err)
			s.Require().Zero(voteResponse.Code, "Vote error: %s", voteResponse.RawLog)
		}

		s.T().Log("Waiting for proposal to be approved..")
		s.Require().Eventually(func() bool {
			proposalQueryResponse, err := govQueryClient.Proposal(context.Background(), &govtypes.QueryProposalRequest{ProposalId: 1})
			s.Require().NoError(err)
			return govtypes.StatusPassed == proposalQueryResponse.Proposal.Status
		}, time.Second*30, time.Second*5, "proposal was never accepted")
		s.T().Log("Proposal approved!")

		bankQueryClient := banktypes.NewQueryClient(val0ClientCtx)
		balanceRes, err := bankQueryClient.AllBalances(context.Background(), &banktypes.QueryAllBalancesRequest{Address: authtypes.NewModuleAddress(types.ModuleName).String()})
		s.Require().NoError(err)
		s.T().Logf("Auction module token balances before bids %v", balanceRes.Balances)

		initialOrchBalanceRes, err := bankQueryClient.AllBalances(context.Background(), &banktypes.QueryAllBalancesRequest{Address: orch0Address})
		s.Require().NoError(err)
		s.T().Logf("Orchestrator 0 balances before bids %v", initialOrchBalanceRes.Balances)

		s.T().Log("Submitting first bid for auction")
		bidRequest1 := types.MsgSubmitBidRequest{
			AuctionId:              uint32(1),
			Signer:                 orch0Address,
			MaxBidInUsomm:          sdk.NewCoin("usomm", sdk.NewIntFromUint64(5000000000)),
			SaleTokenMinimumAmount: sdk.NewCoin("gravity0x3506424f91fd33084466f402d5d97f05f8e3b4af", sdk.NewIntFromUint64(1)),
		}

		_, err = s.chain.sendMsgs(*orchClientCtx, &bidRequest1)
		s.Require().NoError(err)
		s.T().Log("Bid submmitted successfully!")

		auctionResponse, err = auctionQueryClient.QueryActiveAuction(context.Background(), &auctionQuery)
		s.Require().NoError(err)

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
		s.Require().Equal(expectedAuction, *auctionResponse.Auction)
		s.T().Log("Auction updated correctly!")

		// Verify user has funds debited and purchase credited
		s.T().Log("Verifying user funds debited and credited appropriately.")
		balanceRes, err = bankQueryClient.AllBalances(context.Background(), &banktypes.QueryAllBalancesRequest{Address: orch0Address})
		s.Require().NoError(err)
		s.T().Logf("Orchestrator 0 token balances after first bid %v", balanceRes.Balances)

		s.Require().Equal(sdk.NewCoin("gravity0x3506424f91fd33084466f402d5d97f05f8e3b4af", sdk.NewInt(2500000000)), balanceRes.Balances[0])
		s.Require().Equal(sdk.NewCoin("usomm", initialOrchBalanceRes.Balances[1].Amount.Sub(sdk.NewInt(5000000000))), balanceRes.Balances[2])
		s.T().Log("User funds updated correctly!")

		node, err := orchClientCtx.GetNode()
		s.Require().NoError(err)
		status, err := node.Status(context.Background())
		s.Require().NoError(err)
		currentBlockHeight := status.SyncInfo.LatestBlockHeight

		// Verify bid is stored as expected
		expectedBid1 := types.Bid{
			Id:                        uint64(1),
			AuctionId:                 uint32(1),
			Bidder:                    orch0Address,
			MaxBidInUsomm:             sdk.NewCoin("usomm", sdk.NewIntFromUint64(5000000000)),
			SaleTokenMinimumAmount:    sdk.NewCoin("gravity0x3506424f91fd33084466f402d5d97f05f8e3b4af", sdk.NewIntFromUint64(1)),
			TotalFulfilledSaleTokens:  sdk.NewCoin("gravity0x3506424f91fd33084466f402d5d97f05f8e3b4af", sdk.NewIntFromUint64(2500000000)),
			SaleTokenUnitPriceInUsomm: sdk.MustNewDecFromStr("2"),
			TotalUsommPaid:            sdk.NewCoin("usomm", sdk.NewIntFromUint64(5000000000)),
			BlockHeight:               uint64(currentBlockHeight),
		}

		s.T().Log("Verifying bid stored as expected.")
		actualBid, err := auctionQueryClient.QueryBid(context.Background(), &types.QueryBidRequest{BidId: uint64(1), AuctionId: uint32(1)})
		s.Require().NoError(err)
		s.Require().Equal(expectedBid1, *actualBid.Bid)
		s.T().Log("Bid stored correctly!")

		s.T().Log("Submitting another bid...")
		// Submit another bid to be partially fulfilled
		bidRequest2 := types.MsgSubmitBidRequest{
			AuctionId:              uint32(1),
			Signer:                 orch0Address,
			MaxBidInUsomm:          sdk.NewCoin("usomm", sdk.NewIntFromUint64(10000000000)),
			SaleTokenMinimumAmount: sdk.NewCoin("gravity0x3506424f91fd33084466f402d5d97f05f8e3b4af", sdk.NewIntFromUint64(50000000)),
		}

		_, err = s.chain.sendMsgs(*orchClientCtx, &bidRequest2)
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
			Bidder:                    orch0Address,
			MaxBidInUsomm:             sdk.NewCoin("usomm", sdk.NewIntFromUint64(10000000000)),
			SaleTokenMinimumAmount:    sdk.NewCoin("gravity0x3506424f91fd33084466f402d5d97f05f8e3b4af", sdk.NewIntFromUint64(50000000)),
			TotalFulfilledSaleTokens:  sdk.NewCoin("gravity0x3506424f91fd33084466f402d5d97f05f8e3b4af", sdk.NewIntFromUint64(2500000000)),
			SaleTokenUnitPriceInUsomm: sdk.MustNewDecFromStr("2"),
			TotalUsommPaid:            sdk.NewCoin("usomm", sdk.NewIntFromUint64(5000000000)),
			BlockHeight:               uint64(currentBlockHeight),
		}

		actualBid2, err := auctionQueryClient.QueryBid(context.Background(), &types.QueryBidRequest{BidId: uint64(2), AuctionId: uint32(1)})
		s.Require().NoError(err)
		s.Require().Equal(expectedBid2, *actualBid2.Bid)
		s.T().Log("Bid stored correctly!")

		s.T().Log("Verifying user funds debited and credited appropriately.")
		// Verify user has funds debited and purchase credited
		balanceRes, err = bankQueryClient.AllBalances(context.Background(), &banktypes.QueryAllBalancesRequest{Address: orch0Address})
		s.Require().NoError(err)
		s.T().Logf("Orchestrator 0 token balances after second bid %v", balanceRes.Balances)

		expectedSaleTokens := expectedBid1.TotalFulfilledSaleTokens.Amount.Add(expectedBid2.TotalFulfilledSaleTokens.Amount)
		expectedUsommRemaining := initialOrchBalanceRes.Balances[1].Amount.Sub(expectedBid1.TotalUsommPaid.Amount.Add(expectedBid2.TotalUsommPaid.Amount))
		s.Require().Equal(sdk.NewCoin("gravity0x3506424f91fd33084466f402d5d97f05f8e3b4af", expectedSaleTokens), balanceRes.Balances[0])
		s.Require().Equal(sdk.NewCoin("usomm", expectedUsommRemaining), balanceRes.Balances[2])
		s.T().Log("User funds updated correctly!")

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
			EndBlock:                   uint64(currentBlockHeight),
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
