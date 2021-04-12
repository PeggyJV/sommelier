package allocation

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/peggyjv/sommelier/x/allocation/keeper"
	"github.com/peggyjv/sommelier/x/allocation/types"
)

// NewHandler returns a handler for "allocation" type messages
func NewHandler(k keeper.Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) (*sdk.Result, error) {
		ctx = ctx.WithEventManager(sdk.NewEventManager())
		switch msg := msg.(type) {
		case *types.MsgDelegateDecisions:
			res, err := k.DelegateDecisions(sdk.WrapSDKContext(ctx), msg)
			return sdk.WrapServiceResult(ctx, res, err)
		case *types.MsgDecisionPrecommit:
			res, err := k.DecisionPrecommit(sdk.WrapSDKContext(ctx), msg)
			return sdk.WrapServiceResult(ctx, res, err)
		case *types.MsgDecisionCommit:
			res, err := k.DecisionCommit(sdk.WrapSDKContext(ctx), msg)
			return sdk.WrapServiceResult(ctx, res, err)
		default:
			return nil, sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "unrecognized allocation message type: %T", msg)
		}
	}
}