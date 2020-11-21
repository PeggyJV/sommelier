package keeper

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/peggyjv/sommelier/x/oracle/types"
)

// RewardBallotWinners implements
// at the end of every VotePeriod, we give out portion of seigniorage reward(reward-weight) to the
// oracle voters that voted faithfully.
func (k Keeper) RewardBallotWinners(ctx sdk.Context, ballotWinners map[string]types.Claim) {
	// Sum weight of the claims
	ballotPowerSum := int64(0)
	for _, winner := range ballotWinners {
		ballotPowerSum += winner.Weight
	}

	// Exit if the ballot is empty
	if ballotPowerSum == 0 {
		return
	}

	rewardPool := k.GetRewardPool(ctx)

	// return if there's no rewards to give out
	if rewardPool.IsZero() {
		return
	}

	// rewardCoin  = oraclePool * VotePeriod / RewardDistributionWindow
	periodRewards := sdk.NewDecFromInt(rewardPool.AmountOf(types.MicroLunaDenom)).
		MulInt64(k.VotePeriod(ctx)).QuoInt64(k.RewardDistributionWindow(ctx))

	// Dole out rewards
	var distributedReward sdk.Coins
	for _, winner := range ballotWinners {
		rewardCoins := sdk.NewCoins()
		valAddr, err := sdk.ValAddressFromBech32(winner.Recipient)
		if err != nil {
			panic(err)
		}
		rewardeeVal := k.StakingKeeper.Validator(ctx, valAddr)

		// Reflects contribution
		rewardAmt := periodRewards.QuoInt64(ballotPowerSum).MulInt64(winner.Weight).TruncateInt()
		rewardCoins = append(rewardCoins, sdk.NewCoin(types.MicroLunaDenom, rewardAmt))

		// In case absence of the validator, we just skip distribution
		if rewardeeVal != nil && !rewardCoins.IsZero() {
			k.distrKeeper.AllocateTokensToValidator(ctx, rewardeeVal, sdk.NewDecCoinsFromCoins(rewardCoins...))
			distributedReward = distributedReward.Add(rewardCoins...)
		}
	}

	// Move distributed reward to distribution module
	err := k.bankKeeper.SendCoinsFromModuleToModule(ctx, types.ModuleName, k.distrName, distributedReward)
	if err != nil {
		panic(fmt.Sprintf("[oracle] Failed to send coins to distribution module %s", err.Error()))
	}

}
