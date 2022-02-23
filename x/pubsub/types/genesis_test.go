package types

import (
	"testing"

	"github.com/stretchr/testify/require"
)

// TODO(bolten): fill out genesis test

func TestGenesisState_Validate(t *testing.T) {
	for _, tc := range []struct {
		desc     string
		genState GenesisState
		valid    bool
	}{
		{
			desc:     "default is valid",
			genState: DefaultGenesisState(),
			valid:    true,
		},
		{
			desc:     "valid genesis state",
			genState: GenesisState{},
			valid:    true,
		},
	} {
		t.Run(tc.desc, func(t *testing.T) {
			err := tc.genState.Validate()
			if tc.valid {
				require.NoError(t, err)
			} else {
				require.Error(t, err)
			}
		})
	}
}
