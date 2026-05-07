package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/peggyjv/sommelier/v9/x/poa/types"
)

// SetAuthoritySet stores the authority allowlist, replacing any previous set.
func (k Keeper) SetAuthoritySet(ctx sdk.Context, vals []sdk.ValAddress) {
	store := ctx.KVStore(k.storeKey)
	wrapper := types.AuthoritySetWrapper{
		Validators: make([]*types.AuthorityValidator, len(vals)),
	}
	for i, v := range vals {
		wrapper.Validators[i] = &types.AuthorityValidator{OperatorAddress: v.String()}
	}
	store.Set(types.AuthoritySetKey, k.cdc.MustMarshal(&wrapper))
}

// GetAuthoritySet returns the operator addresses in the authority allowlist.
// The order matches the order in which the set was last written.
func (k Keeper) GetAuthoritySet(ctx sdk.Context) []sdk.ValAddress {
	bz := ctx.KVStore(k.storeKey).Get(types.AuthoritySetKey)
	if bz == nil {
		return nil
	}
	var w types.AuthoritySetWrapper
	k.cdc.MustUnmarshal(bz, &w)
	out := make([]sdk.ValAddress, 0, len(w.Validators))
	for _, v := range w.Validators {
		if v == nil {
			continue
		}
		addr, err := sdk.ValAddressFromBech32(v.OperatorAddress)
		if err != nil {
			continue
		}
		out = append(out, addr)
	}
	return out
}

// IsAuthority reports whether `val` is in the current authority allowlist.
func (k Keeper) IsAuthority(ctx sdk.Context, val sdk.ValAddress) bool {
	for _, a := range k.GetAuthoritySet(ctx) {
		if a.Equals(val) {
			return true
		}
	}
	return false
}
