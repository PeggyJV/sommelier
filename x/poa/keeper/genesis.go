package keeper

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/peggyjv/sommelier/v10/x/poa/types"
)

// InitGenesis loads PoA params, the authority allowlist, the activation stamp,
// the safe-mode flag, and the retained multiplier snapshots.
//
// Everything past params/authority-set exists so an export/import round-trip is
// faithful. A restart that silently reset them would (a) reset the activation
// height, making rawSlashPower treat post-activation infractions as pre-PoA and
// slash them on BOOSTED power, and (b) drop the safe-mode flag, unfreezing the
// value-bearing modules for the blocks before the first EndBlocker re-derives
// it.
func InitGenesis(ctx sdk.Context, k Keeper, gs types.GenesisState) {
	k.SetParams(ctx, gs.Params)

	// Prefer an exported activation stamp; otherwise stamp the current height.
	// Idempotent: the v10 upgrade handler may have already set it before
	// RunMigrations ran this.
	if _, ok := k.GetActivationHeight(ctx); !ok {
		if gs.ActivationHeight > 0 && gs.ActivationHeight <= ctx.BlockHeight() {
			// Never persist a zero timestamp. A genesis carrying
			// activation_height without a stamp (hand-written, or exported
			// before ActivationTimeKey existed) would otherwise make
			// estimateBlockNanos measure from the Unix epoch: the block rate
			// explodes, the retention window collapses to a few blocks, and
			// rawSlashPower then refuses every post-activation authority slash
			// whose snapshot was pruned -- disabling authority slashing with
			// nothing but a log line.
			activationTime := gs.ActivationTime
			if activationTime <= 0 {
				activationTime = ctx.BlockTime().UnixNano()
			}
			k.SetActivationStamp(ctx, gs.ActivationHeight, activationTime)
		} else {
			// Either unset, or a stamp from a height space this chain is not
			// in: an export taken at height H relaunched with a lower
			// initial_height carries H forward, and rawSlashPower would then
			// treat every infraction on the new chain as pre-activation and
			// slash authority validators on BOOSTED power -- the exact
			// inversion its comment says it guards against.
			k.SetActivationHeight(ctx, ctx.BlockHeight())
		}
	}

	addrs := make([]sdk.ValAddress, 0, len(gs.AuthoritySet))
	for _, v := range gs.AuthoritySet {
		addr, err := sdk.ValAddressFromBech32(v.OperatorAddress)
		if err != nil {
			// Fail fast: silently skipping a malformed entry could drop an
			// authority and later trip the authority-empty path (safe mode or
			// halt), or silently run the chain below the supermajority floor.
			panic(fmt.Sprintf("poa InitGenesis: invalid authority address %q: %v", v.OperatorAddress, err))
		}
		addrs = append(addrs, addr)
	}
	k.SetAuthoritySet(ctx, addrs)

	for _, s := range gs.MultiplierSnapshots {
		k.SetMultiplierSnapshot(ctx, s)
	}

	if gs.SafeMode {
		k.SetSafeMode(ctx, true)
		if gs.SafeModeThawHeight > 0 {
			// Clamp into this chain's height space. An export taken at height H
			// while frozen carries thaw=H+delay; relaunching at a lower
			// initial_height leaves a thaw the EndBlocker can never reach, so
			// safe mode never lifts even with a fully healthy authority set --
			// and MsgUpdateParams/MsgUpdateAuthoritySet are themselves frozen,
			// so there is no on-chain exit.
			thaw := gs.SafeModeThawHeight
			if max := ctx.BlockHeight() + safeModeThawDelayBlocks; thaw > max {
				thaw = max
			}
			k.setThawHeight(ctx, thaw)
		}
	}
}

// ExportGenesis returns the PoA module's GenesisState.
func ExportGenesis(ctx sdk.Context, k Keeper) types.GenesisState {
	addrs := k.GetAuthoritySet(ctx)
	out := make([]types.AuthorityValidator, len(addrs))
	for i, a := range addrs {
		out[i] = types.AuthorityValidator{OperatorAddress: a.String()}
	}

	activationHeight, _ := k.GetActivationHeight(ctx)
	activationTime, _ := k.GetActivationTime(ctx)
	thaw, _ := k.getThawHeight(ctx)

	return types.GenesisState{
		Params:              k.GetParams(ctx),
		AuthoritySet:        out,
		ActivationHeight:    activationHeight,
		ActivationTime:      activationTime,
		SafeMode:            k.SafeModeActive(ctx),
		SafeModeThawHeight:  thaw,
		MultiplierSnapshots: k.AllMultiplierSnapshots(ctx),
	}
}
