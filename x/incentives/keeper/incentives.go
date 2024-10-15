package keeper

import (
	"fmt"
	"sort"

	abci "github.com/cometbft/cometbft/abci/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type VoterInfo struct {
	Validator *abci.Validator
}

// GetSortedVoterInfosByPower returns the previous block's voter information by validator power in descending order
func GetSortedVoterInfosByPower(votes []abci.VoteInfo) []VoterInfo {
	voterInfos := []VoterInfo{}
	for i := range votes {
		if !votes[i].SignedLastBlock {
			continue
		}

		voterInfos = append(voterInfos, VoterInfo{
			Validator: &votes[i].Validator,
		})
	}

	// Sort voteInfos by descending Power
	sort.Slice(voterInfos, func(i, j int) bool {
		return voterInfos[i].Validator.Power > voterInfos[j].Validator.Power
	})

	return voterInfos
}

// GetApportionments returns a slice of fractions of the passed in value that sums to the value approximately, and
// a remaining value. The sum of the returned slice plus the remaining value will equal the original value.
func getApportionments(numPortions uint64, value sdk.Dec, maxPortionFraction sdk.Dec) ([]sdk.Dec, sdk.Dec, error) {
	// We error check for sanity, even though the arguments should only be coming from validated Param values
	if numPortions == 0 {
		return make([]sdk.Dec, 0), value, nil
	}

	if value.IsNegative() {
		value = sdk.ZeroDec()
	}

	if maxPortionFraction.IsNegative() {
		return nil, sdk.ZeroDec(), fmt.Errorf("max portion cannot be negative")
	} else if maxPortionFraction.GT(sdk.OneDec()) {
		return nil, sdk.ZeroDec(), fmt.Errorf("max portion must be less than or equal to one")
	}

	remainingValue := value
	apportionments := make([]sdk.Dec, numPortions)

	for i := 0; i < len(apportionments); i++ {
		if remainingValue.IsZero() || maxPortionFraction.IsZero() {
			apportionments[i] = sdk.ZeroDec()
			continue
		}

		portion := remainingValue.Mul(maxPortionFraction)
		apportionments[i] = portion
		remainingValue = remainingValue.Sub(portion)
	}

	return apportionments, remainingValue, nil
}
