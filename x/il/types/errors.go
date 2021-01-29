package types

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// Impermanent Loss module errors
var (
	ErrStoplossExists   = sdkerrors.Register(ModuleName, 2, "stoploss already exists for the given address and pair")
	ErrStoplossNotFound = sdkerrors.Register(ModuleName, 3, "stoploss doesn't exist for the given address and pair")
	ErrStoplossInvalid  = sdkerrors.Register(ModuleName, 4, "invalid stoploss")
)
