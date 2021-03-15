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
	ParamStoreKeyContractAddress     = []byte("ContractAddress")
	ParamStoreKeyEthTimeoutBlocks    = []byte("EthTimeoutBlocks")
	ParamStoreKeyEthTimeoutTimestamp = []byte("EthTimeoutBlocksTimestamp")
)

// Default parameter values

var _ paramtypes.ParamSet = &Params{}

// DefaultParams creates default oracle module parameters
func DefaultParams() Params {
	return Params{
		ContractAddress:     common.Address{}.String(), // TODO: define
		EthTimeoutBlocks:    100,                       // 100 blocks
		EthTimeoutTimestamp: 100,                       // 100 seconds
	}
}

// ParamSetPairs returns the parameter set pairs.
func (p *Params) ParamSetPairs() paramtypes.ParamSetPairs {
	return paramtypes.ParamSetPairs{
		paramtypes.NewParamSetPair(ParamStoreKeyContractAddress, &p.ContractAddress, validateContractAddress),
		paramtypes.NewParamSetPair(ParamStoreKeyEthTimeoutBlocks, &p.EthTimeoutBlocks, validateEthTimeout),
		paramtypes.NewParamSetPair(ParamStoreKeyEthTimeoutTimestamp, &p.EthTimeoutTimestamp, validateEthTimeout),
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

	if v < 1 {
		return errors.New("eth timeout value cannot less than 1")
	}

	return nil
}
