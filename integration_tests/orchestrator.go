package integration_tests

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type orchestrator struct {
	index     int
	mnemonic  string
	keyRecord keyring.Record
	keyring   *keyring.Keyring
}

func (o *orchestrator) instanceName() string {
	return fmt.Sprintf("orchestrator%d", o.index)
}

func (o *orchestrator) address() sdk.AccAddress {
	addr, err := o.keyRecord.GetAddress()
	if err != nil {
		panic(err)
	}

	return addr
}
