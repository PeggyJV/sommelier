package keeper

import (
	"fmt"

	"github.com/ethereum/go-ethereum/common"

	"github.com/peggyjv/sommelier/v7/x/axelarcork/types"

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
		// TODO(bolten): might not be a necessary fix, but GetApprovedScheduledAxelarCorks is kind of
		// unexpectedly destructive -- it deletes all of the validator submitted entries for scheduled
		// corks and returns a list of the winning ones, which are then set using SetWinningAxelarCork
		// here, but it seems like an odd side effect of a Get function to delete stuff
		winningScheduledVotes := k.GetApprovedScheduledAxelarCorks(ctx, config.Id)
		if len(winningScheduledVotes) > 0 {
			k.Logger(ctx).Info("marking all winning scheduled cork votes as relayable",
				"winning votes", winningScheduledVotes,
				"chain id", config.Id)
			for _, c := range winningScheduledVotes {
				k.SetWinningAxelarCork(ctx, config.Id, uint64(ctx.BlockHeight()), c)
			}
		}

		k.Logger(ctx).Info("removing timed out approved corks",
			"height", fmt.Sprintf("%d", ctx.BlockHeight()),
			"chain id", config.Id)

		timeoutHeight := uint64(ctx.BlockHeight()) - k.GetParamSet(ctx).CorkTimeoutBlocks
		k.IterateWinningAxelarCorks(ctx, config.Id, func(_ common.Address, blockHeight uint64, cork types.AxelarCork) (stop bool) {
			if blockHeight >= timeoutHeight {
				k.Logger(ctx).Info("deleting expired approved scheduled axelar cork",
					"scheduled height", fmt.Sprintf("%d", blockHeight),
					"target contract address", cork.TargetContractAddress)

				k.DeleteWinningAxelarCorkByBlockheight(ctx, config.Id, blockHeight, cork)
			}
			return false
		})

		return false
	})
}
