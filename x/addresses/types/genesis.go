package types

import fmt "fmt"

// DefaultIndex is the default global index
const DefaultIndex uint64 = 1

// DefaultGenesis returns the default genesis state
func DefaultGenesisState() GenesisState {
	return GenesisState{
		Params:          DefaultParams(),
		AddressMappings: []*AddressMapping{},
	}
}

// Validate performs basic genesis state validation returning an error upon any
// failure.
func (gs GenesisState) Validate() error {
	if err := gs.Params.ValidateBasic(); err != nil {
		return err
	}

	seenMappings := make(map[string]string)
	for _, mapping := range gs.AddressMappings {
		if err := mapping.ValidateBasic(); err != nil {
			return err
		}

		// Check for duplicate mappings
		key := mapping.CosmosAddress + "|" + mapping.EvmAddress
		if _, exists := seenMappings[key]; exists {
			return fmt.Errorf("duplicate address mapping found: %s", key)
		}
		seenMappings[key] = ""
	}

	return nil
}
