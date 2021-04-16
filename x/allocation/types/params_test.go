package types

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
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
			expPass: true,
		},
		{
			name:    "empty",
			params:  Params{},
			expPass: false,
		},
		{
			name: "invalid vote period",
			params: Params{
				VotePeriod: 0,
			},
			expPass: false,
		},
		{
			name: "invalid vote threshold",
			params: Params{
				VotePeriod:    5,
				VoteThreshold: sdk.ZeroDec(),
			},
			expPass: false,
		},
		{
			name: "nil vote threshold",
			params: Params{
				VotePeriod:    5,
				VoteThreshold: sdk.Dec{},
			},
			expPass: false,
		},
		{
			name: "invalid slash window",
			params: Params{
				VotePeriod:    5,
				VoteThreshold: sdk.NewDecFromIntWithPrec(sdk.OneInt(), 2),
				SlashWindow:   0,
			},
			expPass: false,
		},
		{
			name: "invalid min valid window",
			params: Params{
				VotePeriod:        5,
				VoteThreshold:     sdk.NewDecFromIntWithPrec(sdk.OneInt(), 2),
				SlashWindow:       1,
				MinValidPerWindow: sdk.ZeroDec(),
			},
			expPass: false,
		},
		{
			name: "nil min valid window",
			params: Params{
				VotePeriod:        5,
				VoteThreshold:     sdk.NewDecFromIntWithPrec(sdk.NewInt(25), 3),
				SlashWindow:       1,
				MinValidPerWindow: sdk.Dec{},
			},
			expPass: false,
		},
		{
			name: "invalid slash fraction",
			params: Params{
				VotePeriod:        5,
				VoteThreshold:     sdk.NewDecFromIntWithPrec(sdk.NewInt(25), 3),
				SlashWindow:       1,
				MinValidPerWindow: sdk.NewDecFromIntWithPrec(sdk.OneInt(), 2),
				SlashFraction:     sdk.ZeroDec(),
			},
			expPass: false,
		},
		{
			name: "nil slash fraction",
			params: Params{
				VotePeriod:        5,
				VoteThreshold:     sdk.NewDecFromIntWithPrec(sdk.NewInt(25), 3),
				SlashWindow:       1,
				MinValidPerWindow: sdk.NewDecFromIntWithPrec(sdk.OneInt(), 2),
				SlashFraction:     sdk.Dec{},
			},
			expPass: false,
		},
		{
			name: "invalid target threshold",
			params: Params{
				VotePeriod:        5,
				VoteThreshold:     sdk.NewDecFromIntWithPrec(sdk.NewInt(25), 3),
				SlashWindow:       1,
				MinValidPerWindow: sdk.NewDecFromIntWithPrec(sdk.OneInt(), 2),
				SlashFraction:     sdk.NewDecFromIntWithPrec(sdk.OneInt(), 2),
				TargetThreshold:   sdk.ZeroDec(),
			},
			expPass: false,
		},
		{
			name: "nil target threshold",
			params: Params{
				VotePeriod:        5,
				VoteThreshold:     sdk.NewDecFromIntWithPrec(sdk.NewInt(25), 3),
				SlashWindow:       1,
				MinValidPerWindow: sdk.NewDecFromIntWithPrec(sdk.OneInt(), 2),
				SlashFraction:     sdk.NewDecFromIntWithPrec(sdk.OneInt(), 2),
				TargetThreshold:   sdk.Dec{},
			},
			expPass: false,
		},
		{
			name: "invalid data types",
			params: Params{
				VotePeriod:        5,
				VoteThreshold:     sdk.NewDecFromIntWithPrec(sdk.NewInt(25), 3),
				SlashWindow:       1,
				MinValidPerWindow: sdk.NewDecFromIntWithPrec(sdk.OneInt(), 2),
				SlashFraction:     sdk.NewDecFromIntWithPrec(sdk.OneInt(), 2),
				TargetThreshold:   sdk.NewDecFromIntWithPrec(sdk.OneInt(), 2),
				DataTypes:         []string{""},
			},
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
