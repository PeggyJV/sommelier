package types

// cork module event types
const (
	EventTypeCork                  = "cork"
	EventTypeCommitPeriod          = "commit_period"
	EventTypeSubmittedContractCall = "submitted_contract_call"
	// EventTypeCorksDroppedInSafeMode reports approved scheduled corks that were
	// tallied and garbage-collected but NOT submitted, because x/poa was in
	// authority-empty safe mode at their target height. They must be
	// re-scheduled after recovery.
	EventTypeCorksDroppedInSafeMode = "corks_dropped_safe_mode"

	// EventTypeCorkDroppedUnmanagedCellar reports a due cork discarded because
	// its target cellar was removed from the allowlist after scheduling.
	EventTypeCorkDroppedUnmanagedCellar = "cork_dropped_unmanaged_cellar"

	AttributeKeySigner            = "signer"
	AttributeKeyValidator         = "validator"
	AttributeKeyPrevoteHash       = "hash"
	AttributeKeyCork              = "cork"
	AttributeKeyCommitPeriodStart = "commit_period_start"
	AttributeKeyCommitPeriodEnd   = "commit_period_end"
	AttributeKeyBlockHeight       = "block_height"
	AttributeKeyCorkID            = "cork_id"
	AttributeKeyDroppedCount      = "dropped_count"

	AttributeValueCategory = ModuleName
)
