package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/peggyjv/sommelier/v7/x/addresses/types"
)

func (suite *KeeperTestSuite) TestHappyPathsForMsgServer() {
	ctx, addressesKeeper := suite.ctx, suite.addressesKeeper
	require := suite.Require()

	// Test AddAddressMapping
	_, err := sdk.AccAddressFromBech32(cosmosAddrString)
	require.NoError(err)

	_, err = addressesKeeper.AddAddressMapping(sdk.WrapSDKContext(ctx), &types.MsgAddAddressMapping{Signer: cosmosAddrString, EvmAddress: evmAddrString})
	require.NoError(err)

	evmAddr := common.HexToAddress(evmAddrString).Bytes()

	result := addressesKeeper.GetCosmosAddressByEvmAddress(ctx, evmAddr)
	require.NotNil(result)
	actualCosmosAddrString := sdk.AccAddress(result).String()
	require.Equal(cosmosAddrString, actualCosmosAddrString)

	result = addressesKeeper.GetEvmAddressByCosmosAddress(ctx, sdk.AccAddress(result).Bytes())
	require.NotNil(result)

	actualEvmAddrString := common.BytesToAddress(result).Hex()
	require.Equal(evmAddrString, actualEvmAddrString)

	// Test RemoveAddressMapping
	_, err = addressesKeeper.RemoveAddressMapping(sdk.WrapSDKContext(ctx), &types.MsgRemoveAddressMapping{Signer: cosmosAddrString})
	require.NoError(err)

	result = addressesKeeper.GetCosmosAddressByEvmAddress(ctx, evmAddr)
	require.Nil(result)

	result = addressesKeeper.GetEvmAddressByCosmosAddress(ctx, sdk.AccAddress(result).Bytes())
	require.Nil(result)
}

func (suite *KeeperTestSuite) TestUnhappyPathsForMsgServer() {
	ctx, addressesKeeper := suite.ctx, suite.addressesKeeper
	require := suite.Require()

	// Test AddAddressMapping
	// too long evm address
	evmAddrStringInvalid := "0x11111111111111111111111111111111111111111"
	_, err := sdk.AccAddressFromBech32(cosmosAddrStringInvalid)
	require.Error(err)

	_, err = addressesKeeper.AddAddressMapping(sdk.WrapSDKContext(ctx), &types.MsgAddAddressMapping{Signer: cosmosAddrStringInvalid, EvmAddress: evmAddrString})
	require.Error(err)
	require.Contains(err.Error(), "invalid signer address")

	_, err = sdk.AccAddressFromBech32(cosmosAddrString)
	require.NoError(err)
	evmAddr := common.HexToAddress(evmAddrString).Bytes()

	_, err = addressesKeeper.AddAddressMapping(sdk.WrapSDKContext(ctx), &types.MsgAddAddressMapping{Signer: cosmosAddrString, EvmAddress: evmAddrStringInvalid})
	require.Error(err)
	require.Contains(err.Error(), "invalid EVM address")

	result := addressesKeeper.GetCosmosAddressByEvmAddress(ctx, evmAddr)
	require.Nil(result)

	result = addressesKeeper.GetEvmAddressByCosmosAddress(ctx, sdk.AccAddress(result).Bytes())
	require.Nil(result)

	// Test RemoveAddressMapping
	_, err = sdk.AccAddressFromBech32(cosmosAddrStringInvalid)
	require.Error(err)

	_, err = addressesKeeper.RemoveAddressMapping(sdk.WrapSDKContext(ctx), &types.MsgRemoveAddressMapping{Signer: cosmosAddrStringInvalid})
	require.Error(err)
	require.Contains(err.Error(), "invalid signer address")
}
