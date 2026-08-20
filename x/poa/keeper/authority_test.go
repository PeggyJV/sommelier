package keeper_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"

	"github.com/peggyjv/sommelier/v10/x/poa/keeper"
)

func valAddr(t *testing.T, hex string) sdk.ValAddress {
	t.Helper()
	bz := []byte(hex)
	require.Len(t, bz, 20)
	return sdk.ValAddress(bz)
}

func TestAuthoritySet_RoundTrip(t *testing.T) {
	k, ctx := keeper.NewTestKeeper(t)
	a := valAddr(t, "validator-aaaaaaaaaa")
	b := valAddr(t, "validator-bbbbbbbbbb")
	k.SetAuthoritySet(ctx, []sdk.ValAddress{a, b})

	got := k.GetAuthoritySet(ctx)
	require.Len(t, got, 2)
	require.True(t, got[0].Equals(a))
	require.True(t, got[1].Equals(b))

	require.True(t, k.IsAuthority(ctx, a))
	require.True(t, k.IsAuthority(ctx, b))
	require.False(t, k.IsAuthority(ctx, valAddr(t, "outsider-xxxxxxxxxxx")))
}

func TestAuthoritySet_Empty(t *testing.T) {
	k, ctx := keeper.NewTestKeeper(t)
	require.Empty(t, k.GetAuthoritySet(ctx))
	require.False(t, k.IsAuthority(ctx, valAddr(t, "anyone-xxxxxxxxxxxxx")))
}

func TestAuthoritySet_Replace(t *testing.T) {
	k, ctx := keeper.NewTestKeeper(t)
	a := valAddr(t, "validator-aaaaaaaaaa")
	b := valAddr(t, "validator-bbbbbbbbbb")
	c := valAddr(t, "validator-cccccccccc")

	k.SetAuthoritySet(ctx, []sdk.ValAddress{a, b})
	k.SetAuthoritySet(ctx, []sdk.ValAddress{c}) // replace, not merge

	got := k.GetAuthoritySet(ctx)
	require.Len(t, got, 1)
	require.True(t, got[0].Equals(c))
	require.False(t, k.IsAuthority(ctx, a))
}
