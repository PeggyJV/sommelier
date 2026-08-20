package keeper

import (
	"encoding/hex"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"

	_ "github.com/peggyjv/sommelier/v10/app/params"
	v2types "github.com/peggyjv/sommelier/v10/x/cork/types/v2"
)

// Operators drive the wind-down through these queries. If they read the retired
// validator-keyed queue they return empty, leaving no way to see what is
// actually scheduled.
func TestScheduledCorkQueriesReadAuthorityQueue(t *testing.T) {
	k, ctx, _, ctrl := setupCorkKeeper(t)
	defer ctrl.Finish()

	k.SetParams(ctx, v2types.DefaultParams())

	contract := common.HexToAddress("0xabababababababababababababababababababab")
	cork := v2types.Cork{
		TargetContractAddress: contract.String(),
		EncodedContractCall:   []byte{0xfe, 0xed},
	}
	height := uint64(ctx.BlockHeight()) + 3
	id := k.SetAuthorityCork(ctx, height, cork)

	goCtx := sdk.WrapSDKContext(ctx)

	all, err := k.QueryScheduledCorks(goCtx, &v2types.QueryScheduledCorksRequest{})
	require.NoError(t, err)
	require.Len(t, all.Corks, 1, "QueryScheduledCorks must see authority corks")

	byHeight, err := k.QueryScheduledCorksByBlockHeight(goCtx, &v2types.QueryScheduledCorksByBlockHeightRequest{
		BlockHeight: height,
	})
	require.NoError(t, err)
	require.Len(t, byHeight.Corks, 1, "QueryScheduledCorksByBlockHeight must see authority corks")
	require.Equal(t, height, byHeight.Corks[0].BlockHeight)

	byID, err := k.QueryScheduledCorksByID(goCtx, &v2types.QueryScheduledCorksByIDRequest{
		Id: hex.EncodeToString(id),
	})
	require.NoError(t, err)
	require.Len(t, byID.Corks, 1, "QueryScheduledCorksByID must see authority corks")
	require.Equal(t, id, byID.Corks[0].Id)

	// A height with nothing queued returns empty, not the whole queue.
	none, err := k.QueryScheduledCorksByBlockHeight(goCtx, &v2types.QueryScheduledCorksByBlockHeightRequest{
		BlockHeight: height + 1,
	})
	require.NoError(t, err)
	require.Empty(t, none.Corks)
}
