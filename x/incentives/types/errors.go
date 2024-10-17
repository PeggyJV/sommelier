package types

import (
	errorsmod "cosmossdk.io/errors"
)

var (
	ErrInvalidDistributionPerBlock            = errorsmod.Register(ModuleName, 1, "invalid distribution per block")
	ErrInvalidValidatorDistributionPerBlock   = errorsmod.Register(ModuleName, 2, "invalid validator distribution per block")
	ErrInvalidValidatorIncentivesMaxFraction  = errorsmod.Register(ModuleName, 3, "invalid validator incentives max fraction")
	ErrInvalidValidatorIncentivesSetSizeLimit = errorsmod.Register(ModuleName, 4, "invalid validator incentives set size limit")
)
