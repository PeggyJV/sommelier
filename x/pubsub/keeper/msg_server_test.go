package keeper_test

import (
	"context"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
    "github.com/peggyjv/sommelier/v3/x/pubsub/types"
    "github.com/peggyjv/sommelier/v3/x/pubsub/keeper"
    keepertest "github.com/peggyjv/sommelier/v3/testutil/keeper"
)

func setupMsgServer(t testing.TB) (types.MsgServer, context.Context) {
	k, ctx := keepertest.PubsubKeeper(t)
	return keeper.NewMsgServerImpl(*k), sdk.WrapSDKContext(ctx)
}
