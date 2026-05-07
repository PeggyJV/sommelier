package keeper

import (
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
)

// boostedValidator wraps a stakingtypes.ValidatorI and reports tokens scaled
// by the current-block authority multiplier. Token-share math is intentionally
// NOT overridden: delegation share allocation must operate on raw tokens so
// boosted authority validators do not dilute community delegators' shares.
type boostedValidator struct {
	stakingtypes.ValidatorI
	boostedTokens math.Int
}

func (b boostedValidator) GetTokens() math.Int { return b.boostedTokens }

func (b boostedValidator) GetBondedTokens() math.Int {
	if b.IsBonded() {
		return b.boostedTokens
	}
	return math.ZeroInt()
}

func (b boostedValidator) GetConsensusPower(r math.Int) int64 {
	if !b.IsBonded() {
		return 0
	}
	return sdk.TokensToConsensusPower(b.boostedTokens, r)
}
