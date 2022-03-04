package types

import (
	"encoding/hex"
	"math/big"

	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

type ABIEncodedTickRange struct {
	TokenID *big.Int `abi:"tokenId"`
	Upper   *big.Int `abi:"tickUpper"`
	Lower   *big.Int `abi:"tickLower"`
	Weight  *big.Int `abi:"weight"`
}

// ABIEncodedRebalanceBytes gets the checkpoint signature from the given outgoing tx batch
func (rv RebalanceVote) ABIEncodedRebalanceBytes() []byte {
	encodedCall, err := abiJSON()
	if err != nil {
		panic(sdkerrors.Wrap(err, "bad ABI definition in code"))
	}

	ticks := make([]ABIEncodedTickRange, len(rv.Cellar.TickRanges))
	for i, t := range rv.Cellar.TickRanges {
		up := int64(t.Upper)
		lo := int64(t.Lower)
		we := uint64(t.Weight)
		ticks[i] = ABIEncodedTickRange{big.NewInt(0), big.NewInt(up), big.NewInt(lo), new(big.Int).SetUint64(we)}
	}

	abiEncodedCall, err := encodedCall.Pack("rebalance", ticks, new(big.Int).SetUint64(rv.CurrentPrice))
	if err != nil {
		panic(err)
	}
	println("rebalance payload: ", hex.EncodeToString(abiEncodedCall))

	return abiEncodedCall
}

// tuple[] packing example
// https://github.com/ethereum/go-ethereum/blob/master/accounts/abi/packing_test.go#L928

// solidity types to go types examples
// https://github.com/ethereum/go-ethereum/blob/master/accounts/abi/type_test.go#L143

func ABIEncodedCellarTickInfoBytes(index uint) []byte {
	encodedCall, err := abiJSON()
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
	encodedCall, err := abiJSON()
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
