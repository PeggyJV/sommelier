package keeper

import (
	"bytes"
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/peggyjv/sommelier/x/oracle/types"
)

type msgServer struct {
	Keeper
}

// NewMsgServerImpl returns an implementation of the distribution MsgServer interface
// for the provided Keeper.
func NewMsgServerImpl(keeper Keeper) types.MsgServer {
	return &msgServer{Keeper: keeper}
}

// MsgServer is the server API for Msg service.

func (m msgServer) DelegateFeedConsent(c context.Context, msg *types.MsgDelegateFeedConsent) (*types.MsgDelegateFeedConsentResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	signer, err := sdk.ValAddressFromBech32(msg.Operator)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, err.Error())
	}

	// Check the delegator is a validator
	val := m.StakingKeeper.Validator(ctx, signer)
	if val == nil {
		return nil, sdkerrors.Wrap(stakingtypes.ErrNoValidatorFound, signer.String())
	}

	delegate, err := sdk.AccAddressFromBech32(msg.Delegate)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, err.Error())
	}

	// Set the delegation
	m.SetOracleDelegate(ctx, signer, delegate)

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

	return &types.MsgDelegateFeedConsentResponse{}, nil
}

func (m msgServer) AggregateExchangeRatePrevote(c context.Context, msg *types.MsgAggregateExchangeRatePrevote) (*types.MsgAggregateExchangeRatePrevoteResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	feeder, err := sdk.AccAddressFromBech32(msg.Feeder)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, err.Error())
	}
	validator, err := sdk.ValAddressFromBech32(msg.Validator)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, err.Error())
	}
	if !feeder.Equals(validator) {
		delegate := m.GetOracleDelegate(ctx, validator)
		if !delegate.Equals(feeder) {
			return nil, sdkerrors.Wrap(types.ErrNoVotingPermission, msg.Feeder)
		}
	}

	// Check that the given validator exists
	val := m.StakingKeeper.Validator(ctx, validator)
	if val == nil {
		return nil, sdkerrors.Wrap(stakingtypes.ErrNoValidatorFound, msg.Validator)
	}

	aggregatePrevote := types.NewAggregateExchangeRatePrevote(msg.Hash, validator, ctx.BlockHeight())
	m.AddAggregateExchangeRatePrevote(ctx, aggregatePrevote)

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
	return &types.MsgAggregateExchangeRatePrevoteResponse{}, nil
}

func (m msgServer) AggregateExchangeRateVote(c context.Context, msg *types.MsgAggregateExchangeRateVote) (*types.MsgAggregateExchangeRateVoteResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	feeder, err := sdk.AccAddressFromBech32(msg.Feeder)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, err.Error())
	}
	validator, err := sdk.ValAddressFromBech32(msg.Validator)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, err.Error())
	}
	if !feeder.Equals(validator) {
		delegate := m.GetOracleDelegate(ctx, validator)
		if !delegate.Equals(feeder) {
			return nil, sdkerrors.Wrap(types.ErrNoVotingPermission, msg.Feeder)
		}
	}

	// Check that the given validator exists
	val := m.StakingKeeper.Validator(ctx, validator)
	if val == nil {
		return nil, sdkerrors.Wrap(stakingtypes.ErrNoValidatorFound, msg.Validator)
	}

	params := m.GetParams(ctx)

	aggregatePrevote, err := m.GetAggregateExchangeRatePrevote(ctx, validator)
	if err != nil {
		return nil, sdkerrors.Wrap(types.ErrNoAggregatePrevote, msg.Validator)
	}

	// Check a msg is submitted proper period
	if ctx.BlockHeight()-aggregatePrevote.SubmitBlock/params.VotePeriod != 1 {
		return nil, types.ErrRevealPeriodMissMatch
	}

	// NOTE: error checked on msg validation
	exchangeRateTuples, _ := sdk.ParseDecCoins(msg.ExchangeRates)

	// check all denoms are in the vote target
	for _, tuple := range exchangeRateTuples {
		if !m.IsVoteTarget(ctx, tuple.Denom) {
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
		return nil, sdkerrors.Wrapf(types.ErrVerificationFailed, "must be given %s not %s", aggregatePrevote.Hash, hash)
	}

	// Move aggregate prevote to aggregate vote with given exchange rates
	m.AddAggregateExchangeRateVote(ctx, types.NewAggregateExchangeRateVote(exchangeRateTuples, voter))
	m.DeleteAggregateExchangeRatePrevote(ctx, aggregatePrevote)

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

	return &types.MsgAggregateExchangeRateVoteResponse{}, nil
}
