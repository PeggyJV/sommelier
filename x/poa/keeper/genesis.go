package keeper

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/peggyjv/sommelier/v9/x/poa/types"
)

// InitGenesis loads PoA params and the authority allowlist from genesis.
func InitGenesis(ctx sdk.Context, k Keeper, gs types.GenesisState) {
	k.SetParams(ctx, gs.Params)
	// Record the activation height once. This is the first height at which PoA
	// boosting is in effect; infraction heights below it predate the module and
	// carry no boost, so a missing snapshot there is benign. Idempotent: the
	// v10 upgrade handler may have already set it before RunMigrations ran this.
	if _, ok := k.GetActivationHeight(ctx); !ok {
		k.SetActivationHeight(ctx, ctx.BlockHeight())
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
}

// ExportGenesis returns the PoA module's GenesisState.
func ExportGenesis(ctx sdk.Context, k Keeper) types.GenesisState {
	addrs := k.GetAuthoritySet(ctx)
	out := make([]types.AuthorityValidator, len(addrs))
	for i, a := range addrs {
		out[i] = types.AuthorityValidator{OperatorAddress: a.String()}
	}
	return types.GenesisState{
		Params:       k.GetParams(ctx),
		AuthoritySet: out,
	}
}
