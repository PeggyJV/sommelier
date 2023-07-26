package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/peggyjv/sommelier/v6/x/axelarcork/types"
)

// InitGenesis initialize default parameters
// and the keeper's address to pubkey map
func InitGenesis(ctx sdk.Context, k Keeper, gs types.GenesisState) {
	k.SetParams(ctx, *gs.Params)
	// Set the vote period at initialization

	for i, config := range gs.ChainConfigurations.Configurations {
		k.SetChainConfiguration(ctx, config.Id, *config)
		k.SetCellarIDs(ctx, config.Id, *gs.CellarIds[i])
	}

	for _, corkResult := range gs.CorkResults.CorkResults {
		k.SetAxelarCorkResult(
			ctx,
			corkResult.Cork.ChainId,
			corkResult.Cork.IDHash(
				corkResult.Cork.ChainId,
				corkResult.BlockHeight),
			*corkResult,
		)
	}

	for _, scheduledCork := range gs.ScheduledCorks.ScheduledCorks {
		valAddr, err := sdk.ValAddressFromHex(scheduledCork.Validator)
		if err != nil {
			panic(err)
		}

		k.SetScheduledAxelarCork(ctx, scheduledCork.Cork.ChainId, scheduledCork.BlockHeight, valAddr, *scheduledCork.Cork)
	}
}

// ExportGenesis writes the current store values
// to a genesis file, which can be imported again
// with InitGenesis
func ExportGenesis(ctx sdk.Context, k Keeper) types.GenesisState {
	var gs types.GenesisState

	ps := k.GetParamSet(ctx)
	gs.Params = &ps

	k.IterateChainConfigurations(ctx, func(config types.ChainConfiguration) (stop bool) {
		gs.ChainConfigurations.Configurations = append(gs.ChainConfigurations.Configurations, &config)

		cellarIDs := k.GetCellarIDs(ctx, config.Id)
		var cellarIDSet types.CellarIDSet
		for _, id := range cellarIDs {
			cellarIDSet.Ids = append(cellarIDSet.Ids, id.String())
		}
		gs.CellarIds = append(gs.CellarIds, &cellarIDSet)

		gs.ScheduledCorks.ScheduledCorks = append(gs.ScheduledCorks.ScheduledCorks, k.GetScheduledAxelarCorks(ctx, config.Id)...)
		gs.CorkResults.CorkResults = append(gs.CorkResults.CorkResults, k.GetAxelarCorkResults(ctx, config.Id)...)

		return false
	})

	return gs
}
