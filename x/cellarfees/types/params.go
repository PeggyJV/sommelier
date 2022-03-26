package types

import (
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
	"gopkg.in/yaml.v2"
)

var _ paramtypes.ParamSet = &Params{}

// ParamKeyTable returns the parameter key table.
func ParamKeyTable() paramtypes.KeyTable {
	return paramtypes.NewKeyTable().RegisterParamSet(&Params{})
}

// DefaultParams returns default cellarfees parameters
func DefaultParams() Params {
	return Params{}
}

// ParamSetPairs returns the parameter set pairs.
func (p *Params) ParamSetPairs() paramtypes.ParamSetPairs {
	return paramtypes.ParamSetPairs{}
}

// ValidateBasic performs basic validation on cellarfees parameters.
func (p *Params) ValidateBasic() error {
	return nil
}

// String implements the String interface
func (p Params) String() string {
	out, _ := yaml.Marshal(p)
	return string(out)
}
