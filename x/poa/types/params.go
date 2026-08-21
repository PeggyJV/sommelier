package types

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
)

// Param store keys.
var (
	KeyFloorFraction          = []byte("FloorFraction")
	KeyEnabled                = []byte("Enabled")
	KeyHaltWhenAuthorityEmpty = []byte("HaltWhenAuthorityEmpty")

	// DefaultFloorFraction is just over 2/3 to guarantee a strict supermajority.
	DefaultFloorFraction = sdk.MustNewDecFromStr("0.670000000000000001")
)

// ParamKeyTable returns the param key table for the PoA module.
func ParamKeyTable() paramtypes.KeyTable {
	return paramtypes.NewKeyTable().RegisterParamSet(&Params{})
}

// DefaultParams returns the default PoA parameters.
//
// HaltWhenAuthorityEmpty defaults to false, selecting authority-empty SAFE MODE
// over a hard halt: when the bonded authority set collapses the chain keeps
// producing blocks on community stake, while the value-bearing modules
// (gravity, cork, axelarcork) freeze so nothing is committed under the
// untrusted set. Set it to true to instead fail-closed and halt the chain
// (recovery then requires an off-chain coordinated restart).
//
// Recovery from safe mode is by UNJAILING OR RE-BONDING an existing authority
// validator -- NOT by governance. MsgUpdateAuthoritySet and MsgUpdateParams are
// both rejected with ErrSafeModeGovFrozen while safe mode is active, precisely
// because governance is community-only there and Enabled=false would clear the
// freeze with no thaw delay. Choosing false on the assumption that a
// governance escape hatch exists would be a mistake. See docs/poa.md.
func DefaultParams() Params {
	return Params{
		FloorFraction:          DefaultFloorFraction,
		Enabled:                true,
		HaltWhenAuthorityEmpty: false,
	}
}

// ParamSetPairs implements paramtypes.ParamSet.
func (p *Params) ParamSetPairs() paramtypes.ParamSetPairs {
	return paramtypes.ParamSetPairs{
		paramtypes.NewParamSetPair(KeyFloorFraction, &p.FloorFraction, validateFloor),
		paramtypes.NewParamSetPair(KeyEnabled, &p.Enabled, validateBool),
		paramtypes.NewParamSetPair(KeyHaltWhenAuthorityEmpty, &p.HaltWhenAuthorityEmpty, validateBool),
	}
}

// Validate runs all param validators.
func (p Params) Validate() error {
	if err := validateFloor(p.FloorFraction); err != nil {
		return err
	}
	if err := validateBool(p.Enabled); err != nil {
		return err
	}
	return validateBool(p.HaltWhenAuthorityEmpty)
}

func validateFloor(i interface{}) error {
	v, ok := i.(sdk.Dec)
	if !ok {
		return fmt.Errorf("invalid type for floor_fraction: %T", i)
	}
	if v.IsNil() {
		return fmt.Errorf("floor_fraction must not be nil")
	}
	// Enforce the design invariant: the floor must guarantee at least a 2/3
	// supermajority for the authority set. Allowing anything down to >0.5 would
	// let governance silently weaken the core safety property of this module.
	if v.LT(DefaultFloorFraction) || v.GTE(sdk.OneDec()) {
		return ErrInvalidFloor
	}
	return nil
}

func validateBool(i interface{}) error {
	if _, ok := i.(bool); !ok {
		return fmt.Errorf("expected bool, got %T", i)
	}
	return nil
}
