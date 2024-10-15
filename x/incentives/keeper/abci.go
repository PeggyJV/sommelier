package keeper

import (
	abci "github.com/cometbft/cometbft/abci/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	distributiontypes "github.com/cosmos/cosmos-sdk/x/distribution/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k Keeper) BeginBlocker(ctx sdk.Context, req abci.RequestBeginBlock) {
	incentivesParams := k.GetParamSet(ctx)
	cutoffHeight := incentivesParams.ValidatorIncentivesCutoffHeight
	distPerBlock := incentivesParams.ValidatorDistributionPerBlock
	if uint64(ctx.BlockHeight()) >= cutoffHeight || distPerBlock.IsZero() {
		return
	}

	voterInfos := GetSortedVoterInfosByPower(req.LastCommitInfo.GetVotes())
	totalPower := int64(0)
	for _, voterInfo := range voterInfos {
		totalPower += voterInfo.Validator.Power
	}

	// Limit the number of voter info
	// TODO: Make this a function that can be unit tested
	setSizeLimit := incentivesParams.ValidatorIncentivesSetSizeLimit
	if uint64(len(voterInfos)) > setSizeLimit {
		voterInfos = voterInfos[:setSizeLimit]
	}

	distPerBlockDec := sdk.NewDec(distPerBlock.Amount.Int64())
	validatorMaxPortionFraction := sdk.MustNewDecFromStr("0.1")
	apportionments, remaining, err := getApportionments(setSizeLimit, distPerBlockDec, validatorMaxPortionFraction)
	if err != nil {
		ctx.Logger().Error("Error getting apportionments. Are params properly validated?", "error", err)
		return
	}

	// Distribute rewards to each validator
	feePool := k.DistributionKeeper.GetFeePool(ctx)
	newPool, negative := feePool.CommunityPool.SafeSub(sdk.NewDecCoinsFromCoins(distPerBlock))
	if negative {
		k.Logger(ctx).Error("Insufficient coins in community to distribute", "community pool", feePool.CommunityPool)
		return
	}

	for i, voterInfo := range voterInfos {
		recipient := sdk.AccAddress(voterInfo.Validator.Address)
		amount := apportionments[i].TruncateInt()
		if amount.IsZero() {
			continue
		}

		err := k.BankKeeper.SendCoinsFromModuleToAccount(ctx, distributiontypes.ModuleName, recipient, sdk.NewCoins(sdk.NewCoin(distPerBlock.Denom, amount)))
		if err != nil {
			panic(err)
		}
	}

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
