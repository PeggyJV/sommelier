package types

import (
	"errors"
	"fmt"

	bridgetypes "github.com/althea-net/peggy/module/x/peggy/types"
	"github.com/ethereum/go-ethereum/common"

	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
)

// DefaultParamspace defines default space for oracle params
const DefaultParamspace = ModuleName

// Parameter keys
var (
	ParamStoreKeyBatchContractAddress     = []byte("BatchContractAddress")
	ParamStoreKeyLiquidityContractAddress = []byte("LiquidityContractAddress")
	ParamStoreKeyEthTimeoutBlocks         = []byte("EthTimeoutBlocks")
	ParamStoreKeyEthTimeoutTimestamp      = []byte("EthTimeoutBlocksTimestamp")
)

// Default parameter values

var _ paramtypes.ParamSet = &Params{}

// DefaultParams creates default oracle module parameters
func DefaultParams() Params {
	return Params{
		BatchContractAddress:     common.Address{}.String(), // TODO: define
		LiquidityContractAddress: common.Address{}.String(), // TODO: define
		EthTimeoutBlocks:         100,                       // 100 blocks
		EthTimeoutTimestamp:      180,                       // 100 seconds
	}
}

// ParamSetPairs returns the parameter set pairs.
func (p *Params) ParamSetPairs() paramtypes.ParamSetPairs {
	return paramtypes.ParamSetPairs{
		paramtypes.NewParamSetPair(ParamStoreKeyBatchContractAddress, &p.BatchContractAddress, validateContractAddress),
		paramtypes.NewParamSetPair(ParamStoreKeyLiquidityContractAddress, &p.LiquidityContractAddress, validateContractAddress),
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
	if err := validateContractAddress(p.BatchContractAddress); err != nil {
		return fmt.Errorf("invalid batch contract address: %w", err)
	}

	if err := validateContractAddress(p.LiquidityContractAddress); err != nil {
		return fmt.Errorf("invalid liquidity contract address: %w", err)
	}

	if err := validateEthTimeout(p.EthTimeoutBlocks); err != nil {
		return fmt.Errorf("invalid eth timeout blocks height: %w", err)
	}

	if err := validateEthTimeout(p.EthTimeoutTimestamp); err != nil {
		return fmt.Errorf("invalid eth timeout timestamp secods: %w", err)
	}

	return nil
}

func validateContractAddress(i interface{}) error {
	v, ok := i.(string)
	if !ok {
		return fmt.Errorf("invalid contract address parameter type: %T", i)
	}

	return bridgetypes.ValidateEthAddress(v)
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
