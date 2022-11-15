package integration_tests

import (
	"fmt"
	"math/big"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	gravitytypes "github.com/peggyjv/gravity-bridge/module/v2/x/gravity/types"
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
		panic(sdkerrors.Wrap(err, "bad ABI definition in code"))
	}
	abiEncodedCall, err := encodedCall.Pack(method, args...)
	if err != nil {
		panic(sdkerrors.Wrap(err, "error packing calling"))
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

func PackSendToCosmos(tokenContract common.Address, destination sdk.AccAddress, amount sdk.Int) []byte {
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

func UnpackEthUInt(bz []byte) sdk.Int {
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
