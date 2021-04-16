package types

import (
	"crypto/sha256"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	// ModuleName is the module name constant used in many places
	ModuleName = "allocation"

	// StoreKey is the store key string for oracle
	StoreKey = ModuleName

	// RouterKey is the message route for oracle
	RouterKey = ModuleName

	// QuerierRoute is the querier route for oracle
	QuerierRoute = ModuleName
)

// Keys for oracle store, with <prefix><key> -> <value>
var (
	// - 0x01<oracle_data_type_hash><oracle_data_id> -> <OracleData>
	OracleDataKeyPrefix = []byte{0x01} // key for oracle state data

	// - 0x02<oracle_data_id> -> <oracle_data_type>
	OracleDataTypeKeyPrefix = []byte{0x02} //

	// - 0x03<oracle_data_id> -> uint64
	OracleDataHeightKeyPrefix = []byte{0x03}

	// - 0x04<val_address> -> <delegate_address>
	AllocationDelegateKeyPrefix = []byte{0x04} // key for validator allocation delegation

	// - 0x05<val_address> -> <[]hashes>
	AllocationPrecommitKeyPrefix = []byte{0x05} // key for allocation precommits

	// - 0x06<val_address> -> <allocation_commit>
	AllocationCommitKeyPrefix = []byte{0x06} // key for allocation commits

	// - 0x07 -> int64(height)
	CommitPeriodStartKey = []byte{0x07} // key for commit period height start

	// - 0x08<val_address> -> int64(misses)
	MissCounterKeyPrefix = []byte{0x08} // key for validator miss counters

	// - 0x09<oracle_data_type_hash><oracle_data_id> -> <OracleData>
	AggregatedOracleDataKeyPrefix = []byte{0x09} // key for oracle state data
)

// GetAllocationDelegateKey returns the validator for a given delegate key
func GetAllocationDelegateKey(del sdk.AccAddress) []byte {
	return append(AllocationDelegateKeyPrefix, del.Bytes()...)
}

// GetAllocationPrecommitKey returns the key for a validators prevote
func GetAllocationPrecommitKey(val sdk.ValAddress) []byte {
	return append(AllocationPrecommitKeyPrefix, val.Bytes()...)
}

// GetAllocationCommitKey returns the key for a validators vote
func GetAllocationCommitKey(val sdk.ValAddress) []byte {
	return append(AllocationCommitKeyPrefix, val.Bytes()...)
}

func GetAggregatedOracleDataKey(height uint64, dataType, id string) []byte {
	dataTypeHash := sha256.Sum256([]byte(dataType))
	key := append(AggregatedOracleDataKeyPrefix, sdk.Uint64ToBigEndian(height)...)
	key = append(key, dataTypeHash[:]...)
	key = append(key, []byte(id)...)
	return key
}

// GetOracleDataKey returns the key for the stored oracle data
func GetOracleDataKey(dataType, id string) []byte {
	dataTypeHash := sha256.Sum256([]byte(dataType))
	key := append(OracleDataKeyPrefix, dataTypeHash[:]...)
	return append(key, []byte(id)...)
}

// GetOracleDataTypeKey returns the key for the stored oracle data type
func GetOracleDataTypeKey(id string) []byte {
	return append(OracleDataTypeKeyPrefix, []byte(id)...)
}

// GetOracleDataHeightKey returns the key for the latest oracle data height
func GetOracleDataHeightKey(id string) []byte {
	return append(OracleDataHeightKeyPrefix, []byte(id)...)
}

// GetMissCounterKey returns the key for the stored miss counter for a given validator
func GetMissCounterKey(val sdk.ValAddress) []byte {
	return append(MissCounterKeyPrefix, val.Bytes()...)
}
