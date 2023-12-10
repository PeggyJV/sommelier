package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/peggyjv/sommelier/v7/x/axelarcork/types"
)

func (k Keeper) GetChainConfigurationByID(ctx sdk.Context, chainID uint64) (types.ChainConfiguration, bool) {
	bz := ctx.KVStore(k.storeKey).Get(types.ChainConfigurationKey(chainID))
	if len(bz) == 0 {
		return types.ChainConfiguration{}, false
	}

	var chainConfig types.ChainConfiguration
	k.cdc.MustUnmarshal(bz, &chainConfig)
	return chainConfig, true
}

func (k Keeper) SetChainConfiguration(ctx sdk.Context, chainID uint64, config types.ChainConfiguration) {
	bz := k.cdc.MustMarshal(&config)
	ctx.KVStore(k.storeKey).Set(types.ChainConfigurationKey(chainID), bz)
}

func (k Keeper) IterateChainConfigurations(ctx sdk.Context, handler func(config types.ChainConfiguration) (stop bool)) {
	store := ctx.KVStore(k.storeKey)
	iter := sdk.KVStorePrefixIterator(store, []byte{types.ChainConfigurationPrefix})
	defer iter.Close()
	for ; iter.Valid(); iter.Next() {
		var chainConfig types.ChainConfiguration
		k.cdc.MustUnmarshal(iter.Value(), &chainConfig)
		if handler(chainConfig) {
			break
		}
	}
}

func (k Keeper) GetChainConfigurationByName(ctx sdk.Context, chainName string) (types.ChainConfiguration, bool) {
	var chainConfig types.ChainConfiguration
	found := false
	k.IterateChainConfigurations(ctx, func(config types.ChainConfiguration) (stop bool) {
		if config.Name == chainName {
			chainConfig = config
			found = true
			return true
		}

		return false
	})

	return chainConfig, found
}

func (k Keeper) DeleteChainConfigurationByID(ctx sdk.Context, chainID uint64) {
	ctx.KVStore(k.storeKey).Delete(types.ChainConfigurationKey(chainID))
}
