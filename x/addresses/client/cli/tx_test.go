package cli

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAddAddressMapping(t *testing.T) {
	testCases := []struct {
		name string
		args []string
		err  error
	}{
		{
			name: "Valid cmd",
			args: []string{
				"0xdac17f958d2ee523a2206206994597c13d831ec7",
				fmt.Sprintf("--%s=%s", "from", "cosmos16zrkzad482haunrn25ywvwy6fclh3vh7r0hcny"),
			},
			err: fmt.Errorf("key with address cosmos16zrkzad482haunrn25ywvwy6fclh3vh7r0hcny not found: key not found"), // Expect key not found error since this is just a mock request
		},
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
		{
			name: "Missing 'from' field",
			args: []string{
				"0xdac17f958d2ee523a2206206994597c13d831ec7",
			},
			err: fmt.Errorf("empty address string is not allowed: invalid address"),
		},
		{
			name: "Invalid EVM address",
			args: []string{
				"sdlkfjlskdjfsld",
			},
			err: fmt.Errorf("sdlkfjlskdjfsld is not a valid EVM address"),
		},
	}

	for _, tc := range testCases {
		cmd := GetCmdAddAddressMapping()
		cmd.SetArgs(tc.args)
		err := cmd.Execute()

		require.Equal(t, tc.err.Error(), err.Error())
	}
}
