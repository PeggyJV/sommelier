package types

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// x/distribution module sentinel errors
var (
	ErrInvalidHash       = sdkerrors.Register(ModuleName, 2, "invalid sha256 hash")
	ErrInvalidOracleData = sdkerrors.Register(ModuleName, 2, "invalid oracle data hash")
)
