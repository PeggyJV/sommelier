package oracle

import (
	"bytes"
	"fmt"

	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"

	"github.com/peggyjv/sommelier/x/oracle/keeper"
	"github.com/peggyjv/sommelier/x/oracle/types"
)

// NewHandler returns a handler for "oracle" type messages.
func NewHandler(k keeper.Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) (*sdk.Result, error) {
		ctx = ctx.WithEventManager(sdk.NewEventManager())

		switch msg := msg.(type) {
		case *types.MsgDelegateFeedConsent:
			return handleMsgDelegateFeedConsent(ctx, k, *msg)
		case *types.MsgAggregateExchangeRatePrevote:
			return handleMsgAggregateExchangeRatePrevote(ctx, k, *msg)
		case *types.MsgAggregateExchangeRateVote:
			return handleMsgAggregateExchangeRateVote(ctx, k, *msg)
		default:
			return nil, sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "unrecognized distribution message type: %T", msg)
		}
	}
}

// handleMsgDelegateFeedConsent handles a MsgDelegateFeedConsent
func handleMsgDelegateFeedConsent(ctx sdk.Context, keeper keeper.Keeper, msg types.MsgDelegateFeedConsent) (*sdk.Result, error) {
	signer, err := sdk.ValAddressFromBech32(msg.Operator)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, err.Error())
	}

	// Check the delegator is a validator
	val := keeper.StakingKeeper.Validator(ctx, signer)
	if val == nil {
		return nil, sdkerrors.Wrap(stakingtypes.ErrNoValidatorFound, signer.String())
	}

	delegate, err := sdk.AccAddressFromBech32(msg.Delegate)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, err.Error())
	}
	// Set the delegation
	keeper.SetOracleDelegate(ctx, signer, delegate)

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeFeedDelegate,
			sdk.NewAttribute(types.AttributeKeyOperator, msg.Operator),
			sdk.NewAttribute(types.AttributeKeyFeeder, msg.Delegate),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
		),
	})

	return &sdk.Result{Events: ctx.EventManager().Events().ToABCIEvents()}, nil
}

// handleMsgAggregateExchangeRatePrevote handles a MsgAggregateExchangeRatePrevote
func handleMsgAggregateExchangeRatePrevote(ctx sdk.Context, keeper keeper.Keeper, msg types.MsgAggregateExchangeRatePrevote) (*sdk.Result, error) {
	feeder, err := sdk.AccAddressFromBech32(msg.Feeder)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, err.Error())
	}
	validator, err := sdk.ValAddressFromBech32(msg.Validator)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, err.Error())
	}
	if !feeder.Equals(validator) {
		delegate := keeper.GetOracleDelegate(ctx, validator)
		if !delegate.Equals(feeder) {
			return nil, sdkerrors.Wrap(types.ErrNoVotingPermission, msg.Feeder)
		}
	}

	// Check that the given validator exists
	val := keeper.StakingKeeper.Validator(ctx, validator)
	if val == nil {
		return nil, sdkerrors.Wrap(stakingtypes.ErrNoValidatorFound, msg.Validator)
	}

	aggregatePrevote := types.NewAggregateExchangeRatePrevote(msg.Hash, validator, ctx.BlockHeight())
	keeper.AddAggregateExchangeRatePrevote(ctx, aggregatePrevote)

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeAggregatePrevote,
			sdk.NewAttribute(types.AttributeKeyVoter, msg.Validator),
			sdk.NewAttribute(types.AttributeKeyFeeder, msg.Feeder),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
		),
	})

	return &sdk.Result{Events: ctx.EventManager().Events().ToABCIEvents()}, nil
}

// handleMsgAggregateExchangeRateVote handles a MsgAggregateExchangeRateVote
func handleMsgAggregateExchangeRateVote(ctx sdk.Context, keeper keeper.Keeper, msg types.MsgAggregateExchangeRateVote) (*sdk.Result, error) {
	feeder, err := sdk.AccAddressFromBech32(msg.Feeder)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, err.Error())
	}
	validator, err := sdk.ValAddressFromBech32(msg.Validator)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, err.Error())
	}
	if !feeder.Equals(validator) {
		delegate := keeper.GetOracleDelegate(ctx, validator)
		if !delegate.Equals(feeder) {
			return nil, sdkerrors.Wrap(types.ErrNoVotingPermission, msg.Feeder)
		}
	}

	// Check that the given validator exists
	val := keeper.StakingKeeper.Validator(ctx, validator)
	if val == nil {
		return nil, sdkerrors.Wrap(stakingtypes.ErrNoValidatorFound, msg.Validator)
	}

	params := keeper.GetParams(ctx)

	aggregatePrevote, err := keeper.GetAggregateExchangeRatePrevote(ctx, validator)
	if err != nil {
		return nil, sdkerrors.Wrap(types.ErrNoAggregatePrevote, msg.Validator)
	}

	// Check a msg is submitted porper period
	if (ctx.BlockHeight()/params.VotePeriod)-(aggregatePrevote.SubmitBlock/params.VotePeriod) != 1 {
		return nil, types.ErrRevealPeriodMissMatch
	}

	exchangeRateTuples, err := types.ParseExchangeRateTuples(msg.ExchangeRates)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidCoins, err.Error())
	}

	// check all denoms are in the vote target
	for _, tuple := range exchangeRateTuples {
		if !keeper.IsVoteTarget(ctx, tuple.Denom) {
			return nil, sdkerrors.Wrap(types.ErrUnknowDenom, tuple.Denom)
		}
	}

	voter, err := sdk.ValAddressFromBech32(aggregatePrevote.Voter)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, err.Error())
	}

	// Verify a exchange rate with aggregate prevote hash
	hash := types.GetAggregateVoteHash(msg.Salt, msg.ExchangeRates, voter)
	if !bytes.Equal(aggregatePrevote.Hash, hash.Bytes()) {
		return nil, sdkerrors.Wrap(types.ErrVerificationFailed, fmt.Sprintf("must be given %s not %s", aggregatePrevote.Hash, hash))
	}

	// Move aggregate prevote to aggregate vote with given exchange rates
	keeper.AddAggregateExchangeRateVote(ctx, types.NewAggregateExchangeRateVote(exchangeRateTuples, voter))
	keeper.DeleteAggregateExchangeRatePrevote(ctx, aggregatePrevote)

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeAggregateVote,
			sdk.NewAttribute(types.AttributeKeyVoter, msg.Validator),
			sdk.NewAttribute(types.AttributeKeyExchangeRates, msg.ExchangeRates),
			sdk.NewAttribute(types.AttributeKeyFeeder, msg.Feeder),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
		),
	})

	return &sdk.Result{Events: ctx.EventManager().Events().ToABCIEvents()}, nil
}
