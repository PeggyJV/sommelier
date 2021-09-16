package types

import (
	"bytes"
	"encoding/hex"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestABIEncodedRebalanceBytes(t *testing.T) {

	rebalanceHash := Cellar{
		Id: "0x0000000000",
		TickRanges: []*TickRange{
			{-189780, -192480, 160},
			{-192480, -197880, 680},
			{-197880, -200640, 160},
		},
	}.ABIEncodedRebalanceHash()

	// hash from python brownie code cc @stevenj
	testHash, err := hex.DecodeString("0xd0f79d9bfeec64dbc27ccd281a20931cfadc07d87875c3289f55383e59f3ebbc"[2:])
	require.NoError(t, err)
	if !bytes.Equal(testHash, rebalanceHash) {
		t.Errorf("gold hash is not equal to generated hash:\n gold hash: %x\n test hash: %x", testHash, rebalanceHash)
	}
}

func TestABIEncodedCellarTickInfoBytes(t *testing.T) {
	tickInfoHash := ABIEncodedCellarTickInfoBytes(0)
	t.Logf("hash: %b", tickInfoHash)
}
