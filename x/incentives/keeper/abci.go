package keeper

import (
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	distributiontypes "github.com/cosmos/cosmos-sdk/x/distribution/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k Keeper) BeginBlocker(ctx sdk.Context) {}

// EndBlocker defines Distribution of incentives to stakers
//
// 1) Get the amount of coins to distribute
//
// 2) Send the coins to the distribution module

func (k Keeper) EndBlocker(ctx sdk.Context) {
	incentivesParams := k.GetParamSet(ctx)
	// check if incentives are enabled
	if uint64(ctx.BlockHeight()) > incentivesParams.IncentivesCutoffHeight || incentivesParams.DistributionPerBlock.IsZero() {

		return
	}

	distPerBlockCoins := sdk.NewCoins(incentivesParams.DistributionPerBlock)
	feePool := k.DistributionKeeper.GetFeePool(ctx)
	newPool, negative := feePool.CommunityPool.SafeSub(sdk.NewDecCoinsFromCoins(distPerBlockCoins...))
	if negative {
		k.Logger(ctx).Error("Insufficient coins in community to distribute", "community pool", feePool.CommunityPool)
		return
	}

	// Send to fee collector for distribution
	err := k.BankKeeper.SendCoinsFromModuleToModule(ctx, distributiontypes.ModuleName, authtypes.FeeCollectorName, distPerBlockCoins)
	if err != nil {
		panic(err)
	}

	feePool.CommunityPool = newPool
	k.DistributionKeeper.SetFeePool(ctx, feePool)
}
