package keeper

import (
	"sort"

	"cosmossdk.io/math"

	abci "github.com/cometbft/cometbft/abci/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/peggyjv/sommelier/v7/x/incentives/types"
)

type ValidatorInfo struct {
	Validator stakingtypes.ValidatorI
	Power     int64
}

// sortValidatorInfosByPower sorts the validator information by power in descending order
func sortValidatorInfosByPower(valInfos []ValidatorInfo) []ValidatorInfo {
	sort.Slice(valInfos, func(i, j int) bool {
		return valInfos[i].Power > valInfos[j].Power
	})

	return valInfos
}

// GetTotalPower returns the total power of the passed in validatorInfos
func getTotalPower(valInfos *[]ValidatorInfo) int64 {
	totalPower := int64(0)
	for _, valInfo := range *valInfos {
		totalPower += valInfo.Power
	}

	return totalPower
}

// getValidatorInfos returns the validator information for the voters in the last block
func (k Keeper) getValidatorInfos(ctx sdk.Context, req abci.RequestBeginBlock) []ValidatorInfo {
	validatorInfos := []ValidatorInfo{}
	for _, vote := range req.LastCommitInfo.GetVotes() {
		if !vote.SignedLastBlock {
			continue
		}

		validator := k.StakingKeeper.ValidatorByConsAddr(ctx, vote.Validator.Address)
		validatorInfos = append(validatorInfos, ValidatorInfo{
			Validator: validator,
			Power:     vote.Validator.Power,
		})
	}
	return validatorInfos
}

// AllocateTokens performs reward distribution to the provided validators proportionally to their power with a cap
func (k Keeper) AllocateTokens(ctx sdk.Context, totalPreviousPower int64, totalDistribution sdk.DecCoins, qualifyingVoters []ValidatorInfo, maxFraction sdk.Dec) sdk.DecCoins {
	remaining := totalDistribution

	for _, valInfo := range qualifyingVoters {
		validator := valInfo.Validator
		powerFraction := math.LegacyNewDec(valInfo.Power).QuoInt64(totalPreviousPower)

		// Cap at the max fraction
		if powerFraction.GT(maxFraction) {
			powerFraction = maxFraction
		}

		reward := totalDistribution.MulDecTruncate(powerFraction)

		k.AllocateTokensToValidator(ctx, validator, reward)
		remaining = remaining.Sub(reward)
	}

	return remaining
}

// AllocateTokensToValidator allocates tokens to a particular validator.
// All tokens go to the validator.
func (k Keeper) AllocateTokensToValidator(ctx sdk.Context, val stakingtypes.ValidatorI, tokens sdk.DecCoins) {
	// Update validator rewards
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeIncentivesRewards,
			sdk.NewAttribute(sdk.AttributeKeyAmount, tokens.String()),
			sdk.NewAttribute(types.AttributeKeyValidator, val.GetOperator().String()),
		),
	)

	// Update current rewards
	currentRewards := k.DistributionKeeper.GetValidatorCurrentRewards(ctx, val.GetOperator())
	currentRewards.Rewards = currentRewards.Rewards.Add(tokens...)
	k.DistributionKeeper.SetValidatorCurrentRewards(ctx, val.GetOperator(), currentRewards)

	// Update outstanding rewards
	outstanding := k.DistributionKeeper.GetValidatorOutstandingRewards(ctx, val.GetOperator())
	outstanding.Rewards = outstanding.Rewards.Add(tokens...)
	k.DistributionKeeper.SetValidatorOutstandingRewards(ctx, val.GetOperator(), outstanding)
}
