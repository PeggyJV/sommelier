package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
	"gopkg.in/yaml.v2"
)

const (
	DefaultAuctionBlockDelay          uint64 = 201600
	DefaultRewardEmissionPeriod       uint64 = 403200
	DefaultInitialPriceDecreaseRate   uint64 = 347000000000000
	DefaultPriceDecreaseBlockInterval uint64 = 10
)

// Parameter keys
var (
	// Rough number of blocks in 2 weeks, or ~2 fee accrual cycles for one cellar
	KeyAuctionBlockDelay = []byte("AuctionBlockDelay")
	// Rough number of blocks in 28 days, the time it takes to unbond
	KeyRewardEmissionPeriod = []byte("RewardEmissionPeriod")
	// Initial rate at which an auction should decrease the price of the relevant coin from it's starting price
	KeyInitialPriceDecreaseRate = []byte("InitialPriceDecreaseRate")
	// Blocks between each auction price decrease
	KeyPriceDecreaseBlockInterval = []byte("PriceDecreaseBlockInterval")
)

var _ paramtypes.ParamSet = &Params{}

// ParamKeyTable returns the parameter key table.
func ParamKeyTable() paramtypes.KeyTable {
	return paramtypes.NewKeyTable().RegisterParamSet(&Params{})
}

// DefaultParams returns default cellarfees parameters
func DefaultParams() Params {

	return Params{
		AuctionBlockDelay:          DefaultAuctionBlockDelay,
		RewardEmissionPeriod:       DefaultRewardEmissionPeriod,
		InitialPriceDecreaseRate:   sdk.NewDec(int64(DefaultInitialPriceDecreaseRate)),
		PriceDecreaseBlockInterval: DefaultPriceDecreaseBlockInterval,
	}
}

// ParamSetPairs returns the parameter set pairs.
func (p *Params) ParamSetPairs() paramtypes.ParamSetPairs {
	return paramtypes.ParamSetPairs{
		paramtypes.NewParamSetPair(KeyAuctionBlockDelay, &p.AuctionBlockDelay, validateAuctionBlockDelay),
		paramtypes.NewParamSetPair(KeyRewardEmissionPeriod, &p.RewardEmissionPeriod, validateRewardEmissionPeriod),
		paramtypes.NewParamSetPair(KeyInitialPriceDecreaseRate, &p.InitialPriceDecreaseRate, validateInitialPriceDecreaseRate),
		paramtypes.NewParamSetPair(KeyPriceDecreaseBlockInterval, &p.PriceDecreaseBlockInterval, validatePriceDecreaseBlockInterval),
	}
}

// ValidateBasic performs basic validation on cellarfees parameters.
func (p *Params) ValidateBasic() error {
	if err := validateAuctionBlockDelay(p.AuctionBlockDelay); err != nil {
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

func validateAuctionBlockDelay(i interface{}) error {
	blockDelay, ok := i.(uint64)
	if !ok {
		return sdkerrors.Wrapf(ErrInvalidAuctionBlockDelay, "auction block delay: %T", i)
	}

	if blockDelay == 0 {
		return sdkerrors.Wrapf(ErrInvalidAuctionBlockDelay, "auction block delay cannot be zero")
	}

	return nil
}

func validateRewardEmissionPeriod(i interface{}) error {
	emissionPeriod, ok := i.(uint64)
	if !ok {
		return sdkerrors.Wrapf(ErrInvalidRewardEmissionPeriod, "reward emission period: %T", i)
	}

	if emissionPeriod == 0 {
		return sdkerrors.Wrapf(ErrInvalidRewardEmissionPeriod, "reward emission period cannot be zero")
	}

	return nil
}

func validateInitialPriceDecreaseRate(i interface{}) error {
	rate, ok := i.(sdk.Dec)
	if !ok {
		return sdkerrors.Wrapf(ErrInvalidInitialPriceDecreaseRate, "initial price decrease rate: %T", i)
	}

	if rate == sdk.ZeroDec() {
		return sdkerrors.Wrapf(ErrInvalidInitialPriceDecreaseRate, "initial price decrease rate cannot be zero, must be 0 < x < 1")
	}

	if rate == sdk.OneDec() {
		return sdkerrors.Wrapf(ErrInvalidInitialPriceDecreaseRate, "initial price decrease rate cannot be one, must be 0 < x < 1")
	}

	return nil
}

func validatePriceDecreaseBlockInterval(i interface{}) error {
	interval, ok := i.(uint64)
	if !ok {
		return sdkerrors.Wrapf(ErrInvalidPriceDecreaseBlockInterval, "price decrease block interval: %T", i)
	}

	if interval == 0 {
		return sdkerrors.Wrapf(ErrInvalidPriceDecreaseBlockInterval, "price decrease block interval cannot be zero")
	}

	return nil
}

// String implements the String interface
func (p Params) String() string {
	out, _ := yaml.Marshal(p)
	return string(out)
}
