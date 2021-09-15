package types

const DefaultParamspace = ModuleName

// DefaultGenesisState get raw genesis raw message for testing
func DefaultGenesisState() GenesisState {
	return GenesisState{
		Params:  DefaultParams(),
		Cellars: []*Cellar{},
	}
}

// Validate performs a basic stateless validation of the genesis fields.
func (gs GenesisState) Validate() error {
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
