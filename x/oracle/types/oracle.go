package types

import (
	"crypto/sha256"
	"encoding/json"
	fmt "fmt"
	"strconv"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

var (
	_               OracleData = &UniswapData{}
	UniswapDataType            = "uniswap_data"
)

// OracleData represents a data type that is supported by the oracle
type OracleData interface {
	CannonicalJSON() string
	ValidateBasic() error
	Valid(OracleData) bool
	Type() string
	// TODO: Add a parsing function to this and figure out what the signature needs to be
	// Parse()
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
	if len(uds) == 1 {
		return uds[0]
	}
	foo := [][]UniswapPairParsed{}
	for _, od := range uds {
		if ud, ok := od.(*UniswapData); ok {
			parsed, err := ud.Parse()
			if err != nil {
				// TODO: this is validated in the message handler
				panic(err)
			}
			foo = append(foo, parsed)
		}
	}
	return uds[0]
}

// FlattenParsedUniswapData takes an array of arrays of pairs and flattens them
func FlattenParsedUniswapData(in [][]UniswapPairParsed) (out []UniswapPairParsed) {
	byID := make(map[string][]UniswapPairParsed)
	for _, pairs := range in {
		for _, pair := range pairs {
			byID[pair.ID] = append(byID[pair.ID], pair)
			// TODO: might want to validate that the tokens in each ID are the same
			// TODO: also validate that total supply is the same
		}
	}
	for id, pairs := range byID {
		ppout := UniswapPairParsed{ID: id}
		for _, pair := range pairs {
			ppout.Reserve0 += pair.Reserve0
			ppout.Reserve1 += pair.Reserve1
			ppout.ReserveUsd += pair.ReserveUsd
			ppout.Token0Price += pair.Token0Price
			ppout.Token1Price += pair.Token1Price
			ppout.TotalSupply += pair.TotalSupply
			// TODO: validate above
			ppout.Token0 = pair.Token0
			ppout.Token1 = pair.Token1
		}
		// Average all these values
		ppout.Reserve0 = ppout.Reserve0 / float64(len(pairs))
		ppout.Reserve1 = ppout.Reserve1 / float64(len(pairs))
		ppout.ReserveUsd = ppout.ReserveUsd / float64(len(pairs))
		ppout.Token0Price = ppout.Token0Price / float64(len(pairs))
		ppout.Token1Price = ppout.Token1Price / float64(len(pairs))
		ppout.TotalSupply = ppout.TotalSupply / float64(len(pairs))
		out = append(out, ppout)
	}
	return out
}

// Parse parses floats from strings
func (ud *UniswapData) Parse() (out []UniswapPairParsed, err error) {
	for _, pair := range ud.Pairs {
		pp := UniswapPairParsed{}
		pp.ID = pair.Id
		pp.Reserve0, err = strconv.ParseFloat(pair.Reserve0, 64)
		if err != nil {
			return nil, err
		}
		pp.Reserve1, err = strconv.ParseFloat(pair.Reserve1, 64)
		if err != nil {
			return nil, err
		}
		pp.ReserveUsd, err = strconv.ParseFloat(pair.ReserveUsd, 64)
		if err != nil {
			return nil, err
		}
		pp.Token0 = pair.Token0
		pp.Token1 = pair.Token1
		pp.Token0Price, err = strconv.ParseFloat(pair.Token0Price, 64)
		if err != nil {
			return nil, err
		}
		pp.Token1Price, err = strconv.ParseFloat(pair.Token1Price, 64)
		if err != nil {
			return nil, err
		}
		pp.TotalSupply, err = strconv.ParseFloat(pair.TotalSupply, 64)
		if err != nil {
			return nil, err
		}
		out = append(out, pp)
	}
	return out, nil
}

// UniswapPairParsed turns the appropriate strings into floats
type UniswapPairParsed struct {
	ID          string
	Reserve0    float64
	Reserve1    float64
	ReserveUsd  float64
	Token0      UniswapToken
	Token1      UniswapToken
	Token0Price float64
	Token1Price float64
	TotalSupply float64
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

// Type implements OracleData
func (ud *UniswapData) Type() string {
	return UniswapDataType
}

// BlocksTillNextPeriod helper
func (vp *VotePeriod) BlocksTillNextPeriod() int64 {
	return vp.VotePeriodEnd - vp.CurrentHeight
}

// // UnpackInterfaces implements UnpackInterfacesMessage.UnpackInterfaces
// func (p Proposal) UnpackInterfaces(unpacker types.AnyUnpacker) error {
// 	var content Content
// 	return unpacker.UnpackAny(p.Content, &content)
// }
