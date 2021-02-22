package types

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"

	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// FeederVote defines a single oracle data submitted by the feeder address
type FeederVote struct {
	Data    OracleData
	Address sdk.ValAddress
}

// BlocksTillNextPeriod helper
func (vp *VotePeriod) BlocksTillNextPeriod() int64 {
	return vp.VotePeriodEnd - vp.CurrentHeight
}

var (
	_ codectypes.UnpackInterfacesMessage = &OracleFeed{}
	_ codectypes.UnpackInterfacesMessage = &OracleVote{}
)

// Validate performs a basic validation on the Oracle vote fields
func (ov OracleVote) Validate() error {
	if ov.Feed == nil || len(ov.Feed.OracleData) == 0 {
		return sdkerrors.Wrap(ErrInvalidOracleData, "cannot submit empty oracle data")
	}

	if len(ov.Salt) != len(ov.Feed.OracleData) {
		return sdkerrors.Wrapf(ErrInvalidOracleData, "must match salt array length, expected %d, got %d", len(ov.Salt), len(ov.Feed.OracleData))
	}

	for i, salt := range ov.Salt {
		if strings.TrimSpace(salt) == "" {
			return fmt.Errorf("salt string at index %d cannot be blank", i)
		}
	}

	return ov.Feed.Validate()
}

// UnpackInterfaces implements UnpackInterfacesMessage
func (ov *OracleVote) UnpackInterfaces(unpacker codectypes.AnyUnpacker) error {
	return ov.Feed.UnpackInterfaces(unpacker)
}

// MarshalJSON marshals and sorts the returned value
func (of OracleFeed) MarshalJSON() ([]byte, error) {
	bz, err := json.Marshal(of)
	if err != nil {
		return nil, err
	}

	return sdk.SortJSON(bz)
}

// Validate performs a basic validation on the Oracle feed data fields
func (of OracleFeed) Validate() error {
	// NOTE: oracle data from the array MUST have the same type
	var (
		oracleDataType reflect.Type
		seenIds        = make(map[string]bool)
	)

	for i, oracleData := range of.OracleData {
		od, err := UnpackOracleData(oracleData)
		if err != nil {
			return sdkerrors.Wrap(ErrInvalidOracleData, err.Error())
		}

		// check type consistency
		dataType := reflect.TypeOf(od)

		if i == 0 {
			oracleDataType = dataType
		} else if oracleDataType != dataType {
			return sdkerrors.Wrapf(ErrInvalidOracleData, "oracle data type mismatch, expected %v, got %v", oracleDataType, dataType)
		}

		// check if there is a duplicated data entry
		if seenIds[od.GetID()] {
			return sdkerrors.Wrap(ErrDuplicatedOracleData, od.GetID())
		}

		if err = od.Validate(); err != nil {
			return sdkerrors.Wrap(ErrInvalidOracleData, err.Error())
		}

		seenIds[od.GetID()] = true
	}

	return nil
}

// UnpackInterfaces implements UnpackInterfacesMessage
func (of *OracleFeed) UnpackInterfaces(unpacker codectypes.AnyUnpacker) error {
	for _, oracleDataAny := range of.OracleData {
		var od OracleData
		if err := unpacker.UnpackAny(oracleDataAny, &od); err != nil {
			return err
		}
	}

	return nil
}
