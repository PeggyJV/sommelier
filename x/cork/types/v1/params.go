package v1

import (
	"errors"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
)

// Parameter keys
var (
	KeyVotePeriod           = []byte("voteperiod")
	KeyVoteThreshold        = []byte("votethreshold")
	KeyMaxCorksPerValidator = []byte("maxcorkspervalidator")
)

var _ paramtypes.ParamSet = &Params{}

// ParamKeyTable returns the parameter key table.
func ParamKeyTable() paramtypes.KeyTable {
	return paramtypes.NewKeyTable().RegisterParamSet(&Params{})
}

// DefaultParams returns default oracle parameters
func DefaultParams() Params {
	return Params{
		// Deprecated
		VoteThreshold:        sdk.NewDecWithPrec(67, 2), // 67%
		MaxCorksPerValidator: 1000,
	}
}

// ParamSetPairs returns the parameter set pairs.
func (p *Params) ParamSetPairs() paramtypes.ParamSetPairs {
	return paramtypes.ParamSetPairs{
		paramtypes.NewParamSetPair(KeyVoteThreshold, &p.VoteThreshold, validateVoteThreshold),
		paramtypes.NewParamSetPair(KeyMaxCorksPerValidator, &p.MaxCorksPerValidator, validateMaxCorksPerValidator),
	}
}

// ValidateBasic performs basic validation on oracle parameters.
func (p *Params) ValidateBasic() error {
	if err := validateVoteThreshold(p.VoteThreshold); err != nil {
		return err
	}
	if err := validateMaxCorksPerValidator(p.MaxCorksPerValidator); err != nil {
		return err
	}
	return nil
}

func validateVoteThreshold(i interface{}) error {
	voteThreshold, ok := i.(sdk.Dec)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if voteThreshold.IsNil() {
		return errors.New("vote threshold cannot be nil")
	}

	if voteThreshold.LT(sdk.ZeroDec()) || voteThreshold.GT(sdk.OneDec()) {
		return fmt.Errorf("vote threshold value must be within the 0% - 100% range, got: %s", voteThreshold)
	}

	return nil
}

func validateMaxCorksPerValidator(i interface{}) error {
	maxCorksPerValidator, ok := i.(uint64)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if maxCorksPerValidator == 0 {
		return errors.New("max corks per validator cannot be 0")
	}

	return nil
}
