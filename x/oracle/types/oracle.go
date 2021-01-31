package types

import (
	"crypto/sha256"
	"encoding/json"
	fmt "fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

var (
	_               OracleData = &UniswapData{}
	UniswapDataType            = "uniswap_data"
)

// I think we will need the below interface
// type OracleDataCollection interface {
// 	ValidateMember(OracleData) error
// }

// OracleData represents a data type that is supported by the oracle
type OracleData interface {
	CannonicalJSON() string
	ValidateBasic() error
	ValidateGroup([]OracleData) error
	Valid(OracleData) bool
	Type() string

	// TODO: figure out if we need a sorting function here
	// Sort() error
}

// DataHash returns the hash for a precommit given the proper args
func DataHash(salt string, jsn string, signer sdk.AccAddress) []byte {
	h := sha256.New()
	h.Write([]byte(fmt.Sprintf("%s:%s:%s", salt, jsn, signer.String())))
	return h.Sum(nil)
}

// GetAverageFunction registers the collection functions for each of the data types
func GetAverageFunction(typ string) func([]OracleData) OracleData {
	switch typ {
	case UniswapDataType:
		return UniswapDataCollection
	default:
		return nil
	}
}

// UniswapDataCollection averages a collection of uniswap data
func UniswapDataCollection(uds []OracleData) OracleData {
	// TODO: for now just return the first instance
	return uds[0]
}

// Valid uses a cannonical UniswapData instance to validate instances passed in
func (ud *UniswapData) Valid(ud1 OracleData) bool {
	v, ok := ud1.(*UniswapData)
	if !ok {
		return false
	}
	if len(v.Pairs) != len(ud.Pairs) {
		return false
	}
	// TODO: validate more things here
	return true
}

// CannonicalJSON implements OracleData
func (ud *UniswapData) CannonicalJSON() string {
	bz, err := json.Marshal(ud)
	if err != nil {
		panic(err)
	}
	return string(sdk.MustSortJSON(bz))
}

// ValidateBasic implements OracleData
func (ud *UniswapData) ValidateBasic() error {
	// if len(ud.Data) != 1000 {
	// 	return fmt.Errorf("Must input 1000 markets")
	// }
	// TODO: other basic validation
	return nil
}

// ValidateGroup implements OracleData
func (ud *UniswapData) ValidateGroup(unds []OracleData) error {
	// TODO: Ensure that []OracleData is []UniswapData
	// TODO: figure out what metrics an individual vote needs to hit
	// in order to be considered valid
	return nil
}

// Type implements OracleData
func (ud *UniswapData) Type() string {
	return UniswapDataType
}
