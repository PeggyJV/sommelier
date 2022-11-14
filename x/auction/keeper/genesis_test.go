package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/peggyjv/sommelier/v4/app/params"
	auctionTypes "github.com/peggyjv/sommelier/v4/x/auction/types"
)

// Tests Importing of as empty a genesis as possible
func (suite *KeeperTestSuite) TestImportingEmptyGenesis() {
	ctx, auctionKeeper := suite.ctx, suite.auctionKeeper
	require := suite.Require()

	testGenesis := auctionTypes.GenesisState{}

	// Canary to make sure validate basic is being run
	require.Panics(func() { InitGenesis(ctx, auctionKeeper, testGenesis) })

	testGenesis.Params = auctionTypes.DefaultParams()
	require.NotPanics(func() {
		suite.mockGetModuleAccount(ctx)
		InitGenesis(ctx, auctionKeeper, testGenesis)
	})

	activeAuctions := auctionKeeper.GetActiveAuctions(ctx)
	endedAuctions := auctionKeeper.GetEndedAuctions(ctx)
	bids := auctionKeeper.GetBids(ctx)
	tokenPrices := auctionKeeper.GetTokenPrices(ctx)
	lastAuctionID := auctionKeeper.GetLastAuctionID(ctx)
	lastBidID := auctionKeeper.GetLastBidID(ctx)

	require.Len(activeAuctions, 0)
	require.Len(endedAuctions, 0)
	require.Len(bids, 0)
	require.Len(tokenPrices, 0)
	require.Zero(lastAuctionID)
	require.Zero(lastBidID)
}

// Test Importing a populated Genesis
func (suite *KeeperTestSuite) TestImportingPopulatedGenesis() {
	ctx, auctionKeeper := suite.ctx, suite.auctionKeeper
	require := suite.Require()

	expectedEndedAuctions := make([]*auctionTypes.Auction, 1)
	auction0 := &auctionTypes.Auction{
		Id:                         uint32(1),
		StartingTokensForSale:      sdk.NewCoin("weth", sdk.NewInt(17000)),
		StartBlock:                 uint64(1),
		EndBlock:                   uint64(5),
		InitialPriceDecreaseRate:   sdk.MustNewDecFromStr("0.05"),
		CurrentPriceDecreaseRate:   sdk.MustNewDecFromStr("0.05"),
		PriceDecreaseBlockInterval: uint64(20),
		InitialUnitPriceInUsomm:    sdk.MustNewDecFromStr("20"),
		CurrentUnitPriceInUsomm:    sdk.MustNewDecFromStr("10"),
		RemainingTokensForSale:     sdk.NewCoin("weth", sdk.NewInt(10000)),
		FundingModuleAccount:       "madeUpModule1",
		ProceedsModuleAccount:      "madeUpModule1",
	}

	expectedEndedAuctions[0] = auction0

	expectedActiveAuctions := make([]*auctionTypes.Auction, 2)
	auction1 := &auctionTypes.Auction{
		Id:                         uint32(2),
		StartingTokensForSale:      sdk.NewCoin("weth", sdk.NewInt(17000)),
		StartBlock:                 uint64(6),
		EndBlock:                   uint64(0),
		InitialPriceDecreaseRate:   sdk.MustNewDecFromStr("0.05"),
		CurrentPriceDecreaseRate:   sdk.MustNewDecFromStr("0.05"),
		PriceDecreaseBlockInterval: uint64(20),
		InitialUnitPriceInUsomm:    sdk.MustNewDecFromStr("20"),
		CurrentUnitPriceInUsomm:    sdk.MustNewDecFromStr("10"),
		RemainingTokensForSale:     sdk.NewCoin("weth", sdk.NewInt(10000)),
		FundingModuleAccount:       "madeUpModule1",
		ProceedsModuleAccount:      "madeUpModule1",
	}
	auction2 := &auctionTypes.Auction{
		Id:                         uint32(3),
		StartingTokensForSale:      sdk.NewCoin("usdc", sdk.NewInt(5000)),
		StartBlock:                 uint64(5),
		EndBlock:                   uint64(0),
		InitialPriceDecreaseRate:   sdk.MustNewDecFromStr("0.1"),
		CurrentPriceDecreaseRate:   sdk.MustNewDecFromStr("0.1"),
		PriceDecreaseBlockInterval: uint64(10),
		InitialUnitPriceInUsomm:    sdk.MustNewDecFromStr("1"),
		CurrentUnitPriceInUsomm:    sdk.MustNewDecFromStr("0.9"),
		RemainingTokensForSale:     sdk.NewCoin("usdc", sdk.NewInt(2000)),
		FundingModuleAccount:       "madeUpModule2",
		ProceedsModuleAccount:      "madeUpModule2",
	}
	expectedActiveAuctions[0] = auction1
	expectedActiveAuctions[1] = auction2

	allAuctions := make([]*auctionTypes.Auction, 3)
	allAuctions[0] = expectedEndedAuctions[0]
	allAuctions[1] = expectedActiveAuctions[0]
	allAuctions[2] = expectedActiveAuctions[1]

	expectedBids := make([]*auctionTypes.Bid, 2)
	bid1 := &auctionTypes.Bid{
		Id:                        uint64(1),
		AuctionId:                 uint32(1),
		Bidder:                    cosmosAddress1,
		MaxBidInUsomm:             sdk.NewCoin(params.BaseCoinUnit, sdk.NewInt(2000)),
		SaleTokenMinimumAmount:    sdk.NewCoin("weth", sdk.NewInt(20)),
		TotalFulfilledSaleTokens:  sdk.NewCoin("weth", sdk.NewInt(100)),
		SaleTokenUnitPriceInUsomm: sdk.MustNewDecFromStr("20.0"),
		TotalUsommPaid:            sdk.NewCoin(params.BaseCoinUnit, sdk.NewInt(2000)),
	}
	bid2 := &auctionTypes.Bid{
		Id:                        uint64(2),
		AuctionId:                 uint32(2),
		Bidder:                    cosmosAddress2,
		MaxBidInUsomm:             sdk.NewCoin(params.BaseCoinUnit, sdk.NewInt(1000)),
		SaleTokenMinimumAmount:    sdk.NewCoin("usdc", sdk.NewInt(100)),
		TotalFulfilledSaleTokens:  sdk.NewCoin("usdc", sdk.NewInt(1000)),
		SaleTokenUnitPriceInUsomm: sdk.MustNewDecFromStr("1.0"),
		TotalUsommPaid:            sdk.NewCoin(params.BaseCoinUnit, sdk.NewInt(2000)),
	}
	expectedBids[0] = bid1
	expectedBids[1] = bid2

	expectedTokenPrices := make([]*auctionTypes.TokenPrice, 3)
	tokenPrice1 := &auctionTypes.TokenPrice{
		Denom:            "usdc",
		UsdPrice:         sdk.MustNewDecFromStr("1.0"),
		LastUpdatedBlock: uint64(4),
	}
	tokenPrice2 := &auctionTypes.TokenPrice{
		Denom:            "usomm",
		UsdPrice:         sdk.MustNewDecFromStr("0.2"),
		LastUpdatedBlock: uint64(2),
	}
	tokenPrice3 := &auctionTypes.TokenPrice{
		Denom:            "weth",
		UsdPrice:         sdk.MustNewDecFromStr("17.45"),
		LastUpdatedBlock: uint64(3),
	}
	expectedTokenPrices[0] = tokenPrice1
	expectedTokenPrices[1] = tokenPrice2
	expectedTokenPrices[2] = tokenPrice3

	testGenesis := auctionTypes.GenesisState{
		Params:        auctionTypes.DefaultGenesisState().Params,
		Auctions:      allAuctions,
		Bids:          expectedBids,
		TokenPrices:   expectedTokenPrices,
		LastAuctionId: uint32(2),
		LastBidId:     uint64(2),
	}

	require.NotPanics(func() {
		suite.mockGetModuleAccount(ctx)
		InitGenesis(ctx, auctionKeeper, testGenesis)
	})

	// Verify value sets
	foundActiveAuctions := auctionKeeper.GetActiveAuctions(ctx)
	foundEndedAuctions := auctionKeeper.GetEndedAuctions(ctx)
	foundBids := auctionKeeper.GetBids(ctx)
	foundTokenPrices := auctionKeeper.GetTokenPrices(ctx)
	foundLastAuctionID := auctionKeeper.GetLastAuctionID(ctx)
	foundLastBidID := auctionKeeper.GetLastBidID(ctx)

	require.Equal(expectedActiveAuctions, foundActiveAuctions)
	require.Equal(expectedEndedAuctions, foundEndedAuctions)
	require.Equal(expectedBids, foundBids)
	require.Equal(expectedTokenPrices, foundTokenPrices)
	require.Equal(uint32(2), foundLastAuctionID)
	require.Equal(uint64(2), foundLastBidID)
}

// Tests Exportng of an empty/default genesis
func (suite *KeeperTestSuite) TestExportingEmptyGenesis() {
	ctx, auctionKeeper := suite.ctx, suite.auctionKeeper
	require := suite.Require()

	auctionKeeper.setParams(ctx, auctionTypes.DefaultParams())
	exportedGenesis := ExportGenesis(ctx, auctionKeeper)

	require.Equal(auctionTypes.DefaultGenesisState(), exportedGenesis)

	require.Equal(auctionTypes.DefaultGenesisState().Params, exportedGenesis.Params)
	require.Len(exportedGenesis.Auctions, 0)
	require.Len(exportedGenesis.Bids, 0)
	require.Len(exportedGenesis.TokenPrices, 0)
	require.Zero(exportedGenesis.GetAuctions())
	require.Zero(exportedGenesis.LastBidId)
}

// Test Exporting a populated Genesis
func (suite *KeeperTestSuite) TestExportingPopulatedGenesis() {
	ctx, auctionKeeper := suite.ctx, suite.auctionKeeper
	require := suite.Require()

	auctionKeeper.setParams(ctx, auctionTypes.DefaultParams())

	expectedEndedAuctions := make([]*auctionTypes.Auction, 1)
	auction0 := &auctionTypes.Auction{
		Id:                         uint32(1),
		StartingTokensForSale:      sdk.NewCoin("weth", sdk.NewInt(17000)),
		StartBlock:                 uint64(1),
		EndBlock:                   uint64(5),
		InitialPriceDecreaseRate:   sdk.MustNewDecFromStr("0.05"),
		CurrentPriceDecreaseRate:   sdk.MustNewDecFromStr("0.05"),
		PriceDecreaseBlockInterval: uint64(20),
		InitialUnitPriceInUsomm:    sdk.MustNewDecFromStr("20"),
		CurrentUnitPriceInUsomm:    sdk.MustNewDecFromStr("10"),
		RemainingTokensForSale:     sdk.NewCoin("weth", sdk.NewInt(10000)),
		FundingModuleAccount:       "madeUpModule1",
		ProceedsModuleAccount:      "madeUpModule1",
	}
	expectedEndedAuctions[0] = auction0
	auctionKeeper.setEndedAuction(ctx, *auction0)

	expectedActiveAuctions := make([]*auctionTypes.Auction, 2)
	auction1 := &auctionTypes.Auction{
		Id:                         uint32(1),
		StartingTokensForSale:      sdk.NewCoin("weth", sdk.NewInt(17000)),
		StartBlock:                 uint64(6),
		EndBlock:                   uint64(0),
		InitialPriceDecreaseRate:   sdk.MustNewDecFromStr("0.05"),
		CurrentPriceDecreaseRate:   sdk.MustNewDecFromStr("0.05"),
		PriceDecreaseBlockInterval: uint64(20),
		InitialUnitPriceInUsomm:    sdk.MustNewDecFromStr("20"),
		CurrentUnitPriceInUsomm:    sdk.MustNewDecFromStr("10"),
		RemainingTokensForSale:     sdk.NewCoin("weth", sdk.NewInt(10000)),
		FundingModuleAccount:       "madeUpModule1",
		ProceedsModuleAccount:      "madeUpModule1",
	}
	expectedActiveAuctions[0] = auction1
	auctionKeeper.setActiveAuction(ctx, *auction1)

	allAuctions := make([]*auctionTypes.Auction, 2)
	allAuctions[0] = expectedActiveAuctions[0]
	allAuctions[1] = expectedEndedAuctions[0]
	auctionKeeper.setLastAuctionID(ctx, uint32(1))

	expectedBids := make([]*auctionTypes.Bid, 1)
	bid1 := &auctionTypes.Bid{
		Id:                        uint64(1),
		AuctionId:                 uint32(1),
		Bidder:                    cosmosAddress1,
		MaxBidInUsomm:             sdk.NewCoin(params.BaseCoinUnit, sdk.NewInt(2000)),
		SaleTokenMinimumAmount:    sdk.NewCoin("weth", sdk.NewInt(20)),
		TotalFulfilledSaleTokens:  sdk.NewCoin("weth", sdk.NewInt(100)),
		SaleTokenUnitPriceInUsomm: sdk.MustNewDecFromStr("20.0"),
		TotalUsommPaid:            sdk.NewCoin(params.BaseCoinUnit, sdk.NewInt(2000)),
	}
	expectedBids[0] = bid1
	auctionKeeper.setBid(ctx, *bid1)
	auctionKeeper.setLastBidID(ctx, uint64(1))

	expectedTokenPrices := make([]*auctionTypes.TokenPrice, 2)
	tokenPrice1 := &auctionTypes.TokenPrice{
		Denom:            "usomm",
		UsdPrice:         sdk.MustNewDecFromStr("0.2"),
		LastUpdatedBlock: uint64(2),
	}
	tokenPrice2 := &auctionTypes.TokenPrice{
		Denom:            "weth",
		UsdPrice:         sdk.MustNewDecFromStr("17.45"),
		LastUpdatedBlock: uint64(3),
	}
	expectedTokenPrices[0] = tokenPrice1
	expectedTokenPrices[1] = tokenPrice2
	auctionKeeper.setTokenPrice(ctx, *tokenPrice1)
	auctionKeeper.setTokenPrice(ctx, *tokenPrice2)

	expectedGenesis := auctionTypes.GenesisState{
		Params:        auctionTypes.DefaultGenesisState().Params,
		Auctions:      allAuctions,
		Bids:          expectedBids,
		TokenPrices:   expectedTokenPrices,
		LastAuctionId: uint32(1),
		LastBidId:     uint64(1),
	}

	exportedGenesis := ExportGenesis(ctx, auctionKeeper)

	// Verify values
	require.Equal(expectedGenesis, exportedGenesis)

	require.Equal(auctionTypes.DefaultGenesisState().Params, exportedGenesis.Params)
	require.Equal(allAuctions, exportedGenesis.Auctions)
	require.Equal(expectedBids, exportedGenesis.Bids)
	require.Equal(expectedTokenPrices, exportedGenesis.TokenPrices)
	require.Equal(uint32(1), exportedGenesis.LastAuctionId)
	require.Equal(uint64(1), exportedGenesis.LastBidId)
}
