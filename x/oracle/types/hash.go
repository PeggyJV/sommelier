package types

import (
	"encoding/hex"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/tendermint/tendermint/crypto/tmhash"
)

// AggregateVoteHash is hash value to hide vote exchange rates
// which is formatted as hex string in SHA256("{salt}:{exchange rate}{denom},...,{exchange rate}{denom}:{voter}")

// GetAggregateVoteHash computes hash value of ExchangeRateVote
// to avoid redundant DecCoins stringify operation, use string argument
func GetAggregateVoteHash(salt string, exchangeRatesStr string, voter sdk.ValAddress) []byte {
	hash := tmhash.NewTruncated()
	sourceStr := fmt.Sprintf("%s:%s:%s", salt, exchangeRatesStr, voter.String())
	_, err := hash.Write([]byte(sourceStr))
	if err != nil {
		panic(err)
	}
	bz := hash.Sum(nil)
	return bz
}

// AggregateVoteHashFromHexString convert hex string to AggregateVoteHash
func AggregateVoteHashFromHexString(s string) ([]byte, error) {
	h, err := hex.DecodeString(s)
	if err != nil {
		return nil, err
	}
	return h, nil
}
