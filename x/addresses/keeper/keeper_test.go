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

	paramskeeper "github.com/cosmos/cosmos-sdk/x/params/keeper"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
	tmtime "github.com/tendermint/tendermint/types/time"
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

	cosmosAddrString := "cosmos154d0p9xhrruhxvazumej9nq29afeura2alje4u"
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
}

//// Unhappy path tests for BeginAuction
//func (suite *KeeperTestSuite) TestUnhappyPathsForBeginAuction() {
//	ctx, addressesKeeper := suite.ctx, suite.addressesKeeper
//	require := suite.Require()
//
//	// Define basic param(s)
//	addressesParams := addressTypes.DefaultParams()
//	addressesKeeper.setParams(ctx, addressesParams)
//
//	// Setup some token prices
//	sommPrice := addressTypes.TokenPrice{Denom: params.BaseCoinUnit, Exponent: 6, UsdPrice: sdk.MustNewDecFromStr("0.01"), LastUpdatedBlock: 2}
//
//	/* #nosec */
//	saleToken := "gravity0xaaaebe6fe48e54f431b0c390cfaf0b017d09d42d"
//	saleTokenPrice := addressTypes.TokenPrice{Denom: saleToken, Exponent: 6, UsdPrice: sdk.MustNewDecFromStr("0.01"), LastUpdatedBlock: 5}
//	addressesedSaleTokens := sdk.NewCoin(saleToken, sdk.NewInt(10000))
//
//	tests := []struct {
//		name                string
//		beginAuctionRequest BeginAuctionRequest
//		expectedError       error
//		runsBefore          runsBeforeWrapper
//	}{
//		{
//			name: "Unpermissioned funder module account",
//			beginAuctionRequest: BeginAuctionRequest{
//				ctx:                        ctx,
//				startingTokensForSale:      addressesedSaleTokens,
//				initialPriceDecreaseRate:   sdk.MustNewDecFromStr("0.05"),
//				priceDecreaseBlockInterval: uint64(10),
//				fundingModuleAccount:       "cork",
//				proceedsModuleAccount:      permissionedReciever.GetName(),
//			},
//			expectedError: errorsmod.Wrapf(addressTypes.ErrUnauthorizedFundingModule, "Module Account: cork"),
//			runsBefore:    func() {},
//		},
//		{
//			name: "Unpermissioned proceeds module account",
//			beginAuctionRequest: BeginAuctionRequest{
//				ctx:                        ctx,
//				startingTokensForSale:      addressesedSaleTokens,
//				initialPriceDecreaseRate:   sdk.MustNewDecFromStr("0.05"),
//				priceDecreaseBlockInterval: uint64(10),
//				fundingModuleAccount:       permissionedFunder.GetName(),
//				proceedsModuleAccount:      "gravity",
//			},
//			expectedError: errorsmod.Wrapf(addressTypes.ErrUnauthorizedProceedsModule, "Module Account: gravity"),
//			runsBefore:    func() {},
//		},
//		{
//			name: "Starting denom price not found",
//			beginAuctionRequest: BeginAuctionRequest{
//				ctx:                        ctx,
//				startingTokensForSale:      sdk.NewCoin("anvil", sdk.NewInt(7)),
//				initialPriceDecreaseRate:   sdk.MustNewDecFromStr("0.05"),
//				priceDecreaseBlockInterval: uint64(10),
//				fundingModuleAccount:       permissionedFunder.GetName(),
//				proceedsModuleAccount:      permissionedReciever.GetName(),
//			},
//			expectedError: errorsmod.Wrapf(addressTypes.ErrCouldNotFindSaleTokenPrice, "starting amount denom: anvil"),
//			runsBefore:    func() {},
//		},
//		{
//			name: "Starting denom price update too old",
//			beginAuctionRequest: BeginAuctionRequest{
//				ctx:                        ctx.WithBlockHeight(int64(saleTokenPrice.LastUpdatedBlock) + int64(addressesParams.PriceMaxBlockAge) + 1),
//				startingTokensForSale:      addressesedSaleTokens,
//				initialPriceDecreaseRate:   sdk.MustNewDecFromStr("0.05"),
//				priceDecreaseBlockInterval: uint64(10),
//				fundingModuleAccount:       permissionedFunder.GetName(),
//				proceedsModuleAccount:      permissionedReciever.GetName(),
//			},
//			expectedError: errorsmod.Wrapf(addressTypes.ErrLastSaleTokenPriceTooOld, "starting amount denom: %s", saleToken),
//			runsBefore: func() {
//				addressesKeeper.setTokenPrice(ctx, saleTokenPrice)
//			},
//		},
//		{
//			name: "Usomm price not found",
//			beginAuctionRequest: BeginAuctionRequest{
//				ctx:                        ctx,
//				startingTokensForSale:      addressesedSaleTokens,
//				initialPriceDecreaseRate:   sdk.MustNewDecFromStr("0.05"),
//				priceDecreaseBlockInterval: uint64(10),
//				fundingModuleAccount:       permissionedFunder.GetName(),
//				proceedsModuleAccount:      permissionedReciever.GetName(),
//			},
//			expectedError: errorsmod.Wrap(addressTypes.ErrCouldNotFindSommTokenPrice, params.BaseCoinUnit),
//			runsBefore:    func() {},
//		},
//		{
//			name: "Usomm price update too old",
//			beginAuctionRequest: BeginAuctionRequest{
//				ctx:                        ctx.WithBlockHeight(int64(sommPrice.LastUpdatedBlock) + int64(addressesParams.PriceMaxBlockAge) + 1),
//				startingTokensForSale:      addressesedSaleTokens,
//				initialPriceDecreaseRate:   sdk.MustNewDecFromStr("0.05"),
//				priceDecreaseBlockInterval: uint64(10),
//				fundingModuleAccount:       permissionedFunder.GetName(),
//				proceedsModuleAccount:      permissionedReciever.GetName(),
//			},
//			expectedError: errorsmod.Wrap(addressTypes.ErrLastSommTokenPriceTooOld, params.BaseCoinUnit),
//			runsBefore: func() {
//				addressesKeeper.setTokenPrice(ctx, sommPrice)
//			},
//		},
//		{
//			name: "Validate basic canary 1 -- invalid initialPriceDecreaseRate lower bound",
//			beginAuctionRequest: BeginAuctionRequest{
//				ctx:                        ctx,
//				startingTokensForSale:      addressesedSaleTokens,
//				initialPriceDecreaseRate:   sdk.MustNewDecFromStr("0.0"),
//				priceDecreaseBlockInterval: uint64(10),
//				fundingModuleAccount:       permissionedFunder.GetName(),
//				proceedsModuleAccount:      permissionedReciever.GetName(),
//			},
//			expectedError: errorsmod.Wrapf(addressTypes.ErrInvalidInitialDecreaseRate, "Initial price decrease rate 0.000000000000000000"),
//			runsBefore:    func() {},
//		},
//		{
//			name: "Validate basic canary 2 -- invalid initialPriceDecreaseRate upper bound",
//			beginAuctionRequest: BeginAuctionRequest{
//				ctx:                        ctx,
//				startingTokensForSale:      addressesedSaleTokens,
//				initialPriceDecreaseRate:   sdk.MustNewDecFromStr("1.0"),
//				priceDecreaseBlockInterval: uint64(10),
//				fundingModuleAccount:       permissionedFunder.GetName(),
//				proceedsModuleAccount:      permissionedReciever.GetName(),
//			},
//			expectedError: errorsmod.Wrapf(addressTypes.ErrInvalidInitialDecreaseRate, "Initial price decrease rate 1.000000000000000000"),
//			runsBefore:    func() {},
//		},
//		{
//			name: "Cannot have 2 ongoing addressess for the same denom",
//			beginAuctionRequest: BeginAuctionRequest{
//				ctx:                        ctx,
//				startingTokensForSale:      addressesedSaleTokens,
//				initialPriceDecreaseRate:   sdk.MustNewDecFromStr("0.05"),
//				priceDecreaseBlockInterval: uint64(10),
//				fundingModuleAccount:       permissionedFunder.GetName(),
//				proceedsModuleAccount:      permissionedReciever.GetName(),
//			},
//			expectedError: errorsmod.Wrapf(addressTypes.ErrCannotStartTwoAuctionsForSameDenomSimultaneously, "Denom: %s", addressesedSaleTokens.Denom),
//			runsBefore: func() {
//				// Mock initial bank keeper fund transfer
//				suite.mockSendCoinsFromModuleToModule(ctx, permissionedFunder.GetName(), addressTypes.ModuleName, sdk.NewCoins(addressesedSaleTokens))
//
//				// Start addresses
//				decreaseRate := sdk.MustNewDecFromStr("0.05")
//				blockDecreaseInterval := uint64(5)
//				err := addressesKeeper.BeginAuction(ctx, addressesedSaleTokens, decreaseRate, blockDecreaseInterval, permissionedFunder.GetName(), permissionedReciever.GetName())
//				require.Nil(err)
//
//				// Verify addresses got added to active addresses store
//				addressesID := uint32(1)
//				createdAuction, found := addressesKeeper.GetActiveAuctionByID(ctx, addressesID)
//				require.True(found)
//			},
//		},
//	}
//
//	for _, tc := range tests {
//		tc := tc // Redefine variable here due to passing it to function literal below (scopelint)
//		suite.T().Run(fmt.Sprint(tc.name), func(t *testing.T) {
//			// Run expected bank keeper functions, if any
//			tc.runsBefore()
//
//			call := func() error {
//				return addressesKeeper.BeginAuction(
//					tc.beginAuctionRequest.ctx,
//					tc.beginAuctionRequest.startingTokensForSale,
//					tc.beginAuctionRequest.initialPriceDecreaseRate,
//					tc.beginAuctionRequest.priceDecreaseBlockInterval,
//					tc.beginAuctionRequest.fundingModuleAccount,
//					tc.beginAuctionRequest.proceedsModuleAccount,
//				)
//			}
//
//			if tc.name[0:14] == "Validate basic" {
//				require.Panics(func() { call() })
//				return
//			}
//
//			err := call()
//
//			// Verify errors are as expected
//			require.Equal(tc.expectedError.Error(), err.Error())
//		})
//	}
//}
