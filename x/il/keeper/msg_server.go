package keeper

import (
	"context"

	"github.com/armon/go-metrics"
	"github.com/cosmos/cosmos-sdk/telemetry"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/peggyjv/sommelier/x/il/types"
)

var _ types.MsgServer = Keeper{}

// MsgServer is the server API for Msg service.

// CreateStoploss for a given uniswap pair
func (k Keeper) CreateStoploss(c context.Context, msg *types.MsgStoploss) (*types.MsgStoplossResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)

	// NOTE: error checked during msg validation
	address, _ := sdk.AccAddressFromBech32(msg.Address)

	// TODO: check if uniswap pair exists on the oracle
	// oracleData, found := k.oracleKeeper.GetLatestAggregatedOracleData(ctx, oracletypes.UniswapDataType ,msg.Stoploss.UniswapPairId)

	// TODO: ensure that pair ratio is w/in band of pair ratio from oracle

	// check if there's already a position for that pair
	if k.HasStoplossPosition(ctx, address, msg.Stoploss.UniswapPairId) {
		return nil, sdkerrors.Wrapf(types.ErrStoplossExists, "address: %s, uniswap pair id %s", address, msg.Stoploss.UniswapPairId)
	}

	k.SetStoplossPosition(ctx, address, *msg.Stoploss)

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeCreateStoploss,
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Address),
			sdk.NewAttribute(types.AttributeKeyUniswapPair, msg.Stoploss.UniswapPairId),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.ModuleName),
		),
	})

	defer func() {
		telemetry.IncrCounterWithLabels(
			[]string{types.ModuleName, "create-stoploss"},
			1,
			[]metrics.Label{
				{Name: "pair-id", Value: msg.Stoploss.UniswapPairId},
			},
		)
	}()

	return &types.MsgStoplossResponse{}, nil
}

// DeleteStoploss for a given uniswap pair
func (k Keeper) DeleteStoploss(c context.Context, msg *types.MsgDeleteStoploss) (*types.MsgDeleteStoplossResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)

	// NOTE: error checked during msg validation
	address, _ := sdk.AccAddressFromBech32(msg.Address)

	// error if the stoploss doesn't exist
	if !k.HasStoplossPosition(ctx, address, msg.UniswapPairId) {
		return nil, sdkerrors.Wrapf(types.ErrStoplossNotFound, "address: %s, uniswap pair id %s", address, msg.UniswapPairId)
	}

	// delete the stoploss position
	k.DeleteStoplossPosition(ctx, address, msg.UniswapPairId)

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeDeleteStoploss,
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Address),
			sdk.NewAttribute(types.AttributeKeyUniswapPair, msg.UniswapPairId),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.ModuleName),
		),
	})

	defer func() {
		telemetry.IncrCounterWithLabels(
			[]string{types.ModuleName, "delete-stoploss"},
			1,
			[]metrics.Label{
				telemetry.NewLabel("pair-id", msg.UniswapPairId),
			},
		)
	}()

	return &types.MsgDeleteStoplossResponse{}, nil
}
