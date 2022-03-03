package keeper

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/peggyjv/sommelier/v3/x/pubsub/types"
)

// InitGenesis initializes the module's state from a provided genesis
// state.
func InitGenesis(ctx sdk.Context, k Keeper, gs types.GenesisState) {
	k.SetParams(ctx, gs.Params)

	for _, publisher := range gs.Publishers {
		k.SetPublisher(ctx, *publisher)
	}

	for _, subscriber := range gs.Subscribers {
		addr, err := sdk.AccAddressFromBech32(subscriber.Address)
		if err != nil {
			panic(fmt.Errorf("invalid subscriber address in genesis state: %s", err.Error()))
		}
		k.SetSubscriber(ctx, addr, *subscriber)
	}

	for _, publisherIntent := range gs.PublisherIntents {
		k.SetPublisherIntent(ctx, *publisherIntent)
	}

	for _, subscriberIntent := range gs.SubscriberIntents {
		addr, err := sdk.AccAddressFromBech32(subscriberIntent.SubscriberAddress)
		if err != nil {
			panic(fmt.Errorf("invalid subscriber address in genesis state: %s", err.Error()))
		}
		k.SetSubscriberIntent(ctx, addr, *subscriberIntent)
	}
}

// ExportGenesis returns the module's exported genesis.
func ExportGenesis(ctx sdk.Context, k Keeper) types.GenesisState {
	return types.GenesisState{
		Params:            k.GetParams(ctx),
		Publishers:        k.GetPublishers(ctx),
		Subscribers:       k.GetSubscribers(ctx),
		PublisherIntents:  k.GetPublisherIntents(ctx),
		SubscriberIntents: k.GetSubscriberIntents(ctx),
	}
}
