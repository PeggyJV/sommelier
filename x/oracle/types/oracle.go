package types

import (
	"crypto/sha256"
	"encoding/json"
	fmt "fmt"
	"strings"

	proto "github.com/gogo/protobuf/proto"

	tmbytes "github.com/tendermint/tendermint/libs/bytes"

	sdk "github.com/cosmos/cosmos-sdk/types"

	peggytypes "github.com/althea-net/peggy/module/x/peggy/types"
)

var (
	_ OracleData     = &UniswapPair{}
	_ json.Marshaler = &UniswapPair{}
)

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
	// TODO: Add a parsing function to this and figure out what the signature needs to be
	// Parse() error
}

// DataHash returns the hash for a precommit given the proper args
func DataHash(salt, jsonData string, signer sdk.ValAddress) []byte {
	h := sha256.New()
	h.Write([]byte(fmt.Sprintf("%s:%s:%s", salt, jsonData, signer.String())))
	return h.Sum(nil)
}

type uniswapPairPretty struct {
	ID          string             `json:"id,omitempty"`
	Reserve0    string             `json:"reserve0" yaml:"reserve0"`
	Reserve1    string             `json:"reserve1" yaml:"reserve1"`
	ReserveUSD  string             `json:"reserveUSD" yaml:"reserveUSD"`
	Token0      uniswapTokenPretty `json:"token0"`
	Token1      uniswapTokenPretty `json:"token1"`
	Token0Price string             `json:"token0Price" yaml:"token0Price"`
	Token1Price string             `json:"token1Price" yaml:"token1Price"`
	TotalSupply string             `json:"totalSupply" yaml:"totalSupply"`
}

type uniswapTokenPretty struct {
	// token address
	ID string `json:"id,omitempty"`
	// number of decimal positions of the pair token
	Decimals string `json:"decimals"`
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
			ppout.Reserve0 = ppout.Reserve0.Add(pair.Reserve0)
			ppout.Reserve1 = ppout.Reserve1.Add(pair.Reserve1)
			ppout.ReserveUsd = ppout.ReserveUsd.Add(pair.ReserveUsd)
			ppout.Token0Price = ppout.Token0Price.Add(pair.Token0Price)
			ppout.Token1Price = ppout.Token1Price.Add(pair.Token1Price)
			ppout.TotalSupply = ppout.TotalSupply.Add(pair.TotalSupply)
			// TODO: validate above
			ppout.Token0 = pair.Token0
			ppout.Token1 = pair.Token1
		}
		// Average all these values
		ppout.Reserve0 = ppout.Reserve0.Quo(sdk.NewDec(int64(len(pairs))))
		ppout.Reserve1 = ppout.Reserve1.Quo(sdk.NewDec(int64(len(pairs))))
		ppout.ReserveUsd = ppout.ReserveUsd.Quo(sdk.NewDec(int64(len(pairs))))
		ppout.Token0Price = ppout.Token0Price.Quo(sdk.NewDec(int64(len(pairs))))
		ppout.Token1Price = ppout.Token1Price.Quo(sdk.NewDec(int64(len(pairs))))
		ppout.TotalSupply = ppout.TotalSupply.Quo(sdk.NewDec(int64(len(pairs))))
		out = append(out, ppout)
	}
	return out
}

// Parse parses floats from strings
func (ud *UniswapData) Parse() (out []UniswapPairParsed, err error) {
	for _, pair := range ud.Pairs {
		pp := UniswapPairParsed{}
		pp.ID = pair.Id
		pp.Reserve0, err = sdk.NewDecFromStr(normalizeDec(pair.Reserve0))
		if err != nil {
			return nil, err
		}
		pp.Reserve1, err = sdk.NewDecFromStr(normalizeDec(pair.Reserve1))
		if err != nil {
			return nil, err
		}
		pp.ReserveUsd, err = sdk.NewDecFromStr(normalizeDec(pair.ReserveUsd))
		if err != nil {
			return nil, err
		}
		pp.Token0 = pair.Token0
		pp.Token1 = pair.Token1
		pp.Token0Price, err = sdk.NewDecFromStr(normalizeDec(pair.Token0Price))
		if err != nil {
			return nil, err
		}
		pp.Token1Price, err = sdk.NewDecFromStr(normalizeDec(pair.Token1Price))
		if err != nil {
			return nil, err
		}
		pp.TotalSupply, err = sdk.NewDecFromStr(normalizeDec(pair.TotalSupply))
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
	Reserve0    sdk.Dec
	Reserve1    sdk.Dec
	ReserveUsd  sdk.Dec
	Token0      UniswapToken
	Token1      UniswapToken
	Token0Price sdk.Dec
	Token1Price sdk.Dec
	TotalSupply sdk.Dec
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

	if up.ID != aggregatedPair.ID {
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
	if up.Reserve1.Sub(aggregatedPair.Reserve1).Abs().GT(target.Mul(aggregatedPair.Reserve1)) {
		return false
	}

	// |reserveUsd - reserveUsd (agg)| / (reserveUsd (agg)) ≤ target
	if up.ReserveUSD.Sub(aggregatedPair.ReserveUSD).Abs().GT(target.Mul(aggregatedPair.ReserveUSD)) {
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

// MarshalJSON is a custom JSON marshaler that parses the uniswap pair into the format that corresponds
// to the Graph query.
func (up UniswapPair) MarshalJSON() ([]byte, error) {
	upp := uniswapPairPretty{
		ID:         up.ID,
		Reserve0:   up.Reserve0.String(),
		Reserve1:   up.Reserve1.String(),
		ReserveUSD: up.ReserveUSD.String(),
		Token0: uniswapTokenPretty{
			ID:       up.Token0.ID,
			Decimals: strconv.FormatUint(up.Token0.Decimals, 10),
		},
		Token1: uniswapTokenPretty{
			ID:       up.Token1.ID,
			Decimals: strconv.FormatUint(up.Token1.Decimals, 10),
		},
		Token0Price: up.Token0Price.String(),
		Token1Price: up.Token1Price.String(),
		TotalSupply: up.TotalSupply.String(),
	}

	return json.Marshal(upp)
}

// UnmarshalJSON is a custom JSON marshaler that chops the decimals to the
// max precision allowed by the SDK (18).
func (up *UniswapPair) UnmarshalJSON(bz []byte) error {
	var upp uniswapPairPretty

	err := json.Unmarshal(bz, &upp)
	if err != nil {
		return err
	}

	up.ID = upp.ID

	token0dec, err := strconv.ParseUint(upp.Token0.Decimals, 10, 64)
	if err != nil {
		return err
	}

	token1dec, err := strconv.ParseUint(upp.Token1.Decimals, 10, 64)
	if err != nil {
		return err
	}

	up.Token0 = UniswapToken{
		ID:       upp.Token0.ID,
		Decimals: token0dec,
	}
	up.Token1 = UniswapToken{
		ID:       upp.Token1.ID,
		Decimals: token1dec,
	}

	up.Reserve0, err = TruncateDec(upp.Reserve0)
	if err != nil {
		return fmt.Errorf("reserve 0 (%s), pair (%s): %w", upp.Reserve0, upp.ID, err)
	}

	up.Reserve1, err = TruncateDec(upp.Reserve1)
	if err != nil {
		return fmt.Errorf("reserve 1 (%s), pair (%s): %w", upp.Reserve1, upp.ID, err)
	}

	up.ReserveUSD, err = TruncateDec(upp.ReserveUSD)
	if err != nil {
		return fmt.Errorf("reserve USD (%s), pair (%s): %w", upp.ReserveUSD, upp.ID, err)
	}

	up.Token0Price, err = TruncateDec(upp.Token0Price)
	if err != nil {
		return fmt.Errorf("token 0 price (%s), pair (%s): %w", upp.Token0Price, upp.ID, err)
	}

	up.Token1Price, err = TruncateDec(upp.Token1Price)
	if err != nil {
		return fmt.Errorf("token 1 price (%s), pair (%s): %w", upp.Token1Price, upp.ID, err)
	}

	up.TotalSupply, err = TruncateDec(upp.TotalSupply)
	if err != nil {
		return fmt.Errorf("total supply (%s), pair (%s): %w", upp.TotalSupply, upp.ID, err)
	}

	return nil
}

// Validate performs a basic validation of the uniswap token fields.
func (ut UniswapToken) Validate() error {
	if err := peggytypes.ValidateEthAddress(ut.ID); err != nil {
		return fmt.Errorf("invalid token address %s: %w", ut.ID, err)
	}

	// TODO: figure out how to handle higher precision tokens on ETH
	// if ut.Decimals > sdk.Precision {
	// 	return fmt.Errorf("decimal places (%d) exceeds the maximum supported (%d)", ut.Decimals, sdk.Precision)
	// }

func normalizeDec(str string) string {
	spl := strings.Split(str, ".")
	if len(spl) == 1 {
		return str
	}
	// if there are more than 1 period, then just return 0
	if len(spl) > 2 {
		return "0.0"
	}
	return fmt.Sprintf("%s.%.18s", spl[0], spl[1])
}
