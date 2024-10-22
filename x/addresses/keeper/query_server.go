package keeper

import (
	"context"

	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/cosmos/cosmos-sdk/types/bech32"
	"github.com/cosmos/cosmos-sdk/types/query"
	"github.com/ethereum/go-ethereum/common"
	"github.com/peggyjv/sommelier/v8/x/addresses/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var _ types.QueryServer = Keeper{}

// QueryParams implements QueryServer
func (k Keeper) QueryParams(c context.Context, _ *types.QueryParamsRequest) (*types.QueryParamsResponse, error) {
	params := k.GetParamSet(sdk.UnwrapSDKContext(c))

	return &types.QueryParamsResponse{
		Params: &params,
	}, nil
}

func (k Keeper) QueryAddressMappings(c context.Context, request *types.QueryAddressMappingsRequest) (*types.QueryAddressMappingsResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)

	var mappings []*types.AddressMapping
	var err error
	store := ctx.KVStore(k.storeKey)
	prefixStore := prefix.NewStore(store, types.GetCosmosToEvmMapPrefix())

	pageRes, err := query.FilteredPaginate(
		prefixStore,
		&request.Pagination,
		func(key []byte, value []byte, accumulate bool) (bool, error) {
			cosmosAddr, err := sdk.Bech32ifyAddressBytes(sdk.GetConfig().GetBech32AccountAddrPrefix(), key)
			if err != nil {
				return false, err
			}

			evmAddr := common.BytesToAddress(value).Hex()
			mapping := types.AddressMapping{
				CosmosAddress: cosmosAddr,
				EvmAddress:    evmAddr,
			}

			if accumulate {
				mappings = append(mappings, &mapping)
			}
			return true, nil
		},
	)

	if err != nil {
		// this shouldn't be possible if the msg server is doing proper validation
		k.Logger(ctx).Error("failed to paginate cosmos to evm address mappings", "error", err)
		return nil, status.Error(codes.Internal, "error during pagination")
	}

	return &types.QueryAddressMappingsResponse{AddressMappings: mappings, Pagination: *pageRes}, nil
}

func (k Keeper) QueryAddressMappingByCosmosAddress(c context.Context, request *types.QueryAddressMappingByCosmosAddressRequest) (*types.QueryAddressMappingByCosmosAddressResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)

	_, cosmosAddr, err := bech32.DecodeAndConvert(request.GetCosmosAddress())
	if err != nil {
		return &types.QueryAddressMappingByCosmosAddressResponse{}, status.Errorf(codes.InvalidArgument, "failed to parse cosmos address %s as bech32", request.GetCosmosAddress())
	}

	rawEvmAddr := k.GetEvmAddressByCosmosAddress(ctx, cosmosAddr)

	if len(rawEvmAddr) == 0 {
		return &types.QueryAddressMappingByCosmosAddressResponse{}, status.Errorf(codes.NotFound, "no EVM address mapping for cosmos address %s", request.GetCosmosAddress())
	}

	evmAddr := common.BytesToAddress(rawEvmAddr)

	return &types.QueryAddressMappingByCosmosAddressResponse{
		CosmosAddress: request.GetCosmosAddress(),
		EvmAddress:    evmAddr.Hex(),
	}, nil
}

func (k Keeper) QueryAddressMappingByEVMAddress(c context.Context, request *types.QueryAddressMappingByEVMAddressRequest) (*types.QueryAddressMappingByEVMAddressResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)

	if !common.IsHexAddress(request.GetEvmAddress()) {
		return &types.QueryAddressMappingByEVMAddressResponse{}, status.Errorf(codes.InvalidArgument, "invalid EVM address %s", request.GetEvmAddress())
	}

	evmAddr := common.HexToAddress(request.GetEvmAddress()).Bytes()
	rawCosmosAddr := k.GetCosmosAddressByEvmAddress(ctx, evmAddr)

	if len(rawCosmosAddr) == 0 {
		return &types.QueryAddressMappingByEVMAddressResponse{}, status.Errorf(codes.NotFound, "no cosmos address mapping for EVM address %s", request.GetEvmAddress())
	}

	prefix := sdk.GetConfig().GetBech32AccountAddrPrefix()
	cosmosAddr, err := sdk.Bech32ifyAddressBytes(prefix, rawCosmosAddr)
	if err != nil {
		// this shouldn't happen if msg server is doing proper validation
		return &types.QueryAddressMappingByEVMAddressResponse{}, status.Errorf(codes.Internal, "failed to convert cosmos address to bech32: %s", err)
	}

	return &types.QueryAddressMappingByEVMAddressResponse{
		CosmosAddress: cosmosAddr,
		EvmAddress:    request.GetEvmAddress(),
	}, nil
}
