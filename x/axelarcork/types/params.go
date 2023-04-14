package types

import (
	"errors"
	"fmt"

	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
)

// Parameter keys
var (
	KeyIBCChannel = []byte("ibcchannel")
)

var _ paramtypes.ParamSet = &Params{}

// ParamKeyTable returns the parameter key table.
func ParamKeyTable() paramtypes.KeyTable {
	return paramtypes.NewKeyTable().RegisterParamSet(&Params{})
}

// DefaultParams returns default oracle parameters
func DefaultParams() Params {
	return Params{
		IbcChannel: "unimplemented",
	}
}

// ParamSetPairs returns the parameter set pairs.
func (p *Params) ParamSetPairs() paramtypes.ParamSetPairs {
	return paramtypes.ParamSetPairs{
		paramtypes.NewParamSetPair(KeyIBCChannel, &p.IbcChannel, validateIBCChannel),
	}
}

// ValidateBasic performs basic validation on oracle parameters.
func (p *Params) ValidateBasic() error {
	if err := validateIBCChannel(p.IbcChannel); err != nil {
		return err
	}
	return nil
}

func validateIBCChannel(i interface{}) error {
	ibcChannel, ok := i.(string)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if ibcChannel == "" {
		return errors.New("ibc channel cannot be empty")
	}

	return nil
}
