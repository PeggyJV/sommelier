package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// ExchangeRateVotes represents a
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

// NewAggregateExchangeRateVote creates a AggregateExchangeRateVote instance
func NewAggregateExchangeRateVote(tuples sdk.DecCoins, voter sdk.ValAddress) AggregateExchangeRateVote {
	return AggregateExchangeRateVote{
		ExchangeRateTuples: tuples,
		Voter:              voter.String(),
	}
}
