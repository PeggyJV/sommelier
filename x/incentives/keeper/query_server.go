package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/peggyjv/sommelier/v8/x/incentives/types"
)

var _ types.QueryServer = Keeper{}

// QueryParams implements QueryServer
func (k Keeper) QueryParams(c context.Context, _ *types.QueryParamsRequest) (*types.QueryParamsResponse, error) {
	return &types.QueryParamsResponse{
		Params: k.GetParamSet(sdk.UnwrapSDKContext(c)),
	}, nil
}

func (k Keeper) QueryAPY(c context.Context, _ *types.QueryAPYRequest) (*types.QueryAPYResponse, error) {
	return &types.QueryAPYResponse{
		Apy: k.GetAPY(sdk.UnwrapSDKContext(c)).String(),
	}, nil
}
