package oracle

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/peggyjv/sommelier/x/il/keeper"
	"github.com/peggyjv/sommelier/x/il/types"
)

// InitGenesis initialize default parameters and sets the stoploss positions to store
func InitGenesis(ctx sdk.Context, k keeper.Keeper, data types.GenesisState) {
	k.SetParams(ctx, data.Params)
	// set positions for each LP address
	for _, lpStoplossPositions := range data.LpsStoplossPositions {
		for _, position := range lpStoplossPositions.StoplossPositions {
			k.SetStoplossPosition(ctx, sdk.AccAddress(lpStoplossPositions.Address), position)
		}
	}
}

// ExportGenesis writes the current store values
// to a genesis file, which can be imported again
// with InitGenesis
func ExportGenesis(ctx sdk.Context, k keeper.Keeper) types.GenesisState {
	return types.GenesisState{
		Params:               k.GetParams(ctx),
		LpsStoplossPositions: nil, // TODO:
	}
}
