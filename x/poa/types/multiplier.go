package types

import (
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// ComputeMultiplier returns the per-authority-validator power multiplier M
// such that boosted authority share equals at least `floor`:
//
//	M = floor / (1 - floor) * C / B
//
// where B is the aggregate authority power and C is the aggregate community
// power.
//
// Behaviour:
//   - returns 1 when authority is already at or above the floor (no boost
//     needed) or when the community bucket is empty
//   - returns 0 when B is zero (caller must decide halt vs. pass-through)
func ComputeMultiplier(authorityPower, communityPower math.Int, floor sdk.Dec) sdk.Dec {
	if authorityPower.IsZero() {
		return sdk.ZeroDec()
	}
	if communityPower.IsZero() {
		return sdk.OneDec()
	}
	total := authorityPower.Add(communityPower)
	share := sdk.NewDecFromInt(authorityPower).Quo(sdk.NewDecFromInt(total))
	if share.GTE(floor) {
		return sdk.OneDec()
	}
	ratio := floor.Quo(sdk.OneDec().Sub(floor))
	cb := sdk.NewDecFromInt(communityPower).Quo(sdk.NewDecFromInt(authorityPower))
	return ratio.Mul(cb)
}
