package keeper

import (
	"context"

	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/cosmos/cosmos-sdk/types/bech32"
	"github.com/cosmos/cosmos-sdk/types/query"
	"github.com/ethereum/go-ethereum/common"
	"github.com/peggyjv/sommelier/v7/x/addresses/types"
	typesv1 "github.com/peggyjv/sommelier/v7/x/addresses/types/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var _ typesv1.QueryServer = Keeper{}

func (k Keeper) QueryAddressMappings(c context.Context, request *typesv1.QueryAddressMappingsRequest) (*typesv1.QueryAddressMappingsResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)

	var mappings []*typesv1.AddressMapping
	var err error
	store := ctx.KVStore(k.storeKey)
	prefixStore := prefix.NewStore(store, types.GetCosmosToEvmMapPrefix())

	pageRes, err := query.FilteredPaginate(
		prefixStore,
		&request.Pagination,
		func(key []byte, value []byte, accumulate bool) (bool, error) {
			cosmosAddr := sdk.MustBech32ifyAddressBytes(sdk.GetConfig().GetBech32AccountAddrPrefix(), key)
			evmAddr := common.BytesToAddress(value).Hex()
			mapping := typesv1.AddressMapping{
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
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &typesv1.QueryAddressMappingsResponse{AddressMappings: mappings, Pagination: *pageRes}, nil
}

func (k Keeper) QueryAddressMappingByCosmosAddress(c context.Context, request *typesv1.QueryAddressMappingByCosmosAddressRequest) (*typesv1.QueryAddressMappingByCosmosAddressResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)

	_, cosmosAddr, err := bech32.DecodeAndConvert(request.GetCosmosAddress())
	if err != nil {
		return &typesv1.QueryAddressMappingByCosmosAddressResponse{}, status.Errorf(codes.InvalidArgument, "failed to parse cosmos address %s as bech32", request.GetCosmosAddress())
	}

	rawEvmAddr := k.GetEvmAddressByCosmosAddress(ctx, cosmosAddr)

	if rawEvmAddr == nil {
		return &typesv1.QueryAddressMappingByCosmosAddressResponse{}, status.Errorf(codes.NotFound, "no EVM address mappings for cosmos address %s", request.GetCosmosAddress())
	}

	evmAddr := common.BytesToAddress(rawEvmAddr)

	return &typesv1.QueryAddressMappingByCosmosAddressResponse{
		CosmosAddress: request.GetCosmosAddress(),
		EvmAddress:    evmAddr.Hex(),
	}, nil
}

func (k Keeper) QueryAddressMappingByEVMAddress(c context.Context, request *typesv1.QueryAddressMappingByEVMAddressRequest) (*typesv1.QueryAddressMappingByEVMAddressResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)

	if !common.IsHexAddress(request.GetEvmAddress()) {
		return &typesv1.QueryAddressMappingByEVMAddressResponse{}, status.Errorf(codes.InvalidArgument, "invalid EVM address %s", request.GetEvmAddress())
	}

	evmAddr := common.Hex2Bytes(request.GetEvmAddress())
	rawCosmosAddr := k.GetCosmosAddressByEvmAddress(ctx, evmAddr)

	if rawCosmosAddr == nil {
		return &typesv1.QueryAddressMappingByEVMAddressResponse{}, status.Errorf(codes.NotFound, "no cosmos address mapping for EVM address %s", request.GetEvmAddress())
	}

	prefix := sdk.GetConfig().GetBech32AccountAddrPrefix()
	cosmosAddr, err := sdk.Bech32ifyAddressBytes(prefix, rawCosmosAddr)
	if err != nil {
		return &typesv1.QueryAddressMappingByEVMAddressResponse{}, status.Errorf(codes.Internal, "failed to convert cosmos address to bech32: %s", err)
	}

	return &typesv1.QueryAddressMappingByEVMAddressResponse{
		CosmosAddress: cosmosAddr,
		EvmAddress:    request.GetEvmAddress(),
	}, nil
}
