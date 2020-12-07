package oracle

import (
	"fmt"

	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/staking"

	"github.com/peggyjv/sommelier/x/oracle//types"
)

// NewHandler returns a handler for "oracle" type messages.
func NewHandler(k Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) (*sdk.Result, error) {
		ctx = ctx.WithEventManager(sdk.NewEventManager())

		switch msg := msg.(type) {
		case MsgDelegateFeedConsent:
			return handleMsgDelegateFeedConsent(ctx, k, msg)
		case MsgAggregateExchangeRatePrevote:
			return handleMsgAggregateExchangeRatePrevote(ctx, k, msg)
		case MsgAggregateExchangeRateVote:
			return handleMsgAggregateExchangeRateVote(ctx, k, msg)
		default:
			return nil, sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "unrecognized distribution message type: %T", msg)
		}
	}
}

// handleMsgDelegateFeedConsent handles a MsgDelegateFeedConsent
func handleMsgDelegateFeedConsent(ctx sdk.Context, keeper Keeper, msg MsgDelegateFeedConsent) (*sdk.Result, error) {
	signer := msg.Operator

	// Check the delegator is a validator
	val := keeper.StakingKeeper.Validator(ctx, signer)
	if val == nil {
		return nil, sdkerrors.Wrap(staking.ErrNoValidatorFound, signer.String())
	}

	// Set the delegation
	keeper.SetOracleDelegate(ctx, signer, msg.Delegate)

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeFeedDelegate,
			sdk.NewAttribute(types.AttributeKeyOperator, msg.Operator.String()),
			sdk.NewAttribute(types.AttributeKeyFeeder, msg.Delegate.String()),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
		),
	})

	return &sdk.Result{Events: ctx.EventManager().Events()}, nil
}

// handleMsgAggregateExchangeRatePrevote handles a MsgAggregateExchangeRatePrevote
func handleMsgAggregateExchangeRatePrevote(ctx sdk.Context, keeper Keeper, msg MsgAggregateExchangeRatePrevote) (*sdk.Result, error) {
	if !msg.Feeder.Equals(msg.Validator) {
		delegate := keeper.GetOracleDelegate(ctx, msg.Validator)
		if !delegate.Equals(msg.Feeder) {
			return nil, sdkerrors.Wrap(ErrNoVotingPermission, msg.Feeder.String())
		}
	}

	// Check that the given validator exists
	val := keeper.StakingKeeper.Validator(ctx, msg.Validator)
	if val == nil {
		return nil, sdkerrors.Wrap(staking.ErrNoValidatorFound, msg.Validator.String())
	}

	aggregatePrevote := NewAggregateExchangeRatePrevote(msg.Hash, msg.Validator, ctx.BlockHeight())
	keeper.AddAggregateExchangeRatePrevote(ctx, aggregatePrevote)

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeAggregatePrevote,
			sdk.NewAttribute(types.AttributeKeyVoter, msg.Validator.String()),
			sdk.NewAttribute(types.AttributeKeyFeeder, msg.Feeder.String()),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
		),
	})

	return &sdk.Result{Events: ctx.EventManager().Events()}, nil
}

// handleMsgAggregateExchangeRateVote handles a MsgAggregateExchangeRateVote
func handleMsgAggregateExchangeRateVote(ctx sdk.Context, keeper Keeper, msg MsgAggregateExchangeRateVote) (*sdk.Result, error) {
	if !msg.Feeder.Equals(msg.Validator) {
		delegate := keeper.GetOracleDelegate(ctx, msg.Validator)
		if !delegate.Equals(msg.Feeder) {
			return nil, sdkerrors.Wrap(ErrNoVotingPermission, msg.Feeder.String())
		}
	}

	// Check that the given validator exists
	val := keeper.StakingKeeper.Validator(ctx, msg.Validator)
	if val == nil {
		return nil, sdkerrors.Wrap(staking.ErrNoValidatorFound, msg.Validator.String())
	}

	params := keeper.GetParams(ctx)

	aggregatePrevote, err := keeper.GetAggregateExchangeRatePrevote(ctx, msg.Validator)
	if err != nil {
		return nil, sdkerrors.Wrap(ErrNoAggregatePrevote, msg.Validator.String())
	}

	// Check a msg is submitted porper period
	if (ctx.BlockHeight()/params.VotePeriod)-(aggregatePrevote.SubmitBlock/params.VotePeriod) != 1 {
		return nil, ErrRevealPeriodMissMatch
	}

	exchangeRateTuples, err := types.ParseExchangeRateTuples(msg.ExchangeRates)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidCoins, err.Error())
	}

	// check all denoms are in the vote target
	for _, tuple := range exchangeRateTuples {
		if !keeper.IsVoteTarget(ctx, tuple.Denom) {
			return nil, sdkerrors.Wrap(ErrUnknowDenom, tuple.Denom)
		}
	}

	// Verify a exchange rate with aggregate prevote hash
	hash := GetAggregateVoteHash(msg.Salt, msg.ExchangeRates, aggregatePrevote.Voter)
	if !aggregatePrevote.Hash.Equal(hash) {
		return nil, sdkerrors.Wrap(ErrVerificationFailed, fmt.Sprintf("must be given %s not %s", aggregatePrevote.Hash, hash))
	}

	// Move aggregate prevote to aggregate vote with given exchange rates
	keeper.AddAggregateExchangeRateVote(ctx, NewAggregateExchangeRateVote(exchangeRateTuples, aggregatePrevote.Voter))
	keeper.DeleteAggregateExchangeRatePrevote(ctx, aggregatePrevote)

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeAggregateVote,
			sdk.NewAttribute(types.AttributeKeyVoter, msg.Validator.String()),
			sdk.NewAttribute(types.AttributeKeyExchangeRates, msg.ExchangeRates),
			sdk.NewAttribute(types.AttributeKeyFeeder, msg.Feeder.String()),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
		),
	})

	return &sdk.Result{Events: ctx.EventManager().Events()}, nil
}
