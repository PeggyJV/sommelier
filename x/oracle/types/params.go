package types

import (
	fmt "fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
)

// Parameter keys
var (
	PSKVotePeriod        = []byte("voteperiod")
	PSKVoteThreshold     = []byte("votethreshold")
	PSKSlashWindow       = []byte("slashwindow")
	PSKMinValidPerWindow = []byte("minvalidperwindow")
	PSKSlashFraction     = []byte("slashfraction")
	PSKDataTypes         = []byte("datatypes")
)

// ParamKeyTable returns the parameter key table.
func ParamKeyTable() paramtypes.KeyTable {
	return paramtypes.NewKeyTable().RegisterParamSet(&Params{})
}

// DefaultParams returns default distribution parameters
func DefaultParams() Params {
	return Params{
		VotePeriod:        5,
		VoteThreshold:     sdk.NewDecWithPrec(66, 2), // 66%
		SlashWindow:       10000,
		MinValidPerWindow: sdk.NewDecWithPrec(10, 2), // 10%
		SlashFraction:     sdk.NewDecWithPrec(1, 3),  // 0.1%
		DataTypes:         []string{UniswapDataType},
	}
}

// ParamSetPairs returns the parameter set pairs.
func (p *Params) ParamSetPairs() paramtypes.ParamSetPairs {
	return paramtypes.ParamSetPairs{
		paramtypes.NewParamSetPair(PSKVotePeriod, &p.VotePeriod, validateVotePeriod),
		paramtypes.NewParamSetPair(PSKVoteThreshold, &p.VoteThreshold, validateVoteThreshold),
		paramtypes.NewParamSetPair(PSKSlashWindow, &p.SlashWindow, validateSlashWindow),
		paramtypes.NewParamSetPair(PSKMinValidPerWindow, &p.MinValidPerWindow, validateMinValidPerWindow),
		paramtypes.NewParamSetPair(PSKSlashFraction, &p.SlashFraction, validateSlashFraction),
		paramtypes.NewParamSetPair(PSKDataTypes, &p.DataTypes, validateDataTypes),
	}
}

// ValidateBasic performs basic validation on distribution parameters.
func (p *Params) ValidateBasic() error {
	if p.VotePeriod < 4 || p.VotePeriod > 10 {
		return fmt.Errorf(
			"vote period should be between 4 and 10 blocks: %d", p.VotePeriod,
		)
	}

	return nil
}

func validateVotePeriod(i interface{}) error {
	if _, ok := i.(int64); !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}
	return nil
}

func validateVoteThreshold(i interface{}) error {
	if _, ok := i.(sdk.Dec); !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}
	return nil
}

func validateSlashWindow(i interface{}) error {
	if _, ok := i.(int64); !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}
	return nil
}

func validateMinValidPerWindow(i interface{}) error {
	if _, ok := i.(sdk.Dec); !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}
	return nil
}

func validateSlashFraction(i interface{}) error {
	if _, ok := i.(sdk.Dec); !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}
	return nil
}

func validateDataTypes(i interface{}) error {
	if _, ok := i.([]string); !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}
	return nil
}
