package keeper

import (
	"bytes"
	"context"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/peggyjv/sommelier/x/oracle/types"
)

var _ types.MsgServer = Keeper{}

// DelegateFeedConsent implements types.MsgServer
func (k Keeper) DelegateFeedConsent(c context.Context, msg *types.MsgDelegateFeedConsent) (*types.MsgDelegateFeedConsentResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)

	val, del := msg.MustGetValidator(), msg.MustGetDelegate()

	if k.stakingKeeper.Validator(ctx, sdk.ValAddress(val)) == nil {
		return nil, sdkerrors.Wrap(stakingtypes.ErrNoValidatorFound, val.String())
	}

	k.SetValidatorDelegateAddress(ctx, val, del)

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeyAction, types.EventTypeDelegateFeed),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Validator),
			sdk.NewAttribute(types.AttributeKeyValidator, msg.Validator),
			sdk.NewAttribute(types.AttributeKeyDeleagate, msg.Delegate),
		),
	)

	return &types.MsgDelegateFeedConsentResponse{}, nil
}

// OracleDataPrevote implements types.MsgServer
func (k Keeper) OracleDataPrevote(c context.Context, msg *types.MsgOracleDataPrevote) (*types.MsgOracleDataPrevoteResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)

	signer := msg.MustGetSigner()
	validatorAddr := k.GetValidatorAddressFromDelegate(ctx, signer)
	if validatorAddr == nil {
		validator := k.stakingKeeper.Validator(ctx, sdk.ValAddress(signer))
		if validator == nil {
			return nil, sdkerrors.Wrap(types.ErrUnknown, "validator")
		}

		validatorAddr = sdk.AccAddress(validator.GetOperator())
	}

	k.SetOracleDataPrevote(ctx, validatorAddr, msg)

	ctx.EventManager().EmitEvents(
		sdk.Events{
			sdk.NewEvent(
				sdk.EventTypeMessage,
				sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			),
			sdk.NewEvent(
				types.EventTypeOracleDataPrevote,
				sdk.NewAttribute(types.AttributeKeySigner, signer.String()),
				sdk.NewAttribute(types.AttributeKeyValidator, validatorAddr.String()),
				sdk.NewAttribute(types.AttributeKeyHashes, fmt.Sprintf("%x", bytes.Join(msg.Hashes, []byte(",")))),
			),
		},
	)

	return &types.MsgOracleDataPrevoteResponse{}, nil
}

// OracleDataVote implements types.MsgServer
func (k Keeper) OracleDataVote(c context.Context, msg *types.MsgOracleDataVote) (*types.MsgOracleDataVoteResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)

	// Make sure that the message was properly signed
	signer := msg.MustGetSigner()
	validatorAddr := k.GetValidatorAddressFromDelegate(ctx, signer)
	if validatorAddr == nil {
		validator := k.stakingKeeper.Validator(ctx, sdk.ValAddress(signer))
		if validator == nil {
			return nil, sdkerrors.Wrap(types.ErrUnknown, "validator")
		}

		validatorAddr = sdk.AccAddress(validator.GetOperator())
	}

	// Get the prevote for that validator from the store
	prevote := k.GetOracleDataPrevote(ctx, validatorAddr)

	// check that there is a prevote
	if prevote == nil || len(prevote.Hashes) == 0 {
		return nil, sdkerrors.Wrap(types.ErrNoPrevote, validatorAddr.String())
	}

	// ensure that the right number of data is in the msg
	if len(prevote.Hashes) != len(msg.OracleData) {
		return nil, sdkerrors.Wrap(
			types.ErrWrongNumber,
			fmt.Sprintf("oracle data exp(%d) got(%d)", len(prevote.Hashes), len(msg.OracleData)),
		)
	}

	// ensure that the right number of salts is in the msg
	if len(prevote.Hashes) != len(msg.Salt) {
		return nil, sdkerrors.Wrap(
			types.ErrWrongNumber,
			fmt.Sprintf("salt exp(%d) got(%d)", len(prevote.Hashes), len(msg.Salt)),
		)
	}

	allowedTypesMap := make(map[string]bool)
	allowedDataTypes := k.GetParamSet(ctx).DataTypes

	for _, allowedDataType := range allowedDataTypes {
		allowedTypesMap[allowedDataType] = true
	}

	oracleEvents := sdk.Events{
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeyAction, types.EventTypeOracleDataVote),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Signer),
		),
	}

	for i, oracleDataAny := range msg.OracleData {
		salt := msg.Salt[i]

		// unpack the oracle data one by one
		oracleData, err := types.UnpackOracleData(oracleDataAny)
		if err != nil {
			return nil, sdkerrors.Wrap(types.ErrUnpackOracleData, fmt.Sprintf("index %d", i))
		}

		if !allowedTypesMap[oracleData.Type()] {
			return nil, sdkerrors.Wrap(
				types.ErrUnsupportedDataType,
				fmt.Sprintf("%s, allowed %v", oracleData.Type(), allowedDataTypes),
			)
		}

		// parse data to json in order to compute the vote hash
		jsonBz, err := oracleData.MarshalJSON()
		if err != nil {
			return nil, sdkerrors.Wrapf(
				sdkerrors.ErrJSONMarshal,
				"failed to marshal json for oracle data with id: %s", oracleData.GetID(),
			)
		}

		// calculate the vote hash on the server
		voteHash := types.DataHash(salt, string(jsonBz), validatorAddr)

		// compare to prevote hash
		if !bytes.Equal(voteHash, prevote.Hashes[i]) {
			return nil, sdkerrors.Wrap(
				types.ErrHashMismatch,
				fmt.Sprintf("precommit(%x) commit(%x)", prevote.Hashes[i], voteHash),
			)
		}

		oracleEvents = append(
			oracleEvents,
			sdk.NewEvent(
				types.EventTypeOracleDataVote,
				sdk.NewAttribute(types.AttributeKeyOracleDataID, oracleData.GetID()),
				sdk.NewAttribute(types.AttributeKeyOracleDataType, oracleData.Type()),
				sdk.NewAttribute(types.AttributeKeyValidator, validatorAddr.String()),
			),
		)
	}

	// set the vote in the store
	k.SetOracleDataVote(ctx, validatorAddr, msg)
	ctx.EventManager().EmitEvents(oracleEvents)

	return &types.MsgOracleDataVoteResponse{}, nil
}
