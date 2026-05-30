package app

import (
	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
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
// Only top-level messages are inspected. A frozen gravity message nested inside
// an authz MsgExec is not the relayer's normal submission path; extend
// isFrozenGravityMsg to recurse if that surface ever matters.
func NewSafeModeAnteHandler(poa PoaSafeModeReader, next sdk.AnteHandler) sdk.AnteHandler {
	return func(ctx sdk.Context, tx sdk.Tx, simulate bool) (sdk.Context, error) {
		if poa != nil && poa.SafeModeActive(ctx) {
			for _, msg := range tx.GetMsgs() {
				if isFrozenGravityMsg(msg) {
					return ctx, errorsmod.Wrapf(sdkerrors.ErrInvalidRequest,
						"x/poa safe mode active: %s is frozen until the authority set is restored", sdk.MsgTypeURL(msg))
				}
			}
		}
		return next(ctx, tx, simulate)
	}
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
