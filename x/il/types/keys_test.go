package types

import (
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"

	"github.com/cosmos/cosmos-sdk/testutil/testdata"
)

func TestStoplossKey(t *testing.T) {
	_, _, addr := testdata.KeyTestPubAddr()
	pairAddr := common.HexToAddress("0x3041cbd36888becc7bbcbc0045e3b1f144466f5f")

	key := append(StoplossKeyPrefix, StoplossKey(addr, pairAddr.String())...)
	lpAddr := LPAddressFromStoplossKey(key)

	require.Equal(t, addr.String(), lpAddr.String())
}

func TestSubmittedPositionKey(t *testing.T) {
	_, _, addr := testdata.KeyTestPubAddr()
	pairAddr := common.HexToAddress("0x3041cbd36888becc7bbcbc0045e3b1f144466f5f")

	key := append(SubmittedPositionsQueuePrefix, SubmittedPositionKey(10, addr, pairAddr)...)
	height, address, pairAddress := SplitSubmittedStoplossKey(key)

	require.Equal(t, 10, int(height))
	require.Equal(t, pairAddr.String(), pairAddress.String())
	require.Equal(t, addr.String(), address.String())
}
