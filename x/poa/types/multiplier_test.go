package types_test

import (
	"testing"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"

	"github.com/peggyjv/sommelier/v9/x/poa/types"
)

func TestComputeMultiplier_BoostNeeded(t *testing.T) {
	floor := sdk.MustNewDecFromStr("0.67")
	auth := math.NewInt(100)
	com := math.NewInt(300)

	m := types.ComputeMultiplier(auth, com, floor)
	require.True(t, m.GT(sdk.OneDec()), "expected boost > 1, got %s", m)

	// Verify the resulting authority share is at or above the floor.
	boosted := sdk.NewDecFromInt(auth).Mul(m)
	total := boosted.Add(sdk.NewDecFromInt(com))
	share := boosted.Quo(total)
	require.True(t, share.GTE(floor), "post-boost share %s below floor %s", share, floor)
}

func TestComputeMultiplier_AlreadyAboveFloor(t *testing.T) {
	floor := sdk.MustNewDecFromStr("0.67")
	m := types.ComputeMultiplier(math.NewInt(900), math.NewInt(100), floor)
	require.Equal(t, sdk.OneDec(), m)
}

func TestComputeMultiplier_AtFloor(t *testing.T) {
	floor := sdk.MustNewDecFromStr("0.5")
	m := types.ComputeMultiplier(math.NewInt(100), math.NewInt(100), floor)
	require.Equal(t, sdk.OneDec(), m)
}

func TestComputeMultiplier_ZeroCommunity(t *testing.T) {
	floor := sdk.MustNewDecFromStr("0.67")
	m := types.ComputeMultiplier(math.NewInt(500), math.ZeroInt(), floor)
	require.Equal(t, sdk.OneDec(), m)
}

func TestComputeMultiplier_ZeroAuthority(t *testing.T) {
	floor := sdk.MustNewDecFromStr("0.67")
	m := types.ComputeMultiplier(math.ZeroInt(), math.NewInt(100), floor)
	require.True(t, m.IsZero())
}

func TestComputeMultiplier_DefaultFloorRoundingSafe(t *testing.T) {
	// Use the production default floor and verify the post-boost share is
	// strictly above 2/3 even after integer truncation in EndBlocker.
	floor := types.DefaultFloorFraction
	auth := math.NewInt(100)
	com := math.NewInt(1000)

	m := types.ComputeMultiplier(auth, com, floor)
	boosted := sdk.NewDecFromInt(auth).Mul(m).TruncateInt() // simulate floor() in EndBlocker
	total := boosted.Add(com)
	share := sdk.NewDecFromInt(boosted).Quo(sdk.NewDecFromInt(total))
	twoThirds := sdk.MustNewDecFromStr("0.666666666666666666")
	require.True(t, share.GT(twoThirds),
		"post-boost share %s should exceed 2/3 (%s)", share, twoThirds)
}
