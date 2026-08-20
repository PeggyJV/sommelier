package keeper

import (
	"encoding/hex"
	"fmt"

	"github.com/ethereum/go-ethereum/common"

	sdk "github.com/cosmos/cosmos-sdk/types"
	gravitytypes "github.com/peggyjv/gravity-bridge/module/v6/x/gravity/types"
	"github.com/peggyjv/sommelier/v10/x/cork/types"
	v2types "github.com/peggyjv/sommelier/v10/x/cork/types/v2"
)

// BeginBlocker is called at the beginning of every block
func (k Keeper) BeginBlocker(ctx sdk.Context) {}

func (k Keeper) submitContractCall(ctx sdk.Context, cork v2types.Cork) {
	k.Logger(ctx).Info("setting outgoing tx for contract call",
		"address", cork.TargetContractAddress,
		"encoded contract call", cork.EncodedContractCall)
	// increment invalidation nonce
	invalidationNonce := k.IncrementInvalidationNonce(ctx)
	invalidationScope := cork.InvalidationScope()
	// submit contract call to bridge
	k.gravityKeeper.CreateContractCallTx(
		ctx,
		invalidationNonce,
		invalidationScope,
		common.HexToAddress(cork.TargetContractAddress),
		cork.EncodedContractCall,
		[]gravitytypes.ERC20Token{}, // tokens are always zero
		[]gravitytypes.ERC20Token{},
	)

	ctx.EventManager().EmitEvents(
		sdk.Events{
			sdk.NewEvent(types.EventTypeSubmittedContractCall,
				sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
				sdk.NewAttribute(types.AttributeKeyCork, cork.String()),
				sdk.NewAttribute(types.AttributeKeyBlockHeight, fmt.Sprintf("%d", ctx.BlockHeight())),
				sdk.NewAttribute(types.AttributeKeyCorkID, hex.EncodeToString(cork.IDHash(uint64(ctx.BlockHeight())))),
				sdk.NewAttribute(gravitytypes.AttributeKeyContractCallInvalidationScope, fmt.Sprint(invalidationScope)),
				sdk.NewAttribute(gravitytypes.AttributeKeyContractCallInvalidationNonce, fmt.Sprint(invalidationNonce)),
			),
		},
	)
}

// EndBlocker submits every authority-scheduled cork whose target height is the
// current block, then removes it from the queue.
//
// There is no vote tally: a cork is in the queue only because the cork
// authority put it there, and the authority is sufficient on its own.
func (k Keeper) EndBlocker(ctx sdk.Context) {
	height := uint64(ctx.BlockHeight())

	type dueCork struct {
		id       []byte
		contract common.Address
		cork     v2types.Cork
	}
	var due []dueCork

	// Collect before mutating: deleting through a live iterator is undefined.
	k.IterateAuthorityCorksByBlockHeight(ctx, height, func(_ uint64, id []byte, contract common.Address, cork v2types.Cork) bool {
		due = append(due, dueCork{id: id, contract: contract, cork: cork})
		return false
	})

	if len(due) == 0 {
		return
	}

	// Delete unconditionally, including in safe mode. This is the only site
	// that removes due corks, so skipping deletion would leave every cork whose
	// target height elapsed during a freeze permanently stuck in state.
	for _, d := range due {
		k.DeleteAuthorityCork(ctx, height, d.id, d.contract)
	}

	// PoA authority-empty safe mode: do not submit contract calls under an
	// untrusted, community-only validator set. Corks due inside the freeze are
	// DROPPED, not deferred — deferring would mean executing a call through a
	// bridge secured by a set we no longer trust, at an unbounded later time.
	// They must be re-scheduled after recovery.
	if k.inSafeMode(ctx) {
		k.Logger(ctx).Error("x/poa safe mode active: dropping due authority corks without submitting",
			"height", fmt.Sprintf("%d", ctx.BlockHeight()),
			"dropped", len(due))
		ctx.EventManager().EmitEvent(sdk.NewEvent(
			types.EventTypeCorksDroppedInSafeMode,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(types.AttributeKeyBlockHeight, fmt.Sprintf("%d", ctx.BlockHeight())),
			sdk.NewAttribute(types.AttributeKeyDroppedCount, fmt.Sprintf("%d", len(due))),
		))
		return
	}

	k.Logger(ctx).Info("submitting due authority corks as contract calls",
		"height", fmt.Sprintf("%d", ctx.BlockHeight()), "count", len(due))
	// todo: implement batch sends to save on gas
	for _, d := range due {
		k.submitContractCall(ctx, d.cork)
	}
}
