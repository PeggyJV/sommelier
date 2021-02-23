package types

import (
	"crypto/sha256"
	"errors"
	"fmt"

	tmbytes "github.com/tendermint/tendermint/libs/bytes"

	sdk "github.com/cosmos/cosmos-sdk/types"
	proto "github.com/gogo/protobuf/proto"

	peggytypes "github.com/althea-net/peggy/module/x/peggy/types"
)

var _ OracleData = &UniswapPair{}

// UniswapDataType defines the data type for a uniswap pair oracle data
const UniswapDataType = "uniswap"

// DataHashes defines an array of bytes in hex format.
type DataHashes []tmbytes.HexBytes

// OracleData represents a data type that is supported by the oracle
type OracleData interface {
	proto.Message

	// GetID returns the identifier of the data
	GetID() string
	// Type returns the oracle type category
	Type() string
	// Validate performs a stateless validation of the fields
	Validate() error
	// Compare checks that the data is within the target range and the fixed
	// fields match with the aggregated data
	Compare(aggregatedData OracleData, target sdk.Dec) bool
}

// DataHash returns the hash for a precommit given the proper args
func DataHash(salt, jsonData string, signer sdk.ValAddress) []byte {
	h := sha256.New()
	h.Write([]byte(fmt.Sprintf("%s:%s:%s", salt, jsonData, signer.String())))
	return h.Sum(nil)
}

// NewUniswapPair creates a new UniswapPair instance with the fixed values set by args and
// the other fields to their zero values.
func NewUniswapPair(id string, token0, token1 UniswapToken) *UniswapPair {
	return &UniswapPair{
		Id:          id,
		Reserve0:    sdk.ZeroDec(),
		Reserve1:    sdk.ZeroDec(),
		ReserveUsd:  sdk.ZeroDec(),
		Token0:      token0,
		Token1:      token1,
		Token0Price: sdk.ZeroDec(),
		Token1Price: sdk.ZeroDec(),
		TotalSupply: sdk.ZeroDec(),
	}
}

// GetID implements OracleData
func (up UniswapPair) GetID() string {
	return up.Id
}

// Type implements OracleData
func (up *UniswapPair) Type() string {
	return UniswapDataType
}

// Validate implements OracleData
func (up UniswapPair) Validate() error {
	if err := peggytypes.ValidateEthAddress(up.Id); err != nil {
		return fmt.Errorf("invalid uniswap pair id %s: %w", up.Id, err)
	}

	if up.Reserve0.IsNil() {
		return errors.New("reserve 0 cannot be nil")
	}

	if up.Reserve1.IsNil() {
		return errors.New("reserve 1 cannot be nil")
	}

	if up.ReserveUsd.IsNil() {
		return errors.New("reserve USD cannot be nil")
	}

	if up.Token0Price.IsNil() {
		return errors.New("token 0 price cannot be nil")
	}

	if up.Token1Price.IsNil() {
		return errors.New("token 0 price  cannot be nil")
	}

	if up.TotalSupply.IsNil() {
		return errors.New("token supply cannot be nil")
	}

	if up.Reserve0.IsNegative() {
		return fmt.Errorf("reserve 0 value (%s) for uniswap pair %s cannot be negative", up.Reserve0, up.Id)
	}

	if up.Reserve1.IsNegative() {
		return fmt.Errorf("reserve 1 value (%s) for uniswap pair %s cannot be negative", up.Reserve0, up.Id)
	}

	if up.ReserveUsd.IsNegative() {
		return fmt.Errorf("reserve USD value (%s) for uniswap pair %s cannot be negative", up.Reserve0, up.Id)
	}

	if err := up.Token0.Validate(); err != nil {
		return fmt.Errorf("invalid token 0 for uniswap pair %s: %w", up.Id, err)
	}

	if err := up.Token1.Validate(); err != nil {
		return fmt.Errorf("invalid token 1 for uniswap pair %s: %w", up.Id, err)
	}

	if up.Token0Price.IsNegative() {
		return fmt.Errorf("token 0 price (%s) for uniswap pair %s cannot be negative", up.Token0Price, up.Id)
	}

	if up.Token1Price.IsNegative() {
		return fmt.Errorf("token 1 price (%s) for uniswap pair %s cannot be negative", up.Token1Price, up.Id)
	}

	if up.TotalSupply.IsNegative() {
		return fmt.Errorf("total supply (%s) for uniswap pair %s cannot be negative", up.TotalSupply, up.Id)
	}

	return nil
}

// Compare checks that the current pair is within the target range and the fixed
// fields match with the aggregated data.
func (up UniswapPair) Compare(aggregatedData OracleData, target sdk.Dec) bool {
	aggregatedPair, ok := aggregatedData.(*UniswapPair)
	if !ok || up.Type() != aggregatedData.Type() {
		// aggregated data is from a different type
		return false
	}

	// fixed data must match: id, token 0, token 1

	if up.Id != aggregatedPair.Id {
		return false
	}

	if up.Token0 != aggregatedPair.Token0 {
		return false
	}

	if up.Token1 != aggregatedPair.Token1 {
		return false
	}

	// |reserve0 - reserve0 (agg)| / (reserve0 (agg)) ≤ target
	if up.Reserve0.Sub(aggregatedPair.Reserve0).Abs().GT(target.Mul(aggregatedPair.Reserve0)) {
		return false
	}

	// |reserve1 - reserve1 (agg)| / (reserve1 (agg)) ≤ target
	if up.Reserve1.Sub(aggregatedPair.Reserve1).Abs().Quo(aggregatedPair.Reserve1).GT(target.Mul(aggregatedPair.Reserve1)) {
		return false
	}

	// |reserveUsd - reserveUsd (agg)| / (reserveUsd (agg)) ≤ target
	if up.ReserveUsd.Sub(aggregatedPair.ReserveUsd).Abs().GT(target.Mul(aggregatedPair.ReserveUsd)) {
		return false
	}

	// |token0price - token0price (agg)| / (token0price (agg)) ≤ target
	if up.Token0Price.Sub(aggregatedPair.Token0Price).Abs().GT(target.Mul(aggregatedPair.Token0Price)) {
		return false
	}

	// |token1price - token1price (agg)| / (token1price (agg)) ≤ target
	if up.Token1Price.Sub(aggregatedPair.Token1Price).Abs().GT(target.Mul(aggregatedPair.Token1Price)) {
		return false
	}

	// |totalSupply - totalSupply (agg)| / (totalSupply (agg)) ≤ target
	if up.TotalSupply.Sub(aggregatedPair.TotalSupply).Abs().GT(target.Mul(aggregatedPair.TotalSupply)) {
		return false
	}

	return true
}

// Validate performs a basic validation of the uniswap token fields.
func (ut UniswapToken) Validate() error {
	if err := peggytypes.ValidateEthAddress(ut.Id); err != nil {
		return fmt.Errorf("invalid token address %s: %w", ut.Id, err)
	}

	return nil
}
