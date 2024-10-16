package keeper

import (
	abci "github.com/cometbft/cometbft/abci/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	distributiontypes "github.com/cosmos/cosmos-sdk/x/distribution/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// BeginBlocker defines distribution rewards for validators
//
// 1) Subtract the total distribution from the community pool
// 2) Get a list of qualifying validators sorted by descending power
// 3) Allocate tokens to qualifying validators proportionally to their power with a cap
// 4) Add the remaining coins back to the community pool
func (k Keeper) BeginBlocker(ctx sdk.Context, req abci.RequestBeginBlock) {
	incentivesParams := k.GetParamSet(ctx)
	if uint64(ctx.BlockHeight()) >= incentivesParams.ValidatorIncentivesCutoffHeight || incentivesParams.ValidatorDistributionPerBlock.IsZero() {
		return
	}

	// Rewards come from the community pool
	totalDistribution := sdk.NewDecCoinsFromCoins(incentivesParams.ValidatorDistributionPerBlock)
	feePool := k.DistributionKeeper.GetFeePool(ctx)
	newPool, negative := feePool.CommunityPool.SafeSub(totalDistribution)
	if negative {
		k.Logger(ctx).Error("Insufficient coins in community to distribute to validators", "community pool", feePool.CommunityPool)
		return
	}

	// Get a list of qualifying validators sorted by descending power
	valInfos := k.getValidatorInfos(ctx, req)
	sortedValInfos := sortValidatorInfosByPower(valInfos)
	qualifyingVoters := sortedValInfos[:incentivesParams.ValidatorIncentivesSetSizeLimit]

	// Allocate tokens to qualifying validators proportionally to their power with a cap
	totalPower := getTotalPower(&qualifyingVoters)
	remaining := k.AllocateTokens(ctx, totalPower, totalDistribution, qualifyingVoters, incentivesParams.ValidatorIncentivesMaxFraction)

	// Add the remaining coins back to the community pool
	newPool = newPool.Add(remaining...)
	feePool.CommunityPool = newPool
	k.DistributionKeeper.SetFeePool(ctx, feePool)
}

// EndBlocker defines Distribution of incentives to stakers
//
// 1) Get the amount of coins to distribute
//
// 2) Send the coins to the distribution module

func (k Keeper) EndBlocker(ctx sdk.Context) {
	incentivesParams := k.GetParamSet(ctx)
	// check if incentives are enabled
	if uint64(ctx.BlockHeight()) >= incentivesParams.IncentivesCutoffHeight || incentivesParams.DistributionPerBlock.IsZero() {
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
