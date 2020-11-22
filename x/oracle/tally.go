package oracle

import (
	"github.com/peggyjv/sommelier/x/oracle/types"
)

func updateWinnerMap(ballotWinningClaims []types.Claim, validVotesCounterMap map[string]int, winnerMap map[string]types.Claim) {
	// Collect claims of ballot winners
	for _, ballotWinningClaim := range ballotWinningClaims {

		// NOTE: we directly stringify byte to string to prevent unnecessary bech32fy works
		key := string(ballotWinningClaim.Recipient)

		// Update claim
		prevClaim := winnerMap[key]
		prevClaim.Weight += ballotWinningClaim.Weight
		winnerMap[key] = prevClaim

		// Increase valid votes counter
		validVotesCounterMap[key]++
	}
}
