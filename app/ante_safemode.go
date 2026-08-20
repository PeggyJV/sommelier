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
// NOTE: this closes the SUBMISSION path only. A gov v1 proposal already in
// voting when safe mode triggers is not caught here: gov executes a passed
// proposal's messages through the message router in EndBlock, which never
// reaches the ante handler.
//
// For gravity that residual is small: the only reachable message is
// MsgSendToEthereum signed by the gov module account, since the others require
// a registered orchestrator signer. In-repo cork/axelarcork msgs executed via
// gov v1 ARE gated by their own msg servers.
//
// For authz.MsgGrant the residual is NOT small -- a pre-staged proposal can
// still install a standing grant from the governance account. Closing that
// properly requires rejecting gov-account granters inside the authz msg server
// (or a wrapped authz MsgServer), which is a module-wiring change rather than
// an ante-layer one. Until then, audit outstanding grants whose granter is the
// governance module account before and after any safe-mode episode.
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

// isFrozenGravityMsg reports whether msg is frozen while safe mode is active.
//
// Gravity messages that authorize outbound fund movement, inbound minting via
// attestation, or Ethereum-side validator-set rotation. Non-fund messages
// (delegate-keys, height vote, cancel send) are intentionally left enabled.
//
// authz.MsgGrant is frozen too, for a different reason. A grant outlives the
// freeze: a community-only validator set that passes a MsgGrant naming the
// governance module account as granter, with a non-expiring
// GenericAuthorization for poa MsgUpdateAuthoritySet, poa MsgUpdateParams, or
// gov MsgExecLegacyContent, keeps that capability permanently -- converting a
// temporary governance capture into a standing, unilateral one exercisable
// long after the authority set is restored. Blocking grant creation while
// frozen is cheap; nothing legitimate needs a new grant mid-incident.
func isFrozenGravityMsg(msg sdk.Msg) bool {
	switch msg.(type) {
	case *gravitytypes.MsgSendToEthereum,
		*gravitytypes.MsgSubmitEthereumEvent,
		*gravitytypes.MsgSubmitEthereumTxConfirmation:
		return true
	case *authz.MsgGrant:
		return true
	default:
		return false
	}
}
