package types

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// x/distribution module sentinel errors
var (
	ErrInvalidHash       = sdkerrors.Register(ModuleName, 2, "invalid sha256 hash")
	ErrInvalidOracleData = sdkerrors.Register(ModuleName, 3, "invalid oracle data hash")
	ErrUnknown           = sdkerrors.Register(ModuleName, 4, "unknown")
)
