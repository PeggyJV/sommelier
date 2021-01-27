package types

import sdk "github.com/cosmos/cosmos-sdk/types"

const (
	// ModuleName is the name of the impermanent loss module
	ModuleName = "il"

	// StoreKey is the string store representation
	StoreKey = ModuleName

	// RouterKey is the msg router key for the oracle module
	RouterKey = ModuleName

	// QuerierRoute is the query router key for the oracle module
	QuerierRoute = ModuleName
)

// Keys for oracle store
// Items are stored with the following key: values
//
// - 0x01<address_Bytes><uniswap_pair_id>: Stoploss
var (
	// Keys for store prefixes
	StoplossKeyPrefix = []byte{0x01} // prefix for each key to a stoploss
)

// StoplossKey defines the full unprefixed store key for Stoploss
func StoplossKey(address sdk.AccAddress, uniswapPair string) []byte {
	return append(address.Bytes(), []byte(uniswapPair)...)
}

// LPAddressFromStoplossKey d
func LPAddressFromStoplossKey(key []byte) sdk.AccAddress {
	if len(key[1:]) < sdk.AddrLen {
		return nil
	}

	return sdk.AccAddress(key[1 : 1+sdk.AddrLen])
}
