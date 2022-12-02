package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	auctionTypes "github.com/peggyjv/sommelier/v4/x/auction/types"
)

// Tests abci
func (suite *KeeperTestSuite) TestAbci() {
	ctx, auctionKeeper := suite.ctx, suite.auctionKeeper
	require := suite.Require()

	params := auctionTypes.DefaultParams()
	auctionKeeper.setParams(ctx, params)

	// Test base case of no auctions
	// Note BeginBlocker is only run once for completeness, since it has no code in it
	require.NotPanics(func() { auctionKeeper.BeginBlocker(ctx) })
	require.NotPanics(func() { auctionKeeper.EndBlocker(ctx) })

	// Create an auction
	sommPrice := auctionTypes.TokenPrice{Denom: auctionTypes.UsommDenom, UsdPrice: sdk.MustNewDecFromStr("0.01"), LastUpdatedBlock: 5}

	/* #nosec */
	saleToken := "gravity0x95aD61b0a150d79219dCF64E1E6Cc01f0B64C4cE"
	saleTokenPrice := auctionTypes.TokenPrice{Denom: saleToken, UsdPrice: sdk.MustNewDecFromStr("0.01"), LastUpdatedBlock: 5}
	auctionedSaleTokens := sdk.NewCoin(saleToken, sdk.NewInt(10000))

	auctionKeeper.setTokenPrice(ctx, sommPrice)
	auctionKeeper.setTokenPrice(ctx, saleTokenPrice)

	// Mock bank keeper fund transfer
	suite.mockSendCoinsFromModuleToModule(ctx, permissionedFunder.GetName(), auctionTypes.ModuleName, sdk.NewCoins(auctionedSaleTokens))

	// Start auction
	decreaseRate := sdk.MustNewDecFromStr("0.4")
	blockDecreaseInterval := uint64(5)
	err := auctionKeeper.BeginAuction(ctx, auctionedSaleTokens, decreaseRate, blockDecreaseInterval, permissionedFunder.GetName(), permissionedReciever.GetName())
	require.Nil(err)

	auctionID := uint32(1)

	expectedAuction := auctionTypes.Auction{
		Id:                         auctionID,
		StartingTokensForSale:      auctionedSaleTokens,
		StartBlock:                 uint64(ctx.BlockHeight()),
		EndBlock:                   0,
		InitialPriceDecreaseRate:   decreaseRate,
		CurrentPriceDecreaseRate:   decreaseRate,
		PriceDecreaseBlockInterval: blockDecreaseInterval,
		InitialUnitPriceInUsomm:    sdk.NewDec(1),
		CurrentUnitPriceInUsomm:    sdk.NewDec(1),
		RemainingTokensForSale:     auctionedSaleTokens,
		FundingModuleAccount:       permissionedFunder.GetName(),
		ProceedsModuleAccount:      permissionedReciever.GetName(),
	}

	// Run EndBlocker as base case with NO block change
	// Assure expected auction prices have not changed
	require.NotPanics(func() { auctionKeeper.EndBlocker(ctx) })

	foundAuction, found := auctionKeeper.GetActiveAuctionByID(ctx, auctionID)
	require.True(found)
	require.Equal(expectedAuction, foundAuction)

	// Run EndBlocker again with ~insufficient~ block change to induce price decrease
	require.NotPanics(func() {
		auctionKeeper.EndBlocker(ctx.WithBlockHeight(ctx.BlockHeight() + int64(blockDecreaseInterval) - 1))
	})

	foundAuction, found = auctionKeeper.GetActiveAuctionByID(ctx, auctionID)
	require.True(found)
	require.Equal(expectedAuction, foundAuction)

	// Run EndBlocker with block interval ~sufficient~ to induce decrease
	require.NotPanics(func() {
		auctionKeeper.EndBlocker(ctx.WithBlockHeight(ctx.BlockHeight() + int64(blockDecreaseInterval)))
	})

	foundAuction, found = auctionKeeper.GetActiveAuctionByID(ctx, auctionID)
	require.True(found)
	expectedAuction.CurrentUnitPriceInUsomm = sdk.MustNewDecFromStr("0.5996")
	expectedAuction.CurrentPriceDecreaseRate = sdk.MustNewDecFromStr("0.4004")
	require.Equal(expectedAuction, foundAuction)

	// Run EndBlocker with block interval ~sufficient~ to induce decrease again
	require.NotPanics(func() {
		auctionKeeper.EndBlocker(ctx.WithBlockHeight(ctx.BlockHeight() + int64(blockDecreaseInterval)))
	})

	foundAuction, found = auctionKeeper.GetActiveAuctionByID(ctx, auctionID)
	require.True(found)
	expectedAuction.CurrentUnitPriceInUsomm = sdk.MustNewDecFromStr("0.1987996")
	expectedAuction.CurrentPriceDecreaseRate = sdk.MustNewDecFromStr("0.4008004")
	require.Equal(expectedAuction, foundAuction)

	// Run EndBlocker with block interval ~sufficient~ to induce decrease to a ~negative~ price, and thus the auction must finish
	// Since we expect auction to end here without ever transfering any funds to bidders, we need to mock bank balance checks & transfers
	// back to funding module
	finalCtx := ctx.WithBlockHeight(ctx.BlockHeight() + int64(blockDecreaseInterval))
	suite.mockGetBalance(finalCtx, authtypes.NewModuleAddress(auctionTypes.ModuleName), saleToken, auctionedSaleTokens)
	suite.mockSendCoinsFromModuleToModule(finalCtx, auctionTypes.ModuleName, permissionedFunder.GetName(), sdk.NewCoins(auctionedSaleTokens))

	require.NotPanics(func() {
		auctionKeeper.EndBlocker(finalCtx)
	})

	_, found = auctionKeeper.GetActiveAuctionByID(finalCtx, auctionID)
	require.False(found)

	expectedAuction.EndBlock = uint64(finalCtx.BlockHeight())

	// Check expected auction wound up in ended auctions
	foundAuction, found = auctionKeeper.GetEndedAuctionByID(finalCtx, auctionID)
	require.True(found)

	expectedAuction.CurrentPriceDecreaseRate = sdk.MustNewDecFromStr("0.4012012004")
	require.Equal(expectedAuction, foundAuction)
}
