package types

const DefaultParamspace = ModuleName

// DefaultGenesisState get raw genesis raw message for testing
func DefaultGenesisState() GenesisState {
	defaultParams := DefaultParams()
	return GenesisState{
		Params: &defaultParams,
	}
}

// Validate performs a basic stateless validation of the genesis fields.
func (gs GenesisState) Validate() error {
	if err := gs.Params.ValidateBasic(); err != nil {
		return err
	}

	return nil
}
