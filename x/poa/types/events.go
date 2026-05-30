package types

const (
	EventTypeAuthorityRescale       = "authority_rescale"
	EventTypeSlashSkippedNoSnapshot = "slash_skipped_no_snapshot"
	EventTypeAuthoritySetUpdated    = "authority_set_updated"
	EventTypeParamsUpdated          = "params_updated"
	// EventTypeSafeModeEntered / Exited mark transitions into and out of
	// authority-empty safe mode (Option A): the chain keeps producing blocks on
	// community stake while value-bearing modules are frozen.
	EventTypeSafeModeEntered = "authority_safe_mode_entered"
	EventTypeSafeModeExited  = "authority_safe_mode_exited"

	AttributeMultiplier       = "multiplier"
	AttributeAuthorityPower   = "authority_power"
	AttributeCommunityPower   = "community_power"
	AttributeOperator         = "operator"
	AttributeInfractionHeight = "infraction_height"
)
