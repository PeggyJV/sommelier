package types

// auction module event types
const (
	// TODO: fill in remaining event types and their attributes, examples below
	
	EventTypeBid                 			  = "bid"

	AttributeKeyBidId            			  = "bid_id"
	AttributeKeyAuctionId        			  = "auction_id"
	AttributeKeyBidder			 			  = "max_bid"
	AttributeKeyMinimumAmount   			  = "minimum_amount"
	AttributeKeySigner            			  = "signer"
	AttributeKeyFulfilledPrice    			  = "fulfilled_price"
	AttributeKeyFulfilledAmount  			  = "fulfilled_amount"

	EventTypeAuction              			  = "auction"
	AttributeKeyStartBlock       			  = "start_block"
	AttributeKeyInitialDecreaseRate           = "initial_decrease_rate"
	AttributeKeyBlockDecreaseInterval         = "block_decrease_interval"
	AttributeKeyStartingAmount   			  = "starting_amount"
	AttributeKeyStartingDenom   			  = "starting_denom"
	AttributeKeyStartingUsommPrice     		  = "starting_price_in_usomm"

	AttributeValueCategory = ModuleName
)
