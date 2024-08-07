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

	moduletestutil "github.com/peggyjv/sommelier/v7/testutil"
	addressTypes "github.com/peggyjv/sommelier/v7/x/addresses/types"

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

	queryClient addressTypes.QueryClient

	encCfg moduletestutil.TestEncodingConfig
}

func (suite *KeeperTestSuite) SetupTest() {
	key := sdk.NewKVStoreKey(addressTypes.StoreKey)
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

	params.Subspace(addressTypes.ModuleName)
	subSpace, found := params.GetSubspace(addressTypes.ModuleName)
	suite.Assertions.True(found)

	suite.addressesKeeper = *NewKeeper(
		encCfg.Codec,
		key,
		subSpace,
	)

	addressTypes.RegisterInterfaces(encCfg.InterfaceRegistry)

	queryHelper := baseapp.NewQueryServerTestHelper(ctx, encCfg.InterfaceRegistry)
	addressTypes.RegisterQueryServer(queryHelper, suite.addressesKeeper)
	queryClient := addressTypes.NewQueryClient(queryHelper)

	suite.queryClient = queryClient
	suite.encCfg = encCfg
}

func TestKeeperTestSuite(t *testing.T) {
	suite.Run(t, new(KeeperTestSuite))
}

func (suite *KeeperTestSuite) TestSetGetDeleteParams() {
	ctx, addressesKeeper := suite.ctx, suite.addressesKeeper
	require := suite.Require()

	params := addressTypes.DefaultParams()
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
	addressesKeeper.SetAddressMapping(ctx, cosmosAddr, evmAddr)

	// Get
	cosmosResult := addressesKeeper.GetCosmosAddressByEvmAddress(ctx, evmAddr)
	require.Equal(cosmosAddr, cosmosResult)

	evmResult := addressesKeeper.GetEvmAddressByCosmosAddress(ctx, cosmosAddr)
	require.Equal(evmAddr, evmResult)

	// Iterate
	var mappings []*addressTypes.AddressMapping
	addressesKeeper.IterateAddressMappings(ctx, func(cosmosAddr []byte, evmAddr []byte) (stop bool) {
		mapping := addressTypes.AddressMapping{
			CosmosAddress: sdk.MustBech32ifyAddressBytes("cosmos", cosmosAddr),
			EvmAddress:    common.BytesToAddress(evmAddr).Hex(),
		}
		mappings = append(mappings, &mapping)

		return false
	})

	// Delete
	addressesKeeper.DeleteAddressMapping(ctx, cosmosAddr)

	cosmosResult = addressesKeeper.GetCosmosAddressByEvmAddress(ctx, evmAddr)
	require.Nil(cosmosResult)
	evmResult = addressesKeeper.GetEvmAddressByCosmosAddress(ctx, cosmosAddr)
	require.Nil(evmResult)

	// Invalid input
	addressesKeeper.SetAddressMapping(ctx, nil, evmAddr)
}
