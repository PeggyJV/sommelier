package types

import (
	"fmt"

	"github.com/ethereum/go-ethereum/common"
)

func (rv *ReinvestVote) ValidateBasic() error {
	for _, v := range rv.Cellar {
		if err := v.ValidateBasic(); err != nil {
			return err
		}
	}
	return nil
}

func (c *Cellar) ValidateBasic() error {
	if !common.IsHexAddress(c.Id) {
		return fmt.Errorf("needs to be a valid ethereum address")
	}
	return nil
}
