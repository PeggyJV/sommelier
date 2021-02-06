package keeper

import (
	"bytes"
	"context"
	"fmt"

	"github.com/cosmos/cosmos-sdk/telemetry"
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

	if k.Keeper.stakingKeeper.Validator(ctx, sdk.validatorAddress(val)) == nil {
		return nil, sdkerrors.Wrap(stakingtypes.ErrNoValidatorFound, val.String())
	}

	if validator.IsUnbonded() {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrUnauthorized, "validator %s cannot be unbonded", validator.GetOperator())
	}

	if validator.IsJailed() {
		return nil, sdkerrors.Wrap(stakingtypes.ErrValidatorJailed, validator.GetOperator().String())
	}

	// check that the delegate feeder is not a validator, this prevents mirroring and freeloading
	// See https://medium.com/fabric-ventures/decentralised-oracles-a-comprehensive-overview-d3168b9a8841
	if k.stakingKeeper.Validator(ctx, sdk.ValAddress(del)) != nil {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrUnauthorized, "feeder delegate %s cannot be a validator", del)
	}

	k.SetValidatorDelegateAddress(ctx, del, sdk.ValAddress(val))

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
		sval := k.Keeper.stakingKeeper.Validator(ctx, sdk.validatorAddress(signer))
		if sval == nil {
			return nil, sdkerrors.Wrap(types.ErrUnknown, "validator")
		}

		validatorAddr = sdk.AccAddress(sval.GetOperator())
	}

	k.SetOracleDataPrevote(ctx, validatorAddr, msg)

	// NOTE: a validator can prevote multiple times but can only submit a single vote

	// TODO: set prevote for current voting period
	k.SetOracleDataPrevote(ctx, validatorAddr, *msg.Prevote)

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
				sdk.NewAttribute(types.AttributeKeyPrevoteHash, msg.Prevote.Hash.String()),
			),
		},
	)

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeOracleDataPrevote,
			sdk.NewAttribute(types.AttributeKeySigner, signer.String()),
			sdk.NewAttribute(types.AttributeKeyValidator, validatorAddr.String()),
			sdk.NewAttribute(types.AttributeKeyHashes, fmt.Sprintf("%x", bytes.Join(msg.Hashes, []byte(",")))),
		),
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
		validator := k.Keeper.stakingKeeper.Validator(ctx, sdk.ValAddress(signer))
		if validator == nil {
			return nil, sdkerrors.Wrap(types.ErrUnknown, "validator")
		}

		validatorAddr = sdk.AccAddress(validator.GetOperator())
	}

	// Get the prevote for that validator from the store
	prevote := k.GetOracleDataPrevote(ctx, validatorAddr)

	// Get the prevote for that validator from the store
	prevote, found := k.GetOracleDataPrevote(ctx, validatorAddr)
	// check that there is a prevote
	if prevote == nil || len(prevote.Hashes) == 0 {
		return nil, sdkerrors.Wrap(types.ErrNoPrevote, validatorAddr.String())
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

	// parse data to json in order to compute the vote hash and sort
	jsonBz, err := json.Marshal(msg.Vote.Feed.Data)
	if err != nil {
		return nil, sdkerrors.Wrap(
			sdkerrors.ErrJSONMarshal, "failed to marshal json oracle data feed",
		)
	}

	allowedTypesMap := make(map[string]bool)
	allowedDataTypes := k.GetParamSet(ctx).DataTypes

	for _, allowedDataType := range allowedDataTypes {
		allowedTypesMap[allowedDataType] = true
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

		// ensure that the data parses
		uniswapData, ok := oracleData.(*types.UniswapData)
		if !ok {
			return nil, sdkerrors.Wrap(types.ErrInvalidOracleData, "only uniswap data currently supported")
		}

		// calculate the vote hash on the server
		voteHash := types.DataHash(salt, uniswapData.CannonicalJSON(), validatorAddr)

		// compare to prevote hash
		if !bytes.Equal(voteHash, prevote.Hashes[i]) {
			return nil, sdkerrors.Wrap(
				types.ErrHashMismatch,
				fmt.Sprintf("precommit(%x) commit(%x)", prevote.Hashes[i], voteHash),
			)
		}
	}

	// set the vote in the store
	k.SetOracleDataVote(ctx, validatorAddr, msg)

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeyAction, types.EventTypeOracleDataVote),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Signer),
			sdk.NewAttribute(types.AttributeKeyValidator, validatorAddr.String()),
			// TODO: emit other data here?
		),
	)

	return &types.MsgOracleDataVoteResponse{}, nil
}
