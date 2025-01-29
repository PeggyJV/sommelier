package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/peggyjv/sommelier/v9/x/addresses/types"
)

func (suite *KeeperTestSuite) TestImportExportGenesis() {
	ctx, addressesKeeper := suite.ctx, suite.addressesKeeper
	require := suite.Require()

	expectedGenesis := types.DefaultGenesisState()

	InitGenesis(ctx, addressesKeeper, expectedGenesis)

	actualGenesis := ExportGenesis(ctx, addressesKeeper)
	require.Equal(expectedGenesis, actualGenesis)
}

func (suite *KeeperTestSuite) TestGenesisValidation() {
	require := suite.Require()

	genesis := types.DefaultGenesisState()
	require.NoError(genesis.Validate())

	genesis.AddressMappings = append(genesis.AddressMappings, &types.AddressMapping{CosmosAddress: "sldjflslkfjsdf", EvmAddress: "0x0000000000000000000000000000000000000000"})
	require.Error(genesis.Validate())

	genesis.AddressMappings = types.DefaultGenesisState().AddressMappings
	genesis.AddressMappings = append(genesis.AddressMappings, &types.AddressMapping{CosmosAddress: cosmosAddrString, EvmAddress: "zzzz"})
	require.Error(genesis.Validate())

	// Test with invalid EVM address
	genesis.AddressMappings = []*types.AddressMapping{
		{CosmosAddress: cosmosAddrString, EvmAddress: "invalid_evm_address"},
	}
	require.Error(genesis.Validate())

	// Test with duplicate mappings
	genesis.AddressMappings = []*types.AddressMapping{
		{CosmosAddress: cosmosAddrString, EvmAddress: evmAddrString},
		{CosmosAddress: cosmosAddrString, EvmAddress: evmAddrString},
	}
	require.Error(genesis.Validate())

	// Test with empty mappings
	genesis.AddressMappings = []*types.AddressMapping{}
	require.NoError(genesis.Validate())
}

// Add a new test function for InitGenesis
func (suite *KeeperTestSuite) TestInitGenesis() {
	ctx, addressesKeeper := suite.ctx, suite.addressesKeeper
	require := suite.Require()

	genesis := types.GenesisState{
		Params: types.DefaultParams(),
		AddressMappings: []*types.AddressMapping{
			{CosmosAddress: cosmosAddrString, EvmAddress: evmAddrString},
		},
	}

	InitGenesis(ctx, addressesKeeper, genesis)

	// Verify that the address mapping was set
	evmAddr := common.HexToAddress(evmAddrString).Bytes()
	cosmosAddr, err := sdk.AccAddressFromBech32(cosmosAddrString)
	require.NoError(err)

	result := addressesKeeper.GetCosmosAddressByEvmAddress(ctx, evmAddr)
	require.Equal(cosmosAddr.Bytes(), result)

	result = addressesKeeper.GetEvmAddressByCosmosAddress(ctx, cosmosAddr.Bytes())
	require.Equal(evmAddr, result)

	// Verify that params were set
	params := addressesKeeper.GetParamSet(ctx)
	require.Equal(genesis.Params, params)
}
