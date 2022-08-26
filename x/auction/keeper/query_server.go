package keeper

import (
	"context"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/peggyjv/sommelier/v4/x/auction/types"
)

var _ types.QueryServer = Keeper{}

// QueryParams implements QueryServer
func (k Keeper) QueryParams(c context.Context, _ *types.QueryParamsRequest) (*types.QueryParamsResponse, error) {
	return &types.QueryParamsResponse{
		Params: k.GetParamSet(sdk.UnwrapSDKContext(c)),
	}, nil
}

// QueryAuction implements QueryServer
func (k Keeper) QueryAuction(c context.Context, request *types.QueryAuctionRequest) (*types.QueryAuctionResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)

	// Check active auctions first
	activeAuction, found := k.GetActiveAuctionById(ctx, request.GetAuctionId())

	if found {
		return &types.QueryAuctionResponse{&activeAuction}, nil
	} 

	// Check ended auctions first
	endedAuction, found := k.GetEndedAuctionById(ctx, request.GetAuctionId())

	if found {
		return &types.QueryAuctionResponse{&endedAuction}, nil
	} 

	return &types.QueryAuctionResponse{},  fmt.Errorf("No auction found for given id")
}

// QueryCurrentAuctions implements QueryServer
func (k Keeper) QueryCurrentAuctions(c context.Context, _ *types.QueryCurrentAuctionsRequest) (*types.QueryCurrentAuctionsResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	
	auctions := k.GetActiveAuctions(ctx)

	if len(auctions) == 0 {
		return &types.QueryCurrentAuctionsResponse{}, fmt.Errorf("No active auctions found")
	}

	return &types.QueryCurrentAuctionsResponse{auctions}, nil
}

// QueryEndedAuctions implements QueryServer
func (k Keeper) QueryEndedAuctions(c context.Context, _ *types.QueryEndedAuctionsRequest) (*types.QueryEndedAuctionsResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	
	auctions := k.GetEndedAuctions(ctx)

	if len(auctions) == 0 {
		return &types.QueryEndedAuctionsResponse{}, fmt.Errorf("No ended auctions found")
	}

	return &types.QueryEndedAuctionsResponse{auctions}, nil
}

// QueryBid implements QueryServer
func (k Keeper) QueryBid(c context.Context, request *types.QueryBidRequest) (*types.QueryBidResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)

	bid, found := k.GetBid(ctx, request.GetAuctionId() ,request.GetBidId())
	
	if !found {
		return &types.QueryBidResponse{}, fmt.Errorf("No bid found for specified bid id and auction id")
	}

	return &types.QueryBidResponse{&bid}, nil
}

// QueryBidsByAuction implements QueryServer
func (k Keeper) QueryBidsByAuction(c context.Context, request *types.QueryBidsByAuctionRequest) (*types.QueryBidsByAuctionResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	
	bids := k.GetBidsByAuctionId(ctx, request.GetAuctionId())

	if len(bids) == 0 {
		return &types.QueryBidsByAuctionResponse{}, fmt.Errorf("No bids found for given auction id")
	}

	return &types.QueryBidsByAuctionResponse{bids}, nil
}