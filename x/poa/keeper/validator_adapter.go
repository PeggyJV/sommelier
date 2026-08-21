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
	multiplier sdk.Dec
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

func (b boostedValidator) GetConsensusPower(r math.Int) int64 {
	if !b.IsBonded() {
		return 0
	}
	raw := b.ValidatorI.GetConsensusPower(r)
	return sdk.NewDec(raw).Mul(b.multiplier).Ceil().TruncateInt64()
}
