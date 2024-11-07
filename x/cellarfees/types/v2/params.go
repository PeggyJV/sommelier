package v2

import (
	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
	"github.com/peggyjv/sommelier/v8/x/cellarfees/types"
	"gopkg.in/yaml.v2"
)

const (
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
	// Minimum USD value of a token's fees balance to trigger an auction
	// $10,000
	DefaultAuctionThresholdUsdValue = "10000.00"
)

// Parameter keys
var (
	KeyRewardEmissionPeriod       = []byte("RewardEmissionPeriod")
	KeyInitialPriceDecreaseRate   = []byte("InitialPriceDecreaseRate")
	KeyPriceDecreaseBlockInterval = []byte("PriceDecreaseBlockInterval")
	KeyAuctionInterval            = []byte("AuctionInterval")
	KeyAuctionThresholdUsdValue   = []byte("AuctionThresholdUsdValue")
)

var _ paramtypes.ParamSet = &Params{}

// ParamKeyTable returns the parameter key table.
func ParamKeyTable() paramtypes.KeyTable {
	return paramtypes.NewKeyTable().RegisterParamSet(&Params{})
}

// DefaultParams returns default cellarfees parameters
func DefaultParams() Params {
	return Params{
		RewardEmissionPeriod:       DefaultRewardEmissionPeriod,
		InitialPriceDecreaseRate:   sdk.MustNewDecFromStr(DefaultInitialPriceDecreaseRate),
		PriceDecreaseBlockInterval: DefaultPriceDecreaseBlockInterval,
		AuctionInterval:            DefaultAuctionInterval,
		AuctionThresholdUsdValue:   sdk.MustNewDecFromStr(DefaultAuctionThresholdUsdValue),
	}
}

// ParamSetPairs returns the parameter set pairs.
func (p *Params) ParamSetPairs() paramtypes.ParamSetPairs {
	return paramtypes.ParamSetPairs{
		paramtypes.NewParamSetPair(KeyRewardEmissionPeriod, &p.RewardEmissionPeriod, validateRewardEmissionPeriod),
		paramtypes.NewParamSetPair(KeyInitialPriceDecreaseRate, &p.InitialPriceDecreaseRate, validateInitialPriceDecreaseRate),
		paramtypes.NewParamSetPair(KeyPriceDecreaseBlockInterval, &p.PriceDecreaseBlockInterval, validatePriceDecreaseBlockInterval),
		paramtypes.NewParamSetPair(KeyAuctionInterval, &p.AuctionInterval, validateAuctionInterval),
		paramtypes.NewParamSetPair(KeyAuctionThresholdUsdValue, &p.AuctionThresholdUsdValue, validateAuctionThresholdUsdValue),
	}
}

// ValidateBasic performs basic validation on cellarfees parameters.
func (p *Params) ValidateBasic() error {
	if err := validateRewardEmissionPeriod(p.RewardEmissionPeriod); err != nil {
		return err
	}
	if err := validateInitialPriceDecreaseRate(p.InitialPriceDecreaseRate); err != nil {
		return err
	}
	if err := validatePriceDecreaseBlockInterval(p.PriceDecreaseBlockInterval); err != nil {
		return err
	}
	if err := validateAuctionInterval(p.AuctionInterval); err != nil {
		return err
	}
	if err := validateAuctionThresholdUsdValue(p.AuctionThresholdUsdValue); err != nil {
		return err
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

	if rate.LTE(sdk.ZeroDec()) {
		return errorsmod.Wrapf(types.ErrInvalidInitialPriceDecreaseRate, "initial price decrease rate cannot be zero or negative,must be 0 < x < 1")
	}

	if rate.GTE(sdk.OneDec()) {
		return errorsmod.Wrapf(types.ErrInvalidInitialPriceDecreaseRate, "initial price decrease rate cannot be one or greater, must be 0 < x < 1")
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

func validateAuctionThresholdUsdValue(i interface{}) error {
	threshold, ok := i.(sdk.Dec)
	if !ok {
		return errorsmod.Wrapf(types.ErrInvalidAuctionThresholdUsdValue, "auction threshold USD value: %T", i)
	}

	if !threshold.IsPositive() {
		return errorsmod.Wrapf(types.ErrInvalidAuctionThresholdUsdValue, "auction threshold USD value must be greater than zero")
	}

	return nil
}

// String implements the String interface
func (p Params) String() string {
	out, _ := yaml.Marshal(p)
	return string(out)
}
