package keeper

import (
	"fmt"

	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	distributiontypes "github.com/cosmos/cosmos-sdk/x/distribution/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k Keeper) BeginBlocker(ctx sdk.Context) {}

// EndBlocker defines distrbution of incentives to stakers
//
// 1) Get the amount of coins to distribute
//
// 2) Send the coinst to the distribution module

func (k Keeper) EndBlocker(ctx sdk.Context) {
	params := k.GetParamSet(ctx)

	k.Logger(ctx).Info("tallying scheduled cork votes", "height", fmt.Sprintf("%d", ctx.BlockHeight()))

	feePool := k.DistributionKeeper.GetFeePool(ctx)

	newPool, negative := feePool.CommunityPool.SafeSub(sdk.NewDecCoinsFromCoins(params.DistrbutionPerBlock...))
	if negative {
		k.Logger(ctx).Error("Insufficient coins in community to distribute", "community pool", feePool.CommunityPool)
		return
	}

	feePool.CommunityPool = newPool

	// Send to fee collector for distribution
	err := k.bankKeeper.SendCoinsFromModuleToModule(ctx, distributiontypes.ModuleName, authtypes.FeeCollectorName, params.DistrbutionPerBlock)
	if err != nil {
		panic(err)
	}

}
