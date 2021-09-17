package types

import (
	"crypto/sha256"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// DataHash returns the hash for a precommit given the proper args
func DataHash(salt, jsonData string, signer sdk.ValAddress) []byte {
	fmt.Printf("salt: %s, json %s, signer %s", salt, jsonData, signer.String())
	h := sha256.New()
	h.Write([]byte(fmt.Sprintf("%s:%s:%s", salt, jsonData, signer.String())))
	return h.Sum(nil)
}

func (ac *Allocation) ValidateBasic() error {
	if ac.Cellar == nil {
		return fmt.Errorf("no cellar provided for allocation")
	}

	if err := ac.Cellar.ValidateBasic(); err != nil {
		return err
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
