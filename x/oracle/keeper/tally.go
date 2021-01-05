package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/peggyjv/sommelier/x/oracle/types"
)

// ballot for the asset is passing the threshold amount of voting power
func (k Keeper) BallotIsPassing(ctx sdk.Context, ballot types.ExchangeRateBallot) (sdk.Int, bool) {
	totalBondedPower := sdk.TokensToConsensusPower(k.StakingKeeper.TotalBondedTokens(ctx))
	voteThreshold := k.VoteThreshold(ctx)
	thresholdVotes := voteThreshold.MulInt64(totalBondedPower).RoundInt()
	ballotPower := sdk.NewInt(ballot.Power())
	return ballotPower, !ballotPower.IsZero() && ballotPower.GTE(thresholdVotes)
}

// choose Reference Terra with the highest voter turnout
// If the voting power of the two denominations is the same,
// select reference Terra in alphabetical order.
func (k Keeper) PickReferenceTerra(ctx sdk.Context, voteTargets map[string]sdk.Dec, voteMap map[string]types.ExchangeRateBallot) string {
	largestBallotPower := int64(0)
	referenceTerra := ""

	for denom, ballot := range voteMap {
		// If denom is not in the voteTargets, or the ballot for it has failed, then skip
		// and remove it from voteMap for iteration efficiency
		if _, exists := voteTargets[denom]; !exists {
			delete(voteMap, denom)
			continue
		}

		ballotPower := int64(0)

		// If the ballot is not passed, remove it from the voteTargets array
		// to prevent slashing validators who did valid vote.
		if power, ok := k.BallotIsPassing(ctx, ballot); ok {
			ballotPower = power.Int64()
		} else {
			delete(voteTargets, denom)
			delete(voteMap, denom)
			continue
		}

		if ballotPower > largestBallotPower || largestBallotPower == 0 {
			referenceTerra = denom
			largestBallotPower = ballotPower
		} else if largestBallotPower == ballotPower && referenceTerra > denom {
			referenceTerra = denom
		}
	}

	return referenceTerra
}
