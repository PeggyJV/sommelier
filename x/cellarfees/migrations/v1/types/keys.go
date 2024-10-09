package types

const (
	// ModuleName is the module name constant used in many places
	ModuleName = "cellarfees"

	// StoreKey is the store key string for cellarfees
	StoreKey = ModuleName

	// RouterKey is the message route for cellarfees
	RouterKey = ModuleName

	// QuerierRoute is the querier route for cellarfees
	QuerierRoute = ModuleName
)

const (
	// Store keys
	_ = byte(iota)

	// key for storing the reward supply after the latest increase
	LastRewardSupplyPeakKey

	// key for storing fee accrual counts
	FeeAccrualCountersKey
)

// GetLastRewardSupplyPeakKey returns the key prefix
func GetLastRewardSupplyPeakKey() []byte {
	return []byte{LastRewardSupplyPeakKey}
}

// GetFeeAccrualCountersKey returns the key prefix
func GetFeeAccrualCountersKey() []byte {
	return []byte{FeeAccrualCountersKey}
}
