package keeper

import (
	"context"
	"encoding/hex"
	"fmt"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/ethereum/go-ethereum/common"
	gravitytypes "github.com/peggyjv/gravity-bridge/module/v5/x/gravity/types"
	corktypes "github.com/peggyjv/sommelier/v8/x/cork/types"
	types "github.com/peggyjv/sommelier/v8/x/cork/types/v2"
)

var _ types.MsgServer = Keeper{}

// ScheduleCork implements types.MsgServer
func (k Keeper) ScheduleCork(c context.Context, msg *types.MsgScheduleCorkRequest) (*types.MsgScheduleCorkResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)

	signer := msg.MustGetSigner()
	validatorAddr := k.gravityKeeper.GetOrchestratorValidatorAddress(ctx, signer)
	if validatorAddr == nil {
		return nil, errorsmod.Wrapf(sdkerrors.ErrUnauthorized, "signer %s is not a delegate", signer.String())
	}

	params := k.GetParamSet(ctx)
	validatorCorkCount := k.GetValidatorCorkCount(ctx, validatorAddr)
	if validatorCorkCount >= params.MaxCorksPerValidator {
		return nil, corktypes.ErrValidatorCorkCapacityReached
	}

	if !k.HasCellarID(ctx, common.HexToAddress(msg.Cork.TargetContractAddress)) {
		return nil, corktypes.ErrUnmanagedCellarAddress
	}

	if msg.BlockHeight <= uint64(ctx.BlockHeight()) {
		return nil, corktypes.ErrSchedulingInThePast
	}

	corkID := k.SetScheduledCork(ctx, msg.BlockHeight, validatorAddr, *msg.Cork)
	k.IncrementValidatorCorkCount(ctx, validatorAddr)

	invalidationScope := msg.Cork.InvalidationScope()
	// If the vote succeeds, the current invalidation nonce will be incremented
	invalidationNonce := k.GetLatestInvalidationNonce(ctx) + 1

	ctx.EventManager().EmitEvents(
		sdk.Events{
			sdk.NewEvent(
				corktypes.EventTypeCork,
				sdk.NewAttribute(sdk.AttributeKeyModule, corktypes.AttributeValueCategory),
				sdk.NewAttribute(corktypes.AttributeKeySigner, signer.String()),
				sdk.NewAttribute(corktypes.AttributeKeyValidator, validatorAddr.String()),
				sdk.NewAttribute(corktypes.AttributeKeyCork, msg.Cork.String()),
				sdk.NewAttribute(corktypes.AttributeKeyBlockHeight, fmt.Sprintf("%d", msg.BlockHeight)),
				sdk.NewAttribute(corktypes.AttributeKeyCorkID, hex.EncodeToString(corkID)),
				sdk.NewAttribute(gravitytypes.AttributeKeyContractCallInvalidationScope, fmt.Sprint(invalidationScope)),
				sdk.NewAttribute(gravitytypes.AttributeKeyContractCallInvalidationNonce, fmt.Sprint(invalidationNonce)),
			),
		},
	)

	return &types.MsgScheduleCorkResponse{Id: hex.EncodeToString(corkID)}, nil
}
