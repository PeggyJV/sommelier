package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
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

// QueryActiveAuction implements QueryServer
func (k Keeper) QueryActiveAuction(c context.Context, request *types.QueryActiveAuctionRequest) (*types.QueryActiveAuctionResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)

	// Check active auctions
	activeAuction, found := k.GetActiveAuctionByID(ctx, request.GetAuctionId())

	if found {
		return &types.QueryActiveAuctionResponse{Auction: &activeAuction}, nil
	}

	return &types.QueryActiveAuctionResponse{}, status.Errorf(codes.NotFound, "No active auction found for id: %d", request.GetAuctionId())
}

// QueryEndedAuction implements QueryServer
func (k Keeper) QueryEndedAuction(c context.Context, request *types.QueryEndedAuctionRequest) (*types.QueryEndedAuctionResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)

	// Check ended auctions
	endedAuction, found := k.GetEndedAuctionByID(ctx, request.GetAuctionId())

	if found {
		return &types.QueryEndedAuctionResponse{Auction: &endedAuction}, nil
	}

	return &types.QueryEndedAuctionResponse{}, status.Errorf(codes.NotFound, "No ended auction found for id: %d", request.GetAuctionId())
}

// QueryActiveAuctions implements QueryServer
func (k Keeper) QueryActiveAuctions(c context.Context, _ *types.QueryActiveAuctionsRequest) (*types.QueryActiveAuctionsResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)

	auctions := k.GetActiveAuctions(ctx)

	if len(auctions) == 0 {
		return &types.QueryActiveAuctionsResponse{}, status.Error(codes.NotFound, "No active auctions found")
	}

	return &types.QueryActiveAuctionsResponse{Auctions: auctions}, nil
}

// QueryEndedAuctions implements QueryServer
func (k Keeper) QueryEndedAuctions(c context.Context, request *types.QueryEndedAuctionsRequest) (*types.QueryEndedAuctionsResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)

	var auctions []*types.Auction
	store := k.getEndedAuctionsPrefixStore(ctx)
	var err error

	pageRes, err := query.FilteredPaginate(
		store,
		&request.Pagination,
		func(key []byte, value []byte, accumulate bool) (bool, error) {
			var auction types.Auction
			err := auction.Unmarshal(value)
			if err != nil {
				return false, err
			}

			if accumulate {
				auctions = append(auctions, &auction)
			}
			return true, nil
		},
	)

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryEndedAuctionsResponse{Auctions: auctions, Pagination: *pageRes}, nil
}

// QueryBid implements QueryServer
func (k Keeper) QueryBid(c context.Context, request *types.QueryBidRequest) (*types.QueryBidResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)

	bid, found := k.GetBid(ctx, request.GetAuctionId(), request.GetBidId())

	if !found {
		return &types.QueryBidResponse{}, status.Errorf(codes.NotFound, "No bid found for specified bid id: %d, and auction id: %d", request.GetBidId(), request.GetAuctionId())
	}

	return &types.QueryBidResponse{Bid: &bid}, nil
}

// QueryBidsByAuction implements QueryServer
func (k Keeper) QueryBidsByAuction(c context.Context, request *types.QueryBidsByAuctionRequest) (*types.QueryBidsByAuctionResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)

	var bids []*types.Bid
	store := k.getBidsByAuctionPrefixStore(ctx, request.GetAuctionId())
	var err error

	pageRes, err := query.FilteredPaginate(
		store,
		&request.Pagination,
		func(key []byte, value []byte, accumulate bool) (bool, error) {
			var bid types.Bid
			err := bid.Unmarshal(value)
			if err != nil {
				return false, err
			}

			if accumulate {
				bids = append(bids, &bid)
			}
			return true, nil
		},
	)

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryBidsByAuctionResponse{Bids: bids, Pagination: *pageRes}, nil
}
