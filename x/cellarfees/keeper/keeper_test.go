package keeper

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/peggyjv/sommelier/v4/x/cellarfees/types"
	"github.com/stretchr/testify/require"
)

func TestGettingAndSettingCellarFeePool(t *testing.T) {
	env := CreateTestEnv(t)
	ctx := env.Context

	expectedPool := types.CellarFeePool{
		Pool: sdk.Coins{
			{
				Amount: sdk.NewIntFromUint64(100),
				Denom:  "gravity-0x0000000000000000000000000000000000000000",
			},
		},
	}

	env.cellarFeesKeeper.SetCellarFeePool(ctx, expectedPool)
	pool := env.cellarFeesKeeper.GetCellarFeePool(ctx)

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
	require.Equal(t, pool.Pool[0], single)

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
	require.Equal(t, pool.Pool[0], single.Add(multi[0]))
	require.Equal(t, pool.Pool[1], multi[1])
}
