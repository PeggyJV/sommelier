package keeper

import (
	"context"

	"github.com/ethereum/go-ethereum/common"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/peggyjv/sommelier/v4/x/cork/types"
)

var _ types.QueryServer = Keeper{}

// QueryParams implements QueryServer
func (k Keeper) QueryParams(c context.Context, _ *types.QueryParamsRequest) (*types.QueryParamsResponse, error) {
	return &types.QueryParamsResponse{
		Params: k.GetParamSet(sdk.UnwrapSDKContext(c)),
	}, nil
}

func (k Keeper) QueryCellarIDs(c context.Context, _ *types.QueryCellarIDsRequest) (*types.QueryCellarIDsResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)

	response := &types.QueryCellarIDsResponse{}
	for _, id := range k.GetCellarIDs(ctx) {
		response.CellarIds = append(response.CellarIds, id.Hex())
	}

	return response, nil
}

func (k Keeper) QueryScheduledCorks(c context.Context, _ *types.QueryScheduledCorksRequest) (*types.QueryScheduledCorksResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)

	response := types.QueryScheduledCorksResponse{}

	k.IterateScheduledCorks(ctx, func(val sdk.ValAddress, blockHeight uint64, id uint64, cel common.Address, cork types.Cork) (stop bool) {
		response.Corks = append(response.Corks, &types.ScheduledCork{
			Cork:        &cork,
			BlockHeight: blockHeight,
			Validator:   val.String(),
			Id:          id,
		})
		return false
	})
	return &response, nil
}

func (k Keeper) QueryScheduledBlockHeights(c context.Context, _ *types.QueryScheduledBlockHeightsRequest) (*types.QueryScheduledBlockHeightsResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	response := types.QueryScheduledBlockHeightsResponse{}
	response.BlockHeights = k.GetScheduledBlockHeights(ctx)
	return &response, nil
}

func (k Keeper) QueryScheduledCorksByBlockHeight(c context.Context, req *types.QueryScheduledCorksByBlockHeightRequest) (*types.QueryScheduledCorksByBlockHeightResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)

	response := types.QueryScheduledCorksByBlockHeightResponse{}
	response.Corks = k.GetScheduledCorksByBlockHeight(ctx, req.BlockHeight)
	return &response, nil
}

func (k Keeper) QueryScheduledCorksByID(c context.Context, req *types.QueryScheduledCorksByIDRequest) (*types.QueryScheduledCorksByIDResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)

	response := types.QueryScheduledCorksByIDResponse{}
	response.Corks = k.GetScheduledCorksByID(ctx, req.Id)
	return &response, nil
}
