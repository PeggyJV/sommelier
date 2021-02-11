package types

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// x/oracle module sentinel errors
var (
	ErrInvalidPrevote      = sdkerrors.Register(ModuleName, 2, "invalid prevote hashes")
	ErrInvalidOracleData   = sdkerrors.Register(ModuleName, 3, "invalid oracle data hash")
	ErrNoPrevote           = sdkerrors.Register(ModuleName, 4, "no prevote for validator")
	ErrUnpackOracleData    = sdkerrors.Register(ModuleName, 5, "failed to unpack oracle data")
	ErrHashMismatch        = sdkerrors.Register(ModuleName, 6, "precommit hash doesn't match commit hash")
	ErrUnsupportedDataType = sdkerrors.Register(ModuleName, 8, "unsupported data type")
	ErrParseError          = sdkerrors.Register(ModuleName, 9, "failed to parse oracle data")
)
