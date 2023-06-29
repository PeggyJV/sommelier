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
			expPass: false,
		},
		{
			name: "fake-configuration",
			params: Params{
				IbcChannel:      "test-channel",
				IbcPort:         "test-port",
				GmpAccount:      "test-account",
				ExecutorAccount: "test-executor",
				TimeoutDuration: 0,
			},
			expPass: false,
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
