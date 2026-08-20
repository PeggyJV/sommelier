package v10

const (
	// UpgradeName is the on-chain upgrade plan name for the v10 release that
	// introduces the PoA power-floor module.
	UpgradeName = "v10"

	// CorkAuthorityAddress is seeded as the sole cork authority for x/cork and
	// x/axelarcork at the v10 upgrade. It replaces the retired
	// validator-supermajority path: this address alone may schedule corks, and
	// for axelarcork also relay, bump gas, and cancel.
	//
	// Rotated afterwards by governance via ParameterChangeProposal, which is
	// the only recovery path if the key is lost or compromised.
	CorkAuthorityAddress = "somm1lcsjy2d5s33h0sddd8lpuqvwyz5ruz7ju4aeqa"
)

// DefaultAuthorityValidators is the binary-specified initial authority
// allowlist seeded by the v10 upgrade handler. Each entry MUST be a valid
// bech32 sdk.ValAddress (sommvaloper1...).
//
// Operators: replace this slice with the production authority validator set
// before tagging the v10 release. The upgrade handler refuses to run if the
// slice is empty (otherwise the chain would enter authority-empty safe mode on
// the first block after upgrade, freezing the value-bearing modules).
var DefaultAuthorityValidators = []string{
	// Sommelier Foundation -- the existing majority validator. Included so the
	// authority bucket keeps a member that already clears the cork approval
	// threshold (>67% of consensus power), preserving the ability to schedule
	// wind-down corks without coordinating all three new nodes. While the
	// Foundation is bonded its raw share already meets FloorFraction, so
	// ComputeMultiplier returns 1 and PoA redistributes nothing; the boost only
	// engages if the Foundation drops out, at which point the three nodes below
	// are raised to the floor.
	"sommvaloper1rtt69afx4dtj4t3urgm93qq7kxypzzeew4w8t0",
	// Sommelier Authority 1 -- sommelier-authority-1, us-east4-b
	"sommvaloper16zxydy6u5ep50dhs987hgq6lqawkcdpfzefs7j",
	// Sommelier Authority 2 -- sommelier-authority-2, europe-west1-b
	"sommvaloper1rntnc0k545976kzgtn8jm749azry8kpwazcptz",
	// Sommelier Authority 3 -- sommelier-authority-3, us-west1-b
	"sommvaloper1n0qvtedunhxqr84gvz0avar0tc4cnc3xm2wks4",
}
