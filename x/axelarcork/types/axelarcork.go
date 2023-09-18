package types

import (
	"bytes"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

func (c *AxelarCork) IDHash(chainID uint64, blockHeight uint64) []byte {
	blockHeightBytes := sdk.Uint64ToBigEndian(blockHeight)
	chainIDBytes := sdk.Uint64ToBigEndian(chainID)

	address := common.HexToAddress(c.TargetContractAddress)

	return crypto.Keccak256Hash(
		bytes.Join(
			[][]byte{chainIDBytes, blockHeightBytes, address.Bytes(), c.EncodedContractCall, sdk.Uint64ToBigEndian(c.Deadline)},
			[]byte{},
		)).Bytes()
}

func (c *AxelarCork) Equals(other AxelarCork) bool {
	firstAddr := common.HexToAddress(c.TargetContractAddress)
	secondAddr := common.HexToAddress(other.TargetContractAddress)

	if firstAddr != secondAddr {
		return false
	}

	if !bytes.Equal(c.EncodedContractCall, other.EncodedContractCall) {
		return false
	}

	if c.Deadline != other.Deadline {
		return false
	}

	return true
}

func (c *AxelarCork) ValidateBasic() error {
	if len(c.EncodedContractCall) == 0 {
		return ErrEmptyContractCall
	}

	if !common.IsHexAddress(c.TargetContractAddress) {
		return ErrInvalidEVMAddress
	}

	return nil
}
