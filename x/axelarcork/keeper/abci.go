package keeper

import (
	"fmt"

	"github.com/ethereum/go-ethereum/common"

	distributionTypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
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

		currentHeight := uint64(ctx.BlockHeight())
		k.IterateWinningAxelarCorks(ctx, config.Id, func(_ common.Address, blockHeight uint64, cork types.AxelarCork) (stop bool) {
			timeoutHeight := blockHeight + k.GetParamSet(ctx).CorkTimeoutBlocks
			if currentHeight >= timeoutHeight {
				k.Logger(ctx).Info("deleting expired approved scheduled axelar cork",
					"scheduled height", fmt.Sprintf("%d", blockHeight),
					"target contract address", cork.TargetContractAddress)

				k.DeleteWinningAxelarCorkByBlockheight(ctx, config.Id, blockHeight, cork)
			}
			return false
		})

		return false
	})

	// Sweep all axelarcork sender module account balances to the community pool. Because this account is the
	// sender for transfers created by RelayCork calls, funds will not be returned to the caller if the IBC
	// transfer fails or gas is refunded.
	moduleAcct := k.GetSenderAccount(types.ModuleName)
	balances := k.bankKeeper.GetAllBalances(ctx, moduleAcct)
	if err := k.bankKeeper.SendCoinsFromModuleToModule(ctx, types.ModuleName, distributionTypes.ModuleName, balances); err != nil {
		panic(err)
	}

	feePool := k.distributionKeeper.GetFeePool(ctx)
	feePool.CommunityPool = feePool.CommunityPool.Add(sdk.NewDecCoinsFromCoins(balances...)...)
	k.distributionKeeper.SetFeePool(ctx, feePool)
}
