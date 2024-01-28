package cork

import (
	errorsmod "cosmossdk.io/errors"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	govtypesv1beta1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1beta1"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/peggyjv/sommelier/v7/x/cork/keeper"
	types "github.com/peggyjv/sommelier/v7/x/cork/types/v2"
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
			return nil, errorsmod.Wrapf(sdkerrors.ErrUnknownRequest, "unrecognized cork message type: %T", msg)
		}
	}
}

func NewProposalHandler(k keeper.Keeper) govtypesv1beta1.Handler {
	return func(ctx sdk.Context, content govtypesv1beta1.Content) error {
		switch c := content.(type) {
		case *types.AddManagedCellarIDsProposal:
			return keeper.HandleAddManagedCellarsProposal(ctx, k, *c)
		case *types.RemoveManagedCellarIDsProposal:
			return keeper.HandleRemoveManagedCellarsProposal(ctx, k, *c)
		case *types.ScheduledCorkProposal:
			return keeper.HandleScheduledCorkProposal(ctx, k, *c)

		default:
			return errorsmod.Wrapf(sdkerrors.ErrUnknownRequest, "unrecognized cork proposal content type: %T", c)
		}
	}
}
