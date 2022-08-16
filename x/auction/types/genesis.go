package types

const DefaultParamspace = ModuleName

// DefaultGenesisState get raw genesis raw message for testing
func DefaultGenesisState() GenesisState {
	return GenesisState{
		Params: DefaultParams(),
	}
}

// Validate performs a basic stateless validation of the genesis fields.
func (gs GenesisState) Validate() error {
	// TODO: validate auctions, bids, tokenprices

	if err := gs.Params.ValidateBasic(); err != nil {
		return err
	}

	return nil
}
