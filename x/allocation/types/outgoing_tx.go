package types

import (
	"math/big"
	"strings"

	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/crypto"
)

type ABIEncodedTickRange struct {
	TokenID *big.Int `abi:"tokenId"`
	Upper   *big.Int `abi:"tickUpper"`
	Lower   *big.Int `abi:"tickLower"`
	Weight  *big.Int `abi:"weight"`
}

// ABIEncodedRebalanceHash gets the checkpoint signature from the given outgoing tx batch
func (c Cellar) ABIEncodedRebalanceHash() []byte {
	encodedCall, err := abi.JSON(strings.NewReader(rebalanceABI))
	if err != nil {
		panic(sdkerrors.Wrap(err, "bad ABI definition in code"))
	}



	ticks := make([]ABIEncodedTickRange, len(c.TickRanges))
	for _, t := range c.TickRanges {
		up := int64(t.Upper)
		lo := int64(t.Lower)
		we := uint64(t.Weight)
		ticks = append(ticks, ABIEncodedTickRange{big.NewInt(0), big.NewInt(up), big.NewInt(lo), new(big.Int).SetUint64(we)})
	}

	abiEncodedCall, err := encodedCall.Pack("rebalance", ticks)
	if err != nil {
		panic(err)
	}

	return crypto.Keccak256Hash(abiEncodedCall).Bytes()
}

const rebalanceABI = `[{
	"inputs": [
		{
			"components": [
				{ "internalType": "uint184", "name": "tokenId",   "type": "uint184" },
				{ "internalType": "int24",   "name": "tickUpper", "type": "int24"   },
				{ "internalType": "int24",   "name": "tickLower", "type": "int24"   },
				{ "internalType": "uint24",  "name": "weight",    "type": "uint24"  }
			],
			"internalType": "struct ICellarPoolShare.CellarTickInfo[]",
			"name": "_cellarTickInfo",
			"type": "tuple[]"
		}
	],
	"name": "rebalance",
	"outputs": [],
	"stateMutability": "nonpayable",
	"type": "function"
}]`

// tuple[] packing example
// https://github.com/ethereum/go-ethereum/blob/master/accounts/abi/packing_test.go#L928

// solidity types to go types examples
// https://github.com/ethereum/go-ethereum/blob/master/accounts/abi/type_test.go#L143

func ABIEncodedCellarTickInfoBytes(index uint) []byte {
	encodedCall, err := abi.JSON(strings.NewReader(cellarTickInfoABI))
	if err != nil {
		panic(sdkerrors.Wrap(err, "bad ABI definition in code"))
	}

	abiEncodedCall, err := encodedCall.Pack("cellarTickInfo", big.NewInt(int64(index)))
	if err != nil {
		panic(err)
	}

	return abiEncodedCall
}

func BytesToABIEncodedTickRange(bz []byte) (*TickRange, error) {
	encodedCall, err := abi.JSON(strings.NewReader(cellarTickInfoABI))
	if err != nil {
		panic(sdkerrors.Wrap(err, "bad ABI definition in code"))
	}

	var abiEncodedTickRange ABIEncodedTickRange
	if err := encodedCall.UnpackIntoInterface(&abiEncodedTickRange, "cellarTickInfo", bz); err != nil {
		return nil, err
	}

	tr := TickRange{
		Upper:  int32(abiEncodedTickRange.Upper.Int64()),
		Lower:  int32(abiEncodedTickRange.Lower.Int64()),
		Weight: uint32(abiEncodedTickRange.Weight.Uint64()),
	}

	return &tr, nil
}

const cellarTickInfoABI = `[{
    "inputs": [
      {
        "internalType": "uint256",
        "name": "",
        "type": "uint256"
      }
    ],
    "name": "cellarTickInfo",
    "outputs": [
      {
        "internalType": "uint184",
        "name": "tokenId",
        "type": "uint184"
      },
      {
        "internalType": "int24",
        "name": "tickUpper",
        "type": "int24"
      },
      {
        "internalType": "int24",
        "name": "tickLower",
        "type": "int24"
      },
      {
        "internalType": "uint24",
        "name": "weight",
        "type": "uint24"
      }
    ],
    "stateMutability": "view",
    "type": "function"
  }]
`