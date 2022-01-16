package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"
)

const (
	// ModuleName is the module name constant used in many places
	ModuleName = "reinvest"

	// StoreKey is the store key string for oracle
	StoreKey = ModuleName

	// RouterKey is the message route for oracle
	RouterKey = ModuleName

	// QuerierRoute is the querier route for oracle
	QuerierRoute = ModuleName
)

// Keys for reinvest store, with <prefix><key> -> <value>
const (
	_ = byte(iota)

	// ReinvestmentForAddressKeyPrefix - <prefix><val_address><address> -> <reinvestment>
	ReinvestmentForAddressKeyPrefix // key for reinvestments

	// CommitPeriodStartKey - <prefix> -> int64(height)
	CommitPeriodStartKey // key for commit period height start

	// LatestInvalidationNonceKey - <prefix> -> uint64(latestNonce)
	LatestInvalidationNonceKey
)

// GetReinvestmentForValidatorAddressKey returns the key for a validators vote for a given address
func GetReinvestmentForValidatorAddressKey(val sdk.ValAddress, contract common.Address) []byte {
	return append(GetReinvestmentValidatorKeyPrefix(val), contract.Bytes()...)
}

// GetReinvestmentValidatorKeyPrefix returns the key prefix for reinvest commits for a validator
func GetReinvestmentValidatorKeyPrefix(val sdk.ValAddress) []byte {
	return append([]byte{ReinvestmentForAddressKeyPrefix}, val.Bytes()...)
}
