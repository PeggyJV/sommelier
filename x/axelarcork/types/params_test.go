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
			name:    "empty",
			params:  Params{},
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
