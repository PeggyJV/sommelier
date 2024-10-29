package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	paramstypes "github.com/cosmos/cosmos-sdk/x/params/types"
	v1 "github.com/peggyjv/sommelier/v8/x/cellarfees/migrations/v1"
	v1types "github.com/peggyjv/sommelier/v8/x/cellarfees/migrations/v1/types"
	v2types "github.com/peggyjv/sommelier/v8/x/cellarfees/types/v2"
)

// Migrator is a struct for handling in-place store migrations.
type Migrator struct {
	keeper Keeper
}

// NewMigrator returns a new Migrator.
func NewMigrator(keeper Keeper) Migrator {
	return Migrator{keeper: keeper}
}

// Migrate1to2 migrates from consensus version 1 to 2.
func (m Migrator) Migrate1to2(ctx sdk.Context) error {
	// Migrate params
	newParams, err := getNewParams(ctx, m.keeper.paramSpace)
	if err != nil {
		return err
	}

	m.keeper.SetParams(ctx, *newParams)

	ctx.Logger().Info("cellarfees v1 to v2: Params migration complete")

	v1.MigrateStore(ctx, m.keeper.storeKey, m.keeper.cdc)

	return nil
}

func getNewParams(ctx sdk.Context, subspace paramstypes.Subspace) (*v2types.Params, error) {
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
		return nil, err
	}

	return newParamSet, nil
}
