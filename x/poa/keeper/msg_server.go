package keeper

import (
	"context"

	sdkerrors "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/peggyjv/sommelier/v9/x/poa/types"
)

type msgServer struct {
	Keeper
}

// NewMsgServerImpl returns a types.MsgServer implementation backed by the
// PoA Keeper. Both messages are gov-only.
func NewMsgServerImpl(k Keeper) types.MsgServer {
	return &msgServer{Keeper: k}
}

var _ types.MsgServer = msgServer{}

func (s msgServer) UpdateAuthoritySet(goCtx context.Context, msg *types.MsgUpdateAuthoritySet) (*types.MsgUpdateAuthoritySetResponse, error) {
	if msg.Authority != s.authority {
		return nil, sdkerrors.Wrapf(types.ErrUnauthorized, "expected %s, got %s", s.authority, msg.Authority)
	}
	if err := msg.ValidateBasic(); err != nil {
		return nil, err
	}
	ctx := sdk.UnwrapSDKContext(goCtx)
	addrs := make([]sdk.ValAddress, len(msg.Validators))
	for i, v := range msg.Validators {
		addr, err := sdk.ValAddressFromBech32(v)
		if err != nil {
			return nil, err
		}
		addrs[i] = addr
	}
	s.SetAuthoritySet(ctx, addrs)

	ctx.EventManager().EmitEvent(sdk.NewEvent(
		types.EventTypeAuthoritySetUpdated,
		sdk.NewAttribute("size", sdk.NewInt(int64(len(addrs))).String()),
	))
	return &types.MsgUpdateAuthoritySetResponse{}, nil
}

func (s msgServer) UpdateParams(goCtx context.Context, msg *types.MsgUpdateParams) (*types.MsgUpdateParamsResponse, error) {
	if msg.Authority != s.authority {
		return nil, sdkerrors.Wrapf(types.ErrUnauthorized, "expected %s, got %s", s.authority, msg.Authority)
	}
	if err := msg.ValidateBasic(); err != nil {
		return nil, err
	}
	ctx := sdk.UnwrapSDKContext(goCtx)
	s.SetParams(ctx, msg.Params)

	ctx.EventManager().EmitEvent(sdk.NewEvent(
		types.EventTypeParamsUpdated,
		sdk.NewAttribute("floor_fraction", msg.Params.FloorFraction.String()),
	))
	return &types.MsgUpdateParamsResponse{}, nil
}
