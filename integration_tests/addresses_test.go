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
	})
}
