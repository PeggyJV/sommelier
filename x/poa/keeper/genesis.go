package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/peggyjv/sommelier/v9/x/poa/types"
)

// InitGenesis loads PoA params and the authority allowlist from genesis.
func InitGenesis(ctx sdk.Context, k Keeper, gs types.GenesisState) {
	k.SetParams(ctx, gs.Params)
	addrs := make([]sdk.ValAddress, 0, len(gs.AuthoritySet))
	for _, v := range gs.AuthoritySet {
		addr, err := sdk.ValAddressFromBech32(v.OperatorAddress)
		if err != nil {
			continue
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
