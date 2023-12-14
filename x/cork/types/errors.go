package types

import (
	errorsmod "cosmossdk.io/errors"
)

// x/cork module sentinel errors
var (
	ErrInvalidEthereumAddress       = errorsmod.Register(ModuleName, 2, "invalid ethereum address")
	ErrUnmanagedCellarAddress       = errorsmod.Register(ModuleName, 3, "cork sent to address that has not passed governance")
	ErrEmptyContractCall            = errorsmod.Register(ModuleName, 4, "cork has an empty contract call body")
	ErrSchedulingInThePast          = errorsmod.Register(ModuleName, 5, "cork is trying to be scheduled for a block that has already passed")
	ErrInvalidJSON                  = errorsmod.Register(ModuleName, 6, "invalid json")
	ErrValidatorCorkCapacityReached = errorsmod.Register(ModuleName, 7, "validator cork capacity reached")
)
