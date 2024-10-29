package v8

import (
	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	consensusparamkeeper "github.com/cosmos/cosmos-sdk/x/consensus/keeper"
	paramstypes "github.com/cosmos/cosmos-sdk/x/params/types"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"
	ibcexported "github.com/cosmos/ibc-go/v7/modules/core/exported"
	ibckeeper "github.com/cosmos/ibc-go/v7/modules/core/keeper"
	ibctmmigrations "github.com/cosmos/ibc-go/v7/modules/light-clients/07-tendermint/migrations"
	cellarfeeskeeper "github.com/peggyjv/sommelier/v8/x/cellarfees/keeper"
	cellarfeeskeeperv1 "github.com/peggyjv/sommelier/v8/x/cellarfees/migrations/v1/keeper"
	cellarfeestypesv2 "github.com/peggyjv/sommelier/v8/x/cellarfees/types/v2"
)

func CreateUpgradeHandler(
	mm *module.Manager,
	configurator module.Configurator,
	baseAppLegacySS *paramstypes.Subspace,
	consensusParamsKeeper *consensusparamkeeper.Keeper,
	cellarfeesLegacyKeeper *cellarfeeskeeperv1.Keeper,
	cellarfeesKeeper *cellarfeeskeeper.Keeper,
	ibcKeeper *ibckeeper.Keeper,
	cdc codec.BinaryCodec,
	clientKeeper ibctmmigrations.ClientKeeper,
) upgradetypes.UpgradeHandler {
	return func(ctx sdk.Context, plan upgradetypes.Plan, vm module.VersionMap) (module.VersionMap, error) {
		ctx.Logger().Info("v8 upgrade: entering handler and running migrations")

		// Include this when migrating to ibc-go v7 (optional)
		// source: https://github.com/cosmos/ibc-go/blob/v7.2.0/docs/migrations/v6-to-v7.md
		// prune expired tendermint consensus states to save storage space
		if _, err := ibctmmigrations.PruneExpiredConsensusStates(ctx, cdc, clientKeeper); err != nil {
			return nil, err
		}

		// new x/consensus module params migration
		baseapp.MigrateParams(ctx, baseAppLegacySS, consensusParamsKeeper)

		// cellarfees params migration
		newParams := &cellarfeestypesv2.Params{}
		oldParams := cellarfeesLegacyKeeper.GetParams(ctx)
		newParams.AuctionThresholdUsdValue = cellarfeestypesv2.DefaultParams().AuctionThresholdUsdValue
		newParams.AuctionInterval = oldParams.AuctionInterval
		newParams.InitialPriceDecreaseRate = oldParams.InitialPriceDecreaseRate
		newParams.PriceDecreaseBlockInterval = oldParams.PriceDecreaseBlockInterval
		newParams.RewardEmissionPeriod = oldParams.RewardEmissionPeriod

		err := newParams.ValidateBasic()
		if err != nil {
			return nil, err
		}

		cellarfeesKeeper.SetParams(ctx, *newParams)

		// explicitly update the IBC 02-client params, adding the localhost client type
		params := ibcKeeper.ClientKeeper.GetParams(ctx)
		params.AllowedClients = append(params.AllowedClients, ibcexported.Localhost)
		ibcKeeper.ClientKeeper.SetParams(ctx, params)

		return mm.RunMigrations(ctx, configurator, vm)
	}
}
