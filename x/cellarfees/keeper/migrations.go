package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
	v1 "github.com/peggyjv/sommelier/v9/x/cellarfees/migrations/v1"
	v2 "github.com/peggyjv/sommelier/v9/x/cellarfees/types/v2"
)

// Migrator is a struct for handling in-place store migrations.
type Migrator struct {
	keeper         Keeper
	legacySubspace paramtypes.Subspace
}

// NewMigrator returns a new Migrator.
func NewMigrator(keeper Keeper, legacySubspace paramtypes.Subspace) Migrator {
	return Migrator{keeper: keeper, legacySubspace: legacySubspace}
}

// Migrate1to2 migrates from consensus version 1 to 2.
func (m Migrator) Migrate1to2(ctx sdk.Context) error {
	currentParams := m.keeper.GetParamSetIfExists(ctx)
	params := v2.Params{
		AuctionInterval:            currentParams.AuctionInterval,
		InitialPriceDecreaseRate:   currentParams.InitialPriceDecreaseRate,
		PriceDecreaseBlockInterval: currentParams.PriceDecreaseBlockInterval,
		RewardEmissionPeriod:       currentParams.RewardEmissionPeriod,
		AuctionThresholdUsdValue:   sdk.MustNewDecFromStr(v2.DefaultAuctionThresholdUsdValue),
	}

	if err := params.ValidateBasic(); err != nil {
		return err
	}

	m.keeper.SetParams(ctx, params)

	v1.MigrateStore(ctx, m.keeper.storeKey, m.keeper.cdc, m.legacySubspace)

	return nil
}

// Migrate2to3 migrates from consensus version 2 to 3.
func (m Migrator) Migrate2to3(ctx sdk.Context) error {
	ctx.Logger().Info("cellarfees v2 to v3: New param")
	subspace := m.keeper.paramSpace

	if !subspace.Has(ctx, v2.KeyProceedsPortion) {
		subspace.Set(ctx, v2.KeyProceedsPortion, v2.DefaultParams().ProceedsPortion)
	}

	ctx.Logger().Info("cellarfees v2 to v3: Params migration complete")
	return nil
}
