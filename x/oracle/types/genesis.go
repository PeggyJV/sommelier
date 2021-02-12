package types

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// NewGenesisState returns the genesis state struct
func NewGenesisState() *GenesisState {
	return &GenesisState{}
}

// DefaultGenesisState get raw genesis raw message for testing
func DefaultGenesisState() *GenesisState {
	return &GenesisState{
		Params: DefaultParams(),
	}
}

// Validate performs a basic stateless validation of the genesis fields.
func (gs GenesisState) Validate() error {
	for i, feederDelegation := range gs.FeederDelegations {
		if _, err := sdk.AccAddressFromBech32(feederDelegation.Delegate); err != nil {
			return fmt.Errorf("invalid feeder at index %d: %w", i, err)
		}
		if _, err := sdk.AccAddressFromBech32(feederDelegation.Validator); err != nil {
			return fmt.Errorf("invalid feeder at index %d: %w", i, err)
		}
	}

	for i, missCounter := range gs.MissCounters {
		if missCounter.Misses < 0 {
			return fmt.Errorf("miss counter for validator %s cannot be negative: %d", missCounter.Validator, missCounter.Misses)
		}
		if _, err := sdk.AccAddressFromBech32(missCounter.Validator); err != nil {
			return fmt.Errorf("invalid feeder at index %d: %w", i, err)
		}
	}

	return gs.Params.ValidateBasic()
}
