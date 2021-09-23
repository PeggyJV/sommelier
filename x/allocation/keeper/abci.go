package keeper

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	gravitytypes "github.com/peggyjv/gravity-bridge/module/x/gravity/types"
	"github.com/peggyjv/sommelier/x/allocation/types"
)

// BeginBlocker is called at the beginning of every block
func (k Keeper) BeginBlocker(ctx sdk.Context) {
	// if there is not a vote period set, initialize it with current block height
	// TODO: consider removing
	if !k.HasCommitPeriodStart(ctx) {
		k.SetCommitPeriodStart(ctx, ctx.BlockHeight())
	}

	// On begin block, if we are tallying, emit the new vote period data
	params := k.GetParamSet(ctx)
	votePeriodStart, found := k.GetCommitPeriodStart(ctx)
	if !found {
		panic("vote period not set")
	}

	if (ctx.BlockHeight() - votePeriodStart) < params.VotePeriod {
		return
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeCommitPeriod,
			sdk.NewAttribute(types.AttributeKeyCommitPeriodStart, fmt.Sprintf("%d", votePeriodStart)),
			sdk.NewAttribute(types.AttributeKeyCommitPeriodEnd, fmt.Sprintf("%d", votePeriodStart+params.VotePeriod)),
		),
	)
}

// EndBlocker defines the oracle logic that executes at the end of every block:
//
// 0) Checks if the voting period is over and performs a no-op if it's not.
//
// 1) Checks all the votes submitted in the last period and groups them by ID.
// This step also deletes the oracle votes.
//
// 2) Increments the miss counter for validators that didn't vote
//
// 3) Aggregates the data by ID and type
//
// 4) Compares each submitted data with the aggregated result in order to check
// for reward eligibility
//
// 5) Slashes validators that haven't voted in a while
//
// 6) Deletes all prevotes
//
// 7) Sets the new voting period to the next block

func (k Keeper) EndBlocker(ctx sdk.Context) {
	params := k.GetParamSet(ctx)
	votePeriodStart, found := k.GetCommitPeriodStart(ctx)
	if !found {
		panic("vote period start not set")
	}

	// if the vote period has ended, tally the votes
	periodEnded := (ctx.BlockHeight() - votePeriodStart) >= params.VotePeriod
	if !periodEnded {
		return
	}

	k.Logger(ctx).Info("tallying allocation votes", "height", fmt.Sprintf("%d", ctx.BlockHeight()))
	winningVotes := k.GetWinningVotes(ctx, params.VoteThreshold)

	k.Logger(ctx).Info("package all winning allocation votes into contract calls")
	// todo: implement batch sends to save on gas
	for _, wv := range winningVotes {
		k.Logger(ctx).Info("setting outgoing tx for contract call",
			"cellar", wv.String(),
			"tick ranges length", len(wv.TickRanges))
		contractCall := k.gravityKeeper.CreateContractCallTx(
			ctx,
			k.IncrementInvalidationNonce(ctx),
			wv.InvalidationScope(),
			wv.ABIEncodedRebalanceHash(),
			[]gravitytypes.ERC20Token{}, // tokens are always zero
			[]gravitytypes.ERC20Token{})
		k.gravityKeeper.SetOutgoingTx(ctx, contractCall)
	}

	// Reset state prior to next round
	k.DeleteAllPrecommits(ctx)

	// After the tallying is done, reset the vote period start to the next block
	votePeriodStart = ctx.BlockHeight() + 1
	k.SetCommitPeriodStart(ctx, votePeriodStart)

	k.Logger(ctx).Info("vote period set", "height", fmt.Sprintf("%d", votePeriodStart))
}
