package types

import (
	fmt "fmt"

	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
	"gopkg.in/yaml.v2"
)

// Parameter keys
var (
	KeyRewardEmissionPeriod = []byte("rewardemissionperiod")
)

var _ paramtypes.ParamSet = &Params{}

// ParamKeyTable returns the parameter key table.
func ParamKeyTable() paramtypes.KeyTable {
	return paramtypes.NewKeyTable().RegisterParamSet(&Params{})
}

// DefaultParams returns default cellarfees parameters
func DefaultParams() Params {
	return Params{
		// Rough number of blocks in 28 days, the time it takes to unbond
		RewardEmissionPeriod: 403200,
	}
}

// ParamSetPairs returns the parameter set pairs.
func (p *Params) ParamSetPairs() paramtypes.ParamSetPairs {
	return paramtypes.ParamSetPairs{
		paramtypes.NewParamSetPair(KeyRewardEmissionPeriod, &p.RewardEmissionPeriod, validateRewardEmissionPeriod),
	}
}

// ValidateBasic performs basic validation on cellarfees parameters.
func (p *Params) ValidateBasic() error {
	if err := validateRewardEmissionPeriod(p.RewardEmissionPeriod); err != nil {
		return err
	}
	return nil
}

func validateRewardEmissionPeriod(i interface{}) error {
	emissionPeriod, ok := i.(uint64)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if emissionPeriod == 0 {
		return fmt.Errorf(
			"emission period should be greater than 0: %d", emissionPeriod,
		)
	}

	return nil
}

// String implements the String interface
func (p Params) String() string {
	out, _ := yaml.Marshal(p)
	return string(out)
}
