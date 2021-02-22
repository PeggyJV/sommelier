package oracle

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/peggyjv/sommelier/x/oracle/keeper"
	"github.com/peggyjv/sommelier/x/oracle/types"
)

// InitGenesis initialize default parameters
// and the keeper's address to pubkey map
func InitGenesis(ctx sdk.Context, k keeper.Keeper, gs types.GenesisState) {
	k.SetParams(ctx, gs.Params)
	// Set the vote period at initialization
	k.SetVotePeriodStart(ctx, ctx.BlockHeight())

	for _, missCounter := range gs.MissCounters {
		// NOTE: error checked during genesis validation
		valAddress, _ := sdk.ValAddressFromBech32(missCounter.Validator)
		k.SetMissCounter(ctx, valAddress, missCounter.Misses)
	}

	for _, delegation := range gs.FeederDelegations {
		// NOTE: error checked during genesis validation
		valAddress, _ := sdk.ValAddressFromBech32(delegation.Validator)
		delAddress, _ := sdk.AccAddressFromBech32(delegation.Delegate)

		k.SetValidatorDelegateAddress(ctx, delAddress, valAddress)
	}
}

// ExportGenesis writes the current store values
// to a genesis file, which can be imported again
// with InitGenesis
func ExportGenesis(ctx sdk.Context, k keeper.Keeper) types.GenesisState {
	return types.GenesisState{
		Params:            k.GetParamSet(ctx),
		FeederDelegations: k.GetAllFeederDelegations(ctx),
		MissCounters:      k.GetAllMissCounters(ctx),
	}
}
