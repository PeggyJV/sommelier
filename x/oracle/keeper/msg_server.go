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

type msgServer struct {
	Keeper
}

// NewMsgServerImpl returns an implementation of the oracle MsgServer interface
// for the provided Keeper.
func NewMsgServerImpl(keeper Keeper) types.MsgServer {
	return &msgServer{Keeper: keeper}
}

var _ types.MsgServer = msgServer{}

// DelegateFeedConsent implements types.MsgServer
func (k msgServer) DelegateFeedConsent(c context.Context, msg *types.MsgDelegateFeedConsent) (*types.MsgDelegateFeedConsentResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)

	val, del := msg.MustGetValidator(), msg.MustGetDelegate()

	if k.Keeper.stakingKeeper.Validator(ctx, sdk.validatorAddress(val)) == nil {
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
			sdk.NewAttribute(types.AttributeKeyDeleagte, msg.Delegate),
		),
	)

	return &types.MsgDelegateFeedConsentResponse{}, nil
}

// OracleDataPrevote implements types.MsgServer
func (k msgServer) OracleDataPrevote(c context.Context, msg *types.MsgOracleDataPrevote) (*types.MsgOracleDataPrevoteResponse, error) {
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

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
		),
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
func (k msgServer) OracleDataVote(c context.Context, msg *types.MsgOracleDataVote) (*types.MsgOracleDataVoteResponse, error) {
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
