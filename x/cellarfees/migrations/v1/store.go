package v1

import (
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	paramstypes "github.com/cosmos/cosmos-sdk/x/params/types"
	v1types "github.com/peggyjv/sommelier/v7/x/cellarfees/migrations/v1/types"
	"github.com/peggyjv/sommelier/v7/x/cellarfees/types"
	v2types "github.com/peggyjv/sommelier/v7/x/cellarfees/types/v2"
)

func MigrateStore(ctx sdk.Context, storeKey storetypes.StoreKey, cdc codec.BinaryCodec, paramSubspace paramstypes.Subspace) error {
	ctx.Logger().Info("cellarfees v1 to v2: Beginning store migration")

	store := ctx.KVStore(storeKey)
	err := migrateParamStore(ctx, paramSubspace)
	if err != nil {
		return err
	}

	migrateCellarfeesFeeAccrualCounters(ctx, store)

	ctx.Logger().Info("cellarfees v1 to v2: Store migration complete")

	return nil
}

func migrateParamStore(ctx sdk.Context, subspace paramstypes.Subspace) error {
	ctx.Logger().Info("cellarfees v1 to v2: Migrating params")

	oldParamSet := &v1types.Params{}
	subspace.GetParamSet(ctx, oldParamSet)

	newParamSet := &v2types.Params{}
	newParamSet.AuctionThresholdUsdValue = v2types.DefaultParams().AuctionThresholdUsdValue
	newParamSet.AuctionInterval = oldParamSet.AuctionInterval
	newParamSet.InitialPriceDecreaseRate = oldParamSet.InitialPriceDecreaseRate
	newParamSet.PriceDecreaseBlockInterval = oldParamSet.PriceDecreaseBlockInterval
	newParamSet.RewardEmissionPeriod = oldParamSet.RewardEmissionPeriod

	err := newParamSet.ValidateBasic()
	if err != nil {
		return err
	}

	// (Collin): Does this actually delete the old FeeAccrualCounters key? Does it matter?
	subspace.SetParamSet(ctx, newParamSet)

	ctx.Logger().Info("cellarfees v1 to v2: Params migration complete")
	return nil
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
