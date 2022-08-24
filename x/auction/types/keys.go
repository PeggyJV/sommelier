package types

import (
	"encoding/binary"
)

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

	// <prefix><auction_id> 
	ActiveAuctionsPrefix 

	// <prefix><auction_id>
	EndedAuctionsPrefix

	// <prefix><auction_id><bid_id>
	BidsByAuctionPrefix

	// <prefix><denom>
	TokenPricesPrefix
)

// GetActiveAuctionsPrefix returns the key prefix for active auctions
func GetActiveAuctionsPrefix() []byte {
	return []byte{ActiveAuctionsPrefix}
}

// GetActiveAuctionKey returns the key for an active auction
func GetActiveAuctionKey(id uint32) []byte {
	b := make([]byte, 4)
	binary.BigEndian.PutUint32(b, id)
	return append([]byte{ActiveAuctionsPrefix}, b...)
}

// GetEndedAuctionsPrefix returns the key prefix for ended auctions
func GetEndedAuctionsPrefix() []byte {
	return []byte{EndedAuctionsPrefix}
}

// GetEndedAuctionsKey returns the key for an ended auction
func GetEndedAuctionKey(id uint32) []byte {
	b := make([]byte, 4)
	binary.BigEndian.PutUint32(b, id)
	return append([]byte{EndedAuctionsPrefix}, b...)
}

// GetBidsByAuctionPrefix returns the key prefix for bids
func GetBidsByAuctionPrefix() []byte {
	return []byte{BidsByAuctionPrefix}
}

