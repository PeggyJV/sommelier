package keeper

import (
	"context"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/peggyjv/sommelier/v4/x/pubsub/types"
)

var _ types.MsgServer = Keeper{}

func (k Keeper) AddPublisherIntent(c context.Context, msg *types.MsgAddPublisherIntentRequest) (*types.MsgAddPublisherIntentResponse, error) {
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
				sdk.NewAttribute(types.AttributeKeySubscriptionID, publisherIntent.SubscriptionId),
				sdk.NewAttribute(types.AttributeKeyPublisherDomain, publisherIntent.PublisherDomain),
			),
		},
	)

	return &types.MsgAddPublisherIntentResponse{}, nil
}

func (k Keeper) AddSubscriberIntent(c context.Context, msg *types.MsgAddSubscriberIntentRequest) (*types.MsgAddSubscriberIntentResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	subscriberIntent := msg.SubscriberIntent

	if err := subscriberIntent.ValidateBasic(); err != nil {
		return nil, sdkerrors.Wrapf(types.ErrInvalid, "invalid subscriber intent: %s", err.Error())
	}

	subscriptionId := subscriberIntent.SubscriptionId
	subscriberAddress := subscriberIntent.SubscriberAddress
	publisherDomain := subscriberIntent.PublisherDomain

	// ValidateBasic will confirm this is already correctly formatted
	subscriberAccAddress, _ := sdk.AccAddressFromBech32(subscriberIntent.SubscriberAddress)
	subscriber, found := k.GetSubscriber(ctx, subscriberAccAddress)
	if !found {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrNotFound, "no subscriber found with address: %s", subscriberAddress)
	}

	signer := msg.MustGetSigner()
	signerAddress := signer.String()
	if signerAddress != subscriber.Address {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrUnauthorized, "registered subscriber address must be signer: %s", subscriber.Address)
	}

	_, found = k.GetSubscriberIntent(ctx, subscriberIntent.SubscriptionId, subscriberAccAddress)
	if found {
		return nil, sdkerrors.Wrapf(types.ErrAlreadyExists, "subscriber already has intent for this subscription ID, must be removed first")
	}

	publisherIntent, found := k.GetPublisherIntent(ctx, subscriptionId, publisherDomain)
	if !found {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrNotFound, "no publisher intent for domain %s and subscription ID %s found", publisherDomain, subscriptionId)
	}

	if publisherIntent.AllowedSubscribers == types.AllowedSubscribers_VALIDATORS {
		// TODO(bolten): this implementation ends up making the module sort of non-generic but is necessary to allow
		// orchestrator keys to manipulate subscriptions
		var validatorI stakingtypes.ValidatorI
		if validator := k.gravityKeeper.GetOrchestratorValidatorAddress(ctx, signer); validator == nil {
			validatorI = k.stakingKeeper.Validator(ctx, sdk.ValAddress(signer))
		} else {
			validatorI = k.stakingKeeper.Validator(ctx, validator)
		}

		if validatorI == nil {
			return nil, sdkerrors.Wrapf(sdkerrors.ErrUnauthorized, "publisher intent requires subscriber be a validator")
		}
	} else if publisherIntent.AllowedSubscribers == types.AllowedSubscribers_LIST {
		found := false
		for _, allowedAddress := range publisherIntent.AllowedAddresses {
			if allowedAddress == signerAddress {
				found = true
				break
			}
		}

		if !found {
			return nil, sdkerrors.Wrapf(sdkerrors.ErrUnauthorized, "publisher intent requires subscriber to be in authorized list")
		}
	}

	k.SetSubscriberIntent(ctx, subscriberAccAddress, *subscriberIntent)

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
				sdk.NewAttribute(types.AttributeKeySubscriptionID, subscriberIntent.SubscriptionId),
				sdk.NewAttribute(types.AttributeKeySubscriberAddress, subscriberIntent.SubscriberAddress),
				sdk.NewAttribute(types.AttributeKeyPublisherDomain, subscriberIntent.PublisherDomain),
			),
		},
	)

	return &types.MsgAddSubscriberIntentResponse{}, nil
}

func (k Keeper) AddSubscriber(c context.Context, msg *types.MsgAddSubscriberRequest) (*types.MsgAddSubscriberResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	subscriber := msg.Subscriber

	if err := subscriber.ValidateBasic(); err != nil {
		return nil, sdkerrors.Wrapf(types.ErrInvalid, "invalid subscriber: %s", err.Error())
	}

	// ValidateBasic will confirm this is already correctly formatted
	subscriberAccAddress, _ := sdk.AccAddressFromBech32(subscriber.Address)
	_, found := k.GetSubscriber(ctx, subscriberAccAddress)
	if found {
		return nil, sdkerrors.Wrapf(types.ErrAlreadyExists, "subscriber already exists, must be removed first")
	}

	signer := msg.MustGetSigner()
	if signer.String() != subscriber.Address {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrUnauthorized, "subscriber address must be signer: %s", subscriber.Address)
	}

	k.SetSubscriber(ctx, subscriberAccAddress, *subscriber)

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

func (k Keeper) RemovePublisherIntent(c context.Context, msg *types.MsgRemovePublisherIntentRequest) (*types.MsgRemovePublisherIntentResponse, error) {
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

	_, found = k.GetPublisherIntent(ctx, subscriptionId, publisherDomain)
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
				sdk.NewAttribute(types.AttributeKeySubscriptionID, subscriptionId),
				sdk.NewAttribute(types.AttributeKeyPublisherDomain, publisherDomain),
			),
		},
	)

	return &types.MsgRemovePublisherIntentResponse{}, nil
}

func (k Keeper) RemoveSubscriberIntent(c context.Context, msg *types.MsgRemoveSubscriberIntentRequest) (*types.MsgRemoveSubscriberIntentResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	subscriptionId := msg.SubscriptionId
	subscriberAddress := msg.SubscriberAddress

	if err := types.ValidateSubscriptionId(subscriptionId); err != nil {
		return nil, sdkerrors.Wrapf(types.ErrInvalid, "invalid subscription ID: %s", err.Error())
	}

	if err := types.ValidateAddress(subscriberAddress); err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, fmt.Sprintf("invalid subscriber address: %s", err.Error()))
	}

	// ValidateAddress will confirm this is already correctly formatted
	subscriberAccAddress, _ := sdk.AccAddressFromBech32(subscriberAddress)
	subscriber, found := k.GetSubscriber(ctx, subscriberAccAddress)
	if !found {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrNotFound, "no subscriber found with address: %s", subscriberAddress)
	}

	signer := msg.MustGetSigner()
	if signer.String() != subscriber.Address {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrUnauthorized, "registered subscriber address must be signer: %s", subscriber.Address)
	}

	subscriberIntent, found := k.GetSubscriberIntent(ctx, subscriptionId, subscriberAccAddress)
	if !found {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrNotFound, "no subscriber intent for address %s and subscription ID %s found", subscriberAddress, subscriptionId)
	}

	k.DeleteSubscriberIntent(ctx, subscriptionId, subscriberAccAddress, subscriberIntent.PublisherDomain)

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
				sdk.NewAttribute(types.AttributeKeySubscriptionID, subscriptionId),
				sdk.NewAttribute(types.AttributeKeySubscriberAddress, subscriberAddress),
				sdk.NewAttribute(types.AttributeKeyPublisherDomain, subscriberIntent.PublisherDomain),
			),
		},
	)

	return &types.MsgRemoveSubscriberIntentResponse{}, nil
}

func (k Keeper) RemoveSubscriber(c context.Context, msg *types.MsgRemoveSubscriberRequest) (*types.MsgRemoveSubscriberResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	subscriberAddress := msg.SubscriberAddress

	if err := types.ValidateAddress(subscriberAddress); err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, fmt.Sprintf("invalid subscriber address: %s", err.Error()))
	}

	// ValidateAddress will confirm this is already correctly formatted
	subscriberAccAddress, _ := sdk.AccAddressFromBech32(subscriberAddress)
	subscriber, found := k.GetSubscriber(ctx, subscriberAccAddress)
	if !found {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrNotFound, "no subscriber found with address: %s", subscriberAddress)
	}

	signer := msg.MustGetSigner()
	if signer.String() != subscriber.Address {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrUnauthorized, "registered subscriber address must be signer: %s", subscriber.Address)
	}

	k.DeleteSubscriber(ctx, subscriberAccAddress)

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
				sdk.NewAttribute(types.AttributeKeySubscriptionID, subscriberAddress),
			),
		},
	)

	return &types.MsgRemoveSubscriberResponse{}, nil
}

func (k Keeper) RemovePublisher(c context.Context, msg *types.MsgRemovePublisherRequest) (*types.MsgRemovePublisherResponse, error) {
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
