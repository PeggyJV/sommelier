package types

// this line is used by starport scaffolding # genesis/types/import

// DefaultIndex is the default global index
const DefaultIndex uint64 = 1

// DefaultGenesis returns the default genesis state
func DefaultGenesis() *GenesisState {
	return &GenesisState{
		Params:          DefaultParams(),
		AddressMappings: []*AddressMapping{},
	}
}

// Validate performs basic genesis state validation returning an error upon any
// failure.
func (gs GenesisState) Validate() error {
	gs.Params.ValidateBasic()

	for _, mapping := range gs.AddressMappings {
		if err := mapping.ValidateBasic(); err != nil {
			return err
		}
	}

	return nil
}
