package keeper_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"

	"github.com/peggyjv/sommelier/v10/x/poa/keeper"
	"github.com/peggyjv/sommelier/v10/x/poa/types"
)

// testGovAuthority is the gov module account address the test keepers are
// constructed with. Derived rather than hardcoded so it is valid bech32 under
// whatever address prefix the SDK config carries.
var testGovAuthority = sdk.AccAddress([]byte("gov-authority-aaaaaa")).String()

// While the chain is in authority-empty safe mode, the ONLY bonded stake is
// community stake, so governance is decided entirely by the set the freeze
// exists to distrust. If x/poa's own gov messages stayed open there, an
// attacker who induced the authority outage could vote itself the authority
// set (or just flip Enabled=false, which clears safe mode with no thaw delay)
// and own the bridge — exactly the escalation Option A claims to prevent.
func TestUpdateAuthoritySet_FrozenInSafeMode(t *testing.T) {
	k, ctx, fake, _ := newWrapperTestKeeper(t)
	srv := keeper.NewMsgServerImpl(k)

	newAuth := sdk.ValAddress([]byte("new-authority-aaaaa"))
	fake.addValidator(newAuth, sdk.NewInt(1_000_000))

	k.SetSafeMode(ctx, true)

	_, err := srv.UpdateAuthoritySet(sdk.WrapSDKContext(ctx), &types.MsgUpdateAuthoritySet{
		Authority:  testGovAuthority,
		Validators: []string{newAuth.String()},
	})
	require.ErrorIs(t, err, types.ErrSafeModeGovFrozen)
}

func TestUpdateParams_FrozenInSafeMode(t *testing.T) {
	k, ctx, _, _ := newWrapperTestKeeper(t)
	srv := keeper.NewMsgServerImpl(k)

	k.SetSafeMode(ctx, true)

	// Enabled=false is the dangerous one: EndBlocker clears safe mode
	// immediately on a disabled module, with no thaw delay.
	params := types.DefaultParams()
	params.Enabled = false

	_, err := srv.UpdateParams(sdk.WrapSDKContext(ctx), &types.MsgUpdateParams{
		Authority: testGovAuthority,
		Params:    params,
	})
	require.ErrorIs(t, err, types.ErrSafeModeGovFrozen)

	// The params must be untouched.
	require.True(t, k.GetParams(ctx).Enabled)
}

// Outside safe mode the same messages work normally.
func TestUpdateAuthoritySet_AllowedOutsideSafeMode(t *testing.T) {
	k, ctx, fake, _ := newWrapperTestKeeper(t)
	srv := keeper.NewMsgServerImpl(k)

	newAuth := sdk.ValAddress([]byte("new-authority-aaaaa"))
	fake.addValidator(newAuth, sdk.NewInt(1_000_000))

	_, err := srv.UpdateAuthoritySet(sdk.WrapSDKContext(ctx), &types.MsgUpdateAuthoritySet{
		Authority:  testGovAuthority,
		Validators: []string{newAuth.String()},
	})
	require.NoError(t, err)
	require.Equal(t, []sdk.ValAddress{newAuth}, k.GetAuthoritySet(ctx))
}

// A set with no bonded, unjailed member would drop the chain into safe mode on
// the very next EndBlocker, freezing gravity/cork/axelarcork the moment the
// proposal executes. Reject it at the msg server instead.
func TestUpdateAuthoritySet_RejectsSetWithNoBondedValidator(t *testing.T) {
	k, ctx, fake, _ := newWrapperTestKeeper(t)
	srv := keeper.NewMsgServerImpl(k)

	// Known to staking, but jailed.
	jailed := sdk.ValAddress([]byte("jailed-validator-aa"))
	v := fake.addValidator(jailed, sdk.NewInt(1_000_000))
	v.Jailed = true
	fake.setValidator(v)

	// Never registered with staking at all.
	unknown := sdk.ValAddress([]byte("unknown-validator-a"))

	_, err := srv.UpdateAuthoritySet(sdk.WrapSDKContext(ctx), &types.MsgUpdateAuthoritySet{
		Authority:  testGovAuthority,
		Validators: []string{jailed.String(), unknown.String()},
	})
	require.ErrorIs(t, err, types.ErrNoBondedAuthority)

	// Unjailing one of them makes the same message acceptable.
	v.Jailed = false
	fake.setValidator(v)
	_, err = srv.UpdateAuthoritySet(sdk.WrapSDKContext(ctx), &types.MsgUpdateAuthoritySet{
		Authority:  testGovAuthority,
		Validators: []string{jailed.String(), unknown.String()},
	})
	require.NoError(t, err)
}
