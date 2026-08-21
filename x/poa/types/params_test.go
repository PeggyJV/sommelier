package types

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
)

func TestValidateFloor(t *testing.T) {
	cases := []struct {
		name  string
		floor string
		valid bool
	}{
		{"default is valid", DefaultFloorFraction.String(), true},
		{"exactly two-thirds rejected (below strict supermajority)", "0.666666666666666666", false},
		{"just above 0.5 rejected", "0.51", false},
		{"0.6 rejected", "0.6", false},
		{"0.75 accepted", "0.75", true},
		{"one rejected", "1.0", false},
		{"above one rejected", "1.5", false},
		{"zero rejected", "0.0", false},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			err := validateFloor(sdk.MustNewDecFromStr(tc.floor))
			if tc.valid {
				require.NoError(t, err)
			} else {
				require.Error(t, err)
			}
		})
	}

	require.Error(t, validateFloor("not a dec"), "non-Dec type must be rejected")
	require.Error(t, validateFloor(sdk.Dec{}), "nil Dec must be rejected")
}
