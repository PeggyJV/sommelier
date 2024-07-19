package keeper

import (
	"context"
	"encoding/hex"
	"fmt"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/ethereum/go-ethereum/common"
	corktypes "github.com/peggyjv/sommelier/v7/x/cork/types"
	types "github.com/peggyjv/sommelier/v7/x/cork/types/v2"
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

	ctx.EventManager().EmitEvents(
		sdk.Events{
			sdk.NewEvent(
				sdk.EventTypeMessage,
				sdk.NewAttribute(sdk.AttributeKeyModule, corktypes.AttributeValueCategory),
			),
			sdk.NewEvent(
				corktypes.EventTypeCork,
				sdk.NewAttribute(corktypes.AttributeKeySigner, signer.String()),
				sdk.NewAttribute(corktypes.AttributeKeyValidator, validatorAddr.String()),
				sdk.NewAttribute(corktypes.AttributeKeyCork, msg.Cork.String()),
				sdk.NewAttribute(corktypes.AttributeKeyBlockHeight, fmt.Sprintf("%d", msg.BlockHeight)),
			),
		},
	)

	return &types.MsgScheduleCorkResponse{Id: hex.EncodeToString(corkID)}, nil
}
