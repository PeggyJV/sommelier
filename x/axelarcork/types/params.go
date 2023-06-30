package types

import (
	"errors"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"

	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
)

// Parameter keys
var (
	KeyIBCChannel      = []byte("ibcchannel")
	KeyIBCPort         = []byte("ibcport")
	KeyGMPAccount      = []byte("gmpaccount")
	KeyExecutorAccount = []byte("executoraccount")
	KeyTimeoutDuration = []byte("timeoutduration")
)

var _ paramtypes.ParamSet = &Params{}

// ParamKeyTable returns the parameter key table.
func ParamKeyTable() paramtypes.KeyTable {
	return paramtypes.NewKeyTable().RegisterParamSet(&Params{})
}

// DefaultParams returns default oracle parameters
func DefaultParams() Params {
	return Params{
		IbcChannel:      "",
		IbcPort:         "",
		GmpAccount:      "",
		ExecutorAccount: "",
		TimeoutDuration: 0,
	}
}

// ParamSetPairs returns the parameter set pairs.
func (p *Params) ParamSetPairs() paramtypes.ParamSetPairs {
	return paramtypes.ParamSetPairs{
		paramtypes.NewParamSetPair(KeyIBCChannel, &p.IbcChannel, validateIBCChannel),
		paramtypes.NewParamSetPair(KeyIBCPort, &p.IbcPort, validateIBCPort),
		paramtypes.NewParamSetPair(KeyGMPAccount, &p.GmpAccount, validateGMPAccount),
		paramtypes.NewParamSetPair(KeyExecutorAccount, &p.ExecutorAccount, validateExecutorAccount),
		paramtypes.NewParamSetPair(KeyTimeoutDuration, &p.TimeoutDuration, validateTimeoutDuration),
	}
}

// ValidateBasic performs basic validation on oracle parameters.
func (p *Params) ValidateBasic() error {
	if err := validateIBCChannel(p.IbcChannel); err != nil {
		return err
	}
	if err := validateIBCPort(p.IbcPort); err != nil {
		return err
	}
	if err := validateGMPAccount(p.GmpAccount); err != nil {
		return err
	}
	if err := validateExecutorAccount(p.ExecutorAccount); err != nil {
		return err
	}
	if err := validateTimeoutDuration(p.TimeoutDuration); err != nil {
		return err
	}
	return nil
}

func validateIBCChannel(i interface{}) error {
	ibcChannel, ok := i.(string)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if ibcChannel == "" {
		return errors.New("ibc channel cannot be empty")
	}

	return nil
}

func validateIBCPort(i interface{}) error {
	ibcPort, ok := i.(string)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if ibcPort == "" {
		return errors.New("ibc port cannot be empty")
	}

	return nil
}

func validateGMPAccount(i interface{}) error {
	gmpAcc, ok := i.(string)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if gmpAcc == "" {
		return errors.New("gmp account cannot be empty")
	}

	if _, err := sdk.AccAddressFromBech32(gmpAcc); err != nil {
		return err
	}

	return nil
}

func validateExecutorAccount(i interface{}) error {
	execAcc, ok := i.(string)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if execAcc == "" {
		return errors.New("executor account cannot be empty")
	}

	return nil
}

func validateTimeoutDuration(i interface{}) error {
	timeout, ok := i.(uint64)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if timeout == 0 {
		return errors.New("timeout duration cannot be zero")
	}

	return nil
}
