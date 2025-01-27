package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/golang/mock/gomock"
	auctionTypes "github.com/peggyjv/sommelier/v9/x/auction/types"
	cellarfeestypes "github.com/peggyjv/sommelier/v9/x/cellarfees/types"
	cellarfeestypesv2 "github.com/peggyjv/sommelier/v9/x/cellarfees/types/v2"
)

const feeDenom = "testdenom"

func (suite *KeeperTestSuite) TestHappyPathBeginAuction() {
	ctx, cellarfeesKeeper := suite.ctx, suite.cellarfeesKeeper
	require := suite.Require()

	params := cellarfeestypesv2.DefaultParams()
	cellarfeesKeeper.SetParams(ctx, params)

	fees := sdk.NewCoin(feeDenom, sdk.NewInt(1000000))

	// retreiving module account
	suite.accountKeeper.EXPECT().GetModuleAccount(ctx, cellarfeestypes.ModuleName).Return(feesAccount)

	// no active auctions
	suite.auctionKeeper.EXPECT().GetActiveAuctions(ctx).Return([]*auctionTypes.Auction{})

	// begin auction doesn't error
	suite.auctionKeeper.EXPECT().BeginAuction(ctx, fees, gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)

	require.True(cellarfeesKeeper.beginAuction(ctx, fees))
}

func (suite *KeeperTestSuite) TestAuctionFeeBalanceZeroDoesNotStartAuction() {
	ctx, cellarfeesKeeper := suite.ctx, suite.cellarfeesKeeper
	require := suite.Require()

	params := cellarfeestypesv2.DefaultParams()
	cellarfeesKeeper.SetParams(ctx, params)

	fees := sdk.NewCoin(feeDenom, sdk.NewInt(0))

	// retreiving module account
	suite.accountKeeper.EXPECT().GetModuleAccount(ctx, cellarfeestypes.ModuleName).Return(feesAccount)

	// no active auctions
	suite.auctionKeeper.EXPECT().GetActiveAuctions(ctx).Return([]*auctionTypes.Auction{})

	require.False(cellarfeesKeeper.beginAuction(ctx, fees))
}

func (suite *KeeperTestSuite) TestAuctionUnauthorizedFundingModule() {
	ctx, cellarfeesKeeper := suite.ctx, suite.cellarfeesKeeper
	require := suite.Require()

	params := cellarfeestypesv2.DefaultParams()
	cellarfeesKeeper.SetParams(ctx, params)

	fees := sdk.NewCoin(feeDenom, sdk.NewInt(1000000))

	// retreiving module account
	suite.accountKeeper.EXPECT().GetModuleAccount(ctx, cellarfeestypes.ModuleName).Return(feesAccount)

	// no active auctions
	suite.auctionKeeper.EXPECT().GetActiveAuctions(ctx).Return([]*auctionTypes.Auction{})

	// begin auction errors
	suite.auctionKeeper.EXPECT().BeginAuction(ctx, fees, gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(auctionTypes.ErrUnauthorizedFundingModule)

	require.Panics(func() { cellarfeesKeeper.beginAuction(ctx, fees) })
}

func (suite *KeeperTestSuite) TestAuctionUnauthorizedProceedsModule() {
	ctx, cellarfeesKeeper := suite.ctx, suite.cellarfeesKeeper
	require := suite.Require()

	params := cellarfeestypesv2.DefaultParams()
	cellarfeesKeeper.SetParams(ctx, params)

	fees := sdk.NewCoin(feeDenom, sdk.NewInt(1000000))

	// retreiving module account
	suite.accountKeeper.EXPECT().GetModuleAccount(ctx, cellarfeestypes.ModuleName).Return(feesAccount)

	// no active auctions
	suite.auctionKeeper.EXPECT().GetActiveAuctions(ctx).Return([]*auctionTypes.Auction{})

	// begin auction errors
	suite.auctionKeeper.EXPECT().BeginAuction(ctx, fees, gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(auctionTypes.ErrUnauthorizedProceedsModule)

	require.Panics(func() { cellarfeesKeeper.beginAuction(ctx, fees) })
}

func (suite *KeeperTestSuite) TestAuctionNonPanicError() {
	ctx, cellarfeesKeeper := suite.ctx, suite.cellarfeesKeeper
	require := suite.Require()

	params := cellarfeestypesv2.DefaultParams()
	cellarfeesKeeper.SetParams(ctx, params)

	fees := sdk.NewCoin(feeDenom, sdk.NewInt(1000000))

	// retreiving module account
	suite.accountKeeper.EXPECT().GetModuleAccount(ctx, cellarfeestypes.ModuleName).Return(feesAccount)

	// no active auctions
	suite.auctionKeeper.EXPECT().GetActiveAuctions(ctx).Return([]*auctionTypes.Auction{})

	// begin auction errors
	suite.auctionKeeper.EXPECT().BeginAuction(ctx, fees, gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(auctionTypes.ErrDenomCannotBeEmpty)

	require.False(cellarfeesKeeper.beginAuction(ctx, fees))
}

func (suite *KeeperTestSuite) TestAuctionAlreadyActive() {
	ctx, cellarfeesKeeper := suite.ctx, suite.cellarfeesKeeper
	require := suite.Require()

	params := cellarfeestypesv2.DefaultParams()
	cellarfeesKeeper.SetParams(ctx, params)

	fees := sdk.NewCoin(feeDenom, sdk.NewInt(1000000))

	// active auction for our denom
	testAuction := auctionTypes.Auction{
		StartingTokensForSale: fees,
	}
	suite.auctionKeeper.EXPECT().GetActiveAuctions(ctx).Return([]*auctionTypes.Auction{&testAuction})

	require.False(cellarfeesKeeper.beginAuction(ctx, fees))
}
