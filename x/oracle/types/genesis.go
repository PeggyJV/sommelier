package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// NewGenesisState creates a new GenesisState object
func NewGenesisState(
	params Params,
	rates sdk.DecCoins,
	feederDelegations []OracleDelegation,
	missCounters []ValidatorMissCounter,
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
		FeederDelegations:             []OracleDelegation{},
		MissCounters:                  []ValidatorMissCounter{},
		AggregateExchangeRatePrevotes: []AggregateExchangeRatePrevote{},
		AggregateExchangeRateVotes:    []AggregateExchangeRateVote{},
		TobinTaxes:                    make(sdk.DecCoins, 0),
	}
}

// Validate validates the oracle genesis state fields.
func (gs GenesisState) Validate() error {
	return gs.Params.ValidateBasic()
}
