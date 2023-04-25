package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/peggyjv/sommelier/v6/x/cellarfees/types"
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

func (k Keeper) ModuleAccounts(c context.Context, req *types.QueryModuleAccountsRequest) (*types.QueryModuleAccountsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	return &types.QueryModuleAccountsResponse{
		FeesAddress: k.GetFeesAccount(sdk.UnwrapSDKContext(c)).String(),
	}, nil
}
