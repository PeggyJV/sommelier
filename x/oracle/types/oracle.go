package types

import (
	"crypto/sha256"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	proto "github.com/gogo/protobuf/proto"

	sdk "github.com/cosmos/cosmos-sdk/types"

	peggytypes "github.com/althea-net/peggy/module/x/peggy/types"
)

var _ OracleData = &UniswapData{}

// UniswapDataType defines the data type for a uniswap pair oracle data
const UniswapDataType = "uniswap"

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

// GetID implements OracleData
func (ud UniswapData) GetID() string {
	return ud.Id
}

// Type implements OracleData
func (ud *UniswapData) Type() string {
	return UniswapDataType
}

// Validate implements OracleData
func (ud UniswapData) Validate() error {
	if err := peggytypes.ValidateEthAddress(ud.Id); err != nil {
		return fmt.Errorf("invalid uniswap pair id %s: %w", ud.Id, err)
	}

	if ud.Reserve0.IsNegative() {
		return fmt.Errorf("reserve 0 value (%s) for uniswap pair %s cannot be negative", ud.Reserve0, ud.Id)
	}

	if ud.Reserve1.IsNegative() {
		return fmt.Errorf("reserve 1 value (%s) for uniswap pair %s cannot be negative", ud.Reserve0, ud.Id)
	}

	if ud.ReserveUsd.IsNegative() {
		return fmt.Errorf("reserve USD value (%s) for uniswap pair %s cannot be negative", ud.Reserve0, ud.Id)
	}

	if err := ud.Token0.Validate(); err != nil {
		return fmt.Errorf("invalid token 0 for uniswap pair %s: %w", ud.Id, err)
	}

	if err := ud.Token1.Validate(); err != nil {
		return fmt.Errorf("invalid token 1 for uniswap pair %s: %w", ud.Id, err)
	}

	if ud.Token0Price.IsNegative() {
		return fmt.Errorf("token 0 price (%s) for uniswap pair %s cannot be negative", ud.Token0Price, ud.Id)
	}

	if ud.Token1Price.IsNegative() {
		return fmt.Errorf("token 1 price (%s) for uniswap pair %s cannot be negative", ud.Token1Price, ud.Id)
	}

	if ud.TotalSupply.IsNegative() {
		return fmt.Errorf("total supply (%s) for uniswap pair %s cannot be negative", ud.TotalSupply, ud.Id)
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
func (ud UniswapData) MarshalJSON() ([]byte, error) {
	bz, err := json.Marshal(ud)
	if err != nil {
		return nil, err
	}

	return sdk.SortJSON(bz)
}

// BlocksTillNextPeriod helper
func (vp *VotePeriod) BlocksTillNextPeriod() int64 {
	return vp.VotePeriodEnd - vp.CurrentHeight
}
