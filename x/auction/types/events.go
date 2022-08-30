package types

// Auction module event types
const (
	EventTypeBid = "bid"

	AttributeKeyBidId           = "bid_id"
	AttributeKeyAuctionId       = "auction_id"
	AttributeKeyBidder          = "max_bid"
	AttributeKeyMinimumAmount   = "minimum_amount"
	AttributeKeySigner          = "signer"
	AttributeKeyFulfilledPrice  = "fulfilled_price_of_sale_token_in_usomm"
	AttributeKeyFulfilledAmount = "total_sale_token_fulfilled_amount"

	EventTypeNewAuction = "new_auction"

	AttributeKeyStartBlock            = "start_block"
	AttributeKeyInitialDecreaseRate   = "initial_decrease_rate"
	AttributeKeyBlockDecreaseInterval = "block_decrease_interval"
	AttributeKeyStartingAmount        = "starting_amount_of_sale_token"
	AttributeKeyStartingDenom         = "starting_denom"
	AttributeKeyStartingUsommPrice    = "starting_price_in_usomm"

	EventTypeAuctionFinished        = "auction_finished"

	AttributeKeyEndBlock            = "end_block"
	AttributeKeyCurrentDecreaseRate = "current_decrease_rate"
	AttributeKeyAmountRemaining     = "amount_remaining"

	EventTypeAuctionUpdated        = "auction_updated"

	AttributeValueCategory = ModuleName
)
