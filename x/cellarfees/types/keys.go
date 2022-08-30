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

var (
	// key for global cellar fee pool state
	CellarFeePoolKey = []byte{0x00}
	// key for storing the reward supply after the latest increase
	LastHighestRewardSupply = []byte{0x01}
	// key for storing the next scheduled auction height
	ScheduledAuctionHeight = []byte{0x02}
)
