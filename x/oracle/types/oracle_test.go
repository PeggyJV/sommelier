package types

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

var _ OracleData = &mockOracleData{}

type mockOracleData struct{}

func (mockOracleData) GetID() string                        { return "mockOracleData" }
func (mockOracleData) Type() string                         { return "mock" }
func (mockOracleData) Validate() error                      { return nil }
func (mockOracleData) Compare(_ OracleData, _ sdk.Dec) bool { return false }
func (mockOracleData) Reset()                               {}
func (mockOracleData) String() string                       { return "mockOracleData" }
func (mockOracleData) ProtoMessage()                        {}

func TestUniswapPairValidate(t *testing.T) {
	testCases := []struct {
		name    string
		pair    UniswapPair
		expPass bool
	}{
		{
			"valid pair",
			UniswapPair{
				ID:         "0xb4e16d0168e52d35cacd2c6185b44281ec28c9dc",
				Reserve0:   sdk.MustNewDecFromStr("148681992.765143"),
				Reserve1:   sdk.MustNewDecFromStr("97709.503398661101176213"),
				ReserveUSD: sdk.MustNewDecFromStr("297632095.439861032964130850"),
				Token0: UniswapToken{
					Decimals: 6,
					ID:       "0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48",
				},
				Token1: UniswapToken{
					Decimals: 18,
					ID:       "0xc02aaa39b223fe8d0a0e5c4f27ead9083c756cc2",
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
				ID: "0xb4e16d0168e52d35cacd2c6185b44281ec28",
			},
			false,
		},
		{
			"nil reserve 0",
			UniswapPair{
				ID:       "0xb4e16d0168e52d35cacd2c6185b44281ec28c9dc",
				Reserve0: sdk.Dec{},
			},
			false,
		},
		{
			"nil reserve 1",
			UniswapPair{
				ID:       "0xb4e16d0168e52d35cacd2c6185b44281ec28c9dc",
				Reserve0: sdk.OneDec(),
				Reserve1: sdk.Dec{},
			},
			false,
		},
		{
			"nil reserve usd",
			UniswapPair{
				ID:         "0xb4e16d0168e52d35cacd2c6185b44281ec28c9dc",
				Reserve0:   sdk.OneDec(),
				Reserve1:   sdk.OneDec(),
				ReserveUSD: sdk.Dec{},
			},
			false,
		},
		{
			"nil token 0 price",
			UniswapPair{
				ID:          "0xb4e16d0168e52d35cacd2c6185b44281ec28c9dc",
				Reserve0:    sdk.OneDec(),
				Reserve1:    sdk.OneDec(),
				ReserveUSD:  sdk.OneDec(),
				Token0Price: sdk.Dec{},
			},
			false,
		},
		{
			"nil token 1 price",
			UniswapPair{
				ID:          "0xb4e16d0168e52d35cacd2c6185b44281ec28c9dc",
				Reserve0:    sdk.OneDec(),
				Reserve1:    sdk.OneDec(),
				ReserveUSD:  sdk.OneDec(),
				Token0Price: sdk.OneDec(),
				Token1Price: sdk.Dec{},
			},
			false,
		},
		{
			"nil total supply",
			UniswapPair{
				ID:          "0xb4e16d0168e52d35cacd2c6185b44281ec28c9dc",
				Reserve0:    sdk.OneDec(),
				Reserve1:    sdk.OneDec(),
				ReserveUSD:  sdk.OneDec(),
				Token0Price: sdk.OneDec(),
				Token1Price: sdk.OneDec(),
				TotalSupply: sdk.Dec{},
			},
			false,
		},
		{
			"neg reserve 0",
			UniswapPair{
				ID:          "0xb4e16d0168e52d35cacd2c6185b44281ec28c9dc",
				Reserve0:    sdk.NewDec(-1),
				Reserve1:    sdk.OneDec(),
				ReserveUSD:  sdk.OneDec(),
				Token0Price: sdk.OneDec(),
				Token1Price: sdk.OneDec(),
				TotalSupply: sdk.OneDec(),
			},
			false,
		},
		{
			"neg reserve 1",
			UniswapPair{
				ID:          "0xb4e16d0168e52d35cacd2c6185b44281ec28c9dc",
				Reserve0:    sdk.OneDec(),
				Reserve1:    sdk.NewDec(-1),
				ReserveUSD:  sdk.OneDec(),
				Token0Price: sdk.OneDec(),
				Token1Price: sdk.OneDec(),
				TotalSupply: sdk.OneDec(),
			},
			false,
		},
		{
			"neg reserve USD",
			UniswapPair{
				ID:          "0xb4e16d0168e52d35cacd2c6185b44281ec28c9dc",
				Reserve0:    sdk.OneDec(),
				Reserve1:    sdk.OneDec(),
				ReserveUSD:  sdk.NewDec(-1),
				Token0Price: sdk.OneDec(),
				Token1Price: sdk.OneDec(),
				TotalSupply: sdk.OneDec(),
			},
			false,
		},
		{
			"invalid token 0",
			UniswapPair{
				ID:         "0xb4e16d0168e52d35cacd2c6185b44281ec28c9dc",
				Reserve0:   sdk.OneDec(),
				Reserve1:   sdk.OneDec(),
				ReserveUSD: sdk.OneDec(),
				Token0: UniswapToken{
					Decimals: 6,
					ID:       "0x0",
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
				ID:         "0xb4e16d0168e52d35cacd2c6185b44281ec28c9dc",
				Reserve0:   sdk.OneDec(),
				Reserve1:   sdk.OneDec(),
				ReserveUSD: sdk.OneDec(),
				Token0: UniswapToken{
					Decimals: 6,
					ID:       "0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48",
				},
				Token1: UniswapToken{
					Decimals: 18,
					ID:       "0x0",
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
				ID:         "0xb4e16d0168e52d35cacd2c6185b44281ec28c9dc",
				Reserve0:   sdk.OneDec(),
				Reserve1:   sdk.OneDec(),
				ReserveUSD: sdk.OneDec(),
				Token0: UniswapToken{
					Decimals: 6,
					ID:       "0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48",
				},
				Token1: UniswapToken{
					Decimals: 18,
					ID:       "0xc02aaa39b223fe8d0a0e5c4f27ead9083c756cc2",
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
				ID:         "0xb4e16d0168e52d35cacd2c6185b44281ec28c9dc",
				Reserve0:   sdk.OneDec(),
				Reserve1:   sdk.OneDec(),
				ReserveUSD: sdk.OneDec(),
				Token0: UniswapToken{
					Decimals: 6,
					ID:       "0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48",
				},
				Token1: UniswapToken{
					Decimals: 18,
					ID:       "0xc02aaa39b223fe8d0a0e5c4f27ead9083c756cc2",
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
				ID:         "0xb4e16d0168e52d35cacd2c6185b44281ec28c9dc",
				Reserve0:   sdk.OneDec(),
				Reserve1:   sdk.OneDec(),
				ReserveUSD: sdk.OneDec(),
				Token0: UniswapToken{
					Decimals: 6,
					ID:       "0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48",
				},
				Token1: UniswapToken{
					Decimals: 18,
					ID:       "0xc02aaa39b223fe8d0a0e5c4f27ead9083c756cc2",
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
		ID:         "0xb4e16d0168e52d35cacd2c6185b44281ec28c9dc",
		Reserve0:   sdk.NewDec(100),
		Reserve1:   sdk.NewDec(100),
		ReserveUSD: sdk.NewDec(100),
		Token0: UniswapToken{
			Decimals: 6,
			ID:       "0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48",
		},
		Token1: UniswapToken{
			Decimals: 18,
			ID:       "0xc02aaa39b223fe8d0a0e5c4f27ead9083c756cc2",
		},
		Token0Price: sdk.NewDec(100),
		Token1Price: sdk.NewDec(100),
		TotalSupply: sdk.NewDec(100),
	}

	target := sdk.NewDecWithPrec(5, 2) // 0.05
	aggregatePair := NewUniswapPair(pair.ID, pair.Token0, pair.Token1)

	testCases := []struct {
		name           string
		pair           *UniswapPair
		malleate       func()
		isWithinTarget bool
	}{
		{
			"different ID",
			&UniswapPair{
				ID: "0x0",
			},
			func() {},
			false,
		},
		{
			"different token 0",
			&UniswapPair{
				ID: "0xb4e16d0168e52d35cacd2c6185b44281ec28c9dc",
				Token0: UniswapToken{
					Decimals: 0,
					ID:       "0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48",
				},
				Token1: UniswapToken{
					Decimals: 18,
					ID:       "0xc02aaa39b223fe8d0a0e5c4f27ead9083c756cc2",
				},
			},
			func() {},
			false,
		},
		{
			"different token 1",
			&UniswapPair{
				ID: "0xb4e16d0168e52d35cacd2c6185b44281ec28c9dc",
				Token0: UniswapToken{
					Decimals: 6,
					ID:       "0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48",
				},
				Token1: UniswapToken{
					Decimals: 18,
					ID:       "",
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
				aggregatePair.ReserveUSD = pair.ReserveUSD.Add(sdk.NewDec(4))
			},
			false,
		},
		{
			"token 1 price not within target",
			pair,
			func() {
				aggregatePair.Reserve0 = pair.Reserve0.Add(sdk.NewDec(4))
				aggregatePair.Reserve1 = pair.Reserve1.Add(sdk.NewDec(4))
				aggregatePair.ReserveUSD = pair.ReserveUSD.Add(sdk.NewDec(4))
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
				aggregatePair.ReserveUSD = pair.ReserveUSD.Add(sdk.NewDec(4))
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
				aggregatePair.ReserveUSD = pair.ReserveUSD.Add(sdk.NewDec(4))
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

	require.False(t, pair.Compare(&mockOracleData{}, target))
}

func TestPairUnmarshalJSON(t *testing.T) {
	pairJSON := `{
 "id": "0xb4e16d0168e52d35cacd2c6185b44281ec28c9dc",
 "reserve0": "148681992.765143000000000000",
 "reserve1": "97709.503398661101176213",
 "reserveUSD": "297632095.439861032964130850561223123",
 "token0": {
   "id": "0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48",
	 "decimals": "6"
 },
 "token1": {
	 "id": "0xc02aaa39b223fe8d0a0e5c4f27ead9083c756cc2",
	 "decimals": "18"
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
  "decimals": "6"
 },
 "token1": {
  "id": "0xc02aaa39b223fe8d0a0e5c4f27ead9083c756cc2",
  "decimals": "18"
 },
 "token0Price": "1521.673814659673802831",
 "token1Price": "0.000657171064104597",
 "totalSupply": "2.754869216896965436"
}`

	pairTrunc := UniswapPair{
		ID:         "0xb4e16d0168e52d35cacd2c6185b44281ec28c9dc",
		Reserve0:   sdk.MustNewDecFromStr("148681992.765143"),
		Reserve1:   sdk.MustNewDecFromStr("97709.503398661101176213"),
		ReserveUSD: sdk.MustNewDecFromStr("297632095.439861032964130850"),
		Token0: UniswapToken{
			Decimals: 6,
			ID:       "0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48",
		},
		Token1: UniswapToken{
			Decimals: 18,
			ID:       "0xc02aaa39b223fe8d0a0e5c4f27ead9083c756cc2",
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
	_, err := TruncateDec("e")
	require.Error(t, err)

	dec, err := TruncateDec("1")
	require.NoError(t, err)
	require.Equal(t, sdk.OneDec().String(), dec.String())

	dec, err = TruncateDec("1.0")
	require.NoError(t, err)
	require.Equal(t, sdk.OneDec().String(), dec.String())

	dec, err = TruncateDec("195116448.3284569661435357469623931")
	require.NoError(t, err)
	require.Equal(t, "195116448.328456966143535746", dec.String())
}
