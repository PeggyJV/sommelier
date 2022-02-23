package keeper

import (
	"context"

	"github.com/peggyjv/sommelier/v3/x/pubsub/types"
)

var _ types.MsgServer = Keeper{}

// TODO(bolten): implement the msg server functions

func (k Keeper) AddPublisherIntent(c context.Context, msg *types.MsgAddPublisherIntent) (*types.MsgAddPublisherIntentResponse, error) {
	return nil, nil
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
