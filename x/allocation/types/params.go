package types

import (
	"errors"
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
)

// Parameter keys
var (
	KeyVotePeriod    = []byte("voteperiod")
	KeyVoteThreshold = []byte("votethreshold")
)

var _ paramtypes.ParamSet = &Params{}

// ParamKeyTable returns the parameter key table.
func ParamKeyTable() paramtypes.KeyTable {
	return paramtypes.NewKeyTable().RegisterParamSet(&Params{})
}

// DefaultParams returns default oracle parameters
func DefaultParams() Params {
	return Params{
		VotePeriod:    5,
		VoteThreshold: sdk.NewDecWithPrec(67, 2), // 67%
	}
}

// ParamSetPairs returns the parameter set pairs.
func (p *Params) ParamSetPairs() paramtypes.ParamSetPairs {
	return paramtypes.ParamSetPairs{
		paramtypes.NewParamSetPair(KeyVotePeriod, &p.VotePeriod, validateVotePeriod),
		paramtypes.NewParamSetPair(KeyVoteThreshold, &p.VoteThreshold, validateVoteThreshold),
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
	return nil
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

	if voteThreshold.IsNil() {
		return errors.New("vote threshold cannot be nil")
	}

	if voteThreshold.LTE(sdk.ZeroDec()) || voteThreshold.GT(sdk.OneDec()) {
		return fmt.Errorf("vote threshold value must be within the 0% - 100% range, got: %s", voteThreshold)
	}

	return nil
}
