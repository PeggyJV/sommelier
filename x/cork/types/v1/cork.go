package v1

import (
	"bytes"
	"fmt"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	corktypes "github.com/peggyjv/sommelier/v7/x/cork/types"
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

func (c *Cork) IDHash(blockHeight uint64) []byte {
	blockHeightBytes := sdk.Uint64ToBigEndian(blockHeight)
	chainIDBytes := sdk.Uint64ToBigEndian(1) // corks are on eth mainnet
	address := common.HexToAddress(c.TargetContractAddress)

	return crypto.Keccak256Hash(
		bytes.Join(
			[][]byte{blockHeightBytes, chainIDBytes, address.Bytes(), c.EncodedContractCall},
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

func (s *ScheduledCork) ValidateBasic() error {
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

func (c *CorkResult) ValidateBasic() error {
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
	for _, addr := range c.Ids {
		if !common.IsHexAddress(addr) {
			return fmt.Errorf("invalid EVM address: %s", addr)
		}
	}

	return nil
}
