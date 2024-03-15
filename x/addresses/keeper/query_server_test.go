package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/peggyjv/sommelier/v7/x/addresses/types"
)

// Happy path test for query server functions
func (suite *KeeperTestSuite) TestHappyPathsForQueryServer() {
	ctx, addressesKeeper := suite.ctx, suite.addressesKeeper
	require := suite.Require()

	params := types.DefaultParams()
	addressesKeeper.setParams(ctx, params)

	require.Equal(42, len(evmAddrString))
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
}

// Unhappy path test for query server functions
func (suite *KeeperTestSuite) TestUnhappyPathsForQueryServer() {
	ctx, addressesKeeper := suite.ctx, suite.addressesKeeper
	require := suite.Require()

	params := types.DefaultParams()
	addressesKeeper.setParams(ctx, params)

	// invalid length evm address
	evmAddrStringInvalid := "0x11111111111111111111111111111111111111111"
	// invalid checksum cosmos address
	cosmosAddrStringInvalid := "cosmos154d0p9xhrruhxvazumej9nq29afeura2alje41"
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
