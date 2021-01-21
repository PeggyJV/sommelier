package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/peggyjv/sommelier/x/il/types"
)

var _ types.MsgServer = Keeper{}

// MsgServer is the server API for Msg service.

// CreateStoplossPosition for a given uniswap pair
func (k Keeper) CreateStoplossPosition(c context.Context, msg *types.MsgStoploss) (*types.MsgStoplossResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)

	address, err := sdk.AccAddressFromBech32(msg.Address)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, err.Error())
	}

	// Set the delegation
	if k.HasStoplossPosition(ctx, address, msg.Stoploss.UniswapPairId) {
		return nil, sdkerrors.Wrapf(types.ErrStoplossExists, "address: %s, uniswap pair id %s", address, msg.Stoploss.UniswapPairId)
	}

	// Set the delegation
	k.SetStoplossPosition(ctx, address, msg.Stoploss)

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

	return &types.MsgDelegateFeedConsentResponse{}, nil
}
