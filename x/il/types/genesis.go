package types

import (
	fmt "fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	bridgetypes "github.com/cosmos/gravity-bridge/module/x/gravity/types"
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

	seenSubmittedPositions := make(map[string]bool)
	for _, position := range gs.SubmittedPositionsQueue {
		if seenSubmittedPositions[position.Address+position.PairId] {
			return fmt.Errorf("duplicated submitted position for address %s and pair ID %s", position.Address, position.PairId)
		}

		if _, err := sdk.AccAddressFromBech32(position.Address); err != nil {
			return err
		}

		if err := bridgetypes.ValidateEthAddress(position.PairId); err != nil {
			return fmt.Errorf("invalid uniswap pair id: %w", err)
		}

		if position.TimeoutHeight == 0 {
			return fmt.Errorf("eth timeout height cannot be 0 for submitted position pair (%s)", position.PairId)
		}

		seenSubmittedPositions[position.Address+position.PairId] = true
	}

	return gs.Params.ValidateBasic()
}
