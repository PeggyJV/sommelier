package v2

import (
	"bytes"
	"fmt"

	errorsmod "cosmossdk.io/errors"
	tmbytes "github.com/cometbft/cometbft/libs/bytes"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	corktypes "github.com/peggyjv/sommelier/v10/x/cork/types"
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
	// An absent embedded message decodes to nil, so this is reachable from a
	// malformed genesis file or a decoded msg. Reject it rather than
	// dereferencing: InitGenesis dereferences Cork directly afterwards, so a
	// panic here would abort InitChain instead of failing validation.
	if s.Cork == nil {
		return fmt.Errorf("scheduled cork must carry a cork")
	}
	if err := s.Cork.ValidateBasic(); err != nil {
		return err
	}

	if s.BlockHeight == 0 {
		return fmt.Errorf("block height must be non-zero")
	}

	// Validator is vestigial as of v10: corks are scheduled by the cork
	// authority, which is an account, not a validator. The field is retained on
	// the type for wire compatibility and is empty on every cork the chain now
	// produces, so it is validated only when set.
	if s.Validator != "" {
		if _, err := sdk.ValAddressFromBech32(s.Validator); err != nil {
			return errorsmod.Wrap(sdkerrors.ErrInvalidAddress, err.Error())
		}
	}

	// Id is the raw 32-byte keccak256 digest, not its hex encoding. This
	// previously compared against 64 and so could never pass, which made any
	// genesis carrying a scheduled cork unimportable.
	if len(s.Id) != 32 {
		return fmt.Errorf("invalid ID length, must be a 32-byte keccak256 hash")
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
