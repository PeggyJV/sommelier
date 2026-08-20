package keeper

import (
	"encoding/hex"
	"fmt"

	"github.com/ethereum/go-ethereum/common"

	"github.com/peggyjv/sommelier/v10/x/axelarcork/types"

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
	// PoA authority-empty safe mode: do not mark corks relayable or sweep funds
	// under an untrusted, community-only validator set.
	//
	// The freeze is applied per-step rather than by returning early. Collection
	// and the timed-out-cork sweep MUST keep running: the collect-and-delete
	// pass below is the only site that removes due corks, so skipping it would
	// strand every cork whose target height elapsed during the freeze.
	safeMode := k.inSafeMode(ctx)
	height := uint64(ctx.BlockHeight())

	k.IterateChainConfigurations(ctx, func(config types.ChainConfiguration) (stop bool) {
		type dueCork struct {
			id       []byte
			contract common.Address
			cork     types.AxelarCork
		}
		var due []dueCork

		// Collect before mutating: deleting through a live iterator is undefined.
		k.IterateAuthorityAxelarCorksByBlockHeight(ctx, config.Id, height,
			func(id []byte, contract common.Address, cork types.AxelarCork) bool {
				due = append(due, dueCork{id: id, contract: contract, cork: cork})
				return false
			})

		// Delete unconditionally, safe mode included, or they are stranded.
		for _, d := range due {
			k.DeleteAuthorityAxelarCork(ctx, config.Id, height, d.id, d.contract)
		}

		winningScheduledVotes := make([]types.AxelarCork, 0, len(due))
		for _, d := range due {
			winningScheduledVotes = append(winningScheduledVotes, d.cork)
		}

		if safeMode {
			// Corks approved during the freeze are DROPPED, not deferred:
			// marking them relayable would let a strategist relay a call
			// approved by a set we no longer trust, at an unbounded later time.
			if len(winningScheduledVotes) > 0 {
				k.Logger(ctx).Error("x/poa safe mode active: dropping due axelar corks without marking them relayable",
					"height", fmt.Sprintf("%d", ctx.BlockHeight()),
					"chain id", config.Id,
					"dropped", len(winningScheduledVotes))
				ctx.EventManager().EmitEvent(sdk.NewEvent(
					types.EventTypeAxelarCorksDroppedInSafeMode,
					sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
					sdk.NewAttribute(types.AttributeKeyBlockHeight, fmt.Sprintf("%d", ctx.BlockHeight())),
					sdk.NewAttribute(types.AttributeKeyChainID, fmt.Sprintf("%d", config.Id)),
					sdk.NewAttribute(types.AttributeKeyDroppedCount, fmt.Sprintf("%d", len(winningScheduledVotes))),
				))
			}
			winningScheduledVotes = nil
		}
		if len(winningScheduledVotes) > 0 {
			k.Logger(ctx).Info("marking due authority corks as relayable",
				"count", len(winningScheduledVotes),
				"chain id", config.Id)
			for _, c := range winningScheduledVotes {
				k.SetWinningAxelarCork(ctx, config.Id, uint64(ctx.BlockHeight()), c)

				ctx.EventManager().EmitEvents(
					sdk.Events{
						sdk.NewEvent(
							types.EventTypeAxelarCorkApproved,
							sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
							sdk.NewAttribute(types.AttributeKeyCork, c.String()),
							sdk.NewAttribute(types.AttributeKeyBlockHeight, fmt.Sprintf("%d", ctx.BlockHeight())),
							sdk.NewAttribute(types.AttributeKeyCorkID, hex.EncodeToString(c.IDHash(uint64(ctx.BlockHeight())))),
						),
					},
				)
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

	// Value-bearing: moves module-account funds. Skipped entirely in safe mode;
	// balances simply accumulate and are swept once the freeze lifts.
	if !safeMode {
		k.SweepModuleAccountBalances(ctx)
	}
}
