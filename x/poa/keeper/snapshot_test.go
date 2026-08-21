package keeper

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"

	"github.com/peggyjv/sommelier/v10/x/poa/types"
)

func TestSnapshot_RoundTrip(t *testing.T) {
	k, ctx := NewTestKeeper(t)
	op := sdk.ValAddress([]byte("validator-aaaaaaaaaa")).String()

	snap := types.MultiplierSnapshot{
		Height: 100,
		Entries: []*types.MultiplierEntry{
			{OperatorAddress: op, Multiplier: "2.5"},
		},
	}
	k.SetMultiplierSnapshot(ctx, snap)

	got, ok := k.GetMultiplierSnapshot(ctx, 100)
	require.True(t, ok)
	require.Equal(t, int64(100), got.Height)
	require.Len(t, got.Entries, 1)
	require.Equal(t, op, got.Entries[0].OperatorAddress)
	require.Equal(t, "2.5", got.Entries[0].Multiplier)
}

func TestSnapshot_NotFound(t *testing.T) {
	k, ctx := NewTestKeeper(t)
	_, ok := k.GetMultiplierSnapshot(ctx, 999)
	require.False(t, ok)
}

func TestSnapshot_MultiplierForValidator(t *testing.T) {
	k, ctx := NewTestKeeper(t)
	op := sdk.ValAddress([]byte("validator-aaaaaaaaaa"))
	other := sdk.ValAddress([]byte("validator-bbbbbbbbbb"))

	k.SetMultiplierSnapshot(ctx, types.MultiplierSnapshot{
		Height: 50,
		Entries: []*types.MultiplierEntry{
			{OperatorAddress: op.String(), Multiplier: "3.000000000000000000"},
		},
	})

	m, ok := k.MultiplierForValidatorWithStatus(ctx, op, 50)
	require.True(t, ok)
	require.Equal(t, sdk.MustNewDecFromStr("3.0"), m)

	// snapshot exists, but `other` was not boosted that block: returns 1.
	m, ok = k.MultiplierForValidatorWithStatus(ctx, other, 50)
	require.True(t, ok)
	require.Equal(t, sdk.OneDec(), m)

	// snapshot doesn't exist: returns 1, found=false.
	m, ok = k.MultiplierForValidatorWithStatus(ctx, op, 51)
	require.False(t, ok)
	require.Equal(t, sdk.OneDec(), m)
}

func TestSnapshot_Prune(t *testing.T) {
	k, ctx := NewTestKeeper(t)
	for h := int64(1); h <= 50; h++ {
		k.SetMultiplierSnapshot(ctx, types.MultiplierSnapshot{Height: h})
	}
	k.PruneSnapshotsBefore(ctx, 30)

	for h := int64(1); h < 30; h++ {
		_, ok := k.GetMultiplierSnapshot(ctx, h)
		require.False(t, ok, "snapshot %d should have been pruned", h)
	}
	for h := int64(30); h <= 50; h++ {
		_, ok := k.GetMultiplierSnapshot(ctx, h)
		require.True(t, ok, "snapshot %d should have been retained", h)
	}
}
