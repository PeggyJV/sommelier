package types

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

const DefaultParamspace = ModuleName

// DefaultGenesisState get raw genesis raw message for testing
func DefaultGenesisState() GenesisState {
	return GenesisState{
		Params: DefaultParams(),
		Cellars: []*Cellar{},
	}
}

// Validate performs a basic stateless validation of the genesis fields.
func (gs GenesisState) Validate() error {
	seenDelegations := make(map[string]bool)

	for i, feederDelegation := range gs.FeederDelegations {
		if seenDelegations[feederDelegation.Validator] {
			return fmt.Errorf("duplicated feeder delegation for validator %s at index %d", feederDelegation.Validator, i)
		}

		delegateAddr, err := sdk.AccAddressFromBech32(feederDelegation.Delegate)
		if err != nil {
			return fmt.Errorf("invalid feeder delegate at index %d: %w", i, err)
		}

		validatorAddr, err := sdk.ValAddressFromBech32(feederDelegation.Validator)
		if err != nil {
			return fmt.Errorf("invalid feeder validator at index %d: %w", i, err)
		}

		if delegateAddr.Equals(validatorAddr) {
			return fmt.Errorf("delegate address %s cannot be equal to validator address %s", feederDelegation.Delegate, feederDelegation.Validator)
		}
		seenDelegations[feederDelegation.Validator] = true
	}

	for _, cellar := range gs.Cellars {
		if err := cellar.ValidateBasic(); err != nil {
			return err
		}
	}

	if err := gs.Params.ValidateBasic(); err != nil {
		return err
	}

	return nil
}
