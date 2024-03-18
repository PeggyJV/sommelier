package types

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGenesisState_Validate(t *testing.T) {
	def := DefaultGenesisState()
	testCases := []struct {
		desc     string
		genState *GenesisState
		valid    bool
	}{
		{
			desc:     "default is valid",
			genState: &def,
			valid:    true,
		},
		{
			desc:     "valid genesis state",
			genState: &GenesisState{},
			valid:    true,
		},
	}

	for _, tc := range testCases {
		err := tc.genState.Validate()
		if tc.valid {
			require.NoError(t, err)
		} else {
			require.Error(t, err)
		}
	}
}
