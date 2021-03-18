package types

import (
	fmt "fmt"

	bridgetypes "github.com/althea-net/peggy/module/x/peggy/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// NewGenesisState creates a new GenesisState object
func NewGenesisState(
	params Params,
	stoplossPositions LPsStoplossPositions,
	invalidationID uint64,
	submittedPositionsQueue []SubmittedPosition,
) GenesisState {

	return GenesisState{
		Params:                  params,
		LpsStoplossPositions:    stoplossPositions,
		InvalidationID:          invalidationID,
		SubmittedPositionsQueue: submittedPositionsQueue,
	}
}

// DefaultGenesisState is the default GenesisState used by the IL module
func DefaultGenesisState() GenesisState {
	return GenesisState{
		Params:                  DefaultParams(),
		LpsStoplossPositions:    make(LPsStoplossPositions, 0),
		InvalidationID:          0,
		SubmittedPositionsQueue: make([]SubmittedPosition, 0),
	}
}

// Validate validates the oracle genesis state fields.
func (gs GenesisState) Validate() error {
	if err := gs.LpsStoplossPositions.Validate(); err != nil {
		return err
	}

	for _, position := range gs.SubmittedPositionsQueue {
		if _, err := sdk.AccAddressFromBech32(position.Address); err != nil {
			return err
		}

		if err := bridgetypes.ValidateEthAddress(position.PairId); err != nil {
			return fmt.Errorf("invalid uniswap pair id: %w", err)
		}

		if position.TimeoutHeight == 0 {
			return fmt.Errorf("eth timeout height cannot be 0 for submitted position pair (%s)", position.PairId)
		}
	}

	return gs.Params.ValidateBasic()
}
