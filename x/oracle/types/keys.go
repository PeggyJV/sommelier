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

// Keys for oracle store, with <prefix><key> -> <value>
var (
	// - 0x00<oracle_date_type> -> <OracleData>
	OracleDataKeyPrefix = []byte{0x00} // key for oracle state data

	// - 0x01<val_address> -> <delegate_address>
	FeedDelegateKeyPrefix = []byte{0x01} // key for validator feed delegation

	// - 0x02<val_address> -> <[]hashes>
	OracleDataPrevoteKeyPrefix = []byte{0x02} // key for oracle prevotes

	// - 0x03<val_address> -> <oracle_data_vote>
	OracleDataVoteKeyPrefix = []byte{0x03} // key for oracle votes

	// - 0x04 -> int64(height)
	VotePeriodStartKey = []byte{0x04} // key for vote period height start

	// - 0x05<val_address> -> int64(misses)
	MissCounterKeyPrefix = []byte{0x05} // key for validator miss counters
)

// GetFeedDelegateKey returns the validator for a given delegate key
func GetFeedDelegateKey(del sdk.AccAddress) []byte {
	return append(FeedDelegateKeyPrefix, del.Bytes()...)
}

// GetOracleDataPrevoteKey returns the key for a validators prevote
func GetOracleDataPrevoteKey(val sdk.AccAddress) []byte {
	return append(OracleDataPrevoteKeyPrefix, val.Bytes()...)
}

// GetOracleDataVoteKey returns the key for a validators vote
func GetOracleDataVoteKey(val sdk.AccAddress) []byte {
	return append(OracleDataVoteKeyPrefix, val.Bytes()...)
}

// GetOracleDataKey returns the key for the stored oracle data
func GetOracleDataKey(typ string) []byte {
	return append(OracleDataKeyPrefix, []byte(typ)...)
}

// GetMissCounterKey returns the key for the stored miss counter for a given validator
func GetMissCounterKey(val sdk.AccAddress) []byte {
	return append(MissCounterKeyPrefix, val.Bytes()...)
}
