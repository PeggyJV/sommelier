package integration_tests

import (
	"context"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	"github.com/peggyjv/sommelier/v4/x/auction/types"
)

func (s *IntegrationTestSuite) TestAuctionModule() {
	s.Run("Bring up chain, test governance proposal to set token prices, submit some bids, and finish an auction", func() {
		val := s.chain.validators[0]
		kb, err := val.keyring()
		s.Require().NoError(err)
		clientCtx, err := s.chain.clientContext("tcp://localhost:26657", &kb, "val", val.keyInfo.GetAddress())
		s.Require().NoError(err)
		auctionQueryClient := types.NewQueryClient(clientCtx)

		s.T().Logf("create governance proposal to update some token prices")
		orch := s.chain.orchestrators[0]
		clientCtx, err = s.chain.clientContext("tcp://localhost:26657", orch.keyring, "orch", orch.keyInfo.GetAddress())
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
					Denom:    types.UsommDenom,
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
		submitProposalResponse, err := s.chain.sendMsgs(*clientCtx, proposalMsg)
		s.Require().NoError(err)
		s.Require().Zero(submitProposalResponse.Code, "raw log: %s", submitProposalResponse.RawLog)

		s.T().Log("check proposal was submitted correctly")
		govQueryClient := govtypes.NewQueryClient(clientCtx)
		proposalsQueryResponse, err := govQueryClient.Proposals(context.Background(), &govtypes.QueryProposalsRequest{})
		s.Require().NoError(err)
		s.Require().NotEmpty(proposalsQueryResponse.Proposals)
		s.Require().Equal(uint64(1), proposalsQueryResponse.Proposals[0].ProposalId, "not proposal id 1")
		s.Require().Equal(govtypes.StatusVotingPeriod, proposalsQueryResponse.Proposals[0].Status, "proposal not in voting period")

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

		bankQueryClient := banktypes.NewQueryClient(clientCtx)
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

		submitBid, err := s.chain.sendMsgs(*clientCtx, &bidRequest1)
		s.Require().NoError(err)
		s.T().Log(submitBid)


		auctionResponse, err = auctionQueryClient.QueryActiveAuction(context.Background(), &auctionQuery)
		s.Require().NoError(err)
		s.T().Log(auctionResponse.Auction)

		//s.Require().Zero(submitProposalResponse.Code, "raw log: %s", submitProposalResponse.RawLog)

	})
}
