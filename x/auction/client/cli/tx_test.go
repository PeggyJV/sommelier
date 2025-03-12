package cli

import (
	"fmt"
	"os"
	"testing"

	"github.com/peggyjv/sommelier/v9/x/auction/types"

	"github.com/cosmos/cosmos-sdk/testutil"
	moduletestutil "github.com/cosmos/cosmos-sdk/types/module/testutil"
	"github.com/stretchr/testify/require"
)

func TestParseSetTokenPricesProposal(t *testing.T) {
	encodingConfig := moduletestutil.MakeTestEncodingConfig()

	okJSON := testutil.WriteToNewTempFile(t, `
{
	"title": "My token proposal",
	"description":	"Contains a usomm price update",
    "token_prices":	[ { "denom" : "usomm", "exponent": "6", "usd_price" : "4.200000000000000000"} ],
	"deposit": "10000usomm"
}
`)

	proposal := types.SetTokenPricesProposalWithDeposit{}
	contents, err := os.ReadFile(okJSON.Name())
	require.NoError(t, err)

	err = encodingConfig.Codec.UnmarshalJSON(contents, &proposal)
	require.NoError(t, err)

	require.Equal(t, "My token proposal", proposal.Title)
	require.Equal(t, "Contains a usomm price update", proposal.Description)
	require.Equal(t, "denom:\"usomm\" exponent:6 usd_price:\"4200000000000000000\" ", proposal.TokenPrices[0].String())
	require.Equal(t, "10000usomm", proposal.Deposit)
}

func TestSubmitBid(t *testing.T) {
	testCases := []struct {
		name string
		args []string
		err  error
	}{
		// {
		// 	name: "Valid cmd",
		// 	args: []string{
		// 		"1",
		// 		"10000usomm",
		// 		"50000gravity0xdac17f958d2ee523a2206206994597c13d831ec7",
		// 		fmt.Sprintf("--%s=%s", "from", "cosmos16zrkzad482haunrn25ywvwy6fclh3vh7r0hcny"),
		// 	},
		// 	err: fmt.Errorf("cosmos16zrkzad482haunrn25ywvwy6fclh3vh7r0hcny.info: key not found"), // Expect key not found error since this is just a mock request
		// },
		{
			name: "Insufficient args",
			args: []string{},
			err:  fmt.Errorf("accepts 3 arg(s), received 0"),
		},
		{
			name: "Too many args",
			args: []string{
				"1",
				"2",
				"3",
				"4",
			},
			err: fmt.Errorf("accepts 3 arg(s), received 4"),
		},
		{
			name: "Missing 'from' field",
			args: []string{
				"1",
				"10000usomm",
				"50000gravity0xdac17f958d2ee523a2206206994597c13d831ec7",
			},
			err: fmt.Errorf("must include `--from` flag"),
		},
		{
			name: "Invalid bid",
			args: []string{
				"1",
				"10000",
				"50000gravity0xdac17f958d2ee523a2206206994597c13d831ec7",
			},
			err: fmt.Errorf("invalid decimal coin expression: 10000"),
		},
		{
			name: "Invalid minimum amount",
			args: []string{
				"1",
				"10000usomm",
				"50000",
			},
			err: fmt.Errorf("invalid decimal coin expression: 50000"),
		},
	}

	for _, tc := range testCases {
		cmd := GetCmdSubmitBid()
		cmd.SetArgs(tc.args)
		err := cmd.Execute()

		require.Equal(t, tc.err.Error(), err.Error())
	}
}
