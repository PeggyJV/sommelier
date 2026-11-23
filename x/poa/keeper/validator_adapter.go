package keeper

import (
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
)

// boostedValidator wraps a stakingtypes.ValidatorI and reports tokens and
// consensus power scaled by the current-block authority multiplier. Both
// scalings use Ceil to stay consistent with the EndBlocker's
// LastValidatorPower overwrite: floor on read + ceil on write would cause
// boundary under-slashing of authority validators.
//
// Token-share math (TokensFromShares*, SharesFromTokens*) is intentionally
// NOT overridden: delegation share allocation must operate on raw tokens so
// boosted authority validators do not dilute community delegators' shares.
type boostedValidator struct {
	stakingtypes.ValidatorI
	multiplier           sdk.Dec
	scaleDelegatorShares bool
}

func (b boostedValidator) GetTokens() math.Int {
	return sdk.NewDecFromInt(b.ValidatorI.GetTokens()).Mul(b.multiplier).Ceil().TruncateInt()
}

func (b boostedValidator) GetBondedTokens() math.Int {
	if !b.IsBonded() {
		return math.ZeroInt()
	}
	return b.GetTokens()
}

// GetDelegatorShares scales the governance tally's denominator by the exact
// effective token boost (including GetTokens' integer rounding). The SDK tally
// leaves each stored delegation's shares raw, so an external delegation keeps
// raw-stake voting power while the remaining authority boost belongs to the
// validator operator. This override is enabled only by ForGovernance; token-
// share math for staking and every other wrapped consumer stays raw.
func (b boostedValidator) GetDelegatorShares() sdk.Dec {
	rawShares := b.ValidatorI.GetDelegatorShares()
	if !b.scaleDelegatorShares {
		return rawShares
	}
	rawTokens := b.ValidatorI.GetTokens()
	if rawTokens.IsZero() {
		return rawShares
	}
	return rawShares.MulInt(b.GetTokens()).QuoInt(rawTokens)
}

func (b boostedValidator) GetConsensusPower(r math.Int) int64 {
	if !b.IsBonded() {
		return 0
	}
	raw := b.ValidatorI.GetConsensusPower(r)
	return sdk.NewDec(raw).Mul(b.multiplier).Ceil().TruncateInt64()
}
