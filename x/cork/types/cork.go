package types

import (
	"bytes"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	tmbytes "github.com/tendermint/tendermint/libs/bytes"
)

func (c *Cork) InvalidationScope() tmbytes.HexBytes {
	return crypto.Keccak256Hash(c.EncodedContractCall).Bytes()
}

func (c *Cork) Equals(other Cork) bool {
	firstAddr := common.HexToAddress(c.TargetContractAddress)
	secondAddr := common.HexToAddress(other.TargetContractAddress)

	if firstAddr != secondAddr {
		return false
	}

	if !bytes.Equal(c.EncodedContractCall, other.EncodedContractCall) {
		return false
	}

	return true
}
