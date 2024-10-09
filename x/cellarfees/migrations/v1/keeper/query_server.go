package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/peggyjv/sommelier/v7/x/cellarfees/migrations/v1/types"
	v1types "github.com/peggyjv/sommelier/v7/x/cellarfees/migrations/v1/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var _ types.QueryServer = Keeper{}

func (k Keeper) QueryParams(c context.Context, req *v1types.QueryParamsRequest) (*v1types.QueryParamsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	return &v1types.QueryParamsResponse{
		Params: k.GetParams(sdk.UnwrapSDKContext(c)),
	}, nil
}

func (k Keeper) QueryModuleAccounts(c context.Context, req *v1types.QueryModuleAccountsRequest) (*v1types.QueryModuleAccountsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	return &v1types.QueryModuleAccountsResponse{
		FeesAddress: k.GetFeesAccount(sdk.UnwrapSDKContext(c)).GetAddress().String(),
	}, nil
}

func (k Keeper) QueryLastRewardSupplyPeak(c context.Context, req *v1types.QueryLastRewardSupplyPeakRequest) (*v1types.QueryLastRewardSupplyPeakResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	return &v1types.QueryLastRewardSupplyPeakResponse{LastRewardSupplyPeak: k.GetLastRewardSupplyPeak(sdk.UnwrapSDKContext(c))}, nil
}

func (k Keeper) QueryFeeAccrualCounters(c context.Context, req *v1types.QueryFeeAccrualCountersRequest) (*v1types.QueryFeeAccrualCountersResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	return &v1types.QueryFeeAccrualCountersResponse{FeeAccrualCounters: k.GetFeeAccrualCounters(sdk.UnwrapSDKContext(c))}, nil
}

func (k Keeper) QueryAPY(c context.Context, _ *v1types.QueryAPYRequest) (*v1types.QueryAPYResponse, error) {
	return &v1types.QueryAPYResponse{
		Apy: k.GetAPY(sdk.UnwrapSDKContext(c)).String(),
	}, nil
}
