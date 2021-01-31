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

	if err := msg.ValidateBasic(); err != nil {
		return nil, sdkerrors.Wrap(err, "validate basic failed")
	}

	val, del := msg.MustGetValidator(), msg.MustGetDelegate()

	if k.Keeper.StakingKeeper.Validator(ctx, sdk.ValAddress(val)) == nil {
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

	if err := msg.ValidateBasic(); err != nil {
		return nil, sdkerrors.Wrap(err, "validate basic failed")
	}

	signer := msg.MustGetSigner()
	valaddr := k.GetValidatorAddressFromDelegate(ctx, signer)
	if valaddr == nil {
		sval := k.Keeper.StakingKeeper.Validator(ctx, sdk.ValAddress(signer))
		if sval == nil {
			return nil, sdkerrors.Wrap(types.ErrUnknown, "validator")
		}
		valaddr = sdk.AccAddress(sval.GetOperator())
	}

	k.SetOracleDataPrevote(ctx, valaddr, msg.Hashes)

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeyAction, types.EventTypeOracleDataPrevote),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Signer),
			sdk.NewAttribute(types.AttributeKeyValidator, valaddr.String()),
			sdk.NewAttribute(types.AttributeKeyHashes, fmt.Sprintf("%x", bytes.Join(msg.Hashes, []byte(",")))),
		),
	)

	return &types.MsgOracleDataPrevoteResponse{}, nil
}

// OracleDataVote implements types.MsgServer
func (k msgServer) OracleDataVote(c context.Context, msg *types.MsgOracleDataVote) (*types.MsgOracleDataVoteResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	if err := msg.ValidateBasic(); err != nil {
		return nil, sdkerrors.Wrap(err, "validate basic failed")
	}

	// Make sure that the message was properly signed
	signer := msg.MustGetSigner()
	valaddr := k.GetValidatorAddressFromDelegate(ctx, signer)
	if valaddr == nil {
		sval := k.Keeper.StakingKeeper.Validator(ctx, sdk.ValAddress(signer))
		if sval == nil {
			return nil, sdkerrors.Wrap(types.ErrUnknown, "validator")
		}
		valaddr = sdk.AccAddress(sval.GetOperator())
	}

	// Get the prevote for that validator from the store
	hashes := k.GetOracleDataPrevote(ctx, valaddr)

	// check that there is a prevote
	if hashes == nil || len(hashes) == 0 {
		return nil, sdkerrors.Wrap(types.ErrNoPrevote, valaddr.String())
	}

	// ensure that the right number of data is in the msg
	if len(hashes) != len(msg.OracleData) {
		return nil, sdkerrors.Wrap(
			types.ErrWrongNumber,
			fmt.Sprintf("oracle data exp(%d) got(%d)", len(hashes), len(msg.OracleData)),
		)
	}

	// ensure that the right number of salts is in the msg
	if len(hashes) != len(msg.Salt) {
		return nil, sdkerrors.Wrap(
			types.ErrWrongNumber,
			fmt.Sprintf("salt exp(%d) got(%d)", len(hashes), len(msg.Salt)),
		)
	}

	// validate the hashes
	for i := range msg.OracleData {
		salt := msg.Salt[i]
		od, err := types.UnpackOracleData(msg.OracleData[i])
		if err != nil {
			return nil, sdkerrors.Wrap(types.ErrUnpackOracleData, fmt.Sprintf("index %d", i))
		}
		voteHash := types.DataHash(salt, od.CannonicalJSON(), valaddr)
		if !bytes.Equal(voteHash, hashes[i]) {
			return nil, sdkerrors.Wrap(
				types.ErrHashMismatch,
				fmt.Sprintf("precommit(%x) commit(%x)", hashes[i], voteHash),
			)
		}
	}

	// set the vote in the store
	k.SetOracleDataVote(ctx, valaddr, msg)

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeyAction, types.EventTypeOracleDataVote),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Signer),
			sdk.NewAttribute(types.AttributeKeyValidator, valaddr.String()),
			// TODO: emit other data here?
		),
	)

	return &types.MsgOracleDataVoteResponse{}, nil
}
