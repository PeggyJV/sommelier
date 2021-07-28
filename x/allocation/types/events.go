package types

// allocation module event types
const (
	EventTypeDelegateAllocations = "delegate_decisions"
	EventTypeAllocationPrecommit = "decision_precommit"
	EventTypeAllocationCommit    = "allocation_commit"
	EventTypeCommitPeriod      = "allocation_commit"

	AttributeKeySigner    = "signer"
	AttributeKeyDelegate  = "delegate"
	AttributeKeyValidator = "validator"

	AttributeKeyCommitPeriodStart = "commit_period_start"
	AttributeKeyCommitPeriodEnd   = "commit_period_end"

	AttributeValueCategory = ModuleName
)
