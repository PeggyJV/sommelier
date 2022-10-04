package keeper

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/peggyjv/sommelier/v4/x/cellarfees/types"
	"github.com/stretchr/testify/require"
)

func TestGettingSettingFeeAccrualCounters(t *testing.T) {
	env := CreateTestEnv(t)
	ctx := env.Context

	expected := types.DefaultFeeAccrualCounters()
	env.cellarFeesKeeper.SetFeeAccrualCounters(ctx, expected)

	require.Equal(t, expected, env.cellarFeesKeeper.GetFeeAccrualCounters(ctx))
}

func TestGettingSettingLastRewardSupplyPeak(t *testing.T) {
	env := CreateTestEnv(t)
	ctx := env.Context

	expected := sdk.NewInt(10 ^ 19 - 1)
	env.cellarFeesKeeper.SetLastRewardSupplyPeak(ctx, expected)

	require.Equal(t, expected, env.cellarFeesKeeper.GetLastRewardSupplyPeak(ctx))
}
