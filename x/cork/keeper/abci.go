package keeper

import (
	"encoding/hex"
	"fmt"

	"github.com/ethereum/go-ethereum/common"

	sdk "github.com/cosmos/cosmos-sdk/types"
	gravitytypes "github.com/peggyjv/gravity-bridge/module/v5/x/gravity/types"
	"github.com/peggyjv/sommelier/v8/x/cork/types"
	v2types "github.com/peggyjv/sommelier/v8/x/cork/types/v2"
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

// EndBlocker defines the oracle logic that executes at the end of every block:
//
// 1) Collects all winning votes
//
// 2) Submits all winning votes as contract calls via the gravity bridge

func (k Keeper) EndBlocker(ctx sdk.Context) {
	k.Logger(ctx).Info("tallying scheduled cork votes", "height", fmt.Sprintf("%d", ctx.BlockHeight()))
	winningScheduledVotes := k.GetApprovedScheduledCorks(ctx)
	if len(winningScheduledVotes) > 0 {
		k.Logger(ctx).Info("packaging all winning scheduled cork votes into contract calls",
			"winning votes", winningScheduledVotes)
		// todo: implement batch sends to save on gas
		for _, wv := range winningScheduledVotes {
			k.submitContractCall(ctx, wv)
		}
	}
}
