package keeper

import (
	"context"
	"encoding/hex"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/ethereum/go-ethereum/common"

	"github.com/peggyjv/sommelier/v7/x/cork/types"
)

var _ types.MsgServer = Keeper{}

// ScheduleCork implements types.MsgServer
func (k Keeper) ScheduleCork(c context.Context, msg *types.MsgScheduleCorkRequest) (*types.MsgScheduleCorkResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)

	if !k.HasCellarID(ctx, common.HexToAddress(msg.Cork.TargetContractAddress)) {
		return nil, types.ErrUnmanagedCellarAddress
	}

	if msg.BlockHeight <= uint64(ctx.BlockHeight()) {
		return nil, types.ErrSchedulingInThePast
	}

	signer := msg.MustGetSigner()
	validatorAddr := k.gravityKeeper.GetOrchestratorValidatorAddress(ctx, signer)
	if validatorAddr == nil {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrUnauthorized, "signer %s is not a delegate", signer.String())
	}

	if k.GetValidatorCorkCount(ctx, validatorAddr) >= k.GetParamSet(ctx).MaxCorksPerValidator {
		return nil, sdkerrors.Wrapf(types.ErrValidatorCorkCapacityReached, "validator %s has submitted too many uncompleted corks", validatorAddr.String())
	}

	corkID := k.SetScheduledCork(ctx, msg.BlockHeight, validatorAddr, *msg.Cork)
	k.IncrementValidatorCorkCount(ctx, validatorAddr)

	ctx.EventManager().EmitEvents(
		sdk.Events{
			sdk.NewEvent(
				sdk.EventTypeMessage,
				sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			),
			sdk.NewEvent(
				types.EventTypeCork,
				sdk.NewAttribute(types.AttributeKeySigner, signer.String()),
				sdk.NewAttribute(types.AttributeKeyValidator, validatorAddr.String()),
				sdk.NewAttribute(types.AttributeKeyCork, msg.Cork.String()),
				sdk.NewAttribute(types.AttributeKeyBlockHeight, fmt.Sprintf("%d", msg.BlockHeight)),
			),
		},
	)

	return &types.MsgScheduleCorkResponse{Id: hex.EncodeToString(corkID)}, nil
}
