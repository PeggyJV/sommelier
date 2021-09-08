package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/peggyjv/sommelier/x/allocation/types"
)

// InitGenesis initialize default parameters
// and the keeper's address to pubkey map
func InitGenesis(ctx sdk.Context, k Keeper, gs types.GenesisState) {
	k.setParams(ctx, gs.Params)
	// Set the vote period at initialization
	k.SetCommitPeriodStart(ctx, ctx.BlockHeight())

	for _, delegation := range gs.FeederDelegations {
		// NOTE: error checked during genesis validation
		valAddress, _ := sdk.ValAddressFromBech32(delegation.Validator)
		delAddress, _ := sdk.AccAddressFromBech32(delegation.Delegate)

		k.SetValidatorDelegateAddress(ctx, delAddress, valAddress)
	}

	for _, cellar := range gs.Cellars {
		k.SetCellar(ctx, *cellar)
	}
}

// ExportGenesis writes the current store values
// to a genesis file, which can be imported again
// with InitGenesis
func ExportGenesis(ctx sdk.Context, k Keeper) types.GenesisState {
	var cellars []*types.Cellar
	for _, cellar := range k.GetCellars(ctx) {
		cellars = append(cellars, &cellar)
	}

	return types.GenesisState{
		Params:            k.GetParamSet(ctx),
		FeederDelegations: k.GetAllAllocationDelegations(ctx),
		Cellars: cellars,
	}
}
