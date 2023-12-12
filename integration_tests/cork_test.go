package integration_tests

import (
	"context"
	"fmt"
	"strings"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
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
		panic(errorsmod.Wrap(err, "bad ABI definition in code"))
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
		panic(errorsmod.Wrap(err, "bad ABI definition in code"))
	}

	abiEncodedCall, err := encodedCall.Pack("inc")
	if err != nil {
		panic(err)
	}

	return abiEncodedCall
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
