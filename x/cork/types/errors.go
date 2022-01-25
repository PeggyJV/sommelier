package types

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// x/cork module sentinel errors
var (
	ErrInvalidAddress = sdkerrors.Register(ModuleName, 2, "invalid ethereum address")
)
