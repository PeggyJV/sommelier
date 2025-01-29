package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/peggyjv/sommelier/v9/x/auction/types"
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
	ctx.Logger().Info("auction v1 to v2: New params")
	subspace := m.keeper.paramSpace

	if !subspace.Has(ctx, types.KeyAuctionBurnRate) {
		subspace.Set(ctx, types.KeyAuctionBurnRate, types.DefaultParams().AuctionBurnRate)
	}

	ctx.Logger().Info("auction v1 to v2: Params migration complete")
	return nil
}
