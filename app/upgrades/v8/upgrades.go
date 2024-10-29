package v8

import (
	"github.com/cosmos/cosmos-sdk/baseapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	consensusparamkeeper "github.com/cosmos/cosmos-sdk/x/consensus/keeper"
	paramstypes "github.com/cosmos/cosmos-sdk/x/params/types"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"
	ibcexported "github.com/cosmos/ibc-go/v7/modules/core/exported"
	ibckeeper "github.com/cosmos/ibc-go/v7/modules/core/keeper"
)

func CreateUpgradeHandler(
	mm *module.Manager,
	configurator module.Configurator,
	baseAppLegacySS *paramstypes.Subspace,
	consensusParamsKeeper *consensusparamkeeper.Keeper,
	ibcKeeper *ibckeeper.Keeper,
) upgradetypes.UpgradeHandler {
	return func(ctx sdk.Context, plan upgradetypes.Plan, vm module.VersionMap) (module.VersionMap, error) {
		ctx.Logger().Info("v8 upgrade: entering handler and running migrations")

		// new x/consensus module params migration
		baseapp.MigrateParams(ctx, baseAppLegacySS, consensusParamsKeeper)

		// explicitly update the IBC 02-client params, adding the localhost client type
		params := ibcKeeper.ClientKeeper.GetParams(ctx)
		params.AllowedClients = append(params.AllowedClients, ibcexported.Localhost)
		ibcKeeper.ClientKeeper.SetParams(ctx, params)

		return mm.RunMigrations(ctx, configurator, vm)
	}
}
