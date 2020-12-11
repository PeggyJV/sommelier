package types

import (
	"fmt"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

type ExchangeRateVotes []ExchangeRateVote

// NewExchangeRateVote creates a ExchangeRateVote instance
func NewExchangeRateVote(rate sdk.Dec, denom string, voter sdk.ValAddress) ExchangeRateVote {
	return ExchangeRateVote{
		ExchangeRate: rate,
		Denom:        denom,
		Voter:        voter.String(),
	}
}

// NewAggregateExchangeRatePrevote returns AggregateExchangeRatePrevote object
func NewAggregateExchangeRatePrevote(hash AggregateVoteHash, voter sdk.ValAddress, submitBlock int64) AggregateExchangeRatePrevote {
	return AggregateExchangeRatePrevote{
		Hash:        hash,
		Voter:       voter.String(),
		SubmitBlock: submitBlock,
	}
}

// ParseExchangeRateTuples ExchangeRateTuple parser
func ParseExchangeRateTuples(tuplesStr string) (sdk.DecCoins, error) {
	tuplesStr = strings.TrimSpace(tuplesStr)
	if len(tuplesStr) == 0 {
		return nil, nil
	}

	tupleStrs := strings.Split(tuplesStr, ",")
	tuples := make(sdk.DecCoins, len(tupleStrs))
	duplicateCheckMap := make(map[string]bool)
	for i, tupleStr := range tupleStrs {
		decCoin, err := sdk.ParseDecCoin(tupleStr)
		if err != nil {
			return nil, err
		}
		tuples[i] = decCoin
		if _, ok := duplicateCheckMap[decCoin.Denom]; ok {
			return nil, fmt.Errorf("duplicated denom %s", decCoin.Denom)
		}
		duplicateCheckMap[decCoin.Denom] = true
	}
	return tuples, nil
}

// NewAggregateExchangeRateVote creates a AggregateExchangeRateVote instance
func NewAggregateExchangeRateVote(tuples sdk.DecCoins, voter sdk.ValAddress) AggregateExchangeRateVote {
	return AggregateExchangeRateVote{
		ExchangeRateTuples: tuples,
		Voter:              voter.String(),
	}
}
