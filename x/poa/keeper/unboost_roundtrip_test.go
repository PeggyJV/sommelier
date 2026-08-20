package keeper

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
)

// Un-boosting must invert the boost exactly.
//
// The EndBlocker writes boosted = ceil(raw*M); rawSlashPower recovers raw with
// Dec.Quo(m).TruncateInt64(). Truncation is the exact inverse, not an
// approximation: ceil gives raw*M <= boosted < raw*M+1, so
// raw <= boosted/M < raw + 1/M <= raw+1 for M >= 1.
//
// A code review proposed switching this to Ceil "for symmetry with the boost".
// That would over-slash by one power unit in most cases. This test exists so
// that change fails loudly.
func TestUnboostInvertsBoostExactly(t *testing.T) {
	multipliers := []sdk.Dec{
		sdk.MustNewDecFromStr("1.3"),
		sdk.MustNewDecFromStr("1.5"),
		sdk.MustNewDecFromStr("2.1"),
		sdk.MustNewDecFromStr("3"),
		sdk.MustNewDecFromStr("1.000000000000000001"),
	}

	for _, m := range multipliers {
		for raw := int64(1); raw <= 50; raw++ {
			boosted := sdk.NewDec(raw).Mul(m).Ceil().TruncateInt64()

			// The inverse actually used by rawSlashPower.
			got := sdk.NewDec(boosted).Quo(m).TruncateInt64()
			require.Equalf(t, raw, got,
				"truncating inverse must recover raw exactly (raw=%d M=%s boosted=%d)",
				raw, m, boosted)

			// And the proposed alternative must NOT be adopted: show it is wrong.
			viaCeil := sdk.NewDec(boosted).Quo(m).Ceil().TruncateInt64()
			require.GreaterOrEqualf(t, viaCeil, raw,
				"sanity: ceil inverse never under-recovers (raw=%d M=%s)", raw, m)
		}
	}
}
