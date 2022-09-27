package v1

import (
	"bytes"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/ethereum/go-ethereum/common"

	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/peggyjv/sommelier/v4/x/cork/migrations/v1/types"
)

func MigrateStore(ctx sdk.Context, storeKey storetypes.StoreKey, cdc codec.BinaryCodec) error {
	ctx.Logger().Info("Cork v1 to v2: Beginning store migration")

	store := ctx.KVStore(storeKey)

	removeCommitPeriod(ctx, store)
	removeOldCorks(store, cdc)

	ctx.Logger().Info("Cork v1 to v2: Store migration complete")

	return nil
}

func removeCommitPeriod(ctx sdk.Context, store storetypes.KVStore) {
	store.Delete([]byte{types.CommitPeriodStartKey})

	// TODO(bolten): remove commit period param from old state
}

func removeOldCorks(store storetypes.KVStore, cdc codec.BinaryCodec) {
	var validatorCorks []*types.ValidatorCork
	iter := sdk.KVStorePrefixIterator(store, []byte{types.CorkForAddressKeyPrefix})
	defer iter.Close()
	for ; iter.Valid(); iter.Next() {
		keyBytes := bytes.NewBuffer(bytes.TrimPrefix(iter.Key(), []byte{types.CorkForAddressKeyPrefix}))
		val := sdk.ValAddress(keyBytes.Next(20))

		var cork types.Cork
		cdc.MustUnmarshal(iter.Value(), &cork)
		validatorCorks = append(validatorCorks, &types.ValidatorCork{
			Cork:      &cork,
			Validator: val.String(),
		})
	}

	for _, validatorCork := range validatorCorks {
		store.Delete(types.GetCorkForValidatorAddressKey(
			sdk.ValAddress(validatorCork.Validator),
			common.HexToAddress(validatorCork.Cork.TargetContractAddress),
		))
	}
}
