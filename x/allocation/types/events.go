package types

// allocation module event types
const (
	EventTypeDelegateDecisions = "delegate_decisions"
	EventTypeDecisionPrecommit = "decision_precommit"
	EventTypeDecisionCommit = "decision_commit"

	AttributeKeySigner = "signer"
	AttributeKeyDelegate = "delegate"
	AttributeKeyValidator = "validator"

	AttributeValueCategory = ModuleName
)