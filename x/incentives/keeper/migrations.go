package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/peggyjv/sommelier/v8/x/incentives/types"
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
	ctx.Logger().Info("incentives v1 to v2: New params")

	// New params
	subspace := m.keeper.paramSpace

	if !subspace.Has(ctx, types.KeyValidatorIncentivesCutoffHeight) {
		subspace.Set(ctx, types.KeyValidatorIncentivesCutoffHeight, types.DefaultParams().ValidatorIncentivesCutoffHeight)
	}

	if !subspace.Has(ctx, types.KeyValidatorMaxDistributionPerBlock) {
		subspace.Set(ctx, types.KeyValidatorMaxDistributionPerBlock, types.DefaultParams().ValidatorMaxDistributionPerBlock)
	}

	if !subspace.Has(ctx, types.KeyValidatorIncentivesMaxFraction) {
		subspace.Set(ctx, types.KeyValidatorIncentivesMaxFraction, types.DefaultParams().ValidatorIncentivesMaxFraction)
	}

	if !subspace.Has(ctx, types.KeyValidatorIncentivesSetSizeLimit) {
		subspace.Set(ctx, types.KeyValidatorIncentivesSetSizeLimit, types.DefaultParams().ValidatorIncentivesSetSizeLimit)
	}

	ctx.Logger().Info("incentives v1 to v2: Params migration complete")

	return nil
}
