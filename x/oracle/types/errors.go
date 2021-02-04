package types

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// x/distribution module sentinel errors
var (
	ErrInvalidHash       = sdkerrors.Register(ModuleName, 2, "invalid sha256 hash")
	ErrInvalidOracleData = sdkerrors.Register(ModuleName, 3, "invalid oracle data hash")
	ErrUnknown           = sdkerrors.Register(ModuleName, 4, "unknown")
	ErrNoPrevote         = sdkerrors.Register(ModuleName, 5, "no prevote for validator")
	ErrUnpackOracleData  = sdkerrors.Register(ModuleName, 6, "failed to unpack oracle data")
	ErrHashMismatch      = sdkerrors.Register(ModuleName, 7, "precommit hash doesn't match commit hash")
	ErrWrongNumber       = sdkerrors.Register(ModuleName, 8, "wrong number of args")
	ErrWrongDataType     = sdkerrors.Register(ModuleName, 9, "wrong data type")
	ErrParseError        = sdkerrors.Register(ModuleName, 10, "parsing oracle data")
)
