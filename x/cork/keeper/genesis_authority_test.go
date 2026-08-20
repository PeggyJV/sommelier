package keeper

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"

	_ "github.com/peggyjv/sommelier/v10/app/params"
	corktypes "github.com/peggyjv/sommelier/v10/x/cork/types"
	v2types "github.com/peggyjv/sommelier/v10/x/cork/types/v2"
)

// A scheduled authority cork must survive an export/import round trip.
//
// Before this was fixed, ExportGenesis read only the retired validator-keyed
// queue, so every authority cork was silently dropped on a chain restart, and
// InitGenesis wrote back into a queue that nothing drains or executes.
func TestGenesisRoundTripsAuthorityCorks(t *testing.T) {
	k, ctx, _, ctrl := setupCorkKeeper(t)
	defer ctrl.Finish()

	k.SetParams(ctx, v2types.DefaultParams())

	contract := common.HexToAddress("0x8888888888888888888888888888888888888888")
	cork := v2types.Cork{
		TargetContractAddress: contract.String(),
		EncodedContractCall:   []byte{0x0b, 0xad},
	}
	height := uint64(ctx.BlockHeight()) + 7
	id := k.SetAuthorityCork(ctx, height, cork)

	exported := ExportGenesis(ctx, k)
	require.Len(t, exported.ScheduledCorks, 1, "authority corks must be exported")
	require.Equal(t, height, exported.ScheduledCorks[0].BlockHeight)
	require.Equal(t, id, exported.ScheduledCorks[0].Id)
	require.Equal(t, cork.EncodedContractCall, exported.ScheduledCorks[0].Cork.EncodedContractCall)

	// Fresh keeper, import the exported state.
	k2, ctx2, _, ctrl2 := setupCorkKeeper(t)
	defer ctrl2.Finish()
	InitGenesis(ctx2, k2, exported)

	var restored int
	k2.IterateAuthorityCorksByBlockHeight(ctx2, height, func(_ uint64, gotID []byte, gotContract common.Address, c v2types.Cork) bool {
		restored++
		require.Equal(t, id, gotID)
		require.Equal(t, contract, gotContract)
		require.Equal(t, cork.EncodedContractCall, c.EncodedContractCall)
		return false
	})
	require.Equal(t, 1, restored, "authority cork must be restored into the authority queue")

	// The typed legacy accessors are deleted, so check the raw prefix directly.
	legacy := 0
	legacyIter := sdk.KVStorePrefixIterator(ctx2.KVStore(k2.storeKey), []byte{corktypes.ScheduledCorkKeyPrefix})
	for ; legacyIter.Valid(); legacyIter.Next() {
		legacy++
	}
	legacyIter.Close()
	require.Zero(t, legacy, "import must not resurrect the retired validator-keyed queue")
}
