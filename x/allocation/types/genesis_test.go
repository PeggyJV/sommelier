package types

import (
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
)

func TestGenesisValidate(t *testing.T) {

	testCases := []struct {
		name     string
		genState GenesisState
		expPass  bool
	}{
		{
			name:     "default",
			genState: DefaultGenesisState(),
			expPass:  true,
		},
		{
			name: "invalid feeder delegator",
			genState: GenesisState{
				Params: DefaultParams(),
				Cellars: []*Cellar{
					{
						common.HexToAddress("0x6ea5992aB4A78D5720bD12A089D13c073d04B55d").String(),
						[]*TickRange{
							{-189720, -192660, 160},
							{-192660, -198540, 680},
							{-198540, -201540, 160},
						},
					},
				},
			},
			expPass: false,
		},
		{
			name: "invalid feeder validator",
			genState: GenesisState{
				Params: DefaultParams(),
				Cellars: []*Cellar{
					{
						common.HexToAddress("0x6ea5992aB4A78D5720bD12A089D13c073d04B55d").String(),
						[]*TickRange{
							{-189720, -192660, 160},
							{-192660, -198540, 680},
							{-198540, -201540, 160},
						},
					},
				},
			},
			expPass: false,
		},
		{
			name: "equal feeder addresses",
			genState: GenesisState{
				Params: DefaultParams(),
				Cellars: []*Cellar{
					{
						common.HexToAddress("0x6ea5992aB4A78D5720bD12A089D13c073d04B55d").String(),
						[]*TickRange{
							{-189720, -192660, 160},
							{-192660, -198540, 680},
							{-198540, -201540, 160},
						},
					},
				},
			},
			expPass: false,
		},
		{
			name: "dup feeder delegation",
			genState: GenesisState{
				Params: DefaultParams(),
				Cellars: []*Cellar{
					{
						common.HexToAddress("0x6ea5992aB4A78D5720bD12A089D13c073d04B55d").String(),
						[]*TickRange{
							{-189720, -192660, 160},
							{-192660, -198540, 680},
							{-198540, -201540, 160},
						},
					},
				},
			},
			expPass: false,
		},
	}

	for _, tc := range testCases {

		err := tc.genState.Validate()
		if tc.expPass {
			require.NoError(t, err, tc.name)
		} else {
			require.Error(t, err, tc.name)
		}
	}
}
