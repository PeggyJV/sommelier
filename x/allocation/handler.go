package allocation

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	sdk "github.com/cosmos/cosmos-sdk/types"

	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	"github.com/peggyjv/sommelier/v3/x/allocation/keeper"
	"github.com/peggyjv/sommelier/v3/x/allocation/types"
)

// NewHandler returns a handler for "oracle" type messages.
func NewHandler(k keeper.Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) (*sdk.Result, error) {
		ctx = ctx.WithEventManager(sdk.NewEventManager())
		switch msg := msg.(type) {
		case *types.MsgAllocationPrecommit:
			res, err := k.AllocationPrecommit(sdk.WrapSDKContext(ctx), msg)
			return sdk.WrapServiceResult(ctx, res, err)
		case *types.MsgAllocationCommit:
			res, err := k.AllocationCommit(sdk.WrapSDKContext(ctx), msg)
			return sdk.WrapServiceResult(ctx, res, err)
		default:
			return nil, sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "unrecognized oracle message type: %T", msg)
		}
	}
}

func NewUpdateManagedCellarsProposalHandler(k keeper.Keeper) govtypes.Handler {
	return func(ctx sdk.Context, content govtypes.Content) error {
		switch c := content.(type) {
		case *types.AddManagedCellarsProposal:
			return keeper.HandleAddManagedCellarsProposal(ctx, k, *c)
		case *types.RemoveManagedCellarsProposal:
			return keeper.HandleRemoveManagedCellarsProposal(ctx, k, *c)

		default:
			return sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "unrecognized allocation proposal content type: %T", c)
		}
	}
}
