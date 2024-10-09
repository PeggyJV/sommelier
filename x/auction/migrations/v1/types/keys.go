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

	// <prefix>
	LastAuctionIDKey

	// <prefix>
	LastBidIDKey
)

// GetActiveAuctionsPrefix returns the key prefix for active auctions
func GetActiveAuctionsPrefix() []byte {
	return []byte{ActiveAuctionsPrefix}
}

// GetActiveAuctionKey returns the key for an active auction
func GetActiveAuctionKey(id uint32) []byte {
	b := make([]byte, 4)
	binary.BigEndian.PutUint32(b, id)
	return append(GetActiveAuctionsPrefix(), b...)
}

// GetEndedAuctionsPrefix returns the key prefix for ended auctions
func GetEndedAuctionsPrefix() []byte {
	return []byte{EndedAuctionsPrefix}
}

// GetEndedAuctionsKey returns the key for an ended auction
func GetEndedAuctionKey(id uint32) []byte {
	b := make([]byte, 4)
	binary.BigEndian.PutUint32(b, id)
	return append(GetEndedAuctionsPrefix(), b...)
}

// GetBidsByAuctionPrefix returns the key prefix for bids
func GetBidsByAuctionPrefix() []byte {
	return []byte{BidsByAuctionPrefix}
}

// GetBidsByAuctionIDPrefix returns the bids for an auction id
func GetBidsByAuctionIDPrefix(auctionID uint32) []byte {
	b := make([]byte, 4)
	binary.BigEndian.PutUint32(b, auctionID)
	return append(GetBidsByAuctionPrefix(), b...)
}

// GetBidKey returns the bid for an auction and bid id
func GetBidKey(auctionID uint32, bidID uint64) []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, bidID)
	return append(GetBidsByAuctionIDPrefix(auctionID), b...)
}

// GetTokenPricesPrefix returns the key prefix for token prices
func GetTokenPricesPrefix() []byte {
	return []byte{TokenPricesPrefix}
}

// GetTokenPriceKey returns the key for a token price
func GetTokenPriceKey(denom string) []byte {
	return append(GetTokenPricesPrefix(), []byte(denom)...)
}

// GetLastAuctionIDKey returns the key prefix for the last stored auction id
func GetLastAuctionIDKey() []byte {
	return []byte{LastAuctionIDKey}
}

// GetLastBidIDKey returns the key prefix for the last stored bid id
func GetLastBidIDKey() []byte {
	return []byte{LastBidIDKey}
}
