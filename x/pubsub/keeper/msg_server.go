package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/peggyjv/sommelier/v3/x/pubsub/types"
)

var _ types.MsgServer = Keeper{}

// TODO(bolten): implement the msg server functions

func (k Keeper) AddPublisherIntent(c context.Context, msg *types.MsgAddPublisherIntent) (*types.MsgAddPublisherIntentResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	publisherIntent := msg.PublisherIntent

	publisher, found := k.GetPublisher(ctx, publisherIntent.PublisherDomain)
	if !found {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrNotFound, "no publisher found with domain: %s", publisherIntent.PublisherDomain)
	}

	signer := msg.MustGetSigner()
	if signer.String() != publisher.Address {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrUnauthorized, "registered publisher address must be signer: %s", publisher.Address)
	}

	_, found = k.GetPublisherIntent(ctx, publisherIntent.PublisherDomain, publisherIntent.SubscriptionId)
	if found {
		// TODO(bolten): would this UX be better if it removed the old one for you? deleting a publisher intent will clear out
		// any subscriber intents attached to it
		return nil, sdkerrors.Wrapf(types.ErrAlreadyExists, "publisher already has intent for this subscription ID, must be removed first")
	}

	k.SetPublisherIntent(ctx, *publisherIntent)

	ctx.EventManager().EmitEvents(
		sdk.Events{
			sdk.NewEvent(
				sdk.EventTypeMessage,
				sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
				sdk.NewAttribute(sdk.AttributeKeyAction, types.EventTypeAddPublisherIntent),
				sdk.NewAttribute(sdk.AttributeKeySender, msg.Signer),
			),
			sdk.NewEvent(
				types.EventTypeAddPublisherIntent,
				sdk.NewAttribute(types.AttributeKeyPublisherDomain, publisherIntent.PublisherDomain),
				sdk.NewAttribute(types.AttributeKeySubscriptionId, publisherIntent.SubscriptionId),
			),
		},
	)

	return &types.MsgAddPublisherIntentResponse{}, nil
}

func (k Keeper) AddSubscriberIntent(c context.Context, msg *types.MsgAddSubscriberIntent) (*types.MsgAddSubscriberIntentResponse, error) {
	return nil, nil
}

func (k Keeper) AddSubscriber(c context.Context, msg *types.MsgAddSubscriber) (*types.MsgAddSubscriberResponse, error) {
	return nil, nil
}

func (k Keeper) RemovePublisherIntent(c context.Context, msg *types.MsgRemovePublisherIntent) (*types.MsgRemovePublisherIntentResponse, error) {
	return nil, nil
}

func (k Keeper) RemoveSubscriberIntent(c context.Context, msg *types.MsgRemoveSubscriberIntent) (*types.MsgRemoveSubscriberIntentResponse, error) {
	return nil, nil
}

func (k Keeper) RemoveSubscriber(c context.Context, msg *types.MsgRemoveSubscriber) (*types.MsgRemoveSubscriberResponse, error) {
	return nil, nil
}

func (k Keeper) RemovePublisher(c context.Context, msg *types.MsgRemovePublisher) (*types.MsgRemovePublisherResponse, error) {
	return nil, nil
}
