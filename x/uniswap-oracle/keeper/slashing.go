package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// SlashAndResetMissCounters do slash any operator who over criteria & clear all operators miss counter to zero
func (k Keeper) SlashAndResetMissCounters(ctx sdk.Context) {
	height := ctx.BlockHeight()
	distributionHeight := height - sdk.ValidatorUpdateDelay - 1

	votePeriodsPerWindow := sdk.NewDec(k.SlashWindow(ctx)).QuoInt64(k.VotePeriod(ctx)).TruncateInt64()
	k.IterateMissCounters(ctx, func(operator sdk.ValAddress, missCounter int64) bool {

		// Calculate valid vote rate; (SlashWindow - MissCounter)/SlashWindow
		validVoteRate := sdk.NewDecFromInt(
			sdk.NewInt(votePeriodsPerWindow - missCounter)).
			QuoInt64(votePeriodsPerWindow)

		// Penalize the validator whose the valid vote rate is smaller than min threshold
		if validVoteRate.LT(k.MinValidPerWindow(ctx)) {
			validator := k.StakingKeeper.Validator(ctx, operator)
			if validator != nil && validator.IsBonded() && !validator.IsJailed() {
				cons, err := validator.GetConsAddr()
				if err != nil {
					panic(err)
				}
				k.StakingKeeper.Slash(
					ctx, cons,
					distributionHeight, validator.GetConsensusPower(), k.SlashFraction(ctx),
				)
				// TODO: Reenable jailing
				// k.StakingKeeper.Jail(ctx, cons)
			}
		}

		k.DeleteMissCounter(ctx, operator)
		return false
	})
}
