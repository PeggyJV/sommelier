package v8

import (
	"bytes"

	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	consensusparamkeeper "github.com/cosmos/cosmos-sdk/x/consensus/keeper"
	paramstypes "github.com/cosmos/cosmos-sdk/x/params/types"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"
	ibcexported "github.com/cosmos/ibc-go/v7/modules/core/exported"
	ibckeeper "github.com/cosmos/ibc-go/v7/modules/core/keeper"
	ibctmmigrations "github.com/cosmos/ibc-go/v7/modules/light-clients/07-tendermint/migrations"
	cellarfeeskeeper "github.com/peggyjv/sommelier/v8/x/cellarfees/keeper"
	cellarfeestypesv1 "github.com/peggyjv/sommelier/v8/x/cellarfees/migrations/v1/types"
	cellarfeestypes "github.com/peggyjv/sommelier/v8/x/cellarfees/types/v2"
)

func CreateUpgradeHandler(
	mm *module.Manager,
	configurator module.Configurator,
	baseAppLegacySS *paramstypes.Subspace,
	consensusParamsKeeper *consensusparamkeeper.Keeper,
	cellarfeesLegacySS *paramstypes.Subspace,
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

		// delete the old param key
		store := prefix.NewStore(ctx.KVStore(sdk.NewKVStoreKey("params")), append([]byte("params"), '/'))
		store.Delete(bytes.Join([][]byte{[]byte("cellarfees"), {byte('/')}, cellarfeestypesv1.KeyFeeAccrualAuctionThreshold}, []byte{}))

		oldParams := cellarfeestypesv1.Params{}
		cellarfeesLegacySS.GetParamSet(ctx, &oldParams)
		newParams := cellarfeestypes.Params{
			AuctionInterval:            oldParams.AuctionInterval,
			RewardEmissionPeriod:       oldParams.RewardEmissionPeriod,
			PriceDecreaseBlockInterval: oldParams.PriceDecreaseBlockInterval,
			InitialPriceDecreaseRate:   oldParams.InitialPriceDecreaseRate,
			AuctionThresholdUsdValue:   cellarfeestypes.DefaultParams().AuctionThresholdUsdValue,
		}
		cellarfeesKeeper.SetParams(ctx, newParams)

		// explicitly update the IBC 02-client params, adding the localhost client type
		params := ibcKeeper.ClientKeeper.GetParams(ctx)
		params.AllowedClients = append(params.AllowedClients, ibcexported.Localhost)
		ibcKeeper.ClientKeeper.SetParams(ctx, params)

		return mm.RunMigrations(ctx, configurator, vm)
	}
}
