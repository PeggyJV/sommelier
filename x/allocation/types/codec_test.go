package types

import (
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
