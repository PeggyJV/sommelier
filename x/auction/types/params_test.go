package types

import (
	"testing"

	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/stretchr/testify/require"
)

func TestParamsValidate(t *testing.T) {
	testCases := []struct {
		name    string
		params  Params
		expPass bool
		err     error
	}{
		{
			name:    "Happy path -- default params",
			params:  DefaultParams(),
			expPass: true,
			err:     nil,
		},
		{
			name: "Happy path -- custom params",
			params: Params{
				PriceMaxBlockAge: uint64(1000),
			},
			expPass: true,
			err:     nil,
		},
		{
			name: "Max block age cannot be 0",
			params: Params{
				PriceMaxBlockAge: uint64(0),
			},
			expPass: false,
			err:     sdkerrors.Wrapf(ErrTokenPriceMaxBlockAgeMustBePositive, "value: 0"),
		},
	}

	for _, tc := range testCases {
		err := tc.params.ValidateBasic()
		if tc.expPass {
			require.NoError(t, err, tc.name)
			require.Nil(t, err)
		} else {
			require.Error(t, err, tc.name)
			require.Equal(t, tc.err.Error(), err.Error())
		}
	}
}
