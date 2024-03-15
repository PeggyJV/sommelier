package keeper

import "github.com/peggyjv/sommelier/v7/x/addresses/types"

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
	genesis.AddressMappings = append(genesis.AddressMappings, &types.AddressMapping{CosmosAddress: "cosmos1l8n6v5f4j5s8j5l8n6v5f4j5s8j5l8n6v5f4j", EvmAddress: "zzzz"})
	require.Error(genesis.Validate())
}
