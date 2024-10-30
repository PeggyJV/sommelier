package v1

import (
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
	"github.com/peggyjv/sommelier/v8/x/cellarfees/types"
)

func MigrateStore(ctx sdk.Context, storeKey storetypes.StoreKey, cdc codec.BinaryCodec, legacySubspace paramtypes.Subspace) {
	ctx.Logger().Info("cellarfees v1 to v2: Beginning store migration")

	store := ctx.KVStore(storeKey)

	migrateCellarfeesFeeAccrualCounters(ctx, store)

	ctx.Logger().Info("cellarfees v1 to v2: Store migration complete")
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
