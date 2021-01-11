package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/peggyjv/sommelier/x/oracle/types"
)

func (k Keeper) aggregateRates(ctx sdk.Context, votes map[string]types.ExchangeRateBallot, aggregateVoterMap map[string]bool) {
	k.IterateAggregateExchangeRateVotes(ctx, func(vote types.AggregateExchangeRateVote) (stop bool) {
		voter, err := sdk.ValAddressFromBech32(vote.Voter)
		if err != nil {
			// TODO: state machine panics still OK?
			panic(err)
		}

		validator := k.StakingKeeper.Validator(ctx, voter)

		// organize ballot only for the active validators
		if validator == nil || !validator.IsBonded() || validator.IsJailed() {
			return false
		}

		aggregateVoterMap[validator.GetOperator().String()] = true

		power := validator.GetConsensusPower()
		for _, tuple := range vote.ExchangeRateTuples {
			tmpPower := power
			if !tuple.Amount.IsPositive() {
				// Make the power of abstain vote zero
				tmpPower = 0
			}

			votes[tuple.Denom] = append(votes[tuple.Denom],
				types.NewVoteForTally(
					types.NewExchangeRateVote(tuple.Amount, tuple.Denom, voter),
					tmpPower,
				),
			)
		}

		return false
	})
}

func (k Keeper) organizeBallot(ctx sdk.Context, votes map[string]types.ExchangeRateBallot, aggregateVoterMap map[string]bool) {
	k.IterateExchangeRateVotes(ctx, func(vote types.ExchangeRateVote) (stop bool) {
		// block normal vote from the voter who did aggregate vote
		if _, ok := aggregateVoterMap[vote.Voter]; ok {
			return false
		}

		voter, err := sdk.ValAddressFromBech32(vote.Voter)
		if err != nil {
			// TODO: state machine panics still OK?
			panic(err)
		}

		validator := k.StakingKeeper.Validator(ctx, voter)

		// organize ballot only for the active validators
		if validator == nil || !validator.IsBonded() || validator.IsJailed() {
			return false
		}

		power := validator.GetConsensusPower()
		if !vote.ExchangeRate.IsPositive() {
			// Make the power of abstain vote zero
			power = 0
		}

		votes[vote.Denom] = append(votes[vote.Denom],
			types.NewVoteForTally(
				vote,
				power,
			),
		)

		return false
	})
}

// OrganizeBallotByDenom collects all oracle votes for the period, categorized by the votes' denom parameter
func (k Keeper) OrganizeBallotByDenom(ctx sdk.Context) map[string]types.ExchangeRateBallot {
	votes := map[string]types.ExchangeRateBallot{}
	aggregateVoterMap := map[string]bool{}

	k.aggregateRates(ctx, votes, aggregateVoterMap)
	k.organizeBallot(ctx, votes, aggregateVoterMap)
	return votes
}
