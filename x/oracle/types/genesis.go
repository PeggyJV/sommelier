package types

import (
	"fmt"

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
	tobinTaxes sdk.DecCoins,
) GenesisState {

	return GenesisState{
		Params:                        params,
		ExchangeRates:                 rates,
		FeederDelegations:             feederDelegations,
		MissCounters:                  missCounters,
		AggregateExchangeRatePrevotes: aggregateExchangeRatePrevotes,
		AggregateExchangeRateVotes:    aggregateExchangeRateVotes,
		TobinTaxes:                    tobinTaxes,
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
	if !gs.ExchangeRates.IsValid() {
		return fmt.Errorf("invalid exchange rates coins: %s", gs.ExchangeRates)
	}

	for _, delegation := range gs.FeederDelegations {
		_, err := sdk.ValAddressFromBech32(delegation.DelegatorAddress)
		if err != nil {
			return err
		}

		_, err = sdk.AccAddressFromBech32(delegation.DelegateAddress)
		if err != nil {
			return err
		}
	}

	for _, missCounter := range gs.MissCounters {
		_, err := sdk.ValAddressFromBech32(missCounter.ValAddress)
		if err != nil {
			return err
		}
	}

	for _, prevote := range gs.AggregateExchangeRatePrevotes {
		_, err := sdk.ValAddressFromBech32(prevote.Voter)
		if err != nil {
			return err
		}
		if prevote.SubmitBlock < 0 {
			return fmt.Errorf("prevote submit block cannot be negative, got %d", prevote.SubmitBlock)
		}
	}

	for _, vote := range gs.AggregateExchangeRateVotes {
		_, err := sdk.ValAddressFromBech32(vote.Voter)
		if err != nil {
			return err
		}

		if !vote.ExchangeRateTuples.IsValid() {
			return fmt.Errorf("invalid exchange rates tuple coins: %s", vote.ExchangeRateTuples)
		}
	}

	return gs.Params.ValidateBasic()
}
