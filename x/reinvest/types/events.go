package types

// reinvest module event types
const (
	EventTypereinvestPrecommit = "reinvest_precommit"

	AttributeKeySigner            = "signer"
	AttributeKeyValidator         = "validator"
	AttributeKeyPrevoteHash  = "hash"
	AttributeKeyReinvestment = "reinvestment"
	AttributeKeyCommitPeriodEnd   = "commit_period_end"

	AttributeValueCategory = ModuleName
)
