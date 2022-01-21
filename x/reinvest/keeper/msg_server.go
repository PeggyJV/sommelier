package keeper

import (
	"context"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"

	"github.com/peggyjv/sommelier/x/reinvest/types"
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

// SubmitReinvest implements types.MsgServer
func (k Keeper) SubmitReinvest(c context.Context, msg *types.MsgSubmitReinvestRequest) (*types.MsgSubmitReinvestResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)

	signer := msg.MustGetSigner()
	validatorAddr, err := k.signerToValAddr(ctx, signer)
	if err != nil {
		return nil, err
	}

	k.SetReinvestment(ctx, validatorAddr, *msg.Reinvestment)

	ctx.EventManager().EmitEvents(
		sdk.Events{
			sdk.NewEvent(
				sdk.EventTypeMessage,
				sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			),
			sdk.NewEvent(
				types.EventTypeReinvest,
				sdk.NewAttribute(types.AttributeKeySigner, signer.String()),
				sdk.NewAttribute(types.AttributeKeyValidator, validatorAddr.String()),
				sdk.NewAttribute(types.AttributeKeyReinvestment, msg.Reinvestment.String()),
			),
		},
	)

	return &types.MsgSubmitReinvestResponse{}, nil
}
