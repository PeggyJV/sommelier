package types

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
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
				AuctionMaxBlockAge:                   uint64(100),
				AuctionPriceDecreaseAccelerationRate: sdk.MustNewDecFromStr("0.1"),
			},
			expPass: true,
			err:     nil,
		},
		{
			name: "Token price max block age cannot be 0",
			params: Params{
				PriceMaxBlockAge:                     uint64(0),
				MinimumBidInUsomm:                    uint64(500),
				AuctionMaxBlockAge:                   uint64(100),
				AuctionPriceDecreaseAccelerationRate: sdk.MustNewDecFromStr("0.1"),
			},
			expPass: false,
			err:     sdkerrors.Wrapf(ErrTokenPriceMaxBlockAgeMustBePositive, "value: 0"),
		},
		{
			name: "Auction price decrease acceleration rate bounds check lower end",
			params: Params{
				PriceMaxBlockAge:                     uint64(1000),
				MinimumBidInUsomm:                    uint64(500),
				AuctionMaxBlockAge:                   uint64(100),
				AuctionPriceDecreaseAccelerationRate: sdk.MustNewDecFromStr("-0.01"),
			},
			expPass: false,
			err:     sdkerrors.Wrapf(ErrInvalidAuctionPriceDecreaseAccelerationRateParam, "auction price decrease acceleration rate must be between 0 and 1 inclusive (0%% to 100%%)"),
		},
		{
			name: "Auction price decrease acceleration rate bounds check upper end",
			params: Params{
				PriceMaxBlockAge:                     uint64(1000),
				MinimumBidInUsomm:                    uint64(500),
				AuctionMaxBlockAge:                   uint64(100),
				AuctionPriceDecreaseAccelerationRate: sdk.MustNewDecFromStr("1.1"),
			},
			expPass: false,
			err:     sdkerrors.Wrapf(ErrInvalidAuctionPriceDecreaseAccelerationRateParam, "auction price decrease acceleration rate must be between 0 and 1 inclusive (0%% to 100%%)"),
		},
	}

	for _, tc := range testCases {
		err := tc.params.ValidateBasic()
		if tc.expPass {
			require.NoError(t, err, tc.name)
			require.Nil(t, err)
		} else {
			require.Error(t, err, tc.name)
			require.Equal(t, tc.err.Error(), err.Error())
		}
	}
}
