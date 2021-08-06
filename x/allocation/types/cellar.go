package types

import (
	"github.com/ethereum/go-ethereum/common"
	tmbytes "github.com/tendermint/tendermint/libs/bytes"
)

func (c *Cellar) Address() common.Address {
	return common.HexToAddress(c.Id)
}

func (c *Cellar) InvalidationScope() tmbytes.HexBytes {
	panic("implement me")
}