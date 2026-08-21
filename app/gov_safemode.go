package app

import (
	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	govtypesv1beta1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1beta1"
)

// freezeGovHandlerInSafeMode wraps a legacy gov proposal handler so that, while
// x/poa is in authority-empty safe mode, the proposal is rejected. It is used
// for value-bearing routes whose handlers run inside gov EndBlock and therefore
// bypass the ante handler and the module msg servers — notably gravity's
// community-pool Ethereum spend. In-repo cork/axelarcork gate their own
// value-bearing proposals directly in the keeper handlers.
func freezeGovHandlerInSafeMode(poa PoaSafeModeReader, next govtypesv1beta1.Handler) govtypesv1beta1.Handler {
	return func(ctx sdk.Context, content govtypesv1beta1.Content) error {
		if poa != nil && poa.SafeModeActive(ctx) {
			return errorsmod.Wrap(sdkerrors.ErrInvalidRequest,
				"x/poa safe mode active: this proposal is frozen until the authority set is restored")
		}
		return next(ctx, content)
	}
}
