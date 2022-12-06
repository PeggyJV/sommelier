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

func TestQueryScheduledCorksCmd(t *testing.T) {
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
			err: sdkerrors.New("", uint32(1), "unknown command \"1\" for \"scheduled-corks\""),
		},
	}

	for _, tc := range testCases {
		cmd := *queryScheduledCorks()
		cmd.SetArgs(tc.args)
		err := cmd.Execute()

		require.Equal(t, tc.err.Error(), err.Error())
	}
}

func TestQueryCellarIDsCmd(t *testing.T) {
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
			err: sdkerrors.New("", uint32(1), "unknown command \"1\" for \"cellar-ids\""),
		},
	}

	for _, tc := range testCases {
		cmd := *queryCellarIDs()
		cmd.SetArgs(tc.args)
		err := cmd.Execute()

		require.Equal(t, tc.err.Error(), err.Error())
	}
}

func TestQueryScheduledBlockHeightsCmd(t *testing.T) {
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
			err: sdkerrors.New("", uint32(1), "unknown command \"1\" for \"scheduled-block-heights\""),
		},
	}

	for _, tc := range testCases {
		cmd := *queryScheduledBlockHeights()
		cmd.SetArgs(tc.args)
		err := cmd.Execute()

		require.Equal(t, tc.err.Error(), err.Error())
	}
}

func TestQueryScheduledCorksByBlockHeightCmd(t *testing.T) {
	testCases := []struct {
		name string
		args []string
		err  error
	}{
		{
			name: "Insufficient args",
			args: []string{},
			err:  sdkerrors.New("", uint32(1), "accepts 1 arg(s), received 0"),
		},
		{
			name: "Block height overflow",
			args: []string{
				"18446744073709551616",
			},
			err: sdkerrors.New("", uint32(1), "strconv.Atoi: parsing \"18446744073709551616\": value out of range"),
		},
	}

	for _, tc := range testCases {
		cmd := *queryScheduledCorksByBlockHeight()
		cmd.SetArgs(tc.args)
		err := cmd.Execute()

		require.Equal(t, tc.err.Error(), err.Error())
	}
}

func TestQueryScheduledCorksByIDCmd(t *testing.T) {
	testCases := []struct {
		name string
		args []string
		err  error
	}{
		{
			name: "Insufficient args",
			args: []string{},
			err:  sdkerrors.New("", uint32(1), "accepts 1 arg(s), received 0"),
		},
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
			name: "Insufficient args",
			args: []string{},
			err:  sdkerrors.New("", uint32(1), "accepts 1 arg(s), received 0"),
		},
		{
			name: "Invalid ID",
			args: []string{
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

func TestQueryCorkResultsCmd(t *testing.T) {
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
			err: sdkerrors.New("", uint32(1), "unknown command \"1\" for \"cork-results\""),
		},
	}

	for _, tc := range testCases {
		cmd := *queryCorkResults()
		cmd.SetArgs(tc.args)
		err := cmd.Execute()

		require.Equal(t, tc.err.Error(), err.Error())
	}
}
