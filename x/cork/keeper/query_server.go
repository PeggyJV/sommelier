package keeper

import (
	"context"

	"github.com/ethereum/go-ethereum/common"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/peggyjv/sommelier/v3/x/cork/types"
)

var _ types.QueryServer = Keeper{}

// QuerySubmittedCorks implements QueryServer
func (k Keeper) QuerySubmittedCorks(c context.Context, _ *types.QuerySubmittedCorksRequest) (*types.QuerySubmittedCorksResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)

	// get corks
	var corks []*types.Cork
	k.IterateCorks(ctx, func(_ sdk.ValAddress, _ common.Address, cork types.Cork) (stop bool) {
		corks = append(corks, &cork)
		return false
	})

	return &types.QuerySubmittedCorksResponse{
		Corks: corks,
	}, nil
}

// QueryParams implements QueryServer
func (k Keeper) QueryParams(c context.Context, _ *types.QueryParamsRequest) (*types.QueryParamsResponse, error) {
	return &types.QueryParamsResponse{
		Params: k.GetParamSet(sdk.UnwrapSDKContext(c)),
	}, nil
}

// QueryCommitPeriod implements QueryServer
func (k Keeper) QueryCommitPeriod(c context.Context, _ *types.QueryCommitPeriodRequest) (*types.QueryCommitPeriodResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	votePeriodStart, found := k.GetCommitPeriodStart(ctx)
	if !found {
		return nil, status.Error(codes.NotFound, "vote period start not set")
	}

	return &types.QueryCommitPeriodResponse{
		VotePeriodStart: votePeriodStart,
		VotePeriodEnd:   votePeriodStart + k.GetParamSet(ctx).VotePeriod,
		CurrentHeight:   ctx.BlockHeight(),
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

	k.IterateScheduledCorks(ctx, func(val sdk.ValAddress, blockHeight uint64, cel common.Address, cork types.Cork) (stop bool) {
		response.Corks = append(response.Corks, &types.ScheduledCork{
			Cork:        &cork,
			BlockHeight: blockHeight,
			Validator:   val.String(),
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
