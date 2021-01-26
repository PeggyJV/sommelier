package types

import (
	"crypto/sha256"
	fmt "fmt"

	"github.com/cosmos/cosmos-sdk/codec"
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
	CannonicalJSON(cdc codec.JSONMarshaler) string
	ValidateBasic() error
	ValidateGroup([]OracleData) error
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

// CannonicalJSON implements OracleData
func (ud *UniswapData) CannonicalJSON(cdc codec.JSONMarshaler) string {
	// TODO: do we need to sort here?
	bz, err := cdc.MarshalJSON(ud)
	if err != nil {
		panic(err)
	}
	return string(bz)
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
