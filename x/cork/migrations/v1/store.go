package v1

import (
	"bytes"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/ethereum/go-ethereum/common"

	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/peggyjv/sommelier/v7/x/cork/migrations/v1/types"
)

func MigrateStore(ctx sdk.Context, storeKey storetypes.StoreKey, cdc codec.BinaryCodec) error {
	ctx.Logger().Info("Cork v1 to v2: Beginning store migration")

	store := ctx.KVStore(storeKey)

	removeCommitPeriod(store)
	removeOldCorks(store, cdc)

	ctx.Logger().Info("Cork v1 to v2: Store migration complete")

	return nil
}

func removeCommitPeriod(store storetypes.KVStore) {
	store.Delete([]byte{types.CommitPeriodStartKey})
}

// wipe away all existing cork state during upgrade -- there shouldn't be in-transit corks
// during the upgrade
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

	var scheduledCorks []*types.ScheduledCork
	iter = sdk.KVStorePrefixIterator(store, []byte{types.ScheduledCorkKeyPrefix})
	defer iter.Close()
	for ; iter.Valid(); iter.Next() {
		var cork types.Cork
		keyPair := bytes.NewBuffer(iter.Key())
		keyPair.Next(1) // trim prefix byte
		blockHeight := sdk.BigEndianToUint64(keyPair.Next(8))
		val := sdk.ValAddress(keyPair.Next(20))

		cdc.MustUnmarshal(iter.Value(), &cork)
		scheduledCorks = append(scheduledCorks, &types.ScheduledCork{
			Cork:        &cork,
			BlockHeight: blockHeight,
			Validator:   val.String(),
		})
	}

	for _, scheduledCork := range scheduledCorks {
		store.Delete(types.GetScheduledCorkKey(
			scheduledCork.BlockHeight,
			sdk.ValAddress(scheduledCork.Validator),
			common.HexToAddress(scheduledCork.Cork.TargetContractAddress),
		))
	}
}
