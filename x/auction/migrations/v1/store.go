package v1

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	paramstypes "github.com/cosmos/cosmos-sdk/x/params/types"
	"github.com/peggyjv/sommelier/v7/x/auction/types"
)

func MigrateParamStore(ctx sdk.Context, subspace paramstypes.Subspace) error {
	ctx.Logger().Info("auction v1 to v2: Migrating params")

	if !subspace.Has(ctx, types.KeyAuctionBurnRate) {
		subspace.Set(ctx, types.KeyAuctionBurnRate, types.DefaultParams().AuctionBurnRate)
	}

	ctx.Logger().Info("auction v1 to v2: Params migration complete")
	return nil
}
