package keeper

import (
	"context"

	"github.com/peggyjv/sommelier/x/allocation/types"
)

func (k Keeper) DelegateDecisions(c context.Context, msg *types.MsgDelegateDecisions) (*types.MsgDelegateDecisionsResponse, error) {

}

func (k Keeper) DecisionPrecommit(c context.Context, msg *types.MsgDecisionPrecommit) (*types.MsgDecisionPrecommitResponse, error) {

}

func (k Keeper) DecisionCommit(c context.Context, msg *types.MsgDecisionCommit) (*types.MsgDecisionCommitResponse, error) {

}