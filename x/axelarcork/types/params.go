package types

import (
	"errors"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	host "github.com/cosmos/ibc-go/v7/modules/core/24-host"

	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
)

// Parameter keys
var (
	KeyEnabled           = []byte("enabled")
	KeyIBCChannel        = []byte("ibcchannel")
	KeyIBCPort           = []byte("ibcport")
	KeyGMPAccount        = []byte("gmpaccount")
	KeyExecutorAccount   = []byte("executoraccount")
	KeyTimeoutDuration   = []byte("timeoutduration")
	KeyCorkTimeoutBlocks = []byte("corktimeoutblocks")
	KeyCorkAuthority     = []byte("corkauthority")
)

var _ paramtypes.ParamSet = &Params{}

// ParamKeyTable returns the parameter key table.
func ParamKeyTable() paramtypes.KeyTable {
	return paramtypes.NewKeyTable().RegisterParamSet(&Params{})
}

// DefaultParams returns default oracle parameters
func DefaultParams() Params {
	return Params{
		Enabled:           false,
		IbcChannel:        "",
		IbcPort:           "",
		GmpAccount:        "",
		ExecutorAccount:   "",
		TimeoutDuration:   0,
		CorkTimeoutBlocks: 10000,
		CorkAuthority:     "",
	}
}

// ParamSetPairs returns the parameter set pairs.
func (p *Params) ParamSetPairs() paramtypes.ParamSetPairs {
	return paramtypes.ParamSetPairs{
		paramtypes.NewParamSetPair(KeyEnabled, &p.Enabled, validateEnabled),
		paramtypes.NewParamSetPair(KeyIBCChannel, &p.IbcChannel, validateIBCChannel),
		paramtypes.NewParamSetPair(KeyIBCPort, &p.IbcPort, validateIBCPort),
		paramtypes.NewParamSetPair(KeyGMPAccount, &p.GmpAccount, validateGMPAccount),
		paramtypes.NewParamSetPair(KeyExecutorAccount, &p.ExecutorAccount, validateExecutorAccount),
		paramtypes.NewParamSetPair(KeyTimeoutDuration, &p.TimeoutDuration, validateTimeoutDuration),
		paramtypes.NewParamSetPair(KeyCorkTimeoutBlocks, &p.CorkTimeoutBlocks, validateCorkTimeoutBlocks),
		paramtypes.NewParamSetPair(KeyCorkAuthority, &p.CorkAuthority, validateCorkAuthority),
	}
}

// ValidateBasic performs basic validation on oracle parameters.
func (p *Params) ValidateBasic() error {
	if err := validateEnabled(p.Enabled); err != nil {
		return err
	}

	if p.Enabled {
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
		if err := validateCorkTimeoutBlocks(p.CorkTimeoutBlocks); err != nil {
			return err
		}
	}

	// Validated unconditionally, outside the Enabled guard: a malformed
	// authority address is invalid whether or not the module is enabled, and
	// leaving it unchecked while disabled would let a bad value be staged and
	// only surface when the module is turned on.
	if err := validateCorkAuthority(p.CorkAuthority); err != nil {
		return err
	}

	return nil
}

func validateEnabled(i interface{}) error {
	_, ok := i.(bool)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	return nil
}

func validateIBCChannel(i interface{}) error {
	ibcChannel, ok := i.(string)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if err := host.ChannelIdentifierValidator(ibcChannel); err != nil {
		return err
	}

	return nil
}

func validateIBCPort(i interface{}) error {
	ibcPort, ok := i.(string)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if err := host.PortIdentifierValidator(ibcPort); err != nil {
		return err
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

	if _, err := sdk.GetFromBech32(gmpAcc, "axelar"); err != nil {
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

	if _, err := sdk.GetFromBech32(execAcc, "axelar"); err != nil {
		return err
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

func validateCorkTimeoutBlocks(i interface{}) error {
	timeout, ok := i.(uint64)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if timeout == 0 {
		return errors.New("timeout blocks cannot be zero")
	}

	return nil
}

// validateCorkAuthority accepts the empty string (fail-closed: no address may
// act) or a well-formed somm1 account address. Enforcement of the non-empty
// requirement lives in the msg server, so that governance can revoke the
// authority by setting it empty.
func validateCorkAuthority(i interface{}) error {
	authority, ok := i.(string)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}
	if authority == "" {
		return nil
	}
	if _, err := sdk.AccAddressFromBech32(authority); err != nil {
		return fmt.Errorf("invalid cork authority address %q: %w", authority, err)
	}
	return nil
}
