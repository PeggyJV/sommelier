package types

import (
	errorsmod "cosmossdk.io/errors"
)

var (
	ErrInvalidDistributionPerBlock = errorsmod.Register(ModuleName, 1, "invalid distribution per block")
)
