package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/peggyjv/sommelier/v4/x/pubsub/types"
)

// HandleAddPublisherProposal is a handler for executing a passed community publisher addition proposal
func HandleAddPublisherProposal(ctx sdk.Context, k Keeper, p types.AddPublisherProposal) error {
	_, found := k.GetPublisher(ctx, p.Domain)
	if found {
		return sdkerrors.Wrapf(types.ErrAlreadyExists, "publisher already exists with domain: %s", p.Domain)
	}

	publisher := types.Publisher{
		Domain:  p.Domain,
		Address: p.Address,
		CaCert:  p.CaCert,
	}

	k.SetPublisher(ctx, publisher)

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

func HandleAddDefaultSubscriptionProposal(ctx sdk.Context, k Keeper, p types.AddDefaultSubscriptionProposal) error {
	defaultSubscription := types.DefaultSubscription{
		SubscriptionId:  p.SubscriptionId,
		PublisherDomain: p.PublisherDomain,
	}

	k.SetDefaultSubscription(ctx, defaultSubscription)

	return nil
}

func HandleRemoveDefaultSubscriptionProposal(ctx sdk.Context, k Keeper, p types.RemoveDefaultSubscriptionProposal) error {
	_, found := k.GetDefaultSubscription(ctx, p.SubscriptionId)
	if !found {
		return sdkerrors.Wrapf(sdkerrors.ErrNotFound, "no default subscription found with id: %s", p.SubscriptionId)
	}

	k.DeleteDefaultSubscription(ctx, p.SubscriptionId)

	return nil
}
