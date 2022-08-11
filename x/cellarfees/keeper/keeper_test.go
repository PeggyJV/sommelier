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
	env.cellarFeesKeeper.SetCellarFeePool(ctx, types.NewEmptyPool())
	pool := env.cellarFeesKeeper.GetCellarFeePool(ctx).Pool
	require.Equal(t, len(pool), 0)
}

func TestGettingSettingSearchingCellarFeePool(t *testing.T) {
	env := CreateTestEnv(t)
	ctx := env.Context

	// initialize pool
	env.cellarFeesKeeper.SetCellarFeePool(ctx, types.NewEmptyPool())

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
	i, balance := env.cellarFeesKeeper.FindBalanceInPool(ctx, denom, pool.Pool)
	require.Equal(t, -1, i)
	require.Zero(t, balance)
	require.Equal(t, types.CellarFeePool{}, pool)

	env.cellarFeesKeeper.SetCellarFeePool(ctx, expectedPool)
	pool = env.cellarFeesKeeper.GetCellarFeePool(ctx)
	i, _ = env.cellarFeesKeeper.FindBalanceInPool(ctx, denom, pool.Pool)
	require.Equal(t, 0, i)
	require.Equal(t, expectedPool, pool)
}

func TestAddingCoinsToPool(t *testing.T) {
	env := CreateTestEnv(t)
	ctx := env.Context

	// initialize pool
	env.cellarFeesKeeper.SetCellarFeePool(ctx, types.NewEmptyPool())

	denom0 := "gravity-0x0000000000000000000000000000000000000000"
	single := sdk.Coin{
		Amount: sdk.NewIntFromUint64(100),
		Denom:  denom0,
	}

	env.cellarFeesKeeper.AddCoinToPool(ctx, single)
	pool := env.cellarFeesKeeper.GetCellarFeePool(ctx)
	require.Equal(t, single, pool.Pool[0])

	env.cellarFeesKeeper.AddCoinToPool(ctx, single)
	pool = env.cellarFeesKeeper.GetCellarFeePool(ctx)
	expected := single.Add(single)
	require.Equal(t, expected, pool.Pool[0])

	denom1 := "gravity-0x1111111111111111111111111111111111111111"
	multi := sdk.Coins{
		{
			Amount: sdk.NewIntFromUint64(100),
			Denom:  denom0,
		},
		{

			Amount: sdk.NewIntFromUint64(100),
			Denom:  denom1,
		},
	}

	env.cellarFeesKeeper.AddCoinsToPool(ctx, multi)
	pool = env.cellarFeesKeeper.GetCellarFeePool(ctx)
	expected = expected.Add(multi[0])
	require.Equal(t, 2, len(pool.Pool))
	require.Equal(t, expected, pool.Pool[0])
	require.Equal(t, multi[1], pool.Pool[1])
}
