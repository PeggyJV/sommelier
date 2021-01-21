package types

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// Impermanent Loss module errors
var (
	ErrStoplossExists = sdkerrors.Register(ModuleName, 2, "stoploss doesn't exist for the given address and pai")
)
