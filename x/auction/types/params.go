package types

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
)

// Parameter keys
var (
	KeyPriceMaxBlockAge                     = []byte("PriceMaxBlockAge")
	KeyMinimumBidInUsomm                    = []byte("MinimumBidInUsomm")
	KeyAuctionMaxBlockAge                   = []byte("AuctionMaxBlockAge")
	KeyAuctionPriceDecreaseAccelerationRate = []byte("AuctionPriceDecreaseAccelerationRate")
)

var _ paramtypes.ParamSet = &Params{}

// ParamKeyTable returns the parameter key table.
func ParamKeyTable() paramtypes.KeyTable {
	return paramtypes.NewKeyTable().RegisterParamSet(&Params{})
}

// DefaultParams returns default auction parameters
func DefaultParams() Params {
	return Params{
		PriceMaxBlockAge:                     403200,                       // roughly four weeks based on 6 second blocks
		MinimumBidInUsomm:                    1000000,                      // 1 somm
		AuctionMaxBlockAge:                   864000,                       // roughly 60 days based on 6 second blocks
		AuctionPriceDecreaseAccelerationRate: sdk.MustNewDecFromStr("0.1"), // 10%
	}
}

// ParamSetPairs returns the parameter set pairs.
func (p *Params) ParamSetPairs() paramtypes.ParamSetPairs {
	return paramtypes.ParamSetPairs{
		paramtypes.NewParamSetPair(KeyPriceMaxBlockAge, &p.PriceMaxBlockAge, validatePriceMaxBlockAge),
		paramtypes.NewParamSetPair(KeyMinimumBidInUsomm, &p.MinimumBidInUsomm, validateMinimumBidInUsomm),
		paramtypes.NewParamSetPair(KeyAuctionMaxBlockAge, &p.AuctionMaxBlockAge, validateAuctionMaxBlockAge),
		paramtypes.NewParamSetPair(KeyAuctionPriceDecreaseAccelerationRate, &p.AuctionPriceDecreaseAccelerationRate, validateAuctionPriceDecreaseAccelerationRate),
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

	if err := validateAuctionMaxBlockAge(p.AuctionMaxBlockAge); err != nil {
		return err
	}

	if err := validateAuctionPriceDecreaseAccelerationRate(p.AuctionPriceDecreaseAccelerationRate); err != nil {
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

func validateAuctionMaxBlockAge(i interface{}) error {
	auctionMaxBlockAge, ok := i.(uint64)
	if !ok {
		return fmt.Errorf("invalid blocks to not prune parameter type: %T", i)
	}

	if auctionMaxBlockAge == 0 {
		return fmt.Errorf("blocks to not prune must be non-zero")
	}

	return nil
}

func validateAuctionPriceDecreaseAccelerationRate(i interface{}) error {
	auctionPriceDecreaseAccelerationRate, ok := i.(sdk.Dec)
	if !ok {
		return fmt.Errorf("invalid auction price decrease acceleration rate parameter type: %T", i)
	}

	if auctionPriceDecreaseAccelerationRate.LT(sdk.MustNewDecFromStr("0")) || auctionPriceDecreaseAccelerationRate.GT(sdk.MustNewDecFromStr("1.0")) {
		// Acceleration rates could in theory be more than 100% if need be, but we are establishing this as a bound for now
		return fmt.Errorf("auction price decrease acceleration rate must be betwen 0 and 1 inclusive (0%% to 100%%)")
	}

	return nil
}
