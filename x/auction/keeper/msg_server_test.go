package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"

	auctionTypes "github.com/peggyjv/sommelier/v4/x/auction/types"
)

func (suite *KeeperTestSuite) mockSendCoinsFromAccountToModule(ctx sdk.Context, senderAcct sdk.AccAddress, receiverModule string, amt sdk.Coins) {
	suite.bankKeeper.EXPECT().SendCoinsFromAccountToModule(ctx, senderAcct, receiverModule, amt).Return(nil)
}

func (suite *KeeperTestSuite) mockSendCoinsFromModuleToAccount(ctx sdk.Context, senderModule string, receiverAcct sdk.AccAddress, amt sdk.Coins) {
	suite.bankKeeper.EXPECT().SendCoinsFromModuleToAccount(ctx, senderModule, receiverAcct, amt).Return(nil)
}

// Happy path test for submitting a bid ~fully~
func (suite *KeeperTestSuite) TestHappyPathSubmitBidAndFulfillFully() {
	ctx, auctionKeeper := suite.ctx, suite.auctionKeeper
	require := suite.Require()

	// -----> Create an auction we can bid on first
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

	// Submit a bid
	auctionID := uint32(1)
	bidder := "cosmos16zrkzad482haunrn25ywvwy6fclh3vh7r0hcny"
	require.Nil(err)

	bid := sdk.NewCoin("usomm", sdk.NewInt(5000))
	minAmount := sdk.NewCoin(saleToken, sdk.NewInt(1))

	fulfilledBid := sdk.NewCoin(saleToken, sdk.NewInt(2500))

	// Mock out bank keeper calls
	suite.mockGetBalance(ctx, authtypes.NewModuleAddress(auctionTypes.ModuleName), saleToken, auctionedSaleTokens)
	suite.mockSendCoinsFromAccountToModule(ctx, sdk.AccAddress(bidder), auctionTypes.ModuleName, sdk.NewCoins(bid))
	suite.mockSendCoinsFromModuleToAccount(ctx, auctionTypes.ModuleName, sdk.AccAddress(bidder), sdk.NewCoins(fulfilledBid))

	// ~Actually~ submit the bid now
	response, err := auctionKeeper.SubmitBid(sdk.WrapSDKContext(ctx), &auctionTypes.MsgSubmitBidRequest{
		AuctionId:              auctionID,
		Bidder:                 bidder,
		MaxBidInUsomm:          bid,
		SaleTokenMinimumAmount: minAmount,
		Signer:                 sdk.AccAddress(bidder).String(),
	})
	require.Nil(err)

	// Assert bid store and response bids are both equal to expected bid
	expectedBid := auctionTypes.Bid{
		Id:                        uint64(1),
		AuctionId:                 uint32(1),
		Bidder:                    bidder,
		MaxBidInUsomm:             bid,
		SaleTokenMinimumAmount:    minAmount,
		TotalFulfilledSaleTokens:  fulfilledBid,
		SaleTokenUnitPriceInUsomm: sdk.NewDec(2),
		TotalUsommPaid:            bid,
	}
	require.Equal(&expectedBid, response.Bid)

	storedBid, found := auctionKeeper.GetBid(ctx, uint32(1), uint64(1))
	require.True(found)
	require.Equal(expectedBid, storedBid)

	// Assert auction token amounts are updated
	expectedUpdatedAuction := auctionTypes.Auction{
		Id:                         auctionID,
		StartingTokensForSale:      auctionedSaleTokens,
		StartBlock:                 uint64(ctx.BlockHeight()),
		EndBlock:                   0,
		InitialPriceDecreaseRate:   decreaseRate,
		CurrentPriceDecreaseRate:   decreaseRate,
		PriceDecreaseBlockInterval: blockDecreaseInterval,
		InitialUnitPriceInUsomm:    sdk.NewDec(2),
		CurrentUnitPriceInUsomm:    sdk.NewDec(2),
		RemainingTokensForSale:     sdk.NewCoin(saleToken, sdk.NewInt(7500)), // this is the important part, need to make sure it decremented
		FundingModuleAccount:       permissionedFunder.GetName(),
		ProceedsModuleAccount:      permissionedReciever.GetName(),
	}

	activeAuction, found := auctionKeeper.GetActiveAuctionByID(ctx, auctionID)
	require.True(found)
	require.Equal(expectedUpdatedAuction, activeAuction)
}
