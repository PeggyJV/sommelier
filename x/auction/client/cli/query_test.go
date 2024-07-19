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
			err: fmt.Errorf("unknown command \"1\" for \"parameters\""),
		},
	}

	for _, tc := range testCases {
		cmd := *queryParams()
		cmd.SetArgs(tc.args)
		err := cmd.Execute()

		require.Equal(t, tc.err.Error(), err.Error())
	}
}

func TestQueryActiveAuction(t *testing.T) {
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
		cmd := *queryActiveAuction()
		cmd.SetArgs(tc.args)
		err := cmd.Execute()

		require.Equal(t, tc.err.Error(), err.Error())
	}
}

func TestQueryEndedAuction(t *testing.T) {
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
		cmd := *queryEndedAuction()
		cmd.SetArgs(tc.args)
		err := cmd.Execute()

		require.Equal(t, tc.err.Error(), err.Error())
	}
}

func TestActiveAuctionsCmd(t *testing.T) {
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
			err: fmt.Errorf("unknown command \"1\" for \"active-auctions\""),
		},
	}

	for _, tc := range testCases {
		cmd := *queryActiveAuctions()
		cmd.SetArgs(tc.args)
		err := cmd.Execute()

		require.Equal(t, tc.err.Error(), err.Error())
	}
}

func TestEndedAuctionsCmd(t *testing.T) {
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
			err: fmt.Errorf("unknown command \"1\" for \"ended-auctions\""),
		},
	}

	for _, tc := range testCases {
		cmd := *queryEndedAuctions()
		cmd.SetArgs(tc.args)
		err := cmd.Execute()

		require.Equal(t, tc.err.Error(), err.Error())
	}
}

func TestQueryActiveAuctionsByDenom(t *testing.T) {
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
		cmd := *queryActiveAuctionsByDenom()
		cmd.SetArgs(tc.args)
		err := cmd.Execute()

		require.Equal(t, tc.err.Error(), err.Error())
	}
}

func TestQueryEndedAuctionsByDenom(t *testing.T) {
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
		cmd := *queryEndedAuctionsByDenom()
		cmd.SetArgs(tc.args)
		err := cmd.Execute()

		require.Equal(t, tc.err.Error(), err.Error())
	}
}

func TestQueryBids(t *testing.T) {
	testCases := []struct {
		name string
		args []string
		err  error
	}{
		{
			name: "Insufficient args",
			args: []string{
				"1",
			},
			err: fmt.Errorf("accepts 2 arg(s), received 1"),
		},
		{
			name: "Too many args",
			args: []string{
				"1",
				"2",
				"3",
			},
			err: fmt.Errorf("accepts 2 arg(s), received 3"),
		},
	}

	for _, tc := range testCases {
		cmd := *queryBid()
		cmd.SetArgs(tc.args)
		err := cmd.Execute()

		require.Equal(t, tc.err.Error(), err.Error())
	}
}

func TestQueryBidByAuction(t *testing.T) {
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
		cmd := *queryBidsByAuction()
		cmd.SetArgs(tc.args)
		err := cmd.Execute()

		require.Equal(t, tc.err.Error(), err.Error())
	}
}

func TestQueryTokenPrice(t *testing.T) {
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
		cmd := *queryTokenPrice()
		cmd.SetArgs(tc.args)
		err := cmd.Execute()

		require.Equal(t, tc.err.Error(), err.Error())
	}
}

func TestQueryTokenPrices(t *testing.T) {
	testCases := []struct {
		name string
		args []string
		err  error
	}{
		{
			name: "Too many args",
			args: []string{
				"1",
			},
			err: fmt.Errorf("unknown command \"1\" for \"token-prices\""),
		},
	}

	for _, tc := range testCases {
		cmd := *queryTokenPrices()
		cmd.SetArgs(tc.args)
		err := cmd.Execute()

		require.Equal(t, tc.err.Error(), err.Error())
	}
}
