package keeper

import (
	"testing"

	tmprotocrypto "github.com/cometbft/cometbft/proto/tendermint/crypto"
	"github.com/stretchr/testify/require"
)

// mergeUpdatesWithBoost indexes validators by consensus pubkey. Keying on the
// raw bytes alone collapsed every unrecognised key type onto the empty string,
// so distinct validators shared one map slot and a raw ValidatorUpdate with an
// equally unrecognised key would match it and be assigned an unrelated
// validator's boosted power -- then emitted to CometBFT.
func TestConsPubKeyIndex(t *testing.T) {
	ed := tmprotocrypto.PublicKey{Sum: &tmprotocrypto.PublicKey_Ed25519{Ed25519: []byte("aaaa")}}
	secp := tmprotocrypto.PublicKey{Sum: &tmprotocrypto.PublicKey_Secp256K1{Secp256K1: []byte("aaaa")}}

	edKey, ok := consPubKeyIndex(ed)
	require.True(t, ok)
	secpKey, ok := consPubKeyIndex(secp)
	require.True(t, ok)

	require.NotEqual(t, edKey, secpKey,
		"identical bytes under different key types must not alias")

	// Unset / unrecognised: must report not-ok rather than yielding a key that
	// every other unrecognised validator would also produce.
	empty, ok := consPubKeyIndex(tmprotocrypto.PublicKey{})
	require.False(t, ok, "an unrecognised key type must be reported, not indexed")
	require.Empty(t, empty)

	zeroLen, ok := consPubKeyIndex(tmprotocrypto.PublicKey{
		Sum: &tmprotocrypto.PublicKey_Ed25519{Ed25519: []byte{}},
	})
	require.False(t, ok, "a zero-length key must not index either")
	require.Empty(t, zeroLen)
}
