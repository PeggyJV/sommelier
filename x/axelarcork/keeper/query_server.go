package keeper

import (
	"context"
	"encoding/hex"

	"github.com/ethereum/go-ethereum/common"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/peggyjv/sommelier/v6/x/axelarcork/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var _ types.QueryServer = Keeper{}

// QueryParams implements QueryServer
func (k Keeper) QueryParams(c context.Context, req *types.QueryParamsRequest) (*types.QueryParamsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	return &types.QueryParamsResponse{
		Params: k.GetParamSet(sdk.UnwrapSDKContext(c)),
	}, nil
}

func (k Keeper) QueryCellarIDs(c context.Context, req *types.QueryCellarIDsRequest) (*types.QueryCellarIDsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	ctx := sdk.UnwrapSDKContext(c)
	config, err := k.GetChainConfigurationByNameAndID(ctx, req.ChainName, req.ChainId)
	if err != nil {
		return nil, err
	}

	response := &types.QueryCellarIDsResponse{}
	for _, id := range k.GetCellarIDs(ctx, config.Id) {
		response.CellarIds = append(response.CellarIds, id.Hex())
	}

	return response, nil
}

func (k Keeper) QueryScheduledCorks(c context.Context, req *types.QueryScheduledCorksRequest) (*types.QueryScheduledCorksResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	ctx := sdk.UnwrapSDKContext(c)
	config, err := k.GetChainConfigurationByNameAndID(ctx, req.ChainName, req.ChainId)
	if err != nil {
		return nil, err
	}

	response := types.QueryScheduledCorksResponse{}

	k.IterateScheduledCorks(ctx, config.Id, func(val sdk.ValAddress, blockHeight uint64, id []byte, cel common.Address, cork types.Cork) (stop bool) {
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

func (k Keeper) QueryScheduledBlockHeights(c context.Context, req *types.QueryScheduledBlockHeightsRequest) (*types.QueryScheduledBlockHeightsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	ctx := sdk.UnwrapSDKContext(c)
	config, err := k.GetChainConfigurationByNameAndID(ctx, req.ChainName, req.ChainId)
	if err != nil {
		return nil, err
	}

	response := types.QueryScheduledBlockHeightsResponse{}
	response.BlockHeights = k.GetScheduledBlockHeights(ctx, config.Id)
	return &response, nil
}

func (k Keeper) QueryScheduledCorksByBlockHeight(c context.Context, req *types.QueryScheduledCorksByBlockHeightRequest) (*types.QueryScheduledCorksByBlockHeightResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	ctx := sdk.UnwrapSDKContext(c)
	config, err := k.GetChainConfigurationByNameAndID(ctx, req.ChainName, req.ChainId)
	if err != nil {
		return nil, err
	}

	response := types.QueryScheduledCorksByBlockHeightResponse{}
	response.Corks = k.GetScheduledCorksByBlockHeight(ctx, config.Id, req.BlockHeight)
	return &response, nil
}

func (k Keeper) QueryScheduledCorksByID(c context.Context, req *types.QueryScheduledCorksByIDRequest) (*types.QueryScheduledCorksByIDResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	ctx := sdk.UnwrapSDKContext(c)
	config, err := k.GetChainConfigurationByNameAndID(ctx, req.ChainName, req.ChainId)
	if err != nil {
		return nil, err
	}

	id, err := hex.DecodeString(req.Id)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "Failed to decode %s from hexidecimal to bytes", req.Id)
	}

	response := types.QueryScheduledCorksByIDResponse{}
	response.Corks = k.GetScheduledCorksByID(ctx, config.Id, id)
	return &response, nil
}

func (k Keeper) QueryCorkResult(c context.Context, req *types.QueryCorkResultRequest) (*types.QueryCorkResultResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	ctx := sdk.UnwrapSDKContext(c)
	config, err := k.GetChainConfigurationByNameAndID(ctx, req.ChainName, req.ChainId)
	if err != nil {
		return nil, err
	}

	id, err := hex.DecodeString(req.Id)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "Failed to decode %s from hexidecimal to bytes", req.Id)
	}

	response := types.QueryCorkResultResponse{}
	var found bool
	result, found := k.GetCorkResult(ctx, config.Id, id)
	if !found {
		return &types.QueryCorkResultResponse{}, status.Errorf(codes.NotFound, "No cork result found for id: %s", req.GetId())
	}
	response.CorkResult = &result

	return &response, nil
}

func (k Keeper) QueryCorkResults(c context.Context, req *types.QueryCorkResultsRequest) (*types.QueryCorkResultsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	ctx := sdk.UnwrapSDKContext(c)
	config, err := k.GetChainConfigurationByNameAndID(ctx, req.ChainName, req.ChainId)
	if err != nil {
		return nil, err
	}

	response := types.QueryCorkResultsResponse{}
	response.CorkResults = k.GetCorkResults(ctx, config.Id)
	return &response, nil
}

func (k Keeper) QueryChainConfigurations(c context.Context, req *types.QueryChainConfigurationsRequest) (*types.QueryChainConfigurationsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	ctx := sdk.UnwrapSDKContext(c)
	var chainConfigurations []*types.ChainConfiguration
	k.IterateChainConfigurations(ctx, func(config types.ChainConfiguration) (stop bool) {
		chainConfigurations = append(chainConfigurations, &config)
		return false
	})

	return &types.QueryChainConfigurationsResponse{Configurations: chainConfigurations}, nil
}