package integration_tests

import (
	"context"
	"fmt"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	disttypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	auctiontypes "github.com/peggyjv/sommelier/v7/x/auction/types"
	cellarfeestypesv2 "github.com/peggyjv/sommelier/v7/x/cellarfees/types/v2"
	corktypes "github.com/peggyjv/sommelier/v7/x/cork/types/v2"
)

func (s *IntegrationTestSuite) TestCellarFees() {
	s.Run("Bring up chain, send fees from ethereum, observe auction and fee distribution", func() {
		val := s.chain.validators[0]
		ethereumSender := val.ethereumKey.address
		kb, err := val.keyring()
		s.Require().NoError(err)

		clientCtx, err := s.chain.clientContext("tcp://localhost:26657", &kb, "val", val.address())
		s.Require().NoError(err)

		auctionQueryClient := auctiontypes.NewQueryClient(clientCtx)
		bankQueryClient := banktypes.NewQueryClient(clientCtx)
		cellarfeesQueryClient := cellarfeestypesv2.NewQueryClient(clientCtx)
		corkQueryClient := corktypes.NewQueryClient(clientCtx)
		distQueryClient := disttypes.NewQueryClient(clientCtx)

		s.T().Log("Verify that the first validator address is an approved cellar ID")
		idsRes, err := corkQueryClient.QueryCellarIDs(context.Background(), &corktypes.QueryCellarIDsRequest{})
		s.Require().NoError(err)

		var found bool
		for _, id := range idsRes.CellarIds {
			if id == ethereumSender {
				found = true
				break
			}
		}
		s.Require().True(found, "validator ethereum address %s is not an approved cellar ID", ethereumSender)

		s.T().Logf("Verify that the module account's fee balances are zero")
		ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
		defer cancel()

		acctsRes, err := cellarfeesQueryClient.QueryModuleAccounts(ctx, &cellarfeestypesv2.QueryModuleAccountsRequest{})
		s.Require().NoError(err, "Failed to query module accounts")

		feesAddress := acctsRes.FeesAddress
		s.T().Logf("Fees address: %s", feesAddress)
		balanceRes, err := bankQueryClient.Balance(ctx, &banktypes.QueryBalanceRequest{
			Address: feesAddress,
			Denom:   fmt.Sprintf("gravity%s", alphaERC20Contract.Hex()),
		})
		s.Require().NoError(err, "Failed to query fee balance of denom %s", alphaERC20Contract.Hex())
		s.Require().Zero(balanceRes.Balance.Amount.Uint64())

		s.T().Log("Waiting for auctions to start")
		alphaAuctionID := uint32(0)
		s.Require().Eventually(func() bool {
			res, _ := auctionQueryClient.QueryActiveAuctions(ctx, &auctiontypes.QueryActiveAuctionsRequest{})

			if res != nil {
				for _, auction := range res.Auctions {
					if auction.StartingTokensForSale.Denom == alphaFeeDenom {
						return true
					}
				}
			}

			return false
		}, time.Second*30, time.Second*5, "Auctions never started for test fees")

		s.T().Log("Bidding to buy all of the ALPHA fees available")
		orch := s.chain.orchestrators[0]
		bidRequest1 := auctiontypes.MsgSubmitBidRequest{
			AuctionId:              alphaAuctionID,
			Signer:                 orch.address().String(),
			MaxBidInUsomm:          sdk.NewCoin(testDenom, sdk.NewIntFromUint64(300000)),
			SaleTokenMinimumAmount: sdk.NewCoin(alphaFeeDenom, sdk.NewIntFromUint64(150000)),
		}

		orchClientCtx, err := s.chain.clientContext("tcp://localhost:26657", orch.keyring, "orch", orch.address())
		s.Require().NoError(err, "Failed to create client for orchestrator")
		_, err = s.chain.sendMsgs(*orchClientCtx, &bidRequest1)
		s.Require().NoError(err, "Failed to submit bid")

		s.T().Log("Bid submitted. Waiting to confirm auction ended")
		s.Require().Eventually(func() bool {
			_, err := auctionQueryClient.QueryEndedAuction(ctx, &auctiontypes.QueryEndedAuctionRequest{
				AuctionId: alphaAuctionID,
			})

			// a nil error indicates the item was found
			return err == nil
		}, time.Second*10, time.Second, "Auction did not end.")

		s.T().Log("Auction ended. Waiting to receive usomm in fees account")
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

		s.T().Log("Waiting to see distribution rate increase")
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

			return res == nil || res.Balance.Amount.Equal(sdk.ZeroInt())
		}, time.Second*300, time.Second*10, "Reward supply did not exhaust in the provided amount of time")

		s.T().Log("Confirming no auction is started...")
		for i := 0; i < 20; i++ {
			res, _ := auctionQueryClient.QueryActiveAuctions(ctx, &auctiontypes.QueryActiveAuctionsRequest{})
			if res == nil {
				continue
			}

			for _, auction := range res.Auctions {
				s.Require().NotEqual(auction.StartingTokensForSale.Denom, alphaFeeDenom)
			}

			time.Sleep(time.Second)
		}

		s.T().Log("Done!")
	})
}
