package types

import (
	"errors"
	"fmt"

	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
	"github.com/ethereum/go-ethereum/common"
)

// DefaultParamspace defines default space for oracle params
const DefaultParamspace = ModuleName

// Parameter keys
var (
	ParamStoreKeyContractAddress = []byte("ContractAddress")
	ParamStoreKeyEthTimeout      = []byte("EthTimeout")
)

// Default parameter values

var _ paramtypes.ParamSet = &Params{}

// DefaultParams creates default oracle module parameters
func DefaultParams() Params {
	return Params{
		ContractAddress: common.Address{}.String(), // TODO: define
		EthTimeout:      100,
	}
}

// ParamSetPairs returns the parameter set pairs.
func (p *Params) ParamSetPairs() paramtypes.ParamSetPairs {
	return paramtypes.ParamSetPairs{
		paramtypes.NewParamSetPair(ParamStoreKeyContractAddress, &p.ContractAddress, validateContractAddress),
		paramtypes.NewParamSetPair(ParamStoreKeyEthTimeout, &p.EthTimeout, validateEthTimeout),
	}
}

// ParamKeyTable returns the parameter key table.
func ParamKeyTable() paramtypes.KeyTable {
	return paramtypes.NewKeyTable().RegisterParamSet(&Params{})
}

// ValidateBasic performs basic validation on oracle parameters.
func (p Params) ValidateBasic() error {
	return validateContractAddress(p.ContractAddress)
}

func validateContractAddress(i interface{}) error {
	v, ok := i.(string)
	if !ok {
		return fmt.Errorf("invalid contract address parameter type: %T", i)
	}

	if v == "" {
		return errors.New("contract address cannot be blank or zero address")
	}

	return nil
}

func validateEthTimeout(i interface{}) error {
	v, ok := i.(uint64)
	if !ok {
		return fmt.Errorf("invalid ETH timeout parameter type: %T", i)
	}

	if v < 10 {
		return errors.New("eth timeout value cannot less than 10")
	}

	return nil
}
