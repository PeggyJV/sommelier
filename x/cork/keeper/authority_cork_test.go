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
	var seenHeight uint64
	var seenID []byte
	var seenContract common.Address
	k.IterateAuthorityCorksByBlockHeight(ctx, 100, func(h uint64, gotID []byte, gotContract common.Address, c v2types.Cork) bool {
		seen = append(seen, c)
		seenHeight, seenID, seenContract = h, gotID, gotContract
		return false
	})
	require.Len(t, seen, 1)
	require.Equal(t, cork.EncodedContractCall, seen[0].EncodedContractCall)

	// The iterator recovers id and contract by slicing the raw store key. Assert
	// them explicitly: the EndBlocker feeds both straight back into
	// DeleteAuthorityCork, so an offset error that still round-trips the VALUE
	// would leave corks undeletable and re-executing every block.
	require.Equal(t, uint64(100), seenHeight)
	require.Equal(t, id, seenID)
	require.Equal(t, contract, seenContract)

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
