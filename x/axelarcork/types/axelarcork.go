package types

import (
	"bytes"
	"fmt"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

func (c *AxelarCork) IDHash(blockHeight uint64) []byte {
	blockHeightBytes := sdk.Uint64ToBigEndian(blockHeight)
	chainIDBytes := sdk.Uint64ToBigEndian(c.ChainId)
	address := common.HexToAddress(c.TargetContractAddress)

	return crypto.Keccak256Hash(
		bytes.Join(
			[][]byte{blockHeightBytes, chainIDBytes, address.Bytes(), c.EncodedContractCall},
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

	if c.ChainId == 0 {
		return fmt.Errorf("chain ID must be non-zero")
	}

	if !common.IsHexAddress(c.TargetContractAddress) {
		return ErrInvalidEVMAddress
	}

	if c.Deadline == 0 {
		return fmt.Errorf("deadline must be non-zero")
	}

	return nil
}

func (s *ScheduledAxelarCork) ValidateBasic() error {
	if err := s.Cork.ValidateBasic(); err != nil {
		return err
	}

	if s.BlockHeight == 0 {
		return fmt.Errorf("block height must be non-zero")
	}

	if _, err := sdk.ValAddressFromBech32(s.Validator); err != nil {
		return errorsmod.Wrap(sdkerrors.ErrInvalidAddress, err.Error())
	}

	if len(s.Id) != 64 {
		return fmt.Errorf("invalid ID length, must be a keccak256 hash")
	}

	return nil
}

func (c *AxelarCorkResult) ValidateBasic() error {
	if err := c.Cork.ValidateBasic(); err != nil {
		return err
	}

	if c.BlockHeight == 0 {
		return fmt.Errorf("block height must be non-zero")
	}

	if _, err := sdk.NewDecFromStr(c.ApprovalPercentage); err != nil {
		return fmt.Errorf("approval percentage must be a valid Dec")
	}

	return nil
}

func (c *CellarIDSet) ValidateBasic() error {
	if c.ChainId == 0 {
		return fmt.Errorf("chain ID must be non-zero")
	}

	for _, addr := range c.Ids {
		if !common.IsHexAddress(addr) {
			return errorsmod.Wrapf(ErrInvalidEVMAddress, "%s", addr)
		}
	}

	return nil
}
