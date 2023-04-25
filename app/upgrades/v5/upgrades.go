package v5

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"
	appparams "github.com/peggyjv/sommelier/v6/app/params"
	incentiveskeeper "github.com/peggyjv/sommelier/v6/x/incentives/keeper"
	incentivestypes "github.com/peggyjv/sommelier/v6/x/incentives/types"
)

func CreateUpgradeHandler(
	mm *module.Manager,
	configurator module.Configurator,
	incentivesKeeper incentiveskeeper.Keeper,
) upgradetypes.UpgradeHandler {
	return func(ctx sdk.Context, plan upgradetypes.Plan, fromVM module.VersionMap) (module.VersionMap, error) {
		ctx.Logger().Info("v5 upgrade: entering handler")

		// We must manually run InitGenesis for incentives so we can adjust their values during the upgrade process.
		// Setting the consensus version to 1 prevents RunMigrations from running InitGenesis itself.
		fromVM[incentivestypes.ModuleName] = mm.Modules[incentivestypes.ModuleName].ConsensusVersion()

		ctx.Logger().Info("v5 upgrade: initializing incentives genesis state")
		incentivesInitGenesis(ctx, incentivesKeeper)

		ctx.Logger().Info("v5 upgrade: running migrations and exiting handler")
		return mm.RunMigrations(ctx, configurator, fromVM)
	}
}

// Launch the incentives module with 2 SOMM per block distribution and a cutoff height 5 million blocks past
// the upgrade height
func incentivesInitGenesis(ctx sdk.Context, incentivesKeeper incentiveskeeper.Keeper) {
	genesisState := incentivestypes.DefaultGenesisState()

	upgradeHeight := uint64(7766725)
	incentivesBlocks := uint64(5000000)

	params := incentivestypes.Params{
		DistributionPerBlock:   sdk.NewCoin(appparams.BaseCoinUnit, sdk.NewInt(2000000)),
		IncentivesCutoffHeight: upgradeHeight + incentivesBlocks,
	}
	genesisState.Params = params

	incentiveskeeper.InitGenesis(ctx, incentivesKeeper, genesisState)
}
