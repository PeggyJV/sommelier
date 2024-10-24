package keeper

import (
	"testing"

	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/testutil"
	"github.com/ethereum/go-ethereum/common"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/suite"

	sdk "github.com/cosmos/cosmos-sdk/types"

	moduletestutil "github.com/peggyjv/sommelier/v8/testutil"
	"github.com/peggyjv/sommelier/v8/x/addresses/types"

	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"
	tmtime "github.com/cometbft/cometbft/types/time"
	paramskeeper "github.com/cosmos/cosmos-sdk/x/params/keeper"
)

const (
	cosmosAddrString        = "cosmos154d0p9xhrruhxvazumej9nq29afeura2alje4u"
	cosmosAddrStringInvalid = "cosmos154d0p9xhrruhxvazumej9nq29afeura2alje41"
	evmAddrString           = "0x1111111111111111111111111111111111111111"
)

type KeeperTestSuite struct {
	suite.Suite

	ctx             sdk.Context
	addressesKeeper Keeper

	queryClient types.QueryClient

	encCfg moduletestutil.TestEncodingConfig
}

func (suite *KeeperTestSuite) SetupTest() {
	key := sdk.NewKVStoreKey(types.StoreKey)
	tkey := sdk.NewTransientStoreKey("transient_test")
	testCtx := testutil.DefaultContext(key, tkey)
	ctx := testCtx.WithBlockHeader(tmproto.Header{Height: 5, Time: tmtime.Now()})
	encCfg := moduletestutil.MakeTestEncodingConfig()

	// gomock initializations
	ctrl := gomock.NewController(suite.T())
	defer ctrl.Finish()

	suite.ctx = ctx

	params := paramskeeper.NewKeeper(
		encCfg.Codec,
		codec.NewLegacyAmino(),
		key,
		tkey,
	)

	params.Subspace(types.ModuleName)
	subSpace, found := params.GetSubspace(types.ModuleName)
	suite.Assertions.True(found)

	suite.addressesKeeper = *NewKeeper(
		encCfg.Codec,
		key,
		subSpace,
	)

	types.RegisterInterfaces(encCfg.InterfaceRegistry)

	queryHelper := baseapp.NewQueryServerTestHelper(ctx, encCfg.InterfaceRegistry)
	types.RegisterQueryServer(queryHelper, suite.addressesKeeper)
	queryClient := types.NewQueryClient(queryHelper)

	suite.queryClient = queryClient
	suite.encCfg = encCfg
}

func TestKeeperTestSuite(t *testing.T) {
	suite.Run(t, new(KeeperTestSuite))
}

func (suite *KeeperTestSuite) TestSetGetDeleteParams() {
	ctx, addressesKeeper := suite.ctx, suite.addressesKeeper
	require := suite.Require()

	params := types.DefaultParams()
	addressesKeeper.setParams(ctx, params)

	retrievedParams := addressesKeeper.GetParamSet(ctx)
	require.Equal(params, retrievedParams)
}

func (suite *KeeperTestSuite) TestSetGetDeleteAddressMappings() {
	ctx, addressesKeeper := suite.ctx, suite.addressesKeeper
	require := suite.Require()

	evmAddrString := "0x1111111111111111111111111111111111111111"
	require.Equal(42, len(evmAddrString))
	evmAddr := common.HexToAddress(evmAddrString).Bytes()

	acc, err := sdk.AccAddressFromBech32(cosmosAddrString)
	require.NoError(err)

	cosmosAddr := acc.Bytes()

	// Set
	err = addressesKeeper.SetAddressMapping(ctx, cosmosAddr, evmAddr)
	require.NoError(err)

	// Get
	cosmosResult := addressesKeeper.GetCosmosAddressByEvmAddress(ctx, evmAddr)
	require.Equal(cosmosAddr, cosmosResult)

	evmResult := addressesKeeper.GetEvmAddressByCosmosAddress(ctx, cosmosAddr)
	require.Equal(evmAddr, evmResult)

	// Iterate
	var mappings []*types.AddressMapping
	addressesKeeper.IterateAddressMappings(ctx, func(cosmosAddr []byte, evmAddr []byte) (stop bool) {
		mapping := types.AddressMapping{
			CosmosAddress: sdk.MustBech32ifyAddressBytes("cosmos", cosmosAddr),
			EvmAddress:    common.BytesToAddress(evmAddr).Hex(),
		}
		mappings = append(mappings, &mapping)

		return false
	})

	// Invalid input
	err = addressesKeeper.SetAddressMapping(ctx, nil, evmAddr)
	require.Error(err)

	// Test setting multiple mappings
	evmAddr2 := common.HexToAddress("0x2222222222222222222222222222222222222222").Bytes()
	cosmosAddr2, _ := sdk.AccAddressFromBech32("cosmos1y6d5kasehecexf09ka6y0ggl0pxzt6dgk0gnl9")

	err = addressesKeeper.SetAddressMapping(ctx, cosmosAddr2, evmAddr2)
	require.NoError(err)

	// Verify second mapping exists
	result := addressesKeeper.GetCosmosAddressByEvmAddress(ctx, evmAddr2)
	require.Equal(cosmosAddr2.Bytes(), result)

	result = addressesKeeper.GetEvmAddressByCosmosAddress(ctx, cosmosAddr2)
	require.Equal(evmAddr2, result)

	// Test deleting one mapping
	err = addressesKeeper.DeleteAddressMapping(ctx, cosmosAddr)
	require.NoError(err)

	// Verify the deleted mapping is gone but the other remains
	result = addressesKeeper.GetCosmosAddressByEvmAddress(ctx, evmAddr)
	require.Nil(result)
	result = addressesKeeper.GetEvmAddressByCosmosAddress(ctx, cosmosAddr)
	require.Nil(result)

	result = addressesKeeper.GetCosmosAddressByEvmAddress(ctx, evmAddr2)
	require.Equal(cosmosAddr2.Bytes(), result)
	result = addressesKeeper.GetEvmAddressByCosmosAddress(ctx, cosmosAddr2)
	require.Equal(evmAddr2, result)

	// Test IterateAddressMappings
	var mappings2 []*types.AddressMapping
	addressesKeeper.IterateAddressMappings(ctx, func(cosmosAddr []byte, evmAddr []byte) bool {
		mappings2 = append(mappings2, &types.AddressMapping{
			CosmosAddress: sdk.AccAddress(cosmosAddr).String(),
			EvmAddress:    common.BytesToAddress(evmAddr).Hex(),
		})
		return false
	})
	require.Len(mappings2, 1)
	require.Equal(cosmosAddr2.String(), mappings2[0].CosmosAddress)
	require.Equal(common.BytesToAddress(evmAddr2).Hex(), mappings2[0].EvmAddress)

	// Delete second mapping
	err = addressesKeeper.DeleteAddressMapping(ctx, cosmosAddr2)
	require.NoError(err)

	cosmosResult = addressesKeeper.GetCosmosAddressByEvmAddress(ctx, evmAddr2)
	require.Nil(cosmosResult)
	evmResult = addressesKeeper.GetEvmAddressByCosmosAddress(ctx, cosmosAddr2)
	require.Nil(evmResult)
}
