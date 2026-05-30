package v10

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"

	poakeeper "github.com/peggyjv/sommelier/v9/x/poa/keeper"
	poatypes "github.com/peggyjv/sommelier/v9/x/poa/types"
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
) upgradetypes.UpgradeHandler {
	return func(ctx sdk.Context, _ upgradetypes.Plan, vm module.VersionMap) (module.VersionMap, error) {
		ctx.Logger().Info("v10 upgrade: entering handler")

		if len(DefaultAuthorityValidators) == 0 {
			return vm, fmt.Errorf(
				"v10 upgrade refuses to run: DefaultAuthorityValidators is empty. " +
					"Populate app/upgrades/v10/constants.go with the production authority validator " +
					"set before tagging the release, or the chain will enter safe mode on the next " +
					"block (value-bearing modules frozen) because the authority set is empty.",
			)
		}

		addrs := make([]sdk.ValAddress, 0, len(DefaultAuthorityValidators))
		for _, s := range DefaultAuthorityValidators {
			addr, err := sdk.ValAddressFromBech32(s)
			if err != nil {
				return vm, fmt.Errorf("v10: invalid authority validator %q: %w", s, err)
			}
			addrs = append(addrs, addr)
		}

		poaKeeper.SetParams(ctx, poatypes.DefaultParams())
		poaKeeper.SetAuthoritySet(ctx, addrs)
		// Record the activation height: the first block at which PoA boosting is
		// in effect. Slashes for infraction heights below this predate the
		// module (no boost) and pass through; at or above it a missing snapshot
		// indicates corruption and the slash is refused.
		poaKeeper.SetActivationHeight(ctx, ctx.BlockHeight())
		ctx.Logger().Info("v10 upgrade: PoA params and authority set initialised",
			"authority_count", len(addrs), "activation_height", ctx.BlockHeight())

		return mm.RunMigrations(ctx, configurator, vm)
	}
}
