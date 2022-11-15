package integration_tests

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"log"
	"math/big"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	disttypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	ethereumtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	auctiontypes "github.com/peggyjv/sommelier/v4/x/auction/types"
	cellarfeestypes "github.com/peggyjv/sommelier/v4/x/cellarfees/types"
	corktypes "github.com/peggyjv/sommelier/v4/x/cork/types"
)

const ALPHA_FEE_DENOM string = "gravity0x4C4a2f8c81640e47606d3fd77B353E87Ba015584"
const BETA_FEE_DENOM string = "gravity0x21dF544947ba3E8b3c32561399E88B52Dc8b2823"

func (s *IntegrationTestSuite) sendEthTransaction(ethClient *ethclient.Client, ethereumKey *ethereumKey, toAddress common.Address, data []byte) error {
	privateKey, err := crypto.HexToECDSA(ethereumKey.privateKey[2:])
	if err != nil {
		return err
	}

	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Fatal("error casting public key to ECDSA")
	}

	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)
	nonce, err := ethClient.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		return err
	}

	value := big.NewInt(0)
	gasLimit := uint64(1000000)
	gasPrice, err := ethClient.SuggestGasPrice(context.Background())
	if err != nil {
		return err
	}

	tx := ethereumtypes.NewTransaction(nonce, toAddress, value, gasLimit, gasPrice, data)

	chainID, err := ethClient.NetworkID(context.Background())
	if err != nil {
		return err
	}

	signedTx, err := ethereumtypes.SignTx(tx, ethereumtypes.NewEIP155Signer(chainID), privateKey)
	if err != nil {
		return err
	}

	err = ethClient.SendTransaction(context.Background(), signedTx)
	if err != nil {
		return err
	}

	return nil
}

func (s *IntegrationTestSuite) TestCellarFees() {
	s.Run("Bring up chain, send fees from ethereum, observe auction and fee distribution", func() {
		val := s.chain.validators[0]
		ethereumSender := val.ethereumKey.address
		kb, err := val.keyring()
		s.Require().NoError(err)

		ethClient, err := ethclient.Dial(fmt.Sprintf("http://%s", s.ethResource.GetHostPort("8545/tcp")))
		s.Require().NoError(err)

		clientCtx, err := s.chain.clientContext("tcp://localhost:26657", &kb, "val", val.keyInfo.GetAddress())
		s.Require().NoError(err)

		auctionQueryClient := auctiontypes.NewQueryClient(clientCtx)
		bankQueryClient := banktypes.NewQueryClient(clientCtx)
		cellarfeesQueryClient := cellarfeestypes.NewQueryClient(clientCtx)
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

		acctsRes, err := cellarfeesQueryClient.ModuleAccounts(ctx, &cellarfeestypes.QueryModuleAccountsRequest{})
		s.Require().NoError(err, "Failed to query module accounts")

		feesAddress := acctsRes.FeesAddress
		s.T().Logf("Fees address: %s", feesAddress)
		balanceRes, err := bankQueryClient.Balance(ctx, &banktypes.QueryBalanceRequest{
			Address: feesAddress,
			Denom:   fmt.Sprintf("gravity%s", alphaERC20Contract.Hex()),
		})
		s.Require().NoError(err, "Failed to query fee balance of denom %s", alphaERC20Contract.Hex())
		s.Require().Zero(balanceRes.Balance.Amount.Uint64())
		balanceRes, err = bankQueryClient.Balance(ctx, &banktypes.QueryBalanceRequest{
			Address: feesAddress,
			Denom:   fmt.Sprintf("gravity%s", betaERC20Contract.Hex()),
		})
		s.Require().NoError(err, "Failed to query fee balance of denom %s", betaERC20Contract.Hex())
		s.Require().Zero(balanceRes.Balance.Amount.Uint64())

		s.T().Log("Getting baseline rewards increase rate")
		rewardsRes, err := distQueryClient.DelegationRewards(ctx, &disttypes.QueryDelegationRewardsRequest{
			DelegatorAddress: val.keyInfo.GetAddress().String(),
			ValidatorAddress: "sommvaloper199sjfhaw3hempwzljw0lgwsm9kk6r8e5ef3hmp",
		})
		s.Require().NoError(err)
		startAmount := rewardsRes.Rewards.AmountOf(testDenom)
		time.Sleep(time.Second * 12)
		rewardsRes, err = distQueryClient.DelegationRewards(ctx, &disttypes.QueryDelegationRewardsRequest{
			DelegatorAddress: val.keyInfo.GetAddress().String(),
			ValidatorAddress: "sommvaloper199sjfhaw3hempwzljw0lgwsm9kk6r8e5ef3hmp",
		})
		s.Require().NoError(err)
		endAmount := rewardsRes.Rewards.AmountOf(testDenom)
		rewardRateBaseline := endAmount.Sub(startAmount).Quo(sdk.NewDec(12))

		s.T().Logf("Approving Gravity to spend Alpha ERC20")
		approveData := PackApproveERC20(gravityContract)
		err = s.sendEthTransaction(ethClient, &val.ethereumKey, alphaERC20Contract, approveData)
		s.Require().NoError(err, "Error approving spending ALPHA balance for the gravity contract on behalf of the first validator")

		s.T().Logf("Approving Gravity to spend Beta ERC20")
		approveData = PackApproveERC20(gravityContract)
		err = s.sendEthTransaction(ethClient, &val.ethereumKey, betaERC20Contract, approveData)
		s.Require().NoError(err, "Error approving spending BETA balance for the gravity contract on behalf of the first validator")

		s.T().Logf("Waiting for allowance confirmations..")
		s.Require().Eventually(func() bool {
			data := PackAllowance(common.HexToAddress(ethereumSender), gravityContract)

			res, _ := ethClient.CallContract(context.Background(), ethereum.CallMsg{
				From: common.HexToAddress(ethereumSender),
				To:   &alphaERC20Contract,
				Gas:  0,
				Data: data,
			}, nil)

			allowance := UnpackEthUInt(res).BigInt()
			s.T().Logf("Allowance: %v", allowance)

			return sdk.NewIntFromBigInt(allowance).GT(sdk.ZeroInt())
		}, time.Second*10, time.Second, "AlphaERC20 allowance not found")

		s.Require().Eventually(func() bool {
			data := PackAllowance(common.HexToAddress(ethereumSender), gravityContract)

			res, _ := ethClient.CallContract(context.Background(), ethereum.CallMsg{
				From: common.HexToAddress(ethereumSender),
				To:   &betaERC20Contract,
				Gas:  0,
				Data: data,
			}, nil)

			allowance := UnpackEthUInt(res).BigInt()
			s.T().Logf("Allowance: %v", allowance)

			return sdk.NewIntFromBigInt(allowance).GT(sdk.ZeroInt())
		}, time.Second*10, time.Second, "BetaERC20 allowance not found")

		s.T().Log("Sending ALPHA fees to cellarfees module account")
		acc, err := sdk.AccAddressFromBech32(feesAddress)
		s.Require().NoError(err, "Failed to derive fees account address from bech32 string: %s", feesAddress)
		sendData := PackSendToCosmos(alphaERC20Contract, acc, sdk.NewInt(50000))
		err = s.sendEthTransaction(ethClient, &val.ethereumKey, gravityContract, sendData)
		s.Require().NoError(err, "Failed to send fees transaction to Cosmos")

		s.T().Log("Sending BETA fees to cellarfees module account")
		acc, err = sdk.AccAddressFromBech32(feesAddress)
		s.Require().NoError(err, "Failed to derive fees account address from bech32 string: %s", feesAddress)
		sendData = PackSendToCosmos(betaERC20Contract, acc, sdk.NewInt(20000))
		err = s.sendEthTransaction(ethClient, &val.ethereumKey, gravityContract, sendData)
		s.Require().NoError(err, "Failed to send fees transaction to Cosmos")

		s.T().Log("Waiting for fees to be received...")
		s.Require().Eventually(func() bool {
			res, err := bankQueryClient.Balance(context.Background(),
				&banktypes.QueryBalanceRequest{
					Address: feesAddress,
					Denom:   ALPHA_FEE_DENOM,
				})
			s.Require().NoError(err)
			s.T().Logf("fee balance: %s", res.Balance)

			return res.Balance.Amount.GT(sdk.ZeroInt())
		}, time.Second*60, time.Second*6, "ALPHA Fees never received by cellarfees account")

		s.Require().Eventually(func() bool {
			res, err := bankQueryClient.Balance(context.Background(),
				&banktypes.QueryBalanceRequest{
					Address: feesAddress,
					Denom:   BETA_FEE_DENOM,
				})
			s.Require().NoError(err)
			s.T().Logf("fee balance: %s", res.Balance)

			return res.Balance.Amount.GT(sdk.ZeroInt())
		}, time.Second*60, time.Second*6, "BETA Fees never received by cellarfees account")

		s.T().Log("Fees received! Confirming no auction gets started yet...")
		for i := 0; i < 10; i++ {
			res, _ := auctionQueryClient.QueryActiveAuctions(ctx, &auctiontypes.QueryActiveAuctionsRequest{})
			if res == nil {
				continue
			}

			for _, auction := range res.Auctions {
				s.Require().NotEqual(auction.StartingTokensForSale.Denom, ALPHA_FEE_DENOM)
				s.Require().NotEqual(auction.StartingTokensForSale.Denom, BETA_FEE_DENOM)
			}

			time.Sleep(time.Second)
		}

		s.T().Log("Sending ERC20 fees a second time")
		sendData = PackSendToCosmos(alphaERC20Contract, acc, sdk.NewInt(100000))
		err = s.sendEthTransaction(ethClient, &val.ethereumKey, gravityContract, sendData)
		s.Require().NoError(err, "Failed to send fees transaction to Cosmos")

		sendData = PackSendToCosmos(betaERC20Contract, acc, sdk.NewInt(120000))
		err = s.sendEthTransaction(ethClient, &val.ethereumKey, gravityContract, sendData)
		s.Require().NoError(err, "Failed to send fees transaction to Cosmos")

		s.T().Log("Waiting for auctions to start")
		alphaAuctionId, betaAuctionId := uint32(0), uint32(0)
		s.Require().Eventually(func() bool {
			res, _ := auctionQueryClient.QueryActiveAuctions(ctx, &auctiontypes.QueryActiveAuctionsRequest{})

			alpha, beta := false, false
			if res != nil {
				for _, auction := range res.Auctions {
					if auction.StartingTokensForSale.Denom == ALPHA_FEE_DENOM {
						alphaAuctionId = auction.Id
						alpha = true
					} else if auction.StartingTokensForSale.Denom == BETA_FEE_DENOM {
						betaAuctionId = auction.Id
						beta = true
					}

					if alpha && beta {
						break
					}
				}
			}

			return alpha && beta
		}, time.Second*30, time.Second*5, "Auctions never started for test fees")

		s.T().Log("Bidding to buy all of the ALPHA fees available")
		orch := s.chain.orchestrators[0]
		bidRequest1 := auctiontypes.MsgSubmitBidRequest{
			AuctionId:              alphaAuctionId,
			Signer:                 orch.keyInfo.GetAddress().String(),
			MaxBidInUsomm:          sdk.NewCoin("usomm", sdk.NewIntFromUint64(300000)),
			SaleTokenMinimumAmount: sdk.NewCoin(ALPHA_FEE_DENOM, sdk.NewIntFromUint64(150000)),
		}

		orchClientCtx, err := s.chain.clientContext("tcp://localhost:26657", orch.keyring, "orch", orch.keyInfo.GetAddress())
		s.Require().NoError(err, "Failed to create client for orchestrator")
		_, err = s.chain.sendMsgs(*orchClientCtx, &bidRequest1)
		s.Require().NoError(err, "Failed to submit bid")

		s.T().Log("Bid submitted. Waiting to receive usomm in fees account")
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

		s.T().Log("Distribution rate is linear. Increasing the reward supply by bidding on the BETA auction")
		bidRequest2 := auctiontypes.MsgSubmitBidRequest{
			AuctionId:              betaAuctionId,
			Signer:                 orch.keyInfo.GetAddress().String(),
			MaxBidInUsomm:          sdk.NewCoin("usomm", sdk.NewIntFromUint64(1400000)),
			SaleTokenMinimumAmount: sdk.NewCoin(BETA_FEE_DENOM, sdk.NewIntFromUint64(140000)),
		}

		_, err = s.chain.sendMsgs(*orchClientCtx, &bidRequest2)
		s.Require().NoError(err, "Failed to submit bid")

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

		s.T().Log("Distribution rate increased with supply! Getting current reward rate")
		rewardsRes, err = distQueryClient.DelegationRewards(ctx, &disttypes.QueryDelegationRewardsRequest{
			DelegatorAddress: val.keyInfo.GetAddress().String(),
			ValidatorAddress: "sommvaloper199sjfhaw3hempwzljw0lgwsm9kk6r8e5ef3hmp",
		})
		s.Require().NoError(err)

		startAmount = rewardsRes.Rewards.AmountOf(testDenom)

		// let some time elapse so we can calculate an average rate
		time.Sleep(time.Second * 12)

		rewardsRes, err = distQueryClient.DelegationRewards(ctx, &disttypes.QueryDelegationRewardsRequest{
			DelegatorAddress: val.keyInfo.GetAddress().String(),
			ValidatorAddress: "sommvaloper199sjfhaw3hempwzljw0lgwsm9kk6r8e5ef3hmp",
		})
		s.Require().NoError(err)
		endAmount = rewardsRes.Rewards.AmountOf(testDenom)
		rewardRate := (endAmount.Sub(startAmount).Quo(sdk.NewDec(12)))
		s.T().Logf("Baseline reward rate: %d, current reward rate: %d", rewardRateBaseline.RoundInt64(), rewardRate.RoundInt64())
		s.Require().True(rewardRate.GT(rewardRateBaseline), "Reward rate has not increased")

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

		s.T().Log("Verify that the accrual counter reset by sending more ALPHA")
		sendData = PackSendToCosmos(alphaERC20Contract, acc, sdk.NewInt(25000))
		err = s.sendEthTransaction(ethClient, &val.ethereumKey, gravityContract, sendData)
		s.Require().NoError(err, "Failed to send fees transaction to Cosmos")

		s.T().Log("Confirming no auction is started...")
		for i := 0; i < 20; i++ {
			res, _ := auctionQueryClient.QueryActiveAuctions(ctx, &auctiontypes.QueryActiveAuctionsRequest{})
			if res == nil {
				continue
			}

			for _, auction := range res.Auctions {
				s.Require().NotEqual(auction.StartingTokensForSale.Denom, ALPHA_FEE_DENOM)
			}

			time.Sleep(time.Second)
		}

		s.T().Log("Done!")
	})
}
