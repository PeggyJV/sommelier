package cellarfees

import (
	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/peggyjv/sommelier/v8/x/cellarfees/keeper"
)

// NewHandler returns a handler for "cellarfees" type messages.
func NewHandler(k keeper.Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) (*sdk.Result, error) {
		return nil, errorsmod.Wrapf(sdkerrors.ErrUnknownRequest, "unrecognized cellarfees message type: %T", msg)
	}
}
