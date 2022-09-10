package keeper_test

import (
	"testing"

	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/testutil"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/golang/mock/gomock"

	moduletestutil "github.com/peggyjv/sommelier/v4/testutil"
	auctionKeeper "github.com/peggyjv/sommelier/v4/x/auction/keeper"
	auctionTypes "github.com/peggyjv/sommelier/v4/x/auction/types"

	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	paramskeeper "github.com/cosmos/cosmos-sdk/x/params/keeper"
	auctiontestutil "github.com/peggyjv/sommelier/v4/x/auction/testutil"
	"github.com/stretchr/testify/suite"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
	tmtime "github.com/tendermint/tendermint/types/time"
)

var (
	permissionedFunder   = authtypes.NewEmptyModuleAccount("permissionedFunder")
	permissionedReciever = authtypes.NewEmptyModuleAccount("permissionedReciever")
)

type KeeperTestSuite struct {
	suite.Suite

	ctx           sdk.Context
	auctionKeeper auctionKeeper.Keeper
	bankKeeper    *auctiontestutil.MockBankKeeper

	queryClient auctionTypes.QueryClient

	encCfg moduletestutil.TestEncodingConfig
}

func (suite *KeeperTestSuite) SetupTest() {
	key := sdk.NewKVStoreKey(auctionTypes.StoreKey)
	tkey := sdk.NewTransientStoreKey("transient_test")
	testCtx := testutil.DefaultContext(key, tkey)
	ctx := testCtx.WithBlockHeader(tmproto.Header{Height: 5, Time: tmtime.Now()})
	encCfg := moduletestutil.MakeTestEncodingConfig()

	// gomock initializations
	ctrl := gomock.NewController(suite.T())
	suite.bankKeeper = auctiontestutil.NewMockBankKeeper(ctrl)
	suite.ctx = ctx

	params := paramskeeper.NewKeeper(
		encCfg.Codec,
		codec.NewLegacyAmino(),
		key,
		tkey,
	)

	params.Subspace(auctionTypes.ModuleName)
	subSpace, found := params.GetSubspace(auctionTypes.ModuleName)
	suite.Assertions.True(found)

	suite.auctionKeeper = auctionKeeper.NewKeeper(
		encCfg.Codec,
		key,
		subSpace,
		suite.bankKeeper,
		map[string]bool{permissionedFunder.GetName(): true},
		map[string]bool{permissionedReciever.GetName(): true},
	)

	auctionTypes.RegisterInterfaces(encCfg.InterfaceRegistry)

	queryHelper := baseapp.NewQueryServerTestHelper(ctx, encCfg.InterfaceRegistry)
	auctionTypes.RegisterQueryServer(queryHelper, suite.auctionKeeper)
	queryClient := auctionTypes.NewQueryClient(queryHelper)

	suite.queryClient = queryClient
	suite.encCfg = encCfg
}

func TestKeeperTestSuite(t *testing.T) {
	suite.Run(t, new(KeeperTestSuite))
}

func (suite *KeeperTestSuite) mockSendCoinsFromModuleToModule(ctx sdk.Context, sender string, receiver string, amt sdk.Coins) {
	suite.bankKeeper.EXPECT().SendCoinsFromModuleToModule(ctx, sender, receiver, amt).Return(nil)
}

// Happy path for BeginAuction call
func (suite *KeeperTestSuite) TestValidBeginAuction() {
	ctx, auctionKeeper := suite.ctx, suite.auctionKeeper
	require := suite.Require()

	params := auctionTypes.Params{PriceMaxBlockAge: 10}
	auctionKeeper.SetParams(ctx, params)

	sommPrice := auctionTypes.TokenPrice{Denom: "usomm", UsdPrice: sdk.MustNewDecFromStr("0.01"), LastUpdatedBlock: 5}

	saleToken := "gravity0xdac17f958d2ee523a2206206994597c13d831ec7"
	saleTokenPrice := auctionTypes.TokenPrice{Denom: saleToken, UsdPrice: sdk.MustNewDecFromStr("0.02"), LastUpdatedBlock: 5}
	auctionedSaleTokens := sdk.NewCoin(saleToken, sdk.NewInt(10000))

	auctionKeeper.SetTokenPrice(ctx, sommPrice)
	auctionKeeper.SetTokenPrice(ctx, saleTokenPrice)

	// Mock bank keeper fund transfer
	suite.mockSendCoinsFromModuleToModule(ctx, permissionedFunder.GetName(), auctionTypes.ModuleName, sdk.NewCoins(auctionedSaleTokens))

	// Start auction
	decreaseRate := sdk.MustNewDecFromStr("0.05")
	blockDecreaseInterval := uint64(5)
	err := auctionKeeper.BeginAuction(ctx, auctionedSaleTokens, decreaseRate, blockDecreaseInterval, permissionedFunder.GetName(), permissionedReciever.GetName())
	require.Nil(err)

	// Verify auction got added to active auction store
	auctionId := uint32(1)
	createdAuction, found := auctionKeeper.GetActiveAuctionByID(ctx, auctionId)
	require.True(found)

	expectedAuction := auctionTypes.Auction{
		Id:                         auctionId,
		StartingTokensForSale:      auctionedSaleTokens,
		StartBlock:                 uint64(ctx.BlockHeight()),
		EndBlock:                   0,
		InitialPriceDecreaseRate:   decreaseRate,
		CurrentPriceDecreaseRate:   decreaseRate,
		PriceDecreaseBlockInterval: blockDecreaseInterval,
		InitialUnitPriceInUsomm:    sdk.NewDec(2),
		CurrentUnitPriceInUsomm:    sdk.NewDec(2),
		RemainingTokensForSale:     auctionedSaleTokens,
		FundingModuleAccount:       permissionedFunder.GetName(),
		ProceedsModuleAccount:      permissionedReciever.GetName(),
	}

	require.Equal(expectedAuction, createdAuction)
}
