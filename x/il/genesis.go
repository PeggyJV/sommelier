package il

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/peggyjv/sommelier/x/il/keeper"
	"github.com/peggyjv/sommelier/x/il/types"
)

// InitGenesis initialize default parameters and sets the stoploss positions to store
func InitGenesis(ctx sdk.Context, k keeper.Keeper, data types.GenesisState) {
	k.SetParams(ctx, data.Params)
	k.SetInvalidationID(ctx, data.InvalidationID)

	// set positions for each LP address
	for _, lpStoplossPositions := range data.LpsStoplossPositions {
		for _, position := range lpStoplossPositions.StoplossPositions {
			k.SetStoplossPosition(ctx, sdk.AccAddress(lpStoplossPositions.Address), position)
		}
	}

	for _, position := range data.SubmittedPositionsQueue {
		// NOTE: error checked during genesis validation
		address, _ := sdk.AccAddressFromBech32(position.Address)
		k.SetSubmittedPosition(ctx, position.TimeoutHeight, address, position.PairId)
	}
}

// ExportGenesis writes the current store values
// to a genesis file, which can be imported again
// with InitGenesis
func ExportGenesis(ctx sdk.Context, k keeper.Keeper) types.GenesisState {
	return types.GenesisState{
		Params:                  k.GetParams(ctx),
		LpsStoplossPositions:    k.GetLPsStoplossPositions(ctx),
		InvalidationID:          k.GetInvalidationID(ctx),
		SubmittedPositionsQueue: k.GetSubmittedQueue(ctx),
	}
}
