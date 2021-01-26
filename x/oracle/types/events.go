package types

// distribution module event types
const (
	EventTypeDelegateFeed      = "delegate_feed"
	EventTypeOracleDataPrevote = "oracle_data_prevote"
	EventTypeOracleDataVote    = "oracle_data_vote"

	AttributeKeySigner    = "signer"
	AttributeKeyDeleagte  = "delegate"
	AttributeKeyValidator = "validator"
	AttributeKeyHashes    = "hashes"

	AttributeValueCategory = ModuleName
)
