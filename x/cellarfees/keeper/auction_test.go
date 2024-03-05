package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/golang/mock/gomock"
	auctionTypes "github.com/peggyjv/sommelier/v7/x/auction/types"
	cellarfeestypes "github.com/peggyjv/sommelier/v7/x/cellarfees/types"
	cellarfeestypesv2 "github.com/peggyjv/sommelier/v7/x/cellarfees/types/v2"
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

	// nonzero fee balance
	suite.bankKeeper.EXPECT().GetBalance(ctx, feesAccount.GetAddress(), feeDenom).Return(fees)

	// begin auction doesn't error
	suite.auctionKeeper.EXPECT().BeginAuction(ctx, fees, gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)

	require.True(cellarfeesKeeper.beginAuction(ctx, feeDenom))
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

	// zero fee balance
	suite.bankKeeper.EXPECT().GetBalance(ctx, feesAccount.GetAddress(), feeDenom).Return(fees)

	require.False(cellarfeesKeeper.beginAuction(ctx, feeDenom))
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

	// nonzero fee balance
	suite.bankKeeper.EXPECT().GetBalance(ctx, feesAccount.GetAddress(), feeDenom).Return(fees)

	// begin auction errors
	suite.auctionKeeper.EXPECT().BeginAuction(ctx, fees, gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(auctionTypes.ErrUnauthorizedFundingModule)

	require.Panics(func() { cellarfeesKeeper.beginAuction(ctx, feeDenom) })
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

	// nonzero fee balance
	suite.bankKeeper.EXPECT().GetBalance(ctx, feesAccount.GetAddress(), feeDenom).Return(fees)

	// begin auction errors
	suite.auctionKeeper.EXPECT().BeginAuction(ctx, fees, gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(auctionTypes.ErrUnauthorizedProceedsModule)

	require.Panics(func() { cellarfeesKeeper.beginAuction(ctx, feeDenom) })
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

	// nonzero fee balance
	suite.bankKeeper.EXPECT().GetBalance(ctx, feesAccount.GetAddress(), feeDenom).Return(fees)

	// begin auction errors
	suite.auctionKeeper.EXPECT().BeginAuction(ctx, fees, gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(auctionTypes.ErrDenomCannotBeEmpty)

	require.False(cellarfeesKeeper.beginAuction(ctx, feeDenom))
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

	require.False(cellarfeesKeeper.beginAuction(ctx, feeDenom))
}
