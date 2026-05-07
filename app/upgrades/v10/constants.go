package v10

const (
	// UpgradeName is the on-chain upgrade plan name for the v10 release that
	// introduces the PoA power-floor module.
	UpgradeName = "v10"
)

// DefaultAuthorityValidators is the binary-specified initial authority
// allowlist seeded by the v10 upgrade handler. Each entry MUST be a valid
// bech32 sdk.ValAddress (sommvaloper1...).
//
// Operators: replace this slice with the production authority validator set
// before tagging the v10 release. The upgrade handler refuses to run if the
// slice is empty (otherwise the chain would halt on the first block after
// upgrade because Params.HaltWhenAuthorityEmpty defaults to true).
var DefaultAuthorityValidators = []string{
	// "sommvaloper1...",
}
