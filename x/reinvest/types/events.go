package types

// reinvest module event types
const (
	EventTypeReinvest = "reinvest"
	EventTypeCommitPeriod        = "commit_period"

	AttributeKeySigner            = "signer"
	AttributeKeyValidator         = "validator"
	AttributeKeyPrevoteHash  = "hash"
	AttributeKeyReinvestment = "reinvestment"
	AttributeKeyCommitPeriodStart   = "commit_period_start"
	AttributeKeyCommitPeriodEnd   = "commit_period_end"

	AttributeValueCategory = ModuleName
)
