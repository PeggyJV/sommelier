package keeper

import (
	"context"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/ethereum/go-ethereum/common"

	"github.com/peggyjv/sommelier/v3/x/cork/types"
)

var _ types.MsgServer = Keeper{}

func (k Keeper) signerToValAddr(ctx sdk.Context, signer sdk.AccAddress) (sdk.ValAddress, error) {
	validatorAddr := k.gravityKeeper.GetOrchestratorValidatorAddress(ctx, signer)
	if validatorAddr == nil {
		validator := k.stakingKeeper.Validator(ctx, sdk.ValAddress(signer))
		if validator == nil {
			return nil, sdkerrors.Wrap(stakingtypes.ErrNoValidatorFound, sdk.ValAddress(signer).String())
		}

		validatorAddr = validator.GetOperator()
		// NOTE: we set the validator address so we don't have to call look up for the validator
		// everytime a validator feeder submits oracle data
		k.gravityKeeper.SetOrchestratorValidatorAddress(ctx, validatorAddr, signer)
	}
	return validatorAddr, nil
}

// SubmitCork implements types.MsgServer
func (k Keeper) SubmitCork(c context.Context, msg *types.MsgSubmitCorkRequest) (*types.MsgSubmitCorkResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)

	if !k.HasCellarID(ctx, common.HexToAddress(msg.Cork.TargetContractAddress)) {
		return nil, types.ErrUnmanagedCellarAddress
	}

	signer := msg.MustGetSigner()
	validatorAddr, err := k.signerToValAddr(ctx, signer)
	if err != nil {
		return nil, err
	}

	k.SetCork(ctx, validatorAddr, *msg.Cork)

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
			),
		},
	)

	return &types.MsgSubmitCorkResponse{}, nil
}
