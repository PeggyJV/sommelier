package v1

import (
	"bytes"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
	v1types "github.com/peggyjv/sommelier/v8/x/cellarfees/migrations/v1/types"
	"github.com/peggyjv/sommelier/v8/x/cellarfees/types"
)

func MigrateStore(ctx sdk.Context, storeKey storetypes.StoreKey, cdc codec.BinaryCodec, legacySubspace paramtypes.Subspace) {
	ctx.Logger().Info("cellarfees v1 to v2: Beginning store migration")

	store := ctx.KVStore(storeKey)

	deleteFeeAccrualAuctionThresholdParam(ctx, legacySubspace)
	migrateCellarfeesFeeAccrualCounters(ctx, store)

	ctx.Logger().Info("cellarfees v1 to v2: Store migration complete")
}

func deleteFeeAccrualAuctionThresholdParam(ctx sdk.Context, legacySubspace paramtypes.Subspace) {
	ctx.Logger().Info("cellarfees v1 to v2: Deleting FeeAccrualAuctionThreshold param")
	store := ctx.KVStore(sdk.NewKVStoreKey(paramtypes.ModuleName))
	paramStore := prefix.NewStore(store, append([]byte(types.ModuleName), '/'))

	legacySubspace.IterateKeys(ctx, func(key []byte) bool {
		if bytes.Equal(key, v1types.KeyFeeAccrualAuctionThreshold) {
			paramStore.Delete(key)
			return true
		}
		return false
	})

	ctx.Logger().Info("cellarfees v1 to v2: FeeAccrualAuctionThreshold param deleted!")
}

// Removes the old FeeAccrualCounters state
func migrateCellarfeesFeeAccrualCounters(ctx sdk.Context, store sdk.KVStore) {
	ctx.Logger().Info("cellarfees v1 to v2: Migrating FeeAccrualCounters")

	prefixStore := prefix.NewStore(store, []byte{types.FeeAccrualCountersKey})
	iter := prefixStore.Iterator(nil, nil)
	defer iter.Close()

	for ; iter.Valid(); iter.Next() {
		prefixStore.Delete(iter.Key())
	}

	ctx.Logger().Info("cellarfees v1 to v2: FeeAccrualCounters migration complete")
}
