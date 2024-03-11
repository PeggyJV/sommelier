package integration_tests

import (
	"context"
	"math/big"
	"time"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	disttypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	govtypesv1beta1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1beta1"
	auctiontypes "github.com/peggyjv/sommelier/v7/x/auction/types"
	cellarfeestypes "github.com/peggyjv/sommelier/v7/x/cellarfees/types"
	cellarfeestypesv2 "github.com/peggyjv/sommelier/v7/x/cellarfees/types/v2"
)

func (s *IntegrationTestSuite) TestCellarFees() {
	s.Run("Bring up chain, submit TokenPrices, observe auction and fee distribution", func() {
		val := s.chain.validators[0]
		kb, err := val.keyring()
		s.Require().NoError(err)

		clientCtx, err := s.chain.clientContext("tcp://localhost:26657", &kb, "val", val.address())
		s.Require().NoError(err)

		auctionQueryClient := auctiontypes.NewQueryClient(clientCtx)
		bankQueryClient := banktypes.NewQueryClient(clientCtx)
		cellarfeesQueryClient := cellarfeestypesv2.NewQueryClient(clientCtx)
		distQueryClient := disttypes.NewQueryClient(clientCtx)

		s.T().Logf("Verify that the module account's fee balances are not zero")
		ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
		defer cancel()

		acctsRes, err := cellarfeesQueryClient.QueryModuleAccounts(ctx, &cellarfeestypesv2.QueryModuleAccountsRequest{})
		s.Require().NoError(err, "Failed to query module accounts")

		feesAddress := acctsRes.FeesAddress
		s.T().Logf("Fees address: %s", feesAddress)
		balanceRes, err := bankQueryClient.AllBalances(ctx, &banktypes.QueryAllBalancesRequest{
			Address: feesAddress,
		})
		s.Require().NoError(err, "Failed to query fee balance of denom %s", alphaERC20Contract.Hex())

		foundGravityDenom, foundIbcDenom := false, false
		for _, balance := range balanceRes.Balances {
			if balance.Denom == gravityDenom {
				s.Require().NotZero(balance.Amount.Uint64())
				foundGravityDenom = true
			} else if balance.Denom == ibcDenom {
				s.Require().NotZero(balance.Amount.Uint64())
				foundIbcDenom = true
			}
		}

		s.Require().True(foundGravityDenom, "fees account is missing initial gravity denom balance")
		s.Require().True(foundIbcDenom, "fees account is missing initial ibc denom balance")

		// Submit TokenPrices proposal
		orch0 := s.chain.orchestrators[0]
		orch0ClientCtx, err := s.chain.clientContext("tcp://localhost:26657", orch0.keyring, "orch", orch0.address())
		s.Require().NoError(err)
		proposer := s.chain.proposer
		proposerCtx, err := s.chain.clientContext("tcp://localhost:26657", proposer.keyring, "proposer", proposer.address())
		s.Require().NoError(err)
		propID := uint64(1)

		s.T().Log("Submitting TokenPrices proposal")
		tokenPrices := []*auctiontypes.ProposedTokenPrice{
			{
				Denom:    gravityDenom,
				Exponent: 12,
				UsdPrice: sdk.MustNewDecFromStr("10.00"),
			},
			{
				Denom:    ibcDenom,
				Exponent: 6,
				UsdPrice: sdk.MustNewDecFromStr("1.00"),
			},
		}
		addTokenPricesProp := auctiontypes.SetTokenPricesProposal{
			Title:       "add token prices",
			Description: "add token prices",
			TokenPrices: tokenPrices,
		}

		addTokenPricesPropMsg, err := govtypesv1beta1.NewMsgSubmitProposal(
			&addTokenPricesProp,
			sdk.Coins{
				{
					Denom:  testDenom,
					Amount: math.NewInt(2),
				},
			},
			proposer.address(),
		)
		s.Require().NoError(err, "Unable to create governance proposal")

		s.submitAndVoteForProposal(proposerCtx, orch0ClientCtx, propID, addTokenPricesPropMsg)

		s.T().Log("Waiting for gravity denom auction to start")
		var gravityAuctionID uint32
		var ibcAuctionID uint32
		s.Require().Eventually(func() bool {
			res, _ := auctionQueryClient.QueryActiveAuctions(ctx, &auctiontypes.QueryActiveAuctionsRequest{})

			if res != nil {
				for _, auction := range res.Auctions {
					if auction.StartingTokensForSale.Denom == gravityDenom {
						gravityAuctionID = auction.Id
						return true
					}
				}
			}

			return false
		}, time.Second*60, time.Second*5, "Auctions never started for gravity fees")

		// Send ibcDenom tokens from the orch to the fee account to trigger another auction
		orch := s.chain.orchestrators[0]
		orchClientCtx, err := s.chain.clientContext("tcp://localhost:26657", orch.keyring, "orch", orch.address())
		s.Require().NoError(err, "Failed to create client for orchestrator")

		s.T().Log("sending 1 ibc/1 to fees account to trigger auction")
		feesAcct := authtypes.NewModuleAddress(cellarfeestypes.ModuleName)
		sendRequest := banktypes.NewMsgSend(
			orch.address(),
			feesAcct,
			sdk.NewCoins(
				sdk.Coin{
					Denom:  ibcDenom,
					Amount: sdk.NewInt(1),
				},
			),
		)

		_, err = s.chain.sendMsgs(*orchClientCtx, sendRequest)
		s.Require().NoError(err, "Failed to submit send request")

		s.T().Log("Waiting for ibc denom auction to start")
		s.Require().Eventually(func() bool {
			res, _ := auctionQueryClient.QueryActiveAuctions(ctx, &auctiontypes.QueryActiveAuctionsRequest{})

			if res != nil {
				for _, auction := range res.Auctions {
					if auction.StartingTokensForSale.Denom == ibcDenom {
						ibcAuctionID = auction.Id
						return true
					}
				}
			}

			return false
		}, time.Second*120, time.Second*5, "Auctions never started for ibc fees")

		s.T().Log("Bidding to buy all of the gravity fees available")
		bidRequest1 := auctiontypes.MsgSubmitBidRequest{
			AuctionId:              gravityAuctionID,
			Signer:                 orch.address().String(),
			MaxBidInUsomm:          sdk.NewCoin(testDenom, sdk.NewIntFromUint64(1000000000000000)),
			SaleTokenMinimumAmount: sdk.NewCoin(gravityDenom, sdk.NewIntFromBigInt(big.NewInt(100000000000000))),
		}
		_, err = s.chain.sendMsgs(*orchClientCtx, &bidRequest1)
		s.Require().NoError(err, "Failed to submit gravity bid")

		s.T().Log("Bids submitted. Waiting to confirm gravity auction ended")
		s.Require().Eventually(func() bool {
			_, err := auctionQueryClient.QueryEndedAuction(ctx, &auctiontypes.QueryEndedAuctionRequest{
				AuctionId: gravityAuctionID,
			})

			// a nil error indicates the item was found
			return err == nil
		}, time.Second*10, time.Second, "Auction did not end.")

		s.T().Log("Gravity auction ended. Waiting to receive usomm in fees account")
		s.Require().Eventually(func() bool {
			res, err := bankQueryClient.Balance(ctx, &banktypes.QueryBalanceRequest{
				Address: feesAddress,
				Denom:   testDenom,
			})
			s.Require().NoError(err, "Error querying for usomm balance in fees account")
			s.Require().NotNil(res)

			s.T().Logf("usomm balance: %v", res.Balance)
			return res.Balance.Amount.GT(sdk.ZeroInt())
		}, time.Second*60, time.Second*5, "Never received usomm from auction")

		s.T().Log("usomm received! Evaluating distribution rate")
		lastBalanceSeen := sdk.ZeroInt()
		lastDiff := sdk.ZeroInt()
		s.Require().Eventually(func() bool {
			res, err := bankQueryClient.Balance(ctx, &banktypes.QueryBalanceRequest{
				Address: feesAddress,
				Denom:   testDenom,
			})
			s.Require().NoError(err, "Error querying for usomm balance in fees account")
			s.Require().NotNil(res)

			if lastBalanceSeen.Sub(res.Balance.Amount).GT(sdk.ZeroInt()) {
				if lastDiff.GT(sdk.ZeroInt()) {
					// Assuming fee distribution hasn't completed before this moment, this must be true
					// if our distribution rate is linear
					diff := lastBalanceSeen.Sub(res.Balance.Amount)
					s.Require().True(lastDiff.Equal(diff), "Observed reward distribution rate of %d usomm per block", diff.Uint64())
					return true
				}

				lastDiff = lastBalanceSeen.Sub(res.Balance.Amount)
				s.T().Logf("Observed reward distribution rate of %d usomm per block", lastDiff.Uint64())
			}

			lastBalanceSeen = res.Balance.Amount
			return false
		}, time.Second*30, time.Millisecond*400, "Distribution rate was invalid or could not be determined")

		s.T().Log("Distribution rate is nonzero. Submitting bid for all of the ibc fees available")
		bidRequest2 := auctiontypes.MsgSubmitBidRequest{
			AuctionId:              ibcAuctionID,
			Signer:                 orch.address().String(),
			MaxBidInUsomm:          sdk.NewCoin(testDenom, sdk.NewIntFromUint64(100000000)),
			SaleTokenMinimumAmount: sdk.NewCoin(ibcDenom, sdk.NewIntFromBigInt(big.NewInt(100000000))),
		}

		_, err = s.chain.sendMsgs(*orchClientCtx, &bidRequest2)
		s.Require().NoError(err, "Failed to submit ibc bid")

		s.T().Log("Waiting to confirm ibc auction ended")
		s.Require().Eventually(func() bool {
			_, err := auctionQueryClient.QueryEndedAuction(ctx, &auctiontypes.QueryEndedAuctionRequest{
				AuctionId: ibcAuctionID,
			})

			// a nil error indicates the item was found
			return err == nil
		}, time.Second*10, time.Second, "Auction did not end.")

		s.T().Log("IBC auction ended. Waiting to see distribution rate increase")
		lastBalanceSeen = sdk.ZeroInt()
		s.Require().Eventually(func() bool {
			res, err := bankQueryClient.Balance(ctx, &banktypes.QueryBalanceRequest{
				Address: feesAddress,
				Denom:   testDenom,
			})
			s.Require().NoError(err, "Error querying for usomm balance in fees account")
			s.Require().NotNil(res)

			if lastBalanceSeen.Sub(res.Balance.Amount).GT(sdk.ZeroInt()) {
				diff := lastBalanceSeen.Sub(res.Balance.Amount)
				s.T().Logf("Supply: %d usomm, Reward rate: %d usomm per block", res.Balance.Amount.Uint64(), diff.Uint64())
				if diff.GT(lastDiff) {
					return true
				}
			}

			lastBalanceSeen = res.Balance.Amount
			return false
		}, time.Second*30, time.Millisecond*400, "Distribution rate did not increase")

		s.T().Log("Distribution rate increased with supply! Getting current reward rate per validator")

		rewardRateBaseline := sdk.ZeroDec()
		rewardsRes, err := distQueryClient.DelegationRewards(ctx, &disttypes.QueryDelegationRewardsRequest{
			DelegatorAddress: val.address().String(),
			ValidatorAddress: "sommvaloper199sjfhaw3hempwzljw0lgwsm9kk6r8e5ef3hmp",
		})
		s.Require().NoError(err)

		startAmount := rewardsRes.Rewards.AmountOf(testDenom)

		// let some time elapse so we can calculate an average rate
		time.Sleep(time.Second * 12)

		rewardsRes, err = distQueryClient.DelegationRewards(ctx, &disttypes.QueryDelegationRewardsRequest{
			DelegatorAddress: val.address().String(),
			ValidatorAddress: "sommvaloper199sjfhaw3hempwzljw0lgwsm9kk6r8e5ef3hmp",
		})
		s.Require().NoError(err)
		endAmount := rewardsRes.Rewards.AmountOf(testDenom)
		rewardRate := (endAmount.Sub(startAmount).Quo(sdk.NewDec(12)))
		s.T().Logf("Baseline reward rate: %d, current validator reward rate: %d", rewardRateBaseline.RoundInt64(), rewardRate.RoundInt64())
		s.Require().True(rewardRate.GT(rewardRateBaseline), "Rewards have not increased")

		s.T().Log("Reward rate has increased. Waiting for reward supply in the fees account to be exhausted...")
		s.Require().Eventually(func() bool {
			res, err := bankQueryClient.Balance(ctx, &banktypes.QueryBalanceRequest{
				Address: feesAddress,
				Denom:   testDenom,
			})
			if err != nil {
				s.T().Log(err)
				bankQueryClient = banktypes.NewQueryClient(clientCtx)
			}

			return res != nil && res.Balance.Amount.IsZero()
		}, time.Second*300, time.Second*10, "Reward supply did not exhaust in the provided amount of time")

		s.T().Log("Confirming no auction is started...")
		for i := 0; i < 20; i++ {
			res, _ := auctionQueryClient.QueryActiveAuctions(ctx, &auctiontypes.QueryActiveAuctionsRequest{})
			if res == nil {
				continue
			}

			for _, auction := range res.Auctions {
				s.Require().NotEqual(auction.StartingTokensForSale.Denom, gravityDenom)
			}

			time.Sleep(time.Second)
		}

		s.T().Log("Done!")
	})
}
