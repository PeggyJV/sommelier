package integration_tests

import (
	"context"
	"time"

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
		s.Require().Eventually(func() bool {
			queryRes, err := addressesQueryClient.QueryAddressMappings(context.Background(), &types.QueryAddressMappingsRequest{})
			if err != nil {
				s.T().Logf("Error querying address mappings: %s", err)
				return false
			}

			return len(queryRes.AddressMappings) == 1 && evmAddress.Hex() == queryRes.AddressMappings[0].EvmAddress
		}, 20*time.Second, 4*time.Second, "address mapping never found")

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
		s.Require().Eventually(func() bool {
			queryRes, err := addressesQueryClient.QueryAddressMappings(context.Background(), &types.QueryAddressMappingsRequest{})
			if err != nil {
				s.T().Logf("Error querying address mappings: %s", err)
				return false
			}

			return len(queryRes.AddressMappings) == 0
		}, 20*time.Second, 4*time.Second, "address mapping not deleted")

		_, err = addressesQueryClient.QueryAddressMappingByEVMAddress(context.Background(), &types.QueryAddressMappingByEVMAddressRequest{EvmAddress: evmAddress.Hex()})
		s.Require().Error(err)
		s.Require().Contains(err.Error(), "code = NotFound")

		_, err = addressesQueryClient.QueryAddressMappingByCosmosAddress(context.Background(), &types.QueryAddressMappingByCosmosAddressRequest{CosmosAddress: cosmosAddress.String()})
		s.Require().Error(err)
		s.Require().Contains(err.Error(), "code = NotFound")

		// Test error cases

		// Test adding multiple mappings
		s.T().Log("Adding multiple mappings")
		evmAddress2 := common.HexToAddress("0x2345678901234567890123456789012345678901")
		cosmosAddress2 := s.chain.orchestrators[1].address()
		orch1 := s.chain.orchestrators[1]
		orch1ClientCtx, err := s.chain.clientContext("tcp://localhost:26657", orch1.keyring, "orch", orch1.address())
		s.Require().NoError(err)

		addAddressMappingMsg2, err := types.NewMsgAddAddressMapping(evmAddress2, cosmosAddress2)
		s.Require().NoError(err)

		_, err = s.chain.sendMsgs(*orchClientCtx, addAddressMappingMsg)
		s.Require().NoError(err)
		_, err = s.chain.sendMsgs(*orch1ClientCtx, addAddressMappingMsg2)
		s.Require().NoError(err)

		// Query multiple mappings
		s.T().Log("Querying multiple mappings")
		s.Require().Eventually(func() bool {
			queryRes, err := addressesQueryClient.QueryAddressMappings(context.Background(), &types.QueryAddressMappingsRequest{})
			if err != nil {
				s.T().Logf("Error querying address mappings: %s", err)
				return false
			}

			return len(queryRes.AddressMappings) == 2
		}, 20*time.Second, 4*time.Second, "expected two address mappings")

		// Test adding a duplicate mapping
		s.T().Log("Adding a duplicate mapping")
		duplicateMsg, err := types.NewMsgAddAddressMapping(evmAddress, cosmosAddress)
		s.Require().NoError(err)

		_, err = s.chain.sendMsgs(*orchClientCtx, duplicateMsg)
		s.Require().NoError(err)
		_, err = s.chain.sendMsgs(*orchClientCtx, duplicateMsg)
		s.Require().NoError(err)

		// Test removing a non-existent mapping
		s.T().Log("Removing a non-existent mapping")
		nonExistentAddress := s.chain.orchestrators[2].address()
		orch2 := s.chain.orchestrators[2]
		orch2ClientCtx, err := s.chain.clientContext("tcp://localhost:26657", orch2.keyring, "orch", orch2.address())
		removeNonExistentMsg, err := types.NewMsgRemoveAddressMapping(nonExistentAddress)
		s.Require().NoError(err)

		_, err = s.chain.sendMsgs(*orch2ClientCtx, removeNonExistentMsg)
		s.Require().NoError(err)

		// Test querying with invalid addresses
		s.T().Log("Querying with invalid addresses")
		_, err = addressesQueryClient.QueryAddressMappingByEVMAddress(context.Background(), &types.QueryAddressMappingByEVMAddressRequest{EvmAddress: "invalid"})
		s.Require().Error(err)
		s.Require().Contains(err.Error(), "invalid EVM address")

		s.T().Log("Querying with invalid cosmos address")
		_, err = addressesQueryClient.QueryAddressMappingByCosmosAddress(context.Background(), &types.QueryAddressMappingByCosmosAddressRequest{CosmosAddress: "invalid"})
		s.Require().Error(err)
		s.Require().Contains(err.Error(), "failed to parse cosmos address")
	})
}
