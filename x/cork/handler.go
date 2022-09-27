package cork

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/peggyjv/sommelier/v4/x/cork/keeper"
	"github.com/peggyjv/sommelier/v4/x/cork/types"
)

// NewHandler returns a handler for "cork" type messages.
func NewHandler(k keeper.Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) (*sdk.Result, error) {
		ctx = ctx.WithEventManager(sdk.NewEventManager())
		switch msg := msg.(type) {
		case *types.MsgScheduleCorkRequest:
			res, err := k.ScheduleCork(sdk.WrapSDKContext(ctx), msg)
			return sdk.WrapServiceResult(ctx, res, err)
		default:
			return nil, sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "unrecognized cork message type: %T", msg)
		}
	}
}

func NewUpdateCellarIDsProposalHandler(k keeper.Keeper) govtypes.Handler {
	return func(ctx sdk.Context, content govtypes.Content) error {
		switch c := content.(type) {
		case *types.AddManagedCellarIDsProposal:
			return keeper.HandleAddManagedCellarsProposal(ctx, k, *c)
		case *types.RemoveManagedCellarIDsProposal:
			return keeper.HandleRemoveManagedCellarsProposal(ctx, k, *c)

		default:
			return sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "unrecognized cork proposal content type: %T", c)
		}
	}
}
