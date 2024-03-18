package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	query "github.com/cosmos/cosmos-sdk/types/query"
	"github.com/ethereum/go-ethereum/common"
	"github.com/peggyjv/sommelier/v7/x/addresses/types"
)

// Happy path test for query server functions
func (suite *KeeperTestSuite) TestHappyPathsForQueryServer() {
	ctx, addressesKeeper := suite.ctx, suite.addressesKeeper
	require := suite.Require()

	params := types.DefaultParams()
	addressesKeeper.setParams(ctx, params)

	evmAddr := common.HexToAddress(evmAddrString).Bytes()
	acc, err := sdk.AccAddressFromBech32(cosmosAddrString)
	require.NoError(err)

	cosmosAddr := acc.Bytes()

	addressesKeeper.SetAddressMapping(ctx, cosmosAddr, evmAddr)

	// Test QueryParams
	queryParams, err := addressesKeeper.QueryParams(sdk.WrapSDKContext(ctx), &types.QueryParamsRequest{})
	require.NoError(err)
	require.Equal(&params, queryParams.Params)

	// Test QueryAddressMappingByEvmAddress
	addressMappingByEvmAddress, err := addressesKeeper.QueryAddressMappingByEVMAddress(sdk.WrapSDKContext(ctx), &types.QueryAddressMappingByEVMAddressRequest{EvmAddress: evmAddrString})
	require.NoError(err)
	require.Equal(cosmosAddrString, addressMappingByEvmAddress.CosmosAddress)

	// Test QueryAddressMappingByCosmosAddress
	addressMappingByCosmosAddress, err := addressesKeeper.QueryAddressMappingByCosmosAddress(sdk.WrapSDKContext(ctx), &types.QueryAddressMappingByCosmosAddressRequest{CosmosAddress: cosmosAddrString})
	require.NoError(err)
	require.Equal(evmAddrString, addressMappingByCosmosAddress.EvmAddress)

	// Test QueryAddressMappings
	addressMappings, err := addressesKeeper.QueryAddressMappings(sdk.WrapSDKContext(ctx), &types.QueryAddressMappingsRequest{})
	require.NoError(err)
	require.Len(addressMappings.AddressMappings, 1)
	require.Equal(cosmosAddrString, addressMappings.AddressMappings[0].CosmosAddress)
	require.Equal(evmAddrString, addressMappings.AddressMappings[0].EvmAddress)

	// Test QueryAddressMappings with pagination
	addressMappings, err = addressesKeeper.QueryAddressMappings(sdk.WrapSDKContext(ctx), &types.QueryAddressMappingsRequest{Pagination: query.PageRequest{Limit: 1}})
	require.NoError(err)
	require.Len(addressMappings.AddressMappings, 1)
	require.Equal(cosmosAddrString, addressMappings.AddressMappings[0].CosmosAddress)
	require.Equal(evmAddrString, addressMappings.AddressMappings[0].EvmAddress)

	evmAddrString2 := "0x2222222222222222222222222222222222222222"
	// keys stored in ascending order, so this one will be stored before the previous value (cosmos15...)
	cosmosAddrString2 := "cosmos1y6d5kasehecexf09ka6y0ggl0pxzt6dgk0gnl9"
	evmAddr2 := common.HexToAddress(evmAddrString2).Bytes()
	acc2, err := sdk.AccAddressFromBech32(cosmosAddrString2)
	require.NoError(err)

	cosmosAddr2 := acc2.Bytes()

	addressesKeeper.SetAddressMapping(ctx, cosmosAddr2, evmAddr2)

	addressMappings, err = addressesKeeper.QueryAddressMappings(sdk.WrapSDKContext(ctx), &types.QueryAddressMappingsRequest{Pagination: query.PageRequest{Limit: 1, Offset: 1}})
	require.NoError(err)
	require.Len(addressMappings.AddressMappings, 1)
	require.Equal(cosmosAddrString, addressMappings.AddressMappings[0].CosmosAddress)
	require.Equal(evmAddrString, addressMappings.AddressMappings[0].EvmAddress)

	addressMappings, err = addressesKeeper.QueryAddressMappings(sdk.WrapSDKContext(ctx), &types.QueryAddressMappingsRequest{Pagination: query.PageRequest{Limit: 1}})
	require.NoError(err)
	require.Len(addressMappings.AddressMappings, 1)
	require.Equal(cosmosAddrString2, addressMappings.AddressMappings[0].CosmosAddress)
	require.Equal(evmAddrString2, addressMappings.AddressMappings[0].EvmAddress)

	addressMappings, err = addressesKeeper.QueryAddressMappings(sdk.WrapSDKContext(ctx), &types.QueryAddressMappingsRequest{Pagination: query.PageRequest{Key: addressMappings.Pagination.NextKey}})
	require.NoError(err)
	require.Len(addressMappings.AddressMappings, 1)
	require.Equal(cosmosAddrString, addressMappings.AddressMappings[0].CosmosAddress)
	require.Equal(evmAddrString, addressMappings.AddressMappings[0].EvmAddress)

	addressMappings, err = addressesKeeper.QueryAddressMappings(sdk.WrapSDKContext(ctx), &types.QueryAddressMappingsRequest{Pagination: query.PageRequest{Limit: 2}})
	require.NoError(err)
	require.Len(addressMappings.AddressMappings, 2)
	require.Equal(cosmosAddrString, addressMappings.AddressMappings[1].CosmosAddress)
	require.Equal(evmAddrString, addressMappings.AddressMappings[1].EvmAddress)
	require.Equal(cosmosAddrString2, addressMappings.AddressMappings[0].CosmosAddress)
	require.Equal(evmAddrString2, addressMappings.AddressMappings[0].EvmAddress)
}

// Unhappy path test for query server functions
func (suite *KeeperTestSuite) TestUnhappyPathsForQueryServer() {
	ctx, addressesKeeper := suite.ctx, suite.addressesKeeper
	require := suite.Require()

	params := types.DefaultParams()
	addressesKeeper.setParams(ctx, params)

	// invalid length evm address
	evmAddrStringInvalid := "0x11111111111111111111111111111111111111111"
	_, err := sdk.AccAddressFromBech32(cosmosAddrStringInvalid)
	require.Error(err)

	// Test QueryParams
	queryParams, err := addressesKeeper.QueryParams(sdk.WrapSDKContext(ctx), &types.QueryParamsRequest{})
	require.NoError(err)
	require.Equal(&params, queryParams.Params)

	// Test QueryAddressMappingByEvmAddress
	_, err = addressesKeeper.QueryAddressMappingByEVMAddress(sdk.WrapSDKContext(ctx), &types.QueryAddressMappingByEVMAddressRequest{EvmAddress: evmAddrStringInvalid})
	require.Error(err)
	require.Contains(err.Error(), "invalid EVM address")

	// valid evm address
	_, err = addressesKeeper.QueryAddressMappingByEVMAddress(sdk.WrapSDKContext(ctx), &types.QueryAddressMappingByEVMAddressRequest{EvmAddress: evmAddrString})
	require.Error(err)
	require.Contains(err.Error(), "no cosmos address mapping for EVM address")

	// Test QueryAddressMappingByCosmosAddress
	_, err = addressesKeeper.QueryAddressMappingByCosmosAddress(sdk.WrapSDKContext(ctx), &types.QueryAddressMappingByCosmosAddressRequest{CosmosAddress: cosmosAddrStringInvalid})
	require.Error(err)
	require.Contains(err.Error(), "failed to parse cosmos address")

	// valid cosmos address
	_, err = addressesKeeper.QueryAddressMappingByCosmosAddress(sdk.WrapSDKContext(ctx), &types.QueryAddressMappingByCosmosAddressRequest{CosmosAddress: cosmosAddrString})
	require.Error(err)
	require.Contains(err.Error(), "no EVM address mapping for cosmos address")
}
