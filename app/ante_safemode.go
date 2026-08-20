package app

import (
	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/x/authz"
	govtypesv1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1"
	gravitytypes "github.com/peggyjv/gravity-bridge/module/v6/x/gravity/types"
)

// PoaSafeModeReader is the read-only view of x/poa safe mode used by the ante
// handler to freeze value-bearing messages.
type PoaSafeModeReader interface {
	SafeModeActive(ctx sdk.Context) bool
}

// NewSafeModeAnteHandler wraps `next` so that, while x/poa is in authority-empty
// safe mode, gravity-bridge messages that move funds or rotate the Ethereum
// validator set are rejected. Gravity is an external dependency and cannot gate
// its own handlers, so the ante handler is the freeze lever for it; the in-repo
// cork and axelarcork modules instead gate themselves in their msg servers and
// EndBlockers.
//
// Inspection recurses into message wrappers (authz MsgExec, gov v1
// MsgSubmitProposal) — see containsFrozenGravityMsg.
func NewSafeModeAnteHandler(poa PoaSafeModeReader, next sdk.AnteHandler) sdk.AnteHandler {
	return func(ctx sdk.Context, tx sdk.Tx, simulate bool) (sdk.Context, error) {
		if poa != nil && poa.SafeModeActive(ctx) {
			for _, msg := range tx.GetMsgs() {
				if containsFrozenGravityMsg(msg) {
					return ctx, errorsmod.Wrapf(sdkerrors.ErrInvalidRequest,
						"x/poa safe mode active: %s is frozen until the authority set is restored", sdk.MsgTypeURL(msg))
				}
			}
		}
		return next(ctx, tx, simulate)
	}
}

// containsFrozenGravityMsg reports whether msg is a frozen gravity message or a
// wrapper that carries one (recursively): an authz MsgExec, or a gov v1
// MsgSubmitProposal. Gov v1 executes a proposal's embedded messages through the
// message router in gov's EndBlock, bypassing the ante handler, so blocking the
// submission tx is the lever available here (the gov keeper takes a concrete
// MsgServiceRouter that cannot be wrapped). A wrapper whose inner messages
// cannot be decoded is treated as frozen (fail-closed) while in safe mode.
//
// NOTE: this closes the submission path. A gov v1 proposal already in voting
// when safe mode triggers is not caught here; its only reachable gravity msg is
// MsgSendToEthereum signed by the gov module account (others require a
// registered orchestrator signer), which is a negligible residual. In-repo
// cork/axelarcork msgs executed via gov v1 ARE gated by their msg servers.
func containsFrozenGravityMsg(msg sdk.Msg) bool {
	return containsFrozenGravityMsgAtDepth(msg, 0)
}

// safeModeMaxMsgDepth caps wrapper recursion.
//
// This handler is installed OUTSIDE the SDK ante chain (see app.go), so it runs
// before SetUpContextDecorator installs a gas meter and before signature
// verification. Unbounded recursion there is unmetered work an unauthenticated
// sender can trigger during safe mode -- exactly when the chain is least able
// to absorb it -- and each level re-unpacks every Any it contains, so cost
// grows faster than tx size.
//
// Real txs nest at most a couple of levels (authz MsgExec around a gov
// proposal). Anything deeper is treated as frozen rather than walked, matching
// the fail-closed handling of an undecodable wrapper below.
const safeModeMaxMsgDepth = 8

func containsFrozenGravityMsgAtDepth(msg sdk.Msg, depth int) bool {
	if isFrozenGravityMsg(msg) {
		return true
	}
	if depth >= safeModeMaxMsgDepth {
		// Fail closed: refuse rather than recurse further.
		return true
	}
	switch m := msg.(type) {
	case *authz.MsgExec:
		inner, err := m.GetMessages()
		if err != nil {
			return true
		}
		for _, im := range inner {
			if containsFrozenGravityMsgAtDepth(im, depth+1) {
				return true
			}
		}
	case *govtypesv1.MsgSubmitProposal:
		inner, err := m.GetMsgs()
		if err != nil {
			return true
		}
		for _, im := range inner {
			if containsFrozenGravityMsgAtDepth(im, depth+1) {
				return true
			}
		}
	}
	return false
}

// isFrozenGravityMsg reports whether msg is a gravity message that authorizes
// outbound fund movement, inbound minting via attestation, or Ethereum-side
// validator-set rotation. Non-fund messages (delegate-keys, height vote, cancel
// send) are intentionally left enabled.
func isFrozenGravityMsg(msg sdk.Msg) bool {
	switch msg.(type) {
	case *gravitytypes.MsgSendToEthereum,
		*gravitytypes.MsgSubmitEthereumEvent,
		*gravitytypes.MsgSubmitEthereumTxConfirmation:
		return true
	default:
		return false
	}
}
