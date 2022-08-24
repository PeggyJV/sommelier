package auction

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/peggyjv/sommelier/v4/x/auction/keeper"
	"github.com/peggyjv/sommelier/v4/x/auction/types"
)

// NewHandler returns a handler for "auction" type messages.
func NewHandler(k keeper.Keeper) sdk.Handler {
	// TODO: fill in
	return nil
}

func NewSetTokenPricesProposalHandler(k keeper.Keeper) govtypes.Handler {
	return func(ctx sdk.Context, content govtypes.Content) error {
		switch c := content.(type) {
		case *types.SetTokenPricesProposal:
			return keeper.HandleSetTokenPricesProposal(ctx, k, *c)

		default:
			return sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "unrecognized proposal content type: %T", c)
		}
	}
}