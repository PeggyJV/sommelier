package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/peggyjv/sommelier/v4/x/cellarfees/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var _ types.QueryServer = Keeper{}

func (k Keeper) QueryParams(c context.Context, req *types.QueryParamsRequest) (*types.QueryParamsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	return &types.QueryParamsResponse{
		Params: k.GetParams(sdk.UnwrapSDKContext(c)),
	}, nil
}

func (k Keeper) QueryModuleAccounts(c context.Context, req *types.QueryModuleAccountsRequest) (*types.QueryModuleAccountsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	return &types.QueryModuleAccountsResponse{
		FeesAddress: k.GetFeesAccount(sdk.UnwrapSDKContext(c)).GetAddress().String(),
	}, nil
}

func (k Keeper) QueryLastRewardSupplyPeak(c context.Context, req *types.QueryLastRewardSupplyPeakRequest) (*types.QueryLastRewardSupplyPeakResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	return &types.QueryLastRewardSupplyPeakResponse{LastRewardSupplyPeak: k.GetLastRewardSupplyPeak(sdk.UnwrapSDKContext(c))}, nil
}

func (k Keeper) QueryFeeAccrualCounters(c context.Context, req *types.QueryFeeAccrualCountersRequest) (*types.QueryFeeAccrualCountersResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	return &types.QueryFeeAccrualCountersResponse{FeeAccrualCounters: k.GetFeeAccrualCounters(sdk.UnwrapSDKContext(c))}, nil
}

func (k Keeper) QueryAPY(c context.Context, _ *types.QueryAPYRequest) (*types.QueryAPYResponse, error) {
	return &types.QueryAPYResponse{
		Apy: k.GetAPY(sdk.UnwrapSDKContext(c)).String(),
	}, nil
}
