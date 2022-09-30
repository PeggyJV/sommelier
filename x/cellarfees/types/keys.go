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

	// Store keys
	_ = byte(iota)

	// key for global cellar fee pool state
	CellarFeePoolKey

	// key for storing the reward supply after the latest increase
	LastHighestRewardSupplyKey

	// key for storing the next scheduled auction height
	ScheduledAuctionHeightKey
)

// GetCellarFeePoolKey returns the key prefix
func GetCellarFeePoolKey() []byte {
	return []byte{CellarFeePoolKey}
}

// GetLastHighestRewardSupplyKey returns the key prefix
func GetLastHighestRewardSupplyKey() []byte {
	return []byte{LastHighestRewardSupplyKey}
}

// GetScheduledAuctionHeightKey returns the key prefix
func GetScheduledAuctionHeightKey() []byte {
	return []byte{ScheduledAuctionHeightKey}
}
