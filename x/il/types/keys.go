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
	// Stoploss positions key prefix
	StoplossKeyPrefix = []byte{0x01}
	// Invalidation ID prefix for outgoing logic
	InvalidationIDPrefix = []byte{0x02}
	// Submitted positions queue prefix for
	SubmittedPositionsQueuePrefix = []byte{0x03}
)

// StoplossKey defines the full unprefixed store key for Stoploss
func StoplossKey(address sdk.AccAddress, uniswapPair string) []byte {
	return append(address.Bytes(), []byte(uniswapPair)...)
}

// SubmittedPositionKey defines the full unprefixed store key for submitted positions to the bridge
func SubmittedPositionKey(timeoutHeight uint64, address sdk.AccAddress) []byte {
	return append(sdk.Uint64ToBigEndian(timeoutHeight), address.Bytes()...)
}

// LPAddressFromStoplossKey
func LPAddressFromStoplossKey(key []byte) sdk.AccAddress {
	if len(key[1:]) < sdk.AddrLen {
		return nil
	}

	return sdk.AccAddress(key[1 : 1+sdk.AddrLen])
}

// SplitSubmittedStoplossKey
func SplitSubmittedStoplossKey(key []byte) (uint64, sdk.AccAddress) {
	if len(key[1:]) < 8+sdk.AddrLen {
		return 0, nil
	}

	return sdk.BigEndianToUint64(key[1:9]), sdk.AccAddress(key[9 : 9+sdk.AddrLen])
}
