package types

import (
	"crypto/sha256"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// DataHash returns the hash for a precommit given the proper args
func DataHash(salt, jsonData string, signer sdk.ValAddress) []byte {
	h := sha256.New()
	h.Write([]byte(fmt.Sprintf("%s:%s:%s", salt, jsonData, signer.String())))
	return h.Sum(nil)
}

func (ac *Allocation) ValidateBasic() error {
	err := ac.Vote.Cellar.ValidateBasic()
	if err != nil {
		return err
	}

	if ac.Vote.CurrentPrice == 0 {
		return fmt.Errorf("invalid current price of zero")
	}

	return nil
}

func (c *Cellar) ValidateBasic() error {
	if c.Id == "" {
		return fmt.Errorf("no cellar id provided")
	}

	if len(c.TickRanges) == 0 {
		return fmt.Errorf("no tick ranges provided for cellar")
	}

	return nil
}
