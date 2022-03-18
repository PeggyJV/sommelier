package cellarfees

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/peggyjv/sommelier/v3/x/cellarfees/keeper"
)

// NewHandler returns a handler for "cellarfees" type messages.
func NewHandler(k keeper.Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) (*sdk.Result, error) {
		ctx = ctx.WithEventManager(sdk.NewEventManager())
		return nil, sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "unrecognized cellarfees message type: %T", msg)
	}
}
