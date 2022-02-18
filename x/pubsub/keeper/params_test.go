package keeper_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	testkeeper "github.com/peggyjv/sommelier/v3/testutil/keeper"
	"github.com/peggyjv/sommelier/v3/x/pubsub/types"
)

func TestGetParams(t *testing.T) {
	k, ctx := testkeeper.PubsubKeeper(t)
	params := types.DefaultParams()

	k.SetParams(ctx, params)

	require.EqualValues(t, params, k.GetParams(ctx))
}
