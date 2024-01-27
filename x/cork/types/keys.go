package types

import (
	"bytes"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"
)

const (
	// ModuleName is the module name constant used in many places
	ModuleName = "cork"

	// StoreKey is the store key string for oracle
	StoreKey = ModuleName

	// RouterKey is the message route for oracle
	RouterKey = ModuleName

	// QuerierRoute is the querier route for oracle
	QuerierRoute = ModuleName
)

// Keys for cork store, with <prefix><key> -> <value>
const (
	_ = byte(iota)

	// CorkForAddressKeyPrefix - <prefix><val_address><address> -> <cork>
	CorkForAddressKeyPrefix // key for corks -- removed as of cork v2, left to preserve ID values

	// CommitPeriodStartKey - <prefix> -> int64(height)
	CommitPeriodStartKey // key for commit period height start -- removed as of cork v2, left to preserve ID values

	// LatestInvalidationNonceKey - <prefix> -> uint64(latestNonce)
	LatestInvalidationNonceKey

	// CellarIDsKey - <prefix> -> []string
	CellarIDsKey

	// ScheduledCorkKeyPrefix - <prefix><block_height><val_address><address> -> <cork>
	ScheduledCorkKeyPrefix

	// LatestCorkIDKey - <key> -> uint64(latestCorkID)
	LatestCorkIDKey

	// CorkResultPrefix - <prefix><id> -> CorkResult
	CorkResultPrefix

	// ValidatorCorkCountKey - <prefix><val_address> -> uint64(count)
	ValidatorCorkCountKey
)

func MakeCellarIDsKey() []byte {
	return []byte{CellarIDsKey}
}

func GetScheduledCorkKeyPrefix() []byte {
	return []byte{ScheduledCorkKeyPrefix}
}

func GetScheduledCorkKeyByBlockHeightPrefix(blockHeight uint64) []byte {
	return append(GetScheduledCorkKeyPrefix(), sdk.Uint64ToBigEndian(blockHeight)...)
}

func GetScheduledCorkKey(blockHeight uint64, id []byte, val sdk.ValAddress, contract common.Address) []byte {
	blockHeightBytes := sdk.Uint64ToBigEndian(blockHeight)
	return bytes.Join([][]byte{GetScheduledCorkKeyPrefix(), blockHeightBytes, id, val.Bytes(), contract.Bytes()}, []byte{})
}

func GetCorkResultPrefix() []byte {
	return []byte{CorkResultPrefix}
}

func GetCorkResultKey(id []byte) []byte {
	return append(GetCorkResultPrefix(), id...)
}

func GetValidatorCorkCountKey(val sdk.ValAddress) []byte {
	return append([]byte{ValidatorCorkCountKey}, val.Bytes()...)
}
