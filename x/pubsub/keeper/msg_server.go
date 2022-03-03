package keeper

import (
	"context"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/peggyjv/sommelier/v3/x/pubsub/types"
)

var _ types.MsgServer = Keeper{}

func (k Keeper) AddPublisherIntent(c context.Context, msg *types.MsgAddPublisherIntent) (*types.MsgAddPublisherIntentResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	publisherIntent := msg.PublisherIntent

	if err := publisherIntent.ValidateBasic(); err != nil {
		return nil, sdkerrors.Wrapf(types.ErrInvalid, "invalid publisher intent: %s", err.Error())
	}

	publisher, found := k.GetPublisher(ctx, publisherIntent.PublisherDomain)
	if !found {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrNotFound, "no publisher found with domain: %s", publisherIntent.PublisherDomain)
	}

	signer := msg.MustGetSigner()
	if signer.String() != publisher.Address {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrUnauthorized, "registered publisher address must be signer: %s", publisher.Address)
	}

	_, found = k.GetPublisherIntent(ctx, publisherIntent.SubscriptionId, publisherIntent.PublisherDomain)
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
				sdk.NewAttribute(types.AttributeKeySubscriptionId, publisherIntent.SubscriptionId),
				sdk.NewAttribute(types.AttributeKeyPublisherDomain, publisherIntent.PublisherDomain),
			),
		},
	)

	return &types.MsgAddPublisherIntentResponse{}, nil
}

func (k Keeper) AddSubscriberIntent(c context.Context, msg *types.MsgAddSubscriberIntent) (*types.MsgAddSubscriberIntentResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	subscriberIntent := msg.SubscriberIntent

	if err := subscriberIntent.ValidateBasic(); err != nil {
		return nil, sdkerrors.Wrapf(types.ErrInvalid, "invalid subscriber intent: %s", err.Error())
	}

	// ValidateBasic will confirm this is already correctly formatted
	subscriberAddress, _ := sdk.AccAddressFromBech32(subscriberIntent.SubscriberAddress)
	subscriber, found := k.GetSubscriber(ctx, subscriberAddress)
	if !found {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrNotFound, "no subscriber found with address: %s", subscriberIntent.SubscriberAddress)
	}

	signer := msg.MustGetSigner()
	if signer.String() != subscriber.Address {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrUnauthorized, "registered subscriber address must be signer: %s", subscriber.Address)
	}

	_, found = k.GetSubscriberIntent(ctx, subscriberIntent.SubscriptionId, subscriberAddress)
	if found {
		// TODO(bolten): would this UX be better if it removed the old one for you?
		return nil, sdkerrors.Wrapf(types.ErrAlreadyExists, "subscriber already has intent for this subscription ID, must be removed first")
	}

	// TODO(bolten): verify subscriber is authorized

	k.SetSubscriberIntent(ctx, subscriberAddress, *subscriberIntent)

	ctx.EventManager().EmitEvents(
		sdk.Events{
			sdk.NewEvent(
				sdk.EventTypeMessage,
				sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
				sdk.NewAttribute(sdk.AttributeKeyAction, types.EventTypeAddSubscriberIntent),
				sdk.NewAttribute(sdk.AttributeKeySender, msg.Signer),
			),
			sdk.NewEvent(
				types.EventTypeAddSubscriberIntent,
				sdk.NewAttribute(types.AttributeKeySubscriptionId, subscriberIntent.SubscriptionId),
				sdk.NewAttribute(types.AttributeKeySubscriberAddress, subscriberIntent.SubscriberAddress),
				sdk.NewAttribute(types.AttributeKeyPublisherDomain, subscriberIntent.PublisherDomain),
			),
		},
	)

	return &types.MsgAddSubscriberIntentResponse{}, nil
}

func (k Keeper) AddSubscriber(c context.Context, msg *types.MsgAddSubscriber) (*types.MsgAddSubscriberResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	subscriber := msg.Subscriber

	if err := subscriber.ValidateBasic(); err != nil {
		return nil, sdkerrors.Wrapf(types.ErrInvalid, "invalid subscriber: %s", err.Error())
	}

	// ValidateBasic will confirm this is already correctly formatted
	subscriberAddress, _ := sdk.AccAddressFromBech32(subscriber.Address)
	_, found := k.GetSubscriber(ctx, subscriberAddress)
	if found {
		// TODO(bolten): would this UX be better if it removed the old one for you?
		return nil, sdkerrors.Wrapf(types.ErrAlreadyExists, "subscriber already exists, must be removed first")
	}

	k.SetSubscriber(ctx, subscriberAddress, *subscriber)

	ctx.EventManager().EmitEvents(
		sdk.Events{
			sdk.NewEvent(
				sdk.EventTypeMessage,
				sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
				sdk.NewAttribute(sdk.AttributeKeyAction, types.EventTypeAddSubscriber),
				sdk.NewAttribute(sdk.AttributeKeySender, msg.Signer),
			),
			sdk.NewEvent(
				types.EventTypeAddSubscriber,
				sdk.NewAttribute(types.AttributeKeySubscriberAddress, subscriber.Address),
			),
		},
	)

	return &types.MsgAddSubscriberResponse{}, nil
}

func (k Keeper) RemovePublisherIntent(c context.Context, msg *types.MsgRemovePublisherIntent) (*types.MsgRemovePublisherIntentResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	subscriptionId := msg.SubscriptionId
	publisherDomain := msg.PublisherDomain

	if err := types.ValidateSubscriptionId(subscriptionId); err != nil {
		return nil, sdkerrors.Wrapf(types.ErrInvalid, "invalid subscription ID: %s", err.Error())
	}

	if err := types.ValidateDomain(publisherDomain); err != nil {
		return nil, sdkerrors.Wrapf(types.ErrInvalid, "invalid publisher domain: %s", err.Error())
	}

	publisher, found := k.GetPublisher(ctx, publisherDomain)
	if !found {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrNotFound, "no publisher found with domain: %s", publisherDomain)
	}

	signer := msg.MustGetSigner()
	if signer.String() != publisher.Address {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrUnauthorized, "registered publisher address must be signer: %s", publisher.Address)
	}

	_, found = k.GetPublisherIntent(ctx, publisherDomain, subscriptionId)
	if !found {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrNotFound, "no publisher intent for domain %s and subscription ID %s found", publisherDomain, subscriptionId)
	}

	k.DeletePublisherIntent(ctx, subscriptionId, publisherDomain)

	ctx.EventManager().EmitEvents(
		sdk.Events{
			sdk.NewEvent(
				sdk.EventTypeMessage,
				sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
				sdk.NewAttribute(sdk.AttributeKeyAction, types.EventTypeRemovePublisherIntent),
				sdk.NewAttribute(sdk.AttributeKeySender, msg.Signer),
			),
			sdk.NewEvent(
				types.EventTypeRemovePublisherIntent,
				sdk.NewAttribute(types.AttributeKeySubscriptionId, subscriptionId),
				sdk.NewAttribute(types.AttributeKeyPublisherDomain, publisherDomain),
			),
		},
	)

	return &types.MsgRemovePublisherIntentResponse{}, nil
}

func (k Keeper) RemoveSubscriberIntent(c context.Context, msg *types.MsgRemoveSubscriberIntent) (*types.MsgRemoveSubscriberIntentResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	subscriptionId := msg.SubscriptionId

	if err := types.ValidateSubscriptionId(subscriptionId); err != nil {
		return nil, sdkerrors.Wrapf(types.ErrInvalid, "invalid subscription ID: %s", err.Error())
	}

	if err := types.ValidateAddress(msg.SubscriberAddress); err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, fmt.Sprintf("invalid subscriber address: %s", err.Error()))
	}

	// ValidateAddress will confirm this is already correctly formatted
	subscriberAddress, _ := sdk.AccAddressFromBech32(msg.SubscriberAddress)
	subscriber, found := k.GetSubscriber(ctx, subscriberAddress)
	if !found {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrNotFound, "no subscriber found with address: %s", subscriberAddress.String())
	}

	signer := msg.MustGetSigner()
	if signer.String() != subscriber.Address {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrUnauthorized, "registered subscriber address must be signer: %s", subscriber.Address)
	}

	subscriberIntent, found := k.GetSubscriberIntent(ctx, subscriptionId, subscriberAddress)
	if !found {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrNotFound, "no subscriber intent for address %s and subscription ID %s found", subscriberAddress.String(), subscriptionId)
	}

	k.DeleteSubscriberIntent(ctx, subscriptionId, subscriberAddress, subscriberIntent.PublisherDomain)

	ctx.EventManager().EmitEvents(
		sdk.Events{
			sdk.NewEvent(
				sdk.EventTypeMessage,
				sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
				sdk.NewAttribute(sdk.AttributeKeyAction, types.EventTypeRemoveSubscriberIntent),
				sdk.NewAttribute(sdk.AttributeKeySender, msg.Signer),
			),
			sdk.NewEvent(
				types.EventTypeRemoveSubscriberIntent,
				sdk.NewAttribute(types.AttributeKeySubscriptionId, subscriptionId),
				sdk.NewAttribute(types.AttributeKeySubscriberAddress, subscriberAddress.String()),
				sdk.NewAttribute(types.AttributeKeyPublisherDomain, subscriberIntent.PublisherDomain),
			),
		},
	)

	return &types.MsgRemoveSubscriberIntentResponse{}, nil
}

func (k Keeper) RemoveSubscriber(c context.Context, msg *types.MsgRemoveSubscriber) (*types.MsgRemoveSubscriberResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)

	if err := types.ValidateAddress(msg.SubscriberAddress); err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, fmt.Sprintf("invalid subscriber address: %s", err.Error()))
	}

	// ValidateAddress will confirm this is already correctly formatted
	subscriberAddress, _ := sdk.AccAddressFromBech32(msg.SubscriberAddress)
	subscriber, found := k.GetSubscriber(ctx, subscriberAddress)
	if !found {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrNotFound, "no subscriber found with address: %s", subscriberAddress.String())
	}

	signer := msg.MustGetSigner()
	if signer.String() != subscriber.Address {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrUnauthorized, "registered subscriber address must be signer: %s", subscriber.Address)
	}

	k.DeleteSubscriber(ctx, subscriberAddress)

	ctx.EventManager().EmitEvents(
		sdk.Events{
			sdk.NewEvent(
				sdk.EventTypeMessage,
				sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
				sdk.NewAttribute(sdk.AttributeKeyAction, types.EventTypeRemoveSubscriber),
				sdk.NewAttribute(sdk.AttributeKeySender, msg.Signer),
			),
			sdk.NewEvent(
				types.EventTypeRemoveSubscriber,
				sdk.NewAttribute(types.AttributeKeySubscriptionId, subscriberAddress.String()),
			),
		},
	)

	return &types.MsgRemoveSubscriberResponse{}, nil
}

func (k Keeper) RemovePublisher(c context.Context, msg *types.MsgRemovePublisher) (*types.MsgRemovePublisherResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	publisherDomain := msg.PublisherDomain

	if err := types.ValidateDomain(publisherDomain); err != nil {
		return nil, sdkerrors.Wrapf(types.ErrInvalid, "invalid publisher domain: %s", err.Error())
	}

	publisher, found := k.GetPublisher(ctx, publisherDomain)
	if !found {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrNotFound, "no publisher found with domain: %s", publisherDomain)
	}

	signer := msg.MustGetSigner()
	if signer.String() != publisher.Address {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrUnauthorized, "registered publisher address must be signer: %s", publisher.Address)
	}

	k.DeletePublisher(ctx, publisherDomain)

	ctx.EventManager().EmitEvents(
		sdk.Events{
			sdk.NewEvent(
				sdk.EventTypeMessage,
				sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
				sdk.NewAttribute(sdk.AttributeKeyAction, types.EventTypeRemovePublisher),
				sdk.NewAttribute(sdk.AttributeKeySender, msg.Signer),
			),
			sdk.NewEvent(
				types.EventTypeRemovePublisher,
				sdk.NewAttribute(types.AttributeKeyPublisherDomain, publisherDomain),
			),
		},
	)

	return &types.MsgRemovePublisherResponse{}, nil
}
