package types

// allocation module event types
const (
	EventTypeDelegateDecisions = "delegate_decisions"
	EventTypeDecisionPrecommit = "decision_precommit"
	EventTypeDecisionCommit    = "allocation_commit"
	EventTypeCommitPeriod      = "allocation_commit"

	AttributeKeySigner    = "signer"
	AttributeKeyDelegate  = "delegate"
	AttributeKeyValidator = "validator"

	AttributeKeyCommitPeriodStart = "commit_period_start"
	AttributeKeyCommitPeriodEnd   = "commit_period_end"

	AttributeValueCategory = ModuleName
)
