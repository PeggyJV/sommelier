package keeper

import (
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"

	v2types "github.com/peggyjv/sommelier/v10/x/cork/types/v2"
)

func TestAuthorityCorkRoundTrip(t *testing.T) {
	k, ctx, _, ctrl := setupCorkKeeper(t)
	defer ctrl.Finish()

	contract := common.HexToAddress("0x1111111111111111111111111111111111111111")
	cork := v2types.Cork{
		TargetContractAddress: contract.String(),
		EncodedContractCall:   []byte{0xde, 0xad, 0xbe, 0xef},
	}

	id := k.SetAuthorityCork(ctx, 100, cork)
	require.NotEmpty(t, id)

	var seen []v2types.Cork
	k.IterateAuthorityCorksByBlockHeight(ctx, 100, func(_ uint64, _ []byte, _ common.Address, c v2types.Cork) bool {
		seen = append(seen, c)
		return false
	})
	require.Len(t, seen, 1)
	require.Equal(t, cork.EncodedContractCall, seen[0].EncodedContractCall)

	// A different height must not see it.
	var other int
	k.IterateAuthorityCorksByBlockHeight(ctx, 101, func(_ uint64, _ []byte, _ common.Address, _ v2types.Cork) bool {
		other++
		return false
	})
	require.Zero(t, other)

	k.DeleteAuthorityCork(ctx, 100, id, contract)

	var after int
	k.IterateAuthorityCorksByBlockHeight(ctx, 100, func(_ uint64, _ []byte, _ common.Address, _ v2types.Cork) bool {
		after++
		return false
	})
	require.Zero(t, after)
}
