package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/peggyjv/sommelier/v4/app/params"
	"github.com/peggyjv/sommelier/v4/x/cellarfees/types"
)

// BeginBlocker calculates a reward emission based on a constant proportion of the latest peak reward supply.
// This results in a constant emission rate between top-ups that will exhaust the reward supply after a number
// of blocks equal to the RewardEmissionPeriod param.
func (k Keeper) BeginBlocker(ctx sdk.Context) {
	moduleAccount := k.accountKeeper.GetModuleAccount(ctx, types.ModuleName)
	remainingRewardsSupply := k.bankKeeper.GetBalance(ctx, moduleAccount.GetAddress(), params.BaseCoinUnit).Amount

	if remainingRewardsSupply.IsZero() {
		return
	}

	previousSupplyPeak := k.GetLastRewardSupplyPeak(ctx)
	p := k.GetParams(ctx)

	var emissionAmount sdk.Int

	// If this is the first emission after a resupply from zero, the current reward supply is the new peak.
	if previousSupplyPeak.IsZero() {
		k.SetLastRewardSupplyPeak(ctx, remainingRewardsSupply)
		emissionAmount = remainingRewardsSupply.Quo(sdk.NewInt(int64(p.RewardEmissionPeriod)))
	} else {
		emissionAmount = previousSupplyPeak.Quo(sdk.NewInt(int64(p.RewardEmissionPeriod)))
	}

	// Emission should be at least 1usomm and at most the remaining reward supply
	if emissionAmount.IsZero() {
		emissionAmount = sdk.OneInt()
	} else if emissionAmount.GTE(remainingRewardsSupply) {
		// We zero out the previous peak value early to avoid doing it every block when the remaining supply's
		// zero check occurs. We set the final emission to the remaining supply here even though it's potentially
		// redundant because it's less code than having another check where we would also have to zero out the
		// prevoius peak supply.
		k.SetLastRewardSupplyPeak(ctx, sdk.ZeroInt())
		emissionAmount = remainingRewardsSupply
	}

	coin := sdk.NewCoin(params.BaseCoinUnit, emissionAmount)
	emission := sdk.NewCoins(coin)

	// Send to fee collector for distribution
	distributionFeeCollectorAddress := k.accountKeeper.GetModuleAccount(ctx, authtypes.FeeCollectorName).GetAddress().String()
	k.bankKeeper.SendCoinsFromModuleToModule(ctx, moduleAccount.GetAddress().String(), distributionFeeCollectorAddress, emission)
}

// EndBlocker is called at the end of every block
func (k Keeper) EndBlocker(ctx sdk.Context) {}
