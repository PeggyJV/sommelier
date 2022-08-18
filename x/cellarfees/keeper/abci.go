package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/peggyjv/sommelier/v4/app/params"
	"github.com/peggyjv/sommelier/v4/x/cellarfees/types"
)

// approximate number of blocks per two weeks
const emissionPeriod = 172800

// BeginBlocker calculates a reward emission based on a constant proportion of the latest peak reward supply.
// This results in a constant emission rate between top-ups that will exhaust the reward supply after a number
// of blocks equal to `emissionPeriod`
func (k Keeper) BeginBlocker(ctx sdk.Context) {
	moduleAccount := k.accountKeeper.GetModuleAccount(ctx, types.ModuleName)
	remainingRewardsSupply := k.bankKeeper.GetBalance(ctx, moduleAccount.GetAddress(), params.BaseCoinUnit).Amount

	if remainingRewardsSupply.IsZero() {
		return
	}

	previousSupplyPeak := k.GetLastRewardSupplyPeak(ctx)

	var emissionAmount sdk.Int

	// If this is the first emission after a resupply from zero, we can't calculate the emission from
	// the previous peak because it's zeroed out in the store, so we use the current reward supply.
	if previousSupplyPeak.IsZero() {
		k.SetLastRewardSupplyPeak(ctx, remainingRewardsSupply)
		emissionAmount = remainingRewardsSupply.Quo(sdk.NewInt(emissionPeriod))
	} else {
		emissionAmount = previousSupplyPeak.Quo(sdk.NewInt(emissionPeriod))
	}

	// Emission should be at least 1usomm and at most the remaining supply
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
