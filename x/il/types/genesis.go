package types

// NewGenesisState creates a new GenesisState object
func NewGenesisState(
	params Params,
	stoplossPositions []StoplossPosition,
) GenesisState {

	return GenesisState{
		Params:            params,
		StoplossPositions: stoplossPositions,
	}
}

// DefaultGenesisState - default GenesisState used by columbus-2
func DefaultGenesisState() GenesisState {
	return GenesisState{
		Params:            DefaultParams(),
		StoplossPositions: make([]StoplossPosition, 0),
	}
}

// Validate validates the oracle genesis state fields.
func (gs GenesisState) Validate() error {
	if err := gs.StoplossPositions.Validate(); err != nil {
		return err
	}

	return gs.Params.ValidateBasic()
}
