package keeper

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

// estimateBlockNanos previously was a compiled-in 4s constant, justified as a
// "conservative lower bound" on block time. That argument only holds if the
// assumed value really is a lower bound: if actual block times fall below it,
// the derived retention window under-covers the slashable period, snapshots for
// still-slashable heights get pruned, and authority slashes are then silently
// REFUSED by rawSlashPower. These tests pin the self-calibrating replacement.
func TestEstimateBlockNanos(t *testing.T) {
	sec := int64(time.Second)

	t.Run("falls back before the activation stamp exists", func(t *testing.T) {
		k, ctx := NewTestKeeper(t)
		require.Equal(t, fallbackBlockNanos, k.estimateBlockNanos(ctx))
	})

	t.Run("falls back until the sample is large enough", func(t *testing.T) {
		k, ctx := NewTestKeeper(t)
		k.SetActivationStamp(ctx, 0, 0)
		// Only 10 blocks of history: not enough to trust.
		ctx = ctx.WithBlockHeight(10).WithBlockTime(time.Unix(60, 0))
		require.Equal(t, fallbackBlockNanos, k.estimateBlockNanos(ctx))
	})

	t.Run("measures the real rate and shades it down", func(t *testing.T) {
		k, ctx := NewTestKeeper(t)
		k.SetActivationStamp(ctx, 0, 0)
		// 2000 blocks in 12000s => 6s/block, shaded by 4/5 => 4.8s.
		ctx = ctx.WithBlockHeight(2000).WithBlockTime(time.Unix(12000, 0))
		require.Equal(t, 4800*int64(time.Millisecond), k.estimateBlockNanos(ctx))
	})

	t.Run("tracks block times below the old hardcoded 4s bound", func(t *testing.T) {
		k, ctx := NewTestKeeper(t)
		k.SetActivationStamp(ctx, 0, 0)
		// 2000 blocks in 2000s => 1s/block. The old constant would have assumed
		// 4s and retained a quarter of the snapshots actually needed.
		ctx = ctx.WithBlockHeight(2000).WithBlockTime(time.Unix(2000, 0))
		got := k.estimateBlockNanos(ctx)
		require.Equal(t, 800*int64(time.Millisecond), got)
		require.Less(t, got, 4*sec, "must not assume the stale 4s lower bound")
	})

	t.Run("floors at minBlockNanos", func(t *testing.T) {
		k, ctx := NewTestKeeper(t)
		k.SetActivationStamp(ctx, 0, 0)
		// Absurdly fast (or a clock anomaly): 1_000_000 blocks in 1s.
		ctx = ctx.WithBlockHeight(1_000_000).WithBlockTime(time.Unix(1, 0))
		require.Equal(t, minBlockNanos, k.estimateBlockNanos(ctx))
	})

	t.Run("ignores non-monotonic clocks", func(t *testing.T) {
		k, ctx := NewTestKeeper(t)
		k.SetActivationStamp(ctx, 0, int64(time.Hour))
		ctx = ctx.WithBlockHeight(2000).WithBlockTime(time.Unix(0, 0))
		require.Equal(t, fallbackBlockNanos, k.estimateBlockNanos(ctx))
	})
}

// A faster chain must retain MORE snapshots for the same unbonding period.
func TestEstimateBlockNanos_FasterChainRetainsMore(t *testing.T) {
	unbonding := 21 * 24 * time.Hour

	blocksFor := func(height int64, elapsed time.Duration) int64 {
		k, ctx := NewTestKeeper(t)
		k.SetActivationStamp(ctx, 0, 0)
		ctx = ctx.WithBlockHeight(height).WithBlockTime(time.Unix(0, int64(elapsed)))
		bn := k.estimateBlockNanos(ctx)
		return (unbonding.Nanoseconds()+bn-1)/bn + 1
	}

	slow := blocksFor(2000, 12000*time.Second) // 6s/block
	fast := blocksFor(2000, 2000*time.Second)  // 1s/block
	require.Greater(t, fast, slow)
}
