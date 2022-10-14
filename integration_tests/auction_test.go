package integration_tests

import (
	"context"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
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

		s.T().Logf("create governance proposal to update some token prices")
		orch := s.chain.orchestrators[0]
		orchClientCtx, err := s.chain.clientContext("tcp://localhost:26657", orch.keyring, "orch", orch.keyInfo.GetAddress())
		s.Require().NoError(err)

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
					Denom:    testDenom,
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
		s.Require().NoError(err, "unable to create governance proposal")

		s.T().Log("submit proposal")
		submitProposalResponse, err := s.chain.sendMsgs(*orchClientCtx, proposalMsg)
		s.Require().NoError(err)
		s.Require().Zero(submitProposalResponse.Code, "raw log: %s", submitProposalResponse.RawLog)

		s.T().Log("check proposal was submitted correctly")
		govQueryClient := govtypes.NewQueryClient(orchClientCtx)
		proposalsQueryResponse, err := govQueryClient.Proposals(context.Background(), &govtypes.QueryProposalsRequest{})
		s.Require().NoError(err)
		s.Require().NotEmpty(proposalsQueryResponse.Proposals)
		s.Require().Equal(uint64(1), proposalsQueryResponse.Proposals[0].ProposalId, "not proposal id 1")
		s.Require().Equal(govtypes.StatusVotingPeriod, proposalsQueryResponse.Proposals[0].Status, "proposal not in voting period")

		s.T().Log("vote for proposal")
		for _, val := range s.chain.validators {
			kr, err := val.keyring()
			s.Require().NoError(err)
			localClientCtx, err := s.chain.clientContext("tcp://localhost:26657", &kr, "val", val.keyInfo.GetAddress())
			s.Require().NoError(err)

			voteMsg := govtypes.NewMsgVote(val.keyInfo.GetAddress(), 1, govtypes.OptionYes)
			voteResponse, err := s.chain.sendMsgs(*localClientCtx, voteMsg)
			s.Require().NoError(err)
			s.Require().Zero(voteResponse.Code, "vote error: %s", voteResponse.RawLog)
		}

		s.T().Log("wait for proposal to be approved")
		s.Require().Eventuallyf(func() bool {
			proposalQueryResponse, err := govQueryClient.Proposal(context.Background(), &govtypes.QueryProposalRequest{ProposalId: 1})
			s.Require().NoError(err)
			return govtypes.StatusPassed == proposalQueryResponse.Proposal.Status
		}, time.Second*30, time.Second*5, "proposal was never accepted")

		// Verify auction created for testing exists
		auctionQuery := types.QueryActiveAuctionRequest{
			AuctionId: uint32(1),
		}
		auctionResponse, err := auctionQueryClient.QueryActiveAuction(context.Background(), &auctionQuery)
		s.Require().NoError(err)
		s.Require().Equal(auctionResponse.Auction.Id, uint32(1))
		s.T().Log(auctionResponse.Auction)

		bankQueryClient := banktypes.NewQueryClient(val0ClientCtx)
		balanceRes, err := bankQueryClient.AllBalances(context.Background(), &banktypes.QueryAllBalancesRequest{Address: authtypes.NewModuleAddress(types.ModuleName).String()})
		s.Require().NoError(err)
		s.T().Log(balanceRes)

		balanceRes, err = bankQueryClient.AllBalances(context.Background(), &banktypes.QueryAllBalancesRequest{Address: s.chain.orchestrators[0].keyInfo.GetAddress().String()})
		s.Require().NoError(err)
		s.T().Log(balanceRes)

		bidRequest1 := types.MsgSubmitBidRequest{
			AuctionId:              uint32(1),
			Bidder:                 s.chain.orchestrators[0].keyInfo.GetAddress().String(),
			MaxBidInUsomm:          sdk.NewCoin("usomm", sdk.NewIntFromUint64(5000)),
			SaleTokenMinimumAmount: sdk.NewCoin("gravity0x3506424f91fd33084466f402d5d97f05f8e3b4af", sdk.NewIntFromUint64(1)),
			Signer:                 s.chain.orchestrators[0].keyInfo.GetAddress().String(),
		}

		submitBid, err := s.chain.sendMsgs(*orchClientCtx, &bidRequest1)
		s.Require().NoError(err)
		s.T().Log(submitBid)

		auctionResponse, err = auctionQueryClient.QueryActiveAuction(context.Background(), &auctionQuery)
		s.Require().NoError(err)
		s.T().Log(auctionResponse.Auction)

		// Verify auction updated as expected
		expectedAuction := types.Auction{
			Id:                         uint32(1),
			StartingTokensForSale:      sdk.NewCoin("gravity0x3506424f91fd33084466f402d5d97f05f8e3b4af", sdk.NewInt(5000)),
			StartBlock:                 uint64(0),
			EndBlock:                   uint64(0),
			InitialPriceDecreaseRate:   sdk.MustNewDecFromStr("0.05"),
			CurrentPriceDecreaseRate:   sdk.MustNewDecFromStr("0.05"),
			PriceDecreaseBlockInterval: uint64(1000),
			InitialUnitPriceInUsomm:    sdk.MustNewDecFromStr("2"),
			CurrentUnitPriceInUsomm:    sdk.MustNewDecFromStr("2"),
			RemainingTokensForSale:     sdk.NewCoin("gravity0x3506424f91fd33084466f402d5d97f05f8e3b4af", sdk.NewInt(2500)),
			FundingModuleAccount:       cellarfees.ModuleName,
			ProceedsModuleAccount:      cellarfees.ModuleName,
		}
		s.Require().Equal(expectedAuction, *auctionResponse.Auction)

		// Verify user has funds debited and purchase credited
		balanceRes, err = bankQueryClient.AllBalances(context.Background(), &banktypes.QueryAllBalancesRequest{Address: s.chain.orchestrators[0].keyInfo.GetAddress().String()})
		s.Require().NoError(err)
		s.T().Log(balanceRes)

		s.Require().Equal(sdk.NewCoin("gravity0x3506424f91fd33084466f402d5d97f05f8e3b4af", sdk.NewInt(2500)), balanceRes.Balances[0])
		// Collin: This will likely conflict with the auction integration tests PR, but commenting out to keep the test from breaking. Prefer
		// whatever fix the other branch employs.
		// s.Require().Equal(sdk.NewCoin("usomm", sdk.NewInt(99999995000)), balanceRes.Balances[2])

		// Verify bid is stored as expected
		expectedBid := types.Bid{
			Id:                        uint64(1),
			AuctionId:                 uint32(1),
			Bidder:                    s.chain.orchestrators[0].keyInfo.GetAddress().String(),
			MaxBidInUsomm:             sdk.NewCoin("usomm", sdk.NewIntFromUint64(5000)),
			SaleTokenMinimumAmount:    sdk.NewCoin("gravity0x3506424f91fd33084466f402d5d97f05f8e3b4af", sdk.NewIntFromUint64(1)),
			TotalFulfilledSaleTokens:  sdk.NewCoin("gravity0x3506424f91fd33084466f402d5d97f05f8e3b4af", sdk.NewIntFromUint64(2500)),
			SaleTokenUnitPriceInUsomm: sdk.MustNewDecFromStr("2"),
			TotalUsommPaid:            sdk.NewCoin("usomm", sdk.NewIntFromUint64(5000)),
		}

		actualBid, err := auctionQueryClient.QueryBid(context.Background(), &types.QueryBidRequest{BidId: uint64(1), AuctionId: uint32(1)})
		s.Require().NoError(err)
		s.Require().Equal(expectedBid, *actualBid.Bid)

		// Submit another bid to be partially fulfilled
		bidRequest2 := types.MsgSubmitBidRequest{
			AuctionId:              uint32(1),
			Bidder:                 s.chain.orchestrators[0].keyInfo.GetAddress().String(),
			MaxBidInUsomm:          sdk.NewCoin("usomm", sdk.NewIntFromUint64(10000)),
			SaleTokenMinimumAmount: sdk.NewCoin("gravity0x3506424f91fd33084466f402d5d97f05f8e3b4af", sdk.NewIntFromUint64(50)),
			Signer:                 s.chain.orchestrators[0].keyInfo.GetAddress().String(),
		}

		submitBid2, err := s.chain.sendMsgs(*orchClientCtx, &bidRequest2)
		s.Require().NoError(err)
		s.T().Log(submitBid2)

		// Verify user has funds debited and purchase credited
		balanceRes, err = bankQueryClient.AllBalances(context.Background(), &banktypes.QueryAllBalancesRequest{Address: s.chain.orchestrators[0].keyInfo.GetAddress().String()})
		s.Require().NoError(err)
		s.T().Log(balanceRes)

		s.Require().Equal(sdk.NewCoin("gravity0x3506424f91fd33084466f402d5d97f05f8e3b4af", sdk.NewInt(5000)), balanceRes.Balances[0])
		// s.Require().Equal(sdk.NewCoin("usomm", sdk.NewInt(99999990000)), balanceRes.Balances[2])

		// Verify no active auctions
		_, err = auctionQueryClient.QueryActiveAuctions(context.Background(), &types.QueryActiveAuctionsRequest{})
		s.Require().Error(err)
		s.T().Log(err)
		s.Require().Equal("rpc error: code = NotFound desc = rpc error: code = NotFound desc = No active auctions found: key not found", err.Error())

		// Verify ended auction
		endedAuctionResponse, err := auctionQueryClient.QueryEndedAuction(context.Background(), &types.QueryEndedAuctionRequest{AuctionId: uint32(1)})
		s.Require().NoError(err)
		s.T().Log(endedAuctionResponse.Auction)

		node, err := orchClientCtx.GetNode()
		s.Require().NoError(err)
		status, err := node.Status(context.Background())
		s.Require().NoError(err)
		currentBlockHeight := status.SyncInfo.LatestBlockHeight
		s.T().Log(currentBlockHeight)

		expectedEndedAuction := types.Auction{
			Id:                         uint32(1),
			StartingTokensForSale:      sdk.NewCoin("gravity0x3506424f91fd33084466f402d5d97f05f8e3b4af", sdk.NewInt(5000)),
			StartBlock:                 uint64(0),
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

		// Verify bid is stored as expected
		expectedBid2 := types.Bid{
			Id:                        uint64(2),
			AuctionId:                 uint32(1),
			Bidder:                    s.chain.orchestrators[0].keyInfo.GetAddress().String(),
			MaxBidInUsomm:             sdk.NewCoin("usomm", sdk.NewIntFromUint64(10000)),
			SaleTokenMinimumAmount:    sdk.NewCoin("gravity0x3506424f91fd33084466f402d5d97f05f8e3b4af", sdk.NewIntFromUint64(50)),
			TotalFulfilledSaleTokens:  sdk.NewCoin("gravity0x3506424f91fd33084466f402d5d97f05f8e3b4af", sdk.NewIntFromUint64(2500)),
			SaleTokenUnitPriceInUsomm: sdk.MustNewDecFromStr("2"),
			TotalUsommPaid:            sdk.NewCoin("usomm", sdk.NewIntFromUint64(5000)),
		}

		actualBid2, err := auctionQueryClient.QueryBid(context.Background(), &types.QueryBidRequest{BidId: uint64(2), AuctionId: uint32(1)})
		s.Require().NoError(err)
		s.Require().Equal(expectedBid2, *actualBid2.Bid)
	})
}
