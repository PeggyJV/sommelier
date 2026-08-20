package types

// axelarcork module event types
const (
	EventTypeAxelarCorkApproved    = "axelar_cork_approved"
	EventTypeAxelarCorkRelayCalled = "axelar_cork_relay_called"
	// EventTypeAxelarCorksDroppedInSafeMode reports approved corks that were
	// tallied and garbage-collected but NOT marked relayable, because x/poa was
	// in authority-empty safe mode at their target height. They must be
	// re-scheduled after recovery.
	EventTypeAxelarCorksDroppedInSafeMode = "axelar_corks_dropped_safe_mode"

	AttributeKeyCork         = "cork"
	AttributeKeyBlockHeight  = "block_height"
	AttributeKeyCorkID       = "cork_id"
	AttributeKeyDeadline     = "deadline"
	AttributeKeyChainID      = "chain_id"
	AttributeKeyDroppedCount = "dropped_count"

	AttributeValueCategory = ModuleName
)
