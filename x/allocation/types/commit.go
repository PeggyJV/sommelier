package types

import (
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)


// Validate performs a basic validation on the Oracle vote fields
func (ov OracleVote) Validate() error {
	if ov.Feed == nil || len(ov.Feed.Data) == 0 {
		return sdkerrors.Wrap(ErrInvalidOracleVote, "cannot submit empty oracle data")
	}

	if strings.TrimSpace(ov.Salt) == "" {
		return sdkerrors.Wrap(ErrInvalidOracleVote, "salt string cannot be blank")
	}

	return ov.Feed.Validate()
}

// // UnpackInterfaces implements UnpackInterfacesMessage
// func (ov *OracleVote) UnpackInterfaces(unpacker codectypes.AnyUnpacker) error {
// 	return ov.Feed.UnpackInterfaces(unpacker)
// }

// Validate performs a basic validation on the Oracle feed data fields
func (of OracleFeed) Validate() error {
	var seenIds = make(map[string]bool)

	if len(of.Data) == 0 {
		return sdkerrors.Wrap(ErrInvalidOracleData, "cannot submit empty oracle data")
	}

	for i, od := range of.Data {
		if od == nil {
			return sdkerrors.Wrapf(ErrInvalidOracleData, "oracle data at index %d cannot be nil", i)
		}

		// check if there is a duplicated data entry
		if seenIds[od.GetID()] {
			return sdkerrors.Wrap(ErrDuplicatedOracleData, od.GetID())
		}

		if err := od.Validate(); err != nil {
			return sdkerrors.Wrap(ErrInvalidOracleData, err.Error())
		}

		seenIds[od.GetID()] = true
	}

	return nil
}

// // UnpackInterfaces implements UnpackInterfacesMessage
// func (of *OracleFeed) UnpackInterfaces(unpacker codectypes.AnyUnpacker) error {
// 	for _, oracleDataAny := range of.Data {
// 		if err := unpacker.UnpackAny(oracleDataAny, new(OracleData)); err != nil {
// 			return err
// 		}
// 	}

// 	return nil
// }
