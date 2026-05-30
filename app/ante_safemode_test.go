package app

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/authz"
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
