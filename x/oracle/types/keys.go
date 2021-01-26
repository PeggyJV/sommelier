package types

import sdk "github.com/cosmos/cosmos-sdk/types"

const (
	// ModuleName is the module name constant used in many places
	ModuleName = "oracle"

	// StoreKey is the store key string for oracle
	StoreKey = ModuleName

	// RouterKey is the message route for oracle
	RouterKey = ModuleName

	// QuerierRoute is the querier route for oracle
	QuerierRoute = ModuleName
)

// Keys for oracle store
// Items are stored with the following key: values
//
// - 0x00<oracle_date_type><>: FeePol
// - 0x01<val_address> -> <delegate_address>
var (
	OracleDataKey              = []byte{0x00} // key for oracle state data
	FeedDelegateKeyPrefix      = []byte{0x01} // key for validator feed delegation
	OracleDataPrevoteKeyPrefix = []byte{0x02} // key for oracle prevotes
	OracleDataVoteKeyPrefix    = []byte{0x02} // key for oracle votes
)

// GetFeedDelegateKey returns the key for a validators feed delegation key
func GetFeedDelegateKey(val sdk.AccAddress) []byte {
	return append(FeedDelegateKeyPrefix, val.Bytes()...)
}

// GetOracleDataPrevoteKey returns the key for a validators prevote
func GetOracleDataPrevoteKey(val sdk.AccAddress) []byte {
	return append(OracleDataPrevoteKeyPrefix, val.Bytes()...)
}

// GetOracleDataVoteKey returns the key for a validators vote
func GetOracleDataVoteKey(val sdk.AccAddress) []byte {
	return append(OracleDataVoteKeyPrefix, val.Bytes()...)
}
