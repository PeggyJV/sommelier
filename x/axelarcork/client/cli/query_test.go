package cli

import (
	"testing"

	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/stretchr/testify/require"
)

func TestQueryScheduledCorksByIDCmd(t *testing.T) {
	testCases := []struct {
		name string
		args []string
		err  error
	}{
		{
			name: "Invalid ID",
			args: []string{
				"bad",
			},
			err: sdkerrors.New("", uint32(1), "invalid ID length, must be a keccak256 hash"),
		},
	}

	for _, tc := range testCases {
		cmd := *queryScheduledCorksByID()
		cmd.SetArgs(tc.args)
		err := cmd.Execute()

		require.Equal(t, tc.err.Error(), err.Error())
	}
}

func TestQueryCorkResultCmd(t *testing.T) {
	testCases := []struct {
		name string
		args []string
		err  error
	}{
		{
			name: "Invalid ID",
			args: []string{
				"1",
				"bad",
			},
			err: sdkerrors.New("", uint32(1), "invalid ID length, must be a keccak256 hash"),
		},
	}

	for _, tc := range testCases {
		cmd := *queryCorkResult()
		cmd.SetArgs(tc.args)
		err := cmd.Execute()

		require.Equal(t, tc.err.Error(), err.Error())
	}
}
