package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	auctionTypes "github.com/peggyjv/sommelier/v4/x/auction/types"
)

// Happy path test for abci, 
func (suite *KeeperTestSuite) TestAbci() {
	ctx, auctionKeeper := suite.ctx, suite.auctionKeeper
	require := suite.Require()

	// Test base case of no auctions
	// Note BeginBlocker is only run once for completeness, since it has no code in it
	require.NotPanics(func() { auctionKeeper.BeginBlocker(ctx) })
	require.NotPanics(func() { auctionKeeper.EndBlocker(ctx) })

	// Create an auction
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
		InitialUnitPriceInUsomm:    sdk.NewDec(2),
		CurrentUnitPriceInUsomm:    sdk.NewDec(2),
		RemainingTokensForSale:     auctionedSaleTokens,
		FundingModuleAccount:       permissionedFunder.GetName(),
		ProceedsModuleAccount:      permissionedReciever.GetName(),
	}

	// Run EndBlocker as base case with insufficient block change to induce price decrease
	// Assure expected auction prices have not changed
	require.NotPanics(func() { auctionKeeper.EndBlocker(ctx) })

	foundAuction, found := auctionKeeper.GetActiveAuctionByID(ctx, auctionID)
	require.True(found)
	require.Equal(expectedAuction, foundAuction)

	

	// Run EndBlocker with block interval sufficient to induce decrease 

	// Run EndBlocker with block interval sufficient to induce decrease again

	// Run EndBlocker with block interval sufficient to induce decrease to a ~negative~ price, and thus the auction must finish
}
