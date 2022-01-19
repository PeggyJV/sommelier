package types

import (
	"github.com/ethereum/go-ethereum/common"
)

const DefaultParamspace = ModuleName

// DefaultGenesisState get raw genesis raw message for testing
func DefaultGenesisState() GenesisState {
	return GenesisState{
		Params:    DefaultParams(),
		Addresses: []string{},
	}
}

// Validate performs a basic stateless validation of the genesis fields.
func (gs GenesisState) Validate() error {
	for _, address := range gs.Addresses {
		if common.IsHexAddress(address) {
			return ErrInvalidAddress
		}
	}

	if err := gs.Params.ValidateBasic(); err != nil {
		return err
	}

	return nil
}
