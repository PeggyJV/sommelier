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

func TestGettingSettingCellarFeePool(t *testing.T) {
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
	require.True(t, pool.Pool.Empty())

	env.cellarFeesKeeper.SetCellarFeePool(ctx, expectedPool)
	pool = env.cellarFeesKeeper.GetCellarFeePool(ctx)
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

	// adding zeroed coin
	poolBefore := env.cellarFeesKeeper.GetCellarFeePool(ctx)
	env.cellarFeesKeeper.AddCoinToPool(ctx, sdk.Coin{})
	poolAfter := env.cellarFeesKeeper.GetCellarFeePool(ctx)
	require.Equal(t, poolBefore, poolAfter)

	// appending new coin
	env.cellarFeesKeeper.AddCoinToPool(ctx, single)
	pool := env.cellarFeesKeeper.GetCellarFeePool(ctx)
	require.Equal(t, sdk.Coins{single}, pool.Pool)

	// adding to existing coin
	env.cellarFeesKeeper.AddCoinToPool(ctx, single)
	pool = env.cellarFeesKeeper.GetCellarFeePool(ctx)
	expected := single.Add(single)
	require.Equal(t, expected, pool.Pool[0])

	// adding zeroed coins
	poolBefore = env.cellarFeesKeeper.GetCellarFeePool(ctx)
	env.cellarFeesKeeper.AddCoinsToPool(ctx, sdk.Coins{})
	poolAfter = env.cellarFeesKeeper.GetCellarFeePool(ctx)
	require.Equal(t, poolBefore, poolAfter)

	denom1 := "gravity-0x1111111111111111111111111111111111111111"
	multi := sdk.Coins{
		{

			Amount: sdk.NewIntFromUint64(100),
			Denom:  denom1,
		},
		{
			Amount: sdk.NewIntFromUint64(100),
			Denom:  denom0,
		},
	}

	// appends one and adds to existing
	env.cellarFeesKeeper.AddCoinsToPool(ctx, multi)
	pool = env.cellarFeesKeeper.GetCellarFeePool(ctx)
	expectedPool := multi.Sort()
	expectedPool[0] = expectedPool[0].Add(single).Add(single)
	require.Equal(t, expectedPool, pool.Pool)
}

func TestSendFeesToAuction(t *testing.T) {
	env := CreateTestEnv(t)
	ctx := env.Context

	// initialize pool
	env.cellarFeesKeeper.SetCellarFeePool(ctx, types.NewEmptyPool())

	denom0 := "gravity-0x0000000000000000000000000000000000000000"
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

	// simulate coins being added to module account
	env.BankKeeper.MintCoins(ctx, types.ModuleName, multi)

	// nothing should happen because pool is empty
	env.cellarFeesKeeper.SendPoolToAuction(ctx)

	// update pool with moduel account balance
	env.cellarFeesKeeper.AddCoinsToPool(ctx, multi)
	env.cellarFeesKeeper.SendPoolToAuction(ctx)

	pool := env.cellarFeesKeeper.GetCellarFeePool(ctx)
	balances := env.BankKeeper.GetAllBalances(ctx, env.AccountKeeper.GetModuleAddress("auction"))
	require.Equal(t, multi, balances)
	require.Equal(t, 0, len(pool.Pool))
}
