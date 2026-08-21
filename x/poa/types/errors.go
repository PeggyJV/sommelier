package types

import sdkerrors "cosmossdk.io/errors"

var (
	ErrNoBondedAuthority  = sdkerrors.Register(ModuleName, 1, "no bonded, unjailed authority validator: refusing to advance block")
	ErrEmptyAuthoritySet  = sdkerrors.Register(ModuleName, 2, "authority set must not be empty")
	ErrDuplicateAuthority = sdkerrors.Register(ModuleName, 3, "duplicate validator in authority set")
	ErrInvalidFloor       = sdkerrors.Register(ModuleName, 4, "floor_fraction must be in [0.67, 1)")
	ErrUnauthorized       = sdkerrors.Register(ModuleName, 5, "signer is not the configured gov authority")
	ErrSnapshotMissing    = sdkerrors.Register(ModuleName, 6, "multiplier snapshot missing for authority slash")
	// ErrSafeModeGovFrozen guards x/poa's own gov messages while the chain is in
	// authority-empty safe mode, where governance is decided entirely by the
	// community set the freeze exists to distrust. Recovery is by unjailing and
	// re-bonding the existing authority validators, not by governance.
	ErrSafeModeGovFrozen = sdkerrors.Register(ModuleName, 7,
		"x/poa safe mode active: authority-set and param changes are frozen; restore the existing authority set by unjailing and re-bonding")
	ErrInvalidGenesis = sdkerrors.Register(ModuleName, 8, "invalid poa genesis state")
)
