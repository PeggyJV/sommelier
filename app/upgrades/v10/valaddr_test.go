package v10

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/peggyjv/sommelier/v10/app/params"
)

func init() {
	// The bech32 prefix config is process-global and normally installed by the
	// app's init; a bare test binary for this package does not pull that in.
	params.SetAddressPrefixes()
}

func TestDefaultAuthorityValidatorsParse(t *testing.T) {
	if len(DefaultAuthorityValidators) == 0 {
		t.Fatal("DefaultAuthorityValidators is empty; upgrade handler would refuse to run")
	}
	seen := map[string]bool{}
	for _, s := range DefaultAuthorityValidators {
		addr, err := sdk.ValAddressFromBech32(s)
		if err != nil {
			t.Fatalf("%s: %v", s, err)
		}
		if addr.String() != s {
			t.Fatalf("round-trip mismatch: %s != %s", addr.String(), s)
		}
		if seen[s] {
			t.Fatalf("duplicate entry: %s", s)
		}
		seen[s] = true
	}
	t.Logf("%d authority validators parsed", len(DefaultAuthorityValidators))
}
