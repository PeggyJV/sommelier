package types

import (
	"encoding/json"
	"testing"

	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
)

type caseAny struct {
	name    string
	any     *codectypes.Any
	expPass bool
}

func TestOracleDataParser(t *testing.T) {
	ud := []byte(`{"pairs":[{"id":"0xb4e16d0168e52d35cacd2c6185b44281ec28c9dc","reserve0":"104461984.297382","reserve1":"68900.854799992224259823","reserveUSD":"208943132.4491626262213595826180009","token0":{"decimals":"6","id":"0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48"},"token0Price":"1516.120294887746283192558065861314","token1":{"decimals":"18","id":"0xc02aaa39b223fe8d0a0e5c4f27ead9083c756cc2"},"token1Price":"0.0006595782691993052650502521727791989","totalSupply":"1.880242946708099187"},{"id":"0x0d4a11d5eeaac28ec3f61d100daf4d40471f1852","reserve0":"53692.177663232830871307","reserve1":"81435706.765445","reserveUSD":"162849896.1247941563471620169475477","token0":{"decimals":"18","id":"0xc02aaa39b223fe8d0a0e5c4f27ead9083c756cc2"},"token0Price":"0.0006593198462423810897434548766768373","token1":{"decimals":"6","id":"0xdac17f958d2ee523a2206206994597c13d831ec7"},"token1Price":"1516.714544085446922789613770363501","totalSupply":"1.462627249647822595"},{"id":"0xa478c2975ab1ea89e8196811f51a7b7ade33eb11","reserve0":"79621128.460135101988015534","reserve1":"52491.925771368416639614","reserveUSD":"159207795.8677692888272481356485488","token0":{"decimals":"18","id":"0x6b175474e89094c44da98b954eedeac495271d0f"},"token0Price":"1516.826203079869435171582129025643","token1":{"decimals":"18","id":"0xc02aaa39b223fe8d0a0e5c4f27ead9083c756cc2"},"token1Price":"0.000659271311353621422597060710564277","totalSupply":"1547052.230273708855937977"},{"id":"0xbb2b8038a1640196fbe3e38816f3e67cba72d940","reserve0":"3846.53655471","reserve1":"121466.398299929551140843","reserveUSD":"368411027.5777067517007445214099817","token0":{"decimals":"8","id":"0x2260fac5e5542a773aa44fbcfedf7c193bc2c599"},"token0Price":"0.03166749494960723585472081861318512","token1":{"decimals":"18","id":"0xc02aaa39b223fe8d0a0e5c4f27ead9083c756cc2"},"token1Price":"31.57812140149730264224077230074571","totalSupply":"0.189820649890157008"},{"id":"0xd3d2e2692501a5c9ca623199d38826e513033a17","reserve0":"5319718.233921056627144513","reserve1":"84915.982714119045010681","reserveUSD":"257291680.0035612150151144934818004","token0":{"decimals":"18","id":"0x1f9840a85d5af5bf1d1762f925bdaddc4201f984"},"token0Price":"62.64684296041884034640389605526021","token1":{"decimals":"18","id":"0xc02aaa39b223fe8d0a0e5c4f27ead9083c756cc2"},"token1Price":"0.01596249631656321636472293427813239","totalSupply":"338797.452045304481602297"}]}`)

	type parseUD struct {
		Pairs []*UniswapPair `json:"pairs"`
	}

	out := parseUD{}

	require.NoError(t, json.Unmarshal(ud, &out))
	require.Equal(t, 5, len(out.Pairs))
}

func TestOracleDataPacker(t *testing.T) {
	testCases := []struct {
		name       string
		oracleData OracleData
		expPass    bool
	}{
		{
			"uniswap pair",
			&UniswapPair{
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
			"nil",
			nil,
			false,
		},
	}

	testCasesAny := []caseAny{}

	for _, tc := range testCases {
		oracleDataAny, err := PackOracleData(tc.oracleData)
		if tc.expPass {
			require.NoError(t, err, tc.name)
		} else {
			require.Error(t, err, tc.name)
		}

		testCasesAny = append(testCasesAny, caseAny{tc.name, oracleDataAny, tc.expPass})
	}

	for i, tc := range testCasesAny {
		cs, err := UnpackOracleData(tc.any)
		if tc.expPass {
			require.NoError(t, err, tc.name)
			require.Equal(t, testCases[i].oracleData, cs, tc.name)
		} else {
			require.Error(t, err, tc.name)
		}
	}
}
