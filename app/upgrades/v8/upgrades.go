package v8

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"
)

func CreateUpgradeHandler(
	mm *module.Manager,
	configurator module.Configurator,
) upgradetypes.UpgradeHandler {
	return func(ctx sdk.Context, plan upgradetypes.Plan, vm module.VersionMap) (module.VersionMap, error) {
		ctx.Logger().Info("v2 upgrade: entering handler and running migrations")

		fromVM := make(map[string]uint64)
		for moduleName, m := range mm.Modules {
			fromVM[moduleName] = m.(module.HasConsensusVersion).ConsensusVersion()
		}

		return mm.RunMigrations(ctx, configurator, fromVM)
	}
}
