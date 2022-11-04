package types

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/stretchr/testify/require"
)

func TestGenesisValidate(t *testing.T) {
	testCases := []struct {
		name         string
		genesisState GenesisState
		expPass      bool
		err          error
	}{
		{
			name:         "Happy path -- default genesis",
			genesisState: DefaultGenesisState(),
			expPass:      true,
			err:          nil,
		},
		{
			name: "Happy path -- empty genesis",
			genesisState: GenesisState{
				Params:      DefaultParams(),
				Auctions:    []*Auction{},
				Bids:        []*Bid{},
				TokenPrices: []*TokenPrice{},
			},
			expPass: true,
			err:     nil,
		},
		{
			name: "Happy path -- custom populated genesis",
			genesisState: GenesisState{
				Params: DefaultParams(),
				Auctions: []*Auction{
					{
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
				},
				Bids: []*Bid{
					{
						Id:                        uint64(1),
						AuctionId:                 uint32(1),
						Bidder:                    "cosmos12svsksqaakc6r0gyxasf5el84946mp0svdl603",
						MaxBidInUsomm:             sdk.NewCoin("usomm", sdk.NewInt(100)),
						SaleTokenMinimumAmount:    sdk.NewCoin("gravity0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48", sdk.NewInt(50)),
						TotalFulfilledSaleTokens:  sdk.NewCoin("gravity0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48", sdk.NewInt(50)),
						SaleTokenUnitPriceInUsomm: sdk.MustNewDecFromStr("2.0"),
						TotalUsommPaid:            sdk.NewCoin("usomm", sdk.NewInt(100)),
					},
				},
				TokenPrices: []*TokenPrice{
					{
						Denom:            "usomm",
						UsdPrice:         sdk.MustNewDecFromStr("0.0008"),
						LastUpdatedBlock: uint64(123),
					},
					{
						Denom:            "gravity0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48",
						UsdPrice:         sdk.MustNewDecFromStr("0.032"),
						LastUpdatedBlock: uint64(321),
					},
				},
				LastAuctionId: uint32(1),
				LastBidId:     uint64(1),
			},
			expPass: true,
			err:     nil,
		},
		{
			name: "Validate basic canary -- Invalid auction",
			genesisState: GenesisState{
				Params: DefaultParams(),
				Auctions: []*Auction{
					{
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
				},
				Bids: []*Bid{
					{
						Id:                        uint64(1),
						AuctionId:                 uint32(1),
						Bidder:                    "cosmos12svsksqaakc6r0gyxasf5el84946mp0svdl603",
						MaxBidInUsomm:             sdk.NewCoin("usomm", sdk.NewInt(100)),
						SaleTokenMinimumAmount:    sdk.NewCoin("gravity0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48", sdk.NewInt(50)),
						TotalFulfilledSaleTokens:  sdk.NewCoin("gravity0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48", sdk.NewInt(50)),
						SaleTokenUnitPriceInUsomm: sdk.MustNewDecFromStr("2.0"),
						TotalUsommPaid:            sdk.NewCoin("usomm", sdk.NewInt(100)),
					},
				},
				TokenPrices: []*TokenPrice{
					{
						Denom:            "usomm",
						UsdPrice:         sdk.MustNewDecFromStr("0.0008"),
						LastUpdatedBlock: uint64(123),
					},
					{
						Denom:            "gravity0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48",
						UsdPrice:         sdk.MustNewDecFromStr("0.032"),
						LastUpdatedBlock: uint64(321),
					},
				},
				LastAuctionId: uint32(1),
				LastBidId:     uint64(1),
			},
			expPass: false,
			err:     sdkerrors.Wrapf(ErrAuctionIDMustBeNonZero, "id: 0"),
		},
		{
			name: "Validate basic canary -- Invalid bid",
			genesisState: GenesisState{
				Params: DefaultParams(),
				Auctions: []*Auction{
					{
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
				},
				Bids: []*Bid{
					{
						Id:                        uint64(0),
						AuctionId:                 uint32(1),
						Bidder:                    "cosmos12svsksqaakc6r0gyxasf5el84946mp0svdl603",
						MaxBidInUsomm:             sdk.NewCoin("usomm", sdk.NewInt(100)),
						SaleTokenMinimumAmount:    sdk.NewCoin("gravity0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48", sdk.NewInt(50)),
						TotalFulfilledSaleTokens:  sdk.NewCoin("gravity0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48", sdk.NewInt(50)),
						SaleTokenUnitPriceInUsomm: sdk.MustNewDecFromStr("2.0"),
						TotalUsommPaid:            sdk.NewCoin("usomm", sdk.NewInt(100)),
					},
				},
				TokenPrices: []*TokenPrice{
					{
						Denom:            "usomm",
						UsdPrice:         sdk.MustNewDecFromStr("0.0008"),
						LastUpdatedBlock: uint64(123),
					},
					{
						Denom:            "gravity0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48",
						UsdPrice:         sdk.MustNewDecFromStr("0.032"),
						LastUpdatedBlock: uint64(321),
					},
				},
				LastAuctionId: uint32(1),
				LastBidId:     uint64(1),
			},
			expPass: false,
			err:     sdkerrors.Wrapf(ErrBidIDMustBeNonZero, "id: 0"),
		},
		{
			name: "Validate basic canary -- Invalid token price",
			genesisState: GenesisState{
				Params: DefaultParams(),
				Auctions: []*Auction{
					{
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
				},
				Bids: []*Bid{
					{
						Id:                        uint64(1),
						AuctionId:                 uint32(1),
						Bidder:                    "cosmos12svsksqaakc6r0gyxasf5el84946mp0svdl603",
						MaxBidInUsomm:             sdk.NewCoin("usomm", sdk.NewInt(100)),
						SaleTokenMinimumAmount:    sdk.NewCoin("gravity0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48", sdk.NewInt(50)),
						TotalFulfilledSaleTokens:  sdk.NewCoin("gravity0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48", sdk.NewInt(50)),
						SaleTokenUnitPriceInUsomm: sdk.MustNewDecFromStr("2.0"),
						TotalUsommPaid:            sdk.NewCoin("usomm", sdk.NewInt(100)),
					},
				},
				TokenPrices: []*TokenPrice{
					{
						Denom:            "usdc",
						UsdPrice:         sdk.MustNewDecFromStr("0.0008"),
						LastUpdatedBlock: uint64(123),
					},
					{
						Denom:            "gravity0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48",
						UsdPrice:         sdk.MustNewDecFromStr("0.032"),
						LastUpdatedBlock: uint64(321),
					},
				},
				LastAuctionId: uint32(1),
				LastBidId:     uint64(1),
			},
			expPass: false,
			err:     sdkerrors.Wrapf(ErrInvalidTokenPriceDenom, "denom: usdc"),
		},
	}

	for _, tc := range testCases {
		err := tc.genesisState.Validate()
		if tc.expPass {
			require.NoError(t, err, tc.name)
			require.Nil(t, err)
		} else {
			require.Error(t, err, tc.name)
			require.Equal(t, tc.err.Error(), err.Error())
		}
	}
}
