package types

import (
	"fmt"

	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
)

// Parameter keys
var (
	KeyPriceMaxBlockAge = []byte("PriceMaxBlockAge")
	MinimumBidInUsomm   = []byte("MinimumBidInUsomm")
)

var _ paramtypes.ParamSet = &Params{}

// ParamKeyTable returns the parameter key table.
func ParamKeyTable() paramtypes.KeyTable {
	return paramtypes.NewKeyTable().RegisterParamSet(&Params{})
}

// DefaultParams returns default auction parameters
func DefaultParams() Params {
	return Params{
		PriceMaxBlockAge:  403200,  // roughly four weeks based on 6 second blocks
		MinimumBidInUsomm: 1000000, // 1 somm
	}
}

// ParamSetPairs returns the parameter set pairs.
func (p *Params) ParamSetPairs() paramtypes.ParamSetPairs {
	return paramtypes.ParamSetPairs{
		paramtypes.NewParamSetPair(KeyPriceMaxBlockAge, &p.PriceMaxBlockAge, validatePriceMaxBlockAge),
		paramtypes.NewParamSetPair(MinimumBidInUsomm, &p.MinimumBidInUsomm, validateMinimumBidInUsomm),
	}
}

// ValidateBasic performs basic validation on auction parameters.
func (p *Params) ValidateBasic() error {
	if err := validatePriceMaxBlockAge(p.PriceMaxBlockAge); err != nil {
		return err
	}

	if err := validateMinimumBidInUsomm(p.MinimumBidInUsomm); err != nil {
		return err
	}

	return nil
}

func validatePriceMaxBlockAge(i interface{}) error {
	priceMaxBlockAge, ok := i.(uint64)
	if !ok {
		return fmt.Errorf("invalid price max block age parameter type: %T", i)
	}

	if priceMaxBlockAge == 0 {
		return fmt.Errorf(
			"price max block age must be non-zero",
		)
	}

	return nil
}

func validateMinimumBidInUsomm(i interface{}) error {
	_, ok := i.(uint64)
	if !ok {
		return fmt.Errorf("invalid minimum bid in usomm parameter type: %T", i)
	}

	return nil
}
