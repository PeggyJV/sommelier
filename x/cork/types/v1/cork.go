package v1

import (
	"bytes"

	tmbytes "github.com/cometbft/cometbft/libs/bytes"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	corktypes "github.com/peggyjv/sommelier/v7/x/cork/types"
)

func (c *Cork) InvalidationScope() tmbytes.HexBytes {
	addr := common.HexToAddress(c.TargetContractAddress)
	return crypto.Keccak256Hash(
		bytes.Join(
			[][]byte{addr.Bytes(), c.EncodedContractCall},
			[]byte{},
		)).Bytes()
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

func (c *Cork) ValidateBasic() error {
	if len(c.EncodedContractCall) == 0 {
		return corktypes.ErrEmptyContractCall
	}

	if !common.IsHexAddress(c.TargetContractAddress) {
		return corktypes.ErrInvalidEthereumAddress
	}

	return nil
}
