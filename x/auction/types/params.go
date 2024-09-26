package types

import (
	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
)

// Parameter keys
var (
	KeyPriceMaxBlockAge                     = []byte("PriceMaxBlockAge")
	KeyMinimumBidInUsomm                    = []byte("MinimumBidInUsomm")
	KeyMinimumSaleTokensUSDValue            = []byte("MinimumSaleTokensUSDValue")
	KeyAuctionMaxBlockAge                   = []byte("AuctionMaxBlockAge")
	KeyAuctionPriceDecreaseAccelerationRate = []byte("AuctionPriceDecreaseAccelerationRate")
	KeyMinimumAuctionHeight                 = []byte("MinimumAuctionHeight")
	KeyAuctionBurnRate                      = []byte("AuctionBurnRate")
)

var _ paramtypes.ParamSet = &Params{}

// ParamKeyTable returns the parameter key table.
func ParamKeyTable() paramtypes.KeyTable {
	return paramtypes.NewKeyTable().RegisterParamSet(&Params{})
}

// DefaultParams returns default auction parameters
func DefaultParams() Params {
	return Params{
		PriceMaxBlockAge:                     806400,                         // roughly eight weeks based on 6 second blocks
		MinimumBidInUsomm:                    1000000,                        // 1 somm
		MinimumSaleTokensUsdValue:            sdk.MustNewDecFromStr("1.0"),   // unimplemented currently -- minimum value of sale tokens to consider starting an auction
		AuctionMaxBlockAge:                   864000,                         // roughly 60 days based on 6 second blocks
		AuctionPriceDecreaseAccelerationRate: sdk.MustNewDecFromStr("0.001"), // 0.1%
		MinimumAuctionHeight:                 0,                              // do not run auctions before this block height
		AuctionBurnRate:                      sdk.MustNewDecFromStr("0.5"),   // 50%
	}
}

// ParamSetPairs returns the parameter set pairs.
func (p *Params) ParamSetPairs() paramtypes.ParamSetPairs {
	return paramtypes.ParamSetPairs{
		paramtypes.NewParamSetPair(KeyPriceMaxBlockAge, &p.PriceMaxBlockAge, validatePriceMaxBlockAge),
		paramtypes.NewParamSetPair(KeyMinimumBidInUsomm, &p.MinimumBidInUsomm, validateMinimumBidInUsomm),
		paramtypes.NewParamSetPair(KeyMinimumSaleTokensUSDValue, &p.MinimumSaleTokensUsdValue, validateMinimumSaleTokensUSDValue),
		paramtypes.NewParamSetPair(KeyAuctionMaxBlockAge, &p.AuctionMaxBlockAge, validateAuctionMaxBlockAge),
		paramtypes.NewParamSetPair(KeyAuctionPriceDecreaseAccelerationRate, &p.AuctionPriceDecreaseAccelerationRate, validateAuctionPriceDecreaseAccelerationRate),
		paramtypes.NewParamSetPair(KeyMinimumAuctionHeight, &p.MinimumAuctionHeight, validateMinimumAuctionHeight),
		paramtypes.NewParamSetPair(KeyAuctionBurnRate, &p.AuctionBurnRate, validateAuctionBurnRate),
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

	if err := validateMinimumSaleTokensUSDValue(p.MinimumSaleTokensUsdValue); err != nil {
		return err
	}

	if err := validateAuctionMaxBlockAge(p.AuctionMaxBlockAge); err != nil {
		return err
	}

	if err := validateAuctionPriceDecreaseAccelerationRate(p.AuctionPriceDecreaseAccelerationRate); err != nil {
		return err
	}

	if err := validateAuctionBurnRate(p.AuctionBurnRate); err != nil {
		return err
	}

	return nil
}

func validatePriceMaxBlockAge(i interface{}) error {
	priceMaxBlockAge, ok := i.(uint64)
	if !ok {
		return errorsmod.Wrapf(ErrInvalidPriceMaxBlockAgeParameterType, "type: %T", i)
	}

	if priceMaxBlockAge == 0 {
		return errorsmod.Wrapf(ErrTokenPriceMaxBlockAgeMustBePositive, "value: %d", priceMaxBlockAge)
	}

	return nil
}

func validateMinimumBidInUsomm(i interface{}) error {
	_, ok := i.(uint64)
	if !ok {
		return errorsmod.Wrapf(ErrMinimumBidParam, "invalid minimum bid in usomm parameter type: %T", i)
	}

	return nil
}

func validateMinimumSaleTokensUSDValue(i interface{}) error {
	minimumSaleTokensUsdValue, ok := i.(sdk.Dec)
	if !ok {
		return errorsmod.Wrapf(ErrInvalidMinimumSaleTokensUSDValue, "invalid minimum sale tokens USD value parameter type: %T", i)
	}

	if minimumSaleTokensUsdValue.LT(sdk.MustNewDecFromStr("1.0")) {
		// Setting this to a minimum of 1.0 USD to ensure we can realistically charge a non-fractional usomm value
		return errorsmod.Wrapf(ErrInvalidMinimumSaleTokensUSDValue, "minimum sale tokens USD value must be at least 1.0")
	}

	return nil
}

func validateAuctionMaxBlockAge(i interface{}) error {
	_, ok := i.(uint64)
	if !ok {
		return errorsmod.Wrapf(ErrInvalidAuctionMaxBlockAgeParam, "invalid auction max block age parameter type: %T", i)
	}

	return nil
}

func validateAuctionPriceDecreaseAccelerationRate(i interface{}) error {
	auctionPriceDecreaseAccelerationRate, ok := i.(sdk.Dec)
	if !ok {
		return errorsmod.Wrapf(ErrInvalidAuctionPriceDecreaseAccelerationRateParam, "invalid auction price decrease acceleration rate parameter type: %T", i)
	}

	if auctionPriceDecreaseAccelerationRate.LT(sdk.MustNewDecFromStr("0")) || auctionPriceDecreaseAccelerationRate.GT(sdk.MustNewDecFromStr("1.0")) {
		// Acceleration rates could in theory be more than 100% if need be, but we are establishing this as a bound for now
		return errorsmod.Wrapf(ErrInvalidAuctionPriceDecreaseAccelerationRateParam, "auction price decrease acceleration rate must be between 0 and 1 inclusive (0%% to 100%%)")
	}

	return nil
}

func validateMinimumAuctionHeight(i interface{}) error {
	_, ok := i.(uint64)
	if !ok {
		return errorsmod.Wrapf(ErrInvalidMinimumAuctionHeightParam, "invalid minimum auction height parameter type: %T", i)
	}

	return nil
}

func validateAuctionBurnRate(i interface{}) error {
	auctionBurnRate, ok := i.(sdk.Dec)
	if !ok {
		return errorsmod.Wrapf(ErrInvalidAuctionBurnRateParam, "invalid auction burn rate parameter type: %T", i)
	}

	if auctionBurnRate.LT(sdk.MustNewDecFromStr("0")) || auctionBurnRate.GT(sdk.MustNewDecFromStr("1.0")) {
		return errorsmod.Wrapf(ErrInvalidAuctionBurnRateParam, "auction burn rate must be between 0 and 1 inclusive (0%% to 100%%)")
	}

	return nil
}
