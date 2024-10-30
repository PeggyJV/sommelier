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
	currentParams := m.keeper.GetParamSetIfExists(ctx)
	params := types.DefaultParams()
	params.IncentivesCutoffHeight = currentParams.IncentivesCutoffHeight
	params.DistributionPerBlock = currentParams.DistributionPerBlock

	if err := params.ValidateBasic(); err != nil {
		return err
	}

	m.keeper.SetParams(ctx, params)

	return nil
}
