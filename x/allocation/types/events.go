package types

// allocation module event types
const (
	EventTypeDelegateAllocations = "delegate_allocations"
	EventTypeAllocationPrecommit = "allocation_precommit"
	EventTypeAllocationCommit    = "allocation_commit"
	EventTypeCommitPeriod        = "commit_period"

	AttributeKeySigner            = "signer"
	AttributeKeyDeleagate         = "delegate"
	AttributeKeyValidator         = "validator"
	AttributeKeyPrevoteHash       = "hash"
	AttributeKeyCellar            = "cellar"
	AttributePoolAllocations      = "pool_allocations"
	AttributeKeyCommitPeriodStart = "commit_period_start"
	AttributeKeyCommitPeriodEnd   = "commit_period_end"

	AttributeValueCategory = ModuleName
)
