package integration_tests

import (
	"context"

	"github.com/ethereum/go-ethereum/common"
	"github.com/peggyjv/sommelier/v7/x/addresses/types"
)

func (s *IntegrationTestSuite) TestAddresses() {
	s.Run("Bring up chain, submit, query, and remove address mappings", func() {
		s.T().Log("Starting x/addresses tests")
		val := s.chain.validators[0]
		kb, err := val.keyring()
		s.Require().NoError(err)
		val0ClientCtx, err := s.chain.clientContext("tcp://localhost:26657", &kb, "val", val.address())
		s.Require().NoError(err)
		addressesQueryClient := types.NewQueryClient(val0ClientCtx)

		orch := s.chain.orchestrators[0]
		orchClientCtx, err := s.chain.clientContext("tcp://localhost:26657", orch.keyring, "orch", orch.address())
		s.Require().NoError(err)

		evmAddress := common.HexToAddress("0x1234567890123456789012345678901234567890")
		cosmosAddress := orch.address()

		addAddressMappingMsg, err := types.NewMsgAddAddressMapping(evmAddress, cosmosAddress)
		s.Require().NoError(err)

		s.T().Logf("Submitting mapping from %s to %s", evmAddress.Hex(), cosmosAddress.String())
		_, err = s.chain.sendMsgs(*orchClientCtx, addAddressMappingMsg)
		s.Require().NoError(err)

		s.T().Log("Testing queries return expected addresses")
		queryRes, err := addressesQueryClient.QueryAddressMappings(context.Background(), &types.QueryAddressMappingsRequest{})
		s.Require().NoError(err)
		s.Require().Len(queryRes.AddressMappings, 1, "There should be one address mapping")
		s.Require().Equal(evmAddress.Hex(), queryRes.AddressMappings[0].EvmAddress, "EVM address does not match")

		queryByEvmRes, err := addressesQueryClient.QueryAddressMappingByEVMAddress(context.Background(), &types.QueryAddressMappingByEVMAddressRequest{EvmAddress: evmAddress.Hex()})
		s.Require().NoError(err)
		s.Require().Equal(cosmosAddress.String(), queryByEvmRes.CosmosAddress, "Cosmos address does not match")
		s.Require().Equal(evmAddress.Hex(), queryByEvmRes.EvmAddress, "EVM address does not match")

		queryByCosmosRes, err := addressesQueryClient.QueryAddressMappingByCosmosAddress(context.Background(), &types.QueryAddressMappingByCosmosAddressRequest{CosmosAddress: cosmosAddress.String()})
		s.Require().NoError(err)
		s.Require().Equal(cosmosAddress.String(), queryByCosmosRes.CosmosAddress, "Cosmos address does not match")
		s.Require().Equal(evmAddress.Hex(), queryByCosmosRes.EvmAddress, "EVM address does not match")

		s.T().Log("Removing address mapping")
		removeAddressMappingMsg, err := types.NewMsgRemoveAddressMapping(cosmosAddress)
		s.Require().NoError(err)

		_, err = s.chain.sendMsgs(*orchClientCtx, removeAddressMappingMsg)
		s.Require().NoError(err)

		s.T().Log("Testing mappings query returns no addresses")
		queryRes, err = addressesQueryClient.QueryAddressMappings(context.Background(), &types.QueryAddressMappingsRequest{})
		s.Require().NoError(err)
		s.Require().Len(queryRes.AddressMappings, 0, "There should be no address mappings")
	})
}
