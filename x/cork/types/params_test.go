package types

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
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
			name:    "empty",
			params:  Params{},
			expPass: false,
		},
		{
			name: "invalid vote period",
			params: Params{
				VotePeriod: 0,
			},
			expPass: false,
		},
		{
			name: "invalid vote threshold",
			params: Params{
				VotePeriod:    5,
				VoteThreshold: sdk.ZeroDec(),
			},
			expPass: false,
		},
		{
			name: "nil vote threshold",
			params: Params{
				VotePeriod:    5,
				VoteThreshold: sdk.Dec{},
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
