package types

import sdk "github.com/cosmos/cosmos-sdk/types"

const (
	// ModuleName is the module name constant used in many places
	ModuleName = "allocation"

	// StoreKey is the store key string for allocation
	StoreKey = ModuleName

	// RouterKey is the message route for allocation
	RouterKey = ModuleName

	// QuerierRoute is the querier route for allocation
	QuerierRoute = ModuleName
)

// Keys for the allocation store, with <prefix><key> -> <value>
var (
	// - 0x01<val_address> -> <delegate_address>
	AllocationDelegateKeyPrefix = []byte{0x01} // key for validator feed delegation

	// - 0x02<val_address> -> <hash>
	AllocationPrecommitKeyPrefix = []byte{0x02} // key for decision precommits

	// - 0x03<val_address> -> <hash>
	AllocationCommitKeyPrefix = []byte{0x03} // key for decision commits

	// - 0x04 -> int64(height)
	CommitPeriodStartKey = []byte{0x04} // key for vote period height start

	// - 0x05<val_address> -> int64(misses)
	MissCounterKeyPrefix = []byte{0x05} // key for validator miss counters
)

// GetFeedDelegateKey returns the validator for a given delegate key
func GetFeedDelegateKey(del sdk.AccAddress) []byte {
	return append(AllocationDelegateKeyPrefix, del.Bytes()...)
}

// GetAllocationPrecommitKey returns the key for a validators precommit
func GetAllocationPrecommitKey(val sdk.ValAddress) []byte {
	return append(AllocationPrecommitKeyPrefix, val.Bytes()...)
}

// GetAllocationCommitKey returns the key for a validators precommit
func GetAllocationCommitKey(val sdk.ValAddress) []byte {
	return append(AllocationCommitKeyPrefix, val.Bytes()...)
}

// GetMissCounterKey returns the key for the stored miss counter for a given validator
func GetMissCounterKey(val sdk.ValAddress) []byte {
	return append(MissCounterKeyPrefix, val.Bytes()...)
}
