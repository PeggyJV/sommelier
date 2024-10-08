package cli

import (
	"fmt"
	"testing"

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
			err: fmt.Errorf("unknown command \"1\" for \"params\""),
		},
	}

	for _, tc := range testCases {
		cmd := *CmdQueryParams()
		cmd.SetArgs(tc.args)
		err := cmd.Execute()

		require.Equal(t, tc.err.Error(), err.Error())
	}
}

func TestQueryAddressMappings(t *testing.T) {
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
			err: fmt.Errorf("unknown command \"1\" for \"address-mappings\""),
		},
	}

	for _, tc := range testCases {
		cmd := *CmdQueryAddressMappings()
		cmd.SetArgs(tc.args)
		err := cmd.Execute()

		require.Equal(t, tc.err.Error(), err.Error())
	}
}

func TestQueryAddressMappingByCosmosAddress(t *testing.T) {
	testCases := []struct {
		name string
		args []string
		err  error
	}{
		{
			name: "Insufficient args",
			args: []string{},
			err:  fmt.Errorf("accepts 1 arg(s), received 0"),
		},
		{
			name: "Too many args",
			args: []string{
				"1",
				"2",
			},
			err: fmt.Errorf("accepts 1 arg(s), received 2"),
		},
	}

	for _, tc := range testCases {
		cmd := *CmdQueryAddressMappingByCosmosAddress()
		cmd.SetArgs(tc.args)
		err := cmd.Execute()

		require.Equal(t, tc.err.Error(), err.Error())
	}
}

func TestAddressMappingByEVMAddressCmd(t *testing.T) {
	testCases := []struct {
		name string
		args []string
		err  error
	}{
		{
			name: "Insufficient args",
			args: []string{},
			err:  fmt.Errorf("accepts 1 arg(s), received 0"),
		},
		{
			name: "Too many args",
			args: []string{
				"1",
				"2",
			},
			err: fmt.Errorf("accepts 1 arg(s), received 2"),
		},
	}

	for _, tc := range testCases {
		cmd := *CmdQueryAddressMappingByEVMAddress()
		cmd.SetArgs(tc.args)
		err := cmd.Execute()

		require.Equal(t, tc.err.Error(), err.Error())
	}
}
