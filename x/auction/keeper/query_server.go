package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/peggyjv/sommelier/v4/x/auction/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var _ types.QueryServer = Keeper{}

// QueryParams implements QueryServer
func (k Keeper) QueryParams(c context.Context, _ *types.QueryParamsRequest) (*types.QueryParamsResponse, error) {
	return &types.QueryParamsResponse{
		Params: k.GetParamSet(sdk.UnwrapSDKContext(c)),
	}, nil
}

// QueryCurrentAuction implements QueryServer
func (k Keeper) QueryCurrentAuction(c context.Context, request *types.QueryCurrentAuctionRequest) (*types.QueryCurrentAuctionResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)

	// Check active auctions
	activeAuction, found := k.GetActiveAuctionByID(ctx, request.GetAuctionId())

	if found {
		return &types.QueryCurrentAuctionResponse{Auction: &activeAuction}, nil
	}

	return &types.QueryCurrentAuctionResponse{}, status.Error(codes.NotFound, "No ongoing auction found for given id")
}

// QueryEndedAuction implements QueryServer
func (k Keeper) QueryEndedAuction(c context.Context, request *types.QueryEndedAuctionRequest) (*types.QueryEndedAuctionResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)

	// Check ended auctions first
	endedAuction, found := k.GetEndedAuctionByID(ctx, request.GetAuctionId())

	if found {
		return &types.QueryEndedAuctionResponse{Auction: &endedAuction}, nil
	}

	return &types.QueryEndedAuctionResponse{}, status.Error(codes.NotFound, "No completed auction found for given id")
}

// QueryCurrentAuctions implements QueryServer
func (k Keeper) QueryCurrentAuctions(c context.Context, _ *types.QueryCurrentAuctionsRequest) (*types.QueryCurrentAuctionsResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)

	auctions := k.GetActiveAuctions(ctx)

	if len(auctions) == 0 {
		return &types.QueryCurrentAuctionsResponse{}, status.Error(codes.NotFound, "No active auctions found")
	}

	return &types.QueryCurrentAuctionsResponse{Auctions: auctions}, nil
}

// QueryEndedAuctions implements QueryServer
func (k Keeper) QueryEndedAuctions(c context.Context, _ *types.QueryEndedAuctionsRequest) (*types.QueryEndedAuctionsResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)

	auctions := k.GetEndedAuctions(ctx)

	if len(auctions) == 0 {
		return &types.QueryEndedAuctionsResponse{}, status.Error(codes.NotFound, "No ended auctions found")
	}

	return &types.QueryEndedAuctionsResponse{Auctions: auctions}, nil
}

// QueryBid implements QueryServer
func (k Keeper) QueryBid(c context.Context, request *types.QueryBidRequest) (*types.QueryBidResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)

	bid, found := k.GetBid(ctx, request.GetAuctionId(), request.GetBidId())

	if !found {
		return &types.QueryBidResponse{}, status.Error(codes.NotFound, "No bid found for specified bid id and auction id")
	}

	return &types.QueryBidResponse{Bid: &bid}, nil
}

// QueryBidsByAuction implements QueryServer
func (k Keeper) QueryBidsByAuction(c context.Context, request *types.QueryBidsByAuctionRequest) (*types.QueryBidsByAuctionResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)

	bids := k.GetBidsByAuctionID(ctx, request.GetAuctionId())

	if len(bids) == 0 {
		return &types.QueryBidsByAuctionResponse{}, status.Error(codes.NotFound, "No bids found for given auction id")
	}

	return &types.QueryBidsByAuctionResponse{Bids: bids}, nil
}
