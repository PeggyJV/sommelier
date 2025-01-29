package integration_tests

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"log"
	"math/big"
	"strings"

	errorsmod "cosmossdk.io/errors"
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	ethereumtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	gravitytypes "github.com/peggyjv/gravity-bridge/module/v6/x/gravity/types"
)

type EthereumConfig struct {
	ChainID             uint `json:"chainId"`
	HomesteadBlock      uint `json:"homesteadBlock"`
	EIP150Block         uint `json:"eip150Block"`
	EIP155Block         uint `json:"eip155Block"`
	EIP158Block         uint `json:"eip158Block"`
	ByzantiumBlock      uint `json:"byzantiumBlock"`
	ConstantinopleBlock uint `json:"constantinopleBlock"`
	PetersburgBlock     uint `json:"petersburgBlock"`
	IstanbulBlock       uint `json:"istanbulBlock"`
	BerlinBlock         uint `json:"berlinBlock"`
}

type Allocation struct {
	Balance string `json:"balance"`
}

type EthereumGenesis struct {
	Difficulty string                `json:"difficulty"`
	GasLimit   string                `json:"gasLimit"`
	Config     EthereumConfig        `json:"config"`
	Alloc      map[string]Allocation `json:"alloc"`
}

const approveERC20ABIJSON = `
[
	{
      "inputs": [
        {
          "internalType": "address",
          "name": "spender",
          "type": "address"
        },
        {
          "internalType": "uint256",
          "name": "amount",
          "type": "uint256"
        }
      ],
      "name": "approve",
      "outputs": [
        {
          "internalType": "bool",
          "name": "",
          "type": "bool"
        }
      ],
      "stateMutability": "nonpayable",
      "type": "function"
    }
]
`

const balanceOfERC20ABIJSON = `
[
	{
      "inputs": [
        {
          "internalType": "address",
          "name": "account",
          "type": "address"
        }
      ],
      "name": "balanceOf",
      "outputs": [
        {
          "internalType": "uint256",
          "name": "",
          "type": "uint256"
        }
      ],
      "stateMutability": "view",
      "type": "function"
    }
]
`

const allowanceERC20ABIJSON = `
[
	{
      "inputs": [
        {
          "internalType": "address",
          "name": "owner",
          "type": "address"
        },
        {
          "internalType": "address",
          "name": "spender",
          "type": "address"
        }
      ],
      "name": "allowance",
      "outputs": [
        {
          "internalType": "uint256",
          "name": "",
          "type": "uint256"
        }
      ],
      "stateMutability": "view",
      "type": "function"
    }
]
`

const stateLastValsetNonceABIJSON = `[
	{
		"inputs": [],
		"name": "state_lastValsetNonce",
		"outputs": [
		  {
			"internalType": "uint256",
			"name": "",
			"type": "uint256"
		  }
		],
		"stateMutability": "view",
		"type": "function"
	  }
]`

func packCall(abiString, method string, args []interface{}) []byte {
	encodedCall, err := abi.JSON(strings.NewReader(abiString))
	if err != nil {
		panic(errorsmod.Wrap(err, "bad ABI definition in code"))
	}
	abiEncodedCall, err := encodedCall.Pack(method, args...)
	if err != nil {
		panic(errorsmod.Wrap(err, "error packing calling"))
	}
	return abiEncodedCall
}

func PackDeployERC20(denom string, name string, symbol string, decimals uint8) []byte {
	return packCall(gravitytypes.DeployERC20ABIJSON, "deployERC20", []interface{}{
		denom,
		name,
		symbol,
		decimals,
	})
}

func PackSendToCosmos(tokenContract common.Address, destination sdk.AccAddress, amount math.Int) []byte {
	destinationBytes, _ := byteArrayToFixByteArray(destination.Bytes())
	return packCall(gravitytypes.SendToCosmosABIJSON, "sendToCosmos", []interface{}{
		tokenContract,
		destinationBytes,
		amount.BigInt(),
	})
}

func UInt256Max() *big.Int {
	return new(big.Int).Sub(new(big.Int).Lsh(common.Big1, 256), common.Big1)
}

func PackApproveERC20(spender common.Address) []byte {
	return packCall(approveERC20ABIJSON, "approve", []interface{}{
		spender,
		UInt256Max(),
	})
}

func PackBalanceOf(account common.Address) []byte {
	return packCall(balanceOfERC20ABIJSON, "balanceOf", []interface{}{
		account,
	})
}

func PackAllowance(owner common.Address, spender common.Address) []byte {
	return packCall(allowanceERC20ABIJSON, "allowance", []interface{}{
		owner,
		spender,
	})
}

func PackLastValsetNonce() []byte {
	return packCall(stateLastValsetNonceABIJSON, "state_lastValsetNonce", []interface{}{})
}

func UnpackEthUInt(bz []byte) math.Int {
	output := big.NewInt(0)
	output.SetBytes(bz)

	return sdk.NewIntFromBigInt(output)
}

func byteArrayToFixByteArray(b []byte) (out [32]byte, err error) {
	if len(b) > 32 {
		return out, fmt.Errorf("array too long")
	}

	copy(out[12:], b)
	return out, nil
}

func SendEthTransaction(ethClient *ethclient.Client, ethereumKey *ethereumKey, toAddress common.Address, data []byte) error {
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
