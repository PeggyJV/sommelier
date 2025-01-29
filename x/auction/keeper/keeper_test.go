package keeper

import (
	"fmt"
	"testing"

	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/testutil"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/suite"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/peggyjv/sommelier/v9/app/params"
	moduletestutil "github.com/peggyjv/sommelier/v9/testutil"
	auctionTypes "github.com/peggyjv/sommelier/v9/x/auction/types"

	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"
	tmtime "github.com/cometbft/cometbft/types/time"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	paramskeeper "github.com/cosmos/cosmos-sdk/x/params/keeper"
	auctiontestutil "github.com/peggyjv/sommelier/v9/x/auction/testutil"
)

var (
	permissionedFunder          = authtypes.NewEmptyModuleAccount("permissionedFunder")
	permissionedReciever        = authtypes.NewEmptyModuleAccount("permissionedReciever")
	cosmosAddress1       string = "somm16zrkzad482haunrn25ywvwy6fclh3vh70nc5zw"
	cosmosAddress2       string = "somm18ld4633yswcyjdklej3att6aw93nhlf7596qkk"
)

type BeginAuctionRequest struct {
	ctx                        sdk.Context
	startingTokensForSale      sdk.Coin
	initialPriceDecreaseRate   sdk.Dec
	priceDecreaseBlockInterval uint64
	fundingModuleAccount       string
	proceedsModuleAccount      string
}

type KeeperTestSuite struct {
	suite.Suite

	ctx           sdk.Context
	auctionKeeper Keeper
	bankKeeper    *auctiontestutil.MockBankKeeper
	accountKeeper *auctiontestutil.MockAccountKeeper

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
	defer ctrl.Finish()

	suite.bankKeeper = auctiontestutil.NewMockBankKeeper(ctrl)
	suite.accountKeeper = auctiontestutil.NewMockAccountKeeper(ctrl)
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
		suite.accountKeeper,
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

// Happy path for BeginAuction call
func (suite *KeeperTestSuite) TestHappyPathBeginAuction() {
	ctx, auctionKeeper := suite.ctx, suite.auctionKeeper
	require := suite.Require()

	auctionParams := auctionTypes.DefaultParams()
	auctionKeeper.setParams(ctx, auctionParams)

	sommPrice := auctionTypes.TokenPrice{Denom: params.BaseCoinUnit, Exponent: 6, UsdPrice: sdk.MustNewDecFromStr("0.01"), LastUpdatedBlock: 5}

	/* #nosec */
	saleToken := "gravity0xdac17f958d2ee523a2206206994597c13d831ec7"
	saleTokenPrice := auctionTypes.TokenPrice{Denom: saleToken, Exponent: 6, UsdPrice: sdk.MustNewDecFromStr("0.02"), LastUpdatedBlock: 5}
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
	auctionParams := auctionTypes.DefaultParams()
	auctionKeeper.setParams(ctx, auctionParams)

	sommPrice := auctionTypes.TokenPrice{Denom: params.BaseCoinUnit, Exponent: 6, UsdPrice: sdk.MustNewDecFromStr("0.02"), LastUpdatedBlock: 2}

	/* #nosec */
	saleToken := "gravity0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48"
	saleTokenPrice := auctionTypes.TokenPrice{Denom: saleToken, Exponent: 6, UsdPrice: sdk.MustNewDecFromStr("0.01"), LastUpdatedBlock: 2}
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
	amountPaid1 := sdk.NewCoin(params.BaseCoinUnit, sdk.NewInt(2500))
	auctionKeeper.setBid(ctx, auctionTypes.Bid{
		Id:                       1,
		AuctionId:                1,
		Bidder:                   "bidder1",
		MaxBidInUsomm:            sdk.NewCoin(params.BaseCoinUnit, sdk.NewInt(2500)),
		SaleTokenMinimumAmount:   sdk.NewCoin(saleToken, sdk.NewInt(0)),
		TotalFulfilledSaleTokens: sdk.NewCoin(saleToken, sdk.NewInt(5000)),
		TotalUsommPaid:           amountPaid1,
	})

	amountPaid2 := sdk.NewCoin(params.BaseCoinUnit, sdk.NewInt(1250))
	auctionKeeper.setBid(ctx, auctionTypes.Bid{
		Id:                       2,
		AuctionId:                1,
		Bidder:                   "bidder2",
		MaxBidInUsomm:            sdk.NewCoin(params.BaseCoinUnit, sdk.NewInt(1250)),
		SaleTokenMinimumAmount:   sdk.NewCoin(saleToken, sdk.NewInt(0)),
		TotalFulfilledSaleTokens: sdk.NewCoin(saleToken, sdk.NewInt(2500)),
		TotalUsommPaid:           amountPaid2,
	})

	// Burn and send remaining proceeds from bids
	totalBurnExpected := sdk.NewCoin(params.BaseCoinUnit, amountPaid1.Amount.Add(amountPaid2.Amount).Quo(sdk.NewInt(2)))
	suite.mockBurnCoins(ctx, auctionTypes.ModuleName, sdk.NewCoins(totalBurnExpected))
	totalUsommExpected := sdk.NewCoin(params.BaseCoinUnit, amountPaid1.Amount.Add(amountPaid2.Amount).Sub(totalBurnExpected.Amount))
	suite.mockSendCoinsFromModuleToModule(ctx, auctionTypes.ModuleName, permissionedReciever.GetName(), sdk.NewCoins(totalUsommExpected))

	// Change active auction tokens remaining before finishing auction to pretend tokens were sold
	createdAuction.RemainingTokensForSale = remainingSaleTokens

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

// Unhappy path tests for BeginAuction
func (suite *KeeperTestSuite) TestUnhappyPathsForBeginAuction() {
	ctx, auctionKeeper := suite.ctx, suite.auctionKeeper
	require := suite.Require()

	// Define basic param(s)
	auctionParams := auctionTypes.DefaultParams()
	auctionKeeper.setParams(ctx, auctionParams)

	// Setup some token prices
	sommPrice := auctionTypes.TokenPrice{Denom: params.BaseCoinUnit, Exponent: 6, UsdPrice: sdk.MustNewDecFromStr("0.01"), LastUpdatedBlock: 2}

	/* #nosec */
	saleToken := "gravity0xaaaebe6fe48e54f431b0c390cfaf0b017d09d42d"
	saleTokenPrice := auctionTypes.TokenPrice{Denom: saleToken, Exponent: 6, UsdPrice: sdk.MustNewDecFromStr("0.01"), LastUpdatedBlock: 5}
	auctionedSaleTokens := sdk.NewCoin(saleToken, sdk.NewInt(10000))

	tests := []struct {
		name                string
		beginAuctionRequest BeginAuctionRequest
		expectedError       error
		runsBefore          runsBeforeWrapper
	}{
		{
			name: "Unpermissioned funder module account",
			beginAuctionRequest: BeginAuctionRequest{
				ctx:                        ctx,
				startingTokensForSale:      auctionedSaleTokens,
				initialPriceDecreaseRate:   sdk.MustNewDecFromStr("0.05"),
				priceDecreaseBlockInterval: uint64(10),
				fundingModuleAccount:       "cork",
				proceedsModuleAccount:      permissionedReciever.GetName(),
			},
			expectedError: errorsmod.Wrapf(auctionTypes.ErrUnauthorizedFundingModule, "Module Account: cork"),
			runsBefore:    func() {},
		},
		{
			name: "Unpermissioned proceeds module account",
			beginAuctionRequest: BeginAuctionRequest{
				ctx:                        ctx,
				startingTokensForSale:      auctionedSaleTokens,
				initialPriceDecreaseRate:   sdk.MustNewDecFromStr("0.05"),
				priceDecreaseBlockInterval: uint64(10),
				fundingModuleAccount:       permissionedFunder.GetName(),
				proceedsModuleAccount:      "gravity",
			},
			expectedError: errorsmod.Wrapf(auctionTypes.ErrUnauthorizedProceedsModule, "Module Account: gravity"),
			runsBefore:    func() {},
		},
		{
			name: "Starting denom price not found",
			beginAuctionRequest: BeginAuctionRequest{
				ctx:                        ctx,
				startingTokensForSale:      sdk.NewCoin("anvil", sdk.NewInt(7)),
				initialPriceDecreaseRate:   sdk.MustNewDecFromStr("0.05"),
				priceDecreaseBlockInterval: uint64(10),
				fundingModuleAccount:       permissionedFunder.GetName(),
				proceedsModuleAccount:      permissionedReciever.GetName(),
			},
			expectedError: errorsmod.Wrapf(auctionTypes.ErrCouldNotFindSaleTokenPrice, "starting amount denom: anvil"),
			runsBefore:    func() {},
		},
		{
			name: "Starting denom price update too old",
			beginAuctionRequest: BeginAuctionRequest{
				ctx:                        ctx.WithBlockHeight(int64(saleTokenPrice.LastUpdatedBlock) + int64(auctionParams.PriceMaxBlockAge) + 1),
				startingTokensForSale:      auctionedSaleTokens,
				initialPriceDecreaseRate:   sdk.MustNewDecFromStr("0.05"),
				priceDecreaseBlockInterval: uint64(10),
				fundingModuleAccount:       permissionedFunder.GetName(),
				proceedsModuleAccount:      permissionedReciever.GetName(),
			},
			expectedError: errorsmod.Wrapf(auctionTypes.ErrLastSaleTokenPriceTooOld, "starting amount denom: %s", saleToken),
			runsBefore: func() {
				auctionKeeper.setTokenPrice(ctx, saleTokenPrice)
			},
		},
		{
			name: "Usomm price not found",
			beginAuctionRequest: BeginAuctionRequest{
				ctx:                        ctx,
				startingTokensForSale:      auctionedSaleTokens,
				initialPriceDecreaseRate:   sdk.MustNewDecFromStr("0.05"),
				priceDecreaseBlockInterval: uint64(10),
				fundingModuleAccount:       permissionedFunder.GetName(),
				proceedsModuleAccount:      permissionedReciever.GetName(),
			},
			expectedError: errorsmod.Wrap(auctionTypes.ErrCouldNotFindSommTokenPrice, params.BaseCoinUnit),
			runsBefore:    func() {},
		},
		{
			name: "Usomm price update too old",
			beginAuctionRequest: BeginAuctionRequest{
				ctx:                        ctx.WithBlockHeight(int64(sommPrice.LastUpdatedBlock) + int64(auctionParams.PriceMaxBlockAge) + 1),
				startingTokensForSale:      auctionedSaleTokens,
				initialPriceDecreaseRate:   sdk.MustNewDecFromStr("0.05"),
				priceDecreaseBlockInterval: uint64(10),
				fundingModuleAccount:       permissionedFunder.GetName(),
				proceedsModuleAccount:      permissionedReciever.GetName(),
			},
			expectedError: errorsmod.Wrap(auctionTypes.ErrLastSommTokenPriceTooOld, params.BaseCoinUnit),
			runsBefore: func() {
				auctionKeeper.setTokenPrice(ctx, sommPrice)
			},
		},
		{
			name: "Validate basic canary 1 -- invalid initialPriceDecreaseRate lower bound",
			beginAuctionRequest: BeginAuctionRequest{
				ctx:                        ctx,
				startingTokensForSale:      auctionedSaleTokens,
				initialPriceDecreaseRate:   sdk.MustNewDecFromStr("0.0"),
				priceDecreaseBlockInterval: uint64(10),
				fundingModuleAccount:       permissionedFunder.GetName(),
				proceedsModuleAccount:      permissionedReciever.GetName(),
			},
			expectedError: errorsmod.Wrapf(auctionTypes.ErrInvalidInitialDecreaseRate, "Initial price decrease rate 0.000000000000000000"),
			runsBefore:    func() {},
		},
		{
			name: "Validate basic canary 2 -- invalid initialPriceDecreaseRate upper bound",
			beginAuctionRequest: BeginAuctionRequest{
				ctx:                        ctx,
				startingTokensForSale:      auctionedSaleTokens,
				initialPriceDecreaseRate:   sdk.MustNewDecFromStr("1.0"),
				priceDecreaseBlockInterval: uint64(10),
				fundingModuleAccount:       permissionedFunder.GetName(),
				proceedsModuleAccount:      permissionedReciever.GetName(),
			},
			expectedError: errorsmod.Wrapf(auctionTypes.ErrInvalidInitialDecreaseRate, "Initial price decrease rate 1.000000000000000000"),
			runsBefore:    func() {},
		},
		{
			name: "Cannot have 2 ongoing auctions for the same denom",
			beginAuctionRequest: BeginAuctionRequest{
				ctx:                        ctx,
				startingTokensForSale:      auctionedSaleTokens,
				initialPriceDecreaseRate:   sdk.MustNewDecFromStr("0.05"),
				priceDecreaseBlockInterval: uint64(10),
				fundingModuleAccount:       permissionedFunder.GetName(),
				proceedsModuleAccount:      permissionedReciever.GetName(),
			},
			expectedError: errorsmod.Wrapf(auctionTypes.ErrCannotStartTwoAuctionsForSameDenomSimultaneously, "Denom: %s", auctionedSaleTokens.Denom),
			runsBefore: func() {
				// Mock initial bank keeper fund transfer
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
					InitialUnitPriceInUsomm:    sdk.NewDec(1),
					CurrentUnitPriceInUsomm:    sdk.NewDec(1),
					RemainingTokensForSale:     auctionedSaleTokens,
					FundingModuleAccount:       permissionedFunder.GetName(),
					ProceedsModuleAccount:      permissionedReciever.GetName(),
				}
				require.Equal(expectedActiveAuction, createdAuction)
			},
		},
	}

	for _, tc := range tests {
		tc := tc // Redefine variable here due to passing it to function literal below (scopelint)
		suite.T().Run(fmt.Sprint(tc.name), func(t *testing.T) {
			// Run expected bank keeper functions, if any
			tc.runsBefore()

			call := func() error {
				return auctionKeeper.BeginAuction(
					tc.beginAuctionRequest.ctx,
					tc.beginAuctionRequest.startingTokensForSale,
					tc.beginAuctionRequest.initialPriceDecreaseRate,
					tc.beginAuctionRequest.priceDecreaseBlockInterval,
					tc.beginAuctionRequest.fundingModuleAccount,
					tc.beginAuctionRequest.proceedsModuleAccount,
				)
			}

			if tc.name[0:14] == "Validate basic" {
				require.Panics(func() { call() })
				return
			}

			err := call()

			// Verify errors are as expected
			require.Equal(tc.expectedError.Error(), err.Error())
		})
	}
}
