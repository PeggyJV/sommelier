package integration_tests

import (
	"context"
	"fmt"
	"math/big"
	"strings"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/peggyjv/sommelier/v3/x/cork/types"
)

const CounterABI = `
  [
    {
      "inputs": [],
      "name": "count",
      "outputs": [
        {
          "internalType": "uint256",
          "name": "",
          "type": "uint256"
        }
      ],
      "stateMutability": "view",
      "type": "function"
    },
    {
      "inputs": [],
      "name": "dec",
      "outputs": [],
      "stateMutability": "nonpayable",
      "type": "function"
    },
    {
      "inputs": [],
      "name": "get",
      "outputs": [
        {
          "internalType": "uint256",
          "name": "",
          "type": "uint256"
        }
      ],
      "stateMutability": "view",
      "type": "function"
    },
    {
      "inputs": [],
      "name": "inc",
      "outputs": [],
      "stateMutability": "nonpayable",
      "type": "function"
    }
  ]
`

func ABIEncodedGet() []byte {
	encodedCall, err := abi.JSON(strings.NewReader(CounterABI))
	if err != nil {
		panic(sdkerrors.Wrap(err, "bad ABI definition in code"))
	}

	abiEncodedCall, err := encodedCall.Pack("get")
	if err != nil {
		panic(err)
	}

	return abiEncodedCall
}

func ABIEncodedInc() []byte {
	encodedCall, err := abi.JSON(strings.NewReader(CounterABI))
	if err != nil {
		panic(sdkerrors.Wrap(err, "bad ABI definition in code"))
	}

	abiEncodedCall, err := encodedCall.Pack("inc")
	if err != nil {
		panic(err)
	}

	return abiEncodedCall
}

func UnpackEthUInt(bz []byte) sdk.Int {
	output := big.NewInt(0)
	output.SetBytes(bz)

	return sdk.NewIntFromBigInt(output)
}

func (s *IntegrationTestSuite) getCurrentCount() (*sdk.Int, error) {
	ethClient, err := ethclient.Dial(fmt.Sprintf("http://%s", s.ethResource.GetHostPort("8545/tcp")))
	if err != nil {
		return nil, err
	}

	bz, err := ethClient.CallContract(context.Background(), ethereum.CallMsg{
		From: common.HexToAddress(s.chain.validators[0].ethereumKey.address),
		To:   &counterContract,
		Gas:  0,
		Data: ABIEncodedGet(),
	}, nil)
	if err != nil {
		return nil, err
	}

	count := UnpackEthUInt(bz)

	return &count, nil
}

func (s *IntegrationTestSuite) TestCork() {
	s.Run("Bring up chain, and submit a cork call to ethereum", func() {

		// makes sure ethereum can be contacted and counter contract is working
		count, err := s.getCurrentCount()
		s.Require().NoError(err)
		s.Require().Equal(int64(0), count.Int64())

		s.T().Logf("create governance proposal to add counter contract")
		orch := s.chain.orchestrators[0]
		clientCtx, err := s.chain.clientContext("tcp://localhost:26657", orch.keyring, "orch", orch.keyInfo.GetAddress())
		s.Require().NoError(err)

		proposal := types.AddManagedCellarsProposal{
			Title:       "add counter contract in test",
			Description: "test description",
			CellarIds: &types.CellarIDSet{
				Ids: []string{counterContract.Hex()},
			},
		}
		proposalMsg, err := govtypes.NewMsgSubmitProposal(
			&proposal,
			sdk.Coins{
				{
					Denom:  testDenom,
					Amount: stakeAmount.Quo(sdk.NewInt(1000)),
				},
			},
			orch.keyInfo.GetAddress(),
		)
		s.Require().NoError(err, "unable to create governance proposal")

		response, err := s.chain.sendMsgs(*clientCtx, proposalMsg)
		s.Require().NoError(err)
		s.Require().Zero(response.Code, "raw log: %s", response.RawLog)

		s.T().Logf("verify that contract exists in allowed addresses")
		val := s.chain.validators[0]
		s.Require().Eventuallyf(func() bool {
			kb, err := val.keyring()
			s.Require().NoError(err)
			clientCtx, err := s.chain.clientContext("tcp://localhost:26657", &kb, "val", val.keyInfo.GetAddress())
			s.Require().NoError(err)

			queryClient := types.NewQueryClient(clientCtx)
			res, err := queryClient.QueryCellarIDs(context.Background(), &types.QueryCellarIDsRequest{})
			if err != nil {
				return false
			}

			found := false
			for _, id := range res.CellarIds {
				s.T().Logf("managed addresses: %v", res.CellarIds)
				if common.HexToAddress(id) == counterContract {
					found = true
					break
				}
			}

			return found
		}, 10*time.Second, 2*time.Second, "did not find address in managed cellars")

		s.T().Logf("wait for new vote period start")
		val = s.chain.validators[0]
		s.Require().Eventuallyf(func() bool {
			kb, err := val.keyring()
			s.Require().NoError(err)
			clientCtx, err := s.chain.clientContext("tcp://localhost:26657", &kb, "val", val.keyInfo.GetAddress())
			s.Require().NoError(err)

			queryClient := types.NewQueryClient(clientCtx)
			res, err := queryClient.QueryCommitPeriod(context.Background(), &types.QueryCommitPeriodRequest{})
			if err != nil {
				return false
			}
			if res.VotePeriodStart != res.CurrentHeight {
				if res.CurrentHeight%10 == 0 {
					s.T().Logf("current height: %d, period end: %d", res.CurrentHeight, res.VotePeriodEnd)
				}
				return false
			}
			return true
		}, 105*time.Second, 1*time.Second, "new vote period never seen")

		s.T().Logf("sending cork calls")
		for i, orch := range s.chain.orchestrators {
			i := i
			orch := orch
			s.Require().Eventuallyf(func() bool {
				clientCtx, err := s.chain.clientContext("tcp://localhost:26657", orch.keyring, "orch", orch.keyInfo.GetAddress())
				s.Require().NoError(err)

				corkMsg, err := types.NewMsgSubmitCorkRequest(ABIEncodedInc(), counterContract, orch.keyInfo.GetAddress())
				s.Require().NoError(err, "unable to create cork msg")

				response, err := s.chain.sendMsgs(*clientCtx, corkMsg)
				if err != nil {
					s.T().Logf("error: %s", err)
					return false
				}
				if response.Code != 0 {
					if response.Code != 32 {
						s.T().Log(response)
					}
					return false
				}

				s.T().Logf("cork msg for %d node sent successfully", i)
				return true
			}, 10*time.Second, 500*time.Millisecond, "unable to deploy cork msg for node %d", i)
		}

		s.T().Logf("waiting for end of vote period, endblocker to run")
		val = s.chain.validators[0]
		s.Require().Eventuallyf(func() bool {
			kb, err := val.keyring()
			s.Require().NoError(err)
			clientCtx, err := s.chain.clientContext("tcp://localhost:26657", &kb, "val", val.keyInfo.GetAddress())
			s.Require().NoError(err)

			queryClient := types.NewQueryClient(clientCtx)
			res, err := queryClient.QueryCommitPeriod(context.Background(), &types.QueryCommitPeriodRequest{})
			if err != nil {
				return false
			}
			if res.VotePeriodStart != res.CurrentHeight {
				if res.CurrentHeight%10 == 0 {
					s.T().Logf("current height: %d, period end: %d", res.CurrentHeight, res.VotePeriodEnd)
				}
				return false
			}

			return true
		}, 105*time.Second, 1*time.Second, "new vote period never seen")

		s.T().Logf("checking for updated count in contract")
		s.Require().Eventuallyf(func() bool {
			count, err = s.getCurrentCount()
			if err != nil {
				s.T().Logf("got error %e querying count", err)
				return false
			}
			if count.Int64() != int64(1) {
				s.T().Logf("wrong count %s", count.String())
				return false
			}

			return true
		}, 1*time.Minute, 5*time.Second, "count never updated")

	})
}
