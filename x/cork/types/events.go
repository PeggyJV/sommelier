package types

// cork module event types
const (
	EventTypeCork                  = "cork"
	EventTypeCommitPeriod          = "commit_period"
	EventTypeSubmittedContractCall = "submitted_contract_call"

	AttributeKeySigner            = "signer"
	AttributeKeyValidator         = "validator"
	AttributeKeyPrevoteHash       = "hash"
	AttributeKeyCork              = "cork"
	AttributeKeyCommitPeriodStart = "commit_period_start"
	AttributeKeyCommitPeriodEnd   = "commit_period_end"
	AttributeKeyBlockHeight       = "block_height"
	AttributeKeyCorkId            = "cork_id"

	AttributeValueCategory = ModuleName
)
