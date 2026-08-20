package keeper

import (
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"

	v2types "github.com/peggyjv/sommelier/v10/x/cork/types/v2"
)

// Cork IDs recovered from the raw store key must be copies, not sub-slices of
// the iterator's key buffer.
//
// The EndBlocker collects ids during iteration and feeds them to
// DeleteAuthorityCork only after the iterator has closed. If the ids alias a
// reused buffer, every collected id degrades to the last one seen, the deletes
// target the wrong keys, and the corks are stranded in state forever -- nothing
// ever queries that height again.
func TestAuthorityCorkIteratorReturnsIndependentIDs(t *testing.T) {
	k, ctx, _, ctrl := setupCorkKeeper(t)
	defer ctrl.Finish()

	height := uint64(500)
	want := map[string]bool{}
	for _, hexAddr := range []string{
		"0x1111111111111111111111111111111111111111",
		"0x2222222222222222222222222222222222222222",
		"0x3333333333333333333333333333333333333333",
	} {
		cork := v2types.Cork{
			TargetContractAddress: common.HexToAddress(hexAddr).String(),
			EncodedContractCall:   []byte(hexAddr),
		}
		want[string(k.SetAuthorityCork(ctx, height, cork))] = true
	}
	require.Len(t, want, 3, "test setup must produce three distinct cork IDs")

	var got [][]byte
	k.IterateAuthorityCorksByBlockHeight(ctx, height, func(_ uint64, id []byte, _ common.Address, _ v2types.Cork) bool {
		got = append(got, id)
		return false
	})
	require.Len(t, got, 3)

	// Compared AFTER iteration finishes, which is when aliasing would show.
	distinct := map[string]bool{}
	for _, id := range got {
		distinct[string(id)] = true
		require.True(t, want[string(id)], "recovered an ID that was never stored")
	}
	require.Len(t, distinct, 3,
		"collected IDs collapsed to fewer than three: they alias the iterator key buffer")
}
