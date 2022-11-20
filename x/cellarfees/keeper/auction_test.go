package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/golang/mock/gomock"
	auctionTypes "github.com/peggyjv/sommelier/v4/x/auction/types"
	cellarfeesTypes "github.com/peggyjv/sommelier/v4/x/cellarfees/types"
)

func (suite *KeeperTestSuite) TestHappyPathBeginAuction() {
	ctx, cellarfeesKeeper := suite.ctx, suite.cellarfeesKeeper
	require := suite.Require()

	params := cellarfeesTypes.DefaultParams()
	cellarfeesKeeper.SetParams(ctx, params)

	feeDenom := "testdenom"
	fees := sdk.NewCoin(feeDenom, sdk.NewInt(1000000))

	// retreiving module account
	suite.accountKeeper.EXPECT().GetModuleAccount(ctx, cellarfeesTypes.ModuleName).Return(feesAccount)

	// no active auctions
	suite.auctionKeeper.EXPECT().GetActiveAuctions(ctx).Return([]*auctionTypes.Auction{})

	// nonzero fee balance
	suite.bankKeeper.EXPECT().GetBalance(ctx, feesAccount.GetAddress(), feeDenom).Return(fees)

	// begin auction doesn't error
	suite.auctionKeeper.EXPECT().BeginAuction(ctx, fees, gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)

	require.True(cellarfeesKeeper.beginAuction(ctx, feeDenom))
}

func (suite *KeeperTestSuite) TestAuctionFeeBalanceZeroPanics() {
	ctx, cellarfeesKeeper := suite.ctx, suite.cellarfeesKeeper
	require := suite.Require()

	params := cellarfeesTypes.DefaultParams()
	cellarfeesKeeper.SetParams(ctx, params)

	feeDenom := "testdenom"
	fees := sdk.NewCoin(feeDenom, sdk.NewInt(0))

	// retreiving module account
	suite.accountKeeper.EXPECT().GetModuleAccount(ctx, cellarfeesTypes.ModuleName).Return(feesAccount)

	// no active auctions
	suite.auctionKeeper.EXPECT().GetActiveAuctions(ctx).Return([]*auctionTypes.Auction{})

	// zero fee balance
	suite.bankKeeper.EXPECT().GetBalance(ctx, feesAccount.GetAddress(), feeDenom).Return(fees)

	require.Panics(func() { cellarfeesKeeper.beginAuction(ctx, feeDenom) })
}

func (suite *KeeperTestSuite) TestAuctionUnauthorizedFundingModule() {
	ctx, cellarfeesKeeper := suite.ctx, suite.cellarfeesKeeper
	require := suite.Require()

	params := cellarfeesTypes.DefaultParams()
	cellarfeesKeeper.SetParams(ctx, params)

	feeDenom := "testdenom"
	fees := sdk.NewCoin(feeDenom, sdk.NewInt(1000000))

	// retreiving module account
	suite.accountKeeper.EXPECT().GetModuleAccount(ctx, cellarfeesTypes.ModuleName).Return(feesAccount)

	// no active auctions
	suite.auctionKeeper.EXPECT().GetActiveAuctions(ctx).Return([]*auctionTypes.Auction{})

	// nonzero fee balance
	suite.bankKeeper.EXPECT().GetBalance(ctx, feesAccount.GetAddress(), feeDenom).Return(fees)

	// begin auction errors
	suite.auctionKeeper.EXPECT().BeginAuction(ctx, fees, gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(auctionTypes.ErrUnauthorizedFundingModule)

	require.Panics(func() { cellarfeesKeeper.beginAuction(ctx, feeDenom) })
}

func (suite *KeeperTestSuite) TestAuctionUnauthorizedProceedsModule() {
	ctx, cellarfeesKeeper := suite.ctx, suite.cellarfeesKeeper
	require := suite.Require()

	params := cellarfeesTypes.DefaultParams()
	cellarfeesKeeper.SetParams(ctx, params)

	feeDenom := "testdenom"
	fees := sdk.NewCoin(feeDenom, sdk.NewInt(1000000))

	// retreiving module account
	suite.accountKeeper.EXPECT().GetModuleAccount(ctx, cellarfeesTypes.ModuleName).Return(feesAccount)

	// no active auctions
	suite.auctionKeeper.EXPECT().GetActiveAuctions(ctx).Return([]*auctionTypes.Auction{})

	// nonzero fee balance
	suite.bankKeeper.EXPECT().GetBalance(ctx, feesAccount.GetAddress(), feeDenom).Return(fees)

	// begin auction errors
	suite.auctionKeeper.EXPECT().BeginAuction(ctx, fees, gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(auctionTypes.ErrUnauthorizedProceedsModule)

	require.Panics(func() { cellarfeesKeeper.beginAuction(ctx, feeDenom) })
}

func (suite *KeeperTestSuite) TestAuctionNonPanicError() {
	ctx, cellarfeesKeeper := suite.ctx, suite.cellarfeesKeeper
	require := suite.Require()

	params := cellarfeesTypes.DefaultParams()
	cellarfeesKeeper.SetParams(ctx, params)

	feeDenom := "testdenom"
	fees := sdk.NewCoin(feeDenom, sdk.NewInt(1000000))

	// retreiving module account
	suite.accountKeeper.EXPECT().GetModuleAccount(ctx, cellarfeesTypes.ModuleName).Return(feesAccount)

	// no active auctions
	suite.auctionKeeper.EXPECT().GetActiveAuctions(ctx).Return([]*auctionTypes.Auction{})

	// nonzero fee balance
	suite.bankKeeper.EXPECT().GetBalance(ctx, feesAccount.GetAddress(), feeDenom).Return(fees)

	// begin auction errors
	suite.auctionKeeper.EXPECT().BeginAuction(ctx, fees, gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(auctionTypes.ErrDenomCannotBeEmpty)

	require.False(cellarfeesKeeper.beginAuction(ctx, feeDenom))
}

func (suite *KeeperTestSuite) TestAuctionAlreadyActive() {
	ctx, cellarfeesKeeper := suite.ctx, suite.cellarfeesKeeper
	require := suite.Require()

	params := cellarfeesTypes.DefaultParams()
	cellarfeesKeeper.SetParams(ctx, params)

	feeDenom := "testdenom"
	fees := sdk.NewCoin(feeDenom, sdk.NewInt(1000000))

	// active auction for our denom
	testAuction := auctionTypes.Auction{
		StartingTokensForSale: fees,
	}
	suite.auctionKeeper.EXPECT().GetActiveAuctions(ctx).Return([]*auctionTypes.Auction{&testAuction})

	require.False(cellarfeesKeeper.beginAuction(ctx, feeDenom))
}
