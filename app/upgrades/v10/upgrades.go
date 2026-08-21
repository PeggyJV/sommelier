package v10

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"

	storetypes "github.com/cosmos/cosmos-sdk/store/types"

	axelarcorkkeeper "github.com/peggyjv/sommelier/v10/x/axelarcork/keeper"
	corkkeeper "github.com/peggyjv/sommelier/v10/x/cork/keeper"
	poakeeper "github.com/peggyjv/sommelier/v10/x/poa/keeper"
	poatypes "github.com/peggyjv/sommelier/v10/x/poa/types"
)

// CreateUpgradeHandler builds the v10 upgrade handler. It runs module
// migrations and seeds the PoA module's params and authority allowlist from
// DefaultAuthorityValidators.
//
// The handler refuses to proceed when DefaultAuthorityValidators is empty: with
// the default params (authority-empty safe mode), an empty allowlist would put
// the chain into safe mode on the very next block — value-bearing modules
// (gravity/cork/axelarcork) frozen — which is not a usable production state
// (Codex review item 5).
func CreateUpgradeHandler(
	mm *module.Manager,
	configurator module.Configurator,
	poaKeeper poakeeper.Keeper,
	corkKeeper corkkeeper.Keeper,
	axelarcorkKeeper axelarcorkkeeper.Keeper,
	corkStoreKey storetypes.StoreKey,
	axelarcorkStoreKey storetypes.StoreKey,
) upgradetypes.UpgradeHandler {
	return func(ctx sdk.Context, _ upgradetypes.Plan, vm module.VersionMap) (module.VersionMap, error) {
		ctx.Logger().Info("v10 upgrade: entering handler")

		// Fail fast on a malformed constant: a bad authority would be written to
		// params and, being fail-closed, would make cork scheduling impossible
		// until a governance proposal completed.
		if _, err := sdk.AccAddressFromBech32(CorkAuthorityAddress); err != nil {
			return vm, fmt.Errorf("v10: invalid CorkAuthorityAddress %q: %w", CorkAuthorityAddress, err)
		}

		if len(DefaultAuthorityValidators) == 0 {
			return vm, fmt.Errorf(
				"v10 upgrade refuses to run: DefaultAuthorityValidators is empty. " +
					"Populate app/upgrades/v10/constants.go with the production authority validator " +
					"set before tagging the release, or the chain will enter safe mode on the next " +
					"block (value-bearing modules frozen) because the authority set is empty.",
			)
		}

		// Migrations MUST run before the seeding below.
		//
		// x/poa is new in v10, so it is absent from the on-chain version map and
		// RunMigrations therefore calls its InitGenesis with DefaultGenesis, whose
		// AuthoritySet is nil. poa InitGenesis calls SetAuthoritySet
		// unconditionally, so seeding first meant the set was wiped moments later.
		// The chain would then enter authority-empty safe mode on the very next
		// block -- gravity/cork/axelarcork frozen -- with MsgUpdateAuthoritySet and
		// MsgUpdateParams both rejected by ErrSafeModeGovFrozen, i.e. no on-chain
		// recovery. Guarded by
		// TestV10UpgradeHandlerKeepsAuthoritySetThroughRunMigrations.
		vm, err := mm.RunMigrations(ctx, configurator, vm)
		if err != nil {
			return vm, err
		}

		addrs := make([]sdk.ValAddress, 0, len(DefaultAuthorityValidators))
		for _, s := range DefaultAuthorityValidators {
			addr, err := sdk.ValAddressFromBech32(s)
			if err != nil {
				return vm, fmt.Errorf("v10: invalid authority validator %q: %w", s, err)
			}
			addrs = append(addrs, addr)
		}

		// Refuse to seed an authority set with no usable validator.
		//
		// The empty-slice guard above catches a forgotten constant, but not a
		// valid-looking typo, a stale operator address, or every listed
		// validator being jailed or unbonded at the upgrade height. Any of
		// those seeds an allowlist whose members carry no power, the next
		// EndBlocker sees authPower == 0 and enters safe mode, and both
		// MsgUpdateAuthoritySet and MsgUpdateParams are frozen there -- no
		// on-chain recovery. MsgUpdateAuthoritySet already refuses such a set
		// (ErrNoBondedAuthority); the upgrade path must too.
		//
		// Failing here halts the chain at the upgrade height, which operators
		// can fix and restart. Proceeding cannot be fixed on-chain at all.
		sk := poaKeeper.StakingKeeper()
		live := 0
		for _, addr := range addrs {
			v, found := sk.GetValidator(ctx, addr)
			if found && !v.Jailed && v.IsBonded() {
				live++
			}
		}
		if live == 0 {
			return vm, fmt.Errorf(
				"v10 upgrade refuses to run: none of the %d configured authority validators "+
					"is bonded and unjailed at this height. Seeding them would enter "+
					"authority-empty safe mode on the next block, freezing gravity/cork/"+
					"axelarcork with no on-chain recovery (poa governance msgs are frozen "+
					"in safe mode). Check app/upgrades/v10/constants.go against the live "+
					"validator set", len(addrs))
		}
		ctx.Logger().Info("v10 upgrade: authority validator liveness verified",
			"configured", len(addrs), "bonded_and_unjailed", live)

		poaKeeper.SetParams(ctx, poatypes.DefaultParams())
		poaKeeper.SetAuthoritySet(ctx, addrs)
		// Record the activation height: the first block at which PoA boosting is
		// in effect. Slashes for infraction heights below this predate the
		// module (no boost) and pass through; at or above it a missing snapshot
		// indicates corruption and the slash is refused.
		poaKeeper.SetActivationHeight(ctx, ctx.BlockHeight())
		ctx.Logger().Info("v10 upgrade: PoA params and authority set initialised",
			"authority_count", len(addrs), "activation_height", ctx.BlockHeight())

		// Seed the cork authority in both modules.
		//
		// Write the single new key rather than round-tripping the whole param
		// set. cork_authority does not exist in v9 state, and GetParamSet reads
		// every registered pair, so it panics on the missing key
		// ("UnmarshalJSON cannot decode empty bytes") before returning --
		// halting the chain inside this BeginBlocker. Verified against a
		// mainnet-state rehearsal, where the round-trip form did exactly that.
		corkKeeper.SetCorkAuthority(ctx, CorkAuthorityAddress)
		axelarcorkKeeper.SetCorkAuthority(ctx, CorkAuthorityAddress)

		// Drain the retired validator-scheduled queues. MUST happen here: v10
		// removes the power tally that was the only deleter of these prefixes,
		// so anything left behind is permanently undeletable state.
		drained := DrainLegacyCorkQueues(ctx, corkStoreKey, axelarcorkStoreKey)
		ctx.Logger().Info("v10 upgrade: cork authority seeded and legacy queues drained",
			"authority", CorkAuthorityAddress, "drained_keys", drained)

		return vm, nil
	}
}
