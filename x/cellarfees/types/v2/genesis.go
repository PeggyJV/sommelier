package v2

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/peggyjv/sommelier/v9/x/cellarfees/types"
)

const DefaultParamspace = types.ModuleName

// DefaultGenesisState get raw genesis raw message for testing
func DefaultGenesisState() GenesisState {
	return GenesisState{
		Params:               DefaultParams(),
		LastRewardSupplyPeak: sdk.ZeroInt(),
	}
}

// Validate performs a basic stateless validation of the genesis fields.
func (gs GenesisState) Validate() error {
	if err := gs.Params.ValidateBasic(); err != nil {
		return err
	}

	if gs.LastRewardSupplyPeak.LT(sdk.ZeroInt()) {
		return types.ErrInvalidLastRewardSupplyPeak.Wrap("last reward supply peak cannot be less than zero!")
	}

	return nil
}
