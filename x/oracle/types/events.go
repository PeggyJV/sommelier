package types

// distribution module event types
const (
	EventTypeDelegateFeed      = "delegate_feed"
	EventTypeOracleDataPrevote = "oracle_data_prevote"
	EventTypeOracleDataVote    = "oracle_data_vote"
	EventTypeVotePeriod        = "vote_period"

	AttributeKeySigner          = "signer"
	AttributeKeyDeleagte        = "delegate"
	AttributeKeyValidator       = "validator"
	AttributeKeyHashes          = "hashes"
	AttributeKeyVotePeriodStart = "vote_period_start"
	AttributeKeyVotePeriodEnd   = "vote_period_end"

	AttributeValueCategory = ModuleName
)
