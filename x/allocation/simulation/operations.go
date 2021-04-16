package simulation

// DONTCOVER

import (
	"github.com/cosmos/cosmos-sdk/codec"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
	"github.com/cosmos/cosmos-sdk/x/simulation"

	"github.com/peggyjv/sommelier/x/allocation/keeper"
)

// WeightedOperations returns all the operations from the module with their respective weights
func WeightedOperations(
	appParams simtypes.AppParams,
	cdc codec.JSONMarshaler,
	k keeper.Keeper) simulation.WeightedOperations {
	// TODO: reimplement
	return simulation.WeightedOperations{}
}
