package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
	"github.com/peggyjv/sommelier/v6/app/params"
)

// Parameter keys
var (
	KeyDistributionPerBlock   = []byte("DistributionPerBlock")
	KeyIncentivesCutoffHeight = []byte("IncentivesCutoffHeight")
)

var _ paramtypes.ParamSet = &Params{}

// ParamKeyTable returns the parameter key table.
func ParamKeyTable() paramtypes.KeyTable {
	return paramtypes.NewKeyTable().RegisterParamSet(&Params{})
}

// DefaultParams returns default incentives parameters
func DefaultParams() Params {
	return Params{
		DistributionPerBlock: sdk.NewCoin(params.BaseCoinUnit, sdk.ZeroInt()),
		// Anything lower than or equal to current height is "off"
		IncentivesCutoffHeight: 0,
	}
}

// ParamSetPairs returns the parameter set pairs.
func (p *Params) ParamSetPairs() paramtypes.ParamSetPairs {
	return paramtypes.ParamSetPairs{
		paramtypes.NewParamSetPair(KeyDistributionPerBlock, &p.DistributionPerBlock, validateDistributionPerBlock),
		paramtypes.NewParamSetPair(KeyIncentivesCutoffHeight, &p.IncentivesCutoffHeight, validateIncentivesCutoffHeight),
	}
}

// ValidateBasic performs basic validation on incentives parameters.
func (p *Params) ValidateBasic() error {
	if err := validateDistributionPerBlock(p.DistributionPerBlock); err != nil {
		return err
	}
	return nil
}

func validateDistributionPerBlock(i interface{}) error {
	coinsPerBlock, ok := i.(sdk.Coin)
	if !ok {
		return sdkerrors.Wrapf(ErrInvalidDistributionPerBlock, "invalid parameter type: %T", i)
	}

	if coinsPerBlock.IsNil() {
		return sdkerrors.Wrapf(ErrInvalidDistributionPerBlock, "distribution per block cannot be nil")
	}
	if coinsPerBlock.Amount.IsNegative() {
		return sdkerrors.Wrapf(ErrInvalidDistributionPerBlock, "distribution per block cannot be negative")
	}
	if coinsPerBlock.Denom != params.BaseCoinUnit {
		return sdkerrors.Wrapf(ErrInvalidDistributionPerBlock, "distribution per block denom must be %s, got %s", params.BaseCoinUnit, coinsPerBlock.Denom)
	}

	return nil
}

// Since there is no unsigned integer that is an invalid height, this is a no-op.
// Collin: I wasn't sure if it was safe to pass `nil` into the above NewParamSetPair call so I defined this.
func validateIncentivesCutoffHeight(_ interface{}) error {
	return nil
}
