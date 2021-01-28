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
