package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/peggyjv/sommelier/x/oracle/types"
)

// OrganizeBallotByDenom collects all oracle votes for the period, categorized by the votes' denom parameter
func (k Keeper) OrganizeBallotByDenom(ctx sdk.Context) (votes map[string]types.ExchangeRateVote) {
	votes = map[string]types.ExchangeRateVote{}
	aggregateVoterMap := map[string]bool{}

	// Organize aggregate votes
	k.IterateAggregateExchangeRateVotes(ctx, func(vote *types.AggregateExchangeRateVote) (stop bool) {
		validator := k.StakingKeeper.Validator(ctx, string(vote.Voter))

		// organize ballot only for the active validators
		if validator != nil && validator.IsBonded() && !validator.IsJailed() {
			aggregateVoterMap[string(validator.GetOperator().Bytes())] = true

			power := validator.GetConsensusPower()
			for _, tuple := range vote.ExchangeRateTuples {
				tmpPower := power
				if !tuple.ExchangeRate.IsPositive() {
					// Make the power of abstain vote zero
					tmpPower = 0
				}

				votes[tuple.Denom] = append(votes[tuple.Denom],
					types.NewVoteForTally(
						types.NewExchangeRateVote(tuple.ExchangeRate, tuple.Denom, vote.Voter),
						tmpPower,
					),
				)
			}

		}

		return false
	})
	return
}
