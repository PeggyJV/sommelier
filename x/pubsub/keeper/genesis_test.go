package keeper

import (
	"testing"

	"github.com/peggyjv/sommelier/v9/x/pubsub/types"
	"github.com/stretchr/testify/require"
)

// TODO(bolten): fill out genesis test
func TestGenesis(t *testing.T) {
	genesisState := types.GenesisState{
		Params: types.DefaultParams(),
	}

	require.NotNil(t, genesisState)
}
