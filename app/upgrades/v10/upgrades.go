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
// The handler refuses to proceed when DefaultAuthorityValidators is empty:
// PoA's default Params has HaltWhenAuthorityEmpty=true, so an empty allowlist
// would halt the chain on the very next block (Codex review item 5).
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
					"set before tagging the release, or the chain will halt on the next block " +
					"because Params.HaltWhenAuthorityEmpty defaults to true.",
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
		ctx.Logger().Info("v10 upgrade: PoA params and authority set initialised",
			"authority_count", len(addrs))

		return mm.RunMigrations(ctx, configurator, vm)
	}
}
