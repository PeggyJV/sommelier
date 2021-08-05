package types

import (
	"bytes"

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

// Keys for allocation store, with <prefix><key> -> <value>
const (
	_ = byte(iota)
	// CellarKeyPrefix - <prefix><cellar_id> -> Cellar
	CellarKeyPrefix

	// AllocationDelegateKeyPrefix - <prefix><val_address> -> <delegate_address>
	AllocationDelegateKeyPrefix // key for validator allocation delegation

	// AllocationPrecommitKeyPrefix - <prefix><val_address><cel_address> -> <hash>
	AllocationPrecommitKeyPrefix // key for allocation precommits

	// AllocationCommitForCellarKeyPrefix - <prefix><val_address><cel_address> -> <allocation_commit>
	AllocationCommitForCellarKeyPrefix // key for allocation commits

	// CommitPeriodStartKey - <prefix> -> int64(height)
	CommitPeriodStartKey // key for commit period height start

	// LatestInvalidationNonceKey - <prefix> -> uint64(latestNonce)
	LatestInvalidationNonceKey

	// CellarUpdateKey - <prefix><invalidationNonce> -> Cellar
	CellarUpdateKey
)

// GetAllocationDelegateKey returns the validator for a given delegate key
func GetAllocationDelegateKey(del sdk.AccAddress) []byte {
	return append([]byte{AllocationDelegateKeyPrefix}, del.Bytes()...)
}

// GetAllocationPrecommitKey returns the key for a validators prevote for a cellar
func GetAllocationPrecommitKey(val sdk.ValAddress, cel common.Address) []byte {
	return bytes.Join([][]byte{{AllocationPrecommitKeyPrefix}, val.Bytes(), cel.Bytes()}, []byte{})
}

// GetAllocationCommitForCellarKey returns the key for a validators vote for a given cellar
func GetAllocationCommitForCellarKey(val sdk.ValAddress, cel common.Address) []byte {
	return append(GetAllocationCommitKeyPrefix(val), cel.Bytes()...)
}

// GetAllocationCommitKeyPrefix returns the key prefix for allocation commits for a validator
func GetAllocationCommitKeyPrefix(val sdk.ValAddress) []byte {
	return append([]byte{AllocationCommitForCellarKeyPrefix}, val.Bytes()...)
}

// GetCellarUpdateKey - the pending cellar update by invalidation nonce
func GetCellarUpdateKey(invalidationNonce uint64) []byte {
	return append([]byte{CellarKeyPrefix}, sdk.Uint64ToBigEndian(invalidationNonce)...)
}

// GetCellarKey - the cellar by id
func GetCellarKey(address common.Address) []byte {
	return append([]byte{CellarKeyPrefix}, address.Bytes()...)
}