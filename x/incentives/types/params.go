package types

import (
	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
	"github.com/peggyjv/sommelier/v9/app/params"
)

// Parameter keys
var (
	KeyDistributionPerBlock             = []byte("DistributionPerBlock")
	KeyIncentivesCutoffHeight           = []byte("IncentivesCutoffHeight")
	KeyValidatorMaxDistributionPerBlock = []byte("ValidatorMaxDistributionPerBlock")
	KeyValidatorIncentivesCutoffHeight  = []byte("ValidatorIncentivesCutoffHeight")
	KeyValidatorIncentivesMaxFraction   = []byte("ValidatorIncentivesMaxFraction")
	KeyValidatorIncentivesSetSizeLimit  = []byte("ValidatorIncentivesSetSizeLimit")
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
		IncentivesCutoffHeight:           0,
		ValidatorMaxDistributionPerBlock: sdk.NewCoin(params.BaseCoinUnit, sdk.ZeroInt()),
		ValidatorIncentivesCutoffHeight:  0,
		ValidatorIncentivesMaxFraction:   sdk.MustNewDecFromStr("0.1"),
		ValidatorIncentivesSetSizeLimit:  50,
	}
}

// ParamSetPairs returns the parameter set pairs.
func (p *Params) ParamSetPairs() paramtypes.ParamSetPairs {
	return paramtypes.ParamSetPairs{
		paramtypes.NewParamSetPair(KeyDistributionPerBlock, &p.DistributionPerBlock, validateDistributionPerBlock),
		paramtypes.NewParamSetPair(KeyIncentivesCutoffHeight, &p.IncentivesCutoffHeight, validateIncentivesCutoffHeight),
		paramtypes.NewParamSetPair(KeyValidatorMaxDistributionPerBlock, &p.ValidatorMaxDistributionPerBlock, validateValidatorMaxDistributionPerBlock),
		paramtypes.NewParamSetPair(KeyValidatorIncentivesCutoffHeight, &p.ValidatorIncentivesCutoffHeight, validateValidatorIncentivesCutoffHeight),
		paramtypes.NewParamSetPair(KeyValidatorIncentivesMaxFraction, &p.ValidatorIncentivesMaxFraction, validateValidatorIncentivesMaxFraction),
		paramtypes.NewParamSetPair(KeyValidatorIncentivesSetSizeLimit, &p.ValidatorIncentivesSetSizeLimit, validateValidatorIncentivesSetSizeLimit),
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
		return errorsmod.Wrapf(ErrInvalidDistributionPerBlock, "invalid parameter type: %T", i)
	}

	if coinsPerBlock.IsNil() {
		return errorsmod.Wrapf(ErrInvalidDistributionPerBlock, "distribution per block cannot be nil")
	}
	if !coinsPerBlock.IsValid() {
		return errorsmod.Wrapf(ErrInvalidDistributionPerBlock, "distribution per block must be valid")
	}
	if coinsPerBlock.Denom != params.BaseCoinUnit {
		return errorsmod.Wrapf(ErrInvalidDistributionPerBlock, "distribution per block denom must be %s, got %s", params.BaseCoinUnit, coinsPerBlock.Denom)
	}

	return nil
}

// Since there is no unsigned integer that is an invalid height, this is a no-op.
// Collin: I wasn't sure if it was safe to pass `nil` into the above NewParamSetPair call so I defined this.
func validateIncentivesCutoffHeight(_ interface{}) error {
	return nil
}

func validateValidatorMaxDistributionPerBlock(i interface{}) error {
	if err := validateDistributionPerBlock(i); err != nil {
		return errorsmod.Wrapf(ErrInvalidValidatorMaxDistributionPerBlock, "invalid parameter type: %T", i)
	}

	return nil
}

func validateValidatorIncentivesCutoffHeight(i interface{}) error {
	return nil
}

func validateValidatorIncentivesMaxFraction(i interface{}) error {
	dec, ok := i.(sdk.Dec)
	if !ok {
		return errorsmod.Wrapf(ErrInvalidValidatorIncentivesMaxFraction, "invalid parameter type: %T", i)
	}

	if dec.IsNegative() {
		return errorsmod.Wrapf(ErrInvalidValidatorIncentivesMaxFraction, "validator incentives max fraction cannot be negative")
	}

	if dec.GT(sdk.OneDec()) {
		return errorsmod.Wrapf(ErrInvalidValidatorIncentivesMaxFraction, "validator incentives max fraction cannot be greater than one")
	}

	return nil
}

func validateValidatorIncentivesSetSizeLimit(i interface{}) error {
	_, ok := i.(uint64)
	if !ok {
		return errorsmod.Wrapf(ErrInvalidValidatorIncentivesSetSizeLimit, "invalid parameter type: %T", i)
	}

	return nil
}
