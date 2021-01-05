package types

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
)

// DefaultParamspace defines default space for oracle params
const DefaultParamspace = ModuleName

// Parameter keys
var (
	ParamStoreKeyVotePeriod               = []byte("voteperiod")
	ParamStoreKeyVoteThreshold            = []byte("votethreshold")
	ParamStoreKeyRewardBand               = []byte("rewardband")
	ParamStoreKeyRewardDistributionWindow = []byte("rewarddistributionwindow")
	ParamStoreKeyWhitelist                = []byte("whitelist")
	ParamStoreKeySlashFraction            = []byte("slashfraction")
	ParamStoreKeySlashWindow              = []byte("slashwindow")
	ParamStoreKeyMinValidPerWindow        = []byte("minvalidperwindow")
)

// Default parameter values
const (
	BlocksPerMinute = int64(10)
	BlocksPerHour   = BlocksPerMinute * 60
	BlocksPerDay    = BlocksPerHour * 24
	BlocksPerWeek   = BlocksPerDay * 7
	BlocksPerMonth  = BlocksPerDay * 30
	BlocksPerYear   = BlocksPerDay * 365

	DefaultVotePeriod               = BlocksPerMinute / 2 // 30 seconds
	DefaultSlashWindow              = BlocksPerWeek       // window for a week
	DefaultRewardDistributionWindow = BlocksPerYear       // window for a year

	MicroLunaDenom = "uluna"
	MicroUSDDenom  = "uusd"
	MicroKRWDenom  = "ukrw"
	MicroSDRDenom  = "usdr"
	MicroCNYDenom  = "ucny"
	MicroJPYDenom  = "ujpy"
	MicroEURDenom  = "ueur"
	MicroGBPDenom  = "ugbp"
	MicroMNTDenom  = "umnt"

	MicroUnit = int64(1e6)
)

// Default parameter values
var (
	DefaultVoteThreshold = sdk.NewDecWithPrec(50, 2) // 50%
	DefaultRewardBand    = sdk.NewDecWithPrec(2, 2)  // 2% (-1, 1)
	DefaultTobinTax      = sdk.NewDecWithPrec(25, 4) // 0.25%
	// TODO: need a better default whitelist
	DefaultWhitelist = sdk.NewDecCoins(
		sdk.NewDecCoinFromDec(MicroKRWDenom, DefaultTobinTax),
		sdk.NewDecCoinFromDec(MicroSDRDenom, DefaultTobinTax),
		sdk.NewDecCoinFromDec(MicroUSDDenom, DefaultTobinTax),
		sdk.NewDecCoinFromDec(MicroMNTDenom, DefaultTobinTax.MulInt64(8)))
	DefaultSlashFraction     = sdk.NewDecWithPrec(1, 4) // 0.01%
	DefaultMinValidPerWindow = sdk.NewDecWithPrec(5, 2) // 5%
)

var _ paramtypes.ParamSet = &Params{}

// DefaultParams creates default oracle module parameters
func DefaultParams() Params {
	return Params{
		VotePeriod:               DefaultVotePeriod,
		VoteThreshold:            DefaultVoteThreshold,
		RewardBand:               DefaultRewardBand,
		RewardDistributionWindow: DefaultRewardDistributionWindow,
		Whitelist:                DefaultWhitelist,
		SlashFraction:            DefaultSlashFraction,
		SlashWindow:              DefaultSlashWindow,
		MinValidPerWindow:        DefaultMinValidPerWindow,
	}
}

// ParamSetPairs returns the parameter set pairs.
func (p *Params) ParamSetPairs() paramtypes.ParamSetPairs {
	return paramtypes.ParamSetPairs{
		paramtypes.NewParamSetPair(ParamStoreKeyVotePeriod, &p.VotePeriod, validateVotePeriod),
		paramtypes.NewParamSetPair(ParamStoreKeyVoteThreshold, &p.VoteThreshold, validateVoteThreshold),
		paramtypes.NewParamSetPair(ParamStoreKeyRewardBand, &p.RewardBand, validateRewardBand),
		paramtypes.NewParamSetPair(ParamStoreKeyRewardDistributionWindow, &p.RewardDistributionWindow, validateRewardDistributionWindow),
		paramtypes.NewParamSetPair(ParamStoreKeyWhitelist, &p.Whitelist, validateWhitelist),
		paramtypes.NewParamSetPair(ParamStoreKeySlashFraction, &p.SlashFraction, validateSlashFraction),
		paramtypes.NewParamSetPair(ParamStoreKeySlashWindow, &p.SlashWindow, validateSlashWindow),
		paramtypes.NewParamSetPair(ParamStoreKeyMinValidPerWindow, &p.MinValidPerWindow, validateMinValidPerWindow),
	}
}

// ParamKeyTable returns the parameter key table.
func ParamKeyTable() paramtypes.KeyTable {
	return paramtypes.NewKeyTable().RegisterParamSet(&Params{})
}

// ValidateBasic performs basic validation on oracle parameters.
func (p Params) ValidateBasic() error {
	if err := validateVotePeriod(p.VotePeriod); err != nil {
		return err
	}

	if err := validateVoteThreshold(p.VoteThreshold); err != nil {
		return err
	}

	if err := validateRewardBand(p.RewardBand); err != nil {
		return err
	}

	if err := validateSlashFraction(p.SlashFraction); err != nil {
		return err
	}

	if err := validateMinValidPerWindow(p.MinValidPerWindow); err != nil {
		return err
	}

	if p.RewardDistributionWindow < p.VotePeriod {
		return fmt.Errorf("oracle RewardDistributionWindow parameter must be ≥ VotePeriod parameter")
	}

	if p.SlashWindow < p.VotePeriod {
		return fmt.Errorf("oracle SlashWindow parameter must be ≥ VotePeriod parameter")
	}

	return validateWhitelist(p.Whitelist)
}

func validateVotePeriod(i interface{}) error {
	v, ok := i.(int64)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v <= 0 {
		return fmt.Errorf("vote period must be positive: %d", v)
	}

	return nil
}

func validateVoteThreshold(i interface{}) error {
	v, ok := i.(sdk.Dec)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v.LT(sdk.NewDecWithPrec(33, 2)) || v.GT(sdk.OneDec()) {
		return fmt.Errorf("reward band must be between [0.33, 1], got %s", v)
	}

	return nil
}

func validateRewardBand(i interface{}) error {
	v, ok := i.(sdk.Dec)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v.IsNegative() || v.GT(sdk.OneDec()) {
		return fmt.Errorf("reward band must be between [0, 1], got %s", v)
	}

	return nil
}

func validateRewardDistributionWindow(i interface{}) error {
	v, ok := i.(int64)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v <= 0 {
		return fmt.Errorf("reward distribution window must be positive: %d", v)
	}

	return nil
}

func validateWhitelist(i interface{}) error {
	v, ok := i.(sdk.DecCoins)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	for _, d := range v {
		if d.Amount.LT(sdk.ZeroDec()) || d.Amount.GT(sdk.OneDec()) {
			return fmt.Errorf("whitelist denom must have TobinTax between [0, 1]")
		}
	}

	return v.Validate()
}

func validateSlashFraction(i interface{}) error {
	v, ok := i.(sdk.Dec)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v.GT(sdk.OneDec()) || v.IsNegative() {
		return fmt.Errorf("slash fraction must be between [0, 1], got %s", v)
	}

	return nil
}

func validateSlashWindow(i interface{}) error {
	v, ok := i.(int64)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v <= 0 {
		return fmt.Errorf("slash window must be positive: %d", v)
	}

	return nil
}

func validateMinValidPerWindow(i interface{}) error {
	v, ok := i.(sdk.Dec)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v.GT(sdk.NewDecWithPrec(5, 1)) || v.IsNegative() {
		return fmt.Errorf("min valid perWindow must be between [0, 0.5], got %s", v)
	}

	return nil
}
