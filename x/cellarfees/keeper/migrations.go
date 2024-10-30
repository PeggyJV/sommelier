package keeper

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
	v1 "github.com/peggyjv/sommelier/v8/x/cellarfees/migrations/v1"
	v2 "github.com/peggyjv/sommelier/v8/x/cellarfees/types/v2"
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

	panic(fmt.Sprintf("cellarfees: it worked %v", params))

	v1.MigrateStore(ctx, m.keeper.storeKey, m.keeper.cdc, m.legacySubspace)

	return nil
}
