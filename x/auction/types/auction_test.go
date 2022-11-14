package types

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/peggyjv/sommelier/v4/app/params"
	"github.com/stretchr/testify/require"
)

var (
	cosmosAddress1   string = "cosmos16zrkzad482haunrn25ywvwy6fclh3vh7r0hcny"
	cosmosAddress2   string = "cosmos12svsksqaakc6r0gyxasf5el84946mp0svdl603"
)

func TestAuctionValidate(t *testing.T) {
	testCases := []struct {
		name    string
		auction Auction
		expPass bool
		err     error
	}{
		{
			name: "Happy path",
			auction: Auction{
				Id:                         uint32(1),
				StartingTokensForSale:      sdk.NewCoin("gravity0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48", sdk.NewIntFromUint64(1000)),
				StartBlock:                 uint64(200),
				EndBlock:                   uint64(0),
				InitialPriceDecreaseRate:   sdk.MustNewDecFromStr("0.05"),
				CurrentPriceDecreaseRate:   sdk.MustNewDecFromStr("0.05"),
				PriceDecreaseBlockInterval: uint64(10),
				InitialUnitPriceInUsomm:    sdk.MustNewDecFromStr("20.0"),
				CurrentUnitPriceInUsomm:    sdk.MustNewDecFromStr("20.0"),
				RemainingTokensForSale:     sdk.NewCoin("gravity0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48", sdk.NewIntFromUint64(900)),
				FundingModuleAccount:       "someModule",
				ProceedsModuleAccount:      "someModule",
			},
			expPass: true,
			err:     nil,
		},
		{
			name: "Auction ID cannot be 0",
			auction: Auction{
				Id:                         uint32(0),
				StartingTokensForSale:      sdk.NewCoin("gravity0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48", sdk.NewIntFromUint64(1000)),
				StartBlock:                 uint64(200),
				EndBlock:                   uint64(0),
				InitialPriceDecreaseRate:   sdk.MustNewDecFromStr("0.05"),
				CurrentPriceDecreaseRate:   sdk.MustNewDecFromStr("0.05"),
				PriceDecreaseBlockInterval: uint64(10),
				InitialUnitPriceInUsomm:    sdk.MustNewDecFromStr("20.0"),
				CurrentUnitPriceInUsomm:    sdk.MustNewDecFromStr("20.0"),
				RemainingTokensForSale:     sdk.NewCoin("gravity0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48", sdk.NewIntFromUint64(900)),
				FundingModuleAccount:       "someModule",
				ProceedsModuleAccount:      "someModule",
			},
			expPass: false,
			err:     sdkerrors.Wrapf(ErrAuctionIDMustBeNonZero, "id: 0"),
		},
		{
			name: "Starting tokens for sale must be positive",
			auction: Auction{
				Id:                         uint32(1),
				StartingTokensForSale:      sdk.NewCoin("gravity0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48", sdk.NewIntFromUint64(0)),
				StartBlock:                 uint64(200),
				EndBlock:                   uint64(0),
				InitialPriceDecreaseRate:   sdk.MustNewDecFromStr("0.05"),
				CurrentPriceDecreaseRate:   sdk.MustNewDecFromStr("0.05"),
				PriceDecreaseBlockInterval: uint64(10),
				InitialUnitPriceInUsomm:    sdk.MustNewDecFromStr("20.0"),
				CurrentUnitPriceInUsomm:    sdk.MustNewDecFromStr("20.0"),
				RemainingTokensForSale:     sdk.NewCoin("gravity0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48", sdk.NewIntFromUint64(900)),
				FundingModuleAccount:       "someModule",
				ProceedsModuleAccount:      "someModule",
			},
			expPass: false,
			err:     sdkerrors.Wrapf(ErrAuctionStartingAmountMustBePositve, "Starting tokens for sale: %s", sdk.NewCoin("gravity0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48", sdk.NewIntFromUint64(0)).String()),
		},
		{
			name: "Starting tokens for sale cannot be usomm",
			auction: Auction{
				Id:                         uint32(1),
				StartingTokensForSale:      sdk.NewCoin("usomm", sdk.NewIntFromUint64(1000)),
				StartBlock:                 uint64(200),
				EndBlock:                   uint64(0),
				InitialPriceDecreaseRate:   sdk.MustNewDecFromStr("0.05"),
				CurrentPriceDecreaseRate:   sdk.MustNewDecFromStr("0.05"),
				PriceDecreaseBlockInterval: uint64(10),
				InitialUnitPriceInUsomm:    sdk.MustNewDecFromStr("20.0"),
				CurrentUnitPriceInUsomm:    sdk.MustNewDecFromStr("20.0"),
				RemainingTokensForSale:     sdk.NewCoin("gravity0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48", sdk.NewIntFromUint64(900)),
				FundingModuleAccount:       "someModule",
				ProceedsModuleAccount:      "someModule",
			},
			expPass: false,
			err:     sdkerrors.Wrapf(ErrCannotAuctionUsomm, "Starting denom tokens for sale: %s", params.BaseCoinUnit),
		},
		{
			name: "Start block must be positive",
			auction: Auction{
				Id:                         uint32(1),
				StartingTokensForSale:      sdk.NewCoin("gravity0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48", sdk.NewIntFromUint64(1000)),
				StartBlock:                 uint64(0),
				EndBlock:                   uint64(0),
				InitialPriceDecreaseRate:   sdk.MustNewDecFromStr("0.05"),
				CurrentPriceDecreaseRate:   sdk.MustNewDecFromStr("0.05"),
				PriceDecreaseBlockInterval: uint64(10),
				InitialUnitPriceInUsomm:    sdk.MustNewDecFromStr("20.0"),
				CurrentUnitPriceInUsomm:    sdk.MustNewDecFromStr("20.0"),
				RemainingTokensForSale:     sdk.NewCoin("gravity0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48", sdk.NewIntFromUint64(900)),
				FundingModuleAccount:       "someModule",
				ProceedsModuleAccount:      "someModule",
			},
			expPass: false,
			err:     sdkerrors.Wrapf(ErrInvalidStartBlock, "start block: 0"),
		},
		{
			name: "Initial decrease rate cannot be <= 0",
			auction: Auction{
				Id:                         uint32(1),
				StartingTokensForSale:      sdk.NewCoin("gravity0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48", sdk.NewIntFromUint64(1000)),
				StartBlock:                 uint64(200),
				EndBlock:                   uint64(0),
				InitialPriceDecreaseRate:   sdk.MustNewDecFromStr("0.0"),
				CurrentPriceDecreaseRate:   sdk.MustNewDecFromStr("0.05"),
				PriceDecreaseBlockInterval: uint64(10),
				InitialUnitPriceInUsomm:    sdk.MustNewDecFromStr("20.0"),
				CurrentUnitPriceInUsomm:    sdk.MustNewDecFromStr("20.0"),
				RemainingTokensForSale:     sdk.NewCoin("gravity0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48", sdk.NewIntFromUint64(900)),
				FundingModuleAccount:       "someModule",
				ProceedsModuleAccount:      "someModule",
			},
			expPass: false,
			err:     sdkerrors.Wrapf(ErrInvalidInitialDecreaseRate, "Inital price decrease rate %s", sdk.MustNewDecFromStr("0.0").String()),
		},
		{
			name: "Initial decrease rate cannot be >= 0",
			auction: Auction{
				Id:                         uint32(1),
				StartingTokensForSale:      sdk.NewCoin("gravity0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48", sdk.NewIntFromUint64(1000)),
				StartBlock:                 uint64(200),
				EndBlock:                   uint64(0),
				InitialPriceDecreaseRate:   sdk.MustNewDecFromStr("1.0"),
				CurrentPriceDecreaseRate:   sdk.MustNewDecFromStr("0.05"),
				PriceDecreaseBlockInterval: uint64(10),
				InitialUnitPriceInUsomm:    sdk.MustNewDecFromStr("20.0"),
				CurrentUnitPriceInUsomm:    sdk.MustNewDecFromStr("20.0"),
				RemainingTokensForSale:     sdk.NewCoin("gravity0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48", sdk.NewIntFromUint64(900)),
				FundingModuleAccount:       "someModule",
				ProceedsModuleAccount:      "someModule",
			},
			expPass: false,
			err:     sdkerrors.Wrapf(ErrInvalidInitialDecreaseRate, "Inital price decrease rate %s", sdk.MustNewDecFromStr("1.0").String()),
		},
		{
			name: "Current decrease rate cannot be <= 0",
			auction: Auction{
				Id:                         uint32(1),
				StartingTokensForSale:      sdk.NewCoin("gravity0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48", sdk.NewIntFromUint64(1000)),
				StartBlock:                 uint64(200),
				EndBlock:                   uint64(0),
				InitialPriceDecreaseRate:   sdk.MustNewDecFromStr("0.05"),
				CurrentPriceDecreaseRate:   sdk.MustNewDecFromStr("0.0"),
				PriceDecreaseBlockInterval: uint64(10),
				InitialUnitPriceInUsomm:    sdk.MustNewDecFromStr("20.0"),
				CurrentUnitPriceInUsomm:    sdk.MustNewDecFromStr("20.0"),
				RemainingTokensForSale:     sdk.NewCoin("gravity0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48", sdk.NewIntFromUint64(900)),
				FundingModuleAccount:       "someModule",
				ProceedsModuleAccount:      "someModule",
			},
			expPass: false,
			err:     sdkerrors.Wrapf(ErrInvalidCurrentDecreaseRate, "Current price decrease rate %s", sdk.MustNewDecFromStr("0.0").String()),
		},
		{
			name: "Current decrease rate cannot be >= 0",
			auction: Auction{
				Id:                         uint32(1),
				StartingTokensForSale:      sdk.NewCoin("gravity0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48", sdk.NewIntFromUint64(1000)),
				StartBlock:                 uint64(200),
				EndBlock:                   uint64(0),
				InitialPriceDecreaseRate:   sdk.MustNewDecFromStr("0.05"),
				CurrentPriceDecreaseRate:   sdk.MustNewDecFromStr("1.0"),
				PriceDecreaseBlockInterval: uint64(10),
				InitialUnitPriceInUsomm:    sdk.MustNewDecFromStr("20.0"),
				CurrentUnitPriceInUsomm:    sdk.MustNewDecFromStr("20.0"),
				RemainingTokensForSale:     sdk.NewCoin("gravity0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48", sdk.NewIntFromUint64(900)),
				FundingModuleAccount:       "someModule",
				ProceedsModuleAccount:      "someModule",
			},
			expPass: false,
			err:     sdkerrors.Wrapf(ErrInvalidCurrentDecreaseRate, "Current price decrease rate %s", sdk.MustNewDecFromStr("1.0").String()),
		},
		{
			name: "Price decrease block interval cannot be 0",
			auction: Auction{
				Id:                         uint32(1),
				StartingTokensForSale:      sdk.NewCoin("gravity0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48", sdk.NewIntFromUint64(1000)),
				StartBlock:                 uint64(200),
				EndBlock:                   uint64(0),
				InitialPriceDecreaseRate:   sdk.MustNewDecFromStr("0.05"),
				CurrentPriceDecreaseRate:   sdk.MustNewDecFromStr("0.05"),
				PriceDecreaseBlockInterval: uint64(0),
				InitialUnitPriceInUsomm:    sdk.MustNewDecFromStr("20.0"),
				CurrentUnitPriceInUsomm:    sdk.MustNewDecFromStr("20.0"),
				RemainingTokensForSale:     sdk.NewCoin("gravity0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48", sdk.NewIntFromUint64(900)),
				FundingModuleAccount:       "someModule",
				ProceedsModuleAccount:      "someModule",
			},
			expPass: false,
			err:     sdkerrors.Wrapf(ErrInvalidBlockDecreaseInterval, "price decrease block interval: 0"),
		},
		{
			name: "Initial unit price in usomm must be postive",
			auction: Auction{
				Id:                         uint32(1),
				StartingTokensForSale:      sdk.NewCoin("gravity0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48", sdk.NewIntFromUint64(1000)),
				StartBlock:                 uint64(200),
				EndBlock:                   uint64(10),
				InitialPriceDecreaseRate:   sdk.MustNewDecFromStr("0.05"),
				CurrentPriceDecreaseRate:   sdk.MustNewDecFromStr("0.05"),
				PriceDecreaseBlockInterval: uint64(10),
				InitialUnitPriceInUsomm:    sdk.MustNewDecFromStr("0.0"),
				CurrentUnitPriceInUsomm:    sdk.MustNewDecFromStr("20.0"),
				RemainingTokensForSale:     sdk.NewCoin("gravity0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48", sdk.NewIntFromUint64(900)),
				FundingModuleAccount:       "someModule",
				ProceedsModuleAccount:      "someModule",
			},
			expPass: false,
			err:     sdkerrors.Wrapf(ErrPriceMustBePositive, "initial unit price in usomm: %s", sdk.MustNewDecFromStr("0.0").String()),
		},
		{
			name: "Current unit price in usomm must be postive",
			auction: Auction{
				Id:                         uint32(1),
				StartingTokensForSale:      sdk.NewCoin("gravity0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48", sdk.NewIntFromUint64(1000)),
				StartBlock:                 uint64(200),
				EndBlock:                   uint64(10),
				InitialPriceDecreaseRate:   sdk.MustNewDecFromStr("0.05"),
				CurrentPriceDecreaseRate:   sdk.MustNewDecFromStr("0.05"),
				PriceDecreaseBlockInterval: uint64(10),
				InitialUnitPriceInUsomm:    sdk.MustNewDecFromStr("20.0"),
				CurrentUnitPriceInUsomm:    sdk.MustNewDecFromStr("0.0"),
				RemainingTokensForSale:     sdk.NewCoin("gravity0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48", sdk.NewIntFromUint64(900)),
				FundingModuleAccount:       "someModule",
				ProceedsModuleAccount:      "someModule",
			},
			expPass: false,
			err:     sdkerrors.Wrapf(ErrPriceMustBePositive, "current unit price in usomm: %s", sdk.MustNewDecFromStr("0.0").String()),
		},

		{
			name: "Funding Module account cannot be empty",
			auction: Auction{
				Id:                         uint32(1),
				StartingTokensForSale:      sdk.NewCoin("gravity0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48", sdk.NewIntFromUint64(1000)),
				StartBlock:                 uint64(200),
				EndBlock:                   uint64(10),
				InitialPriceDecreaseRate:   sdk.MustNewDecFromStr("0.05"),
				CurrentPriceDecreaseRate:   sdk.MustNewDecFromStr("0.05"),
				PriceDecreaseBlockInterval: uint64(10),
				InitialUnitPriceInUsomm:    sdk.MustNewDecFromStr("20.0"),
				CurrentUnitPriceInUsomm:    sdk.MustNewDecFromStr("20.0"),
				RemainingTokensForSale:     sdk.NewCoin("gravity0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48", sdk.NewIntFromUint64(900)),
				FundingModuleAccount:       "",
				ProceedsModuleAccount:      "someModule",
			},
			expPass: false,
			err:     sdkerrors.Wrapf(ErrUnauthorizedFundingModule, "funding module account: "),
		},
		{
			name: "Proceeds Module account cannot be empty",
			auction: Auction{
				Id:                         uint32(1),
				StartingTokensForSale:      sdk.NewCoin("gravity0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48", sdk.NewIntFromUint64(1000)),
				StartBlock:                 uint64(200),
				EndBlock:                   uint64(10),
				InitialPriceDecreaseRate:   sdk.MustNewDecFromStr("0.05"),
				CurrentPriceDecreaseRate:   sdk.MustNewDecFromStr("0.05"),
				PriceDecreaseBlockInterval: uint64(10),
				InitialUnitPriceInUsomm:    sdk.MustNewDecFromStr("20.0"),
				CurrentUnitPriceInUsomm:    sdk.MustNewDecFromStr("20.0"),
				RemainingTokensForSale:     sdk.NewCoin("gravity0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48", sdk.NewIntFromUint64(900)),
				FundingModuleAccount:       "someModule",
				ProceedsModuleAccount:      "",
			},
			expPass: false,
			err:     sdkerrors.Wrapf(ErrUnauthorizedFundingModule, "proceeds module account: "),
		},
	}

	for _, tc := range testCases {
		err := tc.auction.ValidateBasic()
		if tc.expPass {
			require.NoError(t, err, tc.name)
			require.Nil(t, err)
		} else {
			require.Error(t, err, tc.name)
			require.Equal(t, tc.err.Error(), err.Error())
		}
	}
}

func TestBidValidate(t *testing.T) {
	testCases := []struct {
		name    string
		bid     Bid
		expPass bool
		err     error
	}{
		{
			name: "Happy path",
			bid: Bid{
				Id:                        uint64(1),
				AuctionId:                 uint32(1),
				Bidder:                    cosmosAddress2,
				MaxBidInUsomm:             sdk.NewCoin("usomm", sdk.NewInt(100)),
				SaleTokenMinimumAmount:    sdk.NewCoin("gravity0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48", sdk.NewInt(50)),
				TotalFulfilledSaleTokens:  sdk.NewCoin("gravity0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48", sdk.NewInt(50)),
				SaleTokenUnitPriceInUsomm: sdk.MustNewDecFromStr("2.0"),
				TotalUsommPaid:            sdk.NewCoin("usomm", sdk.NewInt(100)),
			},
			expPass: true,
			err:     nil,
		},
		{
			name: "Bid ID cannot be 0",
			bid: Bid{
				Id:                        uint64(0),
				AuctionId:                 uint32(1),
				Bidder:                    cosmosAddress2,
				MaxBidInUsomm:             sdk.NewCoin("usomm", sdk.NewInt(100)),
				SaleTokenMinimumAmount:    sdk.NewCoin("gravity0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48", sdk.NewInt(50)),
				TotalFulfilledSaleTokens:  sdk.NewCoin("gravity0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48", sdk.NewInt(50)),
				SaleTokenUnitPriceInUsomm: sdk.MustNewDecFromStr("2.0"),
				TotalUsommPaid:            sdk.NewCoin("usomm", sdk.NewInt(100)),
			},
			expPass: false,
			err:     sdkerrors.Wrapf(ErrBidIDMustBeNonZero, "id: 0"),
		},
		{
			name: "Auction ID cannot be 0",
			bid: Bid{
				Id:                        uint64(1),
				AuctionId:                 uint32(0),
				Bidder:                    cosmosAddress2,
				MaxBidInUsomm:             sdk.NewCoin("usomm", sdk.NewInt(100)),
				SaleTokenMinimumAmount:    sdk.NewCoin("gravity0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48", sdk.NewInt(50)),
				TotalFulfilledSaleTokens:  sdk.NewCoin("gravity0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48", sdk.NewInt(50)),
				SaleTokenUnitPriceInUsomm: sdk.MustNewDecFromStr("2.0"),
				TotalUsommPaid:            sdk.NewCoin("usomm", sdk.NewInt(100)),
			},
			expPass: false,
			err:     sdkerrors.Wrapf(ErrAuctionIDMustBeNonZero, "id: 0"),
		},
		{
			name: "Bidder cannot be empty",
			bid: Bid{
				Id:                        uint64(1),
				AuctionId:                 uint32(1),
				Bidder:                    "",
				MaxBidInUsomm:             sdk.NewCoin("usomm", sdk.NewInt(100)),
				SaleTokenMinimumAmount:    sdk.NewCoin("gravity0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48", sdk.NewInt(50)),
				TotalFulfilledSaleTokens:  sdk.NewCoin("gravity0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48", sdk.NewInt(50)),
				SaleTokenUnitPriceInUsomm: sdk.MustNewDecFromStr("2.0"),
				TotalUsommPaid:            sdk.NewCoin("usomm", sdk.NewInt(100)),
			},
			expPass: false,
			err:     sdkerrors.Wrapf(ErrAddressExpected, "bidder: "),
		},
		{
			name: "Bidder must be a valid bech32 address",
			bid: Bid{
				Id:                        uint64(1),
				AuctionId:                 uint32(1),
				Bidder:                    "ironman",
				MaxBidInUsomm:             sdk.NewCoin("usomm", sdk.NewInt(100)),
				SaleTokenMinimumAmount:    sdk.NewCoin("gravity0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48", sdk.NewInt(50)),
				TotalFulfilledSaleTokens:  sdk.NewCoin("gravity0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48", sdk.NewInt(50)),
				SaleTokenUnitPriceInUsomm: sdk.MustNewDecFromStr("2.0"),
				TotalUsommPaid:            sdk.NewCoin("usomm", sdk.NewInt(100)),
			},
			expPass: false,
			err:     sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "decoding bech32 failed: invalid bech32 string length 7"),
		},
		{
			name: "Bid must be positive",
			bid: Bid{
				Id:                        uint64(1),
				AuctionId:                 uint32(1),
				Bidder:                    cosmosAddress2,
				MaxBidInUsomm:             sdk.NewCoin("usomm", sdk.NewInt(0)),
				SaleTokenMinimumAmount:    sdk.NewCoin("gravity0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48", sdk.NewInt(50)),
				TotalFulfilledSaleTokens:  sdk.NewCoin("gravity0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48", sdk.NewInt(50)),
				SaleTokenUnitPriceInUsomm: sdk.MustNewDecFromStr("2.0"),
				TotalUsommPaid:            sdk.NewCoin("usomm", sdk.NewInt(100)),
			},
			expPass: false,
			err:     sdkerrors.Wrapf(ErrBidAmountMustBePositive, "bid amount in usomm: %s", sdk.NewCoin("usomm", sdk.NewInt(0)).String()),
		},
		{
			name: "Bid must be in usomm",
			bid: Bid{
				Id:                        uint64(1),
				AuctionId:                 uint32(1),
				Bidder:                    cosmosAddress2,
				MaxBidInUsomm:             sdk.NewCoin("usdc", sdk.NewInt(100)),
				SaleTokenMinimumAmount:    sdk.NewCoin("gravity0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48", sdk.NewInt(50)),
				TotalFulfilledSaleTokens:  sdk.NewCoin("gravity0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48", sdk.NewInt(50)),
				SaleTokenUnitPriceInUsomm: sdk.MustNewDecFromStr("2.0"),
				TotalUsommPaid:            sdk.NewCoin("usomm", sdk.NewInt(100)),
			},
			expPass: false,
			err:     sdkerrors.Wrapf(ErrBidMustBeInUsomm, "bid: %s", sdk.NewCoin("usdc", sdk.NewInt(100)).String()),
		},
		{
			name: "Sale token must be gravity prefixed",
			bid: Bid{
				Id:                        uint64(1),
				AuctionId:                 uint32(1),
				Bidder:                    cosmosAddress2,
				MaxBidInUsomm:             sdk.NewCoin("usomm", sdk.NewInt(100)),
				SaleTokenMinimumAmount:    sdk.NewCoin("usdc", sdk.NewInt(50)),
				TotalFulfilledSaleTokens:  sdk.NewCoin("gravity0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48", sdk.NewInt(50)),
				SaleTokenUnitPriceInUsomm: sdk.MustNewDecFromStr("2.0"),
				TotalUsommPaid:            sdk.NewCoin("usomm", sdk.NewInt(100)),
			},
			expPass: false,
			err:     sdkerrors.Wrapf(ErrInvalidTokenBeingBidOn, "sale token: %s", sdk.NewCoin("usdc", sdk.NewInt(50)).String()),
		},
		{
			name: "Sale token amount must be positive",
			bid: Bid{
				Id:                        uint64(1),
				AuctionId:                 uint32(1),
				Bidder:                    cosmosAddress2,
				MaxBidInUsomm:             sdk.NewCoin("usomm", sdk.NewInt(100)),
				SaleTokenMinimumAmount:    sdk.NewCoin("gravity0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48", sdk.NewInt(0)),
				TotalFulfilledSaleTokens:  sdk.NewCoin("gravity0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48", sdk.NewInt(50)),
				SaleTokenUnitPriceInUsomm: sdk.MustNewDecFromStr("2.0"),
				TotalUsommPaid:            sdk.NewCoin("usomm", sdk.NewInt(100)),
			},
			expPass: false,
			err:     sdkerrors.Wrapf(ErrMinimumAmountMustBePositive, "sale token amount: %s", sdk.NewCoin("gravity0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48", sdk.NewInt(0)).String()),
		},
		{
			name: "Sale token unit price must be in usomm",
			bid: Bid{
				Id:                        uint64(1),
				AuctionId:                 uint32(1),
				Bidder:                    cosmosAddress2,
				MaxBidInUsomm:             sdk.NewCoin("usomm", sdk.NewInt(100)),
				SaleTokenMinimumAmount:    sdk.NewCoin("gravity0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48", sdk.NewInt(50)),
				TotalFulfilledSaleTokens:  sdk.NewCoin("gravity0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48", sdk.NewInt(50)),
				SaleTokenUnitPriceInUsomm: sdk.MustNewDecFromStr("0.0"),
				TotalUsommPaid:            sdk.NewCoin("usomm", sdk.NewInt(100)),
			},
			expPass: false,
			err:     sdkerrors.Wrapf(ErrBidUnitPriceInUsommMustBePositive, "sale token unit price: %s", sdk.MustNewDecFromStr("0.0").String()),
		},
		{
			name: "Total usomm paid denom must be usomm",
			bid: Bid{
				Id:                        uint64(1),
				AuctionId:                 uint32(1),
				Bidder:                    cosmosAddress2,
				MaxBidInUsomm:             sdk.NewCoin("usomm", sdk.NewInt(100)),
				SaleTokenMinimumAmount:    sdk.NewCoin("gravity0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48", sdk.NewInt(50)),
				TotalFulfilledSaleTokens:  sdk.NewCoin("gravity0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48", sdk.NewInt(50)),
				SaleTokenUnitPriceInUsomm: sdk.MustNewDecFromStr("2.0"),
				TotalUsommPaid:            sdk.NewCoin("usdc", sdk.NewInt(100)),
			},
			expPass: false,
			err:     sdkerrors.Wrapf(ErrBidMustBeInUsomm, "payment denom: usdc"),
		},
	}

	for _, tc := range testCases {
		err := tc.bid.ValidateBasic()
		if tc.expPass {
			require.NoError(t, err, tc.name)
			require.Nil(t, err)
		} else {
			require.Error(t, err, tc.name)
			require.Equal(t, tc.err.Error(), err.Error())
		}
	}
}

func TestTokenPriceValidate(t *testing.T) {
	testCases := []struct {
		name       string
		tokenPrice TokenPrice
		expPass    bool
		err        error
	}{
		{
			name: "Happy path -- usomm",
			tokenPrice: TokenPrice{
				Denom:            "usomm",
				UsdPrice:         sdk.MustNewDecFromStr("0.0008"),
				LastUpdatedBlock: uint64(321),
			},
			expPass: true,
			err:     nil,
		},
		{
			name: "Happy path -- gravity denom",
			tokenPrice: TokenPrice{
				Denom:            "gravity0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48",
				UsdPrice:         sdk.MustNewDecFromStr("0.001"),
				LastUpdatedBlock: uint64(321),
			},
			expPass: true,
			err:     nil,
		},
		{
			name: "Denom cannot be empty",
			tokenPrice: TokenPrice{
				Denom:            "",
				UsdPrice:         sdk.MustNewDecFromStr("0.0008"),
				LastUpdatedBlock: uint64(321),
			},
			expPass: false,
			err:     sdkerrors.Wrapf(ErrDenomCannotBeEmpty, "price denom: "),
		},
		{
			name: "Price must be positive",
			tokenPrice: TokenPrice{
				Denom:            "usomm",
				UsdPrice:         sdk.MustNewDecFromStr("0.0"),
				LastUpdatedBlock: uint64(321),
			},
			expPass: false,
			err:     sdkerrors.Wrapf(ErrPriceMustBePositive, "usd price: %s", sdk.MustNewDecFromStr("0.0").String()),
		},
		{
			name: "Last updated block cannot be 0",
			tokenPrice: TokenPrice{
				Denom:            "usomm",
				UsdPrice:         sdk.MustNewDecFromStr("0.0008"),
				LastUpdatedBlock: uint64(0),
			},
			expPass: false,
			err:     sdkerrors.Wrapf(ErrInvalidLastUpdatedBlock, "block: 0"),
		},
		{
			name: "Token price must be usomm or gravity prefixed",
			tokenPrice: TokenPrice{
				Denom:            "usdc",
				UsdPrice:         sdk.MustNewDecFromStr("1.0"),
				LastUpdatedBlock: uint64(321),
			},
			expPass: false,
			err:     sdkerrors.Wrapf(ErrInvalidTokenPriceDenom, "denom: usdc"),
		},
	}

	for _, tc := range testCases {
		err := tc.tokenPrice.ValidateBasic()
		if tc.expPass {
			require.NoError(t, err, tc.name)
			require.Nil(t, err)
		} else {
			require.Error(t, err, tc.name)
			require.Equal(t, tc.err.Error(), err.Error())
		}
	}
}

func TestProposedTokenPriceValidate(t *testing.T) {
	testCases := []struct {
		name               string
		proposedTokenPrice ProposedTokenPrice
		expPass            bool
		err                error
	}{
		{
			name: "Happy path -- usomm",
			proposedTokenPrice: ProposedTokenPrice{
				Denom:    "usomm",
				UsdPrice: sdk.MustNewDecFromStr("0.0008"),
			},
			expPass: true,
			err:     nil,
		},
		{
			name: "Happy path -- gravity denom",
			proposedTokenPrice: ProposedTokenPrice{
				Denom:    "gravity0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48",
				UsdPrice: sdk.MustNewDecFromStr("0.001"),
			},
			expPass: true,
			err:     nil,
		},
		{
			name: "Denom cannot be empty",
			proposedTokenPrice: ProposedTokenPrice{
				Denom:    "",
				UsdPrice: sdk.MustNewDecFromStr("0.0008"),
			},
			expPass: false,
			err:     sdkerrors.Wrapf(ErrDenomCannotBeEmpty, "price denom: "),
		},
		{
			name: "Price must be positive",
			proposedTokenPrice: ProposedTokenPrice{
				Denom:    "usomm",
				UsdPrice: sdk.MustNewDecFromStr("0.0"),
			},
			expPass: false,
			err:     sdkerrors.Wrapf(ErrPriceMustBePositive, "usd price: %s", sdk.MustNewDecFromStr("0.0").String()),
		},
		{
			name: "Token price must be usomm or gravity prefixed",
			proposedTokenPrice: ProposedTokenPrice{
				Denom:    "usdc",
				UsdPrice: sdk.MustNewDecFromStr("1.0"),
			},
			expPass: false,
			err:     sdkerrors.Wrapf(ErrInvalidTokenPriceDenom, "denom: usdc"),
		},
	}

	for _, tc := range testCases {
		err := tc.proposedTokenPrice.ValidateBasic()
		if tc.expPass {
			require.NoError(t, err, tc.name)
			require.Nil(t, err)
		} else {
			require.Error(t, err, tc.name)
			require.Equal(t, tc.err.Error(), err.Error())
		}
	}
}
