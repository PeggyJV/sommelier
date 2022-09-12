package keeper

import (
	"testing"

	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/testutil"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/golang/mock/gomock"

	moduletestutil "github.com/peggyjv/sommelier/v4/testutil"
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
	auctionKeeper Keeper
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

	suite.auctionKeeper = NewKeeper(
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

func (suite *KeeperTestSuite) mockGetBalance(ctx sdk.Context, addr sdk.AccAddress, denom string, expectedOutput sdk.Coin) {
	suite.bankKeeper.EXPECT().GetBalance(ctx, addr, denom).Return(expectedOutput)
}

// Happy path for BeginAuction call
func (suite *KeeperTestSuite) TestHappyPathBeginAuction() {
	ctx, auctionKeeper := suite.ctx, suite.auctionKeeper
	require := suite.Require()

	params := auctionTypes.Params{PriceMaxBlockAge: 10}
	auctionKeeper.setParams(ctx, params)

	sommPrice := auctionTypes.TokenPrice{Denom: "usomm", UsdPrice: sdk.MustNewDecFromStr("0.01"), LastUpdatedBlock: 5}

	/* #nosec */
	saleToken := "gravity0xdac17f958d2ee523a2206206994597c13d831ec7"
	saleTokenPrice := auctionTypes.TokenPrice{Denom: saleToken, UsdPrice: sdk.MustNewDecFromStr("0.02"), LastUpdatedBlock: 5}
	auctionedSaleTokens := sdk.NewCoin(saleToken, sdk.NewInt(10000))

	auctionKeeper.setTokenPrice(ctx, sommPrice)
	auctionKeeper.setTokenPrice(ctx, saleTokenPrice)

	// Mock bank keeper fund transfer
	suite.mockSendCoinsFromModuleToModule(ctx, permissionedFunder.GetName(), auctionTypes.ModuleName, sdk.NewCoins(auctionedSaleTokens))

	// Start auction
	decreaseRate := sdk.MustNewDecFromStr("0.05")
	blockDecreaseInterval := uint64(5)
	err := auctionKeeper.BeginAuction(ctx, auctionedSaleTokens, decreaseRate, blockDecreaseInterval, permissionedFunder.GetName(), permissionedReciever.GetName())
	require.Nil(err)

	// Verify auction got added to active auction store
	auctionID := uint32(1)
	createdAuction, found := auctionKeeper.GetActiveAuctionByID(ctx, auctionID)
	require.True(found)

	expectedActiveAuction := auctionTypes.Auction{
		Id:                         auctionID,
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

	require.Equal(expectedActiveAuction, createdAuction)
}

// Happy path for FinishAuction (with some remaining funds)
func (suite *KeeperTestSuite) TestHappyPathFinishAuction() {
	ctx, auctionKeeper := suite.ctx, suite.auctionKeeper
	require := suite.Require()

	// 1. --------> Create an auction first so we can finish it
	params := auctionTypes.Params{PriceMaxBlockAge: 10}
	auctionKeeper.setParams(ctx, params)

	sommPrice := auctionTypes.TokenPrice{Denom: "usomm", UsdPrice: sdk.MustNewDecFromStr("0.02"), LastUpdatedBlock: 2}

	/* #nosec */
	saleToken := "gravity0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48"
	saleTokenPrice := auctionTypes.TokenPrice{Denom: saleToken, UsdPrice: sdk.MustNewDecFromStr("0.01"), LastUpdatedBlock: 2}
	auctionedSaleTokens := sdk.NewCoin(saleToken, sdk.NewInt(10000))

	auctionKeeper.setTokenPrice(ctx, sommPrice)
	auctionKeeper.setTokenPrice(ctx, saleTokenPrice)

	// Mock bank keeper fund transfer
	suite.mockSendCoinsFromModuleToModule(ctx, permissionedFunder.GetName(), auctionTypes.ModuleName, sdk.NewCoins(auctionedSaleTokens))

	// Start auction
	decreaseRate := sdk.MustNewDecFromStr("0.05")
	blockDecreaseInterval := uint64(5)
	err := auctionKeeper.BeginAuction(ctx, auctionedSaleTokens, decreaseRate, blockDecreaseInterval, permissionedFunder.GetName(), permissionedReciever.GetName())
	require.Nil(err)

	// 2. --------> We can now attempt to finish the created auction
	auctionID := uint32(1)
	createdAuction, found := auctionKeeper.GetActiveAuctionByID(ctx, auctionID)
	require.True(found)

	// Mock bank keeper balance and transfers (say only 75% got sold (25% remaining) to test funder returns & proceeds transfers)
	remainingSaleTokens := sdk.NewCoin(saleToken, auctionedSaleTokens.Amount.Quo(sdk.NewInt(4)))
	suite.mockGetBalance(ctx, authtypes.NewModuleAddress(auctionTypes.ModuleName), saleToken, remainingSaleTokens)

	// First transfer to return funding tokens
	suite.mockSendCoinsFromModuleToModule(ctx, auctionTypes.ModuleName, permissionedFunder.GetName(), sdk.NewCoins(remainingSaleTokens))

	// Add a couple of fake bids into store (note none of these fields matter for this test aside from TotalUsommPaid)
	amountPaid1 := sdk.NewCoin("usomm", sdk.NewInt(2500))
	auctionKeeper.setBid(ctx, auctionTypes.Bid{
		Id:                       1,
		AuctionId:                1,
		Bidder:                   "bidder1",
		MaxBidInUsomm:            sdk.NewCoin("usomm", sdk.NewInt(2500)),
		SaleTokenMinimumAmount:   sdk.NewCoin(saleToken, sdk.NewInt(0)),
		TotalFulfilledSaleTokens: sdk.NewCoin(saleToken, sdk.NewInt(5000)),
		TotalUsommPaid:           amountPaid1,
	})

	amountPaid2 := sdk.NewCoin("usomm", sdk.NewInt(1250))
	auctionKeeper.setBid(ctx, auctionTypes.Bid{
		Id:                       2,
		AuctionId:                1,
		Bidder:                   "bidder2",
		MaxBidInUsomm:            sdk.NewCoin("usomm", sdk.NewInt(1250)),
		SaleTokenMinimumAmount:   sdk.NewCoin(saleToken, sdk.NewInt(0)),
		TotalFulfilledSaleTokens: sdk.NewCoin(saleToken, sdk.NewInt(2500)),
		TotalUsommPaid:           amountPaid2,
	})

	// Second transfer to return proceeds from bids
	totalUsommExpected := sdk.NewCoin("usomm", amountPaid1.Amount.Add(amountPaid2.Amount))
	suite.mockSendCoinsFromModuleToModule(ctx, auctionTypes.ModuleName, permissionedReciever.GetName(), sdk.NewCoins(totalUsommExpected))

	// Change active auction tokens remaining before finishing auction to pretend tokens were sold
	createdAuction.RemainingTokensForSale = remainingSaleTokens
	auctionKeeper.setEndedAuction(ctx, createdAuction)

	// Finally actually finish the auction
	auctionKeeper.FinishAuction(ctx, &createdAuction)

	// Verify actual ended auction equals expected one
	expectedEndedAuction := auctionTypes.Auction{
		Id:                         auctionID,
		StartingTokensForSale:      auctionedSaleTokens,
		StartBlock:                 createdAuction.StartBlock,
		EndBlock:                   uint64(ctx.BlockHeight()),
		InitialPriceDecreaseRate:   decreaseRate,
		CurrentPriceDecreaseRate:   decreaseRate, // Monotonic decrease rates for now
		PriceDecreaseBlockInterval: blockDecreaseInterval,
		InitialUnitPriceInUsomm:    sdk.MustNewDecFromStr("0.5"),
		CurrentUnitPriceInUsomm:    sdk.MustNewDecFromStr("0.5"), // Started and ended on the same block
		RemainingTokensForSale:     remainingSaleTokens,
		FundingModuleAccount:       permissionedFunder.GetName(),
		ProceedsModuleAccount:      permissionedReciever.GetName(),
	}
	actualEndedAuction, found := auctionKeeper.GetEndedAuctionByID(ctx, auctionID)
	require.True(found)

	require.Equal(expectedEndedAuction, actualEndedAuction)

	// Make sure no active auctions exist anymore
	require.Zero(len(auctionKeeper.GetActiveAuctions(ctx)))
}
