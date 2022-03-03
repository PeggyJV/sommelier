package keeper

import (
	"context"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/peggyjv/sommelier/v3/x/pubsub/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var _ types.QueryServer = Keeper{}

func (k Keeper) Params(c context.Context, req *types.QueryParamsRequest) (*types.QueryParamsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	return &types.QueryParamsResponse{
		Params: k.GetParams(sdk.UnwrapSDKContext(c)),
	}, nil
}

func (k Keeper) QueryPublisher(c context.Context, req *types.QueryPublisherRequest) (*types.QueryPublisherResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	if err := types.ValidateDomain(req.PublisherDomain); err != nil {
		return nil, status.Error(codes.InvalidArgument, fmt.Sprintf("invalid publisher domain: %s", err.Error()))
	}

	publisher, found := k.GetPublisher(sdk.UnwrapSDKContext(c), req.PublisherDomain)
	if !found {
		return nil, status.Error(codes.NotFound, "publisher")
	}

	return &types.QueryPublisherResponse{
		Publisher: &publisher,
	}, nil
}

func (k Keeper) QueryPublishers(c context.Context, _ *types.QueryPublishersRequest) (*types.QueryPublishersResponse, error) {
	return &types.QueryPublishersResponse{
		Publishers: k.GetPublishers(sdk.UnwrapSDKContext(c)),
	}, nil
}

func (k Keeper) QuerySubscriber(c context.Context, req *types.QuerySubscriberRequest) (*types.QuerySubscriberResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	addr, err := sdk.AccAddressFromBech32(req.SubscriberAddress)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, fmt.Sprintf("invalid subscriber address :%s", err.Error()))
	}

	subscriber, found := k.GetSubscriber(sdk.UnwrapSDKContext(c), addr)
	if !found {
		return nil, status.Error(codes.NotFound, "subscriber")
	}

	return &types.QuerySubscriberResponse{
		Subscriber: &subscriber,
	}, nil
}

func (k Keeper) QuerySubscribers(c context.Context, _ *types.QuerySubscribersRequest) (*types.QuerySubscribersResponse, error) {
	return &types.QuerySubscribersResponse{
		Subscribers: k.GetSubscribers(sdk.UnwrapSDKContext(c)),
	}, nil
}

func (k Keeper) QueryPublisherIntent(c context.Context, req *types.QueryPublisherIntentRequest) (*types.QueryPublisherIntentResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	if err := types.ValidateDomain(req.PublisherDomain); err != nil {
		return nil, status.Error(codes.InvalidArgument, fmt.Sprintf("invalid publisher domain: %s", err.Error()))
	}

	if err := types.ValidateSubscriptionId(req.SubscriptionId); err != nil {
		return nil, status.Error(codes.InvalidArgument, fmt.Sprintf("invalid subscription ID: %s", err.Error()))
	}

	publisherIntent, found := k.GetPublisherIntent(sdk.UnwrapSDKContext(c), req.PublisherDomain, req.SubscriptionId)
	if !found {
		return nil, status.Error(codes.NotFound, "publisher intent")
	}

	return &types.QueryPublisherIntentResponse{
		PublisherIntent: &publisherIntent,
	}, nil
}

func (k Keeper) QueryPublisherIntents(c context.Context, req *types.QueryPublisherIntentsRequest) (*types.QueryPublisherIntentsResponse, error) {
	return &types.QueryPublisherIntentsResponse{
		PublisherIntents: k.GetPublisherIntents(sdk.UnwrapSDKContext(c)),
	}, nil
}

func (k Keeper) QueryPublisherIntentsByPublisherDomain(c context.Context, req *types.QueryPublisherIntentsByPublisherDomainRequest) (*types.QueryPublisherIntentsByPublisherDomainResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	if err := types.ValidateDomain(req.PublisherDomain); err != nil {
		return nil, status.Error(codes.InvalidArgument, fmt.Sprintf("invalid publisher domain: %s", err.Error()))
	}

	return &types.QueryPublisherIntentsByPublisherDomainResponse{
		PublisherIntents: k.GetPublisherIntentsByPublisherDomain(sdk.UnwrapSDKContext(c), req.PublisherDomain),
	}, nil
}

func (k Keeper) QueryPublisherIntentsBySubscriptionId(c context.Context, req *types.QueryPublisherIntentsBySubscriptionIdRequest) (*types.QueryPublisherIntentsBySubscriptionIdResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	if err := types.ValidateSubscriptionId(req.SubscriptionId); err != nil {
		return nil, status.Error(codes.InvalidArgument, fmt.Sprintf("invalid subscription ID: %s", err.Error()))
	}

	return &types.QueryPublisherIntentsBySubscriptionIdResponse{
		PublisherIntents: k.GetPublisherIntentsBySubscriptionId(sdk.UnwrapSDKContext(c), req.SubscriptionId),
	}, nil
}

func (k Keeper) QuerySubscriberIntent(c context.Context, req *types.QuerySubscriberIntentRequest) (*types.QuerySubscriberIntentResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	addr, err := sdk.AccAddressFromBech32(req.SubscriberAddress)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, fmt.Sprintf("invalid subscriber address :%s", err.Error()))
	}

	if err := types.ValidateSubscriptionId(req.SubscriptionId); err != nil {
		return nil, status.Error(codes.InvalidArgument, fmt.Sprintf("invalid subscription ID: %s", err.Error()))
	}

	subscriberIntent, found := k.GetSubscriberIntent(sdk.UnwrapSDKContext(c), addr, req.SubscriptionId)
	if !found {
		return nil, status.Error(codes.NotFound, "subscriber intent")
	}

	return &types.QuerySubscriberIntentResponse{
		SubscriberIntent: &subscriberIntent,
	}, nil
}

func (k Keeper) QuerySubscriberIntents(c context.Context, _ *types.QuerySubscriberIntentsRequest) (*types.QuerySubscriberIntentsResponse, error) {
	return &types.QuerySubscriberIntentsResponse{
		SubscriberIntents: k.GetSubscriberIntents(sdk.UnwrapSDKContext(c)),
	}, nil
}

func (k Keeper) QuerySubscriberIntentsBySubscriberAddress(c context.Context, req *types.QuerySubscriberIntentsBySubscriberAddressRequest) (*types.QuerySubscriberIntentsBySubscriberAddressResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	addr, err := sdk.AccAddressFromBech32(req.SubscriberAddress)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, fmt.Sprintf("invalid subscriber address :%s", err.Error()))
	}

	return &types.QuerySubscriberIntentsBySubscriberAddressResponse{
		SubscriberIntents: k.GetSubscriberIntentsBySubscriberAddress(sdk.UnwrapSDKContext(c), addr),
	}, nil
}

func (k Keeper) QuerySubscriberIntentsBySubscriptionId(c context.Context, req *types.QuerySubscriberIntentsBySubscriptionIdRequest) (*types.QuerySubscriberIntentsBySubscriptionIdResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	if err := types.ValidateSubscriptionId(req.SubscriptionId); err != nil {
		return nil, status.Error(codes.InvalidArgument, fmt.Sprintf("invalid subscription ID: %s", err.Error()))
	}

	return &types.QuerySubscriberIntentsBySubscriptionIdResponse{
		SubscriberIntents: k.GetSubscriberIntentsBySubscriptionId(sdk.UnwrapSDKContext(c), req.SubscriptionId),
	}, nil
}
