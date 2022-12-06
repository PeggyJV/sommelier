package cli

import (
	"testing"

	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/stretchr/testify/require"
)

func TestQueryParamsCmd(t *testing.T) {
	testCases := []struct {
		name string
		args []string
		err  error
	}{
		{
			name: "Does not accept args",
			args: []string{
				"1",
			},
			err: sdkerrors.New("", uint32(1), "unknown command \"1\" for \"parameters\""),
		},
	}

	for _, tc := range testCases {
		cmd := *queryParams()
		cmd.SetArgs(tc.args)
		err := cmd.Execute()

		require.Equal(t, tc.err.Error(), err.Error())
	}
}
