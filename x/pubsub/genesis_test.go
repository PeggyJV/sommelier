package pubsub_test

import (
	"testing"

	keepertest "github.com/peggyjv/sommelier/v3/testutil/keeper"
	"github.com/peggyjv/sommelier/v3/testutil/nullify"
	"github.com/peggyjv/sommelier/v3/x/pubsub"
	"github.com/peggyjv/sommelier/v3/x/pubsub/types"
	"github.com/stretchr/testify/require"
)

func TestGenesis(t *testing.T) {
	genesisState := types.GenesisState{
		Params:	types.DefaultParams(),
		
		// this line is used by starport scaffolding # genesis/test/state
	}

	k, ctx := keepertest.PubsubKeeper(t)
	pubsub.InitGenesis(ctx, *k, genesisState)
	got := pubsub.ExportGenesis(ctx, *k)
	require.NotNil(t, got)

	nullify.Fill(&genesisState)
	nullify.Fill(got)

	

	// this line is used by starport scaffolding # genesis/test/assert
}
