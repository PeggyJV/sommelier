package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/peggyjv/sommelier/v3/x/pubsub/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var _ types.QueryServer = Keeper{}

func (k Keeper) Params(c context.Context, req *types.QueryParamsRequest) (*types.QueryParamsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	return &types.QueryParamsResponse{
		Params: k.GetParams(sdk.UnwrapSDKContext(c)),
	}, nil
}

func (k Keeper) QueryPublisher(c context.Context, req *types.QueryPublisherRequest) (*types.QueryPublisherResponse, error) {
	return nil, nil
}

func (k Keeper) QueryPublishers(c context.Context, req *types.QueryPublishersRequest) (*types.QueryPublishersResponse, error) {
	return nil, nil
}

func (k Keeper) QuerySubscriber(c context.Context, req *types.QuerySubscriberRequest) (*types.QuerySubscriberResponse, error) {
	return nil, nil
}

func (k Keeper) QuerySubscribers(c context.Context, req *types.QuerySubscribersRequest) (*types.QuerySubscribersResponse, error) {
	return nil, nil
}

func (k Keeper) QueryPublisherIntent(c context.Context, req *types.QueryPublisherIntentRequest) (*types.QueryPublisherIntentResponse, error) {
	return nil, nil
}

func (k Keeper) QueryPublisherIntents(c context.Context, req *types.QueryPublisherIntentsRequest) (*types.QueryPublisherIntentsResponse, error) {
	return nil, nil
}

func (k Keeper) QueryPublisherIntentsByPublisherDomain(c context.Context, req *types.QueryPublisherIntentsByPublisherDomainRequest) (*types.QueryPublisherIntentsByPublisherDomainResponse, error) {
	return nil, nil
}

func (k Keeper) QueryPublisherIntentsBySubscriptionId(c context.Context, req *types.QueryPublisherIntentsBySubscriptionIdRequest) (*types.QueryPublisherIntentsBySubscriptionIdResponse, error) {
	return nil, nil
}

func (k Keeper) QuerySubscriberIntent(c context.Context, req *types.QuerySubscriberIntentRequest) (*types.QuerySubscriberIntentResponse, error) {
	return nil, nil
}

func (k Keeper) QuerySubscriberIntents(c context.Context, req *types.QuerySubscriberIntentsRequest) (*types.QuerySubscriberIntentsResponse, error) {
	return nil, nil
}

func (k Keeper) QuerySubscriberIntentsBySubscriberAddress(c context.Context, req *types.QuerySubscriberIntentsBySubscriberAddressRequest) (*types.QuerySubscriberIntentsBySubscriberAddressResponse, error) {
	return nil, nil
}

func (k Keeper) QuerySubscriberIntentsBySubscriptionId(c context.Context, req *types.QuerySubscriberIntentsBySubscriptionIdRequest) (*types.QuerySubscriberIntentsBySubscriptionIdResponse, error) {
	return nil, nil
}
