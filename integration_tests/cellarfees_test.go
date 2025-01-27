package integration_tests

import (
	"context"
	"fmt"
	"math/big"
	"time"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	disttypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	govtypesv1beta1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1beta1"
	paramsproposal "github.com/cosmos/cosmos-sdk/x/params/types/proposal"
	auctiontypes "github.com/peggyjv/sommelier/v9/x/auction/types"
	cellarfeestypes "github.com/peggyjv/sommelier/v9/x/cellarfees/types"
	cellarfeestypesv2 "github.com/peggyjv/sommelier/v9/x/cellarfees/types/v2"
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

		s.Require().NotZero(balanceRes.Balances.AmountOf(gravityDenom))
		s.Require().NotZero(balanceRes.Balances.AmountOf(ibcDenom))

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

func (s *IntegrationTestSuite) TestProceeds() {
	s.Run("Bring up chain, submit TokenPrices, observe proceeds address balance(s) and auction starting", func() {
		val := s.chain.validators[0]
		kb, err := val.keyring()
		s.Require().NoError(err)

		clientCtx, err := s.chain.clientContext("tcp://localhost:26657", &kb, "val", val.address())
		s.Require().NoError(err)

		auctionQueryClient := auctiontypes.NewQueryClient(clientCtx)
		bankQueryClient := banktypes.NewQueryClient(clientCtx)
		cellarfeesQueryClient := cellarfeestypesv2.NewQueryClient(clientCtx)
		proposerCtx, err := s.chain.clientContext("tcp://localhost:26657", s.chain.proposer.keyring, "proposer", s.chain.proposer.address())
		s.Require().NoError(err)

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

		s.Require().NotZero(balanceRes.Balances.AmountOf(gravityDenom))
		s.Require().NotZero(balanceRes.Balances.AmountOf(ibcDenom))

		// Ensure the proceeds account exists by sending it a small amount of usomm
		s.T().Log("Ensuring proceeds account exists")
		proceedsAccount := sdk.MustAccAddressFromBech32(proceedsAddress)
		sendMsg := banktypes.NewMsgSend(
			val.address(),
			proceedsAccount,
			sdk.NewCoins(
				sdk.Coin{
					Denom:  testDenom,
					Amount: math.NewInt(2),
				},
			),
		)
		_, err = s.chain.sendMsgs(*clientCtx, sendMsg)
		s.Require().NoError(err, "Failed to send usomm to proceeds account")

		// Checking that proceeds address has no balance
		s.T().Log("Checking that proceeds address has no balance")
		proceedsBalanceRes, err := bankQueryClient.AllBalances(ctx, &banktypes.QueryAllBalancesRequest{
			Address: proceedsAddress,
		})
		s.Require().NoError(err, "Failed to query proceeds address balance")
		s.Require().True(proceedsBalanceRes.Balances.IsZero(), "Proceeds address should have no balance")

		proceedsPortionStr := "0.33"
		proceedsPortion := sdk.MustNewDecFromStr(proceedsPortionStr)
		s.T().Logf("Proceeds portion: %s", proceedsPortionStr)
		s.T().Logf("Submitting proposal to update proceeds portion")
		proposal := paramsproposal.ParameterChangeProposal{
			Title:       "update proceeds portion",
			Description: "updates proceeds portion",
			Changes: []paramsproposal.ParamChange{
				{
					Subspace: "cellarfees",
					Key:      "ProceedsPortion",
					Value:    fmt.Sprintf("\"%s\"", proceedsPortionStr),
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
			s.chain.proposer.address(),
		)
		s.Require().NoError(err, "Unable to create governance proposal")
		s.T().Log("Submitting proposal")
		submitProposalResponse, err := s.chain.sendMsgs(*proposerCtx, proposalMsg)
		s.Require().NoError(err)
		s.Require().Zero(submitProposalResponse.Code, "raw log: %s", submitProposalResponse.RawLog)

		s.T().Log("Checking that proposal was submitted correctly")
		govQueryClient := govtypesv1beta1.NewQueryClient(proposerCtx)

		s.Require().Eventually(func() bool {
			proposalsQueryResponse, err := govQueryClient.Proposals(context.Background(), &govtypesv1beta1.QueryProposalsRequest{})
			if err != nil {
				s.T().Logf("error querying proposals: %e", err)
				return false
			}

			s.T().Logf("Proposals: %+v", proposalsQueryResponse.Proposals)

			return len(proposalsQueryResponse.Proposals) == 1 && proposalsQueryResponse.Proposals[0].ProposalId == 1 && proposalsQueryResponse.Proposals[0].Status == govtypesv1beta1.StatusVotingPeriod
		}, time.Second*30, time.Second*5, "proposal submission was never found")

		s.T().Log("vote for proposal")
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

		s.T().Log("waiting for proposal to be approved..")
		s.Require().Eventually(func() bool {
			proposalQueryResponse, _ := govQueryClient.Proposal(context.Background(), &govtypesv1beta1.QueryProposalRequest{ProposalId: 1})
			return govtypesv1beta1.StatusPassed == proposalQueryResponse.Proposal.Status
		}, time.Second*30, time.Second*5, "proposal was never accepted")
		s.T().Log("proposal approved!")

		s.T().Log("verifying parameter was changed")
		cellarfeesParamsRes, err := cellarfeesQueryClient.QueryParams(ctx, &cellarfeestypesv2.QueryParamsRequest{})
		s.Require().NoError(err)
		s.Require().Equal(cellarfeesParamsRes.Params.ProceedsPortion, proceedsPortion)

		expectedGravityProceeds := sdk.NewCoin(gravityDenom, balanceRes.Balances.AmountOf(gravityDenom).ToLegacyDec().Mul(proceedsPortion).TruncateInt())
		expectedIbcProceeds := sdk.NewCoin(ibcDenom, balanceRes.Balances.AmountOf(ibcDenom).ToLegacyDec().Add(sdk.OneDec()).Mul(proceedsPortion).TruncateInt())

		s.Require().NotZero(expectedGravityProceeds.Amount)
		s.Require().NotZero(expectedIbcProceeds.Amount)

		// Submit TokenPrices proposal
		orch0 := s.chain.orchestrators[0]
		orch0ClientCtx, err := s.chain.clientContext("tcp://localhost:26657", orch0.keyring, "orch", orch0.address())
		s.Require().NoError(err)
		propID := uint64(2)

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
			s.chain.proposer.address(),
		)
		s.Require().NoError(err, "Unable to create governance proposal")

		s.submitAndVoteForProposal(proposerCtx, orch0ClientCtx, propID, addTokenPricesPropMsg)

		s.T().Log("Waiting for gravity denom auction to start")
		s.Require().Eventually(func() bool {
			res, _ := auctionQueryClient.QueryActiveAuctions(ctx, &auctiontypes.QueryActiveAuctionsRequest{})

			if res != nil {
				for _, auction := range res.Auctions {
					if auction.StartingTokensForSale.Denom == gravityDenom {
						return true
					}
				}
			}

			return false
		}, time.Second*120, time.Second*5, "Auctions never started for gravity fees")

		// Check that proceeds address has the expected balance
		s.Require().Eventually(func() bool {
			proceedsBalanceRes, err := bankQueryClient.AllBalances(ctx, &banktypes.QueryAllBalancesRequest{
				Address: proceedsAddress,
			})
			s.Require().NoError(err, "Failed to query proceeds address balance")

			actualGravityProceeds := proceedsBalanceRes.Balances.AmountOf(gravityDenom)
			s.T().Logf("Actual gravity proceeds: %s, expected gravity proceeds: %s", actualGravityProceeds.String(), expectedGravityProceeds.Amount.String())
			return actualGravityProceeds.Equal(expectedGravityProceeds.Amount)
		}, time.Second*60, time.Second*5, "Proceeds address balance did not match expected")

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
						return true
					}
				}
			}

			return false
		}, time.Second*120, time.Second*5, "Auctions never started for ibc fees")

		s.Require().Eventually(func() bool {
			proceedsBalanceRes, err := bankQueryClient.AllBalances(ctx, &banktypes.QueryAllBalancesRequest{
				Address: proceedsAddress,
			})
			s.Require().NoError(err, "Failed to query proceeds address balance")
			actualIbcProceeds := proceedsBalanceRes.Balances.AmountOf(ibcDenom)
			s.T().Logf("Actual ibc proceeds: %s, expected ibc proceeds: %s", actualIbcProceeds.String(), expectedIbcProceeds.Amount.String())
			return actualIbcProceeds.Equal(expectedIbcProceeds.Amount)
		}, time.Second*60, time.Second*5, "Proceeds address balance did not match expected")

		s.T().Log("Done!")
	})
}
