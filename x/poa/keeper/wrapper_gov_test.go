package keeper_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/stretchr/testify/require"

	"github.com/peggyjv/sommelier/v10/x/poa/types"
)

// x/gov builds its tally numerator from per-validator GetBondedTokens (which
// this wrapper boosts) and its quorum denominator from TotalBondedTokens. If
// the denominator stayed raw, turnout would be overstated by the boost factor
// and could exceed 100%. Both sides must be boosted the same way.
func TestWrapper_TotalBondedTokens_MatchesBoostedTally(t *testing.T) {
	k, ctx, fake, w := newWrapperTestKeeper(t)

	auth := sdk.ValAddress([]byte("auth-validator-aaaa"))
	com := sdk.ValAddress([]byte("com-validator-aaaaa"))
	fake.addValidator(auth, sdk.NewInt(1_000_000))
	fake.addValidator(com, sdk.NewInt(3_000_000))
	fake.bondedOrder = []sdk.ValAddress{auth, com}

	k.SetAuthoritySet(ctx, []sdk.ValAddress{auth})
	k.SetMultiplierSnapshot(ctx, types.MultiplierSnapshot{
		Height: ctx.BlockHeight(),
		Entries: []*types.MultiplierEntry{
			{OperatorAddress: auth.String(), Multiplier: "6.0"},
		},
	})

	// Sum exactly the per-validator quantity x/gov's tally numerator uses.
	tallyNumeratorBasis := sdk.ZeroInt()
	w.IterateBondedValidatorsByPower(ctx, func(_ int64, v stakingtypes.ValidatorI) bool {
		tallyNumeratorBasis = tallyNumeratorBasis.Add(v.GetBondedTokens())
		return false
	})

	require.Equal(t, tallyNumeratorBasis, w.TotalBondedTokens(ctx),
		"quorum denominator must equal the sum of the boosted per-validator bonded tokens")

	// Sanity: boosted (6,000,000 + 3,000,000), not the raw 4,000,000.
	require.Equal(t, sdk.NewInt(9_000_000), w.TotalBondedTokens(ctx))
}
