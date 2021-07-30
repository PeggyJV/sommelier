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
	// PoolAllocationKeyPrefix - <prefix><val_address><cel_address> -> <[]pool_allocation>
	PoolAllocationKeyPrefix = []byte{0x01} //

	// AllocationDelegateKeyPrefix - <prefix><val_address> -> <delegate_address>
	AllocationDelegateKeyPrefix = []byte{0x02} // key for validator allocation delegation

	// AllocationPrecommitKeyPrefix - <prefix><val_address><cel_address> -> <hash>
	AllocationPrecommitKeyPrefix = []byte{0x03} // key for allocation precommits

	// AllocationCommitForCellarKeyPrefix - <prefix><val_address><cel_address> -> <allocation_commit>
	AllocationCommitForCellarKeyPrefix = []byte{0x04} // key for allocation commits

	// CommitPeriodStartKey - <prefix> -> int64(height)
	CommitPeriodStartKey = []byte{0x05} // key for commit period height start

	// MissCounterKeyPrefix - <prefix><val_address> -> int64(misses)
	MissCounterKeyPrefix = []byte{0x06} // key for validator miss counters
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

// GetPoolAllocationKey returns the key for pool allocations for a given cellar
func GetPoolAllocationKey(val sdk.ValAddress, cel common.Address) []byte {
	key := append(PoolAllocationKeyPrefix, val.Bytes()...)
	return append(key, cel.Bytes()...)
}

// GetMissCounterKey returns the key for the stored miss counter for a given validator
func GetMissCounterKey(val sdk.ValAddress) []byte {
	return append(MissCounterKeyPrefix, val.Bytes()...)
}
