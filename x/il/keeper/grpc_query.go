package keeper

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"

	"github.com/peggyjv/sommelier/x/il/types"
)

var _ types.QueryServer = Keeper{}

// Stoploss implements QueryServer.Stoploss
func (k Keeper) Stoploss(c context.Context, req *types.QueryStoplossRequest) (*types.QueryStoplossResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	address, err := sdk.AccAddressFromBech32(req.Address)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	ctx := sdk.UnwrapSDKContext(c)
	stoploss, found := k.GetStoplossPosition(ctx, address, req.UniswapPair)
	if !found {
		return nil, status.Errorf(codes.NotFound, "stoploss position for address %s and pair %s", req.Address, req.UniswapPair)
	}

	return &types.QueryStoplossResponse{
		Stoploss: stoploss,
	}, nil
}

// StoplossPositions implements QueryServer.StoplossPositions
func (k Keeper) StoplossPositions(c context.Context, req *types.QueryStoplossPositionsRequest) (*types.QueryStoplossPositionsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	address, err := sdk.AccAddressFromBech32(req.Address)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	ctx := sdk.UnwrapSDKContext(c)
	store := prefix.NewStore(ctx.KVStore(k.storeKey), append(types.StoplossKeyPrefix, address.Bytes()...))

	stoplossPositions := []types.Stoploss{}
	pageRes, err := query.Paginate(store, req.Pagination, func(_, value []byte) error {
		var stoploss types.Stoploss
		k.cdc.MustUnmarshalBinaryBare(value, &stoploss)

		stoplossPositions = append(stoplossPositions, stoploss)
		return nil
	})

	if err != nil {
		return nil, err
	}

	return &types.QueryStoplossPositionsResponse{
		StoplossPositions: stoplossPositions,
		Pagination:        pageRes,
	}, nil
}

// Parameters implements QueryServer.Parameters
func (k Keeper) Parameters(c context.Context, _ *types.QueryParametersRequest) (*types.QueryParametersResponse, error) {
	return &types.QueryParametersResponse{
		Params: k.GetParams(sdk.UnwrapSDKContext(c)),
	}, nil
}
