package types

const DefaultParamspace = ModuleName

// DefaultGenesisState get raw genesis raw message for testing
func DefaultGenesisState() GenesisState {
	defaultParams := DefaultParams()
	return GenesisState{
		Params:              &defaultParams,
		ChainConfigurations: ChainConfigurations{},
		CellarIds:           []*CellarIDSet{},
		ScheduledCorks:      []*ScheduledAxelarCorks{},
		CorkResults:         []*AxelarCorkResults{},
	}
}

// Validate performs a basic stateless validation of the genesis fields.
func (gs GenesisState) Validate() error {
	if err := gs.Params.ValidateBasic(); err != nil {
		return err
	}

	for _, cc := range gs.ChainConfigurations.Configurations {
		if err := cc.ValidateBasic(); err != nil {
			return err
		}
	}

	// todo (mvid): flesh out initial validation

	return nil
}
