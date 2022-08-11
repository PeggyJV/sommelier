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
