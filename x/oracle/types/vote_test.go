package types

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func TestUniswapDataUnmarshal(t *testing.T) {
	graphFeed := `{"pairs":[{"id":"0xb4e16d0168e52d35cacd2c6185b44281ec28c9dc","reserve0":"102942816.711844","reserve1":"65782.700042319083534616","reserveUSD":"205592160.9899610693250813158688189","token0":{"decimals":"6","id":"0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48"},"token0Price":"1564.89193428696614018979257090934","token1":{"decimals":"18","id":"0xc02aaa39b223fe8d0a0e5c4f27ead9083c756cc2"},"token1Price":"0.0006390217612410687980959120303082985","totalSupply":"1.831060361989902063"},{"id":"0x0d4a11d5eeaac28ec3f61d100daf4d40471f1852","reserve0":"48876.595478310639720902","reserve1":"76160219.812644","reserveUSD":"152739148.4135332779055124985116295","token0":{"decimals":"18","id":"0xc02aaa39b223fe8d0a0e5c4f27ead9083c756cc2"},"token0Price":"0.0006417601682157464592031298440374988","token1":{"decimals":"6","id":"0xdac17f958d2ee523a2206206994597c13d831ec7"},"token1Price":"1558.214500566854655310623371178052","totalSupply":"1.35475319056044026"},{"id":"0xa478c2975ab1ea89e8196811f51a7b7ade33eb11","reserve0":"76291257.87137489931615349","reserve1":"48769.877601615775380997","reserveUSD":"152427479.6216520380144593118342682","token0":{"decimals":"18","id":"0x6b175474e89094c44da98b954eedeac495271d0f"},"token0Price":"1564.311038353873663446706027856293","token1":{"decimals":"18","id":"0xc02aaa39b223fe8d0a0e5c4f27ead9083c756cc2"},"token1Price":"0.0006392590574904471494462260550608225","totalSupply":"1465832.764634346662785829"},{"id":"0xbb2b8038a1640196fbe3e38816f3e67cba72d940","reserve0":"3696.03693218","reserve1":"114558.251001995826169239","reserveUSD":"358005984.0245343506561444966654255","token0":{"decimals":"8","id":"0x2260fac5e5542a773aa44fbcfedf7c193bc2c599"},"token0Price":"0.03226338478330651095870213243843367","token1":{"decimals":"18","id":"0xc02aaa39b223fe8d0a0e5c4f27ead9083c756cc2"},"token1Price":"30.99488806634488097082059174219645","totalSupply":"0.180798275422738022"},{"id":"0xd3d2e2692501a5c9ca623199d38826e513033a17","reserve0":"5424546.656394144713446016","reserve1":"82873.000902416935145038","reserveUSD":"259339240.828553424748387936033602","token0":{"decimals":"18","id":"0x1f9840a85d5af5bf1d1762f925bdaddc4201f984"},"token0Price":"65.45613887907299699256692700832052","token1":{"decimals":"18","id":"0xc02aaa39b223fe8d0a0e5c4f27ead9083c756cc2"},"token1Price":"0.01527740586482577142170823317438549","totalSupply":"338230.831670328023456484"}]}`

	type wrapper struct {
		Pairs []*UniswapPair `json:"pairs"`
	}

	w := wrapper{}
	require.NoError(t, json.Unmarshal([]byte(graphFeed), &w))
	require.Equal(t, 5, len(w.Pairs))
}

func TestOracleVoteUnpacker(t *testing.T) {
	registry := codectypes.NewInterfaceRegistry()
	cdc := codec.NewProtoCodec(registry)

	pair := &UniswapPair{
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

	anyPair, err := PackOracleData(pair)
	require.NoError(t, err)

	vote := &OracleVote{
		Salt: []string{"salt"},
		Feed: &OracleFeed{
			Data: []*codectypes.Any{anyPair},
		},
	}

	bz, err := cdc.MarshalJSON(vote)
	require.NoError(t, err)
	require.Equal(t, "", string(bz))
}

func TestOracleVoteMarshalJSON(t *testing.T) {
	pair := &UniswapPair{
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

	pair2 := &UniswapPair{
		ID:         "0xa478c2975ab1ea89e8196811f51a7b7ade33eb11",
		Reserve0:   sdk.MustNewDecFromStr("69453224.061579510781012891"),
		Reserve1:   sdk.MustNewDecFromStr("45584.711379804929448746"),
		ReserveUSD: MustTruncateDec("138883455.9382328581978198800889056"),
		Token0: UniswapToken{
			Decimals: 18,
			ID:       "0x6b175474e89094c44da98b954eedeac495271d0f",
		},
		Token1: UniswapToken{
			Decimals: 18,
			ID:       "0xc02aaa39b223fe8d0a0e5c4f27ead9083c756cc2",
		},
		Token0Price: MustTruncateDec("1523.607849195440530703001829222234"),
		Token1Price: MustTruncateDec("0.0006563368655051639699748790273756903"),
		TotalSupply: MustTruncateDec("1387139.630260982742563912"),
	}

	anyPair, err := PackOracleData(pair)
	require.NoError(t, err)

	anyPair2, err := PackOracleData(pair2)
	require.NoError(t, err)

	vote := &OracleVote{
		Salt: []string{"salt", "salt2"},
		Feed: &OracleFeed{
			Data: []*codectypes.Any{anyPair, anyPair2},
		},
	}

	bz, err := json.Marshal(vote)
	require.NoError(t, err)
	require.Equal(t, "", string(bz))
}

func TestOracleVoteValidate(t *testing.T) {
	pair := &UniswapPair{
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

	anyPair, err := PackOracleData(pair)
	require.NoError(t, err)

	testCases := []struct {
		name    string
		vote    OracleVote
		expPass bool
	}{
		{
			"length missmatch",
			OracleVote{
				Salt: []string{"salt", "salt"},
				Feed: &OracleFeed{
					Data: []*codectypes.Any{nil},
				},
			},
			false,
		},
		{
			"empty salt",
			OracleVote{
				Salt: []string{" "},
				Feed: &OracleFeed{
					Data: []*codectypes.Any{nil},
				},
			},
			false,
		},
		{
			"empty vote",
			OracleVote{},
			false,
		},
		{
			"valid vote",
			OracleVote{
				Salt: []string{"salt"},
				Feed: &OracleFeed{
					Data: []*codectypes.Any{anyPair},
				},
			},
			true,
		},
	}

	for _, tc := range testCases {

		err := tc.vote.Validate()
		if tc.expPass {
			require.NoError(t, err, tc.name)
		} else {
			require.Error(t, err, tc.name)
		}
	}

}

func TestOracleFeedValidate(t *testing.T) {
	pair := &UniswapPair{
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

	pair2 := &UniswapPair{
		ID:         "0xa478c2975ab1ea89e8196811f51a7b7ade33eb11",
		Reserve0:   sdk.MustNewDecFromStr("69453224.061579510781012891"),
		Reserve1:   sdk.MustNewDecFromStr("45584.711379804929448746"),
		ReserveUSD: MustTruncateDec("138883455.9382328581978198800889056"),
		Token0: UniswapToken{
			Decimals: 18,
			ID:       "0x6b175474e89094c44da98b954eedeac495271d0f",
		},
		Token1: UniswapToken{
			Decimals: 18,
			ID:       "0xc02aaa39b223fe8d0a0e5c4f27ead9083c756cc2",
		},
		Token0Price: MustTruncateDec("1523.607849195440530703001829222234"),
		Token1Price: MustTruncateDec("0.0006563368655051639699748790273756903"),
		TotalSupply: MustTruncateDec("1387139.630260982742563912"),
	}

	invalidPair := &UniswapPair{}

	anyInvPair, err := PackOracleData(invalidPair)
	require.NoError(t, err)

	anyPair, err := PackOracleData(pair)
	require.NoError(t, err)

	anyPair2, err := PackOracleData(pair2)
	require.NoError(t, err)

	anyMockData, err := PackOracleData(&mockOracleData{})
	require.NoError(t, err)

	testCases := []struct {
		name    string
		feed    OracleFeed
		expPass bool
	}{
		{
			"single uniswap pair",
			OracleFeed{
				Data: []*codectypes.Any{anyPair},
			},
			true,
		},
		{
			"multiple uniswap pairs",
			OracleFeed{
				Data: []*codectypes.Any{anyPair, anyPair2},
			},
			true,
		},
		{
			"dup uniswap pair",
			OracleFeed{
				Data: []*codectypes.Any{anyPair, anyPair},
			},
			false,
		},
		{
			"different data types",
			OracleFeed{
				Data: []*codectypes.Any{anyPair, anyMockData},
			},
			false,
		},
		{
			"invalid oracle data",
			OracleFeed{
				Data: []*codectypes.Any{anyInvPair},
			},
			false,
		},
		{
			"empty oracle feed",
			OracleFeed{},
			false,
		},
	}

	for _, tc := range testCases {
		err := tc.feed.Validate()

		if tc.expPass {
			require.NoError(t, err, tc.name)
		} else {
			require.Error(t, err, tc.name)
		}
	}
}
