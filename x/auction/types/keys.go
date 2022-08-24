package types

const (
	// ModuleName is the module name constant used in many places
	ModuleName = "auction"

	// StoreKey is the store key string for auction
	StoreKey = ModuleName

	// RouterKey is the message route for auction
	RouterKey = ModuleName

	// QuerierRoute is the querier route for auction
	QuerierRoute = ModuleName
)

// Keys for auction store, with <prefix><key> -> <value>
const (
	_ = byte(iota)

	// CorkForAddressKeyPrefix - <prefix><val_address><address> -> <cork>
	CurrentAuctionsPrefix 

	// 
	EndedAuctionsPrefix

	// <prefix><denom>
	BidsPrefix

	// <prefix><denom>
	TokenPricesPrefix
)
