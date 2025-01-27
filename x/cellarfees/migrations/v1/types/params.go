package types

import (
	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
	"gopkg.in/yaml.v2"

	types "github.com/peggyjv/sommelier/v9/x/cellarfees/types"
)

const (
	DefaultFeeAccrualAuctionThreshold uint64 = 2
	// Rough number of blocks in 28 days, the time it takes to unbond
	DefaultRewardEmissionPeriod uint64 = 403200
	// Initial rate at which an auction should decrease the price of the relevant coin from it's starting price.
	// This value was determined experimentally. It is the initial rate at which it takes ~48 hours for the unit
	// price to hit 0 usomm, assuming a decrease acceleration rate of 0.001.
	DefaultInitialPriceDecreaseRate string = "0.0000648"
	// Blocks between each auction price decrease
	DefaultPriceDecreaseBlockInterval uint64 = 10
	// Blocks between each auction
	DefaultAuctionInterval uint64 = 15000
)

// Parameter keys
var (
	KeyFeeAccrualAuctionThreshold = []byte("FeeAccrualAuctionThreshold")
	KeyRewardEmissionPeriod       = []byte("RewardEmissionPeriod")
	KeyInitialPriceDecreaseRate   = []byte("InitialPriceDecreaseRate")
	KeyPriceDecreaseBlockInterval = []byte("PriceDecreaseBlockInterval")
	KeyAuctionInterval            = []byte("AuctionInterval")
)

var _ paramtypes.ParamSet = &Params{}

// ParamKeyTable returns the parameter key table.
func ParamKeyTable() paramtypes.KeyTable {
	return paramtypes.NewKeyTable().RegisterParamSet(&Params{})
}

// DefaultParams returns default cellarfees parameters
func DefaultParams() Params {
	return Params{
		FeeAccrualAuctionThreshold: DefaultFeeAccrualAuctionThreshold,
		RewardEmissionPeriod:       DefaultRewardEmissionPeriod,
		InitialPriceDecreaseRate:   sdk.MustNewDecFromStr(DefaultInitialPriceDecreaseRate),
		PriceDecreaseBlockInterval: DefaultPriceDecreaseBlockInterval,
		AuctionInterval:            DefaultAuctionInterval,
	}
}

// ParamSetPairs returns the parameter set pairs.
func (p *Params) ParamSetPairs() paramtypes.ParamSetPairs {
	return paramtypes.ParamSetPairs{
		paramtypes.NewParamSetPair(KeyFeeAccrualAuctionThreshold, &p.FeeAccrualAuctionThreshold, validateFeeAccrualAuctionThreshold),
		paramtypes.NewParamSetPair(KeyRewardEmissionPeriod, &p.RewardEmissionPeriod, validateRewardEmissionPeriod),
		paramtypes.NewParamSetPair(KeyInitialPriceDecreaseRate, &p.InitialPriceDecreaseRate, validateInitialPriceDecreaseRate),
		paramtypes.NewParamSetPair(KeyPriceDecreaseBlockInterval, &p.PriceDecreaseBlockInterval, validatePriceDecreaseBlockInterval),
		paramtypes.NewParamSetPair(KeyAuctionInterval, &p.AuctionInterval, validateAuctionInterval),
	}
}

// ValidateBasic performs basic validation on cellarfees parameters.
func (p *Params) ValidateBasic() error {
	if err := validateFeeAccrualAuctionThreshold(p.FeeAccrualAuctionThreshold); err != nil {
		return err
	}
	if err := validateRewardEmissionPeriod(p.RewardEmissionPeriod); err != nil {
		return err
	}
	if err := validateInitialPriceDecreaseRate(p.InitialPriceDecreaseRate); err != nil {
		return err
	}
	if err := validatePriceDecreaseBlockInterval(p.PriceDecreaseBlockInterval); err != nil {
		return err
	}
	return nil
}

func validateFeeAccrualAuctionThreshold(i interface{}) error {
	threshold, ok := i.(uint64)
	if !ok {
		return errorsmod.Wrapf(ErrInvalidFeeAccrualAuctionThreshold, "fee accrual auction threshold: %T", i)
	}

	if threshold == 0 {
		return errorsmod.Wrapf(ErrInvalidFeeAccrualAuctionThreshold, "fee accrual auction threshold cannot be zero")
	}

	return nil
}

func validateRewardEmissionPeriod(i interface{}) error {
	emissionPeriod, ok := i.(uint64)
	if !ok {
		return errorsmod.Wrapf(types.ErrInvalidRewardEmissionPeriod, "reward emission period: %T", i)
	}

	if emissionPeriod == 0 {
		return errorsmod.Wrapf(types.ErrInvalidRewardEmissionPeriod, "reward emission period cannot be zero")
	}

	return nil
}

func validateInitialPriceDecreaseRate(i interface{}) error {
	rate, ok := i.(sdk.Dec)
	if !ok {
		return errorsmod.Wrapf(types.ErrInvalidInitialPriceDecreaseRate, "initial price decrease rate: %T", i)
	}

	if rate == sdk.ZeroDec() {
		return errorsmod.Wrapf(types.ErrInvalidInitialPriceDecreaseRate, "initial price decrease rate cannot be zero, must be 0 < x < 1")
	}

	if rate == sdk.OneDec() {
		return errorsmod.Wrapf(types.ErrInvalidInitialPriceDecreaseRate, "initial price decrease rate cannot be one, must be 0 < x < 1")
	}

	return nil
}

func validatePriceDecreaseBlockInterval(i interface{}) error {
	interval, ok := i.(uint64)
	if !ok {
		return errorsmod.Wrapf(types.ErrInvalidPriceDecreaseBlockInterval, "price decrease block interval: %T", i)
	}

	if interval == 0 {
		return errorsmod.Wrapf(types.ErrInvalidPriceDecreaseBlockInterval, "price decrease block interval cannot be zero")
	}

	return nil
}

func validateAuctionInterval(i interface{}) error {
	interval, ok := i.(uint64)
	if !ok {
		return errorsmod.Wrapf(types.ErrInvalidAuctionInterval, "auction interval: %T", i)
	}

	if interval == 0 {
		return errorsmod.Wrapf(types.ErrInvalidAuctionInterval, "auction interval cannot be zero")
	}

	return nil
}

// String implements the String interface
func (p Params) String() string {
	out, _ := yaml.Marshal(p)
	return string(out)
}
