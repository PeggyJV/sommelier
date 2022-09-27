package v1

import (
	"github.com/cosmos/cosmos-sdk/codec"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/peggyjv/sommelier/v4/x/cork/migrations/v1/types"
)

func MigrateStore(ctx sdk.Context, storeKey storetypes.StoreKey, cdc codec.BinaryCodec) error {
	ctx.Logger().Info("Cork v1 to v2: Beginning store migration")

	store := ctx.KVStore(storeKey)

	removeCommitPeriod(store)

	ctx.Logger().Info("Cork v1 to v2: Store migration complete")

	return nil
}

func removeCommitPeriod(store storetypes.KVStore) {
	store.Delete([]byte{types.CommitPeriodStartKey})
}
