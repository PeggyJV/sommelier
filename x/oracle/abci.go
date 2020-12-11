package oracle

import (
	"github.com/peggyjv/sommelier/x/oracle/keeper"
	"github.com/peggyjv/sommelier/x/oracle/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
)

// EndBlocker is called at the end of every block
func EndBlocker(ctx sdk.Context, k keeper.Keeper) {
	params := k.GetParams(ctx)

	// Not yet time for a tally
	if !IsPeriodLastBlock(ctx, params.VotePeriod) {
		return
	}

	// Build valid votes counter and winner map over all validators in active set
	validVotesCounterMap := make(map[string]int)
	winnerMap := make(map[string]types.Claim)
	k.StakingKeeper.IterateValidators(ctx, func(i int64, validator stakingtypes.ValidatorI) bool {
		// Exclude not bonded validator or jailed validators from tallying
		if validator.IsBonded() && !validator.IsJailed() {

			// NOTE: we directly stringify byte to string to prevent unnecessary bech32fy works
			valAddr := validator.GetOperator()
			validVotesCounterMap[valAddr.String()] = 0
			winnerMap[valAddr.String()] = types.NewClaim(0, valAddr)
		}

		return false
	})

	// Denom-TobinTax map
	voteTargets := make(map[string]sdk.Dec)
	k.IterateTobinTaxes(ctx, func(denom string, tobinTax sdk.Dec) bool {
		voteTargets[denom] = tobinTax
		return false
	})

	// Clear all exchange rates
	k.IterateUSDExchangeRates(ctx, func(denom string, _ sdk.Dec) (stop bool) {
		k.DeleteUSDExchangeRate(ctx, denom)
		return false
	})

	// Organize votes to ballot by denom
	// NOTE: **Filter out inactive or jailed validators**
	// NOTE: **Make abstain votes to have zero vote power**
	voteMap := k.OrganizeBallotByDenom(ctx)

	if referenceTerra := pickReferenceTerra(ctx, k, voteTargets, voteMap); referenceTerra != "" {
		// make voteMap of Reference Terra to calculate cross exchange rates
		ballotRT := voteMap[referenceTerra]
		voteMapRT := ballotRT.ToMap()
		exchangeRateRT := ballotRT.WeightedMedian()

		// Iterate through ballots and update exchange rates; drop if not enough votes have been achieved.
		for denom, ballot := range voteMap {

			// Convert ballot to cross exchange rates
			if denom != referenceTerra {
				ballot = ballot.ToCrossRate(voteMapRT)
			}

			// Get weighted median of cross exchange rates
			exchangeRate, ballotWinningClaims := tally(ctx, ballot, params.RewardBand)

			// Update winnerMap, validVotesCounterMap using ballotWinningClaims of cross exchange rate ballot
			updateWinnerMap(ballotWinningClaims, validVotesCounterMap, winnerMap)

			// Transform into the original form uluna/stablecoin
			if denom != referenceTerra {
				exchangeRate = exchangeRateRT.Quo(exchangeRate)
			}

			// Set the exchange rate, emit ABCI event
			k.SetUSDExchangeRateWithEvent(ctx, denom, exchangeRate)
		}
	}

	//---------------------------
	// Do miss counting & slashing
	voteTargetsLen := len(voteTargets)
	for operatorAddrByteStr, count := range validVotesCounterMap {
		// Skip abstain & valid voters
		if count == voteTargetsLen {
			continue
		}

		// Increase miss counter
		operator, err := sdk.ValAddressFromBech32(operatorAddrByteStr)
		if err != nil {
			panic(err)
		}
		k.SetMissCounter(ctx, operator, k.GetMissCounter(ctx, operator)+1)
	}

	// Do slash who did miss voting over threshold and
	// reset miss counters of all validators at the last block of slash window
	if IsPeriodLastBlock(ctx, params.SlashWindow) {
		SlashAndResetMissCounters(ctx, k)
	}

	// Distribute rewards to ballot winners
	k.RewardBallotWinners(ctx, winnerMap)

	// Clear the ballot
	clearBallots(ctx, k, params.VotePeriod)

	// Update vote targets and tobin tax
	applyWhitelist(ctx, k, params.Whitelist, voteTargets)

	return
}

// clearBallots clears all tallied prevotes and votes from the store
func clearBallots(ctx sdk.Context, k keeper.Keeper, votePeriod int64) {
	// Clear all prevotes
	k.IterateExchangeRatePrevotes(ctx, func(prevote types.ExchangeRatePrevote) (stop bool) {
		if ctx.BlockHeight() > prevote.SubmitBlock+votePeriod {
			k.DeleteExchangeRatePrevote(ctx, prevote)
		}

		return false
	})

	// Clear all votes
	k.IterateExchangeRateVotes(ctx, func(vote types.ExchangeRateVote) (stop bool) {
		k.DeleteExchangeRateVote(ctx, vote)
		return false
	})

	// Clear all aggregate prevotes
	k.IterateAggregateExchangeRatePrevotes(ctx, func(aggregatePrevote types.AggregateExchangeRatePrevote) (stop bool) {
		if ctx.BlockHeight() > aggregatePrevote.SubmitBlock+votePeriod {
			k.DeleteAggregateExchangeRatePrevote(ctx, aggregatePrevote)
		}

		return false
	})

	// Clear all aggregate votes
	k.IterateAggregateExchangeRateVotes(ctx, func(vote types.AggregateExchangeRateVote) (stop bool) {
		k.DeleteAggregateExchangeRateVote(ctx, vote)
		return false
	})
}

// applyWhitelist update vote target denom list and set tobin tax with params whitelist
func applyWhitelist(ctx sdk.Context, k keeper.Keeper, whitelist sdk.DecCoins, voteTargets map[string]sdk.Dec) {

	// check is there any update in whitelist params
	updateRequired := false
	if len(voteTargets) != len(whitelist) {
		updateRequired = true
	} else {
		for _, item := range whitelist {
			if tobinTax, ok := voteTargets[item.Denom]; !ok || !tobinTax.Equal(item.Amount) {
				updateRequired = true
				break
			}
		}
	}

	if updateRequired {
		k.ClearTobinTaxes(ctx)

		for _, item := range whitelist {
			k.SetTobinTax(ctx, item.Denom, item.Amount)
		}
	}
}

// IsPeriodLastBlock returns true if we are at the last block of the period
func IsPeriodLastBlock(ctx sdk.Context, blocksPerPeriod int64) bool {
	return (ctx.BlockHeight()+1)%blocksPerPeriod == 0
}
