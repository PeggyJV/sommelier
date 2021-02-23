package types

import (
	"testing"

	"github.com/stretchr/testify/require"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func TestUniswapPairValidate(t *testing.T) {
	testCases := []struct {
		name    string
		pair    UniswapPair
		expPass bool
	}{
		// TODO: chop precision from string
		{
			"valid pair",
			UniswapPair{
				Id:         "0xb4e16d0168e52d35cacd2c6185b44281ec28c9dc",
				Reserve0:   sdk.MustNewDecFromStr("148681992.765143"),
				Reserve1:   sdk.MustNewDecFromStr("97709.503398661101176213"),
				ReserveUsd: sdk.MustNewDecFromStr("297632095.439861032964130850"),
				Token0: UniswapToken{
					Decimals: 6,
					Id:       "0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48",
				},
				Token1: UniswapToken{
					Decimals: 18,
					Id:       "0xc02aaa39b223fe8d0a0e5c4f27ead9083c756cc2",
				},
				Token0Price: sdk.MustNewDecFromStr("1521.673814659673802831"),
				Token1Price: sdk.MustNewDecFromStr("0.000657171064104597"),
				TotalSupply: sdk.MustNewDecFromStr("2.754869216896965436"),
			},
			true,
		},
		{
			"empty", UniswapPair{}, false,
		},
		{
			"zero fields",
			*NewUniswapPair("", UniswapToken{"0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48", 6}, UniswapToken{"0xc02aaa39b223fe8d0a0e5c4f27ead9083c756cc2", 18}),
			false,
		},
		{
			"invalid pair id",
			UniswapPair{
				Id: "0xb4e16d0168e52d35cacd2c6185b44281ec28",
			},
			false,
		},
		{
			"nil reserve 0",
			UniswapPair{
				Id:       "0xb4e16d0168e52d35cacd2c6185b44281ec28c9dc",
				Reserve0: sdk.Dec{},
			},
			false,
		},
		{
			"nil reserve 1",
			UniswapPair{
				Id:       "0xb4e16d0168e52d35cacd2c6185b44281ec28c9dc",
				Reserve0: sdk.OneDec(),
				Reserve1: sdk.Dec{},
			},
			false,
		},
		{
			"nil reserve usd",
			UniswapPair{
				Id:         "0xb4e16d0168e52d35cacd2c6185b44281ec28c9dc",
				Reserve0:   sdk.OneDec(),
				Reserve1:   sdk.OneDec(),
				ReserveUsd: sdk.Dec{},
			},
			false,
		},
		{
			"nil token 0 price",
			UniswapPair{
				Id:          "0xb4e16d0168e52d35cacd2c6185b44281ec28c9dc",
				Reserve0:    sdk.OneDec(),
				Reserve1:    sdk.OneDec(),
				ReserveUsd:  sdk.OneDec(),
				Token0Price: sdk.Dec{},
			},
			false,
		},
		{
			"nil token 1 price",
			UniswapPair{
				Id:          "0xb4e16d0168e52d35cacd2c6185b44281ec28c9dc",
				Reserve0:    sdk.OneDec(),
				Reserve1:    sdk.OneDec(),
				ReserveUsd:  sdk.OneDec(),
				Token0Price: sdk.OneDec(),
				Token1Price: sdk.Dec{},
			},
			false,
		},
		{
			"nil total supply",
			UniswapPair{
				Id:          "0xb4e16d0168e52d35cacd2c6185b44281ec28c9dc",
				Reserve0:    sdk.OneDec(),
				Reserve1:    sdk.OneDec(),
				ReserveUsd:  sdk.OneDec(),
				Token0Price: sdk.OneDec(),
				Token1Price: sdk.OneDec(),
				TotalSupply: sdk.Dec{},
			},
			false,
		},
		{
			"neg reserve 0",
			UniswapPair{
				Id:          "0xb4e16d0168e52d35cacd2c6185b44281ec28c9dc",
				Reserve0:    sdk.NewDec(-1),
				Reserve1:    sdk.OneDec(),
				ReserveUsd:  sdk.OneDec(),
				Token0Price: sdk.OneDec(),
				Token1Price: sdk.OneDec(),
				TotalSupply: sdk.OneDec(),
			},
			false,
		},
		{
			"neg reserve 1",
			UniswapPair{
				Id:          "0xb4e16d0168e52d35cacd2c6185b44281ec28c9dc",
				Reserve0:    sdk.OneDec(),
				Reserve1:    sdk.NewDec(-1),
				ReserveUsd:  sdk.OneDec(),
				Token0Price: sdk.OneDec(),
				Token1Price: sdk.OneDec(),
				TotalSupply: sdk.OneDec(),
			},
			false,
		},
		{
			"neg reserve USD",
			UniswapPair{
				Id:          "0xb4e16d0168e52d35cacd2c6185b44281ec28c9dc",
				Reserve0:    sdk.OneDec(),
				Reserve1:    sdk.OneDec(),
				ReserveUsd:  sdk.NewDec(-1),
				Token0Price: sdk.OneDec(),
				Token1Price: sdk.OneDec(),
				TotalSupply: sdk.OneDec(),
			},
			false,
		},
		{
			"invalid token 0",
			UniswapPair{
				Id:         "0xb4e16d0168e52d35cacd2c6185b44281ec28c9dc",
				Reserve0:   sdk.OneDec(),
				Reserve1:   sdk.OneDec(),
				ReserveUsd: sdk.OneDec(),
				Token0: UniswapToken{
					Decimals: 6,
					Id:       "0x0",
				},
				Token0Price: sdk.OneDec(),
				Token1Price: sdk.OneDec(),
				TotalSupply: sdk.OneDec(),
			},
			false,
		},
		{
			"invalid token 0",
			UniswapPair{
				Id:         "0xb4e16d0168e52d35cacd2c6185b44281ec28c9dc",
				Reserve0:   sdk.OneDec(),
				Reserve1:   sdk.OneDec(),
				ReserveUsd: sdk.OneDec(),
				Token0: UniswapToken{
					Decimals: 6,
					Id:       "0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48",
				},
				Token1: UniswapToken{
					Decimals: 18,
					Id:       "0x0",
				},
				Token0Price: sdk.OneDec(),
				Token1Price: sdk.OneDec(),
				TotalSupply: sdk.OneDec(),
			},
			false,
		},
		{
			"neg token 0 price",
			UniswapPair{
				Id:         "0xb4e16d0168e52d35cacd2c6185b44281ec28c9dc",
				Reserve0:   sdk.OneDec(),
				Reserve1:   sdk.OneDec(),
				ReserveUsd: sdk.OneDec(),
				Token0: UniswapToken{
					Decimals: 6,
					Id:       "0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48",
				},
				Token1: UniswapToken{
					Decimals: 18,
					Id:       "0xc02aaa39b223fe8d0a0e5c4f27ead9083c756cc2",
				},
				Token0Price: sdk.NewDec(-1),
				Token1Price: sdk.OneDec(),
				TotalSupply: sdk.OneDec(),
			},
			false,
		},
		{
			"neg token 1 price",
			UniswapPair{
				Id:         "0xb4e16d0168e52d35cacd2c6185b44281ec28c9dc",
				Reserve0:   sdk.OneDec(),
				Reserve1:   sdk.OneDec(),
				ReserveUsd: sdk.OneDec(),
				Token0: UniswapToken{
					Decimals: 6,
					Id:       "0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48",
				},
				Token1: UniswapToken{
					Decimals: 18,
					Id:       "0xc02aaa39b223fe8d0a0e5c4f27ead9083c756cc2",
				},
				Token0Price: sdk.OneDec(),
				Token1Price: sdk.NewDec(-1),
				TotalSupply: sdk.OneDec(),
			},
			false,
		},
		{
			"neg total supply",
			UniswapPair{
				Id:         "0xb4e16d0168e52d35cacd2c6185b44281ec28c9dc",
				Reserve0:   sdk.OneDec(),
				Reserve1:   sdk.OneDec(),
				ReserveUsd: sdk.OneDec(),
				Token0: UniswapToken{
					Decimals: 6,
					Id:       "0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48",
				},
				Token1: UniswapToken{
					Decimals: 18,
					Id:       "0xc02aaa39b223fe8d0a0e5c4f27ead9083c756cc2",
				},
				Token0Price: sdk.OneDec(),
				Token1Price: sdk.OneDec(),
				TotalSupply: sdk.NewDec(-1),
			},
			false,
		},
	}

	for _, tc := range testCases {

		err := tc.pair.Validate()
		if tc.expPass {
			require.NoError(t, err, tc.name)
		} else {
			require.Error(t, err, tc.name)
		}

	}
}

func TestUniswapPairCompare(t *testing.T) {
	pair := &UniswapPair{
		Id:         "0xb4e16d0168e52d35cacd2c6185b44281ec28c9dc",
		Reserve0:   sdk.MustNewDecFromStr("148681992.765143"),
		Reserve1:   sdk.MustNewDecFromStr("97709.503398661101176213"),
		ReserveUsd: sdk.MustNewDecFromStr("297632095.439861032964130850"),
		Token0: UniswapToken{
			Decimals: 6,
			Id:       "0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48",
		},
		Token1: UniswapToken{
			Decimals: 18,
			Id:       "0xc02aaa39b223fe8d0a0e5c4f27ead9083c756cc2",
		},
		Token0Price: sdk.MustNewDecFromStr("1521.673814659673802831"),
		Token1Price: sdk.MustNewDecFromStr("0.000657171064104597"),
		TotalSupply: sdk.MustNewDecFromStr("2.754869216896965436"),
	}

	target := sdk.NewDecWithPrec(5, 3) // 0.05
	aggregatePair := NewUniswapPair(pair.Id, pair.Token0, pair.Token1)

	testCases := []struct {
		name           string
		pair           *UniswapPair
		isWithinTarget bool
	}{
		{
			"default pair",
			pair,
			false,
		},
		{
			"different ID",
			&UniswapPair{
				Id: "0x0",
			},
			false,
		},
		{
			"different token 0",
			&UniswapPair{
				Id: "0xb4e16d0168e52d35cacd2c6185b44281ec28c9dc",
				Token0: UniswapToken{
					Decimals: 0,
					Id:       "0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48",
				},
				Token1: UniswapToken{
					Decimals: 18,
					Id:       "0xc02aaa39b223fe8d0a0e5c4f27ead9083c756cc2",
				},
			},
			false,
		},
		{
			"different token 1",
			&UniswapPair{
				Id: "0xb4e16d0168e52d35cacd2c6185b44281ec28c9dc",
				Token0: UniswapToken{
					Decimals: 6,
					Id:       "0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48",
				},
				Token1: UniswapToken{
					Decimals: 18,
					Id:       "",
				},
			},
			false,
		},
	}

	for _, tc := range testCases {
		require.Equal(t, tc.isWithinTarget, tc.pair.Compare(aggregatePair, target), tc.name)
	}
}
