package keeper

import (
	"fmt"
	"github.com/peggyjv/sommelier/v6/x/axelarcork/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// BeginBlocker is called at the beginning of every block
func (k Keeper) BeginBlocker(ctx sdk.Context) {}

// EndBlocker defines the oracle logic that executes at the end of every block:
//
// 1) Collects all winning votes
//
// 2) Stores all winning votes as corks that strategists are allowed to relay via Axelar

func (k Keeper) EndBlocker(ctx sdk.Context) {
	k.IterateChainConfigurations(ctx, func(config types.ChainConfiguration) (stop bool) {
		k.Logger(ctx).Info("tallying scheduled cork votes",
			"height", fmt.Sprintf("%d", ctx.BlockHeight()),
			"chain id", config.Id)
		winningScheduledVotes := k.GetApprovedScheduledCorks(ctx, config.Id)
		if len(winningScheduledVotes) > 0 {
			k.Logger(ctx).Info("marking all winning scheduled cork votes as relayable",
				"winning votes", winningScheduledVotes,
				"chain id", config.Id)
			for _, c := range winningScheduledVotes {
				k.SetWinningCork(ctx, config.Id, c)
			}
		}

		return false
	})
}
