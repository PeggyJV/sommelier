package keeper_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"

	"github.com/peggyjv/sommelier/v10/x/poa/keeper"
	"github.com/peggyjv/sommelier/v10/x/poa/types"
)

// A genesis export/import round-trip must preserve the activation stamp, the
// safe-mode flag, and the multiplier snapshots.
//
// Losing the activation height is the dangerous one: rawSlashPower slashes
// infraction heights BELOW activation verbatim (they predate any boost), so a
// reset activation height would make a post-activation infraction get slashed
// on boosted power instead of raw stake. Losing the safe-mode flag would
// unfreeze the value-bearing modules for the blocks before the first EndBlocker
// re-derives it.
func TestGenesis_RoundTripPreservesActivationSafeModeAndSnapshots(t *testing.T) {
	k, ctx, _, _ := newWrapperTestKeeper(t)

	auth := sdk.ValAddress([]byte("auth-validator-aaaa"))
	k.SetAuthoritySet(ctx, []sdk.ValAddress{auth})
	k.SetActivationStamp(ctx, 4242, 1_700_000_000_000_000_000)
	k.SetSafeMode(ctx, true)
	k.SetMultiplierSnapshot(ctx, types.MultiplierSnapshot{
		Height: 4300,
		Entries: []*types.MultiplierEntry{
			{OperatorAddress: auth.String(), Multiplier: "3.5"},
		},
	})
	k.SetMultiplierSnapshot(ctx, types.MultiplierSnapshot{Height: 4301})

	exported := keeper.ExportGenesis(ctx, k)
	require.Equal(t, int64(4242), exported.ActivationHeight)
	require.Equal(t, int64(1_700_000_000_000_000_000), exported.ActivationTime)
	require.True(t, exported.SafeMode)
	require.Len(t, exported.MultiplierSnapshots, 2)
	require.NoError(t, exported.Validate())

	// Import into a fresh keeper and confirm the security-relevant state landed.
	k2, ctx2, _, _ := newWrapperTestKeeper(t)
	keeper.InitGenesis(ctx2, k2, exported)

	height, ok := k2.GetActivationHeight(ctx2)
	require.True(t, ok)
	require.Equal(t, int64(4242), height)

	at, ok := k2.GetActivationTime(ctx2)
	require.True(t, ok)
	require.Equal(t, int64(1_700_000_000_000_000_000), at)

	require.True(t, k2.SafeModeActive(ctx2))

	m, found := k2.MultiplierForValidatorWithStatus(ctx2, auth, 4300)
	require.True(t, found)
	require.Equal(t, sdk.MustNewDecFromStr("3.5"), m)

	// The un-boosted height round-trips as a present-but-empty snapshot, which
	// is what keeps its slashes on the pass-through path rather than the
	// refuse path.
	_, found = k2.MultiplierForValidatorWithStatus(ctx2, auth, 4301)
	require.True(t, found)
}

// A genesis with no activation stamp (a genuine first launch) stamps the
// current height rather than leaving the module unactivated.
func TestGenesis_FreshLaunchStampsCurrentHeight(t *testing.T) {
	k, ctx, _, _ := newWrapperTestKeeper(t)
	ctx = ctx.WithBlockHeight(77)

	keeper.InitGenesis(ctx, k, *types.DefaultGenesis())

	height, ok := k.GetActivationHeight(ctx)
	require.True(t, ok)
	require.Equal(t, int64(77), height)
	require.False(t, k.SafeModeActive(ctx))
}
