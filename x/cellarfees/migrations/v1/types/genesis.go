package types

import (
	"sort"

	sdk "github.com/cosmos/cosmos-sdk/types"
	types "github.com/peggyjv/sommelier/v7/x/cellarfees/types"
)

const DefaultParamspace = ModuleName

// DefaultGenesisState get raw genesis raw message for testing
func DefaultGenesisState() GenesisState {
	return GenesisState{
		Params:               DefaultParams(),
		FeeAccrualCounters:   DefaultFeeAccrualCounters(),
		LastRewardSupplyPeak: sdk.ZeroInt(),
	}
}

// Validate performs a basic stateless validation of the genesis fields.
func (gs GenesisState) Validate() error {
	if err := gs.Params.ValidateBasic(); err != nil {
		return err
	}

	if gs.FeeAccrualCounters.Counters == nil {
		return ErrInvalidFeeAccrualCounters.Wrap("counters cannot be nil!")
	}

	counters := gs.FeeAccrualCounters
	counters.Counters = append([]FeeAccrualCounter{}, gs.FeeAccrualCounters.Counters...)
	sort.Sort(counters)
	for i := range counters.Counters {
		if counters.Counters[i].Denom != gs.FeeAccrualCounters.Counters[i].Denom {
			return ErrInvalidFeeAccrualCounters.Wrapf("counters are unsorted! expected: %T, actual: %T", counters.Counters, gs.FeeAccrualCounters.Counters)
		}
	}

	if gs.LastRewardSupplyPeak.LT(sdk.ZeroInt()) {
		return types.ErrInvalidLastRewardSupplyPeak.Wrap("last reward supply peak cannot be less than zero!")
	}

	return nil
}
