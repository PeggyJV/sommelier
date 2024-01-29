package keeper

import (
	"context"
	"encoding/hex"
	"fmt"

	"github.com/ethereum/go-ethereum/common"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/peggyjv/sommelier/v7/x/axelarcork/types"
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
	response := &types.QueryCellarIDsResponse{
		CellarIds: []*types.CellarIDSet{},
	}
	k.IterateChainConfigurations(ctx, func(config types.ChainConfiguration) (stop bool) {
		set := types.CellarIDSet{ChainId: config.Id, Ids: []string{}}
		for _, id := range k.GetCellarIDs(ctx, config.Id) {
			set.Ids = append(set.Ids, id.String())
		}

		response.CellarIds = append(response.CellarIds, &set)

		return false
	})

	return response, nil
}

func (k Keeper) QueryCellarIDsByChainID(c context.Context, req *types.QueryCellarIDsByChainIDRequest) (*types.QueryCellarIDsByChainIDResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	ctx := sdk.UnwrapSDKContext(c)
	config, ok := k.GetChainConfigurationByID(ctx, req.ChainId)
	if !ok {
		return nil, fmt.Errorf("chain by id %d not found", req.ChainId)
	}

	response := &types.QueryCellarIDsByChainIDResponse{}
	for _, id := range k.GetCellarIDs(ctx, config.Id) {
		response.CellarIds = append(response.CellarIds, id.String())
	}

	return response, nil
}

func (k Keeper) QueryScheduledCorks(c context.Context, req *types.QueryScheduledCorksRequest) (*types.QueryScheduledCorksResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	ctx := sdk.UnwrapSDKContext(c)
	config, ok := k.GetChainConfigurationByID(ctx, req.ChainId)
	if !ok {
		return nil, fmt.Errorf("chain by id %d not found", req.ChainId)
	}

	response := types.QueryScheduledCorksResponse{}

	k.IterateScheduledAxelarCorks(ctx, config.Id, func(val sdk.ValAddress, blockHeight uint64, id []byte, cel common.Address, cork types.AxelarCork) (stop bool) {
		response.Corks = append(response.Corks, &types.ScheduledAxelarCork{
			Cork:        &cork,
			BlockHeight: blockHeight,
			Validator:   val.String(),
			Id:          hex.EncodeToString(id),
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
	config, ok := k.GetChainConfigurationByID(ctx, req.ChainId)
	if !ok {
		return nil, fmt.Errorf("chain by id %d not found", req.ChainId)
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
	config, ok := k.GetChainConfigurationByID(ctx, req.ChainId)
	if !ok {
		return nil, fmt.Errorf("chain by id %d not found", req.ChainId)
	}

	response := types.QueryScheduledCorksByBlockHeightResponse{}
	response.Corks = k.GetScheduledAxelarCorksByBlockHeight(ctx, config.Id, req.BlockHeight)
	return &response, nil
}

func (k Keeper) QueryScheduledCorksByID(c context.Context, req *types.QueryScheduledCorksByIDRequest) (*types.QueryScheduledCorksByIDResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	ctx := sdk.UnwrapSDKContext(c)
	config, ok := k.GetChainConfigurationByID(ctx, req.ChainId)
	if !ok {
		return nil, fmt.Errorf("chain by id %d not found", req.ChainId)
	}

	id, err := hex.DecodeString(req.Id)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "Failed to decode %s from hexadecimal to bytes", req.Id)
	}

	response := types.QueryScheduledCorksByIDResponse{}
	response.Corks = k.GetScheduledAxelarCorksByID(ctx, config.Id, id)
	return &response, nil
}

func (k Keeper) QueryCorkResult(c context.Context, req *types.QueryCorkResultRequest) (*types.QueryCorkResultResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	ctx := sdk.UnwrapSDKContext(c)
	config, ok := k.GetChainConfigurationByID(ctx, req.ChainId)
	if !ok {
		return nil, fmt.Errorf("chain by id %d not found", req.ChainId)
	}

	id, err := hex.DecodeString(req.Id)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "Failed to decode %s from hexadecimal to bytes", req.Id)
	}

	response := types.QueryCorkResultResponse{}
	var found bool
	result, found := k.GetAxelarCorkResult(ctx, config.Id, id)
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
	config, ok := k.GetChainConfigurationByID(ctx, req.ChainId)
	if !ok {
		return nil, fmt.Errorf("chain by id %d not found", req.ChainId)
	}

	response := types.QueryCorkResultsResponse{}
	response.CorkResults = k.GetAxelarCorkResults(ctx, config.Id)
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

func (k Keeper) QueryAxelarContractCallNonces(c context.Context, req *types.QueryAxelarContractCallNoncesRequest) (*types.QueryAxelarContractCallNoncesResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	ctx := sdk.UnwrapSDKContext(c)

	response := types.QueryAxelarContractCallNoncesResponse{}

	nonces := []*types.AxelarContractCallNonce{}
	k.IterateAxelarContractCallNonces(ctx, func(chainID uint64, address common.Address, nonce uint64) (stop bool) {
		nonces = append(nonces, &types.AxelarContractCallNonce{
			ChainId:         chainID,
			ContractAddress: address.String(),
			Nonce:           nonce,
		})

		return false
	})

	response.ContractCallNonces = nonces

	return &response, nil
}

func (k Keeper) QueryAxelarProxyUpgradeData(c context.Context, req *types.QueryAxelarProxyUpgradeDataRequest) (*types.QueryAxelarProxyUpgradeDataResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	ctx := sdk.UnwrapSDKContext(c)

	response := types.QueryAxelarProxyUpgradeDataResponse{}

	upgradeData := []*types.AxelarUpgradeData{}
	k.IterateAxelarProxyUpgradeData(ctx, func(chainID uint64, data types.AxelarUpgradeData) (stop bool) {
		upgradeData = append(upgradeData, &data)
		return false
	})

	response.ProxyUpgradeData = upgradeData

	return &response, nil
}

func (k Keeper) QueryWinningAxelarCork(c context.Context, req *types.QueryWinningAxelarCorkRequest) (*types.QueryWinningAxelarCorkResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	ctx := sdk.UnwrapSDKContext(c)

	response := types.QueryWinningAxelarCorkResponse{}

	height, cork, found := k.GetWinningAxelarCork(ctx, req.ChainId, common.HexToAddress(req.ContractAddress))
	if !found {
		return &types.QueryWinningAxelarCorkResponse{}, status.Errorf(codes.NotFound, "No winning cork found for chain id: %d", req.GetChainId())
	}

	response.Cork = &cork
	response.BlockHeight = height

	return &response, nil
}

func (k Keeper) QueryWinningAxelarCorks(c context.Context, req *types.QueryWinningAxelarCorksRequest) (*types.QueryWinningAxelarCorksResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	ctx := sdk.UnwrapSDKContext(c)

	response := types.QueryWinningAxelarCorksResponse{}

	winning := []*types.WinningAxelarCork{}
	k.IterateWinningAxelarCorks(ctx, req.ChainId, func(contract common.Address, blockHeight uint64, cork types.AxelarCork) (stop bool) {
		winning = append(winning, &types.WinningAxelarCork{
			Cork:        &cork,
			BlockHeight: blockHeight,
		})
		return false
	})

	response.WinningAxelarCorks = winning

	return &response, nil
}
