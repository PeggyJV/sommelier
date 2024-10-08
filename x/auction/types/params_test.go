package types

import (
	"testing"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
)

func TestParamsValidate(t *testing.T) {
	testCases := []struct {
		name    string
		params  Params
		expPass bool
		err     error
	}{
		{
			name:    "Happy path -- default params",
			params:  DefaultParams(),
			expPass: true,
			err:     nil,
		},
		{
			name: "Happy path -- custom params",
			params: Params{
				PriceMaxBlockAge:                     uint64(1000),
				MinimumBidInUsomm:                    uint64(500),
				MinimumSaleTokensUsdValue:            sdk.MustNewDecFromStr("1.0"),
				AuctionMaxBlockAge:                   uint64(100),
				AuctionPriceDecreaseAccelerationRate: sdk.MustNewDecFromStr("0.1"),
				AuctionBurnRate:                      sdk.MustNewDecFromStr("0.1"),
			},
			expPass: true,
			err:     nil,
		},
		{
			name: "Token price max block age cannot be 0",
			params: Params{
				PriceMaxBlockAge:                     uint64(0),
				MinimumBidInUsomm:                    uint64(500),
				MinimumSaleTokensUsdValue:            sdk.MustNewDecFromStr("1.0"),
				AuctionMaxBlockAge:                   uint64(100),
				AuctionPriceDecreaseAccelerationRate: sdk.MustNewDecFromStr("0.1"),
				AuctionBurnRate:                      sdk.MustNewDecFromStr("0.1"),
			},
			expPass: false,
			err:     errorsmod.Wrapf(ErrTokenPriceMaxBlockAgeMustBePositive, "value: 0"),
		},
		{
			name: "Auction price decrease acceleration rate bounds check lower end",
			params: Params{
				PriceMaxBlockAge:                     uint64(1000),
				MinimumBidInUsomm:                    uint64(500),
				MinimumSaleTokensUsdValue:            sdk.MustNewDecFromStr("1.0"),
				AuctionMaxBlockAge:                   uint64(100),
				AuctionPriceDecreaseAccelerationRate: sdk.MustNewDecFromStr("-0.01"),
				AuctionBurnRate:                      sdk.MustNewDecFromStr("0.1"),
			},
			expPass: false,
			err:     errorsmod.Wrapf(ErrInvalidAuctionPriceDecreaseAccelerationRateParam, "auction price decrease acceleration rate must be between 0 and 1 inclusive (0%% to 100%%)"),
		},
		{
			name: "Auction price decrease acceleration rate bounds check upper end",
			params: Params{
				PriceMaxBlockAge:                     uint64(1000),
				MinimumBidInUsomm:                    uint64(500),
				MinimumSaleTokensUsdValue:            sdk.MustNewDecFromStr("1.0"),
				AuctionMaxBlockAge:                   uint64(100),
				AuctionPriceDecreaseAccelerationRate: sdk.MustNewDecFromStr("1.1"),
				AuctionBurnRate:                      sdk.MustNewDecFromStr("0.1"),
			},
			expPass: false,
			err:     errorsmod.Wrapf(ErrInvalidAuctionPriceDecreaseAccelerationRateParam, "auction price decrease acceleration rate must be between 0 and 1 inclusive (0%% to 100%%)"),
		},
		{
			name: "Auction burn rate lower bound (0)",
			params: Params{
				PriceMaxBlockAge:                     uint64(1000),
				MinimumBidInUsomm:                    uint64(500),
				MinimumSaleTokensUsdValue:            sdk.MustNewDecFromStr("1.0"),
				AuctionMaxBlockAge:                   uint64(100),
				AuctionPriceDecreaseAccelerationRate: sdk.MustNewDecFromStr("0.1"),
				AuctionBurnRate:                      sdk.ZeroDec(),
			},
			expPass: true,
			err:     nil,
		},
		{
			name: "Auction burn rate upper bound (1)",
			params: Params{
				PriceMaxBlockAge:                     uint64(1000),
				MinimumBidInUsomm:                    uint64(500),
				MinimumSaleTokensUsdValue:            sdk.MustNewDecFromStr("1.0"),
				AuctionMaxBlockAge:                   uint64(100),
				AuctionPriceDecreaseAccelerationRate: sdk.MustNewDecFromStr("0.1"),
				AuctionBurnRate:                      sdk.OneDec(),
			},
			expPass: true,
			err:     nil,
		},
		{
			name: "Auction burn rate slightly below lower bound",
			params: Params{
				PriceMaxBlockAge:                     uint64(1000),
				MinimumBidInUsomm:                    uint64(500),
				MinimumSaleTokensUsdValue:            sdk.MustNewDecFromStr("1.0"),
				AuctionMaxBlockAge:                   uint64(100),
				AuctionPriceDecreaseAccelerationRate: sdk.MustNewDecFromStr("0.1"),
				AuctionBurnRate:                      sdk.MustNewDecFromStr("-0.000000000000000001"),
			},
			expPass: false,
			err:     errorsmod.Wrapf(ErrInvalidAuctionBurnRateParam, "auction burn rate must be between 0 and 1 inclusive (0%% to 100%%)"),
		},
		{
			name: "Auction burn rate slightly above upper bound",
			params: Params{
				PriceMaxBlockAge:                     uint64(1000),
				MinimumBidInUsomm:                    uint64(500),
				MinimumSaleTokensUsdValue:            sdk.MustNewDecFromStr("1.0"),
				AuctionMaxBlockAge:                   uint64(100),
				AuctionPriceDecreaseAccelerationRate: sdk.MustNewDecFromStr("0.1"),
				AuctionBurnRate:                      sdk.MustNewDecFromStr("1.000000000000000001"),
			},
			expPass: false,
			err:     errorsmod.Wrapf(ErrInvalidAuctionBurnRateParam, "auction burn rate must be between 0 and 1 inclusive (0%% to 100%%)"),
		},
	}

	for _, tc := range testCases {
		err := tc.params.ValidateBasic()
		tc := tc
		if tc.expPass {
			require.NoError(t, err, tc.name)
			require.Nil(t, err)
		} else {
			require.Error(t, err, tc.name)
			require.Equal(t, tc.err.Error(), err.Error())
		}
	}
}

func TestParamsValidateBasicUnhappyPath(t *testing.T) {
	testCases := []struct {
		name   string
		params Params
		expErr error
	}{
		{
			name: "Invalid minimum sale tokens USD value",
			params: Params{
				PriceMaxBlockAge:                     uint64(1000),
				MinimumBidInUsomm:                    uint64(500),
				MinimumSaleTokensUsdValue:            sdk.MustNewDecFromStr("0.5"),
				AuctionMaxBlockAge:                   uint64(100),
				AuctionPriceDecreaseAccelerationRate: sdk.MustNewDecFromStr("0.1"),
				AuctionBurnRate:                      sdk.MustNewDecFromStr("0.1"),
			},
			expErr: ErrInvalidMinimumSaleTokensUSDValue,
		},
		{
			name: "Invalid auction burn rate (negative)",
			params: Params{
				PriceMaxBlockAge:                     uint64(1000),
				MinimumBidInUsomm:                    uint64(500),
				MinimumSaleTokensUsdValue:            sdk.MustNewDecFromStr("1.0"),
				AuctionMaxBlockAge:                   uint64(100),
				AuctionPriceDecreaseAccelerationRate: sdk.MustNewDecFromStr("0.1"),
				AuctionBurnRate:                      sdk.MustNewDecFromStr("-0.1"),
			},
			expErr: ErrInvalidAuctionBurnRateParam,
		},
		{
			name: "Invalid auction burn rate (greater than 1)",
			params: Params{
				PriceMaxBlockAge:                     uint64(1000),
				MinimumBidInUsomm:                    uint64(500),
				MinimumSaleTokensUsdValue:            sdk.MustNewDecFromStr("1.0"),
				AuctionMaxBlockAge:                   uint64(100),
				AuctionPriceDecreaseAccelerationRate: sdk.MustNewDecFromStr("0.1"),
				AuctionBurnRate:                      sdk.MustNewDecFromStr("1.1"),
			},
			expErr: ErrInvalidAuctionBurnRateParam,
		},
		{
			name: "Nil MinimumSaleTokensUsdValue",
			params: Params{
				PriceMaxBlockAge:                     uint64(1000),
				MinimumBidInUsomm:                    uint64(500),
				MinimumSaleTokensUsdValue:            sdk.Dec{},
				AuctionMaxBlockAge:                   uint64(100),
				AuctionPriceDecreaseAccelerationRate: sdk.MustNewDecFromStr("0.1"),
				AuctionBurnRate:                      sdk.MustNewDecFromStr("0.1"),
			},
			expErr: ErrInvalidMinimumSaleTokensUSDValue,
		},
		{
			name: "Nil AuctionPriceDecreaseAccelerationRate",
			params: Params{
				PriceMaxBlockAge:                     uint64(1000),
				MinimumBidInUsomm:                    uint64(500),
				MinimumSaleTokensUsdValue:            sdk.MustNewDecFromStr("1.0"),
				AuctionMaxBlockAge:                   uint64(100),
				AuctionPriceDecreaseAccelerationRate: sdk.Dec{},
				AuctionBurnRate:                      sdk.MustNewDecFromStr("0.1"),
			},
			expErr: ErrInvalidAuctionPriceDecreaseAccelerationRateParam,
		},
		{
			name: "Nil AuctionBurnRate",
			params: Params{
				PriceMaxBlockAge:                     uint64(1000),
				MinimumBidInUsomm:                    uint64(500),
				MinimumSaleTokensUsdValue:            sdk.MustNewDecFromStr("1.0"),
				AuctionMaxBlockAge:                   uint64(100),
				AuctionPriceDecreaseAccelerationRate: sdk.MustNewDecFromStr("0.1"),
				AuctionBurnRate:                      sdk.Dec{},
			},
			expErr: ErrInvalidAuctionBurnRateParam,
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			err := tc.params.ValidateBasic()
			require.Error(t, err)
			require.ErrorIs(t, err, tc.expErr)
		})
	}
}
