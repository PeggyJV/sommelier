package types

import (
	"encoding/json"
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
		Reserve0:   sdk.NewDec(100),
		Reserve1:   sdk.NewDec(100),
		ReserveUsd: sdk.NewDec(100),
		Token0: UniswapToken{
			Decimals: 6,
			Id:       "0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48",
		},
		Token1: UniswapToken{
			Decimals: 18,
			Id:       "0xc02aaa39b223fe8d0a0e5c4f27ead9083c756cc2",
		},
		Token0Price: sdk.NewDec(100),
		Token1Price: sdk.NewDec(100),
		TotalSupply: sdk.NewDec(100),
	}

	target := sdk.NewDecWithPrec(5, 2) // 0.05
	aggregatePair := NewUniswapPair(pair.Id, pair.Token0, pair.Token1)

	testCases := []struct {
		name           string
		pair           *UniswapPair
		malleate       func()
		isWithinTarget bool
	}{
		{
			"different ID",
			&UniswapPair{
				Id: "0x0",
			},
			func() {},
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
			func() {},
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
			func() {},
			false,
		},
		{
			"zero fields aggregated data",
			pair,
			func() {},
			false,
		},
		{
			"reserve 1 not within target",
			pair,
			func() {
				aggregatePair.Reserve0 = pair.Reserve0.Add(sdk.NewDec(4))
			},
			false,
		},
		{
			"reserve USD not within target",
			pair,
			func() {
				aggregatePair.Reserve0 = pair.Reserve0.Add(sdk.NewDec(4))
				aggregatePair.Reserve1 = pair.Reserve1.Add(sdk.NewDec(4))
			},
			false,
		},
		{
			"token 0 price not within target",
			pair,
			func() {
				aggregatePair.Reserve0 = pair.Reserve0.Add(sdk.NewDec(4))
				aggregatePair.Reserve1 = pair.Reserve1.Add(sdk.NewDec(4))
				aggregatePair.ReserveUsd = pair.ReserveUsd.Add(sdk.NewDec(4))
			},
			false,
		},
		{
			"token 1 price not within target",
			pair,
			func() {
				aggregatePair.Reserve0 = pair.Reserve0.Add(sdk.NewDec(4))
				aggregatePair.Reserve1 = pair.Reserve1.Add(sdk.NewDec(4))
				aggregatePair.ReserveUsd = pair.ReserveUsd.Add(sdk.NewDec(4))
				aggregatePair.Token0Price = pair.Token0Price.Add(sdk.NewDec(4))
			},
			false,
		},
		{
			"total supply not within target",
			pair,
			func() {
				aggregatePair.Reserve0 = pair.Reserve0.Add(sdk.NewDec(4))
				aggregatePair.Reserve1 = pair.Reserve1.Add(sdk.NewDec(4))
				aggregatePair.ReserveUsd = pair.ReserveUsd.Add(sdk.NewDec(4))
				aggregatePair.Token0Price = pair.Token0Price.Add(sdk.NewDec(4))
				aggregatePair.Token1Price = pair.Token1Price.Add(sdk.NewDec(4))
			},
			false,
		},
		{
			"valid pair within target",
			pair,
			func() {
				aggregatePair.Reserve0 = pair.Reserve0.Add(sdk.NewDec(4))
				aggregatePair.Reserve1 = pair.Reserve1.Add(sdk.NewDec(4))
				aggregatePair.ReserveUsd = pair.ReserveUsd.Add(sdk.NewDec(4))
				aggregatePair.Token0Price = pair.Token0Price.Add(sdk.NewDec(4))
				aggregatePair.Token1Price = pair.Token1Price.Add(sdk.NewDec(4))
				aggregatePair.TotalSupply = pair.TotalSupply.Add(sdk.NewDec(4))
			},
			true,
		},
	}

	for _, tc := range testCases {
		tc.malleate()

		require.Equal(t, tc.isWithinTarget, tc.pair.Compare(aggregatePair, target), tc.name)
	}
}

func TestPairUnmarshalJSON(t *testing.T) {
	pairJSON := `{
 "id": "0xb4e16d0168e52d35cacd2c6185b44281ec28c9dc",
 "reserve0": "148681992.765143000000000000",
 "reserve1": "97709.503398661101176213",
 "reserveUSD": "297632095.439861032964130850561223123",
 "token0": {
   "id": "0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48",
	 "decimals": 6
 },
 "token1": {
	 "id": "0xc02aaa39b223fe8d0a0e5c4f27ead9083c756cc2",
	 "decimals": 18
 },
 "token0Price": "1521.67381465967380283112313",
 "token1Price": "0.000657171064104597123123",
 "totalSupply": "2.754869216896965436123123123"
}`

	pairJSONTrunc := `{
 "id": "0xb4e16d0168e52d35cacd2c6185b44281ec28c9dc",
 "reserve0": "148681992.765143000000000000",
 "reserve1": "97709.503398661101176213",
 "reserveUSD": "297632095.439861032964130850",
 "token0": {
  "id": "0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48",
  "decimals": 6
 },
 "token1": {
  "id": "0xc02aaa39b223fe8d0a0e5c4f27ead9083c756cc2",
  "decimals": 18
 },
 "token0Price": "1521.673814659673802831",
 "token1Price": "0.000657171064104597",
 "totalSupply": "2.754869216896965436"
}`

	pairTrunc := UniswapPair{
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

	var pair UniswapPair
	err := json.Unmarshal([]byte(pairJSON), &pair)
	require.NoError(t, err)

	require.Equal(t, pairTrunc, pair)

	bz, err := json.MarshalIndent(pairTrunc, "", " ")
	require.NoError(t, err)

	require.Equal(t, pairJSONTrunc, string(bz))
}

func TestTruncateDec(t *testing.T) {
	_, err := truncateDec("1")
	require.Error(t, err)

	dec, err := truncateDec("1.0")
	require.NoError(t, err)
	require.Equal(t, sdk.OneDec().String(), dec.String())

	dec, err = truncateDec("195116448.3284569661435357469623931")
	require.NoError(t, err)
	require.Equal(t, "195116448.328456966143535746", dec.String())
}
