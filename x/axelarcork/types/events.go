package types

// axelarcork module event types
const (
	EventTypeAxelarCorkApproved    = "axelar_cork_approved"
	EventTypeAxelarCorkRelayCalled = "axelar_cork_relay_called"

	AttributeKeyCork        = "cork"
	AttributeKeyBlockHeight = "block_height"
	AttributeKeyCorkID      = "cork_id"
	AttributeKeyDeadline    = "deadline"

	AttributeValueCategory = ModuleName
)
