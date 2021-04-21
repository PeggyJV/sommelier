package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"
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

	// - 0x02<val_address><cel_address> -> <[]tick_weight>
	AllocationTickWeightKeyPrefix = []byte{0x02} //

	// - 0x03<oracle_data_id> -> uint64
	OracleDataHeightKeyPrefix = []byte{0x03}

	// - 0x04<val_address> -> <delegate_address>
	AllocationDelegateKeyPrefix = []byte{0x04} // key for validator allocation delegation

	// - 0x05<val_address><cel_address> -> <hash>
	AllocationPrecommitKeyPrefix = []byte{0x05} // key for allocation precommits

	// - 0x06<val_address><cel_address> -> <allocation_commit>
	AllocationCommitForCellarKeyPrefix = []byte{0x06} // key for allocation commits

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

// GetAllocationPrecommitKey returns the key for a validators prevote for a cellar
func GetAllocationPrecommitKey(val sdk.ValAddress, cel common.Address) []byte {
	key := append(AllocationPrecommitKeyPrefix, val.Bytes()...)
	return append(key, cel.Bytes()...)
}

// GetAllocationCommitForCellarKey returns the key for a validators vote for a given cellar
func GetAllocationCommitForCellarKey(val sdk.ValAddress, cel common.Address) []byte {
	key := GetAllocationCommitKeyPrefix(val)
	return append(key, cel.Bytes()...)
}

// GetAllocationCommitKeyPrefix returns the key prefix for allocation commits for a validator
func GetAllocationCommitKeyPrefix(val sdk.ValAddress) []byte {
	return append(AllocationCommitForCellarKeyPrefix, val.Bytes()...)
}

// GetAllocationTickWeightKey returns the key for tick_weights for a given cellar
func GetAllocationTickWeightKey(val sdk.ValAddress, cel common.Address) []byte {
	key := append(AllocationTickWeightKeyPrefix, val.Bytes()...)
	return append(key, cel.Bytes()...)
}

// GetMissCounterKey returns the key for the stored miss counter for a given validator
func GetMissCounterKey(val sdk.ValAddress) []byte {
	return append(MissCounterKeyPrefix, val.Bytes()...)
}
