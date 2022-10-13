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
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	ethereumtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	auctiontypes "github.com/peggyjv/sommelier/v4/x/auction/types"
	cellarfeestypes "github.com/peggyjv/sommelier/v4/x/cellarfees/types"
	corktypes "github.com/peggyjv/sommelier/v4/x/cork/types"
)

const CELLAR_FEE_DENOM string = "gravity0x4C4a2f8c81640e47606d3fd77B353E87Ba015584"

func (s *IntegrationTestSuite) SendEthTransaction(ethereumKey *ethereumKey, toAddress common.Address, data []byte) error {
	ethClient, err := ethclient.Dial(fmt.Sprintf("http://%s", s.ethResource.GetHostPort("8545/tcp")))
	if err != nil {
		return err
	}

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
		s.T().Log("Verify that the first validator address is an approved cellar ID")

		val := s.chain.validators[0]
		ethereumSender := val.ethereumKey.address
		kb, err := val.keyring()
		s.Require().NoError(err)

		clientCtx, err := s.chain.clientContext("tcp://localhost:26657", &kb, "val", val.keyInfo.GetAddress())
		s.Require().NoError(err)

		queryClient := corktypes.NewQueryClient(clientCtx)
		idsRes, err := queryClient.QueryCellarIDs(context.Background(), &corktypes.QueryCellarIDsRequest{})
		s.Require().NoError(err)

		var found bool
		for _, id := range idsRes.CellarIds {
			if id == ethereumSender {
				found = true
				break
			}
		}
		s.Require().True(found, "validator ethereum address %s is not an approved cellar ID", ethereumSender)

		s.T().Logf("Verify that the module account's balance of fee denom %s is zero", testERC20Contract.Hex())

		cellarfeesQueryClient := cellarfeestypes.NewQueryClient(clientCtx)
		ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
		defer cancel()

		acctsRes, err := cellarfeesQueryClient.ModuleAccounts(ctx, &cellarfeestypes.QueryModuleAccountsRequest{})
		s.Require().NoError(err, "Failed to query module accounts")

		feesAddress := acctsRes.FeesAddress
		s.T().Logf("Fees address: %s", feesAddress)
		bankQueryClient := banktypes.NewQueryClient(clientCtx)
		balanceRes, err := bankQueryClient.Balance(ctx, &banktypes.QueryBalanceRequest{
			Address: feesAddress,
			Denom:   fmt.Sprintf("gravity%s", testERC20Contract.Hex()),
		})
		s.Require().NoError(err, "Failed to query fee balance of denom %s", testERC20Contract.Hex())
		s.Require().Zero(balanceRes.Balance.Amount.Uint64())

		s.T().Logf("Approving Gravity to spend ERC 20")
		approveData := PackApproveERC20(gravityContract)
		err = s.SendEthTransaction(&val.ethereumKey, testERC20Contract, approveData)
		s.Require().NoError(err, "Error approving spending balance for the gravity contract on behalf of the first validator")

		s.T().Logf("Waiting for allowance confirmation...")
		s.Require().Eventually(func() bool {
			ethClient, err := ethclient.Dial(fmt.Sprintf("http://%s", s.ethResource.GetHostPort("8545/tcp")))
			if err != nil {
				return false
			}

			data := PackAllowance(common.HexToAddress(ethereumSender), gravityContract)

			response, _ := ethClient.CallContract(context.Background(), ethereum.CallMsg{
				From: common.HexToAddress(ethereumSender),
				To:   &testERC20Contract,
				Gas:  0,
				Data: data,
			}, nil)

			allowance := UnpackEthUInt(response).BigInt()
			s.T().Logf("Allowance: %v", allowance)

			return sdk.NewIntFromBigInt(allowance).GT(sdk.ZeroInt())
		}, time.Second*10, time.Second, "TestERC20 allowance not found")

		s.T().Log("Sending ERC20 fees to cellarfees module account")
		acc, err := sdk.AccAddressFromBech32(feesAddress)
		s.Require().NoError(err, "Failed to derive fees account address from bech32 string: %s", feesAddress)
		sendData := PackSendToCosmos(testERC20Contract, acc, sdk.NewInt(50000))
		err = s.SendEthTransaction(&val.ethereumKey, gravityContract, sendData)
		s.Require().NoError(err, "Failed to send fees transaction to Cosmos")

		s.T().Log("Waiting for fees to be received...")
		// var balance sdk.Int
		s.Require().Eventually(func() bool {
			kb, err := val.keyring()
			s.Require().NoError(err)
			clientCtx, err := s.chain.clientContext("tcp://localhost:26657", &kb, "cellarfees", sdk.AccAddress(feesAddress))
			s.Require().NoError(err)

			queryClient := banktypes.NewQueryClient(clientCtx)
			res, err := queryClient.Balance(context.Background(),
				&banktypes.QueryBalanceRequest{
					Address: feesAddress,
					Denom:   CELLAR_FEE_DENOM,
				})
			s.Require().NoError(err)
			s.T().Logf("fee balance: %s", res.Balance)

			if res.Balance.Amount.GT(sdk.ZeroInt()) {
				// balance = res.Balance.Amount
				return true
			}

			return false
		}, time.Second*60, time.Second*6, "Fees never received by cellarfees account")

		s.T().Log("Fees received! Confirming no auction gets started yet...")
		for i := 0; i < 10; i++ {
			kb, err := val.keyring()
			s.Require().NoError(err)
			clientCtx, err := s.chain.clientContext("tcp://localhost:26657", &kb, "cellarfees", sdk.AccAddress(feesAddress))
			s.Require().NoError(err)

			queryClient := auctiontypes.NewQueryClient(clientCtx)
			res, _ := queryClient.QueryActiveAuctions(ctx, &auctiontypes.QueryActiveAuctionsRequest{})
			s.Require().Nil(res)

			time.Sleep(time.Second)
		}

		s.T().Log("Sending ERC20 fees a second time")
		sendData = PackSendToCosmos(testERC20Contract, acc, sdk.NewInt(100000))
		err = s.SendEthTransaction(&val.ethereumKey, gravityContract, sendData)
		s.Require().NoError(err, "Failed to send fees transaction to Cosmos")

		s.T().Log("Waiting for auction to start")
		var auctionId uint32
		s.Require().Eventually(func() bool {
			kb, err := val.keyring()
			s.Require().NoError(err)
			clientCtx, err := s.chain.clientContext("tcp://localhost:26657", &kb, "cellarfees", sdk.AccAddress(feesAddress))
			s.Require().NoError(err)

			queryClient := auctiontypes.NewQueryClient(clientCtx)
			res, _ := queryClient.QueryActiveAuctions(ctx, &auctiontypes.QueryActiveAuctionsRequest{})

			if res != nil {
				auctionId = res.Auctions[0].Id
				return true
			}

			return false
		}, time.Second*30, time.Second*5, "Auction never started for test fees")

		s.T().Log("Bidding to buy half of the fees available")
		orch := s.chain.orchestrators[0]
		bidRequest1 := auctiontypes.MsgSubmitBidRequest{
			AuctionId:              auctionId,
			Bidder:                 orch.keyInfo.GetAddress().String(),
			MaxBidInUsomm:          sdk.NewCoin("usomm", sdk.NewIntFromUint64(150000)),
			SaleTokenMinimumAmount: sdk.NewCoin(CELLAR_FEE_DENOM, sdk.NewIntFromUint64(75000)),
			Signer:                 orch.keyInfo.GetAddress().String(),
		}
		s.T().Logf("Bid: %v", bidRequest1)
		orchClientCtx, err := s.chain.clientContext("tcp://localhost:26657", orch.keyring, "orch", orch.keyInfo.GetAddress())
		s.Require().NoError(err, "Failed to create client for orchestrator")
		submitBid, err := s.chain.sendMsgs(*orchClientCtx, &bidRequest1)
		s.Require().NoError(err, "Failed to submit bid")
		s.T().Log(submitBid)

		s.T().Log("Bid submitted. Waiting to see SOMM in cellarfees module account...")
		s.Require().Eventually(func() bool {
			kb, err := val.keyring()
			s.Require().NoError(err)
			clientCtx, err := s.chain.clientContext("tcp://localhost:26657", &kb, "cellarfees", sdk.AccAddress(feesAddress))
			s.Require().NoError(err)

			queryClient := banktypes.NewQueryClient(clientCtx)
			res, err := queryClient.Balance(ctx, &banktypes.QueryBalanceRequest{
				Address: feesAddress,
				Denom:   testDenom,
			})
			s.Require().NoError(err, "Error querying for SOMM balance in fees account")
			s.Require().NotNil(res)

			s.T().Logf("SOMM balance: %v", res.Balance)
			return res.Balance.Amount.GT(sdk.ZeroInt())
		}, time.Second*60, time.Second*5, "Never received SOMM from auction")

		s.T().Log("Done!")
	})
}
