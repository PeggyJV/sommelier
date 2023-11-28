package v6

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"
	ica "github.com/cosmos/ibc-go/v3/modules/apps/27-interchain-accounts"
	icacontrollertypes "github.com/cosmos/ibc-go/v3/modules/apps/27-interchain-accounts/controller/types"
	icahosttypes "github.com/cosmos/ibc-go/v3/modules/apps/27-interchain-accounts/host/types"
	icatypes "github.com/cosmos/ibc-go/v3/modules/apps/27-interchain-accounts/types"
	transfertypes "github.com/cosmos/ibc-go/v3/modules/apps/transfer/types"
)

func CreateUpgradeHandler(
	mm *module.Manager,
	configurator module.Configurator,
) upgradetypes.UpgradeHandler {
	return func(ctx sdk.Context, plan upgradetypes.Plan, fromVM module.VersionMap) (module.VersionMap, error) {
		ctx.Logger().Info("v6 upgrade: entering handler")

		// set the ICS27 consensus version so InitGenesis is not run
		fromVM[icatypes.ModuleName] = mm.Modules[icatypes.ModuleName].ConsensusVersion()

		// create ICS27 Controller submodule params, controller module not enabled.
		controllerParams := icacontrollertypes.Params{}

		// create ICS27 Host submodule params
		hostParams := icahosttypes.Params{
			HostEnabled: true,
			AllowMessages: []string{
				sdk.MsgTypeURL(&banktypes.MsgSend{}),
				sdk.MsgTypeURL(&banktypes.MsgMultiSend{}),
				sdk.MsgTypeURL(&distrtypes.MsgSetWithdrawAddress{}),
				sdk.MsgTypeURL(&distrtypes.MsgWithdrawDelegatorReward{}),
				sdk.MsgTypeURL(&distrtypes.MsgWithdrawValidatorCommission{}),
				sdk.MsgTypeURL(&distrtypes.MsgFundCommunityPool{}),
				sdk.MsgTypeURL(&govtypes.MsgVote{}),
				sdk.MsgTypeURL(&govtypes.MsgVoteWeighted{}),
				sdk.MsgTypeURL(&stakingtypes.MsgDelegate{}),
				sdk.MsgTypeURL(&stakingtypes.MsgUndelegate{}),
				sdk.MsgTypeURL(&stakingtypes.MsgBeginRedelegate{}),
				sdk.MsgTypeURL(&transfertypes.MsgTransfer{}),
			},
		}

		// initialize ICS27 module
		ctx.Logger().Info("v6 upgrade: initializing ICA module")
		icamodule, correctTypecast := mm.Modules[icatypes.ModuleName].(ica.AppModule)
		if !correctTypecast {
			panic("mm.Modules[icatypes.ModuleName] is not of type ica.AppModule")
		}

		icamodule.InitModule(ctx, controllerParams, hostParams)

		ctx.Logger().Info("v6 upgrade: running migrations and exiting handler")
		return mm.RunMigrations(ctx, configurator, fromVM)
	}
}
