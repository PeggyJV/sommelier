package types

import (
	"crypto/sha256"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

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

	GetID() string
	Type() string
	Validate() error
	MarshalJSON() ([]byte, error)
}

// DataHash returns the hash for a precommit given the proper args
func DataHash(salt, jsonData string, signer sdk.AccAddress) []byte {
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

// Validate performs a basic validation of the uniswap token fields.
func (ut UniswapToken) Validate() error {
	if strings.TrimSpace(ut.Id) == "" {
		return errors.New("token id cannot be blank")
	}

	return nil
}

// MarshalJSON marshals and sorts the returned value
func (up UniswapPair) MarshalJSON() ([]byte, error) {
	bz, err := json.Marshal(up)
	if err != nil {
		return nil, err
	}

	return sdk.SortJSON(bz)
}

// BlocksTillNextPeriod helper
func (vp *VotePeriod) BlocksTillNextPeriod() int64 {
	return vp.VotePeriodEnd - vp.CurrentHeight
}
