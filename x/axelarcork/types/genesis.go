package types

const DefaultParamspace = ModuleName

// DefaultGenesisState get raw genesis raw message for testing
func DefaultGenesisState() GenesisState {
	defaultParams := DefaultParams()
	return GenesisState{
		Params:                   &defaultParams,
		ChainConfigurations:      ChainConfigurations{},
		CellarIds:                []*CellarIDSet{},
		ScheduledCorks:           &ScheduledAxelarCorks{},
		CorkResults:              &AxelarCorkResults{},
		AxelarContractCallNonces: []*AxelarContractCallNonce{},
		AxelarUpgradeData:        []*AxelarUpgradeData{},
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

	for _, sc := range gs.ScheduledCorks.ScheduledCorks {
		if err := sc.Cork.ValidateBasic(); err != nil {
			return err
		}
	}

	for _, cr := range gs.CorkResults.CorkResults {
		if err := cr.Cork.ValidateBasic(); err != nil {
			return err
		}
	}

	for _, ccn := range gs.AxelarContractCallNonces {
		if err := ccn.ValidateBasic(); err != nil {
			return err
		}
	}

	for _, aud := range gs.AxelarUpgradeData {
		if err := aud.ValidateBasic(); err != nil {
			return err
		}
	}

	return nil
}
