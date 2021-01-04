package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// NewGenesisState creates a new GenesisState object
func NewGenesisState(
	params Params,
	rates sdk.DecCoins,
	feederDelegations map[string]string,
	missCounters map[string]int64,
	aggregateExchangeRatePrevotes []AggregateExchangeRatePrevote,
	aggregateExchangeRateVotes []AggregateExchangeRateVote,
	TobinTaxes sdk.DecCoins,
) GenesisState {

	return GenesisState{
		Params:                        params,
		ExchangeRates:                 rates,
		FeederDelegations:             feederDelegations,
		MissCounters:                  missCounters,
		AggregateExchangeRatePrevotes: aggregateExchangeRatePrevotes,
		AggregateExchangeRateVotes:    aggregateExchangeRateVotes,
		TobinTaxes:                    TobinTaxes,
	}
}

// DefaultGenesisState - default GenesisState used by columbus-2
func DefaultGenesisState() GenesisState {
	return GenesisState{
		Params:                        DefaultParams(),
		ExchangeRates:                 make(sdk.DecCoins, 0),
		FeederDelegations:             make(map[string]string),
		MissCounters:                  make(map[string]int64),
		AggregateExchangeRatePrevotes: []AggregateExchangeRatePrevote{},
		AggregateExchangeRateVotes:    []AggregateExchangeRateVote{},
		TobinTaxes:                    make(sdk.DecCoins, 0),
	}
}

// ValidateGenesis validates the oracle genesis parameters
func ValidateGenesis(data GenesisState) error {
	return data.Params.ValidateBasic()
}
