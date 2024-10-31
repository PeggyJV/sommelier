package v8

import (
	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	authkeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	consensusparamkeeper "github.com/cosmos/cosmos-sdk/x/consensus/keeper"
	paramstypes "github.com/cosmos/cosmos-sdk/x/params/types"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"
	ibcexported "github.com/cosmos/ibc-go/v7/modules/core/exported"
	ibckeeper "github.com/cosmos/ibc-go/v7/modules/core/keeper"
	ibctmmigrations "github.com/cosmos/ibc-go/v7/modules/light-clients/07-tendermint/migrations"
	auctiontypes "github.com/peggyjv/sommelier/v8/x/auction/types"
)

func CreateUpgradeHandler(
	mm *module.Manager,
	configurator module.Configurator,
	baseAppLegacySS *paramstypes.Subspace,
	consensusParamsKeeper *consensusparamkeeper.Keeper,
	ibcKeeper *ibckeeper.Keeper,
	cdc codec.BinaryCodec,
	clientKeeper ibctmmigrations.ClientKeeper,
	accountKeeper *authkeeper.AccountKeeper,
) upgradetypes.UpgradeHandler {
	return func(ctx sdk.Context, plan upgradetypes.Plan, vm module.VersionMap) (module.VersionMap, error) {
		ctx.Logger().Info("v8 upgrade: entering handler and running migrations")

		// add burner permission to auction account
		newAcct := authtypes.NewEmptyModuleAccount(auctiontypes.ModuleName, authtypes.Burner)
		oldAcctI := accountKeeper.GetModuleAccount(ctx, auctiontypes.ModuleName)
		newAcct.AccountNumber = oldAcctI.GetAccountNumber()
		newAcct.Address = oldAcctI.GetAddress().String()
		newAcct.Sequence = oldAcctI.GetSequence()
		newAcct.Permissions = oldAcctI.GetPermissions()
		newAcct.Name = oldAcctI.GetName()
		newAcctI := (accountKeeper.NewAccount(ctx, newAcct)).(authtypes.ModuleAccountI)
		accountKeeper.SetModuleAccount(ctx, newAcctI)

		// Include this when migrating to ibc-go v7 (optional)
		// source: https://github.com/cosmos/ibc-go/blob/v7.2.0/docs/migrations/v6-to-v7.md
		// prune expired tendermint consensus states to save storage space
		if _, err := ibctmmigrations.PruneExpiredConsensusStates(ctx, cdc, clientKeeper); err != nil {
			return nil, err
		}

		// new x/consensus module params migration
		baseapp.MigrateParams(ctx, baseAppLegacySS, consensusParamsKeeper)

		// explicitly update the IBC 02-client params, adding the localhost client type
		params := ibcKeeper.ClientKeeper.GetParams(ctx)
		params.AllowedClients = append(params.AllowedClients, ibcexported.Localhost)
		ibcKeeper.ClientKeeper.SetParams(ctx, params)

		return mm.RunMigrations(ctx, configurator, vm)
	}
}
