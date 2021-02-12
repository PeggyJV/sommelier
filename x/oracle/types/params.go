package types

import (
	"errors"
	fmt "fmt"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
)

// Parameter keys
var (
	KeyVotePeriod        = []byte("voteperiod")
	KeyVoteThreshold     = []byte("votethreshold")
	KeySlashWindow       = []byte("slashwindow")
	KeyMinValidPerWindow = []byte("minvalidperwindow")
	KeySlashFraction     = []byte("slashfraction")
	KeyDataTypes         = []byte("datatypes")
)

var _ paramtypes.ParamSet = &Params{}

// ParamKeyTable returns the parameter key table.
func ParamKeyTable() paramtypes.KeyTable {
	return paramtypes.NewKeyTable().RegisterParamSet(&Params{})
}

// DefaultParams returns default oracle parameters
func DefaultParams() Params {
	return Params{
		VotePeriod:        5,
		VoteThreshold:     sdk.NewDecWithPrec(66, 2), // 66%
		SlashWindow:       10000,
		MinValidPerWindow: sdk.NewDecWithPrec(10, 2), // 10%
		SlashFraction:     sdk.NewDecWithPrec(1, 3),  // 0.1%
		TargetThreshold:   sdk.NewDecWithPrec(5, 3),  // 0.5%,
		DataTypes:         []string{UniswapDataType},
	}
}

// ParamSetPairs returns the parameter set pairs.
func (p *Params) ParamSetPairs() paramtypes.ParamSetPairs {
	return paramtypes.ParamSetPairs{
		paramtypes.NewParamSetPair(KeyVotePeriod, &p.VotePeriod, validateVotePeriod),
		paramtypes.NewParamSetPair(KeyVoteThreshold, &p.VoteThreshold, validateVoteThreshold),
		paramtypes.NewParamSetPair(KeySlashWindow, &p.SlashWindow, validateSlashWindow),
		paramtypes.NewParamSetPair(KeyMinValidPerWindow, &p.MinValidPerWindow, validateMinValidPerWindow),
		paramtypes.NewParamSetPair(KeySlashFraction, &p.SlashFraction, validateSlashFraction),
		paramtypes.NewParamSetPair(KeyDataTypes, &p.DataTypes, validateDataTypes),
	}
}

// ValidateBasic performs basic validation on oracle parameters.
func (p *Params) ValidateBasic() error {
	if err := validateVotePeriod(p.VotePeriod); err != nil {
		return err
	}
	if err := validateVoteThreshold(p.VoteThreshold); err != nil {
		return err
	}
	if err := validateSlashWindow(p.SlashWindow); err != nil {
		return err
	}
	if err := validateMinValidPerWindow(p.MinValidPerWindow); err != nil {
		return err
	}
	if err := validateSlashFraction(p.SlashFraction); err != nil {
		return err
	}
	return validateDataTypes(p.DataTypes)
}

func validateVotePeriod(i interface{}) error {
	votePeriod, ok := i.(int64)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if votePeriod < 4 || votePeriod > 10 {
		return fmt.Errorf(
			"vote period should be between 4 and 10 blocks: %d", votePeriod,
		)
	}

	return nil
}

func validateVoteThreshold(i interface{}) error {
	voteThreshold, ok := i.(sdk.Dec)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if voteThreshold.LTE(sdk.ZeroDec()) || voteThreshold.GT(sdk.OneDec()) {
		return fmt.Errorf("vote threshold value must be within the 0% - 100% range, got: %s", voteThreshold)
	}

	return nil
}

func validateSlashWindow(i interface{}) error {
	slashWindow, ok := i.(int64)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if slashWindow < 1 {
		return fmt.Errorf("slashing window can't be zero or negative: %d", slashWindow)
	}

	return nil
}

func validateMinValidPerWindow(i interface{}) error {
	minValidPerWindow, ok := i.(sdk.Dec)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if minValidPerWindow.LTE(sdk.ZeroDec()) || minValidPerWindow.GT(sdk.OneDec()) {
		return fmt.Errorf("min valid per window value must be within the 0% - 100% range, got: %s", minValidPerWindow)
	}

	return nil
}

func validateSlashFraction(i interface{}) error {
	slashFraction, ok := i.(sdk.Dec)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if slashFraction.LTE(sdk.ZeroDec()) || slashFraction.GT(sdk.OneDec()) {
		return fmt.Errorf("slash fraction value must be within the 0% - 100% range, got: %s", slashFraction)
	}

	return nil
}

func validateDataTypes(i interface{}) error {
	dataTypes, ok := i.([]string)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	for i, dataType := range dataTypes {
		if strings.TrimSpace(dataType) == "" {
			return fmt.Errorf("oracle data type at index %d cannot be blank", i)
		}
	}

	return nil
}
