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

func (c *Cellar) abiEncodedTicks() (uppers []int32, lowers []int32, weights []uint32) {
	for _, tick := range c.TickRanges {
		uppers = append(uppers, tick.Upper)
		lowers = append(lowers, tick.Lower)
		weights = append(weights, tick.Weight)
	}

	return
}