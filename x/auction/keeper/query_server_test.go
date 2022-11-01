package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	auctionTypes "github.com/peggyjv/sommelier/v4/x/auction/types"
)

// Happy path test for query server functions
func (suite *KeeperTestSuite) TestHappyPathsForQueryServer() {
	ctx, auctionKeeper := suite.ctx, suite.auctionKeeper
	require := suite.Require()

	params := auctionTypes.DefaultParams()
	auctionKeeper.setParams(ctx, params)

	// Create some active auctions
	activeAuction1 := &auctionTypes.Auction{
		Id:                         uint32(2),
		StartingTokensForSale:      sdk.NewCoin("weth", sdk.NewInt(17000)),
		StartBlock:                 uint64(1),
		EndBlock:                   uint64(0),
		InitialPriceDecreaseRate:   sdk.MustNewDecFromStr("0.05"),
		CurrentPriceDecreaseRate:   sdk.MustNewDecFromStr("0.05"),
		PriceDecreaseBlockInterval: uint64(20),
		InitialUnitPriceInUsomm:    sdk.MustNewDecFromStr("20"),
		CurrentUnitPriceInUsomm:    sdk.MustNewDecFromStr("10"),
		RemainingTokensForSale:     sdk.NewCoin("weth", sdk.NewInt(10000)),
		FundingModuleAccount:       permissionedFunder.GetName(),
		ProceedsModuleAccount:      permissionedReciever.GetName(),
	}
	activeAuction2 := &auctionTypes.Auction{
		Id:                         uint32(3),
		StartingTokensForSale:      sdk.NewCoin("usdc", sdk.NewInt(9000)),
		StartBlock:                 uint64(3),
		EndBlock:                   uint64(0),
		InitialPriceDecreaseRate:   sdk.MustNewDecFromStr("0.02"),
		CurrentPriceDecreaseRate:   sdk.MustNewDecFromStr("0.02"),
		PriceDecreaseBlockInterval: uint64(10),
		InitialUnitPriceInUsomm:    sdk.MustNewDecFromStr("1"),
		CurrentUnitPriceInUsomm:    sdk.MustNewDecFromStr("1"),
		RemainingTokensForSale:     sdk.NewCoin("usdc", sdk.NewInt(5000)),
		FundingModuleAccount:       permissionedFunder.GetName(),
		ProceedsModuleAccount:      permissionedReciever.GetName(),
	}
	auctionKeeper.setActiveAuction(ctx, *activeAuction1)
	auctionKeeper.setActiveAuction(ctx, *activeAuction2)

	// Create an ended auction
	endedAuction := &auctionTypes.Auction{
		Id:                         uint32(1),
		StartingTokensForSale:      sdk.NewCoin("matic", sdk.NewInt(3000)),
		StartBlock:                 uint64(1),
		EndBlock:                   uint64(5),
		InitialPriceDecreaseRate:   sdk.MustNewDecFromStr("0.07"),
		CurrentPriceDecreaseRate:   sdk.MustNewDecFromStr("0.02"),
		PriceDecreaseBlockInterval: uint64(3),
		InitialUnitPriceInUsomm:    sdk.MustNewDecFromStr("17"),
		CurrentUnitPriceInUsomm:    sdk.MustNewDecFromStr("12.5"),
		RemainingTokensForSale:     sdk.NewCoin("weth", sdk.NewInt(1000)),
		FundingModuleAccount:       permissionedFunder.GetName(),
		ProceedsModuleAccount:      permissionedReciever.GetName(),
	}
	auctionKeeper.setEndedAuction(ctx, *endedAuction)

	// Create some bids for active auctions
	bid1 := &auctionTypes.Bid{
		Id:                        uint64(2),
		AuctionId:                 uint32(2),
		Bidder:                    cosmos_address_1,
		MaxBidInUsomm:             sdk.NewCoin(auctionTypes.UsommDenom, sdk.NewInt(2000)),
		SaleTokenMinimumAmount:    sdk.NewCoin("weth", sdk.NewInt(20)),
		TotalFulfilledSaleTokens:  sdk.NewCoin("weth", sdk.NewInt(100)),
		SaleTokenUnitPriceInUsomm: sdk.MustNewDecFromStr("20.0"),
		TotalUsommPaid:            sdk.NewCoin(auctionTypes.UsommDenom, sdk.NewInt(2000)),
	}
	bid2 := &auctionTypes.Bid{
		Id:                        uint64(3),
		AuctionId:                 uint32(2),
		Bidder:                    cosmos_address_2,
		MaxBidInUsomm:             sdk.NewCoin(auctionTypes.UsommDenom, sdk.NewInt(1500)),
		SaleTokenMinimumAmount:    sdk.NewCoin("weth", sdk.NewInt(10)),
		TotalFulfilledSaleTokens:  sdk.NewCoin("weth", sdk.NewInt(500)),
		SaleTokenUnitPriceInUsomm: sdk.MustNewDecFromStr("10.07"),
		TotalUsommPaid:            sdk.NewCoin(auctionTypes.UsommDenom, sdk.NewInt(1370)),
	}
	bid3 := &auctionTypes.Bid{
		Id:                        uint64(4),
		AuctionId:                 uint32(3),
		Bidder:                    cosmos_address_2,
		MaxBidInUsomm:             sdk.NewCoin(auctionTypes.UsommDenom, sdk.NewInt(500)),
		SaleTokenMinimumAmount:    sdk.NewCoin("usdc", sdk.NewInt(1)),
		TotalFulfilledSaleTokens:  sdk.NewCoin("usdc", sdk.NewInt(20)),
		SaleTokenUnitPriceInUsomm: sdk.MustNewDecFromStr("20.0"),
		TotalUsommPaid:            sdk.NewCoin(auctionTypes.UsommDenom, sdk.NewInt(20)),
	}

	auctionKeeper.setBid(ctx, *bid1)
	auctionKeeper.setBid(ctx, *bid2)
	auctionKeeper.setBid(ctx, *bid3)

	// Create a bids for the ended auction
	bid0 := &auctionTypes.Bid{
		Id:                        uint64(1),
		AuctionId:                 uint32(1),
		Bidder:                    cosmos_address_1,
		MaxBidInUsomm:             sdk.NewCoin(auctionTypes.UsommDenom, sdk.NewInt(1500)),
		SaleTokenMinimumAmount:    sdk.NewCoin("matic", sdk.NewInt(100)),
		TotalFulfilledSaleTokens:  sdk.NewCoin("matic", sdk.NewInt(1500)),
		SaleTokenUnitPriceInUsomm: sdk.MustNewDecFromStr("1.0"),
		TotalUsommPaid:            sdk.NewCoin(auctionTypes.UsommDenom, sdk.NewInt(1500)),
	}
	auctionKeeper.setBid(ctx, *bid0)

	auctionKeeper.setLastAuctionID(ctx, uint32(3))
	auctionKeeper.setLastBidID(ctx, uint64(4))

	// -- Actually begin testing

	// QueryParams
	paramsResponse, err := auctionKeeper.QueryParams(sdk.WrapSDKContext(ctx), &auctionTypes.QueryParamsRequest{})
	require.Nil(err)
	require.Equal(&auctionTypes.QueryParamsResponse{Params: auctionTypes.DefaultParams()}, paramsResponse)

	// QueryActiveAuction
	activeAuctionResponse, err := auctionKeeper.QueryActiveAuction(sdk.WrapSDKContext(ctx), &auctionTypes.QueryActiveAuctionRequest{AuctionId: uint32(3)})
	require.Nil(err)
	require.Equal(&auctionTypes.QueryActiveAuctionResponse{Auction: activeAuction2}, activeAuctionResponse)

	// QueryEndedAuction
	endedAuctionResponse, err := auctionKeeper.QueryEndedAuction(sdk.WrapSDKContext(ctx), &auctionTypes.QueryEndedAuctionRequest{AuctionId: uint32(1)})
	require.Nil(err)
	require.Equal(&auctionTypes.QueryEndedAuctionResponse{Auction: endedAuction}, endedAuctionResponse)

	// QueryActiveAuctions
	activeAuctionsResponse, err := auctionKeeper.QueryActiveAuctions(sdk.WrapSDKContext(ctx), &auctionTypes.QueryActiveAuctionsRequest{})
	require.Nil(err)
	require.Equal(&auctionTypes.QueryActiveAuctionsResponse{Auctions: []*auctionTypes.Auction{activeAuction1, activeAuction2}}, activeAuctionsResponse)

	// QueryEndedAuctions
	endedAuctionsResponse, err := auctionKeeper.QueryEndedAuctions(sdk.WrapSDKContext(ctx), &auctionTypes.QueryEndedAuctionsRequest{})
	require.Nil(err)
	require.Equal(&auctionTypes.QueryEndedAuctionsResponse{Auctions: []*auctionTypes.Auction{endedAuction}}, endedAuctionsResponse)

	// QueryBid -- active auction
	activeBidResponse, err := auctionKeeper.QueryBid(sdk.WrapSDKContext(ctx), &auctionTypes.QueryBidRequest{BidId: uint64(4), AuctionId: uint32(3)})
	require.Nil(err)
	require.Equal(&auctionTypes.QueryBidResponse{Bid: bid3}, activeBidResponse)

	// QueryBid -- ended auction
	endedBidResponse, err := auctionKeeper.QueryBid(sdk.WrapSDKContext(ctx), &auctionTypes.QueryBidRequest{BidId: uint64(1), AuctionId: uint32(1)})
	require.Nil(err)
	require.Equal(&auctionTypes.QueryBidResponse{Bid: bid0}, endedBidResponse)

	// QueryBidsByAuction -- active auction
	activeBidsResponse, err := auctionKeeper.QueryBidsByAuction(sdk.WrapSDKContext(ctx), &auctionTypes.QueryBidsByAuctionRequest{AuctionId: uint32(2)})
	require.Nil(err)
	require.Equal(&auctionTypes.QueryBidsByAuctionResponse{Bids: []*auctionTypes.Bid{bid1, bid2}}, activeBidsResponse)

	// QueryBidsByAuction -- ended auction
	endedBidsResponse, err := auctionKeeper.QueryBidsByAuction(sdk.WrapSDKContext(ctx), &auctionTypes.QueryBidsByAuctionRequest{AuctionId: uint32(1)})
	require.Nil(err)
	require.Equal(&auctionTypes.QueryBidsByAuctionResponse{Bids: []*auctionTypes.Bid{bid0}}, endedBidsResponse)
}

// Unhappy path test for query server functions
func (suite *KeeperTestSuite) TestUnhappPathsForQueryServer() {
	ctx, auctionKeeper := suite.ctx, suite.auctionKeeper
	require := suite.Require()

	// Note we dont unhappy path test QueryParams bc it cannot return an error currently

	// QueryActiveAuction
	activeAuctionResponse, err := auctionKeeper.QueryActiveAuction(sdk.WrapSDKContext(ctx), &auctionTypes.QueryActiveAuctionRequest{AuctionId: uint32(3)})
	require.Equal(status.Errorf(codes.NotFound, "No active auction found for id: 3"), err)
	require.Equal(&auctionTypes.QueryActiveAuctionResponse{}, activeAuctionResponse)

	// QueryEndedAuction
	endedAuctionResponse, err := auctionKeeper.QueryEndedAuction(sdk.WrapSDKContext(ctx), &auctionTypes.QueryEndedAuctionRequest{AuctionId: uint32(1)})
	require.Equal(status.Errorf(codes.NotFound, "No ended auction found for id: 1"), err)
	require.Equal(&auctionTypes.QueryEndedAuctionResponse{}, endedAuctionResponse)

	// QueryActiveAuctions
	activeAuctionsResponse, err := auctionKeeper.QueryActiveAuctions(sdk.WrapSDKContext(ctx), &auctionTypes.QueryActiveAuctionsRequest{})
	require.Equal(status.Error(codes.NotFound, "No active auctions found"), err)
	require.Equal(&auctionTypes.QueryActiveAuctionsResponse{}, activeAuctionsResponse)

	// QueryEndedAuctions
	endedAuctionsResponse, err := auctionKeeper.QueryEndedAuctions(sdk.WrapSDKContext(ctx), &auctionTypes.QueryEndedAuctionsRequest{})
	require.Equal(status.Error(codes.NotFound, "No ended auctions found"), err)
	require.Equal(&auctionTypes.QueryEndedAuctionsResponse{}, endedAuctionsResponse)

	// QueryBid
	bidResponse, err := auctionKeeper.QueryBid(sdk.WrapSDKContext(ctx), &auctionTypes.QueryBidRequest{BidId: uint64(4), AuctionId: uint32(3)})
	require.Equal(status.Errorf(codes.NotFound, "No bid found for specified bid id: 4, and auction id: 3"), err)
	require.Equal(&auctionTypes.QueryBidResponse{}, bidResponse)

	// QueryBidsByAuction
	bidsResponse, err := auctionKeeper.QueryBidsByAuction(sdk.WrapSDKContext(ctx), &auctionTypes.QueryBidsByAuctionRequest{AuctionId: uint32(1)})
	require.Equal(status.Errorf(codes.NotFound, "No bids found for auction id: 1"), err)
	require.Equal(&auctionTypes.QueryBidsByAuctionResponse{}, bidsResponse)

}
