package integration_tests

import (
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
)

type proposer struct {
	mnemonic string
	keyInfo  keyring.Info
	keyring  *keyring.Keyring
}
