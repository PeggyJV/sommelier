package types

import (
	"math/rand"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/tendermint/tendermint/crypto"
	"github.com/tendermint/tendermint/crypto/ed25519"
	"github.com/tendermint/tendermint/crypto/secp256k1"
)

// TODO: rewrite these hashing tests

// GeneratePrivKeyAddressPairs generates a total of n private key, address
// pairs.
// TODO: find a better spot for this
func GeneratePrivKeyAddressPairs(n int) (keys []crypto.PrivKey, addrs []sdk.AccAddress) {
	keys = make([]crypto.PrivKey, n)
	addrs = make([]sdk.AccAddress, n)
	for i := 0; i < n; i++ {
		if rand.Int63()%2 == 0 {
			keys[i] = secp256k1.GenPrivKey()
		} else {
			keys[i] = ed25519.GenPrivKey()
		}
		addrs[i] = sdk.AccAddress(keys[i].PubKey().Address())
	}
	return
}
