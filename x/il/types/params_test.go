package types

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParamsValidate(t *testing.T) {
	testCases := []struct {
		name    string
		params  Params
		expPass bool
	}{
		{
			name:    "default",
			params:  DefaultParams(),
			expPass: true,
		},
		{
			name: "invalid timeout timestamp secs",
			params: Params{
				BatchContractAddress:     "0xf784709d2317D872237C4bC22f867d1BAe2913AB",
				LiquidityContractAddress: "0xb4e16d0168e52d35cacd2c6185b44281ec28c9dc",
				EthTimeoutBlocks:         10,
				EthTimeoutTimestamp:      180,
			},
			expPass: true,
		},
		{
			name:    "empty",
			params:  Params{},
			expPass: false,
		},
		{
			name: "invalid batch contract",
			params: Params{
				BatchContractAddress: "",
			},
			expPass: false,
		},
		{
			name: "invalid application contract",
			params: Params{
				BatchContractAddress:     "0xf784709d2317D872237C4bC22f867d1BAe2913AB",
				LiquidityContractAddress: "",
			},
			expPass: false,
		},
		{
			name: "invalid timeout blocks",
			params: Params{
				BatchContractAddress:     "0xf784709d2317D872237C4bC22f867d1BAe2913AB",
				LiquidityContractAddress: "0xb4e16d0168e52d35cacd2c6185b44281ec28c9dc",
				EthTimeoutBlocks:         0,
			},
			expPass: false,
		},
		{
			name: "invalid timeout timestamp secs",
			params: Params{
				BatchContractAddress:     "0xf784709d2317D872237C4bC22f867d1BAe2913AB",
				LiquidityContractAddress: "0xb4e16d0168e52d35cacd2c6185b44281ec28c9dc",
				EthTimeoutBlocks:         1,
				EthTimeoutTimestamp:      0,
			},
			expPass: false,
		},
	}

	for _, tc := range testCases {

		err := tc.params.ValidateBasic()
		if tc.expPass {
			require.NoError(t, err, tc.name)
		} else {
			require.Error(t, err, tc.name)
		}
	}
}
