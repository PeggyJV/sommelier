package pubsub

import (
	"fmt"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	govtypesv1beta1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1beta1"
	"github.com/peggyjv/sommelier/v7/x/pubsub/keeper"
	"github.com/peggyjv/sommelier/v7/x/pubsub/types"
)

// NewHandler returns a handler for "pubsub" type messages
func NewHandler(k keeper.Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) (*sdk.Result, error) {
		ctx = ctx.WithEventManager(sdk.NewEventManager())

		switch msg := msg.(type) {
		case *types.MsgRemovePublisherRequest:
			res, err := k.RemovePublisher(sdk.WrapSDKContext(ctx), msg)
			return sdk.WrapServiceResult(ctx, res, err)
		case *types.MsgAddSubscriberRequest:
			res, err := k.AddSubscriber(sdk.WrapSDKContext(ctx), msg)
			return sdk.WrapServiceResult(ctx, res, err)
		case *types.MsgRemoveSubscriberRequest:
			res, err := k.RemoveSubscriber(sdk.WrapSDKContext(ctx), msg)
			return sdk.WrapServiceResult(ctx, res, err)
		case *types.MsgAddPublisherIntentRequest:
			res, err := k.AddPublisherIntent(sdk.WrapSDKContext(ctx), msg)
			return sdk.WrapServiceResult(ctx, res, err)
		case *types.MsgRemovePublisherIntentRequest:
			res, err := k.RemovePublisherIntent(sdk.WrapSDKContext(ctx), msg)
			return sdk.WrapServiceResult(ctx, res, err)
		case *types.MsgAddSubscriberIntentRequest:
			res, err := k.AddSubscriberIntent(sdk.WrapSDKContext(ctx), msg)
			return sdk.WrapServiceResult(ctx, res, err)
		case *types.MsgRemoveSubscriberIntentRequest:
			res, err := k.RemoveSubscriberIntent(sdk.WrapSDKContext(ctx), msg)
			return sdk.WrapServiceResult(ctx, res, err)
		default:
			errMsg := fmt.Sprintf("unrecognized %s message type: %T", types.ModuleName, msg)
			return nil, errorsmod.Wrap(sdkerrors.ErrUnknownRequest, errMsg)
		}
	}
}

// NewPubsubProposalHandler returns a handler for "pubsub" governance proposals
func NewPubsubProposalHandler(k keeper.Keeper) govtypesv1beta1.Handler {
	return func(ctx sdk.Context, content govtypesv1beta1.Content) error {
		switch c := content.(type) {
		case *types.AddPublisherProposal:
			return keeper.HandleAddPublisherProposal(ctx, k, *c)
		case *types.RemovePublisherProposal:
			return keeper.HandleRemovePublisherProposal(ctx, k, *c)
		case *types.AddDefaultSubscriptionProposal:
			return keeper.HandleAddDefaultSubscriptionProposal(ctx, k, *c)
		case *types.RemoveDefaultSubscriptionProposal:
			return keeper.HandleRemoveDefaultSubscriptionProposal(ctx, k, *c)
		default:
			errMsg := fmt.Sprintf("unrecognized %s proposal type: %T", types.ModuleName, c)
			return errorsmod.Wrap(sdkerrors.ErrUnknownRequest, errMsg)
		}
	}
}
