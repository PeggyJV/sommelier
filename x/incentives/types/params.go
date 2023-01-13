package types

import (
	"errors"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
)

// Parameter keys
var (
	KeyDistributionPerBlock = []byte("DistributionPerBlock")
)

var _ paramtypes.ParamSet = &Params{}

// ParamKeyTable returns the parameter key table.
func ParamKeyTable() paramtypes.KeyTable {
	return paramtypes.NewKeyTable().RegisterParamSet(&Params{})
}

// DefaultParams returns default oracle parameters
func DefaultParams() Params {
	return Params{
		// 2 somm per block
		DistributionPerBlock: sdk.NewCoin("usomm", sdk.NewInt(2_000_000)),
	}
}

// ParamSetPairs returns the parameter set pairs.
func (p *Params) ParamSetPairs() paramtypes.ParamSetPairs {
	return paramtypes.ParamSetPairs{
		paramtypes.NewParamSetPair(KeyDistributionPerBlock, &p.DistributionPerBlock, validateDistributionPerBlock),
	}
}

// ValidateBasic performs basic validation on oracle parameters.
func (p *Params) ValidateBasic() error {
	if err := validateDistributionPerBlock(p.DistributionPerBlock); err != nil {
		return err
	}
	return nil
}

func validateDistributionPerBlock(i interface{}) error {
	coinsPerBlock, ok := i.(sdk.Coin)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if coinsPerBlock.IsNil() {
		return errors.New("distribution per block cannot be nil")
	}

	return nil
}
