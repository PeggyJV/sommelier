package integration_tests

import (
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type proposer struct {
	mnemonic  string
	keyRecord keyring.Record
	keyring   *keyring.Keyring
}

func (o *proposer) address() sdk.AccAddress {
	addr, err := o.keyRecord.GetAddress()
	if err != nil {
		panic(err)
	}

	return addr
}
