package keeper

import (
	"github.com/peggyjv/sommelier/v3/x/pubsub/types"
)

var _ types.QueryServer = Keeper{}
