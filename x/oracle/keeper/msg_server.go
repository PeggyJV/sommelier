package keeper

import (
	"bytes"
	"context"
	"fmt"
	"sort"

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

	k.SetOracleDataPrevote(ctx, valaddr, msg)

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
	prevote := k.GetOracleDataPrevote(ctx, valaddr)

	// check that there is a prevote
	if prevote == nil || len(prevote.Hashes) == 0 {
		return nil, sdkerrors.Wrap(types.ErrNoPrevote, valaddr.String())
	}

	// ensure that the right number of data is in the msg
	if len(prevote.Hashes) != len(msg.OracleData) {
		return nil, sdkerrors.Wrap(
			types.ErrWrongNumber,
			fmt.Sprintf("oracle data exp(%d) got(%d)", len(prevote.Hashes), len(msg.OracleData)),
		)
	}
	fmt.Println("Right number")

	// ensure that the right number of salts is in the msg
	if len(prevote.Hashes) != len(msg.Salt) {
		return nil, sdkerrors.Wrap(
			types.ErrWrongNumber,
			fmt.Sprintf("salt exp(%d) got(%d)", len(prevote.Hashes), len(msg.Salt)),
		)
	}

	fmt.Println("Has salt, right number")
	// validate the hashes
	got := []string{}
	for i := range msg.OracleData {
		salt := msg.Salt[i]
		od, err := types.UnpackOracleData(msg.OracleData[i])
		if err != nil {
			return nil, sdkerrors.Wrap(types.ErrUnpackOracleData, fmt.Sprintf("index %d", i))
		}
		voteHash := types.DataHash(salt, od.CannonicalJSON(), valaddr)
		if !bytes.Equal(voteHash, prevote.Hashes[i]) {
			return nil, sdkerrors.Wrap(
				types.ErrHashMismatch,
				fmt.Sprintf("precommit(%x) commit(%x)", prevote.Hashes[i], voteHash),
			)
		}
		got = append(got, od.Type())
	}
	fmt.Println("Hashes correct")

	// ensure that the right number of data types have been submitted
	exp := k.GetParamSet(ctx).DataTypes
	if len(exp) != len(got) {
		return nil, sdkerrors.Wrap(
			types.ErrWrongNumber,
			fmt.Sprintf("oracle data types exp(%d) got(%d)", len(exp), len(got)),
		)
	}
	fmt.Println("datatypes correct")

	// ensure that all of the right data types have been submitted
	sort.Strings(exp)
	sort.Strings(got)
	for i := range exp {
		if exp[i] != got[i] {
			return nil, sdkerrors.Wrap(
				types.ErrWrongDataType,
				fmt.Sprintf("exp(%s) got(%s)", exp[i], got[i]),
			)
		}
	}
	fmt.Println("datatypes submitted")

	// set the vote in the store
	k.SetOracleDataVote(ctx, valaddr, msg)
	fmt.Println("vote set")

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
