package pubsub

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	"github.com/peggyjv/sommelier/v3/x/pubsub/keeper"
	"github.com/peggyjv/sommelier/v3/x/pubsub/types"
)

// NewHandler returns a handler for "pubsub" type messages
func NewHandler(k keeper.Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) (*sdk.Result, error) {
		ctx = ctx.WithEventManager(sdk.NewEventManager())

		switch msg := msg.(type) {
		case *types.MsgRemovePublisher:
			res, err := k.RemovePublisher(sdk.WrapSDKContext(ctx), msg)
			return sdk.WrapServiceResult(ctx, res, err)
		case *types.MsgAddSubscriber:
			res, err := k.AddSubscriber(sdk.WrapSDKContext(ctx), msg)
			return sdk.WrapServiceResult(ctx, res, err)
		case *types.MsgRemoveSubscriber:
			res, err := k.RemoveSubscriber(sdk.WrapSDKContext(ctx), msg)
			return sdk.WrapServiceResult(ctx, res, err)
		case *types.MsgAddPublisherIntent:
			res, err := k.AddPublisherIntent(sdk.WrapSDKContext(ctx), msg)
			return sdk.WrapServiceResult(ctx, res, err)
		case *types.MsgRemovePublisherIntent:
			res, err := k.RemovePublisherIntent(sdk.WrapSDKContext(ctx), msg)
			return sdk.WrapServiceResult(ctx, res, err)
		case *types.MsgAddSubscriberIntent:
			res, err := k.AddSubscriberIntent(sdk.WrapSDKContext(ctx), msg)
			return sdk.WrapServiceResult(ctx, res, err)
		case *types.MsgRemoveSubscriberIntent:
			res, err := k.RemoveSubscriberIntent(sdk.WrapSDKContext(ctx), msg)
			return sdk.WrapServiceResult(ctx, res, err)
		default:
			errMsg := fmt.Sprintf("unrecognized %s message type: %T", types.ModuleName, msg)
			return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, errMsg)
		}
	}
}

// NewUpdatePublishersProposalHandler returns a handler for "pubsub" governance proposals
func NewUpdatePublishersProposalHandler(k keeper.Keeper) govtypes.Handler {
	return func(ctx sdk.Context, content govtypes.Content) error {
		switch c := content.(type) {
		case *types.AddPublisherProposal:
			return keeper.HandleAddPublisherProposal(ctx, k, *c)
		case *types.RemovePublisherProposal:
			return keeper.HandleRemovePublisherProposal(ctx, k, *c)
		default:
			errMsg := fmt.Sprintf("unrecognized %s proposal type: %T", types.ModuleName, c)
			return sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, errMsg)
		}
	}
}
