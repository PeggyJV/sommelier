package v2

import types "github.com/peggyjv/sommelier/v8/x/cork/types"

const DefaultParamspace = types.ModuleName

// DefaultGenesisState get raw genesis raw message for testing
func DefaultGenesisState() GenesisState {
	return GenesisState{
		Params:            DefaultParams(),
		CellarIds:         CellarIDSet{},
		InvalidationNonce: 0,
		ScheduledCorks:    []*ScheduledCork{},
		CorkResults:       []*CorkResult{},
	}
}

// Validate performs a basic stateless validation of the genesis fields.
func (gs GenesisState) Validate() error {
	if err := gs.Params.ValidateBasic(); err != nil {
		return err
	}

	if err := gs.CellarIds.ValidateBasic(); err != nil {
		return err
	}

	for _, scheduledCork := range gs.ScheduledCorks {
		if err := scheduledCork.ValidateBasic(); err != nil {
			return err
		}
	}

	for _, corkResult := range gs.CorkResults {
		if err := corkResult.ValidateBasic(); err != nil {
			return err
		}
	}

	return nil
}
