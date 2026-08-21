package app

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/authz"
	govtypesv1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1"
	gravitytypes "github.com/peggyjv/gravity-bridge/module/v6/x/gravity/types"
	"github.com/stretchr/testify/require"
)

type stubSafeMode struct{ active bool }

func (s stubSafeMode) SafeModeActive(sdk.Context) bool { return s.active }

type stubTx struct{ msgs []sdk.Msg }

func (t stubTx) GetMsgs() []sdk.Msg                { return t.msgs }
func (t stubTx) ValidateBasic() error              { return nil }
func (t stubTx) GetMsgsV2() ([]interface{}, error) { return nil, nil } //nolint:unused

func TestSafeModeAnteHandler(t *testing.T) {
	called := false
	next := func(ctx sdk.Context, tx sdk.Tx, simulate bool) (sdk.Context, error) {
		called = true
		return ctx, nil
	}
	ctx := sdk.Context{}

	frozen := []sdk.Msg{
		&gravitytypes.MsgSendToEthereum{},
		&gravitytypes.MsgSubmitEthereumEvent{},
		&gravitytypes.MsgSubmitEthereumTxConfirmation{},
	}
	allowed := []sdk.Msg{
		&gravitytypes.MsgDelegateKeys{},
		&gravitytypes.MsgCancelSendToEthereum{},
		&gravitytypes.MsgEthereumHeightVote{},
	}

	// Safe mode OFF: everything passes through, including the value-bearing msgs.
	h := NewSafeModeAnteHandler(stubSafeMode{active: false}, next)
	for _, m := range append(append([]sdk.Msg{}, frozen...), allowed...) {
		called = false
		_, err := h(ctx, stubTx{msgs: []sdk.Msg{m}}, false)
		require.NoError(t, err)
		require.True(t, called, "next must run when safe mode is off")
	}

	// Safe mode ON: frozen msgs are rejected, allowed msgs still pass.
	h = NewSafeModeAnteHandler(stubSafeMode{active: true}, next)
	for _, m := range frozen {
		called = false
		_, err := h(ctx, stubTx{msgs: []sdk.Msg{m}}, false)
		require.Error(t, err, "%T must be frozen", m)
		require.False(t, called, "next must NOT run for a frozen msg")
	}
	for _, m := range allowed {
		called = false
		_, err := h(ctx, stubTx{msgs: []sdk.Msg{m}}, false)
		require.NoError(t, err, "%T must remain allowed", m)
		require.True(t, called)
	}

	// A tx mixing an allowed and a frozen msg is rejected as a whole.
	called = false
	_, err := h(ctx, stubTx{msgs: []sdk.Msg{&gravitytypes.MsgDelegateKeys{}, &gravitytypes.MsgSendToEthereum{}}}, false)
	require.Error(t, err)
	require.False(t, called)
}

// A frozen gravity message wrapped in an authz MsgExec must also be rejected in
// safe mode; an allowed message wrapped the same way must pass.
func TestSafeModeAnteHandler_AuthzExecRecursion(t *testing.T) {
	next := func(ctx sdk.Context, tx sdk.Tx, simulate bool) (sdk.Context, error) { return ctx, nil }
	ctx := sdk.Context{}
	grantee := sdk.AccAddress([]byte("grantee-aaaaaaaaaaaa"))

	h := NewSafeModeAnteHandler(stubSafeMode{active: true}, next)

	frozenExec := authz.NewMsgExec(grantee, []sdk.Msg{&gravitytypes.MsgSendToEthereum{}})
	_, err := h(ctx, stubTx{msgs: []sdk.Msg{&frozenExec}}, false)
	require.Error(t, err, "authz-wrapped frozen gravity msg must be rejected")

	allowedExec := authz.NewMsgExec(grantee, []sdk.Msg{&gravitytypes.MsgDelegateKeys{}})
	_, err = h(ctx, stubTx{msgs: []sdk.Msg{&allowedExec}}, false)
	require.NoError(t, err, "authz-wrapped allowed msg must pass")
}

// A gov v1 MsgSubmitProposal carrying a frozen gravity message is rejected at
// submission while in safe mode (gov executes embedded msgs in EndBlock,
// bypassing the ante, so the submission tx is the lever).
func TestSafeModeAnteHandler_GovSubmitProposalRecursion(t *testing.T) {
	next := func(ctx sdk.Context, tx sdk.Tx, simulate bool) (sdk.Context, error) { return ctx, nil }
	ctx := sdk.Context{}
	h := NewSafeModeAnteHandler(stubSafeMode{active: true}, next)

	frozenProp, err := govtypesv1.NewMsgSubmitProposal(
		[]sdk.Msg{&gravitytypes.MsgSendToEthereum{}}, sdk.NewCoins(), "proposer", "", "t", "s")
	require.NoError(t, err)
	_, err = h(ctx, stubTx{msgs: []sdk.Msg{frozenProp}}, false)
	require.Error(t, err, "gov proposal embedding a frozen gravity msg must be rejected at submission")

	okProp, err := govtypesv1.NewMsgSubmitProposal(
		[]sdk.Msg{&gravitytypes.MsgDelegateKeys{}}, sdk.NewCoins(), "proposer", "", "t", "s")
	require.NoError(t, err)
	_, err = h(ctx, stubTx{msgs: []sdk.Msg{okProp}}, false)
	require.NoError(t, err, "gov proposal with only allowed msgs must pass")
}

// buildNestedExec wraps an innocuous message in `depth` layers of authz.MsgExec.
func buildNestedExec(t *testing.T, depth int) sdk.Msg {
	t.Helper()

	inner := authz.NewMsgExec(sdk.AccAddress("grantee"), []sdk.Msg{&gravitytypes.MsgEthereumHeightVote{}})
	cur := &inner
	for i := 0; i < depth; i++ {
		next := authz.NewMsgExec(sdk.AccAddress("grantee"), []sdk.Msg{cur})
		cur = &next
	}
	return cur
}

// This handler is installed OUTSIDE the SDK ante chain (app.go), so it runs
// before SetUpContextDecorator installs a gas meter and before signature
// verification. Unbounded wrapper recursion would be unmetered work an
// unauthenticated sender could trigger during safe mode -- precisely when the
// chain is least able to absorb it.
//
// Past the depth cap the handler must fail closed rather than keep walking.
func TestSafeModeAnteCapsWrapperDepth(t *testing.T) {
	called := false
	next := func(ctx sdk.Context, tx sdk.Tx, simulate bool) (sdk.Context, error) {
		called = true
		return ctx, nil
	}

	deep := buildNestedExec(t, safeModeMaxMsgDepth+5)
	h := NewSafeModeAnteHandler(stubSafeMode{active: true}, next)

	_, err := h(sdk.Context{}, stubTx{msgs: []sdk.Msg{deep}}, false)
	require.Error(t, err, "wrapper nested past the cap must be rejected")
	require.False(t, called, "must not fall through to the rest of the ante chain")

	// A realistically shallow wrapper carrying nothing frozen still passes.
	called = false
	shallow := buildNestedExec(t, 1)
	_, err = h(sdk.Context{}, stubTx{msgs: []sdk.Msg{shallow}}, false)
	require.NoError(t, err, "ordinary nesting must still pass")
	require.True(t, called)

	// And the cap must not run at all when safe mode is off.
	called = false
	off := NewSafeModeAnteHandler(stubSafeMode{active: false}, next)
	_, err = off(sdk.Context{}, stubTx{msgs: []sdk.Msg{deep}}, false)
	require.NoError(t, err, "depth cap only applies while frozen")
	require.True(t, called)
}

// authz grants are frozen during safe mode.
//
// A grant outlives the freeze. A community-only validator set that passes a
// MsgGrant naming the governance module account as granter, with a
// non-expiring GenericAuthorization for poa MsgUpdateAuthoritySet /
// MsgUpdateParams or gov MsgExecLegacyContent, retains that capability
// permanently -- turning a temporary governance capture into a standing
// unilateral one, exercisable after the authority set is restored.
func TestSafeModeAnteFreezesAuthzGrants(t *testing.T) {
	called := false
	next := func(ctx sdk.Context, tx sdk.Tx, simulate bool) (sdk.Context, error) {
		called = true
		return ctx, nil
	}

	grant := &authz.MsgGrant{Granter: "somm1granter", Grantee: "somm1grantee"}

	// Frozen while safe mode is active...
	h := NewSafeModeAnteHandler(stubSafeMode{active: true}, next)
	_, err := h(sdk.Context{}, stubTx{msgs: []sdk.Msg{grant}}, false)
	require.Error(t, err)
	require.False(t, called)

	// ...including when wrapped in a gov proposal submission.
	called = false
	wrapped, err := govtypesv1.NewMsgSubmitProposal(
		[]sdk.Msg{grant}, sdk.NewCoins(), "somm1proposer", "", "t", "s",
	)
	require.NoError(t, err)
	_, err = h(sdk.Context{}, stubTx{msgs: []sdk.Msg{wrapped}}, false)
	require.Error(t, err, "a grant hidden inside a proposal submission must also be refused")
	require.False(t, called)

	// ...and unaffected when the chain is healthy.
	called = false
	off := NewSafeModeAnteHandler(stubSafeMode{active: false}, next)
	_, err = off(sdk.Context{}, stubTx{msgs: []sdk.Msg{grant}}, false)
	require.NoError(t, err)
	require.True(t, called)
}
