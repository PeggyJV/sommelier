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
	//
	// The import context is placed ABOVE the exported activation height, which
	// is what a real restart looks like: `sommelier export` writes
	// initial_height = H+1, so InitChain runs at a height past the stamp.
	// Importing at a lower height is a different case, covered by
	// TestGenesis_RejectsActivationStampFromAHigherHeightSpace.
	k2, ctx2, _, _ := newWrapperTestKeeper(t)
	ctx2 = ctx2.WithBlockHeight(4243)
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

// activation_height and safe_mode_thaw_height are absolute heights. Relaunching
// a chain from an export with a lower initial_height carries them into a height
// space they do not belong to.
//
// For activation_height the consequence is an inverted slash: rawSlashPower
// treats infractionHeight < activation as pre-PoA and returns the caller's
// power unchanged -- but that power came from LastValidatorPower, which the
// EndBlocker overwrote with the BOOSTED value. An authority validator that
// double-signs would be slashed on M times its real stake.
//
// For safe_mode_thaw_height the consequence is a permanent freeze: the
// EndBlocker can never reach a thaw height above the chain's own, and both
// MsgUpdateParams and MsgUpdateAuthoritySet are frozen in safe mode, so there
// is no on-chain exit.
func TestGenesis_RejectsActivationStampFromAHigherHeightSpace(t *testing.T) {
	k, ctx, _, _ := newWrapperTestKeeper(t)
	ctx = ctx.WithBlockHeight(1)

	keeper.InitGenesis(ctx, k, types.GenesisState{
		Params:             types.DefaultParams(),
		ActivationHeight:   15_000_000,
		ActivationTime:     1_700_000_000_000_000_000,
		SafeMode:           true,
		SafeModeThawHeight: 15_000_002,
	})

	height, ok := k.GetActivationHeight(ctx)
	require.True(t, ok)
	require.Equal(t, int64(1), height,
		"an activation stamp above the current height must be re-stamped, or every "+
			"infraction on the new chain is treated as pre-PoA and slashed on boosted power")

	// Read the thaw back through the public export path.
	round := keeper.ExportGenesis(ctx, k)
	require.True(t, round.SafeMode)
	require.LessOrEqual(t, round.SafeModeThawHeight, ctx.BlockHeight()+2,
		"an imported thaw height must be clamped into this chain's height space, "+
			"or safe mode never lifts and there is no on-chain recovery")
}
