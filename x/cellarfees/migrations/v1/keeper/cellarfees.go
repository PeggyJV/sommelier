package keeper

import (
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/peggyjv/sommelier/v9/app/params"
	"github.com/peggyjv/sommelier/v9/x/cellarfees/types"
)

// Getter for module account that holds the fee pool funds
func (k Keeper) GetFeesAccount(ctx sdk.Context) authtypes.ModuleAccountI {
	return k.accountKeeper.GetModuleAccount(ctx, types.ModuleName)
}

func (k Keeper) GetEmission(ctx sdk.Context, remainingRewardsSupply math.Int) sdk.Coins {
	previousSupplyPeak := k.GetLastRewardSupplyPeak(ctx)
	cellarfeesParams := k.GetParams(ctx)

	var emissionAmount math.Int
	if remainingRewardsSupply.GT(previousSupplyPeak) {
		k.SetLastRewardSupplyPeak(ctx, remainingRewardsSupply)
		emissionAmount = remainingRewardsSupply.Quo(sdk.NewInt(int64(cellarfeesParams.RewardEmissionPeriod)))
	} else {
		emissionAmount = previousSupplyPeak.Quo(sdk.NewInt(int64(cellarfeesParams.RewardEmissionPeriod)))
	}

	// Emission should be at least 1usomm and at most the remaining reward supply
	if emissionAmount.IsZero() {
		emissionAmount = sdk.OneInt()
	} else if emissionAmount.GTE(remainingRewardsSupply) {
		// We zero out the previous peak value here to avoid doing it every block. We set the final emission
		// to the remaining supply here even though it's potentially redundant because it's less code than
		// having another check where we would also have to zero out the prevoius peak supply.
		k.SetLastRewardSupplyPeak(ctx, sdk.ZeroInt())
		emissionAmount = remainingRewardsSupply
	}

	return sdk.NewCoins(sdk.NewCoin(params.BaseCoinUnit, emissionAmount))
}
