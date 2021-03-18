package types

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/cosmos/cosmos-sdk/testutil/testdata"
)

func TestStoplossKey(t *testing.T) {
	_, _, addr := testdata.KeyTestPubAddr()

	key := append(StoplossKeyPrefix, StoplossKey(addr, "random_pair")...)
	lpAddr := LPAddressFromStoplossKey(key)

	require.Equal(t, addr.String(), lpAddr.String())
}

func TestSubmittedPositionKey(t *testing.T) {
	_, _, addr := testdata.KeyTestPubAddr()

	key := append(SubmittedPositionsQueuePrefix, SubmittedPositionKey(10, addr)...)
	timeoutHeight, address := SplitSubmittedStoplossKey(key)

	require.Len(t, key, 29)
	require.Equal(t, 10, int(timeoutHeight))
	require.Equal(t, addr.String(), address.String())
}
