package types

// NewGenesisState returns the genesis state struct
func NewGenesisState() *GenesisState {
	return &GenesisState{}
}

// // get raw genesis raw message for testing
// func DefaultGenesisState() *GenesisState {
// 	return &GenesisState{,
// 		Params:                          DefaultParams(),
// 	}
// }

// // ValidateGenesis validates the genesis state of distribution genesis input
// func ValidateGenesis(gs *GenesisState) error {
// 	if err := gs.Params.ValidateBasic(); err != nil {
// 		return err
// 	}
// 	return gs.FeePool.ValidateGenesis()
// }
