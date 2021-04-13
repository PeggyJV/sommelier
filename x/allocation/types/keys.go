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
	AllocationDecisionPrecommitKeyPrefix = []byte{0x02} // key for decision precommits

	// - 0x03<val_address> -> <hash>
	AllocationDecisionCommitKeyPrefix = []byte{0x03} // key for decision precommits
)

// GetFeedDelegateKey returns the validator for a given delegate key
func GetFeedDelegateKey(del sdk.AccAddress) []byte {
	return append(AllocationDelegateKeyPrefix, del.Bytes()...)
}

// GetDecisionPrecommitKey returns the key for a validators precommit
func GetDecisionPrecommitKey(val sdk.ValAddress) []byte {
	return append(AllocationDecisionPrecommitKeyPrefix, val.Bytes()...)
}


// GetDecisionCommitKey returns the key for a validators precommit
func GetDecisionCommitKey(val sdk.ValAddress) []byte {
	return append(AllocationDecisionCommitKeyPrefix, val.Bytes()...)
}