package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/peggyjv/sommelier/v3/x/pubsub/types"
)

// HandleAddPublisherProposal is a handler for executing a passed community publisher addition proposal
func HandleAddPublisherProposal(ctx sdk.Context, k Keeper, p types.AddPublisherProposal) error {
	// TODO(bolten): I assume that for a proposal to reach this point it has already have ValidateBasic called on it?
	_, found := k.GetPublisher(ctx, p.Domain)
	if found {
		// TODO(bolten): should we just overwrite an existing publisher if governance passes like this, or should we
		// require that the existing publisher is removed first?
		return sdkerrors.Wrapf(types.ErrAlreadyExists, "publisher already exists with domain: %s", p.Domain)
	}

	publisher := types.Publisher{
		Domain:   p.Domain,
		Address:  p.Address,
		ProofUrl: p.ProofUrl,
		CaCert:   p.CaCert,
	}

	k.SetPublisher(ctx, publisher)

	// TODO(bolten): should events be emitted for these proposal handlers?

	return nil
}

// HandleRemovePublisherProposal is a handler for executing a passed community publisher removal proposal
func HandleRemovePublisherProposal(ctx sdk.Context, k Keeper, p types.RemovePublisherProposal) error {
	_, found := k.GetPublisher(ctx, p.Domain)
	if !found {
		return sdkerrors.Wrapf(sdkerrors.ErrNotFound, "no publisher found with domain: %s", p.Domain)
	}

	k.DeletePublisher(ctx, p.Domain)

	return nil
}
