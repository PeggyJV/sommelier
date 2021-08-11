package types

import (
	"math/big"
	"strings"

	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/crypto"
)

// GetCheckpoint gets the checkpoint signature from the given outgoing tx batch
func (c Cellar) GetCheckpoint() []byte {
	encodedCall, err := abi.JSON(strings.NewReader(rebalanceABI))
	if err != nil {
		panic(sdkerrors.Wrap(err, "bad ABI definition in code"))
	}

	type packTick struct {
		TokenId *big.Int `abi:"tokenId"`
		Upper   *big.Int `abi:"tickUpper"`
		Lower   *big.Int `abi:"tickLower"`
		Weight  *big.Int `abi:"weight"`
	}

	var ticks []packTick
	for _, t := range c.TickRanges {
		up := int64(t.Upper)
		lo := int64(t.Lower)
		we := uint64(t.Weight)
		ticks = append(ticks, packTick{big.NewInt(0), big.NewInt(up), big.NewInt(lo), new(big.Int).SetUint64(we)})
	}

	// the methodName needs to be the same as the 'name' above in the checkpointAbiJson
	// but other than that it's a constant that has no impact on the output. This is because
	// it gets encoded as a function name which we must then discard.
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
