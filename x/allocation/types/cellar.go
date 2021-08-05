package types

import "github.com/ethereum/go-ethereum/common"

func (c *Cellar) Address() common.Address {
	return common.HexToAddress(c.Id)
}