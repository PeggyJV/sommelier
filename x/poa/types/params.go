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
func DefaultParams() Params {
	return Params{
		FloorFraction:          DefaultFloorFraction,
		Enabled:                true,
		HaltWhenAuthorityEmpty: true,
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
	if !v.GT(sdk.MustNewDecFromStr("0.5")) || v.GTE(sdk.OneDec()) {
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
