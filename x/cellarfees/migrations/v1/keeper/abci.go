package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/peggyjv/sommelier/v9/app/params"
)

// BeginBlocker emits rewards each block they are available by sending them to the distribution module's fee collector
// account. Emissions are a constant value based on the last peak supply of distributable fees so that the reward supply
// will decrease linearly until exhausted.
func (k Keeper) BeginBlocker(ctx sdk.Context) {

	// Handle fee auctions
	cellarfeesParams := k.GetParams(ctx)

	counters := k.GetFeeAccrualCounters(ctx)

	modulus := ctx.BlockHeader().Height % int64(cellarfeesParams.AuctionInterval)

	for _, counter := range counters.Counters {

		if counter.Count >= cellarfeesParams.FeeAccrualAuctionThreshold && modulus == 0 {
			started := k.beginAuction(ctx, counter.Denom)
			if started {
				counters.ResetCounter(counter.Denom)
			}
		}

	}
	k.SetFeeAccrualCounters(ctx, counters)

	// Handle reward emissions
	moduleAccount := k.GetFeesAccount(ctx)
	remainingRewardsSupply := k.bankKeeper.GetBalance(ctx, moduleAccount.GetAddress(), params.BaseCoinUnit).Amount

	if remainingRewardsSupply.IsZero() {
		return
	}

	emission := k.GetEmission(ctx, remainingRewardsSupply)

	// Send to fee collector for distribution
	err := k.bankKeeper.SendCoinsFromModuleToModule(ctx, moduleAccount.GetName(), authtypes.FeeCollectorName, emission)
	if err != nil {
		panic(err)
	}

}

// EndBlocker is called at the end of every block
func (k Keeper) EndBlocker(ctx sdk.Context) {}
