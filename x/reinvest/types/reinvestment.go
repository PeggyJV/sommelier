package types

import (
	"bytes"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	tmbytes "github.com/tendermint/tendermint/libs/bytes"
)

func (rv *Reinvestment) InvalidationScope() tmbytes.HexBytes {
	return crypto.Keccak256Hash(rv.Body).Bytes()
}

func (rv *Reinvestment) Equals(other Reinvestment) bool {
	firstAddr := common.HexToAddress(rv.Address)
	secondAddr := common.HexToAddress(other.Address)

	if firstAddr.Hex() != secondAddr.Hex() {
		return false
	}

	if !bytes.Equal(rv.Body, other.Body) {
		return false
	}

	return true
}
