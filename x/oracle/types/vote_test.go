package types

import (
	"testing"

	"github.com/stretchr/testify/require"

	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func TestOracleFeedValidate(t *testing.T) {
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

	pair2 := &UniswapPair{
		Id:         "0xa478c2975ab1ea89e8196811f51a7b7ade33eb11",
		Reserve0:   sdk.MustNewDecFromStr("69453224.061579510781012891"),
		Reserve1:   sdk.MustNewDecFromStr("45584.711379804929448746"),
		ReserveUsd: MustTruncateDec("138883455.9382328581978198800889056"),
		Token0: UniswapToken{
			Decimals: 18,
			Id:       "0x6b175474e89094c44da98b954eedeac495271d0f",
		},
		Token1: UniswapToken{
			Decimals: 18,
			Id:       "0xc02aaa39b223fe8d0a0e5c4f27ead9083c756cc2",
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
				OracleData: []*codectypes.Any{anyPair},
			},
			true,
		},
		{
			"multiple uniswap pairs",
			OracleFeed{
				OracleData: []*codectypes.Any{anyPair, anyPair2},
			},
			true,
		},
		{
			"dup uniswap pair",
			OracleFeed{
				OracleData: []*codectypes.Any{anyPair, anyPair},
			},
			false,
		},
		{
			"different data types",
			OracleFeed{
				OracleData: []*codectypes.Any{anyPair, anyMockData},
			},
			false,
		},
		{
			"invalid oracle data",
			OracleFeed{
				OracleData: []*codectypes.Any{anyInvPair},
			},
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
