package types

import (
	"bytes"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	tmbytes "github.com/tendermint/tendermint/libs/bytes"
)

func (c *Cork) InvalidationScope() tmbytes.HexBytes {
	addr := common.HexToAddress(c.TargetContractAddress)
	return crypto.Keccak256Hash(
		bytes.Join(
			[][]byte{addr.Bytes(), c.EncodedContractCall},
			[]byte{},
		)).Bytes()
}

func (c *Cork) IdHash(blockHeight uint64) []byte {
	blockHeightBytes := sdk.Uint64ToBigEndian(blockHeight)

	return crypto.Keccak256Hash(
		bytes.Join(
			[][]byte{blockHeightBytes, []byte(c.TargetContractAddress), c.EncodedContractCall},
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
		return ErrEmptyContractCall
	}

	if !common.IsHexAddress(c.TargetContractAddress) {
		return ErrInvalidAddress
	}

	return nil
}
