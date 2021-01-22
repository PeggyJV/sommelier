package types

// NewGenesisState creates a new GenesisState object
func NewGenesisState(
	params Params,
	stoplossPositions LPsStoplossPositions,
) GenesisState {

	return GenesisState{
		Params:               params,
		LpsStoplossPositions: stoplossPositions,
	}
}

// DefaultGenesisState - default GenesisState used by columbus-2
func DefaultGenesisState() GenesisState {
	return GenesisState{
		Params:               DefaultParams(),
		LpsStoplossPositions: make(LPsStoplossPositions, 0),
	}
}

// Validate validates the oracle genesis state fields.
func (gs GenesisState) Validate() error {
	if err := gs.LpsStoplossPositions.Validate(); err != nil {
		return err
	}

	return gs.Params.ValidateBasic()
}
