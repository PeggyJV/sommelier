package types

import (
	errorsmod "cosmossdk.io/errors"
)

// x/cork module sentinel errors
var (
	ErrInvalidEVMAddress                  = errorsmod.Register(ModuleName, 2, "invalid evm address")
	ErrUnmanagedCellarAddress             = errorsmod.Register(ModuleName, 3, "cork sent to address that has not passed governance")
	ErrEmptyContractCall                  = errorsmod.Register(ModuleName, 4, "cork has an empty contract call body")
	ErrSchedulingInThePast                = errorsmod.Register(ModuleName, 5, "cork is trying to be scheduled for a block that has already passed")
	ErrInvalidJSON                        = errorsmod.Register(ModuleName, 6, "invalid json")
	ErrValuelessSend                      = errorsmod.Register(ModuleName, 7, "transferring an empty token amount")
	ErrDisabled                           = errorsmod.Register(ModuleName, 8, "axelar disabled")
	ErrValidatorAxelarCorkCapacityReached = errorsmod.Register(ModuleName, 9, "validator Axelar cork capacity reached")
)
