package types

import (
	"fmt"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// NewAggregateExchangeRatePrevote returns AggregateExchangeRatePrevote object
func NewAggregateExchangeRatePrevote(hash []byte, voter sdk.ValAddress, submitBlock int64) *AggregateExchangeRatePrevote {
	return &AggregateExchangeRatePrevote{
		Hash:        hash,
		SubmitBlock: submitBlock,
		Voter:       voter.String(),
	}
}

// NewAggregateExchangeRateVote creates a AggregateExchangeRateVote instance
func NewAggregateExchangeRateVote(tuples ExchangeRateTuples, voter sdk.ValAddress) *AggregateExchangeRateVote {
	return &AggregateExchangeRateVote{
		ExchangeRateTuples: tuples,
		Voter:              voter.String(),
	}
}

// ExchangeRateTuples - array of ExchangeRateTuple
type ExchangeRateTuples []*ExchangeRateTuple

// ParseExchangeRateTuples ExchangeRateTuple parser
func ParseExchangeRateTuples(tuplesStr string) (ExchangeRateTuples, error) {
	tuplesStr = strings.TrimSpace(tuplesStr)
	if len(tuplesStr) == 0 {
		return nil, nil
	}

	tupleStrs := strings.Split(tuplesStr, ",")
	tuples := make(ExchangeRateTuples, len(tupleStrs))
	duplicateCheckMap := make(map[string]bool)
	for i, tupleStr := range tupleStrs {
		decCoin, err := sdk.ParseDecCoin(tupleStr)
		if err != nil {
			return nil, err
		}

		tuples[i] = &ExchangeRateTuple{
			Denom:        decCoin.Denom,
			ExchangeRate: decCoin.Amount,
		}

		if _, ok := duplicateCheckMap[decCoin.Denom]; ok {
			return nil, fmt.Errorf("duplicated denom %s", decCoin.Denom)
		}

		duplicateCheckMap[decCoin.Denom] = true
	}

	return tuples, nil
}
