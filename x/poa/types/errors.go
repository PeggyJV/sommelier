package types

import sdkerrors "cosmossdk.io/errors"

var (
	ErrNoBondedAuthority  = sdkerrors.Register(ModuleName, 1, "no bonded, unjailed authority validator: refusing to advance block")
	ErrEmptyAuthoritySet  = sdkerrors.Register(ModuleName, 2, "authority set must not be empty")
	ErrDuplicateAuthority = sdkerrors.Register(ModuleName, 3, "duplicate validator in authority set")
	ErrInvalidFloor       = sdkerrors.Register(ModuleName, 4, "floor_fraction must be in [0.67, 1)")
	ErrUnauthorized       = sdkerrors.Register(ModuleName, 5, "signer is not the configured gov authority")
	ErrSnapshotMissing    = sdkerrors.Register(ModuleName, 6, "multiplier snapshot missing for authority slash")
)
