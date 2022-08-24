package types

// auction module event types
const (
	// TODO: fill in remaining event types and their attributes, examples below
	
	EventTypeBid                  = "bid"

	AttributeKeyBidId             = "bid_id"
	AttributeKeyAuctionId         = "auction_id"
	AttributeKeyBidder			  = "max_bid"
	AttributeKeyMinimumAmount     = "minimum_amount"
	AttributeKeySigner            = "signer"
	AttributeKeyFulfilledPrice    = "fulfilled_price"
	AttributeKeyFulfilledAmount   = "fulfilled_amount"

	AttributeValueCategory = ModuleName
)
