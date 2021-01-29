package types

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/cosmos/cosmos-sdk/testutil/testdata"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func TestLPsStoplossPositionsValidate(t *testing.T) {
	_, _, addr1 := testdata.KeyTestPubAddr()
	_, _, addr2 := testdata.KeyTestPubAddr()

	testCases := []struct {
		msg               string
		stoplossPositions LPsStoplossPositions
		expPass           bool
	}{
		{
			"valid single lps position",
			LPsStoplossPositions{
				{
					Address: addr1.String(),
					StoplossPositions: []Stoploss{
						{
							UniswapPairId:       "0x3041cbd36888becc7bbcbc0045e3b1f144466f5f",
							LiquidityPoolShares: 10,
							MaxSlippage:         sdk.MustNewDecFromStr("0.05"),
							ReferencePairRatio:  sdk.MustNewDecFromStr("0.1"),
						},
					},
				},
			},
			true,
		},
		{
			"valid multiple lps positions",
			LPsStoplossPositions{
				{
					Address: addr1.String(),
					StoplossPositions: []Stoploss{
						{
							UniswapPairId:       "0x3041cbd36888becc7bbcbc0045e3b1f144466f5f",
							LiquidityPoolShares: 10,
							MaxSlippage:         sdk.MustNewDecFromStr("0.05"),
							ReferencePairRatio:  sdk.MustNewDecFromStr("0.1"),
						}, {
							UniswapPairId:       "0x66e33d2605c5fb25ebb7cd7528e7997b0afa55e8",
							LiquidityPoolShares: 1,
							MaxSlippage:         sdk.MustNewDecFromStr("0.03333"),
							ReferencePairRatio:  sdk.MustNewDecFromStr("0.3"),
						},
					},
				},
				{
					Address: addr2.String(),
					StoplossPositions: []Stoploss{
						{
							UniswapPairId:       "0x3041cbd36888becc7bbcbc0045e3b1f144466f5f",
							LiquidityPoolShares: 1000,
							MaxSlippage:         sdk.MustNewDecFromStr("0.01"),
							ReferencePairRatio:  sdk.MustNewDecFromStr("0.5"),
						},
					},
				},
			},
			true,
		},
		{
			"empty position",
			LPsStoplossPositions{{}},
			false,
		},
		{
			"dup positions for same LP address",
			LPsStoplossPositions{
				{
					Address: addr1.String(),
					StoplossPositions: []Stoploss{
						{
							UniswapPairId:       "0x3041cbd36888becc7bbcbc0045e3b1f144466f5f",
							LiquidityPoolShares: 10,
							MaxSlippage:         sdk.MustNewDecFromStr("0.05"),
							ReferencePairRatio:  sdk.MustNewDecFromStr("0.1"),
						},
					},
				},
				{
					Address: addr1.String(),
					StoplossPositions: []Stoploss{
						{
							UniswapPairId:       "0x3041cbd36888becc7bbcbc0045e3b1f144466f5f",
							LiquidityPoolShares: 10,
							MaxSlippage:         sdk.MustNewDecFromStr("0.05"),
							ReferencePairRatio:  sdk.MustNewDecFromStr("0.1"),
						},
					},
				},
			},
			false,
		},
		{
			"dup stoploss for same uniswap pair",
			LPsStoplossPositions{
				{
					Address: addr1.String(),
					StoplossPositions: []Stoploss{
						{
							UniswapPairId:       "0x3041cbd36888becc7bbcbc0045e3b1f144466f5f",
							LiquidityPoolShares: 10,
							MaxSlippage:         sdk.MustNewDecFromStr("0.05"),
							ReferencePairRatio:  sdk.MustNewDecFromStr("0.1"),
						},
						{
							UniswapPairId:       "0x3041cbd36888becc7bbcbc0045e3b1f144466f5f",
							LiquidityPoolShares: 10,
							MaxSlippage:         sdk.MustNewDecFromStr("0.05"),
							ReferencePairRatio:  sdk.MustNewDecFromStr("0.1"),
						},
					},
				},
			},
			false,
		},
		{
			"invalid position",
			LPsStoplossPositions{
				{
					Address:           addr1.String(),
					StoplossPositions: []Stoploss{{}},
				},
			},
			false,
		},
	}

	for _, tc := range testCases {

		err := tc.stoplossPositions.Validate()

		if tc.expPass {
			require.NoError(t, err)
		} else {
			require.Error(t, err)
		}
	}
}

func TestStoplossValidate(t *testing.T) {
	testCases := []struct {
		msg      string
		stoploss Stoploss
		expPass  bool
	}{
		{
			"valid stoploss position",
			Stoploss{
				UniswapPairId:       "0x3041cbd36888becc7bbcbc0045e3b1f144466f5f",
				LiquidityPoolShares: 10,
				MaxSlippage:         sdk.MustNewDecFromStr("0.05"),
				ReferencePairRatio:  sdk.MustNewDecFromStr("0.1"),
			},
			true,
		},
		{
			"invalid pair address",
			Stoploss{
				UniswapPairId: "",
			},
			false,
		},
		{
			"invalid max slippage",
			Stoploss{
				UniswapPairId: "0x3041cbd36888becc7bbcbc0045e3b1f144466f5f",
				MaxSlippage:   sdk.ZeroDec(),
			},
			false,
		},
		{
			"invalid shares",
			Stoploss{
				UniswapPairId:       "0x3041cbd36888becc7bbcbc0045e3b1f144466f5f",
				MaxSlippage:         sdk.OneDec(),
				LiquidityPoolShares: 0,
			},
			false,
		},
		{
			"invalid reference pair ratio",
			Stoploss{
				UniswapPairId:       "0x3041cbd36888becc7bbcbc0045e3b1f144466f5f",
				MaxSlippage:         sdk.OneDec(),
				LiquidityPoolShares: 1,
				ReferencePairRatio:  sdk.MustNewDecFromStr("1.1"),
			},
			false,
		},
	}

	for _, tc := range testCases {

		err := tc.stoploss.Validate()

		if tc.expPass {
			require.NoError(t, err)
		} else {
			require.Error(t, err)
		}
	}
}
