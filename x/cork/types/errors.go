package types

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// x/cork module sentinel errors
var (
	ErrInvalidEthereumAddress = sdkerrors.Register(ModuleName, 2, "invalid ethereum address")
	ErrUnmanagedCellarAddress = sdkerrors.Register(ModuleName, 3, "cork sent to address that has not passed governance")
	ErrEmptyContractCall      = sdkerrors.Register(ModuleName, 4, "cork has an empty contract call body")
	ErrSchedulingInThePast    = sdkerrors.Register(ModuleName, 5, "cork is trying to be scheduled for a block that has already passed")
	ErrInvalidJson            = sdkerrors.Register(ModuleName, 6, "invalid json")
)
