package keeper

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/peggyjv/sommelier/v4/x/cellarfees/types"
	"github.com/stretchr/testify/require"
)

func TestInitializingPool(t *testing.T) {
	env := CreateTestEnv(t)
	ctx := env.Context

	// initialize pool
	env.cellarFeesKeeper.SetCellarFeePool(ctx, types.DefaultCellarFeePool())
	pool := env.cellarFeesKeeper.GetCellarFeePool(ctx).Pool
	require.Equal(t, len(pool), 0)
}

func TestGettingSettingCellarFeePool(t *testing.T) {
	env := CreateTestEnv(t)
	ctx := env.Context

	// initialize pool
	env.cellarFeesKeeper.SetCellarFeePool(ctx, types.DefaultCellarFeePool())

	denom := "gravity-0x0000000000000000000000000000000000000000"
	expectedPool := types.CellarFeePool{
		Pool: sdk.Coins{
			{
				Amount: sdk.NewIntFromUint64(100),
				Denom:  denom,
			},
		},
	}

	pool := env.cellarFeesKeeper.GetCellarFeePool(ctx)
	require.True(t, pool.Pool.Empty())

	env.cellarFeesKeeper.SetCellarFeePool(ctx, expectedPool)
	pool = env.cellarFeesKeeper.GetCellarFeePool(ctx)
	require.Equal(t, expectedPool, pool)
}

func TestAddingCoinsToPool(t *testing.T) {
	env := CreateTestEnv(t)
	ctx := env.Context

	// initialize pool
	env.cellarFeesKeeper.SetCellarFeePool(ctx, types.DefaultCellarFeePool())

	denom0 := "gravity-0x0000000000000000000000000000000000000000"
	single := sdk.Coin{
		Amount: sdk.NewIntFromUint64(100),
		Denom:  denom0,
	}

	// adding zeroed coin
	poolBefore := env.cellarFeesKeeper.GetCellarFeePool(ctx)
	env.cellarFeesKeeper.addCoinToPool(ctx, sdk.Coin{})
	poolAfter := env.cellarFeesKeeper.GetCellarFeePool(ctx)
	require.Equal(t, poolBefore, poolAfter)

	// appending new coin
	env.cellarFeesKeeper.addCoinToPool(ctx, single)
	pool := env.cellarFeesKeeper.GetCellarFeePool(ctx)
	require.Equal(t, sdk.Coins{single}, pool.Pool)

	// adding to existing coin
	env.cellarFeesKeeper.addCoinToPool(ctx, single)
	pool = env.cellarFeesKeeper.GetCellarFeePool(ctx)
	expected := single.Add(single)
	require.Equal(t, expected, pool.Pool[0])
}

func TestGettingSettingLastRewardSupplyPeak(t *testing.T) {
	env := CreateTestEnv(t)
	ctx := env.Context

	expected := sdk.NewInt(10 ^ 19 - 1)
	env.cellarFeesKeeper.SetLastRewardSupplyPeak(ctx, expected)

	require.Equal(t, expected, env.cellarFeesKeeper.GetLastRewardSupplyPeak(ctx))
}

func TestGettingSettingScheduledAuctionHeight(t *testing.T) {
	env := CreateTestEnv(t)
	ctx := env.Context

	expected := uint64(10000)
	env.cellarFeesKeeper.SetScheduledAuctionHeight(ctx, expected)

	require.Equal(t, expected, env.cellarFeesKeeper.GetScheduledAuctionHeight(ctx))
}
