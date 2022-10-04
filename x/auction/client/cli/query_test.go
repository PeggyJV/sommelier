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

func TestQueryActiveAuction(t *testing.T) {
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
			name: "Too many args",
			args: []string{
				"1",
				"2",
			},
			err: sdkerrors.New("", uint32(1), "accepts 1 arg(s), received 2"),
		},
		{
			name: "Auction ID overflow",
			args: []string{
				"4294967296",
			},
			err: sdkerrors.New("", uint32(1), "strconv.ParseUint: parsing \"4294967296\": value out of range"),
		},
		{
			name: "Auction ID invalid type",
			args: []string{
				"one hundred and twenty",
			},
			err: sdkerrors.New("", uint32(1), "strconv.ParseUint: parsing \"one hundred and twenty\": invalid syntax"),
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
			err:  sdkerrors.New("", uint32(1), "accepts 1 arg(s), received 0"),
		},
		{
			name: "Too many args",
			args: []string{
				"1",
				"2",
			},
			err: sdkerrors.New("", uint32(1), "accepts 1 arg(s), received 2"),
		},
		{
			name: "Auction ID invalid type",
			args: []string{
				"one hundred and twenty",
			},
			err: sdkerrors.New("", uint32(1), "strconv.ParseUint: parsing \"one hundred and twenty\": invalid syntax"),
		},
		{
			name: "Auction ID overflow",
			args: []string{
				"4294967296",
			},
			err: sdkerrors.New("", uint32(1), "strconv.ParseUint: parsing \"4294967296\": value out of range"),
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
			err: sdkerrors.New("", uint32(1), "unknown command \"1\" for \"active-auctions\""),
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
			err: sdkerrors.New("", uint32(1), "unknown command \"1\" for \"ended-auctions\""),
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
			err:  sdkerrors.New("", uint32(1), "accepts 1 arg(s), received 0"),
		},
		{
			name: "Too many args",
			args: []string{
				"1",
				"2",
			},
			err: sdkerrors.New("", uint32(1), "accepts 1 arg(s), received 2"),
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
			err:  sdkerrors.New("", uint32(1), "accepts 1 arg(s), received 0"),
		},
		{
			name: "Too many args",
			args: []string{
				"1",
				"2",
			},
			err: sdkerrors.New("", uint32(1), "accepts 1 arg(s), received 2"),
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
			err: sdkerrors.New("", uint32(1), "accepts 2 arg(s), received 1"),
		},
		{
			name: "Too many args",
			args: []string{
				"1",
				"2",
				"3",
			},
			err: sdkerrors.New("", uint32(1), "accepts 2 arg(s), received 3"),
		},
		{
			name: "Auction ID overflow",
			args: []string{
				"4294967296",
				"2",
			},
			err: sdkerrors.New("", uint32(1), "strconv.ParseUint: parsing \"4294967296\": value out of range"),
		},
		{
			name: "Bid ID overflow",
			args: []string{
				"1",
				"18446744073709551616",
			},
			err: sdkerrors.New("", uint32(1), "strconv.ParseUint: parsing \"18446744073709551616\": value out of range"),
		},
		{
			name: "Auction ID invalid type",
			args: []string{
				"one hundred and twenty",
				"2",
			},
			err: sdkerrors.New("", uint32(1), "strconv.ParseUint: parsing \"one hundred and twenty\": invalid syntax"),
		},
		{
			name: "Bid ID invalid type",
			args: []string{
				"1",
				"four",
			},
			err: sdkerrors.New("", uint32(1), "strconv.ParseUint: parsing \"four\": invalid syntax"),
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
			err:  sdkerrors.New("", uint32(1), "accepts 1 arg(s), received 0"),
		},
		{
			name: "Too many args",
			args: []string{
				"1",
				"2",
			},
			err: sdkerrors.New("", uint32(1), "accepts 1 arg(s), received 2"),
		},
		{
			name: "Auction ID overflow",
			args: []string{
				"4294967296",
			},
			err: sdkerrors.New("", uint32(1), "strconv.ParseUint: parsing \"4294967296\": value out of range"),
		},
		{
			name: "Auction ID invalid type",
			args: []string{
				"one hundred and twenty",
			},
			err: sdkerrors.New("", uint32(1), "strconv.ParseUint: parsing \"one hundred and twenty\": invalid syntax"),
		},
	}

	for _, tc := range testCases {
		cmd := *queryBidsByAuction()
		cmd.SetArgs(tc.args)
		err := cmd.Execute()

		require.Equal(t, tc.err.Error(), err.Error())
	}
}
