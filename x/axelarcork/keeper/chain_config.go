package keeper

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/peggyjv/sommelier/v6/x/axelarcork/types"
)

func (k Keeper) GetChainConfigurationByNameAndID(ctx sdk.Context, chainName string, chainID uint64) (*types.ChainConfiguration, error) {
	var config types.ChainConfiguration
	if chainID != 0 {
		idConfig, ok := k.GetChainConfigurationByID(ctx, chainID)
		if !ok {
			return nil, fmt.Errorf("chain configuration for ID %d does not exist", chainID)
		}
		if chainName != "" {
			if chainName != idConfig.Name {
				return nil, fmt.Errorf("chain ID %d and name %s do not match", chainID, chainName)
			}
		}
		config = idConfig
	} else if chainName != "" {
		nameConfig, ok := k.GetChainConfigurationByName(ctx, chainName)
		if !ok {
			return nil, fmt.Errorf("chain configuration for name %s does not exist", chainName)
		}
		if chainID != 0 {
			if nameConfig.Id != chainID {
				return nil, fmt.Errorf("chain ID %d and name %s do not match", chainID, chainName)
			}
		}
		config = nameConfig
	} else {
		return nil, fmt.Errorf("neither chain name nor ID were submitted with this proposal")
	}

	return &config, nil
}

func (k Keeper) GetChainConfigurationByID(ctx sdk.Context, chainID uint64) (types.ChainConfiguration, bool) {
	bz := ctx.KVStore(k.storeKey).Get(types.ChainConfigurationKey(chainID))
	if len(bz) == 0 {
		return types.ChainConfiguration{}, false
	}

	var chainConfig types.ChainConfiguration
	k.cdc.MustUnmarshal(bz, &chainConfig)
	return chainConfig, true
}

func (k Keeper) SetChainConfigurationByID(ctx sdk.Context, chainID uint64, config types.ChainConfiguration) {
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
