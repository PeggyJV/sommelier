package keeper_test

import (
	"testing"

	abci "github.com/cometbft/cometbft/abci/types"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	cryptocodec "github.com/cosmos/cosmos-sdk/crypto/codec"
	"github.com/cosmos/cosmos-sdk/crypto/keys/ed25519"
	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/stretchr/testify/require"

	"github.com/peggyjv/sommelier/v9/x/poa/keeper"
)

// noopStakingEndBlocker is a StakingEndBlockerFn that does nothing — used by
// EndBlocker tests so they exercise the rescaling logic in isolation from
// staking's own state machine.
func noopStakingEndBlocker(_ sdk.Context) []abci.ValidatorUpdate { return nil }

// addValidatorWithPubkey is like fakeStakingKeeper.addValidator but also
// attaches a deterministic ed25519 pubkey so mergeUpdatesWithBoost can
// extract a TmConsPublicKey.
func (f *fakeStakingKeeper) addValidatorWithPubkey(t *testing.T, op sdk.ValAddress, tokens sdk.Int) stakingtypes.Validator {
	t.Helper()
	v := f.addValidator(op, tokens)
	pk := ed25519.GenPrivKeyFromSecret(op).PubKey()
	any, err := codectypes.NewAnyWithValue(pk)
	require.NoError(t, err)
	v.ConsensusPubkey = any
	f.validators[op.String()] = v
	return v
}

// init ensures crypto codec is registered so AnyWithValue serialises pubkeys.
func init() {
	registry := codectypes.NewInterfaceRegistry()
	cryptocodec.RegisterInterfaces(registry)
}

func TestEndBlocker_BoostsAuthorityToFloor(t *testing.T) {
	k, ctx, fake, _ := newWrapperTestKeeper(t)

	auth := sdk.ValAddress([]byte("auth-validator-aaaa"))
	com := sdk.ValAddress([]byte("com-validator-aaaaa"))
	fake.addValidatorWithPubkey(t, auth, sdk.NewInt(100*1_000_000))
	fake.addValidatorWithPubkey(t, com, sdk.NewInt(300*1_000_000))
	fake.bondedOrder = []sdk.ValAddress{auth, com}

	k.SetAuthoritySet(ctx, []sdk.ValAddress{auth})

	ctx = ctx.WithBlockHeight(42)
	keeper.EndBlocker(ctx, k, noopStakingEndBlocker)

	// authority should now have boosted LastValidatorPower
	authPower := fake.GetLastValidatorPower(ctx, auth)
	comPower := fake.GetLastValidatorPower(ctx, com)
	total := authPower + comPower
	require.True(t, total > 0)
	authShare := sdk.NewDec(authPower).Quo(sdk.NewDec(total))
	require.True(t, authShare.GTE(sdk.MustNewDecFromStr("0.67")),
		"authority share %s below floor", authShare)

	// snapshot was recorded
	snap, ok := k.GetMultiplierSnapshot(ctx, 42)
	require.True(t, ok)
	require.NotEmpty(t, snap.Entries)
	require.Equal(t, auth.String(), snap.Entries[0].OperatorAddress)
}

func TestEndBlocker_AlreadyAboveFloor_NoBoost(t *testing.T) {
	k, ctx, fake, _ := newWrapperTestKeeper(t)

	auth := sdk.ValAddress([]byte("auth-validator-aaaa"))
	com := sdk.ValAddress([]byte("com-validator-aaaaa"))
	// authority already at 90%
	fake.addValidatorWithPubkey(t, auth, sdk.NewInt(900*1_000_000))
	fake.addValidatorWithPubkey(t, com, sdk.NewInt(100*1_000_000))
	fake.bondedOrder = []sdk.ValAddress{auth, com}

	k.SetAuthoritySet(ctx, []sdk.ValAddress{auth})

	ctx = ctx.WithBlockHeight(42)
	authPowerBefore := fake.GetLastValidatorPower(ctx, auth)
	keeper.EndBlocker(ctx, k, noopStakingEndBlocker)
	authPowerAfter := fake.GetLastValidatorPower(ctx, auth)

	require.Equal(t, authPowerBefore, authPowerAfter, "no boost expected when above floor")

	// snapshot is still written, but with no entries
	snap, ok := k.GetMultiplierSnapshot(ctx, 42)
	require.True(t, ok)
	require.Empty(t, snap.Entries)
}

func TestEndBlocker_HaltOnEmptyAuthority(t *testing.T) {
	k, ctx, fake, _ := newWrapperTestKeeper(t)

	com := sdk.ValAddress([]byte("com-validator-aaaaa"))
	fake.addValidatorWithPubkey(t, com, sdk.NewInt(100*1_000_000))
	fake.bondedOrder = []sdk.ValAddress{com}

	// authority address that is NOT bonded
	missing := sdk.ValAddress([]byte("absent-aaaaaaaaaaaaa"))
	k.SetAuthoritySet(ctx, []sdk.ValAddress{missing})

	// Opt into fail-closed behavior (the default is safe mode).
	params := k.GetParams(ctx)
	params.HaltWhenAuthorityEmpty = true
	k.SetParams(ctx, params)

	require.Panics(t, func() {
		keeper.EndBlocker(ctx, k, noopStakingEndBlocker)
	})
}

// With the default params (HaltWhenAuthorityEmpty=false), an empty authority
// set enters safe mode instead of halting: blocks keep flowing on community
// stake and the safe-mode flag is set.
func TestEndBlocker_SafeModeOnEmptyAuthority(t *testing.T) {
	k, ctx, fake, _ := newWrapperTestKeeper(t)

	com := sdk.ValAddress([]byte("com-validator-aaaaa"))
	fake.addValidatorWithPubkey(t, com, sdk.NewInt(100*1_000_000))
	fake.bondedOrder = []sdk.ValAddress{com}

	missing := sdk.ValAddress([]byte("absent-aaaaaaaaaaaaa"))
	k.SetAuthoritySet(ctx, []sdk.ValAddress{missing})

	require.NotPanics(t, func() {
		keeper.EndBlocker(ctx, k, noopStakingEndBlocker)
	})
	require.True(t, k.SafeModeActive(ctx), "empty authority must enter safe mode by default")

	// Community validator runs unboosted; an empty snapshot is recorded.
	snap, ok := k.GetMultiplierSnapshot(ctx, ctx.BlockHeight())
	require.True(t, ok)
	require.Empty(t, snap.Entries)
}

// When a bonded authority validator returns, safe mode is held through the
// validator-update propagation delay (so value-bearing modules don't act in a
// block still secured by the old community-only set), then clears.
func TestEndBlocker_SafeModeThawsAfterDelay(t *testing.T) {
	k, ctx, fake, _ := newWrapperTestKeeper(t)

	auth := sdk.ValAddress([]byte("auth-validator-aaaa"))
	com := sdk.ValAddress([]byte("com-validator-aaaaa"))
	fake.addValidatorWithPubkey(t, com, sdk.NewInt(100*1_000_000))
	fake.bondedOrder = []sdk.ValAddress{com}
	k.SetAuthoritySet(ctx, []sdk.ValAddress{auth}) // auth not bonded yet

	ctx = ctx.WithBlockHeight(10)
	keeper.EndBlocker(ctx, k, noopStakingEndBlocker)
	require.True(t, k.SafeModeActive(ctx))

	// Authority validator bonds at height 11. Safe mode schedules a thaw but
	// stays active through the delay window.
	fake.addValidatorWithPubkey(t, auth, sdk.NewInt(50*1_000_000))
	fake.bondedOrder = []sdk.ValAddress{auth, com}
	keeper.EndBlocker(ctx.WithBlockHeight(11), k, noopStakingEndBlocker)
	require.True(t, k.SafeModeActive(ctx), "must stay frozen on the re-bond block (thaw scheduled)")
	keeper.EndBlocker(ctx.WithBlockHeight(12), k, noopStakingEndBlocker)
	require.True(t, k.SafeModeActive(ctx), "must stay frozen within the thaw delay")

	// At/after the thaw height (11 + 2 = 13), safe mode clears.
	keeper.EndBlocker(ctx.WithBlockHeight(13), k, noopStakingEndBlocker)
	require.False(t, k.SafeModeActive(ctx), "safe mode must clear once the restored set is securing consensus")
}

// If the authority set re-empties during the thaw window, the freeze is held
// (pending thaw cancelled) rather than thawing on schedule.
func TestEndBlocker_SafeModeReFreezesIfAuthorityReEmpties(t *testing.T) {
	k, ctx, fake, _ := newWrapperTestKeeper(t)

	auth := sdk.ValAddress([]byte("auth-validator-aaaa"))
	com := sdk.ValAddress([]byte("com-validator-aaaaa"))
	fake.addValidatorWithPubkey(t, com, sdk.NewInt(100*1_000_000))
	fake.bondedOrder = []sdk.ValAddress{com}
	k.SetAuthoritySet(ctx, []sdk.ValAddress{auth})

	keeper.EndBlocker(ctx.WithBlockHeight(10), k, noopStakingEndBlocker)
	require.True(t, k.SafeModeActive(ctx))

	// Authority bonds at 11 (thaw scheduled for 13)...
	fake.addValidatorWithPubkey(t, auth, sdk.NewInt(50*1_000_000))
	fake.bondedOrder = []sdk.ValAddress{auth, com}
	keeper.EndBlocker(ctx.WithBlockHeight(11), k, noopStakingEndBlocker)

	// ...but is unbonded again at 12 (dropped from the bonded set), before the
	// thaw completes.
	fake.bondedOrder = []sdk.ValAddress{com}
	keeper.EndBlocker(ctx.WithBlockHeight(12), k, noopStakingEndBlocker)
	require.True(t, k.SafeModeActive(ctx))

	// The cancelled thaw must not fire at the old scheduled height.
	keeper.EndBlocker(ctx.WithBlockHeight(13), k, noopStakingEndBlocker)
	require.True(t, k.SafeModeActive(ctx), "freeze must hold; the prior thaw was cancelled when authority re-emptied")
}

func TestEndBlocker_DisabledIsNoop(t *testing.T) {
	k, ctx, fake, _ := newWrapperTestKeeper(t)

	auth := sdk.ValAddress([]byte("auth-validator-aaaa"))
	com := sdk.ValAddress([]byte("com-validator-aaaaa"))
	fake.addValidatorWithPubkey(t, auth, sdk.NewInt(100*1_000_000))
	fake.addValidatorWithPubkey(t, com, sdk.NewInt(300*1_000_000))
	fake.bondedOrder = []sdk.ValAddress{auth, com}
	k.SetAuthoritySet(ctx, []sdk.ValAddress{auth})

	params := k.GetParams(ctx)
	params.Enabled = false
	k.SetParams(ctx, params)

	authPowerBefore := fake.GetLastValidatorPower(ctx, auth)
	keeper.EndBlocker(ctx, k, noopStakingEndBlocker)
	authPowerAfter := fake.GetLastValidatorPower(ctx, auth)
	require.Equal(t, authPowerBefore, authPowerAfter)

	// When disabled, an empty snapshot is still written so any future slash
	// for a height in this window does not hit the missing-snapshot refuse
	// path in WrappedStakingKeeper.Slash.
	snap, ok := k.GetMultiplierSnapshot(ctx, ctx.BlockHeight())
	require.True(t, ok)
	require.Empty(t, snap.Entries)
}
