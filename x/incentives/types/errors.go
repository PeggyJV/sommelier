package types

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var (
	ErrInvalidDistributionPerBlock = sdkerrors.Register(ModuleName, 1, "invalid distribution per block")
)
