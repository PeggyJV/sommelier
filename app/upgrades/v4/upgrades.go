package v4

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"
	gravitytypes "github.com/peggyjv/gravity-bridge/module/x/gravity/types"
)

func CreateUpgradeHandler(mm *module.Manager, configurator module.Configurator) upgradetypes.UpgradeHandler {
	fmt.Println("v4 upgrade: entering handler")
	return func(ctx sdk.Context, plan upgradetypes.Plan, vm module.VersionMap) (module.VersionMap, error) {
		// Since this is the first in-place upgrade and InitChainer was not set up for this at genesis
		// time, we must initialize the VM map ourselves.
		fromVM := make(map[string]uint64)
		for moduleName, module := range mm.Modules {
			fromVM[moduleName] = module.ConsensusVersion()
		}

		// Overwrite the gravity module's version back to 1 so the migration will run to v2
		fromVM[gravitytypes.ModuleName] = 1

		ctx.Logger().Info("v4 upgrade: running migrations and exiting handler")
		return mm.RunMigrations(ctx, configurator, fromVM)
	}
}
