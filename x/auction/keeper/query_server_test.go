package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/peggyjv/sommelier/v8/app/params"
	auctionTypes "github.com/peggyjv/sommelier/v8/x/auction/types"
)

// Happy path test for query server functions
func (suite *KeeperTestSuite) TestHappyPathsForQueryServer() {
	ctx, auctionKeeper := suite.ctx, suite.auctionKeeper
	require := suite.Require()

	auctionParams := auctionTypes.DefaultParams()
	auctionKeeper.setParams(ctx, auctionParams)

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
		Bidder:                    cosmosAddress1,
		MaxBidInUsomm:             sdk.NewCoin(params.BaseCoinUnit, sdk.NewInt(2000)),
		SaleTokenMinimumAmount:    sdk.NewCoin("weth", sdk.NewInt(20)),
		TotalFulfilledSaleTokens:  sdk.NewCoin("weth", sdk.NewInt(100)),
		SaleTokenUnitPriceInUsomm: sdk.MustNewDecFromStr("20.0"),
		TotalUsommPaid:            sdk.NewCoin(params.BaseCoinUnit, sdk.NewInt(2000)),
	}
	bid2 := &auctionTypes.Bid{
		Id:                        uint64(3),
		AuctionId:                 uint32(2),
		Bidder:                    cosmosAddress2,
		MaxBidInUsomm:             sdk.NewCoin(params.BaseCoinUnit, sdk.NewInt(1500)),
		SaleTokenMinimumAmount:    sdk.NewCoin("weth", sdk.NewInt(10)),
		TotalFulfilledSaleTokens:  sdk.NewCoin("weth", sdk.NewInt(500)),
		SaleTokenUnitPriceInUsomm: sdk.MustNewDecFromStr("10.07"),
		TotalUsommPaid:            sdk.NewCoin(params.BaseCoinUnit, sdk.NewInt(1370)),
	}
	bid3 := &auctionTypes.Bid{
		Id:                        uint64(4),
		AuctionId:                 uint32(3),
		Bidder:                    cosmosAddress2,
		MaxBidInUsomm:             sdk.NewCoin(params.BaseCoinUnit, sdk.NewInt(500)),
		SaleTokenMinimumAmount:    sdk.NewCoin("usdc", sdk.NewInt(1)),
		TotalFulfilledSaleTokens:  sdk.NewCoin("usdc", sdk.NewInt(20)),
		SaleTokenUnitPriceInUsomm: sdk.MustNewDecFromStr("20.0"),
		TotalUsommPaid:            sdk.NewCoin(params.BaseCoinUnit, sdk.NewInt(20)),
	}

	auctionKeeper.setBid(ctx, *bid1)
	auctionKeeper.setBid(ctx, *bid2)
	auctionKeeper.setBid(ctx, *bid3)

	// Create a bids for the ended auction
	bid0 := &auctionTypes.Bid{
		Id:                        uint64(1),
		AuctionId:                 uint32(1),
		Bidder:                    cosmosAddress1,
		MaxBidInUsomm:             sdk.NewCoin(params.BaseCoinUnit, sdk.NewInt(1500)),
		SaleTokenMinimumAmount:    sdk.NewCoin("matic", sdk.NewInt(100)),
		TotalFulfilledSaleTokens:  sdk.NewCoin("matic", sdk.NewInt(1500)),
		SaleTokenUnitPriceInUsomm: sdk.MustNewDecFromStr("1.0"),
		TotalUsommPaid:            sdk.NewCoin(params.BaseCoinUnit, sdk.NewInt(1500)),
	}
	auctionKeeper.setBid(ctx, *bid0)

	auctionKeeper.setLastAuctionID(ctx, uint32(3))
	auctionKeeper.setLastBidID(ctx, uint64(4))

	// Create some token prices
	tokenPrice1 := &auctionTypes.TokenPrice{
		Denom:            "Shmoo",
		Exponent:         uint64(6),
		UsdPrice:         sdk.MustNewDecFromStr("0.5"),
		LastUpdatedBlock: uint64(1),
	}
	tokenPrice2 := &auctionTypes.TokenPrice{
		Denom:            "Weth",
		Exponent:         uint64(18),
		UsdPrice:         sdk.MustNewDecFromStr("2000"),
		LastUpdatedBlock: uint64(1),
	}
	tokenPrices := []*auctionTypes.TokenPrice{tokenPrice1, tokenPrice2}

	auctionKeeper.setTokenPrice(ctx, *tokenPrice1)
	auctionKeeper.setTokenPrice(ctx, *tokenPrice2)

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
	require.Equal(&auctionTypes.QueryEndedAuctionsResponse{Auctions: []*auctionTypes.Auction{endedAuction}, Pagination: query.PageResponse{Total: 1}}, endedAuctionsResponse)

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
	require.Equal(&auctionTypes.QueryBidsByAuctionResponse{Bids: []*auctionTypes.Bid{bid1, bid2}, Pagination: query.PageResponse{Total: 2}}, activeBidsResponse)

	// QueryBidsByAuction -- ended auction
	endedBidsResponse, err := auctionKeeper.QueryBidsByAuction(sdk.WrapSDKContext(ctx), &auctionTypes.QueryBidsByAuctionRequest{AuctionId: uint32(1)})
	require.Nil(err)
	require.Equal(&auctionTypes.QueryBidsByAuctionResponse{Bids: []*auctionTypes.Bid{bid0}, Pagination: query.PageResponse{Total: 1}}, endedBidsResponse)

	// QueryTokenPrice
	tokenPriceResponse, err := auctionKeeper.QueryTokenPrice(sdk.WrapSDKContext(ctx), &auctionTypes.QueryTokenPriceRequest{Denom: "Shmoo"})
	require.Nil(err)
	require.Equal(&auctionTypes.QueryTokenPriceResponse{TokenPrice: tokenPrice1}, tokenPriceResponse)

	// QueryTokenPrices
	tokenPricesResponse, err := auctionKeeper.QueryTokenPrices(sdk.WrapSDKContext(ctx), &auctionTypes.QueryTokenPricesRequest{})
	require.Nil(err)
	require.Equal(&auctionTypes.QueryTokenPricesResponse{TokenPrices: tokenPrices}, tokenPricesResponse)
}

// Unhappy path test for query server functions
func (suite *KeeperTestSuite) TestUnhappyPathsForQueryServer() {
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

	// QueryActiveAuctions -- empty check
	activeAuctionsResponse, err := auctionKeeper.QueryActiveAuctions(sdk.WrapSDKContext(ctx), &auctionTypes.QueryActiveAuctionsRequest{})
	require.NoError(err)
	require.Zero(len(activeAuctionsResponse.Auctions))

	// QueryEndedAuctions -- empty check
	endedAuctionsResponse, err := auctionKeeper.QueryEndedAuctions(sdk.WrapSDKContext(ctx), &auctionTypes.QueryEndedAuctionsRequest{})
	require.NoError(err)
	require.Zero(len(endedAuctionsResponse.Auctions))

	// QueryBid
	bidResponse, err := auctionKeeper.QueryBid(sdk.WrapSDKContext(ctx), &auctionTypes.QueryBidRequest{BidId: uint64(4), AuctionId: uint32(3)})
	require.Equal(status.Errorf(codes.NotFound, "No bid found for specified bid id: 4, and auction id: 3"), err)
	require.Equal(&auctionTypes.QueryBidResponse{}, bidResponse)

	// QueryBidsByAuction -- empty check
	bidsResponse, err := auctionKeeper.QueryBidsByAuction(sdk.WrapSDKContext(ctx), &auctionTypes.QueryBidsByAuctionRequest{AuctionId: uint32(1)})
	require.NoError(err)
	require.Zero(len(bidsResponse.Bids))

	// QueryTokenPrice
	tokenPriceResponse, err := auctionKeeper.QueryTokenPrice(sdk.WrapSDKContext(ctx), &auctionTypes.QueryTokenPriceRequest{Denom: "Shmoo"})
	require.Equal(status.Errorf(codes.NotFound, "No token price found for denom: Shmoo"), err)
	require.Equal(&auctionTypes.QueryTokenPriceResponse{}, tokenPriceResponse)

	// QueryTokenPrices
	tokenPricesResponse, err := auctionKeeper.QueryTokenPrices(sdk.WrapSDKContext(ctx), &auctionTypes.QueryTokenPricesRequest{})
	require.NoError(err)
	require.Equal(&auctionTypes.QueryTokenPricesResponse{}, tokenPricesResponse)
}
